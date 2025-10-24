package collector_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmontesinos91/collector/internal/repositories/alarmold/alarmoldmocks"
	"github.com/jmontesinos91/collector/internal/repositories/facilitylocationsold/facilitylocationsoldmocks"
	"github.com/jmontesinos91/collector/internal/repositories/locationsold/locationsoldmocks"
	"github.com/jmontesinos91/collector/internal/repositories/router"
	"github.com/jmontesinos91/collector/internal/repositories/router/routermock"
	"github.com/jmontesinos91/collector/internal/repositories/routerold"
	"github.com/jmontesinos91/collector/internal/repositories/routerold/routeroldmocks"
	otraffic "github.com/jmontesinos91/collector/internal/repositories/traffic"
	"github.com/jmontesinos91/collector/internal/repositories/traffic/trafficmocks"
	"github.com/jmontesinos91/collector/internal/repositories/unitsold"
	"github.com/jmontesinos91/collector/internal/repositories/unitsold/unitsoldmocks"
	"github.com/jmontesinos91/collector/internal/services/collector"
	"github.com/jmontesinos91/oevents/broker/brokermock"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/osecurity/sts"
	"github.com/jmontesinos91/terrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCollect(t *testing.T) {
	var ctxBack context.Context
	log := logger.NewContextLogger("GO-STARTER-TEMPLATE-UNIT-TEST", "debug", logger.TextFormat)
	ctxBack = context.Background()
	ctxBack = context.WithValue(ctxBack, &sts.Claim, sts.Claims{
		UserID: 123,
		Role:   "UnitTest",
	})

	ctxBack = context.WithValue(ctxBack, middleware.RequestIDKey, "unit-test-request-id")
	type fields struct { //nolint:wsl
		routerClient     *routermock.IClient
		routerClientFunc func() *routermock.IClient
		streamClient     *brokermock.MessagingBrokerProvider
		streamClientFunc func() *brokermock.MessagingBrokerProvider
	}
	type repositoryOpts struct { //nolint:wsl
		trafficRepo               *trafficmocks.IRepository
		trafficRepoFunc           func() *trafficmocks.IRepository
		oldAlarmRepo              *alarmoldmocks.IRepository
		oldAlarmRepoFunc          func() *alarmoldmocks.IRepository
		oldRouterRepo             *routeroldmocks.IRepository
		oldRouterRepoFunc         func() *routeroldmocks.IRepository
		oldLocationsRepo          *locationsoldmocks.IRepository
		oldLocationsRepoFunc      func() *locationsoldmocks.IRepository
		oldUnitsRepo              *unitsoldmocks.IRepository
		oldUnitsRepoFunc          func() *unitsoldmocks.IRepository
		facilityLocationsRepo     *facilitylocationsoldmocks.IRepository
		facilityLocationsRepoFunc func() *facilitylocationsoldmocks.IRepository
	}
	type args struct { //nolint:wsl
		ctx     context.Context
		collect *collector.Payload
		model   otraffic.Model
	}
	type assertsParams struct { //nolint:wsl
		args
		fields
		repositoryOpts
	}
	cases := []struct { //nolint:wsl
		name           string
		fields         fields
		repositoryOpts repositoryOpts
		args           args
		err            bool
		asserts        func(*testing.T, error, assertsParams) bool
	}{
		{
			name: "Happy Path",
			fields: fields{
				routerClientFunc: func() *routermock.IClient {
					routerMock := &routermock.IClient{}
					routerMock.On("ValidateIMEI", mock.Anything, mock.Anything).
						Return(&router.Response{Success: true}, nil)
					return routerMock
				},
				streamClientFunc: func() *brokermock.MessagingBrokerProvider {
					streamClientMock := new(brokermock.MessagingBrokerProvider)
					streamClientMock.On("Publish", mock.Anything, mock.Anything, mock.Anything).
						Return(true)
					return streamClientMock
				},
			},
			repositoryOpts: repositoryOpts{
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("FindByIMEI", mock.Anything, mock.Anything, mock.Anything).
						Return(false, nil)
					repositoryMock.On("Create", mock.Anything, mock.Anything).
						Return(nil)
					return repositoryMock
				},
			},
			args: args{
				collect: &collector.Payload{
					Request:      "P,12,12,861585041440544,12,12,123456789,123456789,00,00,00,1",
					IP:           "192.168.100.1",
					IMEI:         "861585041440544",
					Latitude:     "123456789",
					Longitude:    "123456789",
					Attending:    "0",
					ConfirmPanic: "1",
					Scare:        "P",
					GPRS:         "",
				},
				model: otraffic.Model{
					ID:        "12345",
					Request:   "qwertyui128765431425",
					IMEI:      "861585042478659",
					Ip:        "192.168.100.1",
					IsAlarm:   true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			err: false,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				return assert.NoError(t, err) &&
					ap.routerClient.AssertExpectations(t) &&
					ap.routerClient.AssertCalled(t, "ValidateIMEI", mock.Anything, mock.Anything) &&
					ap.trafficRepo.AssertExpectations(t) &&
					ap.trafficRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything) &&
					ap.streamClient.AssertExpectations(t) &&
					ap.streamClient.AssertCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error on create alarm",
			fields: fields{
				routerClientFunc: func() *routermock.IClient {
					routerMock := &routermock.IClient{}
					routerMock.On("ValidateIMEI", mock.Anything, mock.Anything).
						Return(&router.Response{Success: true}, nil)
					return routerMock
				},
				streamClientFunc: func() *brokermock.MessagingBrokerProvider {
					streamClientMock := new(brokermock.MessagingBrokerProvider)
					streamClientMock.On("Publish", mock.Anything, mock.Anything, mock.Anything).
						Return(true)
					return streamClientMock
				},
			},
			repositoryOpts: repositoryOpts{
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("FindByIMEI", mock.Anything, mock.Anything, mock.Anything).
						Return(false, nil)
					repositoryMock.On("Create", mock.Anything, mock.Anything).
						Return(terrors.New(terrors.ErrConflict, terrors.MsgBadRequest, nil))
					return repositoryMock
				},
			},
			args: args{
				collect: &collector.Payload{
					Request:      "P,12,12,,,12,123456789,123456789,00,00,00,1",
					IP:           "192.168.100.1",
					IMEI:         "",
					Latitude:     "123456789",
					Longitude:    "123456789",
					Attending:    "0",
					ConfirmPanic: "1",
					Scare:        "P",
					GPRS:         "",
				},
				model: otraffic.Model{
					ID:        "12345",
					Request:   "qwertyui128765431425",
					IMEI:      "",
					Ip:        "192.168.100.1",
					IsAlarm:   true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			err: true,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, err, &terr) &&
					ap.routerClient.AssertExpectations(t) &&
					ap.routerClient.AssertCalled(t, "ValidateIMEI", mock.Anything, mock.Anything) &&
					ap.trafficRepo.AssertExpectations(t) &&
					ap.trafficRepo.AssertCalled(t, "Create", mock.Anything, mock.Anything)
			},
		},
		{
			name: "Error on create traffic",
			fields: fields{
				routerClientFunc: func() *routermock.IClient {
					routerMock := &routermock.IClient{}
					routerMock.On("ValidateIMEI", mock.Anything, mock.Anything).
						Return(&router.Response{Success: true}, nil)
					return routerMock
				},
			},
			repositoryOpts: repositoryOpts{
				oldRouterRepoFunc: func() *routeroldmocks.IRepository {
					repositoryMock := &routeroldmocks.IRepository{}
					repositoryMock.On("UpdateLatAndLong",
						mock.Anything, mock.Anything, mock.Anything, mock.Anything).
						Return(nil)
					repositoryMock.On("FindByIMEI", mock.Anything, mock.Anything).
						Return(&routerold.RouterModel{}, nil)
					return repositoryMock
				},
				oldUnitsRepoFunc: func() *unitsoldmocks.IRepository {
					repositoryMock := &unitsoldmocks.IRepository{}
					repositoryMock.On("FindByRouterID", mock.Anything, 0).
						Return(&unitsold.UnitsModel{ID: 53438, IsVehicle: true, RouterID: 1}, nil)
					return repositoryMock
				},
				oldAlarmRepoFunc: func() *alarmoldmocks.IRepository {
					repositoryMock := &alarmoldmocks.IRepository{}
					repositoryMock.On("FindByRouterID", mock.Anything, mock.Anything).
						Return(true, "1", nil)
					return repositoryMock
				},
				oldLocationsRepoFunc: func() *locationsoldmocks.IRepository {
					repositoryMock := &locationsoldmocks.IRepository{}
					repositoryMock.On("Create", mock.Anything, mock.Anything).
						Return(1, nil)
					return repositoryMock
				},
				facilityLocationsRepoFunc: func() *facilitylocationsoldmocks.IRepository {
					repositoryMock := &facilitylocationsoldmocks.IRepository{}
					repositoryMock.On("Create", mock.Anything, mock.Anything).
						Return(nil)
					return repositoryMock
				},
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("FindByIMEI", mock.Anything, mock.Anything, mock.Anything).
						Return(false, nil)
					repositoryMock.On("Create", mock.Anything, mock.Anything).
						Return(nil)
					return repositoryMock
				},
			},
			args: args{
				collect: &collector.Payload{
					Request:      "0000002c0,12,12,,,12,123456789,123456789,00,00,00,1",
					IP:           "192.168.100.1",
					IMEI:         "",
					Latitude:     "123456789",
					Longitude:    "123456789",
					Attending:    "0",
					ConfirmPanic: "1",
					Scare:        "0",
					GPRS:         "",
				},
				model: otraffic.Model{
					ID:        "12345",
					Request:   "qwertyui128765431425",
					IMEI:      "",
					Ip:        "192.168.100.1",
					IsAlarm:   true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			err: false,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				return assert.NoError(t, err) &&
					ap.oldRouterRepo.AssertExpectations(t) &&
					ap.oldRouterRepo.AssertCalled(t, "UpdateLatAndLong", ap.ctx,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything) &&
					ap.oldUnitsRepo.AssertExpectations(t) &&
					ap.oldUnitsRepo.AssertCalled(t, "FindByRouterID", ap.ctx, 0) &&
					ap.oldLocationsRepo.AssertCalled(t, "Create", ap.ctx,
						mock.Anything) &&
					ap.facilityLocationsRepo.AssertExpectations(t) &&
					ap.facilityLocationsRepo.AssertCalled(t, "Create", ap.ctx,
						mock.Anything)
			},
		},
		{
			name: "Happy Path Exist With Out Alarm",
			repositoryOpts: repositoryOpts{
				oldRouterRepoFunc: func() *routeroldmocks.IRepository {
					repositoryMock := &routeroldmocks.IRepository{}
					repositoryMock.On("FindByIMEI", mock.Anything, mock.Anything).
						Return(&routerold.RouterModel{ID: 1}, nil)
					repositoryMock.On("UpdateLatAndLong",
						mock.Anything, mock.Anything, mock.Anything, mock.Anything).
						Return(nil)
					return repositoryMock
				},
				oldUnitsRepoFunc: func() *unitsoldmocks.IRepository {
					repositoryMock := &unitsoldmocks.IRepository{}
					repositoryMock.On("FindByRouterID", mock.Anything, mock.Anything).
						Return(&unitsold.UnitsModel{ID: 1, IsVehicle: true}, nil)
					return repositoryMock
				},
				oldAlarmRepoFunc: func() *alarmoldmocks.IRepository {
					repositoryMock := &alarmoldmocks.IRepository{}
					repositoryMock.On("FindByRouterID", mock.Anything, mock.Anything).
						Return(true, "1", nil)
					return repositoryMock
				},
				oldLocationsRepoFunc: func() *locationsoldmocks.IRepository {
					repositoryMock := &locationsoldmocks.IRepository{}
					repositoryMock.On("Create", mock.Anything, mock.Anything).
						Return(1, nil)
					return repositoryMock
				},
				facilityLocationsRepoFunc: func() *facilitylocationsoldmocks.IRepository {
					repositoryMock := &facilitylocationsoldmocks.IRepository{}
					repositoryMock.On("Create", mock.Anything, mock.Anything).
						Return(nil)
					return repositoryMock
				},
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("FindByIMEI", mock.Anything, mock.Anything, mock.Anything).
						Return(false, nil)
					repositoryMock.On("Create", mock.Anything, mock.Anything).
						Return(nil)
					return repositoryMock
				},
			},
			args: args{
				collect: &collector.Payload{
					Request:      "P,12,12,861585041440544,12,12,123456789,123456789,00,00,00,0",
					IP:           "192.168.100.1",
					IMEI:         "861585041440544",
					Latitude:     "123456789",
					Longitude:    "123456789",
					Attending:    "0",
					ConfirmPanic: "0",
					Scare:        "P",
					GPRS:         "",
				},
				model: otraffic.Model{
					ID:        "12345",
					Request:   "qwertyui128765431425",
					IMEI:      "IMEI_Edificio_Principal",
					Ip:        "192.168.100.1",
					IsAlarm:   false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			err: false,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				return assert.NoError(t, err) &&
					ap.oldRouterRepo.AssertExpectations(t) &&
					ap.oldRouterRepo.AssertCalled(t, "FindByIMEI", ap.ctx, mock.Anything) &&
					ap.oldRouterRepo.AssertCalled(t, "UpdateLatAndLong", ap.ctx,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything) &&
					ap.oldUnitsRepo.AssertExpectations(t) &&
					ap.oldUnitsRepo.AssertCalled(t, "FindByRouterID", ap.ctx,
						mock.Anything) &&
					ap.oldAlarmRepo.AssertExpectations(t) &&
					ap.oldAlarmRepo.AssertCalled(t, "FindByRouterID", ap.ctx,
						mock.Anything) &&
					ap.oldLocationsRepo.AssertExpectations(t) &&
					ap.oldLocationsRepo.AssertCalled(t, "Create", ap.ctx,
						mock.Anything) &&
					ap.facilityLocationsRepo.AssertExpectations(t) &&
					ap.facilityLocationsRepo.AssertCalled(t, "Create", ap.ctx,
						mock.Anything)
			},
		},
		{
			name: "Error on update db",
			fields: fields{
				routerClientFunc: func() *routermock.IClient {
					routerMock := &routermock.IClient{}
					routerMock.On("ValidateIMEI", mock.Anything, mock.Anything).
						Return(&router.Response{Success: true}, nil)
					return routerMock
				},
				streamClientFunc: func() *brokermock.MessagingBrokerProvider {
					streamClientMock := new(brokermock.MessagingBrokerProvider)
					streamClientMock.On("Publish", mock.Anything, mock.Anything, mock.Anything).
						Return(true)
					return streamClientMock
				},
			},
			repositoryOpts: repositoryOpts{
				oldRouterRepoFunc: func() *routeroldmocks.IRepository {
					repositoryMock := &routeroldmocks.IRepository{}
					repositoryMock.On("FindByIMEI", mock.Anything, mock.Anything).
						Return(&routerold.RouterModel{ID: 1}, nil)
					repositoryMock.On("UpdateLatAndLong",
						mock.Anything, mock.Anything, mock.Anything, mock.Anything).
						Return(nil)
					repositoryMock.On("FindByIMEI", mock.Anything, mock.Anything).
						Return(&routerold.RouterModel{}, nil)
					return repositoryMock
				},
				oldUnitsRepoFunc: func() *unitsoldmocks.IRepository {
					repositoryMock := &unitsoldmocks.IRepository{}
					repositoryMock.On("FindByRouterID", mock.Anything, mock.Anything).
						Return(&unitsold.UnitsModel{ID: 1, IsVehicle: true}, nil)
					return repositoryMock
				},
				oldAlarmRepoFunc: func() *alarmoldmocks.IRepository {
					repositoryMock := &alarmoldmocks.IRepository{}
					repositoryMock.On("FindByRouterID", mock.Anything, mock.Anything).
						Return(true, "1", nil)
					return repositoryMock
				},
				oldLocationsRepoFunc: func() *locationsoldmocks.IRepository {
					repositoryMock := &locationsoldmocks.IRepository{}
					repositoryMock.On("Create", mock.Anything, mock.Anything).
						Return(1, nil)
					return repositoryMock
				},
				facilityLocationsRepoFunc: func() *facilitylocationsoldmocks.IRepository {
					repositoryMock := &facilitylocationsoldmocks.IRepository{}
					repositoryMock.On("Create", mock.Anything, mock.Anything).
						Return(nil)
					return repositoryMock
				},
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("FindByIMEI", mock.Anything, mock.Anything, mock.Anything).
						Return(true, nil)
					repositoryMock.On("UpdateByIMEI", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
						Return(
							terrors.New(terrors.ErrBadRequest, "Internal error service", map[string]string{}))
					return repositoryMock
				},
			},
			args: args{
				collect: &collector.Payload{
					Request:      "P,12,12,861585041440544,12,12,123456789,123456789,00,00,0,1",
					IP:           "192.168.100.1",
					IMEI:         "861585041440544",
					Latitude:     "123456789",
					Longitude:    "123456789",
					Attending:    "0",
					ConfirmPanic: "1",
					Scare:        "P",
					GPRS:         "",
				},
				model: otraffic.Model{
					ID:        "12345",
					Request:   "qwertyui128765431425",
					IMEI:      "IMEI_Edificio_Principal",
					Ip:        "192.168.100.1",
					IsAlarm:   false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			err: false,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				return assert.NoError(t, err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.args.ctx == nil {
				tc.args.ctx = ctxBack
			}

			if tc.fields.streamClientFunc != nil {
				tc.fields.streamClient = tc.fields.streamClientFunc()
			}

			if tc.fields.routerClientFunc != nil {
				tc.fields.routerClient = tc.fields.routerClientFunc()
			}

			if tc.repositoryOpts.trafficRepoFunc != nil {
				tc.repositoryOpts.trafficRepo = tc.repositoryOpts.trafficRepoFunc()
			}

			if tc.repositoryOpts.oldAlarmRepoFunc != nil {
				tc.repositoryOpts.oldAlarmRepo = tc.repositoryOpts.oldAlarmRepoFunc()
			}

			if tc.repositoryOpts.oldRouterRepoFunc != nil {
				tc.repositoryOpts.oldRouterRepo = tc.repositoryOpts.oldRouterRepoFunc()
			}

			if tc.repositoryOpts.oldLocationsRepoFunc != nil {
				tc.repositoryOpts.oldLocationsRepo = tc.repositoryOpts.oldLocationsRepoFunc()
			}

			if tc.repositoryOpts.oldUnitsRepoFunc != nil {
				tc.repositoryOpts.oldUnitsRepo = tc.repositoryOpts.oldUnitsRepoFunc()
			}

			if tc.repositoryOpts.facilityLocationsRepoFunc != nil {
				tc.repositoryOpts.facilityLocationsRepo = tc.repositoryOpts.facilityLocationsRepoFunc()
			}

			repoOpts := collector.RepositoryOpts{
				TrafficRepo:       tc.repositoryOpts.trafficRepo,
				OldRouter:         tc.repositoryOpts.oldRouterRepo,
				OldAlarm:          tc.repositoryOpts.oldAlarmRepo,
				OldLocations:      tc.repositoryOpts.oldLocationsRepo,
				OldUnits:          tc.repositoryOpts.oldUnitsRepo,
				FacilityLocations: tc.repositoryOpts.facilityLocationsRepo,
			}

			collectorService := collector.NewDefaultService(log,
				repoOpts,
				tc.fields.routerClient,
				tc.fields.streamClient)

			err := collectorService.Collector(tc.args.ctx, tc.args.collect)
			if (err != nil) != tc.err {
				t.Errorf("DefaultService.FindByID() error = %v, wantErr %v", err, tc.err)
			}

			assertsParams := assertsParams{
				fields:         tc.fields,
				repositoryOpts: tc.repositoryOpts,
				args:           tc.args,
			}

			if !tc.asserts(t, err, assertsParams) {
				t.Errorf("Assert error on test = '%v'", tc.name)
			}
		})
	}
}
