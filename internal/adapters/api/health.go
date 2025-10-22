package api

import (
	"net/http"

	"github.com/jmontesinos91/ologs/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HealthController Handles all health related routes
type HealthController struct {
	log *logger.ContextLogger
}

// NewHealthController Creates a new instance
func NewHealthController(server *HTTPServer) *HealthController {
	hc := &HealthController{
		log: server.Logger,
	}

	// Loads routes
	server.Router.Get("/health/live", hc.handleLivenessCheck)
	server.Router.Get("/health/ready", hc.handleReadinessCheck)
	server.Router.Get("/metrics", promhttp.Handler().ServeHTTP)

	return hc
}

func (hc *HealthController) handleLivenessCheck(w http.ResponseWriter, r *http.Request) {
	RenderJSON(r.Context(), w, http.StatusOK, map[string]string{"status": "ok"})
}

func (hc *HealthController) handleReadinessCheck(w http.ResponseWriter, r *http.Request) {
	RenderJSON(r.Context(), w, http.StatusOK, map[string]string{"status": "ok"})
}
