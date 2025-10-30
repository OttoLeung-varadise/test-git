package db

import (
	"fmt"
	"os"
	"test-git/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // 全局数据库实例

// 数据库配置结构体
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// 从环境变量加载配置
func loadDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

// Init 初始化数据库连接
func Init() error {
	cfg := loadDBConfig()
	// PostgreSQL 连接参数（根据实际情况修改）
	// dsn := "postgres://postgres:123456@127.0.0.1:5432/bookstore"
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	// 连接数据库
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v, %s", err, dsn)
	}

	// 自动迁移：根据模型创建/更新表结构（生产环境建议手动管理迁移）
	err = DB.AutoMigrate(&model.Book{})
	if err != nil {
		return fmt.Errorf("表结构迁移失败: %v", err)
	}

	return nil
}
