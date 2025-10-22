package traffic

import (
	"github.com/jmontesinos91/collector/domains/pagination"
	"net/http"
	"net/url"
	"testing"
	"time"

	otraffic "github.com/jmontesinos91/collector/internal/repositories/traffic"
	"github.com/stretchr/testify/assert"
)

func TestToResponse(t *testing.T) {

	data := map[string]interface{}{
		"id":   123,
		"name": "test",
	}
	status := "success"
	message := "operation completed successfully"

	response := ToResponse(data, status, message)

	assert.Equal(t, status, response.Status, "Expected status to be 'success'")
	assert.Equal(t, message, response.Message, "Expected message to match")
	assert.Equal(t, data, response.Data, "Expected data to match")
}

func TestToPaginatedResponse(t *testing.T) {

	data := []string{"item1", "item2", "item3"}
	currentPage := 1
	pages := 10
	total := 30

	response := ToPaginatedResponse(data, currentPage, pages, total)

	assert.Equal(t, data, response.Data, "Expected data to match")
	assert.Equal(t, currentPage, response.CurrentPage, "Expected currentPage to be 1")
	assert.Equal(t, pages, response.Pages, "Expected pages to be 10")
	assert.Equal(t, total, response.Total, "Expected total to be 30")
}

func TestParseFilterRequest(t *testing.T) {
	isAlarmTrue := true
	counter := 1
	tests := []struct {
		name        string
		queryParams map[string]string
		expected    *FilterRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Happy path valid parameters",
			queryParams: map[string]string{
				"q":            "test-query",
				"id":           "12345",
				"request":      "test-request",
				"imei":         "test-imei",
				"ip":           "192.168.0.1",
				"alarm":        "1",
				"counter":      "1",
				"createdAtMin": "2024-09-16T15:04",
				"createdAtMax": "2024-09-17T15:04",
				"updatedAtMin": "2024-09-18T15:04",
				"updatedAtMax": "2024-09-19T15:04",
				"sortBy":       "ip",
				"sortDesc":     "true",
				"size":         "20",
				"page":         "2",
				"action":       "list",
			},
			expected: &FilterRequest{
				QParam:  "test-query",
				ID:      "12345",
				Request: "test-request",
				IMEI:    "test-imei",
				Ip:      "192.168.0.1",
				Counter: &counter,
				IsAlarm: &isAlarmTrue,
				Action:  "list",
				Filter: pagination.Filter{
					SortBy:   "ip",
					SortDesc: true,
					Size:     20,
					Page:     2,
				},
			},
			expectError: false,
		},
		{
			name: "Invalid size parameter",
			queryParams: map[string]string{
				"size": "invalid-size",
			},
			expected:    nil,
			expectError: true,
			errorMsg:    "Invalid size parameter",
		},
		{
			name: "Invalid sortDesc parameter",
			queryParams: map[string]string{
				"sortDesc": "invalid-boolean",
			},
			expected:    nil,
			expectError: true,
			errorMsg:    "Invalid sortDesc parameter",
		},
		{
			name: "Invalid counter parameter",
			queryParams: map[string]string{
				"counter": "invalid-integer",
			},
			expected:    nil,
			expectError: true,
			errorMsg:    "Invalid counter parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			query := url.Values{}
			for key, value := range tt.queryParams {
				query.Set(key, value)
			}
			req := &http.Request{
				URL: &url.URL{RawQuery: query.Encode()},
			}

			fr, err := ParseFilterRequest(req)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, fr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, fr)
			}
		})
	}
}

func TestToTrafficSlice(t *testing.T) {
	// Test Data
	trafficModels := []otraffic.Model{
		{
			ID:        "1",
			Request:   "request1",
			IMEI:      "imei1",
			Ip:        "192.168.0.1",
			IsAlarm:   true,
			Counter:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "2",
			Request:   "request2",
			IMEI:      "imei2",
			Ip:        "192.168.0.2",
			IsAlarm:   true,
			Counter:   2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	result := ToTrafficSlice(trafficModels)
	assert.Len(t, result, len(trafficModels), "Expected length to match")
	for i, model := range trafficModels {
		assert.Equal(t, model.ID, result[i].ID, "Expected ID to match for element %d", i)
		assert.Equal(t, model.Request, result[i].Request, "Expected Request to match for element %d", i)
		assert.Equal(t, model.IMEI, result[i].IMEI, "Expected IMEI to match for element %d", i)
		assert.Equal(t, model.Ip, result[i].Ip, "Expected IP to match for element %d", i)
		assert.Equal(t, model.IsAlarm, result[i].IsAlarm, "Expected Alarm to match for element %d", i)
		assert.Equal(t, model.Counter, result[i].Counter, "Expected Count to match for element %d", i)
		assert.Equal(t, model.CreatedAt, result[i].CreatedAt, "Expected CreatedAt to match for element %d", i)
		assert.Equal(t, model.UpdatedAt, result[i].UpdatedAt, "Expected UpdatedAt to match for element %d", i)
	}
}

func TestToTraffic(t *testing.T) {
	model := otraffic.Model{
		ID:        "1",
		Request:   "request1",
		IMEI:      "imei1",
		Ip:        "192.168.0.1",
		IsAlarm:   true,
		Counter:   2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := ToTraffic(model)
	assert.Equal(t, model.ID, result.ID, "Expected ID to match")
	assert.Equal(t, model.Request, result.Request, "Expected Request to match")
	assert.Equal(t, model.IMEI, result.IMEI, "Expected IMEI to match")
	assert.Equal(t, model.Ip, result.Ip, "Expected IP to match")
	assert.Equal(t, model.IsAlarm, result.IsAlarm, "Expected Alarm to match")
	assert.Equal(t, model.Counter, result.Counter, "Expected Counter to match")
	assert.Equal(t, model.CreatedAt, result.CreatedAt, "Expected CreatedAt to match")
	assert.Equal(t, model.UpdatedAt, result.UpdatedAt, "Expected UpdatedAt to match")
}

func TestToMetadata(t *testing.T) {
	isAlarmTrue := true
	filterRequest := &FilterRequest{
		QParam:  "query",
		ID:      "1",
		Request: "request1",
		IMEI:    "imei1",
		Ip:      "192.168.0.1",
		IsAlarm: &isAlarmTrue,
		Action:  "list",
	}

	result := ToMetadata(filterRequest)

	assert.Equal(t, filterRequest.QParam, result.Qparam, "Expected QParam to match")
	assert.Equal(t, filterRequest.ID, result.ID, "Expected ID to match")
	assert.Equal(t, filterRequest.Request, result.Request, "Expected Request to match")
	assert.Equal(t, filterRequest.IMEI, result.IMEI, "Expected IMEI to match")
	assert.Equal(t, filterRequest.Ip, result.Ip, "Expected IP to match")
	assert.Equal(t, filterRequest.IsAlarm, result.IsAlarm, "Expected Alarm to match")
}
