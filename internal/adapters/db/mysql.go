package db

import (
	"database/sql"
	"fmt"
	"github.com/jmontesinos91/collector/config"
	"time"

	_ "github.com/go-sql-driver/mysql" //eslint-disable
	"github.com/jmontesinos91/ologs/logger"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

// NewDatabaseMySQLConnection Initializes a connection pool to the database
func NewDatabaseMySQLConnection(logger *logger.ContextLogger, config config.DatabaseConfigurations) *bun.DB {

	sqlDb, err := sql.Open("mysql", config.Dsn)
	if err != nil {
		logger.Error(logrus.FatalLevel, "DatabaseMySQLConnection", "DB connection error ->", err)
	}

	sqlDb.SetMaxOpenConns(config.Pool)
	sqlDb.SetMaxIdleConns(50)
	sqlDb.SetConnMaxLifetime(1 * time.Minute)

	err = sqlDb.Ping()
	if err != nil {
		logger.Error(logrus.FatalLevel, "DatabaseConnection", "DB connection error -> ", err)
	}

	db := bun.NewDB(sqlDb, mysqldialect.New())
	db.AddQueryHook(bundebug.NewQueryHook())

	logger.Log(logrus.InfoLevel, "Start", fmt.Sprintf("Old Database connected successfully. Connections opened: %d", db.Stats().OpenConnections))

	return db
}
