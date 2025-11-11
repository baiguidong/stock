package model

// MatchResult 匹配结果
// 存储单个历史匹配的完整信息
type MatchResult struct {
	// 日期信息
	TargetDate    string `json:"target_date"`    // 目标日期（用于匹配的日期）
	HistoryDate   string `json:"history_date"`   // 历史匹配日期
	PredictDate   string `json:"predict_date"`   // 预测参考日期（历史日的次日）

	// 偏移信息
	Offset        int    `json:"offset"`         // 平移偏移量

	// 统计信息
	PositiveCount int     `json:"positive_count"` // 正向匹配数量
	NegativeCount int     `json:"negative_count"` // 反向匹配数量
	Similarity    float64 `json:"similarity"`     // 相似度分数

	// 走势信息
	TargetTrend   TrendType `json:"target_trend"`   // 目标日走势类型
	HistoryTrend  TrendType `json:"history_trend"`  // 历史日走势类型
	IsQualified   bool      `json:"is_qualified"`   // 是否合格
	IsReversed    bool      `json:"is_reversed"`    // 是否反转

	// 图表数据
	TargetData    []Point `json:"target_data"`    // 目标日数据点
	HistoryData   []Point `json:"history_data"`   // 历史日数据点
	PredictData   []Point `json:"predict_data"`   // 预测日数据点
}

// Point 图表数据点
type Point struct {
	X float64 `json:"x"` // X坐标（分钟索引）
	Y float64 `json:"y"` // Y坐标（价格）
}

// PredictRequest 预测请求参数
type PredictRequest struct {
	Date          string `json:"date"`           // 预测日期
	OffsetStart   int    `json:"offset_start"`   // 平移起始值
	OffsetEnd     int    `json:"offset_end"`     // 平移结束值
	ResultCount   int    `json:"result_count"`   // 输出结果数量
	CompareLength int    `json:"compare_length"` // 对比长度（分钟数）
}

// PredictResponse 预测响应结果
type PredictResponse struct {
	Success      bool           `json:"success"`       // 是否成功
	Message      string         `json:"message"`       // 消息
	Date         string         `json:"date"`          // 预测日期
	TotalMatches int            `json:"total_matches"` // 总匹配数
	Qualified    int            `json:"qualified"`     // 合格数量
	Unqualified  int            `json:"unqualified"`   // 不合格数量
	Results      []MatchResult  `json:"results"`       // 匹配结果列表
	ProcessTime  float64        `json:"process_time"`  // 处理时间（秒）
}

// FeatureSet 特征值集合
// 存储开高低收的特征值（1-8）
type FeatureSet struct {
	Open  int `json:"open"`  // 开盘特征
	High  int `json:"high"`  // 最高特征
	Low   int `json:"low"`   // 最低特征
	Close int `json:"close"` // 收盘特征
}

// String 返回特征集的字符串表示
func (f FeatureSet) String() string {
	return fmt.Sprintf("[开:%d 高:%d 低:%d 收:%d]", f.Open, f.High, f.Low, f.Close)
}

// Equals 判断两个特征集是否相等
func (f FeatureSet) Equals(other FeatureSet) bool {
	return f.Open == other.Open &&
		f.High == other.High &&
		f.Low == other.Low &&
		f.Close == other.Close
}

// MatchOpen 判断开盘特征是否匹配
func (f FeatureSet) MatchOpen(other FeatureSet) bool {
	return f.Open == other.Open
}

// MatchClose 判断收盘特征是否匹配
func (f FeatureSet) MatchClose(other FeatureSet) bool {
	return f.Close == other.Close
}

// MatchOpenClose 判断开盘和收盘特征是否都匹配
func (f FeatureSet) MatchOpenClose(other FeatureSet) bool {
	return f.MatchOpen(other) && f.MatchClose(other)
}

// DataStatistics 数据统计信息
type DataStatistics struct {
	DailyDataCount    int      `json:"daily_data_count"`    // 日线数据数量
	MinuteDataCount   int      `json:"minute_data_count"`   // 分钟数据数量
	DateRange         []string `json:"date_range"`          // 数据日期范围 [start, end]
	LoadTime          float64  `json:"load_time"`           // 加载时间（秒）
}

// SystemInfo 系统信息
type SystemInfo struct {
	Version       string          `json:"version"`        // 版本号
	DataPath      string          `json:"data_path"`      // 数据路径
	Statistics    DataStatistics  `json:"statistics"`     // 数据统计
	ServerStatus  string          `json:"server_status"`  // 服务器状态
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Success bool   `json:"success"` // 总是false
	Error   string `json:"error"`   // 错误信息
}

// 需要导入fmt包
import "fmt"
