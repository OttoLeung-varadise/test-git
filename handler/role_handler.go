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

func GetRoleHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetRoleHandler not implemented yet"})
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

func UpdateRoleHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateRoleHandler not implemented yet"})
}

func DeleteRoleHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DeleteRoleHandler not implemented yet"})
}

func decodeFile(input interface{}, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(input)
}
