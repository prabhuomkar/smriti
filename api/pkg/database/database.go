package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init ...
func Init(logLevel, host string, port int, username, password, name string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, name)
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(logLevel)),
	})
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func getLogLevel(logLevel string) logger.LogLevel {
	switch logLevel {
	case "ERROR":
		return logger.Error
	case "WARN":
		return logger.Warn
	case "INFO":
		return logger.Info
	}
	return logger.Silent
}
