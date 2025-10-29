package handler

import (
	"net/http"
	"strconv"
	"test-git/model"
	"test-git/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateBookHandler 创建书籍接口
// @Summary 创建新书籍
// @Description 新增一本图书到数据库
// @Accept json
// @Produce json
// @Param book body CreateBookRequest true "书籍信息"
// @Success 200 {object} BookResponse
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "服务器内部错误"
// @Router /books [post]
func CreateBookHandler(c *gin.Context) {
	var req CreateBookRequest
	// 绑定并验证请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：" + err.Error()})
		return
	}

	// 转换为 model.Book
	book := &model.Book{
		Title:       req.Title,
		Author:      req.Author,
		Price:       req.Price,
		Description: req.Description,
	}

	// 调用 service 创建书籍
	if err := service.CreateBook(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建书籍失败：" + err.Error()})
		return
	}

	// 返回创建成功的书籍信息
	c.JSON(http.StatusCreated, toBookResponse(*book))
}

// GetBookHandler 根据ID查询书籍接口
// @Summary 查询单本书籍
// @Description 通过ID查询书籍详情
// @Produce json
// @Param id path int true "书籍ID"
// @Success 200 {object} BookResponse
// @Failure 400 {string} string "ID格式错误"
// @Failure 404 {string} string "书籍不存在"
// @Failure 500 {string} string "服务器内部错误"
// @Router /books/{id} [get]
func GetBookHandler(c *gin.Context) {
	// 从URL路径中获取ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	// 调用 service 查询书籍
	book, err := service.GetBookByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "书籍不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败：" + err.Error()})
		}
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, toBookResponse(*book))
}

// ListBooksHandler 查询书籍列表接口（支持分页）
// @Summary 查询书籍列表
// @Description 分页查询所有书籍
// @Produce json
// @Param page query int false "页码（默认1）"
// @Param pageSize query int false "每页条数（默认10）"
// @Success 200 {object} gin.H{ "total": int, "list": []BookResponse }
// @Failure 500 {string} string "服务器内部错误"
// @Router /books [get]
func ListBooksHandler(c *gin.Context) {
	// 获取分页参数（默认第1页，每页10条）
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 调用 service 查询列表
	books, total, err := service.GetAllBooks(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询列表失败：" + err.Error()})
		return
	}

	// 转换为响应格式
	var respList []BookResponse
	for _, book := range books {
		respList = append(respList, toBookResponse(book))
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"total": total, // 总条数
		"list":  respList,
	})
}

// UpdateBookHandler 更新书籍接口
// @Summary 更新书籍信息
// @Description 根据ID更新书籍信息
// @Accept json
// @Produce json
// @Param id path int true "书籍ID"
// @Param book body UpdateBookRequest true "更新的书籍信息"
// @Success 200 {string} string "更新成功"
// @Failure 400 {string} string "请求参数错误或ID格式错误"
// @Failure 404 {string} string "书籍不存在"
// @Failure 500 {string} string "服务器内部错误"
// @Router /books/{id} [put]
func UpdateBookHandler(c *gin.Context) {
	// 解析ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	// 绑定请求体
	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：" + err.Error()})
		return
	}

	// 转换为 model.Book（只更新非空字段）
	updatedBook := &model.Book{
		Title:       req.Title,
		Author:      req.Author,
		Price:       req.Price,
		Description: req.Description,
	}

	// 调用 service 更新书籍
	if err := service.UpdateBook(uint(id), updatedBook); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "書籍記錄不存在:" + err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败：" + err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "更新成功"})
}

// DeleteBookHandler 删除书籍接口
// @Summary 删除书籍
// @Description 根据ID软删除书籍
// @Produce json
// @Param id path int true "书籍ID"
// @Success 200 {string} string "删除成功"
// @Failure 400 {string} string "ID格式错误"
// @Failure 404 {string} string "书籍不存在"
// @Failure 500 {string} string "服务器内部错误"
// @Router /books/{id} [delete]
func DeleteBookHandler(c *gin.Context) {
	// 解析ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	// 调用 service 删除书籍
	if err := service.DeleteBook(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "删除成功"})
}
