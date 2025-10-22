package collector

// Payload payload example
type Payload struct {
	GPRS         string `json:"gprs"`
	Scare        string `json:"scare"`
	IMEI         string `json:"imei"`
	UnitID       string `json:"unitID"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	Attending    string `json:"attending"`
	ConfirmPanic string `json:"confirmPanic"`
	IP           string `json:"ip"`
	Request      string `json:"request"`
}

type AlarmPayload struct {
	IMEI      string
	Latitude  string
	Longitude string
	AlarmType string
	Waiting   string
	Attending string
}
