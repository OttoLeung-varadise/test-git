package service

import (
	"test-git/db"
	"test-git/model"
)

func CreateBook(book *model.Book) error {
	return db.DB.Create(book).Error
}

func GetBookByID(id uint) (*model.Book, error) {
	var book model.Book
	result := db.DB.First(&book, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func GetAllBooks(page, pageSize int) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	if err := db.DB.Model(&model.Book{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.DB.Offset(offset).Limit(pageSize).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func UpdateBook(id uint, updatedBook *model.Book) error {
	var book model.Book
	if err := db.DB.First(&book, id).Error; err != nil {
		return err
	}

	return db.DB.Model(&book).Updates(updatedBook).Error
}

func DeleteBook(id uint) error {
	return db.DB.Delete(&model.Book{}, id).Error
}

func HardDeleteBook(id uint) error {
	return db.DB.Unscoped().Delete(&model.Book{}, id).Error
}
