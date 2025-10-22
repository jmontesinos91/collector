package traffic

import (
	"github.com/jmontesinos91/collector/domains/pagination"
	"time"
)

// Response struct to response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Alarm struct {
	IMEI      string
	Latitude  string
	Longitude string
	AlarmType string
	Waiting   string
	Attending string
}

// FilterRequest holds the http request params
type FilterRequest struct {
	QParam  string            `json:"q,omitempty"`
	ID      string            `json:"id,omitempty"`
	Request string            `json:"request,omitempty"`
	IMEI    string            `json:"imei,omitempty"`
	Ip      string            `json:"ip,omitempty"`
	IsAlarm *bool             `json:"alarm,omitempty"`
	Counter *int              `json:"counter,omitempty"`
	Action  string            `json:"action,omitempty"`
	Filter  pagination.Filter `json:"filter,omitempty"`
}

// Traffic item
type Traffic struct {
	ID        string    `json:"id"`
	Request   string    `json:"request"`
	IMEI      string    `json:"imei"`
	Ip        string    `json:"ip"`
	IsAlarm   bool      `json:"alarm"`
	Counter   int       `json:"counter"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
