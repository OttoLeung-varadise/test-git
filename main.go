package main

import (
	"fmt"
	"test-git/common"
	"test-git/db"
	_ "test-git/docs"
	"test-git/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	if err := db.Init(); err != nil {
		fmt.Printf("database init fails: %v\n", err)
		return
	}
	fmt.Println("database connet succ")

	logDB, logErr := db.InitLogDB()
	if logErr != nil {
		fmt.Printf("log database init fails: %v\n", logErr)
	}

	r := gin.Default()

	// 注册 Swagger 路由（关键：让服务启动后能访问 Swagger 页面）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(common.HeaderMiddleware())

	if logErr == nil {
		// 啟動日志異步寫入協程
		go common.StartLogWriter(logDB)
		// 註冊請求日志中間件
		r.Use(common.RequestLogMiddleware())
	}

	roleGroup := r.Group("/roles")
	{
		roleGroup.GET("", handler.ListRoleHandler)           // 獲取角色列表
		roleGroup.GET("/:id", handler.GetRoleHandler)        // 查詢角色詳情
		roleGroup.POST("", handler.PreviewRoleHandler)       // 預覽角色卡
		roleGroup.POST("/create", handler.CreateRoleHandler) // 創建角色
		roleGroup.PUT("/:id", handler.UpdateRoleHandler)     // 更新角色
		roleGroup.DELETE("/:id", handler.DeleteRoleHandler)  // 刪除角色
	}

	fmt.Println("service started up, listen no port: 8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("service start fails: %v\n", err)
	}
}
