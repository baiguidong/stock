package algorithm

import (
	"fmt"
	"sort"
	"stock_predict/model"
)

// HistoricalMatcher 历史数据匹配器
// 用于根据特征值匹配历史数据
type HistoricalMatcher struct {
	FeatureExtractor *FeatureExtractor // 特征提取器
	HighLowDetector  *HighLowDetector  // 高低点检测器
}

// NewHistoricalMatcher 创建历史数据匹配器
func NewHistoricalMatcher(featureExtractor *FeatureExtractor, highLowDetector *HighLowDetector) *HistoricalMatcher {
	return &HistoricalMatcher{
		FeatureExtractor: featureExtractor,
		HighLowDetector:  highLowDetector,
	}
}

// MatchResult 单个匹配结果
type MatchResult struct {
	Date           string               // 匹配日期
	NextDate       string               // 下一个交易日
	CurrentData    *model.StockDataSet  // 当天数据
	NextData       *model.StockDataSet  // 下一天数据
	PositiveCount  int                  // 正向匹配数（高点对高点，低点对低点）
	NegativeCount  int                  // 反向匹配数（高点对低点）
	Similarity     float64              // 相似度
	IsReversed     bool                 // 是否反转
	TrendType      model.TrendType      // 走势类型
}

// MatchResults 匹配结果集合
type MatchResults struct {
	Results        []MatchResult
	PositiveTotal  int // 总正向匹配数
	NegativeTotal  int // 总反向匹配数
}

// FindMatches 查找匹配的历史数据
//
// 算法说明：
//   1. 提取目标日期的特征值（开、高、低、收）
//   2. 遍历历史数据，找出开盘和收盘特征相同的日期
//   3. 进一步筛选高或低特征也相同的日期
//   4. 对匹配的日期，计算分钟级高低点匹配度
//   5. 判断趋势一致性
//   6. 如果反向匹配多于正向匹配，则反转数据
//
// 参数：
//   targetDate    - 目标日期
//   targetDataSet - 目标日期的分钟数据
//   dailyData     - 历史日线数据
//   minuteDataMap - 历史分钟数据映射 (日期 -> 分钟数据)
//
// 返回：
//   *MatchResults - 匹配结果集合
//   error         - 错误信息
func (hm *HistoricalMatcher) FindMatches(
	targetDate string,
	targetDataSet *model.StockDataSet,
	dailyData []model.DailyData,
	minuteDataMap map[string]*model.StockDataSet,
) (*MatchResults, error) {

	// 计算目标日期的高低点
	hm.HighLowDetector.DetectHighLowPoints(targetDataSet)
	hm.HighLowDetector.CalculateSessionHighLow(targetDataSet)

	// 提取目标日期的OHLC
	open, high, low, close := hm.HighLowDetector.GetOHLC(targetDataSet)

	// 提取目标日期的特征值
	targetFeature := hm.FeatureExtractor.ExtractFromOHLC(open, high, low, close)

	fmt.Printf("目标日期 %s 的特征: [开:%d, 高:%d, 低:%d, 收:%d]\n",
		targetDate, targetFeature.Open, targetFeature.High, targetFeature.Low, targetFeature.Close)

	results := &MatchResults{
		Results: []MatchResult{},
	}

	positiveTotal := 0
	negativeTotal := 0

	// 遍历历史日线数据
	for i, daily := range dailyData {
		// 跳过目标日期本身
		if daily.Date == targetDate {
			continue
		}

		// 提取历史日期的特征值
		histFeature := hm.FeatureExtractor.ExtractFromOHLC(
			daily.Open, daily.High, daily.Low, daily.Close)

		// 第一步筛选：开盘和收盘特征必须相同
		if histFeature.Open != targetFeature.Open || histFeature.Close != targetFeature.Close {
			continue
		}

		// 第二步筛选：高或低特征至少有一个相同
		if histFeature.High != targetFeature.High && histFeature.Low != targetFeature.Low {
			continue
		}

		// 获取该日期的分钟数据
		histDataSet := minuteDataMap[daily.Date]
		if histDataSet == nil || histDataSet.Data == nil {
			continue
		}

		// 检查数据完整性
		if !hm.HighLowDetector.isDataComplete(histDataSet) {
			continue
		}

		// 计算分钟数据的高低点
		hm.HighLowDetector.DetectHighLowPoints(histDataSet)
		hm.HighLowDetector.CalculateSessionHighLow(histDataSet)

		// 比较高低点匹配度
		positiveCount, negativeCount := hm.compareHighLowPoints(targetDataSet, histDataSet)

		// 判断趋势
		targetTrend := hm.HighLowDetector.GetTrendType(targetDataSet)
		histTrend := hm.HighLowDetector.GetTrendType(histDataSet)

		// 趋势一致性判断
		trendMatch := hm.isTrendMatch(targetTrend, histTrend)

		// 如果趋势不匹配且不是震荡，则跳过
		if !trendMatch && targetTrend != model.TrendFlat && histTrend != model.TrendFlat {
			continue
		}

		// 获取下一个交易日
		nextDate := ""
		var nextData *model.StockDataSet
		if i+1 < len(dailyData) {
			nextDate = dailyData[i+1].Date
			nextData = minuteDataMap[nextDate]
		}

		// 判断是否需要反转
		isReversed := negativeCount > positiveCount

		// 创建匹配结果
		result := MatchResult{
			Date:          daily.Date,
			NextDate:      nextDate,
			CurrentData:   histDataSet,
			NextData:      nextData,
			PositiveCount: positiveCount,
			NegativeCount: negativeCount,
			IsReversed:    isReversed,
			TrendType:     histTrend,
		}

		results.Results = append(results.Results, result)

		if isReversed {
			negativeTotal++
		} else {
			positiveTotal++
		}
	}

	results.PositiveTotal = positiveTotal
	results.NegativeTotal = negativeTotal

	fmt.Printf("找到 %d 个匹配日期，正向 %d 个，反向 %d 个\n",
		len(results.Results), positiveTotal, negativeTotal)

	return results, nil
}

