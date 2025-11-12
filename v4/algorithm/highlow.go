package algorithm

import "stock_predict/model"

// HighLowDetector 高低点检测器
// 用于识别分钟数据中的高点和低点
type HighLowDetector struct {
	TradingMinutes   int // 每日交易分钟数（默认240）
	MorningEndMinute int // 上午结束分钟（默认120）
}

// NewHighLowDetector 创建高低点检测器
func NewHighLowDetector(tradingMinutes, morningEndMinute int) *HighLowDetector {
	return &HighLowDetector{
		TradingMinutes:   tradingMinutes,
		MorningEndMinute: morningEndMinute,
	}
}

// DetectHighLowPoints 识别数据集中的高低点
//
// 算法说明：
//   对于第1个点和最后一个点，比较相邻点判断
//   对于中间的点，比较前后两个点：
//   - 高点：比前后都高
//   - 低点：比前后都低
//
// 参数：
//   dataSet - 股票数据集
//
// 返回：
//   highCount - 高点数量
//   lowCount  - 低点数量
func (hld *HighLowDetector) DetectHighLowPoints(dataSet *model.StockDataSet) (int, int) {
	// 检查数据完整性
	if !hld.isDataComplete(dataSet) {
		return 0, 0
	}

	highCount := 0
	lowCount := 0

	// 处理第1个点（索引1）
	if dataSet.Data[1].Close > dataSet.Data[2].Close {
		dataSet.Data[1].IsHighPoint = true
		highCount++
	} else {
		dataSet.Data[1].IsLowPoint = true
		lowCount++
	}

	// 处理最后一个点（索引240）
	lastIdx := hld.TradingMinutes
	if dataSet.Data[lastIdx].Close > dataSet.Data[lastIdx-1].Close {
		dataSet.Data[lastIdx].IsHighPoint = true
		highCount++
	} else {
		dataSet.Data[lastIdx].IsLowPoint = true
		lowCount++
	}

	// 处理中间的点（索引2到239）
	for i := 2; i < lastIdx; i++ {
		prev := dataSet.Data[i-1]
		curr := dataSet.Data[i]
		next := dataSet.Data[i+1]

		// 跳过空数据
		if prev == nil || curr == nil || next == nil {
			continue
		}

		// 判断高点：比前后都高
		if curr.Close > prev.Close && curr.Close > next.Close {
			curr.IsHighPoint = true
			highCount++
		}

		// 判断低点：比前后都低
		if curr.Close < prev.Close && curr.Close < next.Close {
			curr.IsLowPoint = true
			lowCount++
		}
	}

	return highCount, lowCount
}

// CalculateSessionHighLow 计算上午和下午的最高最低价
//
// 算法说明：
//   - 前120分钟（上午）：计算最高价和最低价
//   - 后120分钟（下午）：计算最高价和最低价
//
// 参数：
//   dataSet - 股票数据集
func (hld *HighLowDetector) CalculateSessionHighLow(dataSet *model.StockDataSet) {
	// 检查数据完整性
	if !hld.isDataComplete(dataSet) {
		return
	}

	// 初始化上午最高最低价（使用第1分钟的收盘价）
	dataSet.MorningHigh = dataSet.Data[1].Close
	dataSet.MorningLow = dataSet.Data[1].Close

	// 初始化下午最高最低价（使用第121分钟的收盘价）
	dataSet.AfternoonHigh = dataSet.Data[hld.MorningEndMinute+1].Close
	dataSet.AfternoonLow = dataSet.Data[hld.MorningEndMinute+1].Close

	// 遍历所有分钟数据
	for i := 1; i <= hld.TradingMinutes; i++ {
		data := dataSet.Data[i]
		if data == nil {
			continue
		}

		// 上午时段（前120分钟）
		if i <= hld.MorningEndMinute {
			if data.Close > dataSet.MorningHigh {
				dataSet.MorningHigh = data.Close
			}
			if data.Close < dataSet.MorningLow {
				dataSet.MorningLow = data.Close
			}
		} else {
			// 下午时段（后120分钟）
			if data.Close > dataSet.AfternoonHigh {
				dataSet.AfternoonHigh = data.Close
			}
			if data.Close < dataSet.AfternoonLow {
				dataSet.AfternoonLow = data.Close
			}
		}
	}
}

// GetDayHighLow 获取全天的最高最低价
//
// 参数：
//   dataSet - 股票数据集
//
// 返回：
//   high - 全天最高价
//   low  - 全天最低价
func (hld *HighLowDetector) GetDayHighLow(dataSet *model.StockDataSet) (float64, float64) {
	if !hld.isDataComplete(dataSet) {
		return 0, 0
	}

	high := dataSet.Data[1].High
	low := dataSet.Data[1].Low

	// 遍历所有分钟数据
	for i := 1; i <= hld.TradingMinutes; i++ {
		data := dataSet.Data[i]
		if data == nil {
			continue
		}

		if data.High > high {
			high = data.High
		}
		if data.Low < low {
			low = data.Low
		}
	}

	return high, low
}

// GetOHLC 获取开高低收价格
//
// 参数：
//   dataSet - 股票数据集
//
// 返回：
//   open  - 开盘价（第1分钟的开盘价）
//   high  - 最高价
//   low   - 最低价
//   close - 收盘价（第240分钟的收盘价）
func (hld *HighLowDetector) GetOHLC(dataSet *model.StockDataSet) (float64, float64, float64, float64) {
	if !hld.isDataComplete(dataSet) {
		return 0, 0, 0, 0
	}

	open := dataSet.Data[1].Open
	close := dataSet.Data[hld.TradingMinutes].Close

	high, low := hld.GetDayHighLow(dataSet)

	return open, high, low, close
}

// isDataComplete 检查数据集是否完整
//
// 参数：
//   dataSet - 股票数据集
//
// 返回：
//   bool - 数据是否完整
func (hld *HighLowDetector) isDataComplete(dataSet *model.StockDataSet) bool {
	if dataSet == nil || dataSet.Data == nil {
		return false
	}

	// 检查所有分钟数据是否存在
	for i := 1; i <= hld.TradingMinutes; i++ {
		if dataSet.Data[i] == nil {
			return false
		}
	}

	return true
}

// IsTrendUp 判断是否为上涨走势
//
// 算法说明：
//   下午的最高价和最低价都高于上午的最高价和最低价
//
// 参数：
//   dataSet - 股票数据集
//
// 返回：
//   bool - 是否为上涨走势
func (hld *HighLowDetector) IsTrendUp(dataSet *model.StockDataSet) bool {
	return dataSet.AfternoonHigh > dataSet.MorningHigh &&
		dataSet.AfternoonLow > dataSet.MorningLow
}

// IsTrendDown 判断是否为下跌走势
//
// 算法说明：
//   下午的最高价和最低价都低于上午的最高价和最低价
//
// 参数：
//   dataSet - 股票数据集
//
// 返回：
//   bool - 是否为下跌走势
func (hld *HighLowDetector) IsTrendDown(dataSet *model.StockDataSet) bool {
	return dataSet.AfternoonHigh < dataSet.MorningHigh &&
		dataSet.AfternoonLow < dataSet.MorningLow
}

// GetTrendType 获取走势类型
//
// 参数：
//   dataSet - 股票数据集
//
// 返回：
//   TrendType - 走势类型（上涨/下跌/震荡）
func (hld *HighLowDetector) GetTrendType(dataSet *model.StockDataSet) model.TrendType {
	if hld.IsTrendUp(dataSet) {
		return model.TrendUp
	}
	if hld.IsTrendDown(dataSet) {
		return model.TrendDown
	}
	return model.TrendFlat
}
