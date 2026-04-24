package utils

import (
	"fmt"
	"time"
)

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}

// Paginate 生成分页 SQL 语句
func Paginate(page, pageSize int) (offset, limit int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	limit = pageSize
	offset = (page - 1) * pageSize
	return offset, limit
}

// GenerateRandomString 生成随机字符串（简单实现）
func GenerateRandomString(length int) string {
	// TODO: 实现更安全的随机字符串生成
	return fmt.Sprintf("random_%d", time.Now().UnixNano())
}
