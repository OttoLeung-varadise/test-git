package common

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // 处理请求
		end := time.Now()

		// 构造日志字段（按需添加）
		logData := map[string]interface{}{
			"timestamp":  end.Format(time.RFC3339Nano), // 时间戳（带时区）
			"method":     c.Request.Method,             // 请求方法（GET/POST）
			"path":       c.Request.URL.Path,           // 请求路径
			"status":     c.Writer.Status(),            // 响应状态码
			"latency":    end.Sub(start).Seconds(),     // 响应时间（秒）
			"client_ip":  c.ClientIP(),                 // 客户端 IP
			"user_agent": c.Request.UserAgent(),        // 用户代理（可选）
		}

		// 输出 JSON 到控制台（Fluent Bit 会采集此日志）
		logJSON, _ := json.Marshal(logData)
		println(string(logJSON))
	}
}
