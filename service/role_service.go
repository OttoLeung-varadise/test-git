package service

import (
	"errors"
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

func GetRoleByID(id uint, userID string) (*model.Role, error) {
	var role model.Role
	result := db.DB.First(&role, id)
	if result.Error != nil {
		return nil, result.Error
	}
	if role.WxUserId != userID {
		return nil, errors.New("无权限查看该角色")
	}
	return &role, nil
}

func CreateRole(role *model.Role) error {
	return db.DB.Create(role).Error
}

func UpdateRole(id uint, updateRole *model.Role) error {
	var role model.Role
	if err := db.DB.First(&role, id).Error; err != nil {
		return err
	}

	if role.WxUserId != updateRole.WxUserId {
		return errors.New("无权限修改该角色")
	}
	return db.DB.Model(&role).Updates(updateRole).Error
}

func DeleteRole(id uint) error {
	return db.DB.Delete(&model.Role{}, id).Error
}

func HardDeleteRole(id uint) error {
	return db.DB.Unscoped().Delete(&model.Role{}, id).Error
}
