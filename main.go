package main

import (
	"fmt"
	"math"
	"sort"
)

type XYs []XY

type XY struct{ X, Y float64 }

type NewPoint struct {
	X, Y float64
	g, d int
}

// —— 工具函数：从 XYs 中计算所有高低点 ——
func extractPoints(data XYs) []NewPoint {
	n := len(data)
	res := make([]NewPoint, 0)

	for i := 0; i < n; i++ {
		cur := data[i].Y
		if i == 0 {
			next := data[i+1].Y
			if cur > next {
				res = append(res, NewPoint{
					X: data[i].X,
					Y: cur,
					g: 1,
				})
			} else if cur < next {
				res = append(res, NewPoint{
					X: data[i].X,
					Y: cur,
					d: 1,
				})
			} else {
				res = append(res, NewPoint{
					X: data[i].X,
					Y: cur,
				})
			}
			continue
		}
		if i == n-1 {
			prev := data[i-1].Y
			if cur > prev {
				res = append(res, NewPoint{
					X: data[i].X,
					Y: cur,
					g: 1,
				})
			} else if cur < prev {
				res = append(res, NewPoint{
					X: data[i].X,
					Y: cur,
					d: 1,
				})
			} else {
				res = append(res, NewPoint{
					X: data[i].X,
					Y: cur,
				})
			}
			continue
		}

		prev := data[i-1].Y
		next := data[i+1].Y

		if cur > prev && cur > next { // 高点
			res = append(res, NewPoint{
				X: data[i].X,
				Y: cur,
				g: 1,
			})
		} else if cur < prev && cur < next { // 低点
			res = append(res, NewPoint{
				X: data[i].X,
				Y: cur,
				d: 1,
			})
		} else {
			res = append(res, NewPoint{
				X: data[i].X,
				Y: cur,
			})
		}
	}
	return res
}

// —— 工具：从高低点里筛出高点
func filterHigh(arr []NewPoint) []NewPoint {
	res := make([]NewPoint, 0)
	for _, v := range arr {
		if v.g == 1 {
			res = append(res, v)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Y > res[j].Y // 降序
	})
	return res
}

// —— 工具：从高低点里筛出低点
func filterLow(arr []NewPoint) []NewPoint {
	res := make([]NewPoint, 0)
	for _, v := range arr {
		if v.d == 1 {
			res = append(res, v)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Y < res[j].Y // 升序
	})
	return res
}

// 输入点 points
// 后面 N 个
// 高点个数 countHigh
// 低点个数 countLow
func selectLastNHighLow(points []NewPoint, N int, countHigh int, countLow int) (hs, ls []NewPoint) {
	n := len(points)
	start := n - N
	if start < 0 {
		start = 0
	}

	seg := points[start:]

	hsAll := filterHigh(seg)
	lsAll := filterLow(seg)

	if len(hsAll) > countHigh {
		hs = hsAll[:countHigh]
	} else {
		hs = hsAll
	}

	if len(lsAll) > countLow {
		ls = lsAll[:countLow]
	} else {
		ls = lsAll
	}
	return
}

// 前一天所有点
// 高低点
// high 匹配高点
// low 匹配低点
func matchNearby(points []NewPoint, x float64, high, low bool) bool {
	for _, p := range points {
		if math.Abs(p.X-x) <= 4 { // X 在 ±4 范围内即可认为匹配
			if high && p.g == 1 {
				return true
			}
			if low && p.d == 1 {
				return true
			}
		}
	}
	return false
}

// 参数 dt1 dt 为 61 个点（X 180-240）
// 1.计算 最后 60 个点的高低点,生成 NewPoint 数组
// 高低点逻辑如下:
// 高点: 比前一个点高,比后一个点高
// 低点: 比前一个点低,比后一个点低
// 最后一个点只和前一点比较
// 2.取dt 后60个点中5个高点,和 dt1 中附近位置（X 坐标 +4）和（X坐标 -4） 是否存在高点，是的话得分5
// 3.取dt 后60个点中5个低点,和 dt1 中附近位置（X 坐标 +4）和（X坐标 -4） 是否存在低点，是的话得分5
// 4.取dt 后40个点中2个高点,和 dt1 中附近位置（X 坐标 +4）和（X坐标 -4） 是否存在高点，是的话得分10
// 5.取dt 后40个点中2个低点,和 dt1 中附近位置（X 坐标 +4）和（X坐标 -4） 是否存在低点，是的话得分10
// 6.取dt 后30个点中2个高点,和 dt1 中附近位置（X 坐标 +4）和（X坐标 -4） 是否存在高点，是的话得分10
// 7.取dt 后30个点中2个低点,和 dt1 中附近位置（X 坐标 +4）和（X坐标 -4） 是否存在低点，是的话得分10

