package service

import (
	"test-git/db"
	"test-git/model"
)

// CreateBook 添加新书籍
func CreateBook(book *model.Book) error {
	return db.DB.Create(book).Error
}

// GetBookByID 根据 ID 查询书籍
func GetBookByID(id uint) (*model.Book, error) {
	var book model.Book
	result := db.DB.First(&book, id) // First 查询单条记录（按 ID）
	if result.Error != nil {
		return nil, result.Error // 若不存在，返回 gorm.ErrRecordNotFound
	}
	return &book, nil
}

// GetAllBooks 查询所有书籍（支持分页）
func GetAllBooks(page, pageSize int) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	// 先查询总数
	if err := db.DB.Model(&model.Book{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询：offset 偏移量 = (page-1)*pageSize
	offset := (page - 1) * pageSize
	if err := db.DB.Offset(offset).Limit(pageSize).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

// UpdateBook 更新书籍信息
func UpdateBook(id uint, updatedBook *model.Book) error {
	// 先查询原记录是否存在
	var book model.Book
	if err := db.DB.First(&book, id).Error; err != nil {
		return err
	}

	// 更新字段（只更新非零值，或指定字段）
	return db.DB.Model(&book).Updates(updatedBook).Error
}

// DeleteBook 删除书籍（软删除，通过 DeletedAt 标记）
func DeleteBook(id uint) error {
	return db.DB.Delete(&model.Book{}, id).Error
}

// HardDeleteBook 硬删除（直接从数据库删除，谨慎使用）
func HardDeleteBook(id uint) error {
	return db.DB.Unscoped().Delete(&model.Book{}, id).Error
}
