package main

import (
	"fmt"
	"test-git/db"
	"test-git/handler"

	_ "test-git/docs" // 导入生成的 docs 包

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 1. 初始化数据库连接
	if err := db.Init(); err != nil {
		fmt.Printf("数据库初始化失败: %v\n", err)
		return
	}
	fmt.Println("数据库连接成功")

	// 2. 创建 Gin 引擎
	r := gin.Default() // 默认包含日志和恢复中间件

	// 注册 Swagger 路由（关键：让服务启动后能访问 Swagger 页面）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 3. 注册书籍相关路由
	bookGroup := r.Group("/books")
	{
		bookGroup.POST("", handler.CreateBookHandler)       // 创建书籍
		bookGroup.GET("/:id", handler.GetBookHandler)       // 查询单本书籍
		bookGroup.GET("", handler.ListBooksHandler)         // 查询书籍列表
		bookGroup.PUT("/:id", handler.UpdateBookHandler)    // 更新书籍
		bookGroup.DELETE("/:id", handler.DeleteBookHandler) // 删除书籍
	}

	// 4. 启动服务（默认端口8080）
	fmt.Println("服务启动成功，监听端口: 8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("服务启动失败: %v\n", err)
	}
}

func Add(a, b int) int {
	return a + b
}
