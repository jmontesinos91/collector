package traffic_test

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmontesinos91/collector/domains/pagination"
	otraffic "github.com/jmontesinos91/collector/internal/repositories/traffic"
	"github.com/jmontesinos91/collector/internal/repositories/traffic/trafficmocks"
	"github.com/jmontesinos91/collector/internal/services/traffic"
	straffic "github.com/jmontesinos91/collector/internal/services/traffic"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/osecurity/sts"
	"github.com/jmontesinos91/terrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRetrieve(t *testing.T) {
	ctxBack := context.Background()
	ctxBack = context.WithValue(ctxBack, middleware.RequestIDKey, "unit-test-request-id")
	ctxBack = context.WithValue(ctxBack, &sts.Claim, sts.Claims{
		UserID: 0,
		Role:   "unit-test-role",
	})
	log := logger.NewContextLogger("Retrieve", "debug", logger.TextFormat)

	type repositoryOpts struct { //nolint:wsl
		trafficRepo     *trafficmocks.IRepository
		trafficRepoFunc func() *trafficmocks.IRepository
	}

	type args struct { //nolint:wsl
		ctx      context.Context
		filter   *straffic.FilterRequest
		expected pagination.PaginatedRes
	}

	type assertsParams struct { //nolint:wsl
		args
		repositoryOpts
		result pagination.PaginatedRes
	}

	cases := []struct { //nolint:wsl
		name           string
		repositoryOpts repositoryOpts
		args           args
		err            bool
		asserts        func(*testing.T, error, assertsParams) bool
	}{
		{
			name: "Happy path",
			repositoryOpts: repositoryOpts{
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("Retrieve", mock.Anything, mock.Anything).
						Return([]otraffic.Model{}, 1, 10, nil)
					return repositoryMock
				},
			},
			args: args{
				filter: &straffic.FilterRequest{
					ID: "209e7c87-84f9-41d0-a5b7-002f5d8d886ds",
					Ip: "192.168.1.100",
					Filter: pagination.Filter{
						Page:     1,
						Size:     10,
						Offset:   0,
						SortBy:   "created_at",
						SortDesc: true,
					},
				},
				expected: pagination.PaginatedRes{
					Data:        []straffic.Traffic(nil),
					CurrentPage: 1,
					Pages:       1,
					Total:       10,
				},
			},
			err: false,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				return assert.NotNil(t, ap.result) &&
					assert.NoError(t, err) &&
					ap.trafficRepo.AssertExpectations(t)
			},
		},
		{
			name: "Happy path with sanitize",
			repositoryOpts: repositoryOpts{
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("Retrieve", mock.Anything, mock.Anything).
						Return([]otraffic.Model{}, 0, 0, nil)

					return repositoryMock
				},
			},
			args: args{
				filter: &straffic.FilterRequest{
					ID: "209e7c87-84f9-41d0-a5b7-002f5d8d886ds",
					Ip: "192.168.1.100",
					Filter: pagination.Filter{
						Page:   1,
						Size:   100,
						Offset: 100,
					},
				},
				expected: pagination.PaginatedRes{},
			},
			err: false,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				return assert.NotNil(t, ap.result) &&
					assert.NoError(t, err) &&
					ap.trafficRepo.AssertExpectations(t)
			},
		},
		{
			name: "Empty filters",
			repositoryOpts: repositoryOpts{
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("Retrieve", mock.Anything, mock.Anything).
						Return([]otraffic.Model{}, 1, 10, nil)
					return repositoryMock
				},
			},
			args: args{
				filter: &straffic.FilterRequest{},
				expected: pagination.PaginatedRes{
					Data:        []straffic.Traffic(nil),
					CurrentPage: 1,
					Pages:       1,
					Total:       10,
				},
			},
			err: false,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				return assert.NotNil(t, ap.result) &&
					assert.NoError(t, err) &&
					ap.trafficRepo.AssertExpectations(t)
			},
		},
		{
			name: "Error on retrieve",
			repositoryOpts: repositoryOpts{
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("Retrieve", mock.Anything, mock.Anything).
						Return([]otraffic.Model{}, 0, 0, terrors.New(terrors.ErrBadRequest, "Failed to retrieve traffics", map[string]string{}))

					return repositoryMock
				},
			},
			args: args{
				filter:   &straffic.FilterRequest{},
				expected: pagination.PaginatedRes{},
			},
			err: true,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, err, &terr) &&
					assert.Equal(t, "Failed to retrieve traffics", terr.Message) &&
					assert.Equal(t, ap.expected, ap.result)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.args.ctx == nil {
				tc.args.ctx = ctxBack
			}

			if tc.repositoryOpts.trafficRepoFunc != nil {
				tc.repositoryOpts.trafficRepo = tc.repositoryOpts.trafficRepoFunc()
			}

			trafficSvc := traffic.NewDefaultService(log, tc.repositoryOpts.trafficRepo)

			result, err := trafficSvc.HandleRetrieve(tc.args.ctx, tc.args.filter)
			if (err != nil) != tc.err {
				t.Errorf("DefaultService.HandleRetrieve() error = %v, wantErr %v", err, tc.err)
			}

			assertsParams := assertsParams{
				repositoryOpts: tc.repositoryOpts,
				args:           tc.args,
				result:         result,
			}

			if !tc.asserts(t, err, assertsParams) {
				t.Errorf("Assert error on test = '%v'", tc.name)
			}

		})
	}

}

