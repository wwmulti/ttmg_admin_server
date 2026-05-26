package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// RandString 生成随机字符串
func RandString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	rand.Seed(time.Now().UnixNano()) // 每次生成不同随机数
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// MapSlice 将多行的模型struct映射成map，只保留指定字段
func MapSlice[T any, R any](list []T, mapper func(T) R) []R {
	result := make([]R, 0, len(list))
	for _, v := range list {
		result = append(result, mapper(v))
	}
	return result
}

// MapSliceWithError 将多行的模型struct映射成map，只保留指定字段，支持返回error
func MapSliceWithError[T any, R any](list []T, mapper func(T) (R, error)) ([]R, error) {
	result := make([]R, 0, len(list))

	for _, v := range list {
		item, err := mapper(v)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

// FtoExInt64 将float64类型*1000转换成int64，保留3位小数
func FtoExInt64(v float64) int64 {
	return int64(v * 1000)
}

// ItoExFloat64 将int64类型/1000转换成float64，保留3位小数
func ItoExFloat64(v int64) float64 {
	return float64(v) / 1000
}

// GetTodayTimestamp 获取今天零点时间戳 比如现在是2026-01-26 10:10:10，那么返回2026-01-26 00:00:00的时间戳
func GetTodayTimestamp() int64 {
	loc := time.Now().Location()
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Unix()
	return todayZero
}

// GetMondayZeroTimestamp 获取指定时间戳的周一零点时间戳
func GetMondayZeroTimestamp(timestamp int64) int64 {
	t := time.Unix(timestamp, 0)

	// 计算距离周一的天数
	weekday := t.Weekday()
	daysToMonday := int(time.Monday - weekday)
	if daysToMonday <= 0 { // 如果是周一或之后，计算下周
		daysToMonday += 7
	}

	// 获取周一的零点
	monday := time.Date(t.Year(), t.Month(), t.Day()+daysToMonday, 0, 0, 0, 0, t.Location())
	return monday.Unix()
}

// GetFirstDayNextMonthZeroTimestamp 获取指定时间戳的次月1号零点时间戳
func GetFirstDayNextMonthZeroTimestamp(timestamp int64) int64 {
	t := time.Unix(timestamp, 0)

	// 计算下个月1号
	nextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return nextMonth.Unix()
}

// GetYesterdayTimestamp 获取昨天零点时间戳
func GetYesterdayTimestamp() int64 {
	loc := time.Now().Location()
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Unix()
	return todayZero - 86400
}

// GetTomorrowTimestamp 获取昨天零点时间戳
func GetTomorrowTimestamp() int64 {
	loc := time.Now().Location()
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Unix()
	return todayZero + 86400
}

// GetZeroTimeTimestamp 获取指定时间戳的零点时间戳
func GetZeroTimeTimestamp(timestamp int64) int64 {
	t := time.Unix(timestamp, 0).In(time.Local)
	zeroTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return zeroTime.Unix()
}

// GetNextZeroTimeTimestamp 获取指定时间戳的第二天零点时间戳
func GetNextZeroTimeTimestamp(timestamp int64) int64 {
	t := time.Unix(timestamp, 0).In(time.Local)
	zeroTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return zeroTime.Unix() + 86400
}

// GetNextMondayZeroTimeTimestamp 获取下周一零点时间戳
func GetNextMondayZeroTimeTimestamp(timestamp int64) int64 {
	t := time.Unix(timestamp, 0).In(time.Local)
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	diffDays := 8 - weekday
	zeroTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return zeroTime.Unix() + int64(diffDays)*86400
}

// GetNextMonthFirstDayZeroTimeTimestamp 获取下个月1号零点时间戳
func GetNextMonthFirstDayZeroTimeTimestamp(timestamp int64) int64 {
	t := time.Unix(timestamp, 0).In(time.Local)
	zeroTime := time.Date(
		t.Year(),
		t.Month()+1,
		1,
		0, 0, 0, 0,
		t.Location(),
	)
	return zeroTime.Unix()
}

// StringToTimestamp 将不同格式的时间字符串转换为 Unix 时间戳（秒）
// 支持格式：2006-01-02 15:04:05, 2006-01-02, 2006/01/02, 20060102
func StringToTimestamp(timeStr string) (int64, error) {
	timeStr = strings.TrimSpace(timeStr)
	if timeStr == "" {
		return 0, errors.New("empty time string")
	}

	var layout string
	// 根据长度和字符粗略判断格式
	switch len(timeStr) {
	case 19: // 2006-01-02 15:04:05
		if strings.Contains(timeStr, "/") {
			layout = "2006/01/02 15:04:05"
		} else {
			layout = "2006-01-02 15:04:05"
		}
	case 10: // 2006-01-02 或 2006/01/02
		if strings.Contains(timeStr, "/") {
			layout = "2006/01/02"
		} else {
			layout = "2006-01-02"
		}
	case 8: // 20060102
		layout = "20060102"
	default:
		return 0, errors.New("unsupported time format")
	}

	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

// GenerateImageHashFromReader 从 io.Reader（如上传的文件流）生成 SHA256
func GenerateImageHashFromReader(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	// 返回前 16 位通常就足够去重了，或者返回全量 64 位
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// IntSliceToString 将 []int 转换为指定分隔符的字符串
func IntSliceToString(ids []int, separator string) string {
	if len(ids) == 0 {
		return ""
	}

	strIds := make([]string, len(ids))
	for i, id := range ids {
		strIds[i] = strconv.Itoa(id)
	}
	return strings.Join(strIds, separator)
}

// 将字符串转换成[]int
func StringSliceToInt(str string, separator string) []int {
	idsStr := strings.Split(str, separator)

	// 转换为 []int
	idsInt := make([]int, 0, len(idsStr))
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			// 处理错误
			continue
		}
		idsInt = append(idsInt, id)
	}

	return idsInt
}

// 辅助函数：计算两个 int 切片的交集
func IntersectSlice(slice1, slice2 []int) []int {
	// 将 slice2 转换为 map 用于快速查找
	lookup := make(map[int]bool)
	for _, v := range slice2 {
		lookup[v] = true
	}

	// 找出交集
	result := make([]int, 0)
	for _, v := range slice1 {
		if lookup[v] {
			result = append(result, v)
		}
	}
	return result
}

// IsValidUrlFormat 是否是正确格式的地址
func IsValidUrlFormat(url string) (string, bool) {
	// 去除空格
	url = strings.TrimSpace(url)

	// 正则表达式验证
	// 支持: http://ip:port, https://domain, http://domain:port
	pattern := `^(http|https)://([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z]{2,}(:\d{1,5})?(/.*)?$|^(http|https)://((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(:\d{1,5})?(/.*)?$`

	matched, _ := regexp.MatchString(pattern, url)
	return url, matched
}

// UniqueInt int数组去重
func UniqueInt(arr []int) []int {
	result := make([]int, 0, len(arr))
	temp := make(map[int]struct{})

	for _, v := range arr {
		if _, ok := temp[v]; !ok {
			temp[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

func ConvertLangCode(code string) string {
	parts := strings.Split(code, "-")
	if len(parts) != 2 {
		return code
	}
	return parts[0] + "-" + strings.ToUpper(parts[1])
}

// ToJson 序列化为 JSON 字符串
func ToJson(v interface{}) string {
	if v == nil {
		return ""
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		// 游戏后端通常记录日志后返回空，避免程序崩溃
		return ""
	}
	return string(bytes)
}

// FromJson 反序列化 JSON 字符串为结构体
func FromJson[T any](str string) (T, error) {
	var target T
	if str == "" {
		return target, fmt.Errorf("json string is empty")
	}

	err := json.Unmarshal([]byte(str), &target)
	if err != nil {
		return target, err
	}
	return target, nil
}

// MD5Join Md5拼接
func MD5Join(parts ...interface{}) string {
	var builder strings.Builder

	for _, part := range parts {
		builder.WriteString(fmt.Sprintf("%v", part))
	}

	sum := md5.Sum([]byte(builder.String()))
	return hex.EncodeToString(sum[:])
}

// SignMap 返回：新map + token
func SignMap(data map[string]interface{}, secret string) (map[string]interface{}, string) {
	// 先算 token
	token := BuildToken(data, secret)

	// 拷贝原 map
	result := make(map[string]interface{}, len(data)+1)
	for k, v := range data {
		result[k] = v
	}

	// 加入 token
	result["token"] = token

	return result, token
}

// BuildToken 构建 token
func BuildToken(data map[string]interface{}, secret string) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		if k == "token" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteString("&")
		}
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(formatValue(data[k]))
	}
	builder.WriteString("&secret=")
	builder.WriteString(secret)
	sum := md5.Sum([]byte(builder.String()))
	return hex.EncodeToString(sum[:])
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case float64:
		// JSON数字统一按整数输出（你的时间戳场景）
		return strconv.FormatInt(int64(val), 10)
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}

// VerifyToken 校验 token
func VerifyToken(data map[string]interface{}, secret string) bool {
	token, ok := data["token"]
	if !ok {
		return false
	}

	expected := BuildToken(data, secret)
	return fmt.Sprintf("%v", token) == expected
}

// DivideAndRound 两个整数相除，保留指定位数小数（四舍五入）
// a: 分子, b: 分母, precision: 保留小数位数
func DivideAndRound(a, b int64, precision ...int) float64 {
	if b == 0 {
		return 0
	}

	// 设置默认精度
	pVal := 2
	if len(precision) > 0 {
		pVal = precision[0]
	}

	res := float64(a) / float64(b)

	// 动态计算 10 的 n 次方 (如精度 2 则是 100)
	p := math.Pow10(pVal)

	// 四舍五入
	return math.Round(res*p) / p
}

// SlicePage 泛型分页函数
// T 代表任意类型，list 可以是 []int, []string 或 []*models.DailyStat 等
func SlicePage[T any](page, limit int, list []T) []T {
	total := len(list)

	// 基本校验
	if total == 0 || limit <= 0 {
		return []T{}
	}

	// 计算起始位置
	start := (page - 1) * limit
	if start < 0 {
		start = 0
	}

	// 越界检查
	if start >= total {
		return []T{}
	}

	// 计算结束位置
	end := start + limit
	if end > total {
		end = total
	}

	return list[start:end]
}

type TimePair struct {
	Day           int
	RegisterStart int64
	RegisterEnd   int64
}

type TimeResult struct {
	LoginStart   int64
	LoginEnd     int64
	HistoryPairs []TimePair
}

// BuildTimePoints 生成留存时间点,返回指定时间昨天的开始结束时间和前天开始的留存时间数组
func BuildTimePoints(baseTimestamp int64, days []int) (int64, int64, []TimePair) {
	baseTime := time.Unix(baseTimestamp, 0)
	loc := baseTime.Location()
	todayStart := time.Date(
		baseTime.Year(),
		baseTime.Month(),
		baseTime.Day(),
		0, 0, 0, 0,
		loc,
	)

	yesterdayStart := todayStart.AddDate(0, 0, -1)
	pairs := make([]TimePair, 0, len(days))
	for _, d := range days {
		// 整体往前偏移 1 天
		start := todayStart.AddDate(0, 0, -(d + 1))
		end := todayStart.AddDate(0, 0, -d)

		pairs = append(pairs, TimePair{
			Day:           d,
			RegisterStart: start.Unix(),
			RegisterEnd:   end.Unix(),
		})
	}
	return yesterdayStart.Unix(), todayStart.Unix(), pairs
}

// StringToIntSlice 逗号字符串转数组
func StringToIntSlice(str string) []int {
	str = strings.TrimSpace(str)
	if str == "" {
		return []int{}
	}
	str = strings.Trim(str, ",")
	if str == "" {
		return []int{}
	}
	parts := strings.Split(str, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if num, err := strconv.Atoi(p); err == nil {
			result = append(result, num)
		}
	}
	return result
}

// InArray 判断元素是否在数组中
func InArray[T comparable](target T, arr []T) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

// GetNextDaysZeroTimestamps 返回从 start+1天 到 end+1天 的每天 00:00:00 时间戳，过滤掉大于今天的日期
func GetNextDaysZeroTimestamps(startStr string, endStr string) ([]int64, error) {
	layout := "2006-01-02 15:04:05"
	start, err := time.Parse(layout, startStr)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(layout, endStr)
	if err != nil {
		return nil, err
	}
	var result []int64
	cur := time.Date(
		start.Year(), start.Month(), start.Day(),
		0, 0, 0, 0,
		start.Location(),
	).AddDate(0, 0, 1)
	endLimit := time.Date(
		end.Year(), end.Month(), end.Day(),
		0, 0, 0, 0,
		end.Location(),
	).AddDate(0, 0, 1)
	now := time.Now()
	tomorrowStart := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		now.Location(),
	).AddDate(0, 0, 1)
	for !cur.After(endLimit) {
		// 保留 < 明天 的，也就是 今天及以前
		if cur.Before(tomorrowStart) {
			result = append(result, cur.Unix())
		}

		cur = cur.AddDate(0, 0, 1)
	}
	return result, nil
}

// BuildInCondition 通用 IN 查询构建器（多表union查询）
func BuildInCondition(field string, slice []int, conditions []string, tableArgs []interface{}) ([]string, []interface{}) {
	if len(slice) == 0 {
		return conditions, tableArgs
	}

	// 根据切片长度生成占位符 "?"
	placeholders := make([]string, len(slice))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	// 拼接成 "field IN (?,?,...)"
	condition := fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ","))
	conditions = append(conditions, condition)

	// 将切片值展开并追加到 tableArgs
	for _, v := range slice {
		tableArgs = append(tableArgs, v)
	}

	return conditions, tableArgs
}
