package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	tservice "github.com/jmontesinos91/collector/internal/services/traffic"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/osecurity/sts"
	"github.com/jmontesinos91/terrors"
	"github.com/sirupsen/logrus"
)

type TrafficController struct {
	log        *logger.ContextLogger
	validate   *validator.Validate
	trafficSvc tservice.IService
	stsClient  sts.ISTSClient
}

func NewTrafficController(server *HTTPServer, validator *validator.Validate, ts tservice.IService, sts sts.ISTSClient) *TrafficController {

	sc := &TrafficController{
		log:        server.Logger,
		validate:   validator,
		trafficSvc: ts,
		stsClient:  sts,
	}

	// Endpoint secure
	server.Router.Group(func(r chi.Router) {
		r.Use(JwtVerifyMiddleware(server.Logger, sts))
		r.Get("/v1/traffic", sc.handleRetrieve)
		r.Post("/v1/traffic/{id}", sc.handleDelete)
		r.Post("/v1/traffic/counter/reset/{id}", sc.handleCounterReset)
	})

	return sc
}

func (tc *TrafficController) handleRetrieve(w http.ResponseWriter, r *http.Request) {
	tc.log.Log(logrus.InfoLevel, "handleRetrieve", "Incoming request to handleRetrieve")

	filters, err := tservice.ParseFilterRequest(r)
	if err != nil {
		tc.log.Error(logrus.ErrorLevel, "handleRetrieve", "Invalid request parameters", err)
		RenderError(r.Context(), w, err)
		return
	}

	data, err := tc.trafficSvc.HandleRetrieve(r.Context(), filters)
	if err != nil {
		tc.log.Error(logrus.ErrorLevel, "handleRetrieve", "Failed to retrieve traffics", err)
		RenderError(r.Context(), w, terrors.InternalService("internal_error", "Failed to retrieve traffics", map[string]string{}))
		return
	}

	RenderJSON(r.Context(), w, http.StatusOK, data)
}

func (tc *TrafficController) handleDelete(w http.ResponseWriter, r *http.Request) {
	tc.log.Log(logrus.InfoLevel, "handleDelete", "Incoming request to handleDelete")

	trafficID := chi.URLParam(r, "id")

	err := tc.trafficSvc.HandleDelete(r.Context(), trafficID)
	if err != nil {
		tc.log.Error(logrus.ErrorLevel, "handleDelete", "Failed to delete traffic resource", err)
		RenderError(r.Context(), w, terrors.InternalService("internal_error", "Failed to delete traffic resource", map[string]string{}))
		return
	}

	RenderJSON(r.Context(), w, http.StatusAccepted, nil)
}

func (tc *TrafficController) handleCounterReset(w http.ResponseWriter, r *http.Request) {
	tc.log.Log(logrus.InfoLevel, "handleCounterReset", "Incoming request to handleRetrieve")

	trafficID := chi.URLParam(r, "id")

	err := tc.trafficSvc.HandleResetCounter(r.Context(), trafficID)
	if err != nil {
		tc.log.Error(logrus.ErrorLevel, "handleCounterReset", "Failed to reset traffic counter", err)
		RenderError(r.Context(), w, terrors.InternalService("internal_error", "Failed to reset traffic counter", map[string]string{}))
		return
	}

	RenderJSON(r.Context(), w, http.StatusAccepted, nil)
}
