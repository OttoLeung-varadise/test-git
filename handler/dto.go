package handler

import "test-git/model" // 替换为实际的 model 包路径

// CreateBookRequest 创建书籍的请求体
type CreateBookRequest struct {
	Title       string  `json:"title" binding:"required"`  // 书名（必填）
	Author      string  `json:"author" binding:"required"` // 作者（必填）
	Price       float64 `json:"price"`                     // 价格
	Description string  `json:"description"`               // 描述
}

// UpdateBookRequest 更新书籍的请求体
type UpdateBookRequest struct {
	Title       string  `json:"title"`       // 书名（可选，不填则不更新）
	Author      string  `json:"author"`      // 作者（可选）
	Price       float64 `json:"price"`       // 价格（可选）
	Description string  `json:"description"` // 描述（可选）
}

// BookResponse 书籍的响应体（返回给前端的数据）
type BookResponse struct {
	ID          uint    `json:"id"`          // 主键ID
	Title       string  `json:"title"`       // 书名
	Author      string  `json:"author"`      // 作者
	Price       float64 `json:"price"`       // 价格
	Description string  `json:"description"` // 描述
	CreatedAt   string  `json:"created_at"`  // 创建时间
	UpdatedAt   string  `json:"updated_at"`  // 更新时间
}

// 把 model.Book 转换为 BookResponse（格式化时间）
func toBookResponse(book model.Book) BookResponse {
	return BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Price:       book.Price,
		Description: book.Description,
		CreatedAt:   book.CreatedAt.Format("2006-01-02 15:04:05"), // 格式化时间
		UpdatedAt:   book.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
