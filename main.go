package main

import (
	"fmt"
	"test-git/db"
	"test-git/handler"

	_ "test-git/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 初始化数据库连接 & migration
	if err := db.Init(); err != nil {
		fmt.Printf("database init fails: %v\n", err)
		return
	}
	fmt.Println("database connet succ")

	r := gin.Default()

	// 注册 Swagger 路由（关键：让服务启动后能访问 Swagger 页面）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	bookGroup := r.Group("/books")
	{
		bookGroup.POST("", handler.CreateBookHandler)       // 创建书籍
		bookGroup.GET("/:id", handler.GetBookHandler)       // 查询单本书籍
		bookGroup.GET("", handler.ListBooksHandler)         // 查询书籍列表
		bookGroup.PUT("/:id", handler.UpdateBookHandler)    // 更新书籍
		bookGroup.DELETE("/:id", handler.DeleteBookHandler) // 删除书籍
	}

	fmt.Println("service started up, listen no port: 8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("service start fails: %v\n", err)
	}
}
