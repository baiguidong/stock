package image

import (
	"fmt"
	"image/color"
	"stock_predict/model"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// ChartGenerator 图表生成器
// 用于生成股票分时图等图表
type ChartGenerator struct {
	Width      int     // 图像宽度
	Height     int     // 图像高度
	LineWidth  float64 // 线条宽度
}

// NewChartGenerator 创建图表生成器
func NewChartGenerator(width, height int, lineWidth float64) *ChartGenerator {
	return &ChartGenerator{
		Width:     width,
		Height:    height,
		LineWidth: lineWidth,
	}
}

// GenerateMinuteChart 生成分钟数据图表
//
// 参数：
//   dataSet  - 股票数据集
//   filename - 输出文件名
//   title    - 图表标题
//
// 返回：
//   error - 错误信息
func (cg *ChartGenerator) GenerateMinuteChart(
	dataSet *model.StockDataSet,
	filename string,
	title string,
) error {

	// 创建新图表
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = "分钟"
	p.Y.Label.Text = "价格"

	// 准备数据点
	points := make(plotter.XYs, 0)
	for i := 1; i <= 240; i++ {
		data := dataSet.Data[i]
		if data != nil {
			points = append(points, plotter.XY{
				X: float64(i),
				Y: data.Close,
			})
		} else {
			points = append(points, plotter.XY{
				X: float64(i),
				Y: 0,
			})
		}
	}

	// 创建折线
	line, err := plotter.NewLine(points)
	if err != nil {
		return fmt.Errorf("创建折线失败: %w", err)
	}

	// 设置线条样式
	line.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255} // 蓝色
	line.Width = vg.Points(cg.LineWidth)

	// 添加到图表
	p.Add(line)

	// 保存图表
	if err := p.Save(vg.Points(float64(cg.Width)), vg.Points(float64(cg.Height)), filename); err != nil {
		return fmt.Errorf("保存图表失败: %w", err)
	}

	fmt.Printf("图表已保存: %s\n", filename)
	return nil
}

// GenerateComparisonChart 生成对比图表
//
// 参数：
//   dataSet1  - 数据集1（目标数据）
//   dataSet2  - 数据集2（匹配数据）
//   filename  - 输出文件名
//   title     - 图表标题
//
// 返回：
//   error - 错误信息
func (cg *ChartGenerator) GenerateComparisonChart(
	dataSet1 *model.StockDataSet,
	dataSet2 *model.StockDataSet,
	filename string,
	title string,
) error {

	// 创建新图表
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = "分钟"
	p.Y.Label.Text = "价格"

	// 准备数据集1的数据点
	points1 := make(plotter.XYs, 0)
	for i := 1; i <= 240; i++ {
		data := dataSet1.Data[i]
		if data != nil {
			points1 = append(points1, plotter.XY{
				X: float64(i),
				Y: data.Close,
			})
		}
	}

	// 准备数据集2的数据点
	points2 := make(plotter.XYs, 0)
	for i := 1; i <= 240; i++ {
		data := dataSet2.Data[i]
		if data != nil {
			points2 = append(points2, plotter.XY{
				X: float64(i),
				Y: data.Close,
			})
		}
	}

	// 创建折线1
	line1, err := plotter.NewLine(points1)
	if err != nil {
		return fmt.Errorf("创建折线1失败: %w", err)
	}
	line1.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255} // 蓝色
	line1.Width = vg.Points(cg.LineWidth)

	// 创建折线2
	line2, err := plotter.NewLine(points2)
	if err != nil {
		return fmt.Errorf("创建折线2失败: %w", err)
	}
	line2.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // 红色
	line2.Width = vg.Points(cg.LineWidth)

	// 添加到图表
	p.Add(line1, line2)

	// 添加图例
	p.Legend.Add("目标数据", line1)
	p.Legend.Add("匹配数据", line2)
	p.Legend.Top = true

	// 保存图表
	if err := p.Save(vg.Points(float64(cg.Width)), vg.Points(float64(cg.Height)), filename); err != nil {
		return fmt.Errorf("保存图表失败: %w", err)
	}

	fmt.Printf("对比图表已保存: %s\n", filename)
	return nil
}

// GenerateHighLowPointsChart 生成带高低点标记的图表
//
// 参数：
//   dataSet  - 股票数据集
//   filename - 输出文件名
//   title    - 图表标题
//
// 返回：
//   error - 错误信息
func (cg *ChartGenerator) GenerateHighLowPointsChart(
	dataSet *model.StockDataSet,
	filename string,
	title string,
) error {

	// 创建新图表
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = "分钟"
	p.Y.Label.Text = "价格"

	// 准备数据点
	points := make(plotter.XYs, 0)
	highPoints := make(plotter.XYs, 0)
	lowPoints := make(plotter.XYs, 0)

	for i := 1; i <= 240; i++ {
		data := dataSet.Data[i]
		if data != nil {
			points = append(points, plotter.XY{
				X: float64(i),
				Y: data.Close,
			})

			// 标记高点
			if data.IsHighPoint {
				highPoints = append(highPoints, plotter.XY{
					X: float64(i),
					Y: data.Close,
				})
			}

			// 标记低点
			if data.IsLowPoint {
				lowPoints = append(lowPoints, plotter.XY{
					X: float64(i),
					Y: data.Close,
				})
			}
		}
	}

	// 创建主折线
	line, err := plotter.NewLine(points)
	if err != nil {
		return fmt.Errorf("创建折线失败: %w", err)
	}
	line.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	line.Width = vg.Points(cg.LineWidth)

	p.Add(line)

	// 添加高点标记
	if len(highPoints) > 0 {
		highScatter, err := plotter.NewScatter(highPoints)
		if err == nil {
			highScatter.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // 红色
			p.Add(highScatter)
			p.Legend.Add("高点", highScatter)
		}
	}

	// 添加低点标记
	if len(lowPoints) > 0 {
		lowScatter, err := plotter.NewScatter(lowPoints)
		if err == nil {
			lowScatter.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255} // 绿色
			p.Add(lowScatter)
			p.Legend.Add("低点", lowScatter)
		}
	}

	p.Legend.Top = true

	// 保存图表
	if err := p.Save(vg.Points(float64(cg.Width)), vg.Points(float64(cg.Height)), filename); err != nil {
		return fmt.Errorf("保存图表失败: %w", err)
	}

	fmt.Printf("高低点图表已保存: %s\n", filename)
	return nil
}
