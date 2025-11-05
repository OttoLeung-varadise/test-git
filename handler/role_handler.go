package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"
	"test-git/common"
	"test-git/model"
	"test-git/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	err := decodeFile(&roleCard, files[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件解析失敗：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, roleCard)
}

// ListRoleHandler 查询角色列表接口（支持分页）
//
//	@Summary		查询角色列表
//	@Description	分页查询所有角色
//	@Produce		json
//	@Param			page		query		int	false	"页码（默认1）"
//	@Param			pageSize	query		int	false	"每页条数（默认10）"
//	@Success		200			{object}	RoleListResponse
//	@Failure		500			{string}	string	"服务器内部错误"
//	@Router			/roles [get]
func ListRoleHandler(c *gin.Context) {
	userID := common.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	roles, total, err := service.GetAllRoles(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询列表失败：" + err.Error()})
		return
	}

	var respList []RoleResponse
	for _, role := range roles {
		respList = append(respList, toRoleResponse(role, false))
	}

	resp := RoleListResponse{
		Total: int(total),
		List:  respList,
	}
	c.JSON(http.StatusOK, resp)
}

// GetRoleHandler 查询角色详情接口
//
//	@Summary		查询角色详情
//	@Description	根据ID查询角色详情
//	@Produce		json
//	@Param			id		path		int	true	"角色ID"
//	@Success		200		{object}	RoleResponse
//	@Failure		400		{string}	string	"请求参数错误或ID格式错误"
//	@Failure		404		{string}	string	"角色不存在"
//	@Failure		500		{string}	string	"服务器内部错误"
//	@Router			/roles/{id} [get]
func GetRoleHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	role, err := service.GetRoleByID(uint(id), common.GetUserID(c))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败：" + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, toRoleResponse(*role, true))
}

// CreateRoleHandler 创建角色接口
//
//	@Summary		创建新角色
//	@Description	新增一個角色到数据库
//	@Accept			json
//	@Produce		json
//	@Param			book	body		CreateRoleHandler	true	"角色信息"
//	@Success		201		{object}	RoleResponse
//	@Failure		400		{string}	string	"请求参数错误"
//	@Failure		500		{string}	string	"服务器内部错误"
//	@Router			/roles/create [post]
func CreateRoleHandler(c *gin.Context) {
	var req CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：" + err.Error()})
		return
	}

	roleCard := req.RoleData

	if req.AvatarUrl != roleCard.BasicInfo.AvatarURL && req.AvatarUrl != "" {
		roleCard.BasicInfo.AvatarURL = req.AvatarUrl
	}

	roleJSON, err := json.Marshal(&roleCard)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色卡格式错误：" + err.Error()})
		return
	}

	role := &model.Role{
		Name:        roleCard.BasicInfo.RoleName,
		WxUserId:    common.GetUserID(c),
		AvatarUrl:   req.AvatarUrl,
		Description: getRoleDesc(&roleCard),
		RoleData:    roleJSON,
	}

	if err := service.CreateRole(role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建角色失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toRoleResponse(*role, true))
}

// UpdateRoleHandler 更新角色接口
//
//	@Summary		更新角色信息
//	@Description	根据ID更新角色信息
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"角色ID"
//	@Param			role	body		UpdateRoleRequest	true	"更新的角色信息"
//	@Success		204		{string}	string				"更新成功"
//	@Failure		400		{string}	string				"请求参数错误或ID格式错误"
//	@Failure		404		{string}	string				"角色不存在"
//	@Failure		500		{string}	string				"服务器内部错误"
//	@Router			/roles/{id} [put]
func UpdateRoleHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：" + err.Error()})
		return
	}

	roleCard := req.RoleData

	if req.AvatarUrl != roleCard.BasicInfo.AvatarURL && req.AvatarUrl != "" {
		roleCard.BasicInfo.AvatarURL = req.AvatarUrl
	}

	roleJSON, err := json.Marshal(&roleCard)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色卡格式错误：" + err.Error()})
		return
	}

	updatedRole := &model.Role{
		Name:        roleCard.BasicInfo.RoleName,
		WxUserId:    common.GetUserID(c),
		Description: getRoleDesc(&roleCard),
		AvatarUrl:   req.AvatarUrl,
	}
	updatedRole.RoleData = roleJSON

	if err := service.UpdateRole(uint(id), updatedRole); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "角色記錄不存在:" + err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败：" + err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "更新成功"})
}

// DeleteRoleHandler 删除角色接口
//
//	@Summary		删除角色
//	@Description	根据ID软删除角色
//	@Produce		json
//	@Param			id	path		int		true	"角色ID"
//	@Success		204	{string}	string	"删除成功"
//	@Failure		400	{string}	string	"ID格式错误"
//	@Failure		404	{string}	string	"角色不存在"
//	@Failure		500	{string}	string	"服务器内部错误"
//	@Router			/roles/{id} [delete]
func DeleteRoleHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	err = service.DeleteRole(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除失敗失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, "")
}

func decodeFile(input interface{}, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(input)
}
