package db

import (
	"fmt"
	"test-git/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // 全局数据库实例

// Init 初始化数据库连接
func Init() error {
	// PostgreSQL 连接参数（根据实际情况修改）
	dsn := "postgres://postgres:123456@127.0.0.1:5432/bookstore"

	// 连接数据库
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 自动迁移：根据模型创建/更新表结构（生产环境建议手动管理迁移）
	err = DB.AutoMigrate(&model.Book{})
	if err != nil {
		return fmt.Errorf("表结构迁移失败: %v", err)
	}

	return nil
}
