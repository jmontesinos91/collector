package main

import (
	"github.com/jmontesinos91/collector/config"
	"github.com/jmontesinos91/collector/internal/adapters/api"
	"github.com/jmontesinos91/collector/internal/adapters/db"
	"github.com/jmontesinos91/collector/internal/adapters/stream"
	"github.com/jmontesinos91/collector/internal/repositories/alarmold" //nolint:goimports
	"github.com/jmontesinos91/collector/internal/repositories/facilitylocationsold"
	"github.com/jmontesinos91/collector/internal/repositories/locationsold"
	"github.com/jmontesinos91/collector/internal/repositories/router"
	"github.com/jmontesinos91/collector/internal/repositories/routerold" //nolint:goimports
	trepository "github.com/jmontesinos91/collector/internal/repositories/traffic"
	"github.com/jmontesinos91/collector/internal/repositories/unitsold" //nolint:goimports
	"github.com/jmontesinos91/collector/internal/services/collector"
	"github.com/jmontesinos91/collector/internal/services/traffic"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/osecurity/services/omnibackend"
	"github.com/jmontesinos91/osecurity/sts"

	"github.com/go-playground/validator/v10"
)

func main() {
	// Logger
	contextLogger := logger.NewContextLogger("COLLECTOR", "debug", logger.TextFormat)

	// Configs
	configs := config.LoadConfig(contextLogger)

	// -- Start dependency injection section --
	// STS Client
	omniService := omnibackend.NewOmniViewService(contextLogger, configs.OmniView)
	stsClient := sts.NewDefaultISTSClient(contextLogger, omniService)

	// Http Router
	httpServer := api.NewHTTPServer(contextLogger, configs.Server, configs.Service, stsClient)

	// Validator
	validate := validator.New()

	// DB Connection
	conn := db.NewDatabaseConnection(contextLogger, configs.Database)

	// DB Connection to old DB
	oldConn := db.NewDatabaseMySQLConnection(contextLogger, configs.OldDatabase)

	// Kafka
	kafka, closer := stream.NewKafkaConnection(contextLogger, configs.Kafka)
	defer closer()

	// Alarm Client
	rClient := router.NewRouterService(contextLogger, configs.OmniView)

	// - Initialize repository -
	trafficRepo := trepository.NewDatabaseRepository(contextLogger, conn)
	oldLocations := locationsold.NewDatabaseRepository(contextLogger, oldConn)
	oldRouter := routerold.NewDatabaseRepository(contextLogger, oldConn)
	oldFacilityLocations := facilitylocationsold.NewDatabaseRepository(contextLogger, oldConn)
	oldUnits := unitsold.NewDatabaseRepository(contextLogger, oldConn)
	oldAlarm := alarmold.NewDatabaseRepository(contextLogger, oldConn)

	repositoryOpts := collector.RepositoryOpts{
		TrafficRepo:  trafficRepo,
		OldAlarm:     oldAlarm,
		OldRouter:    oldRouter,
		OldLocations: oldLocations,
		OldUnits:     oldUnits,

		FacilityLocations: oldFacilityLocations,
	}

	// - Initialize service -
	collectorSvc := collector.NewDefaultService(contextLogger, repositoryOpts, rClient, kafka)
	trafficSvc := traffic.NewDefaultService(contextLogger, trafficRepo)

	api.NewHealthController(httpServer)
	api.NewCollectorController(httpServer, validate, collectorSvc, stsClient)
	api.NewTrafficController(httpServer, validate, trafficSvc, stsClient)

	// -- End dependency injection section --

	// Let the party started!
	httpServer.Start()
}
