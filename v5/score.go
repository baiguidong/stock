package main

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/plot/plotter"
)

type NewPoint struct {
	X, Y float64
	g, d int
}

// —— 工具函数：从 XYs 中计算所有高低点 ——
func extractPoints(data plotter.XYs) []NewPoint {
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

func computeScore(dt1, dt plotter.XYs) (score int) {
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
func splitSegments(dt plotter.XYs, seg int) []int {
	n := len(dt)
	segSize := float64(n) / float64(seg)

	res := make([]int, seg)
	for i := 0; i < seg; i++ {
		start := int(segSize * float64(i))
		end := int(segSize*float64(i+1)) - 1
		if start < 0 {
			start = 0
		}
		res[i] = segType(dt[start].Y, dt[end].Y)
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

func computeScoreSeg(dt1, dt plotter.XYs) (score int) {
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

var arr1 = []plotter.XY{
	{180, 3949.83}, {181, 3942.15}, {182, 3944.18}, {183, 3943.63},
	{184, 3944.65}, {185, 3942.54}, {186, 3940.99}, {187, 3939.16},
	{188, 3938.37}, {189, 3938.56}, {190, 3942.28}, {191, 3942.73},
	{192, 3941.82}, {193, 3939.64}, {194, 3937.21}, {195, 3937.23},
	{196, 3938.54}, {197, 3938.73}, {198, 3939.49}, {199, 3940.65},
	{200, 3942.29}, {201, 3946.91}, {202, 3948.24}, {203, 3951.84},
	{204, 3951}, {205, 3950.41}, {206, 3948.17}, {207, 3946.76},
	{208, 3944.61}, {209, 3943.86}, {210, 3941.36}, {211, 3937.97},
	{212, 3938.02}, {213, 3941}, {214, 3938.96}, {215, 3937.97},
	{216, 3936.48}, {217, 3935.36}, {218, 3935.3}, {219, 3933.78},
	{220, 3931.84}, {221, 3930.54}, {222, 3931.09}, {223, 3931.63},
	{224, 3930.66}, {225, 3928.58}, {226, 3929.04}, {227, 3933.94},
	{228, 3937.27}, {229, 3936.14}, {230, 3932.01}, {231, 3937.28},
	{232, 3938.32}, {233, 3938.94}, {234, 3940.21}, {235, 3939.17},
	{236, 3938.58}, {237, 3939.57}, {238, 3940.23}, {239, 3940.56},
	{240, 3939.81},
}

var arr2 = []plotter.XY{
	{180, 2166.369}, {181, 2166.668}, {182, 2166.649}, {183, 2166.672},
	{184, 2166.726}, {185, 2166.716}, {186, 2166.781}, {187, 2166.897},
	{188, 2166.634}, {189, 2166.761}, {190, 2166.599}, {191, 2166.547},
	{192, 2166.513}, {193, 2166.304}, {194, 2166.097}, {195, 2165.752},
	{196, 2165.474}, {197, 2165.309}, {198, 2164.95}, {199, 2164.802},
	{200, 2164.459}, {201, 2164.341}, {202, 2164.302}, {203, 2164.244},
	{204, 2163.896}, {205, 2163.613}, {206, 2163.432}, {207, 2163.314},
	{208, 2163.311}, {209, 2163.162}, {210, 2162.934}, {211, 2162.811},
	{212, 2162.824}, {213, 2162.871}, {214, 2163.057}, {215, 2163.015},
	{216, 2163.206}, {217, 2163.223}, {218, 2163.458}, {219, 2163.57},
	{220, 2163.649}, {221, 2163.566}, {222, 2163.683}, {223, 2163.676},
	{224, 2163.454}, {225, 2163.399}, {226, 2162.897}, {227, 2162.339},
	{228, 2161.597}, {229, 2160.969}, {230, 2160.545}, {231, 2159.817},
	{232, 2159.66}, {233, 2159.381}, {234, 2159.249}, {235, 2159.121},
	{236, 2159.169}, {237, 2159.849}, {238, 2160.723}, {239, 2161.34},
	{240, 2161.34},
}

func Score(dt1, dt plotter.XYs) int {
	score1 := computeScore(dt1, dt)
	score2 := computeScoreSeg(dt1, dt)
	return score1 + score2
}