func computeScore(dt1, dt XYs) (score int) {
	// —— 1. 提取 dt 和 dt1 中全部高低点 ——
	dtPoints := extractPoints(dt)
	dt1Points := extractPoints(dt1)

	addScore := func(hs, ls []NewPoint, highScore, lowScore int) {
		for _, h := range hs {
			if matchNearby(dt1Points, h.X, true, false) {
				score += highScore
			}
		}
		for _, l := range ls {
			if matchNearby(dt1Points, l.X, false, true) {
				score += lowScore
			}
		}
	}

	// 2. 后 60 点：取 5 个高点、5 个低点，每个匹配得 5 分
	hs60, ls60 := selectLastNHighLow(dtPoints, 60, 5, 5)
	addScore(hs60, ls60, 5, 5)
	fmt.Println("hs60:", hs60, ls60)

	// 4. 后 40 点：2 个高点、2 个低点，每个匹配得 10 分
	hs40, ls40 := selectLastNHighLow(dtPoints, 40, 2, 2)
	addScore(hs40, ls40, 10, 10)
	fmt.Println("hs40:", hs40, ls40)

	//6. 后 30 点：2 个高点、2 个低点，每个匹配得 10 分
	hs30, ls30 := selectLastNHighLow(dtPoints, 30, 2, 2)
	addScore(hs30, ls30, 10, 10)
	fmt.Println("hs30:", hs30, ls30)
	return
}

// 段的结果：1 = 涨, -1 = 跌, 0 = 平
func segType(a, b float64) int {
	if b > a {
		return 1 // 涨
	}
	if b < a {
		return -1 // 跌
	}
	return 0 // 平
}

// 把 XYs 分 n 段，返回每段的 涨/跌/平 结果
func splitSegments(dt XYs, seg int) []int {

	n := len(dt)
	fmt.Println("splitSegments:", n, seg)
	segSize := float64(n) / float64(seg)

	res := make([]int, seg)
	for i := 0; i < seg; i++ {
		start := int(segSize*float64(i)) - 1
		end := int(segSize*float64(i+1)) - 1
		if start < 0 {
			start = 0
		}
		res[i] = segType(dt[start].Y, dt[end].Y)
		fmt.Println("segType:", start, end, res[i])
	}
	return res
}

// 比较两个模式，完全相同则加分
func compareSegs(a, b []int, scorePerSeg int) int {
	score := 0
	for i := range a {
		if a[i] == b[i] {
			score += scorePerSeg
		}
	}
	fmt.Println(score)
	return score
}

// 参数 dt1 dt 为 61 个点（X 180-240,X 按照顺序的)
// 计算波段的分
// 涨 尾部>头部
// 跌 尾部<头部
// 平 尾部=头部
// 1. 取dt最后60个点，分成 12 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的5分
// 2. 取dt最后60个点，分成 10 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的5分
// 3. 取dt最后60个点，分成 6 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的5分
// 4. 取dt最后60个点，分成 5 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的5分
// 5. 取dt最后60个点，分成 4 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的5分

// 6. 取dt最后40个点，分成 10 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的15分
// 7. 取dt最后40个点，分成 8 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的15分
// 8. 取dt最后40个点，分成 5 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的15分
// 9. 取dt最后40个点，分成 4 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的15分
// 10. 取dt最后40个点，分成 2 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的15分

// 11. 取dt最后30个点，分成 6 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的15分
// 12. 取dt最后30个点，分成 5 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的15分
// 13. 取dt最后30个点，分成 3 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的15分
// 14. 取dt最后30个点，分成 2 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的15分
// 15. 取dt最后30个点，分成 1 段,每一段从头到尾,判读 涨/跌/平 三种模式，对比dt1 同样分段的 涨跌平，相同的15分

