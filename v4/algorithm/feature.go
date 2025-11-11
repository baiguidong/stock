package algorithm

import (
	"fmt"
	"stock_predict/model"
)

// FeatureExtractor 特征提取器
// 用于从价格中提取数字特征（余数值）
type FeatureExtractor struct {
	Modulo int // 取模基数，默认8
}

// NewFeatureExtractor 创建特征提取器
func NewFeatureExtractor(modulo int) *FeatureExtractor {
	return &FeatureExtractor{Modulo: modulo}
}

// ExtractFromPrice 从价格提取特征值
//
// 算法说明：
//   将价格数字逐位相加后除以模数，取余数作为特征值
//   这是本系统的核心特征提取算法
//
// 参数：
//   price - 股票价格（保留两位小数）
//
// 返回：
//   特征值，范围[1,模数]，用于历史匹配
//
// 示例：
//   3324.62 → 3+3+2+4+6+2=20 → 20%8=4 → 返回4
//   3328.08 → 3+3+2+8+0+8=24 → 24%8=0 → 返回8（整除按模数处理）
//   3300.01 → 3+3+0+0+0+1=7 → 7<8 → 返回8（小于模数按模数处理）
func (fe *FeatureExtractor) ExtractFromPrice(price float64) int {
	// 转换为固定格式的字符串（保留两位小数）
	priceStr := fmt.Sprintf("%.2f", price)

	// 累加所有数字位（忽略小数点）
	digitSum := int64(0)
	for _, char := range priceStr {
		if char >= '0' && char <= '9' {
			digit := int64(char - '0') // 字符转数字
			digitSum += digit
		}
	}

	// 对模数取余获得余数
	remainder := digitSum % int64(fe.Modulo)

	// 特殊情况：余数为0时按模数处理
	if remainder == 0 {
		remainder = int64(fe.Modulo)
	}

	return int(remainder)
}

// ExtractFromOHLC 从OHLC提取特征组
//
// 参数：
//   open  - 开盘价
//   high  - 最高价
//   low   - 最低价
//   close - 收盘价
//
// 返回：
//   FeatureSet 包含四个价格的特征值
func (fe *FeatureExtractor) ExtractFromOHLC(open, high, low, close float64) model.FeatureSet {
	return model.FeatureSet{
		Open:  fe.ExtractFromPrice(open),
		High:  fe.ExtractFromPrice(high),
		Low:   fe.ExtractFromPrice(low),
		Close: fe.ExtractFromPrice(close),
	}
}

// GenerateMatchPatterns 生成匹配模式
//
// 算法说明：
//   基于当前特征，生成15组变化的匹配模式
//   - 开盘和收盘保持不变
//   - 最高价从1到模数循环（8次）
//   - 最低价从1到模数循环（8次）
//   - 去重后得到15组独特的模式
//
// 参数：
//   baseFeature - 基础特征值
//
// 返回：
//   []FeatureSet 15组匹配模式
func (fe *FeatureExtractor) GenerateMatchPatterns(baseFeature model.FeatureSet) []model.FeatureSet {
	patterns := make([]model.FeatureSet, 0, 15)
	seen := make(map[string]bool)

	// 最高价循环（最低价保持原值）
	for h := 1; h <= fe.Modulo; h++ {
		pattern := model.FeatureSet{
			Open:  baseFeature.Open,
			High:  h,
			Low:   baseFeature.Low,
			Close: baseFeature.Close,
		}

		key := pattern.String()
		if !seen[key] {
			patterns = append(patterns, pattern)
			seen[key] = true
		}
	}

	// 最低价循环（最高价保持原值）
	for l := 1; l <= fe.Modulo; l++ {
		pattern := model.FeatureSet{
			Open:  baseFeature.Open,
			High:  baseFeature.High,
			Low:   l,
			Close: baseFeature.Close,
		}

		key := pattern.String()
		if !seen[key] {
			patterns = append(patterns, pattern)
			seen[key] = true
		}
	}

	return patterns
}
