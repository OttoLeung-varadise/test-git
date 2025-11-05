package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	WxUserId    string         `gorm:"type:varchar(255);not null" json:"wx_user_id"`
	AvatarUrl   string         `gorm:"type:varchar(255);not null" json:"avatar_url"`
	Description string         `gorm:"type:varchar(255);not null" json:"description"`
	RoleData    datatypes.JSON `gorm:"column:role_data;not null" json:"role_data"`
}
