package model

import "gorm.io/gorm"

// Book 书籍模型，映射数据库中的 books 表
type Book struct {
	gorm.Model          // 嵌入 GORM 内置模型，自动包含 ID、CreatedAt、UpdatedAt、DeletedAt 字段
	Title       string  `gorm:"type:varchar(255);not null" json:"title"`  // 书名，非空
	Author      string  `gorm:"type:varchar(100);not null" json:"author"` // 作者，非空
	Price       float64 `gorm:"type:decimal(10,2)" json:"price"`          // 价格，保留2位小数
	Description string  `gorm:"type:text" json:"description"`             // 描述，长文本
}
