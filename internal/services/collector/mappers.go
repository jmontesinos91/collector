package collector

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmontesinos91/collector/internal/repositories/traffic"
	"github.com/jmontesinos91/terrors"
)

// ParsePayload Build the model expected for repository
func (p *Payload) ParsePayload(r *http.Request) error {
	query := r.URL.Query()
	collect := strings.ReplaceAll(query.Get("router"), " ", "")

	if collect == "" {
		collect = strings.ReplaceAll(chi.URLParam(r, "str"), " ", "")
	}

	collectString := strings.Split(collect, ",")

	if len(collectString) < 12 {
		return terrors.New(terrors.ErrBadRequest, "Invalid Request String", nil)
	}

	p.Request = collect

	if collectString[2] != "" {
		p.IP = collectString[2]
	} else if ip := strings.TrimSuffix(r.Header.Get("Referer"), "/"); ip != "" {
		p.IP = ip
	} else {
		p.IP = r.RemoteAddr
	}

	p.GPRS = collectString[0]
	if collectString[3] != "" {
		p.IMEI = collectString[3]
	} else if collectString[4] != "" {
		p.UnitID = collectString[4]
	} else {
		return terrors.New(terrors.ErrBadRequest, "Invalid Request String", nil)
	}

	// Validate Latitude param
	if collectString[6] != "" {
		p.Latitude = collectString[6]
	} else {
		p.Latitude = "0"
	}

	// Validate Longitude param
	if collectString[7] != "" {
		p.Longitude = collectString[7]
	} else {
		p.Longitude = "0"
	}

	// Validate the attending param
	if len(collectString) >= 13 {
		p.Attending = collectString[12]
	} else {
		p.Attending = "0"
	}

	// Collect Confirm Panic param
	p.ConfirmPanic = collectString[11]

	// Collect scare param
	p.Scare = strings.ToUpper(p.GPRS[len(p.GPRS)-1:])

	return nil
}

func (p *Payload) ParseAlarmPayload(alarmType, waiting string) AlarmPayload {
	return AlarmPayload{
		IMEI:      p.IMEI,
		Latitude:  p.Latitude,
		Longitude: p.Longitude,
		AlarmType: alarmType,
		Attending: p.Attending,
		Waiting:   waiting,
	}
}

func (p *Payload) ToModel(isAlarm bool) traffic.Model {
	return traffic.Model{
		ID:         uuid.NewString(),
		Request:    p.Request,
		IMEI:       p.IMEI,
		Ip:         p.IP,
		IsAlarm:    isAlarm,
		IsNotified: false,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}
