package db

import (
	"fmt"
	"log"
	"time"

	"github.com/arjnep/gyanpass/config"
	"github.com/arjnep/gyanpass/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	err error
)

func SetupPostgres() {
	cfg := config.GetConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error Connecting to Database: %v", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error Getting sqlDB: %v", err)
		return
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLifetime) * time.Second)

	if err := migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return
	}

}

func migrate() error {
	return db.AutoMigrate(&entity.User{}, &entity.Book{}, &entity.ExchangeRequest{}, &entity.Notification{})
}

func GetDB() *gorm.DB {
	return db
}
