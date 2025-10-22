package router

// Request Holds the response for a created payout
type Request struct {
	IMEI      string `json:"imei"`
	UnitID    string `json:"unitID"`
	AlarmType string `json:"input"`
}

// Response Holds the response for a created payout
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
