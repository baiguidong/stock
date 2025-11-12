package dataloader

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"stock_predict/model"
	"stock_predict/utils"
)

// DailyDataLoader 日线数据加载器
// 负责从文件中加载日线(K线)数据
type DailyDataLoader struct {
	FilePath string // 数据文件路径
}

// NewDailyDataLoader 创建日线数据加载器
func NewDailyDataLoader(filePath string) *DailyDataLoader {
	return &DailyDataLoader{
		FilePath: filePath,
	}
}

// Load 加载日线数据
//
// 文件格式：
//   日期\t开盘价\t最高价\t最低价\t收盘价
//   2018/11/1\t3324.62\t3352.08\t3312.44\t3333.02
//
// 返回：
//   []model.DailyData - 日线数据切片
//   error - 错误信息
func (l *DailyDataLoader) Load() ([]model.DailyData, error) {
	// 检查文件是否存在
	if _, err := os.Stat(l.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("数据文件不存在: %s", l.FilePath)
	}

	// 读取文件内容
	data, err := ioutil.ReadFile(l.FilePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 检查文件是否为空
	if len(data) == 0 {
		return nil, fmt.Errorf("数据文件为空: %s", l.FilePath)
	}

	// 按行分割
	lines := strings.Split(string(data), "\n")
	var dailyData []model.DailyData

	// 逐行解析
	for lineNum, line := range lines {
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, " ", "", -1)
		line = strings.TrimSpace(line)

		// 跳过空行
		if line == "" {
			continue
		}

		// 按制表符分割
		fields := strings.Split(line, "\t")
		if len(fields) < 5 {
			continue
		}

		// 解析各字段
		date := utils.FormatDate(fields[0])
		open, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			fmt.Printf("警告: 第%d行开盘价解析失败: %v\n", lineNum+1, err)
			continue
		}

		high, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			fmt.Printf("警告: 第%d行最高价解析失败: %v\n", lineNum+1, err)
			continue
		}

		low, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			fmt.Printf("警告: 第%d行最低价解析失败: %v\n", lineNum+1, err)
			continue
		}

		close, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			fmt.Printf("警告: 第%d行收盘价解析失败: %v\n", lineNum+1, err)
			continue
		}

		// 创建日线数据对象
		dailyData = append(dailyData, model.DailyData{
			Date:  date,
			Open:  open,
			High:  high,
			Low:   low,
			Close: close,
		})
	}

	// 检查是否有有效数据
	if len(dailyData) == 0 {
		return nil, fmt.Errorf("未能解析出有效数据")
	}

	fmt.Printf("成功加载 %d 条日线数据\n", len(dailyData))
	return dailyData, nil
}
