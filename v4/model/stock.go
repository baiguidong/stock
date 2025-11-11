package model

import "time"

// StockData 单根K线数据
// 包含单根K线的完整信息：日期、时间、开高低收价格以及高低点标记
type StockData struct {
	Date  string  `json:"date"`  // 交易日期，格式：YYYY-MM-DD
	Time  string  `json:"time"`  // 交易时间，格式：HH:MM
	Open  float64 `json:"open"`  // 开盘价
	High  float64 `json:"high"`  // 最高价
	Low   float64 `json:"low"`   // 最低价
	Close float64 `json:"close"` // 收盘价

	IsHighPoint bool `json:"is_high_point"` // 是否为高点（比前后都高）
	IsLowPoint  bool `json:"is_low_point"`  // 是否为低点（比前后都低）
}

// StockDataSet K线数据集合
// 包含一天的完整分钟数据及相关统计信息
type StockDataSet struct {
	Data map[int]*StockData `json:"data"` // 分钟索引 -> 数据（1-240）

	MorningHigh   float64 `json:"morning_high"`   // 上午最高价
	MorningLow    float64 `json:"morning_low"`    // 上午最低价
	AfternoonHigh float64 `json:"afternoon_high"` // 下午最高价
	AfternoonLow  float64 `json:"afternoon_low"`  // 下午最低价

	PositiveCount int `json:"positive_count"` // 正向匹配数量
	NegativeCount int `json:"negative_count"` // 反向匹配数量
}

// TrendType 走势类型
type TrendType int

const (
	TrendUp   TrendType = 1 // 上涨
	TrendDown TrendType = 2 // 下跌
	TrendFlat TrendType = 3 // 震荡
)

// String 返回走势类型的字符串表示
func (t TrendType) String() string {
	switch t {
	case TrendUp:
		return "上涨"
	case TrendDown:
		return "下跌"
	case TrendFlat:
		return "震荡"
	default:
		return "未知"
	}
}

// InitData 初始化数据集合
func (s *StockDataSet) InitData() {
	s.Data = make(map[int]*StockData)
}

// Add 添加单条K线数据
func (s *StockDataSet) Add(data *StockData, index int) {
	if s.Data == nil {
		s.InitData()
	}
	s.Data[index] = data
}

// GetTrend 判断走势类型
// 规则：
//   - 上午以11:30（第120分钟）为分界线
//   - 上涨：下午最高>上午最高 且 下午最低>上午最低
//   - 下跌：下午最高<上午最高 且 下午最低<上午最低
//   - 震荡：其他情况
func (s *StockDataSet) GetTrend() TrendType {
	if s.AfternoonHigh > s.MorningHigh && s.AfternoonLow > s.MorningLow {
		return TrendUp
	}
	if s.AfternoonHigh < s.MorningHigh && s.AfternoonLow < s.MorningLow {
		return TrendDown
	}
	return TrendFlat
}

// CalculateTrendData 计算涨跌震荡数据
// 分析上午下午的最高最低价，用于判断走势类型
func (s *StockDataSet) CalculateTrendData() {
	// 检查数据完整性
	if !s.isDataComplete() {
		return
	}

	// 初始化上午数据（第1分钟）
	s.MorningHigh = s.Data[1].Close
	s.MorningLow = s.Data[1].Close

	// 初始化下午数据（第121分钟）
	s.AfternoonHigh = s.Data[121].Close
	s.AfternoonLow = s.Data[121].Close

	// 计算上午最高最低（1-120分钟）
	for i := 1; i <= 120; i++ {
		if data, ok := s.Data[i]; ok {
			if data.Close > s.MorningHigh {
				s.MorningHigh = data.Close
			}
			if data.Close < s.MorningLow {
				s.MorningLow = data.Close
			}
		}
	}

	// 计算下午最高最低（121-240分钟）
	for i := 121; i <= 240; i++ {
		if data, ok := s.Data[i]; ok {
			if data.Close > s.AfternoonHigh {
				s.AfternoonHigh = data.Close
			}
			if data.Close < s.AfternoonLow {
				s.AfternoonLow = data.Close
			}
		}
	}
}

