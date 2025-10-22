package db

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmontesinos91/collector/config"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"go.elastic.co/apm/module/apmsql/v2"
)

// NewDatabaseConnection Initializes a connection pool to the database
func NewDatabaseConnection(logger *logger.ContextLogger, config config.DatabaseConfigurations) *bun.DB {

	apmsql.Register("postgres", &stdlib.Driver{})
	sqlDb, err := apmsql.Open("postgres", config.Dsn)
	if err != nil {
		logger.Error(logrus.FatalLevel, "DatabaseConnection", "DB connection error ->", err)
	}

	sqlDb.SetMaxOpenConns(config.Pool + 1)
	sqlDb.SetMaxIdleConns(20)
	sqlDb.SetConnMaxLifetime(3 * time.Minute)

	err = sqlDb.Ping()
	if err != nil {
		logger.Error(logrus.FatalLevel, "DatabaseConnection", "DB connection error -> ", err)
	}

	db := bun.NewDB(sqlDb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook())

	logger.Log(logrus.InfoLevel, "Start", fmt.Sprintf("Database connected successfully. Connections opened: %d", db.Stats().OpenConnections))

	return db
}
