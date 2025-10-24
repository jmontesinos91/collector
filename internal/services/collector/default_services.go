package collector

import (
	"context"
	"fmt"
	straffic "github.com/jmontesinos91/collector/internal/services/traffic"
	"github.com/jmontesinos91/oevents"
	"github.com/jmontesinos91/oevents/eventfactory"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jmontesinos91/collector/internal/repositories/alarmold"
	"github.com/jmontesinos91/collector/internal/repositories/facilitylocationsold"
	"github.com/jmontesinos91/collector/internal/repositories/locationsold"
	"github.com/jmontesinos91/collector/internal/repositories/router"
	"github.com/jmontesinos91/collector/internal/repositories/routerold"
	"github.com/jmontesinos91/collector/internal/repositories/traffic"
	"github.com/jmontesinos91/collector/internal/repositories/unitsold"
	"github.com/jmontesinos91/oevents/broker"
	"github.com/jmontesinos91/ologs/logger"
	tracekey "github.com/jmontesinos91/ologs/logger/v2"
	"github.com/jmontesinos91/terrors"
	"github.com/sirupsen/logrus"
)

// RepositoryOpts ...
type RepositoryOpts struct {
	TrafficRepo       traffic.IRepository
	OldAlarm          alarmold.IRepository
	OldRouter         routerold.IRepository
	OldLocations      locationsold.IRepository
	OldUnits          unitsold.IRepository
	FacilityLocations facilitylocationsold.IRepository
}

// DefaultService struct
type DefaultService struct {
	log               *logger.ContextLogger
	trafficRepo       traffic.IRepository
	oldAlarm          alarmold.IRepository
	oldRouter         routerold.IRepository
	oldLocations      locationsold.IRepository
	oldUnits          unitsold.IRepository
	facilityLocations facilitylocationsold.IRepository
	alarmClient       router.IClient
	streamClient      broker.MessagingBrokerProvider
}

// NewDefaultService creates a new instance of DefaultService Payout
func NewDefaultService(l *logger.ContextLogger, r RepositoryOpts, a router.IClient, bc broker.MessagingBrokerProvider) *DefaultService {
	return &DefaultService{
		log:               l,
		trafficRepo:       r.TrafficRepo,
		oldAlarm:          r.OldAlarm,
		oldRouter:         r.OldRouter,
		oldLocations:      r.OldLocations,
		oldUnits:          r.OldUnits,
		facilityLocations: r.FacilityLocations,
		alarmClient:       a,
		streamClient:      bc,
	}
}

// Collector routers of service of get byID
func (s *DefaultService) Collector(ctx context.Context, payload *Payload) error {
	requestID := ctx.Value(middleware.RequestIDKey).(string)
	var alarmType = "0"
	var isAlarm = false

	if payload.Scare == "P" && (payload.ConfirmPanic == "1" || payload.ConfirmPanic == "2") {

		if payload.ConfirmPanic == "2" {
			alarmType = "3"
		}

		request := router.Request{
			IMEI:      payload.IMEI,
			AlarmType: alarmType,
		}

		//Call to API //wait for the endpoint with IMEI
		response, err := s.alarmClient.ValidateIMEI(ctx, request)
		if err != nil {
			s.log.WithContext(logrus.ErrorLevel,
				"Collector",
				"Error when validate IME I",
				logger.Context{
					tracekey.TrackingID: requestID,
					"IMEI":              payload.IMEI,
				},
				nil)
		}

		if response.Success {
			isAlarm = true
			waiting := "0"
			if payload.Attending == "0" {
				waiting = "1"
			}

			alarm := straffic.Alarm{
				IMEI:      payload.IMEI,
				Latitude:  payload.Latitude,
				Longitude: payload.Longitude,
				AlarmType: alarmType,
				Attending: payload.Attending,
				Waiting:   waiting,
			}

			eventID, err := s.publishAlarmEvent(ctx, alarm, requestID)
			if err != nil {
				s.log.WithContext(
					logrus.ErrorLevel,
					"Collector",
					"The alarm event could not be published:",
					logger.Context{}, err)
			}

			s.log.WithContext(
				logrus.InfoLevel,
				"Collector",
				"Alarm requested event published",
				logger.Context{
					"EventID": eventID,
				}, err)
		}

		errM := s.createOrUpdateTraffic(ctx, payload, isAlarm, requestID)
		if errM != nil {
			return terrors.New(terrors.ErrBadRequest, terrors.MsgBadRequest, map[string]string{})
		}
	} else {
		//Validate UnitID or IMEI
		IsVehicle, routerID, unitID := s.validateRouter(ctx, payload)
		if IsVehicle {
			existAlarm, alarmID, _ := s.oldAlarm.FindByRouterID(ctx, routerID)
			err := s.updateRouterPosition(ctx, routerID, unitID, alarmID, payload.Latitude, payload.Longitude, existAlarm)
			if err != nil {
				return terrors.New(terrors.ErrBadRequest, terrors.MsgBadRequest, map[string]string{})
			}
		}

		err := s.createOrUpdateTraffic(ctx, payload, isAlarm, requestID)
		if err != nil {
			return terrors.New(terrors.ErrBadRequest, terrors.MsgBadRequest, map[string]string{})
		}
	}

	return nil
}

