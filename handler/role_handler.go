package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateRoleHandler 創建角色接口
//
//	@Summary		创建新角色
//	@Description	讀取一個角色excel，返回json
//
//	@Accept			multipart/form-data
//
//	@Produce		json
//	@Param			file	formData	file	true	"要上传的文件"
//	@Success		200		{object}	COCRoleCard
//	@Failure		400		{string}	string	"请求参数错误"
//	@Failure		500		{string}	string	"服务器内部错误"
//	@Router			/roles [post]
func PreviewRoleHandler(c *gin.Context) {

	var roleCard COCRoleCard
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "解析表单失败: " + err.Error(),
		})
		return
	}

	if c.Request.MultipartForm == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "表单数据为空",
		})
		return
	}

	files, ok := c.Request.MultipartForm.File["file"]
	if !ok || len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "未找到名为 'file' 的文件",
		})
		return
	}

	err := DecodeFile(&roleCard, files[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件解析失敗：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, roleCard)
}

func DecodeFile(input interface{}, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(input)
}
