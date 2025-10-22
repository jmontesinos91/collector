package collector

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePayload(t *testing.T) {

	tests := []struct {
		name        string
		queryParams map[string]string
		path        string
		headers     http.Header
		remoteAddr  string
		expected    *Payload
		expectError bool
		errorMsg    string
	}{
		{
			name: "Happy path router parameters",
			queryParams: map[string]string{
				"router": "P,12,192.168.100.1,861585041440544,12,12,123456789,123456789,00,00,00,1",
			},
			expected: &Payload{
				Request:      "P,12,192.168.100.1,861585041440544,12,12,123456789,123456789,00,00,00,1",
				IP:           "192.168.100.1",
				IMEI:         "861585041440544",
				Latitude:     "123456789",
				Longitude:    "123456789",
				Attending:    "0",
				ConfirmPanic: "1",
				Scare:        "P",
				GPRS:         "P",
			},
			expectError: false,
		},
		{
			name: "Biggest length parameter",
			queryParams: map[string]string{
				"router": "P,12,192.168.100.1,861585041440544,12,12,123456789,123456789,00,00,00,1,3",
			},
			expected: &Payload{
				Request:      "P,12,192.168.100.1,861585041440544,12,12,123456789,123456789,00,00,00,1,3",
				IP:           "192.168.100.1",
				IMEI:         "861585041440544",
				Latitude:     "123456789",
				Longitude:    "123456789",
				Attending:    "3",
				ConfirmPanic: "1",
				Scare:        "P",
				GPRS:         "P",
			},
			expectError: false,
		},
		{
			name: "Wrong length parameter",
			queryParams: map[string]string{
				"router": "P,12,192.168.100.1,00,00,1",
			},
			expected:    &Payload{},
			expectError: true,
			errorMsg:    "bad_request: Invalid Request String",
		},
		{
			name: "EmptyIP in Request",
			queryParams: map[string]string{
				"router": "P,12,,861585041440544,12,12,123456789,123456789,00,00,00,1",
			},
			expected: &Payload{
				Request:      "P,12,,861585041440544,12,12,123456789,123456789,00,00,00,1",
				IP:           "192.168.100.1",
				IMEI:         "861585041440544",
				Latitude:     "123456789",
				Longitude:    "123456789",
				Attending:    "0",
				ConfirmPanic: "1",
				Scare:        "P",
				GPRS:         "P",
			},
			headers: http.Header{
				"Referer": []string{"192.168.100.1"},
			},
			expectError: false,
		},
		{
			name: "Remote address",
			queryParams: map[string]string{
				"router": "P,12,,861585041440544,12,12,123456789,123456789,00,00,00,1",
			},
			expected: &Payload{
				Request:      "P,12,,861585041440544,12,12,123456789,123456789,00,00,00,1",
				IP:           "192.168.100.2",
				IMEI:         "861585041440544",
				Latitude:     "123456789",
				Longitude:    "123456789",
				Attending:    "0",
				ConfirmPanic: "1",
				Scare:        "P",
				GPRS:         "P",
			},
			remoteAddr:  "192.168.100.2",
			expectError: false,
		},
		{
			name: "Latitude empty",
			queryParams: map[string]string{
				"router": "P,12,192.168.100.1,861585041440544,12,12,,123456789,00,00,00,1",
			},
			expected: &Payload{
				Request:      "P,12,192.168.100.1,861585041440544,12,12,,123456789,00,00,00,1",
				IP:           "192.168.100.1",
				IMEI:         "861585041440544",
				Latitude:     "0",
				Longitude:    "123456789",
				Attending:    "0",
				ConfirmPanic: "1",
				Scare:        "P",
				GPRS:         "P",
			},
			expectError: false,
		},
		{
			name: "Longitude empty",
			queryParams: map[string]string{
				"router": "P,12,192.168.100.1,861585041440544,12,12,123456789,,00,00,00,1",
			},
			expected: &Payload{
				Request:      "P,12,192.168.100.1,861585041440544,12,12,123456789,,00,00,00,1",
				IP:           "192.168.100.1",
				IMEI:         "861585041440544",
				Latitude:     "123456789",
				Longitude:    "0",
				Attending:    "0",
				ConfirmPanic: "1",
				Scare:        "P",
				GPRS:         "P",
			},
			expectError: false,
		},
		{
			name: "Empty IMEI",
			queryParams: map[string]string{
				"router": "P,12,192.168.100.1,,53438,12,123456789,123456789,00,00,00,1",
			},
			expected: &Payload{
				Request:      "P,12,192.168.100.1,,53438,12,123456789,123456789,00,00,00,1",
				IP:           "192.168.100.1",
				IMEI:         "",
				UnitID:       "53438",
				Latitude:     "123456789",
				Longitude:    "123456789",
				Attending:    "0",
				ConfirmPanic: "1",
				Scare:        "P",
				GPRS:         "P",
			},
			expectError: false,
		},
		{
			name: "Empty IMEI and UnitID",
			queryParams: map[string]string{
				"router": "P,12,192.168.100.1,,,12,123456789,123456789,00,00,00,1",
			},
			expected:    &Payload{},
			expectError: true,
			errorMsg:    "bad_request: Invalid Request String",
		},
		{
			name:        "Empty Value",
			expected:    &Payload{},
			expectError: true,
			errorMsg:    "bad_request: Invalid Request String",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := &Payload{}
			query := url.Values{}
			for key, value := range tt.queryParams {
				query.Set(key, value)
			}

			req := &http.Request{
				URL: &url.URL{
					Path:     tt.path,
					RawQuery: query.Encode(),
				},
				Header: tt.headers,
			}

			req.RemoteAddr = tt.remoteAddr

			err := payload.ParsePayload(req)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, payload)
			}
		})
	}
}