// compareHighLowPoints 比较两个数据集的高低点匹配度
//
// 算法说明：
//   遍历每一分钟的数据，比较高低点标记
//   - 正向匹配：高点对高点，低点对低点
//   - 反向匹配：高点对低点，低点对高点
//
// 参数：
//   target  - 目标数据集
//   history - 历史数据集
//
// 返回：
//   positiveCount - 正向匹配数
//   negativeCount - 反向匹配数
func (hm *HistoricalMatcher) compareHighLowPoints(
	target *model.StockDataSet,
	history *model.StockDataSet,
) (int, int) {

	positiveCount := 0
	negativeCount := 0

	// 遍历每一分钟
	for i := 1; i <= hm.HighLowDetector.TradingMinutes; i++ {
		targetData := target.Data[i]
		histData := history.Data[i]

		if targetData == nil || histData == nil {
			continue
		}

		// 如果目标数据是高点
		if targetData.IsHighPoint {
			if histData.IsHighPoint {
				positiveCount++ // 高点对高点
			} else if histData.IsLowPoint {
				negativeCount++ // 高点对低点
			}
		}

		// 如果目标数据是低点
		if targetData.IsLowPoint {
			if histData.IsLowPoint {
				positiveCount++ // 低点对低点
			} else if histData.IsHighPoint {
				negativeCount++ // 低点对高点
			}
		}
	}

	return positiveCount, negativeCount
}

// isTrendMatch 判断两个趋势是否匹配
//
// 参数：
//   trend1 - 趋势1
//   trend2 - 趋势2
//
// 返回：
//   bool - 是否匹配
func (hm *HistoricalMatcher) isTrendMatch(trend1, trend2 model.TrendType) bool {
	// 完全相同
	if trend1 == trend2 {
		return true
	}

	// 如果有一个是震荡，则认为匹配
	if trend1 == model.TrendFlat || trend2 == model.TrendFlat {
		return true
	}

	return false
}

// SortByPositiveCount 按正向匹配数排序
func (mr *MatchResults) SortByPositiveCount() {
	sort.Slice(mr.Results, func(i, j int) bool {
		return mr.Results[i].PositiveCount > mr.Results[j].PositiveCount
	})
}

// SortBySimilarity 按相似度排序
func (mr *MatchResults) SortBySimilarity() {
	sort.Slice(mr.Results, func(i, j int) bool {
		return mr.Results[i].Similarity > mr.Results[j].Similarity
	})
}

// GetTopN 获取前N个结果
//
// 参数：
//   n - 数量
//
// 返回：
//   []MatchResult - 前N个结果
func (mr *MatchResults) GetTopN(n int) []MatchResult {
	if n > len(mr.Results) {
		n = len(mr.Results)
	}
	return mr.Results[:n]
}

// ReverseData 反转数据（镜像）
//
// 算法说明：
//   当反向匹配数大于正向匹配数时，需要反转数据
//   反转方法：将曲线上下翻转
//
// 参数：
//   dataSet - 数据集
//   pivot   - 翻转轴（中心价格）
//
// 返回：
//   *model.StockDataSet - 反转后的数据集
func ReverseData(dataSet *model.StockDataSet, pivot float64) *model.StockDataSet {
	reversed := &model.StockDataSet{
		Data: make(map[int]*model.StockData),
	}

	for i, data := range dataSet.Data {
		if data == nil {
			continue
		}

		// 翻转价格：新价格 = 2 * pivot - 原价格
		reversedData := &model.StockData{
			Date:  data.Date,
			Time:  data.Time,
			Open:  2*pivot - data.Open,
			High:  2*pivot - data.Low,  // 注意：高和低互换
			Low:   2*pivot - data.High, // 注意：高和低互换
			Close: 2*pivot - data.Close,

			// 高低点也互换
			IsHighPoint: data.IsLowPoint,
			IsLowPoint:  data.IsHighPoint,
		}

		reversed.Data[i] = reversedData
	}

	// 翻转上午下午的高低价
	reversed.MorningHigh = 2*pivot - dataSet.MorningLow
	reversed.MorningLow = 2*pivot - dataSet.MorningHigh
	reversed.AfternoonHigh = 2*pivot - dataSet.AfternoonLow
	reversed.AfternoonLow = 2*pivot - dataSet.AfternoonHigh

	reversed.PositiveCount = dataSet.NegativeCount
	reversed.NegativeCount = dataSet.PositiveCount

	return reversed
}
