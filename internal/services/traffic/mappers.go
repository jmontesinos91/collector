package traffic

import (
	"github.com/jmontesinos91/collector/internal/repositories/traffic"
	"github.com/jmontesinos91/terrors"
	"net/http"
	"strconv"
)

// ParseFilterRequest builds a single filter object given http params
func ParseFilterRequest(r *http.Request) (*FilterRequest, error) {
	fr := FilterRequest{}
	query := r.URL.Query()

	if QParam := query.Get("q"); QParam != "" {
		fr.QParam = QParam
	}

	if id := query.Get("id"); id != "" {
		fr.ID = id
	}

	if counterStr := query.Get("counter"); counterStr != "" {
		if counter, err := strconv.Atoi(query.Get("counter")); err != nil {
			return nil, terrors.New(terrors.ErrBadRequest, "Invalid counter parameter", map[string]string{})
		} else {
			fr.Counter = &counter
		}
	}

	if request := query.Get("request"); request != "" {
		fr.Request = request
	}

	if imei := query.Get("imei"); imei != "" {
		fr.IMEI = imei
	}

	if ip := query.Get("ip"); ip != "" {
		fr.Ip = ip
	}

	if alarm := query.Get("alarm"); alarm != "" {
		isAlarm, err := strconv.ParseBool(query.Get("alarm"))
		if err != nil {
			return nil, terrors.New(terrors.ErrBadRequest, "Invalid page parameter", map[string]string{})
		}
		fr.IsAlarm = &isAlarm
	}

	if action := query.Get("action"); action != "" {
		switch action {
		case "list":
			fr.Action = "list"
		case "export":
			fr.Action = "export"
		default:
			fr.Action = "list"
		}
	}

	return &fr, nil
}

// ToResponse builds a response object given the argument values
func ToResponse(data interface{}, status string, message string) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// ToTrafficSlice converts a traffic model slice into a serializable slice
func ToTrafficSlice(trafficModels []traffic.Model) []Traffic {
	var traffics []Traffic
	for _, model := range trafficModels {
		traffics = append(traffics, ToTraffic(model))
	}

	return traffics
}

// ToTraffic converts a model to a Traffic struct to be serialized
func ToTraffic(model traffic.Model) Traffic {
	return Traffic{
		ID:        model.ID,
		Request:   model.Request,
		IMEI:      model.IMEI,
		Ip:        model.Ip,
		IsAlarm:   model.IsAlarm,
		Counter:   model.Counter,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// ToMetadata maps the properties of the service filter into repo filter
func ToMetadata(filterRequest *FilterRequest) *traffic.Metadata {

	return &traffic.Metadata{
		Qparam:  filterRequest.QParam,
		ID:      filterRequest.ID,
		Request: filterRequest.Request,
		IMEI:    filterRequest.IMEI,
		Ip:      filterRequest.Ip,
		IsAlarm: filterRequest.IsAlarm,
		Counter: filterRequest.Counter,
	}
}
