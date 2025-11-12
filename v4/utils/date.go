package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// FormatDate 格式化日期字符串
//
// 功能说明：
//   将各种格式的日期字符串统一转换为标准格式: YYYY-MM-DD
//
// 支持的输入格式：
//   - 2018/11/1 -> 2018-11-01
//   - 2018/1/1  -> 2018-01-01
//   - 2018-11-01 (保持不变)
//
// 参数：
//   dateStr - 输入的日期字符串
//
// 返回：
//   string - 格式化后的日期字符串 (YYYY-MM-DD)
func FormatDate(dateStr string) string {
	// 如果包含斜杠,则转换
	if strings.Contains(dateStr, "/") {
		parts := strings.Split(dateStr, "/")
		if len(parts) != 3 {
			return dateStr
		}

		year, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		day, _ := strconv.Atoi(parts[2])

		return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	}

	return dateStr
}

// NextTradingDay 获取下一个交易日
//
// 功能说明：
//   跳过周末,返回下一个工作日
//   - 周五 -> 下周一 (+3天)
//   - 周六 -> 下周一 (+2天)
//   - 其他 -> 次日 (+1天)
//
// 参数：
//   dateStr - 当前日期字符串 (YYYY-MM-DD)
//
// 返回：
//   string - 下一个交易日 (YYYY-MM-DD)
func NextTradingDay(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err.Error()
	}

	// 根据星期几决定跳过的天数
	switch t.Weekday() {
	case time.Friday:
		// 周五 -> 下周一
		t = t.AddDate(0, 0, 3)
	case time.Saturday:
		// 周六 -> 下周一
		t = t.AddDate(0, 0, 2)
	default:
		// 其他 -> 次日
		t = t.AddDate(0, 0, 1)
	}

	return t.Format("2006-01-02")
}

// PreviousTradingDay 获取前一个交易日
//
// 功能说明：
//   跳过周末,返回前一个工作日
//   - 周一 -> 上周五 (-3天)
//   - 其他 -> 前一天 (-1天)
//
// 参数：
//   dateStr - 当前日期字符串 (YYYY-MM-DD)
//
// 返回：
//   string - 前一个交易日 (YYYY-MM-DD)
func PreviousTradingDay(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err.Error()
	}

	// 根据星期几决定回退的天数
	if t.Weekday() == time.Monday {
		// 周一 -> 上周五
		t = t.AddDate(0, 0, -3)
	} else {
		// 其他 -> 前一天
		t = t.AddDate(0, 0, -1)
	}

	return t.Format("2006-01-02")
}

// ParseDate 解析日期字符串
//
// 功能说明：
//   尝试多种日期格式进行解析
//
// 支持的格式：
//   - 2006-01-02
//   - 2006/01/02
//   - 2006/1/2
//   - 2006-1-2
//
// 参数：
//   dateStr - 日期字符串
//
// 返回：
//   time.Time - 解析后的时间对象
//   error - 错误信息
func ParseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"2006/1/2",
		"2006-1-2",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析日期: %s", dateStr)
}

// GetHomeDir 获取用户主目录
//
// 返回：
//   string - 用户主目录路径
func GetHomeDir() string {
	// 在Linux/Mac下使用环境变量
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	// 在Windows下使用USERPROFILE
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home
	}

	// 默认返回空字符串
	return ""
}
