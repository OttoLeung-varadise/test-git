package common

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// 統計請求的耗時和內存使用量的中間件
func StartPicker() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 記錄請求開始時間（用於計算耗時）
		startTime := time.Now()

		// 2. 記錄請求開始時的內存狀態
		var startMem runtime.MemStats
		runtime.ReadMemStats(&startMem)

		// 處理請求（執行後續中間件和業務邏輯）
		c.Next()

		// 3. 請求結束後，計算耗時
		duration := time.Since(startTime)

		// 4. 記錄請求結束時的內存狀態，計算內存使用量差值
		var endMem runtime.MemStats
		runtime.ReadMemStats(&endMem)

		// 關注的內存指標（單位：字節）
		// 堆內存分配總量（累計分配的字節數，包括已釋放的）
		allocBytes := endMem.TotalAlloc - startMem.TotalAlloc
		// 當前堆內存使用量（未釋放的）
		heapInuseBytes := endMem.HeapInuse - startMem.HeapInuse

		// 5. 輸出統計結果（可根據需求調整格式，如寫入日志、打印到控制台等）
		c.JSON(200, gin.H{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"duration":    duration.String(),       // 耗時（字符串格式，如 "1.23ms"）
			"duration_ms": duration.Milliseconds(), // 耗時（毫秒，便於後續統計）
			"alloc_bytes": allocBytes,              // 本次請求期間分配的堆內存總量
			"heap_inuse":  heapInuseBytes,          // 本次請求結束後堆內存使用量（相對差值）
			"status_code": c.Writer.Status(),       // 響應狀態碼
		})

		// 若不需要在響應體中返回，可改為打印到日志
		// log.Printf(
		// 	"method=%s path=%s status=%d duration=%s alloc=%d bytes heap_inuse=%d bytes",
		// 	c.Request.Method, c.Request.URL.Path, c.Writer.Status(),
		// 	duration, allocBytes, heapInuseBytes,
		// )
	}
}
