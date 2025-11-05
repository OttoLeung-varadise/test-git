package service

import (
	"test-git/db"
	"test-git/model"
)

func GetAllRoles(userID string, page, pageSize int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	if err := db.DB.Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.DB.Where("wx_user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func CreateRole(role *model.Role) error {
	return db.DB.Create(role).Error
}
