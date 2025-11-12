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

// MinuteDataLoader 分钟数据加载器
// 负责从文件中加载分钟级别的K线数据
type MinuteDataLoader struct {
	FilePath string // 数据文件路径
}

// NewMinuteDataLoader 创建分钟数据加载器
func NewMinuteDataLoader(filePath string) *MinuteDataLoader {
	return &MinuteDataLoader{
		FilePath: filePath,
	}
}

// Load 加载分钟数据
//
// 文件格式：
//   日期-时间\t开盘价\t最高价\t最低价\t收盘价
//   2018-11-01-9:31\t3324.62\t3326.08\t3320.44\t3325.02
//
// 返回：
//   map[string]*model.StockDataSet - 日期 -> 数据集合的映射
//   error - 错误信息
func (l *MinuteDataLoader) Load() (map[string]*model.StockDataSet, error) {
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

	// 创建结果映射
	dataMap := make(map[string]*model.StockDataSet)

	// 按行分割
	lines := strings.Split(string(data), "\n")

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

		// 解析日期时间
		dateTimeParts := strings.Split(fields[0], "-")
		if len(dateTimeParts) != 2 {
			fmt.Printf("警告: 第%d行日期时间格式错误: %s\n", lineNum+1, fields[0])
			continue
		}

		date := utils.FormatDate(dateTimeParts[0])
		time := dateTimeParts[1]

		// 解析价格
		open, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			continue
		}

		high, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			continue
		}

		low, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			continue
		}

		close, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			continue
		}

		// 获取时间对应的索引
		timeIndex := model.GetTimeIndex(time)
		if timeIndex == 0 {
			fmt.Printf("警告: 第%d行无效的时间: %s\n", lineNum+1, time)
			continue
		}

		// 创建K线数据
		stockData := &model.StockData{
			Date:  date,
			Time:  time,
			Open:  open,
			High:  high,
			Low:   low,
			Close: close,
		}

		// 添加到映射中
		if _, exists := dataMap[date]; !exists {
			dataMap[date] = &model.StockDataSet{}
			dataMap[date].InitData()
		}

		dataMap[date].Add(stockData, timeIndex)
	}

	// 检查是否有有效数据
	if len(dataMap) == 0 {
		return nil, fmt.Errorf("未能解析出有效数据")
	}

	fmt.Printf("成功加载 %d 天的分钟数据\n", len(dataMap))
	return dataMap, nil
}

// LoadAll 加载目录下的所有分钟数据文件
//
// 参数：
//   无，使用FilePath作为目录路径
//
// 返回：
//   map[string]*model.StockDataSet - 日期 -> 数据集合的映射
//   error - 错误信息
func (l *MinuteDataLoader) LoadAll() (map[string]*model.StockDataSet, error) {
	// 检查路径是否是目录
	fileInfo, err := os.Stat(l.FilePath)
	if err != nil {
		return nil, fmt.Errorf("路径不存在: %s, %w", l.FilePath, err)
	}

	// 如果是文件，直接加载
	if !fileInfo.IsDir() {
		return l.Load()
	}

	// 如果是目录，遍历所有文件
	files, err := ioutil.ReadDir(l.FilePath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	// 合并所有文件的数据
	allDataMap := make(map[string]*model.StockDataSet)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 只处理.txt文件
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".txt") {
			continue
		}

		// 创建新的加载器加载单个文件
		filePath := l.FilePath + "/" + fileName
		loader := NewMinuteDataLoader(filePath)
		dataMap, err := loader.Load()
		if err != nil {
			fmt.Printf("警告: 加载文件 %s 失败: %v\n", fileName, err)
			continue
		}

		// 合并数据
		for date, dataSet := range dataMap {
			allDataMap[date] = dataSet
		}
	}

	if len(allDataMap) == 0 {
		return nil, fmt.Errorf("目录 %s 中没有有效的分钟数据", l.FilePath)
	}

	fmt.Printf("成功从目录 %s 加载 %d 天的分钟数据\n", l.FilePath, len(allDataMap))
	return allDataMap, nil
}
