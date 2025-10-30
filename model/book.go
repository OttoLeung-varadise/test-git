package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string  `gorm:"type:varchar(255);not null" json:"title"`
	Author      string  `gorm:"type:varchar(100);not null" json:"author"`
	Price       float64 `gorm:"type:decimal(10,2)" json:"price"`
	Description string  `gorm:"type:text" json:"description"`
}