// IdentifyHighLowPoints 识别高低点
// 规则：
//   - 高点：比前一分钟高 且 比后一分钟高
//   - 低点：比前一分钟低 且 比后一分钟低
//   - 返回：高点数量，低点数量
func (s *StockDataSet) IdentifyHighLowPoints() (highCount, lowCount int) {
	// 检查数据完整性
	if !s.isDataComplete() {
		return 0, 0
	}

	// 处理第一个点
	if s.Data[1].Close > s.Data[2].Close {
		s.Data[1].IsHighPoint = true
		highCount++
	} else {
		s.Data[1].IsLowPoint = true
		lowCount++
	}

	// 处理最后一个点
	if s.Data[240].Close > s.Data[239].Close {
		s.Data[240].IsHighPoint = true
		highCount++
	} else {
		s.Data[240].IsLowPoint = true
		lowCount++
	}

	// 处理中间的点
	for i := 2; i < 240; i++ {
		prev := s.Data[i-1]
		curr := s.Data[i]
		next := s.Data[i+1]

		if prev == nil || curr == nil || next == nil {
			continue
		}

		// 判断高点
		if curr.Close > prev.Close && curr.Close > next.Close {
			curr.IsHighPoint = true
			highCount++
		}

		// 判断低点
		if curr.Close < prev.Close && curr.Close < next.Close {
			curr.IsLowPoint = true
			lowCount++
		}
	}

	return highCount, lowCount
}

// GetOHLC 获取开高低收价格
// 返回：开盘价、最高价、最低价、收盘价
func (s *StockDataSet) GetOHLC() (open, high, low, close float64) {
	if !s.isDataComplete() {
		return 0, 0, 0, 0
	}

	open = s.Data[1].Open
	close = s.Data[240].Close
	high = s.Data[1].High
	low = s.Data[1].Low

	// 遍历所有数据找最高最低
	for _, data := range s.Data {
		if data.High > high {
			high = data.High
		}
		if data.Low < low {
			low = data.Low
		}
	}

	return open, high, low, close
}

// isDataComplete 检查数据是否完整（240个数据点）
func (s *StockDataSet) isDataComplete() bool {
	if len(s.Data) != 240 {
		return false
	}
	for i := 1; i <= 240; i++ {
		if s.Data[i] == nil {
			return false
		}
	}
	return true
}

// DailyData 日线数据
type DailyData struct {
	Date  string  `json:"date"`  // 日期
	Open  float64 `json:"open"`  // 开盘价
	High  float64 `json:"high"`  // 最高价
	Low   float64 `json:"low"`   // 最低价
	Close float64 `json:"close"` // 收盘价
}