func TestParseAlarmPayload(t *testing.T) {
	type args struct { //nolint:wsl
		payload   *Payload
		alarmType string
		waiting   string
		expected  AlarmPayload
	}
	tests := []struct {
		name string
		args
	}{
		{
			name: "Happy path",
			args: args{
				payload: &Payload{
					IMEI:      "861585041440544",
					Latitude:  "123456789",
					Longitude: "123456789",
					Attending: "1",
				},
				alarmType: "3",
				waiting:   "1",
				expected: AlarmPayload{
					IMEI:      "861585041440544",
					Latitude:  "123456789",
					Longitude: "123456789",
					Attending: "1",
					AlarmType: "3",
					Waiting:   "1",
				},
			},
		},
		{
			name: "Happy path with Out IMEI",
			args: args{
				payload: &Payload{
					IMEI:      "",
					Latitude:  "123456789",
					Longitude: "123456789",
					Attending: "1",
				},
				alarmType: "3",
				waiting:   "1",
				expected: AlarmPayload{
					IMEI:      "",
					Latitude:  "123456789",
					Longitude: "123456789",
					Attending: "1",
					AlarmType: "3",
					Waiting:   "1",
				},
			},
		},
		{
			name: "Happy path with out data",
			args: args{
				payload:  &Payload{},
				expected: AlarmPayload{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			response := tt.payload.ParseAlarmPayload(tt.alarmType, tt.waiting)
			assert.Equal(t, tt.expected, response)
		})
	}
}

func TestToModel(t *testing.T) {
	payload := &Payload{
		Request:      "P,12,192.168.100.1,861585041440544,12,12,123456789,123456789,00,00,00,1",
		IP:           "192.168.100.1",
		IMEI:         "861585041440544",
		Latitude:     "123456789",
		Longitude:    "123456789",
		Attending:    "0",
		ConfirmPanic: "1",
		Scare:        "P",
		GPRS:         "P",
	}

	model := payload.ToModel(true)
	assert.Equal(t, payload.Request, model.Request, "Expected Request to match")
	assert.Equal(t, payload.IMEI, model.IMEI, "Expected IMEI to match")
	assert.Equal(t, payload.IP, model.Ip, "Expected IP to match")
	assert.Equal(t, true, model.IsAlarm, "Expected Alarm to match")
}
