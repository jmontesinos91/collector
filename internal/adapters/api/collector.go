package api

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jmontesinos91/collector/internal/services/collector"
	scollector "github.com/jmontesinos91/collector/internal/services/collector"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/osecurity/sts"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// CollectorController controller struct
type CollectorController struct {
	log           *logger.ContextLogger
	validate      *validator.Validate
	collectorSv   collector.IService
	stsClient     sts.ISTSClient
	counterMetric prometheus.Counter
}

// NewCollectorController Constructor
func NewCollectorController(server *HTTPServer, validator *validator.Validate, ss collector.IService, sts sts.ISTSClient) *CollectorController {
	sc := &CollectorController{
		log:         server.Logger,
		validate:    validator,
		collectorSv: ss,
		stsClient:   sts,
		counterMetric: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "collector_reqs_total",
			Namespace:   "collector",
			Subsystem:   "api",
			ConstLabels: map[string]string{},
			Help:        "The total number of requests to routers endpoints",
		}),
	}

	// Endpoints without secure
	// Deprecated: We will remove this endpoint for new usages
	server.Router.Get("/v2/routers/", sc.handleCollector)

	return sc
}

func (sc *CollectorController) handleReceive(w http.ResponseWriter, r *http.Request) {
	sc.log.Log(logrus.InfoLevel, "handleReceive", "Incoming request to handleReceive")

	// Create a context with a deadline of 2 seconds.
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	payload := &scollector.Payload{}
	err := payload.ParsePayload(r)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	err = sc.collectorSv.Collector(ctx, payload)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	RenderJSON(r.Context(), w, http.StatusOK, nil)
}

func (sc *CollectorController) handleCollector(w http.ResponseWriter, r *http.Request) {
	sc.log.Log(logrus.InfoLevel, "handleCollector", "Incoming request to handleCollector")

	// Increment metric
	sc.counterMetric.Inc()

	// Create a context with a deadline of 2 seconds.
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	payload := &scollector.Payload{}
	err := payload.ParsePayload(r)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}

	err = sc.collectorSv.Collector(ctx, payload)
	if err != nil {
		RenderError(r.Context(), w, err)
	}

	RenderJSON(r.Context(), w, http.StatusOK, nil)
}