// TimeIndex 时间索引映射
// 将时间字符串映射到分钟索引（1-240）
var TimeIndexMap = map[string]int{
	"9:31": 1, "9:32": 2, "9:33": 3, "9:34": 4, "9:35": 5,
	"9:36": 6, "9:37": 7, "9:38": 8, "9:39": 9, "9:40": 10,
	"9:41": 11, "9:42": 12, "9:43": 13, "9:44": 14, "9:45": 15,
	"9:46": 16, "9:47": 17, "9:48": 18, "9:49": 19, "9:50": 20,
	"9:51": 21, "9:52": 22, "9:53": 23, "9:54": 24, "9:55": 25,
	"9:56": 26, "9:57": 27, "9:58": 28, "9:59": 29,
	"09:31": 1, "09:32": 2, "09:33": 3, "09:34": 4, "09:35": 5,
	"09:36": 6, "09:37": 7, "09:38": 8, "09:39": 9, "09:40": 10,
	"09:41": 11, "09:42": 12, "09:43": 13, "09:44": 14, "09:45": 15,
	"09:46": 16, "09:47": 17, "09:48": 18, "09:49": 19, "09:50": 20,
	"09:51": 21, "09:52": 22, "09:53": 23, "09:54": 24, "09:55": 25,
	"09:56": 26, "09:57": 27, "09:58": 28, "09:59": 29,
	"10:00": 30, "10:01": 31, "10:02": 32, "10:03": 33, "10:04": 34,
	"10:05": 35, "10:06": 36, "10:07": 37, "10:08": 38, "10:09": 39,
	"10:10": 40, "10:11": 41, "10:12": 42, "10:13": 43, "10:14": 44,
	"10:15": 45, "10:16": 46, "10:17": 47, "10:18": 48, "10:19": 49,
	"10:20": 50, "10:21": 51, "10:22": 52, "10:23": 53, "10:24": 54,
	"10:25": 55, "10:26": 56, "10:27": 57, "10:28": 58, "10:29": 59,
	"10:30": 60, "10:31": 61, "10:32": 62, "10:33": 63, "10:34": 64,
	"10:35": 65, "10:36": 66, "10:37": 67, "10:38": 68, "10:39": 69,
	"10:40": 70, "10:41": 71, "10:42": 72, "10:43": 73, "10:44": 74,
	"10:45": 75, "10:46": 76, "10:47": 77, "10:48": 78, "10:49": 79,
	"10:50": 80, "10:51": 81, "10:52": 82, "10:53": 83, "10:54": 84,
	"10:55": 85, "10:56": 86, "10:57": 87, "10:58": 88, "10:59": 89,
	"11:00": 90, "11:01": 91, "11:02": 92, "11:03": 93, "11:04": 94,
	"11:05": 95, "11:06": 96, "11:07": 97, "11:08": 98, "11:09": 99,
	"11:10": 100, "11:11": 101, "11:12": 102, "11:13": 103, "11:14": 104,
	"11:15": 105, "11:16": 106, "11:17": 107, "11:18": 108, "11:19": 109,
	"11:20": 110, "11:21": 111, "11:22": 112, "11:23": 113, "11:24": 114,
	"11:25": 115, "11:26": 116, "11:27": 117, "11:28": 118, "11:29": 119,
	"11:30": 120,
	"13:00": 120, "13:01": 121, "13:02": 122, "13:03": 123, "13:04": 124,
	"13:05": 125, "13:06": 126, "13:07": 127, "13:08": 128, "13:09": 129,
	"13:10": 130, "13:11": 131, "13:12": 132, "13:13": 133, "13:14": 134,
	"13:15": 135, "13:16": 136, "13:17": 137, "13:18": 138, "13:19": 139,
	"13:20": 140, "13:21": 141, "13:22": 142, "13:23": 143, "13:24": 144,
	"13:25": 145, "13:26": 146, "13:27": 147, "13:28": 148, "13:29": 149,
	"13:30": 150, "13:31": 151, "13:32": 152, "13:33": 153, "13:34": 154,
	"13:35": 155, "13:36": 156, "13:37": 157, "13:38": 158, "13:39": 159,
	"13:40": 160, "13:41": 161, "13:42": 162, "13:43": 163, "13:44": 164,
	"13:45": 165, "13:46": 166, "13:47": 167, "13:48": 168, "13:49": 169,
	"13:50": 170, "13:51": 171, "13:52": 172, "13:53": 173, "13:54": 174,
	"13:55": 175, "13:56": 176, "13:57": 177, "13:58": 178, "13:59": 179,
	"14:00": 180, "14:01": 181, "14:02": 182, "14:03": 183, "14:04": 184,
	"14:05": 185, "14:06": 186, "14:07": 187, "14:08": 188, "14:09": 189,
	"14:10": 190, "14:11": 191, "14:12": 192, "14:13": 193, "14:14": 194,
	"14:15": 195, "14:16": 196, "14:17": 197, "14:18": 198, "14:19": 199,
	"14:20": 200, "14:21": 201, "14:22": 202, "14:23": 203, "14:24": 204,
	"14:25": 205, "14:26": 206, "14:27": 207, "14:28": 208, "14:29": 209,
	"14:30": 210, "14:31": 211, "14:32": 212, "14:33": 213, "14:34": 214,
	"14:35": 215, "14:36": 216, "14:37": 217, "14:38": 218, "14:39": 219,
	"14:40": 220, "14:41": 221, "14:42": 222, "14:43": 223, "14:44": 224,
	"14:45": 225, "14:46": 226, "14:47": 227, "14:48": 228, "14:49": 229,
	"14:50": 230, "14:51": 231, "14:52": 232, "14:53": 233, "14:54": 234,
	"14:55": 235, "14:56": 236, "14:57": 237, "14:58": 238, "14:59": 239,
	"15:00": 240,
}

// GetTimeIndex 获取时间对应的分钟索引
func GetTimeIndex(timeStr string) int {
	if index, ok := TimeIndexMap[timeStr]; ok {
		return index
	}
	return 0
}

// ParseDate 解析日期字符串
// 支持格式: 2018/11/1 -> 2018-11-01
func ParseDate(dateStr string) (time.Time, error) {
	// 尝试多种格式
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

	return time.Time{}, nil
}

// FormatDate 格式化日期为标准格式 YYYY-MM-DD
func FormatDate(dateStr string) string {
	t, err := ParseDate(dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("2006-01-02")
}