func computeScoreSeg(dt1, dt XYs) (score int) {
	// 所有任务规则定义：{使用末尾点数, 段数, 每段分值}
	rules := [][3]int{
		{60, 12, 5}, {60, 10, 5}, {60, 6, 5}, {60, 5, 5}, {60, 4, 5},
		{40, 10, 15}, {40, 8, 15}, {40, 5, 15}, {40, 4, 15}, {40, 2, 15},
		{30, 6, 15}, {30, 5, 15}, {30, 3, 15}, {30, 2, 15}, {30, 1, 15},
	}

	for _, r := range rules {
		lastN := r[0]
		seg := r[1]
		segScore := r[2]

		// 截取最后 N 个点
		sub1 := dt1[len(dt1)-lastN:]
		sub2 := dt[len(dt)-lastN:]

		// 计算各自的分段模式
		segA := splitSegments(sub1, seg)
		segB := splitSegments(sub2, seg)

		// 对比 + 计分
		score += compareSegs(segA, segB, segScore)

	}

	return
}

var dt1 = XYs{
	{180, 2601.74}, {181, 2601.85}, {182, 2600.99}, {183, 2603.34}, {184, 2605.61},
	{185, 2607.41}, {186, 2607.34}, {187, 2608.65}, {188, 2607.12}, {189, 2607.27},
	{190, 2607.23}, {191, 2608.78}, {192, 2610.46}, {193, 2612.44}, {194, 2610.3},
	{195, 2607.86}, {196, 2606.52}, {197, 2604.68}, {198, 2604.87}, {199, 2604.56},
	{200, 2604.57}, {201, 2605.02}, {202, 2605.72}, {203, 2606.57}, {204, 2607.45},
	{205, 2608.9}, {206, 2609.93}, {207, 2608.8}, {208, 2607.1}, {209, 2607.75},
	{210, 2607.03}, {211, 2606.1}, {212, 2603.53}, {213, 2600.88}, {214, 2600.54},
	{215, 2600.4}, {216, 2601.21}, {217, 2600.45}, {218, 2601.71}, {219, 2601.05},
	{220, 2599.56}, {221, 2598.65}, {222, 2595.84}, {223, 2596.41}, {224, 2596.62},
	{225, 2597.05}, {226, 2598.4}, {227, 2599.52}, {228, 2600.73}, {229, 2599.98},
	{230, 2600.63}, {231, 2600.53}, {232, 2601.88}, {233, 2602.23}, {234, 2602.21},
	{235, 2603.21}, {236, 2603.21}, {237, 2602.8}, {238, 2603.37}, {239, 2603.37},
	{240, 2603.3},
}

var dt = XYs{
	{180, 3868.68}, {181, 3869.21}, {182, 3868.16}, {183, 3865.89}, {184, 3866.32},
	{185, 3865.25}, {186, 3864.2}, {187, 3865.19}, {188, 3865.95}, {189, 3867.11},
	{190, 3866.28}, {191, 3867.12}, {192, 3865.67}, {193, 3867.97}, {194, 3869.49},
	{195, 3869.53}, {196, 3867.6}, {197, 3866.29}, {198, 3866.32}, {199, 3866.75},
	{200, 3867.23}, {201, 3867.11}, {202, 3866.42}, {203, 3866.68}, {204, 3867.4},
	{205, 3867.89}, {206, 3868.48}, {207, 3868.33}, {208, 3867.58}, {209, 3867.68},
	{210, 3867.03}, {211, 3869.48}, {212, 3867.81}, {213, 3866.55}, {214, 3866.65},
	{215, 3867.29}, {216, 3868.02}, {217, 3867.33}, {218, 3866.75}, {219, 3866.42},
	{220, 3866.11}, {221, 3865.25}, {222, 3864.38}, {223, 3863.13}, {224, 3863.03},
	{225, 3861.98}, {226, 3862.33}, {227, 3862.99}, {228, 3863.04}, {229, 3862.46},
	{230, 3862.34}, {231, 3862.33}, {232, 3861.98}, {233, 3862.19}, {234, 3862.76},
	{235, 3863}, {236, 3863.72}, {237, 3864.83}, {238, 3865.03}, {239, 3865.03},
	{240, 3864.18},
}

func PrintXYsForExcel(xs []NewPoint) {
	// 打印表头
	fmt.Println("X\tY")

	// 打印内容
	for _, v := range xs {
		fmt.Printf("%v\t%v\t%v\t%v\n", v.X, v.Y, v.g, v.d)
	}
}
func main() {
	PrintXYsForExcel(extractPoints(dt1))
	fmt.Println("===>")
	PrintXYsForExcel(extractPoints(dt))
	//computeScore(dt1, dt)
	score := computeScore(dt1, dt)
	fmt.Println(score)
	score = computeScoreSeg(dt1, dt)
	fmt.Println(score)
}
