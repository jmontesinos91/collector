package traffic

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	otraffic "github.com/jmontesinos91/collector/internal/repositories/traffic"
	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"
	"github.com/jmontesinos91/osecurity/sts"
	"github.com/jmontesinos91/terrors"
	"github.com/sirupsen/logrus"
)

type DefaultService struct {
	log         *logger.ContextLogger
	trafficRepo otraffic.IRepository
}

func NewDefaultService(l *logger.ContextLogger, tr otraffic.IRepository) *DefaultService {
	return &DefaultService{
		log:         l,
		trafficRepo: tr,
	}
}

func (s *DefaultService) HandleRetrieve(ctx context.Context, filter *FilterRequest) ([]Traffic, error) {
	requestID := ctx.Value(middleware.RequestIDKey).(string)
	claims := ctx.Value(&sts.Claim).(sts.Claims)

	repoFilter := ToMetadata(filter)
	trafficModels, err := s.trafficRepo.Retrieve(ctx, repoFilter)
	if err != nil {
		s.log.WithContext(logrus.ErrorLevel,
			"HandleRetrieve",
			"Failed to retrieve traffics",
			logger.Context{
				tracekey.TrackingID: requestID,
				tracekey.UserID:     claims.UserID,
				tracekey.Role:       claims.Role,
			},
			err)
		return []Traffic{}, err
	}

	return ToTrafficSlice(trafficModels), nil
}

func (s *DefaultService) HandleDelete(ctx context.Context, trafficID string) error {
	requestID := ctx.Value(middleware.RequestIDKey).(string)
	claims := ctx.Value(&sts.Claim).(sts.Claims)
	if trafficID == "" {
		return terrors.New(terrors.ErrBadRequest, "Invalid trafficID", map[string]string{})
	}

	err := s.trafficRepo.DeleteByID(ctx, trafficID)
	if err != nil {
		s.log.WithContext(logrus.ErrorLevel,
			"HandleResetCounter",
			"Failed to delete traffic resource",
			logger.Context{
				tracekey.TrackingID: requestID,
				tracekey.UserID:     claims.UserID,
				tracekey.Role:       claims.Role,
			},
			err)
		return err
	}

	return nil
}

func (s *DefaultService) HandleResetCounter(ctx context.Context, trafficID string) error {
	requestID := ctx.Value(middleware.RequestIDKey).(string)
	claims := ctx.Value(&sts.Claim).(sts.Claims)
	if trafficID == "" {
		return terrors.New(terrors.ErrBadRequest, "Invalid trafficID", map[string]string{})
	}

	err := s.trafficRepo.ResetCounter(ctx, trafficID)
	if err != nil {
		s.log.WithContext(logrus.ErrorLevel,
			"HandleResetCounter",
			"Failed to reset traffic counter",
			logger.Context{
				tracekey.TrackingID: requestID,
				tracekey.UserID:     claims.UserID,
				tracekey.Role:       claims.Role,
			},
			err)
		return err
	}

	return nil
}