func (s *DefaultService) validateRouter(ctx context.Context, payload *Payload) (bool, int, int) {
	routerModel, err := s.oldRouter.FindByIMEI(ctx, payload.IMEI)
	if err != nil {
		return false, 0, 0
	}

	unit, err := s.oldUnits.FindByRouterID(ctx, routerModel.ID)
	if err != nil {
		return false, 0, 0
	}
	return unit.IsVehicle, routerModel.ID, unit.ID
}

func (s *DefaultService) createOrUpdateTraffic(ctx context.Context, payload *Payload, isAlarm bool, requestID string) error {
	IMEI := payload.IMEI

	finder, err := s.trafficRepo.FindByIMEI(ctx, IMEI, isAlarm)
	if err != nil {
		finder = false
	}

	if !finder {
		model := &traffic.Model{
			ID:         uuid.NewString(),
			Request:    payload.Request,
			IMEI:       IMEI,
			Ip:         payload.IP,
			IsAlarm:    isAlarm,
			IsNotified: false,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		errM := s.trafficRepo.Create(ctx, model)
		if errM != nil {
			return errM
		}
	} else {
		err := s.trafficRepo.UpdateByIMEI(ctx, IMEI, payload.Request, isAlarm)
		if err != nil {
			s.log.WithContext(logrus.ErrorLevel,
				"Collector",
				"Error when try to update traffic",
				logger.Context{
					tracekey.TrackingID: requestID,
					"IMEI":              payload.IMEI,
				}, err)
		}
	}
	return nil
}

func (s *DefaultService) updateRouterPosition(ctx context.Context, routerID, unitID int,
	alarmID, lat, long string, existAlarm bool) error {
	err := s.oldRouter.UpdateLatAndLong(ctx, routerID, lat, long)
	if err != nil {
		return err
	}

	ID, _ := strconv.Atoi(alarmID)
	locationID := 0
	if existAlarm {
		updatedAt := time.Now().UTC()
		locationID, err = s.oldLocations.Create(ctx, &locationsold.LocationsModel{
			AlarmID:   ID,
			Latitude:  lat,
			Longitude: long,
			UpdatedAt: &updatedAt,
		})
		if err != nil {
			return err
		}
	} else {
		locationID, err = s.oldLocations.Create(ctx, &locationsold.LocationsModel{
			AlarmID:   ID,
			Latitude:  lat,
			Longitude: long,
		})
		if err != nil {
			return err
		}
	}

	errFL := s.facilityLocations.Create(ctx, &facilitylocationsold.FacilityLocationsModel{
		LocationID: locationID,
		UnitID:     unitID,
	})
	if errFL != nil {
		return errFL
	}

	return nil
}

func (s *DefaultService) publishAlarmEvent(ctx context.Context, alarm straffic.Alarm, requestID string) (string, error) {

	alarmEvent, err := eventfactory.NewAlarmAcceptedEvent(eventfactory.SourceCollector, ToEventAlarmPayload(alarm, requestID, time.Now().UTC().Format(time.RFC3339)))
	if err != nil {
		return "", err
	}

	ok := s.streamClient.Publish(ctx, oevents.WebHookOmniViewTopic, *alarmEvent)
	if !ok {
		return "", fmt.Errorf("event [%s] could not be published", alarmEvent.EventType)
	}

	return alarmEvent.ID, nil
}
