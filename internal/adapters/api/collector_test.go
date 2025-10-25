package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	scollector "github.com/jmontesinos91/collector/internal/services/collector"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/terrors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type mockCollectorSvc struct {
	called bool
	err    error
}

func (m *mockCollectorSvc) Collector(ctx context.Context, payload *scollector.Payload) error {
	m.called = true
	return m.err
}

func TestHandleCollector_TableDriven(t *testing.T) {
	cases := []struct {
		name         string
		body         string
		contentType  string
		queryParam   string
		statusCode   int
		svcErr       error
		expectCalled bool
	}{
		{
			name:         "ParseErrorDoesNotCallService",
			body:         `{invalid json`,
			contentType:  "application/json",
			svcErr:       nil,
			statusCode:   http.StatusBadRequest,
			expectCalled: false,
		},
		{
			name:         "ServiceCalledOnSuccess",
			body:         `{}`,
			queryParam:   "P,12,12,861585041440544,12,12,123456789,123456789,00,00,00,1",
			contentType:  "application/json",
			svcErr:       nil,
			statusCode:   http.StatusOK,
			expectCalled: true,
		},
		{
			name:         "ServiceErrorReturnsNon200",
			body:         `{}`,
			contentType:  "application/json",
			svcErr:       terrors.New(terrors.ErrBadRequest, terrors.MsgBadRequest, map[string]string{}),
			statusCode:   http.StatusBadRequest,
			expectCalled: false,
		},
		{
			name:         "Invalid RouterString",
			body:         `{}`,
			queryParam:   "P,12,12,861585041440544,12",
			contentType:  "application/json",
			svcErr:       terrors.New(terrors.ErrBadRequest, "Invalid Request String", nil),
			statusCode:   http.StatusBadRequest,
			expectCalled: false,
		},
	}

	for _, tc := range cases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockCollectorSvc{err: tc.svcErr}
			sc := &CollectorController{
				log:      logger.NewContextLogger("COLLECTOR-UNIT-TEST", "debug", logger.TextFormat),
				validate: nil,
				counterMetric: prometheus.NewCounter(prometheus.CounterOpts{
					Name:      "collector_reqs_total_test_" + tc.name,
					Namespace: "collector",
					Subsystem: "api",
					Help:      "test metric",
				}),
				collectorSv: mockSvc,
				stsClient:   nil,
			}

			body := bytes.NewBufferString(tc.body)
			q := url.Values{}
			q.Set("router", tc.queryParam)
			req := httptest.NewRequest(http.MethodPost, "/v2/routers/?"+q.Encode(), body)
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rec := httptest.NewRecorder()

			sc.handleCollector(rec, req)

			if rec.Code != tc.statusCode {
				t.Fatalf("[%s] expected status %d, got %d", tc.name, tc.statusCode, rec.Code)
			}

			if mockSvc.called != tc.expectCalled {
				t.Fatalf("[%s] expected called=%v, got %v", tc.name, tc.expectCalled, mockSvc.called)
			}

			if testutil.ToFloat64(sc.counterMetric) != 1 {
				t.Fatalf("[%s] metric was not incremented, want 1 got %f", tc.name, testutil.ToFloat64(sc.counterMetric))
			}

		})
	}
}