func TestDeleteTrafficByID(t *testing.T) {
	ctxBack := context.Background()
	ctxBack = context.WithValue(ctxBack, middleware.RequestIDKey, "unit-test-request-id")
	ctxBack = context.WithValue(ctxBack, &sts.Claim, sts.Claims{
		UserID: 0,
		Role:   "unit-test-role",
	})
	log := logger.NewContextLogger("Retrieve", "debug", logger.TextFormat)

	type repositoryOpts struct { //nolint:wsl
		trafficRepo     *trafficmocks.IRepository
		trafficRepoFunc func() *trafficmocks.IRepository
	}

	type args struct { //nolint:wsl
		ctx       context.Context
		trafficID string
	}

	type assertsParams struct { //nolint:wsl
		args
		repositoryOpts
		result error
	}

	cases := []struct { //nolint:wsl
		name           string
		repositoryOpts repositoryOpts
		args           args
		err            bool
		asserts        func(*testing.T, error, assertsParams) bool
	}{
		{
			name: "Happy path",
			repositoryOpts: repositoryOpts{
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("DeleteByID", mock.Anything, mock.Anything).
						Return(nil)
					return repositoryMock
				},
			},
			args: args{
				trafficID: "unit-test-traffic-id",
			},
			err: false,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				return assert.Nil(t, ap.result) &&
					assert.NoError(t, err) &&
					ap.trafficRepo.AssertExpectations(t)
			},
		},
		{
			name: "Empty traffic id",
			args: args{
				trafficID: "",
			},
			err: true,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, err, &terr) &&
					assert.Equal(t, "Invalid trafficID", terr.Message)
			},
		},
		{
			name: "Error on delete",
			repositoryOpts: repositoryOpts{
				trafficRepoFunc: func() *trafficmocks.IRepository {
					repositoryMock := &trafficmocks.IRepository{}
					repositoryMock.On("DeleteByID", mock.Anything, mock.Anything).
						Return(terrors.InternalService("delete_traffic", "Failed delete traffic from the database", map[string]string{}))

					return repositoryMock
				},
			},
			args: args{
				trafficID: "unit-test-traffic-id",
			},
			err: true,
			asserts: func(t *testing.T, err error, ap assertsParams) bool {
				var terr *terrors.Error
				return assert.ErrorAs(t, err, &terr) &&
					assert.Equal(t, "Failed delete traffic from the database", terr.Message)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.args.ctx == nil {
				tc.args.ctx = ctxBack
			}

			if tc.repositoryOpts.trafficRepoFunc != nil {
				tc.repositoryOpts.trafficRepo = tc.repositoryOpts.trafficRepoFunc()
			}

			trafficSvc := traffic.NewDefaultService(log, tc.repositoryOpts.trafficRepo)

			err := trafficSvc.HandleDelete(tc.args.ctx, tc.args.trafficID)
			if (err != nil) != tc.err {
				t.Errorf("DefaultService.HandleResetCounter() error = %v, wantErr %v", err, tc.err)
			}

			assertsParams := assertsParams{
				repositoryOpts: tc.repositoryOpts,
				args:           tc.args,
			}

			if !tc.asserts(t, err, assertsParams) {
				t.Errorf("Assert error on test = '%v'", tc.name)
			}

		})
	}

}
