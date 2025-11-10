package db

import (
	"fmt"
	"os"
	"test-git/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func loadDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

func Init() error {
	cfg := loadDBConfig()
	// dsn := "postgres://postgres:123456@127.0.0.1:5432/bookstore"
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("database connetion fails: %v, %s", err, dsn)
	}

	err = DB.AutoMigrate(&model.Role{})
	if err != nil {
		return fmt.Errorf("migrates fails: %v", err)
	}

	return nil
}

func InitLogDB() (*gorm.DB, error) {
	cfg := loadDBConfig()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, "request-log",
	)
	return initLogDB(dsn)
}

func initLogDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get underlying sql.DB")
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	err = model.AutoMigrateRequestLog(db)
	if err != nil {
		return nil, fmt.Errorf("migrates fails: %v", err)
	}
	return db, nil
}
