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
//
//	@Summary		创建新书籍
//	@Description	新增一本图书到数据库
//	@Accept			json
//	@Produce		json
//	@Param			book	body		CreateBookRequest	true	"书籍信息"
//	@Success		201		{object}	BookResponse
//	@Failure		400		{string}	string	"请求参数错误"
//	@Failure		500		{string}	string	"服务器内部错误"
//	@Router			/books [post]
func CreateBookHandler(c *gin.Context) {
	var req CreateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：" + err.Error()})
		return
	}

	book := &model.Book{
		Title:       req.Title,
		Author:      req.Author,
		Price:       req.Price,
		Description: req.Description,
	}

	if err := service.CreateBook(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建书籍失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toBookResponse(*book))
}

// GetBookHandler 根据ID查询书籍接口
//
//	@Summary		查询单本书籍
//	@Description	通过ID查询书籍详情
//	@Produce		json
//	@Param			id	path		int	true	"书籍ID"
//	@Success		200	{object}	BookResponse
//	@Failure		400	{string}	string	"ID格式错误"
//	@Failure		404	{string}	string	"书籍不存在"
//	@Failure		500	{string}	string	"服务器内部错误"
//	@Router			/books/{id} [get]
func GetBookHandler(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	book, err := service.GetBookByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "书籍不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败：" + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, toBookResponse(*book))
}

// ListBooksHandler 查询书籍列表接口（支持分页）
//
//	@Summary		查询书籍列表
//	@Description	分页查询所有书籍
//	@Produce		json
//	@Param			page		query		int	false	"页码（默认1）"
//	@Param			pageSize	query		int	false	"每页条数（默认10）"
//	@Success		200			{object}	BookListResponse
//	@Failure		500			{string}	string	"服务器内部错误"
//	@Router			/books [get]
func ListBooksHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	books, total, err := service.GetAllBooks(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询列表失败：" + err.Error()})
		return
	}

	var respList []BookResponse
	for _, book := range books {
		respList = append(respList, toBookResponse(book))
	}

	resp := BookListResponse{
		Total: int(total),
		List:  respList,
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateBookHandler 更新书籍接口
//
//	@Summary		更新书籍信息
//	@Description	根据ID更新书籍信息
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"书籍ID"
//	@Param			book	body		UpdateBookRequest	true	"更新的书籍信息"
//	@Success		204		{string}	string				"更新成功"
//	@Failure		400		{string}	string				"请求参数错误或ID格式错误"
//	@Failure		404		{string}	string				"书籍不存在"
//	@Failure		500		{string}	string				"服务器内部错误"
//	@Router			/books/{id} [put]
func UpdateBookHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：" + err.Error()})
		return
	}

	updatedBook := &model.Book{
		Title:       req.Title,
		Author:      req.Author,
		Price:       req.Price,
		Description: req.Description,
	}

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
//
//	@Summary		删除书籍
//	@Description	根据ID软删除书籍
//	@Produce		json
//	@Param			id	path		int		true	"书籍ID"
//	@Success		204	{string}	string	"删除成功"
//	@Failure		400	{string}	string	"ID格式错误"
//	@Failure		404	{string}	string	"书籍不存在"
//	@Failure		500	{string}	string	"服务器内部错误"
//	@Router			/books/{id} [delete]
func DeleteBookHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	if err := service.DeleteBook(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "删除成功"})
}
