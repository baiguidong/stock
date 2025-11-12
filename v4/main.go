package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"stock_predict/algorithm"
	"stock_predict/config"
	"stock_predict/dataloader"
	"stock_predict/image"
)

// 命令行参数
var (
	targetDateFlag = flag.String("day", "2019-05-28", "目标日期，格式：YYYY-MM-DD")
	offsetStartFlag = flag.Int("py", 0, "偏移起始值")
	offsetEndFlag   = flag.Int("py2", 2, "偏移结束值")
	topNFlag        = flag.Int("x", 3, "显示前N个匹配结果")
	configPathFlag  = flag.String("config", "config.json", "配置文件路径")
)

func main() {
	// 解析命令行参数
	flag.Parse()

	fmt.Println("========================================")
	fmt.Println("股票分时预测系统（重构版）")
	fmt.Println("========================================")

	// 加载配置
	cfg, err := config.LoadConfig(*configPathFlag)
	if err != nil {
		fmt.Printf("错误: 加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		fmt.Printf("错误: 配置无效: %v\n", err)
		os.Exit(1)
	}

	// 创建输出目录
	if err := os.MkdirAll(cfg.Data.OutputDir, 0755); err != nil {
		fmt.Printf("警告: 无法创建输出目录: %v\n", err)
	}

	fmt.Printf("\n目标日期: %s\n", *targetDateFlag)
	fmt.Printf("偏移范围: %d - %d\n", *offsetStartFlag, *offsetEndFlag)
	fmt.Printf("显示数量: 前 %d 个匹配结果\n\n", *topNFlag)

	// 1. 加载日线数据
	fmt.Println("步骤 1: 加载日线数据...")
	dailyLoader := dataloader.NewDailyDataLoader(cfg.Data.DailyDataPath)
	dailyData, err := dailyLoader.Load()
	if err != nil {
		fmt.Printf("错误: 加载日线数据失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 加载分钟数据
	fmt.Println("\n步骤 2: 加载分钟数据...")
	minuteLoader := dataloader.NewMinuteDataLoader(cfg.Data.M1DataPath)
	minuteDataMap, err := minuteLoader.LoadAll()
	if err != nil {
		fmt.Printf("错误: 加载分钟数据失败: %v\n", err)
		os.Exit(1)
	}

	// 3. 获取目标日期的分钟数据
	targetDataSet := minuteDataMap[*targetDateFlag]
	if targetDataSet == nil {
		fmt.Printf("错误: 未找到目标日期 %s 的分钟数据\n", *targetDateFlag)
		os.Exit(1)
	}

	// 4. 创建算法组件
	fmt.Println("\n步骤 3: 初始化算法组件...")
	featureExtractor := algorithm.NewFeatureExtractor(cfg.Calculation.Modulo)
	highLowDetector := algorithm.NewHighLowDetector(
		cfg.Calculation.TradingMinutes,
		cfg.Calculation.MorningEndMinute,
	)
	historicalMatcher := algorithm.NewHistoricalMatcher(featureExtractor, highLowDetector)

	// 5. 生成目标日期的图表
	fmt.Println("\n步骤 4: 生成目标日期图表...")
	chartGenerator := image.NewChartGenerator(
		cfg.Image.Width,
		cfg.Image.Height,
		float64(cfg.Image.LineWidth),
	)

	targetChartPath := filepath.Join(cfg.Data.OutputDir, fmt.Sprintf("%s_target.png", *targetDateFlag))
	if err := chartGenerator.GenerateHighLowPointsChart(targetDataSet, targetChartPath, fmt.Sprintf("目标日期: %s", *targetDateFlag)); err != nil {
		fmt.Printf("警告: 生成目标图表失败: %v\n", err)
	}

	// 6. 执行历史匹配
	fmt.Println("\n步骤 5: 执行历史匹配...")
	matchResults, err := historicalMatcher.FindMatches(
		*targetDateFlag,
		targetDataSet,
		dailyData,
		minuteDataMap,
	)
	if err != nil {
		fmt.Printf("错误: 历史匹配失败: %v\n", err)
		os.Exit(1)
	}

	// 7. 排序并获取前N个结果
	fmt.Println("\n步骤 6: 排序匹配结果...")
	matchResults.SortByPositiveCount()
	topMatches := matchResults.GetTopN(*topNFlag)

	// 8. 显示结果
	fmt.Println("\n========================================")
	fmt.Println("匹配结果")
	fmt.Println("========================================")
	fmt.Printf("总匹配数: %d\n", len(matchResults.Results))
	fmt.Printf("正向匹配: %d\n", matchResults.PositiveTotal)
	fmt.Printf("反向匹配: %d\n", matchResults.NegativeTotal)
	fmt.Println()

	if len(topMatches) == 0 {
		fmt.Println("未找到匹配的历史数据")
		return
	}

	fmt.Printf("前 %d 个匹配结果:\n\n", len(topMatches))
	for i, match := range topMatches {
		fmt.Printf("%d. 日期: %s\n", i+1, match.Date)
		fmt.Printf("   正向匹配: %d  反向匹配: %d\n", match.PositiveCount, match.NegativeCount)
		fmt.Printf("   走势类型: %s\n", match.TrendType.String())
		if match.IsReversed {
			fmt.Printf("   状态: 已反转\n")
		}
		fmt.Printf("   下一交易日: %s\n", match.NextDate)
		fmt.Println()

		// 生成对比图表
		compareChartPath := filepath.Join(
			cfg.Data.OutputDir,
			fmt.Sprintf("%s_vs_%s.png", *targetDateFlag, match.Date),
		)
		if err := chartGenerator.GenerateComparisonChart(
			targetDataSet,
			match.CurrentData,
			compareChartPath,
			fmt.Sprintf("对比: %s vs %s", *targetDateFlag, match.Date),
		); err != nil {
			fmt.Printf("   警告: 生成对比图表失败: %v\n", err)
		}

		// 如果有下一天数据，生成预测图表
		if match.NextData != nil {
			predictChartPath := filepath.Join(
				cfg.Data.OutputDir,
				fmt.Sprintf("%s_predict_%d.png", *targetDateFlag, i+1),
			)
			if err := chartGenerator.GenerateMinuteChart(
				match.NextData,
				predictChartPath,
				fmt.Sprintf("预测参考 %d: %s", i+1, match.NextDate),
			); err != nil {
				fmt.Printf("   警告: 生成预测图表失败: %v\n", err)
			}
		}
	}

	fmt.Println("========================================")
	fmt.Println("分析完成！")
	fmt.Printf("输出目录: %s\n", cfg.Data.OutputDir)
	fmt.Println("========================================")
}
