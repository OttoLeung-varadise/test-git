package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// RequestLog 請求日志Model，對應數據庫表request_logs
type RequestLog struct {
	ID              uint64         `gorm:"column:id;type:serial;primaryKey" json:"id"`                                        // 主鍵，自增
	RequestID       string         `gorm:"column:request_id;type:varchar(64);not null;index" json:"request_id"`               // 請求唯一標識
	Method          string         `gorm:"column:method;type:varchar(10);not null" json:"method"`                             // HTTP方法
	Path            string         `gorm:"column:path;type:varchar(255);not null;index" json:"path"`                          // 請求路徑
	QueryString     string         `gorm:"column:query_string;type:text" json:"query_string"`                                 // 查詢參數
	StatusCode      int            `gorm:"column:status_code;not null" json:"status_code"`                                    // 響應狀態碼
	RemoteIP        string         `gorm:"column:remote_ip;type:varchar(45);not null" json:"remote_ip"`                       // 客戶端IP
	UserAgent       string         `gorm:"column:user_agent;type:text" json:"user_agent"`                                     // 用戶代理
	RequestTime     float64        `gorm:"column:request_time;not null" json:"request_time"`                                  // 請求耗時（秒）
	CreatedAt       time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:now();index" json:"created_at"` // 日志創建時間（帶時區）
	FileName        string         `gorm:"column:file_name;type:varchar(255)" json:"file_name"`                               // 上傳文件名
	FileSize        int64          `gorm:"column:file_size" json:"file_size"`                                                 // 文件大小（字節）
	FileContentJSON JSONRawMessage `gorm:"column:file_content_json;type:jsonb" json:"file_content_json"`                      // 文件内容（JSON格式）
}

func (RequestLog) TableName() string {
	return "request_logs"
}

type JSONRawMessage json.RawMessage

// Value 實現driver.Valuer接口，將JSONRawMessage轉為驅動可識別的類型
func (j JSONRawMessage) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return string(j), nil
}

// Scan 實現sql.Scanner接口，將數據庫返回的JSON字符串轉為JSONRawMessage
func (j *JSONRawMessage) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

func AutoMigrateRequestLog(db *gorm.DB) error {
	return db.AutoMigrate(&RequestLog{})
}
