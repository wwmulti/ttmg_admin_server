package utils

import (
	"strconv"
	"sync"
)

// GlobalCache 全局并发安全缓存
var GlobalCache sync.Map

// Set 存入任意类型的数据 (string, int, []string, map 等)
func Set(key string, value interface{}) {
	GlobalCache.Store(key, value)
}

// GetItem 泛型获取方法。用法示例：val := GetItem[[]string]("key")
func GetItem[T any](key string) T {
	var zero T
	val, ok := GlobalCache.Load(key)
	if !ok {
		return zero
	}
	// 尝试断言类型
	if res, ok := val.(T); ok {
		return res
	}
	return zero
}

// GetString 快捷获取字符串
func GetString(key string) string {
	return GetItem[string](key)
}

// GetInt64 快捷获取 int64
func GetInt64(key string) int64 {
	return GetItem[int64](key)
}

// GetInt 强制转换获取 int
func GetInt(key string) int {
	val, ok := GlobalCache.Load(key)
	if !ok {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case string:
		i, _ := strconv.Atoi(v)
		return i
	default:
		return 0
	}
}

// GetFloat64 强制转换获取 float64
func GetFloat64(key string) float64 {
	val, ok := GlobalCache.Load(key)
	if !ok {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

//	growthMin := utils.GetInt(config.System.GrowthMin) // 取值
//  utils.Set(key, defaultValue) // 存值
