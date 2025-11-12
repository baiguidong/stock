package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config 系统配置结构
// 包含所有可配置的参数，从config.json文件加载
type Config struct {
	Data struct {
		DailyDataPath  string `json:"daily_data_path"`  // 日线数据文件路径
		M1DataPath     string `json:"m1_data_path"`     // 分钟数据目录路径
		HistoryDataDir string `json:"history_data_dir"` // 历史数据目录
		OutputDir      string `json:"output_dir"`       // 输出目录
	} `json:"data"`

	Image struct {
		Width     int `json:"width"`      // 图像宽度
		Height    int `json:"height"`     // 图像高度
		LineWidth int `json:"line_width"` // 线条宽度
	} `json:"image"`

	Calculation struct {
		Modulo           int `json:"modulo"`             // 取模基数（默认8）
		TradingMinutes   int `json:"trading_minutes"`    // 每日交易分钟数（默认240）
		MorningEndMinute int `json:"morning_end_minute"` // 上午结束分钟（默认120）
	} `json:"calculation"`

	Display struct {
		SubImageWidth  int `json:"sub_image_width"`  // 子图宽度
		SubImageHeight int `json:"sub_image_height"` // 子图高度
		ShowCount      int `json:"show_count"`       // 显示数量
	} `json:"display"`
}

// DefaultConfig 返回默认配置
// 当配置文件不存在或加载失败时使用
func DefaultConfig() *Config {
	cfg := &Config{}

	// 数据路径配置
	cfg.Data.DailyDataPath = "data/日线.txt"
	cfg.Data.M1DataPath = "data/M1"
	cfg.Data.HistoryDataDir = "data/history"
	cfg.Data.OutputDir = "output"

	// 图像配置
	cfg.Image.Width = 600
	cfg.Image.Height = 300
	cfg.Image.LineWidth = 2

	// 计算参数配置
	cfg.Calculation.Modulo = 8
	cfg.Calculation.TradingMinutes = 240
	cfg.Calculation.MorningEndMinute = 120

	// 显示配置
	cfg.Display.SubImageWidth = 400
	cfg.Display.SubImageHeight = 200
	cfg.Display.ShowCount = 40

	return cfg
}

// LoadConfig 从文件加载配置
//
// 参数：
//   path - 配置文件路径
//
// 返回：
//   *Config - 配置对象
//   error   - 错误信息
//
// 注意：
//   如果文件不存在或加载失败，将返回默认配置并输出警告
func LoadConfig(path string) (*Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("警告: 配置文件不存在 %s，使用默认配置\n", path)
		return DefaultConfig(), nil
	}

	// 打开配置文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("警告: 无法打开配置文件 %s: %v，使用默认配置\n", path, err)
		return DefaultConfig(), nil
	}
	defer file.Close()

	// 解析JSON
	config := DefaultConfig() // 先加载默认配置，JSON会覆盖对应的字段
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	fmt.Printf("成功加载配置文件: %s\n", path)
	return config, nil
}

// Save 保存配置到文件
//
// 参数：
//   path - 配置文件路径
//
// 返回：
//   error - 错误信息
func (c *Config) Save(path string) error {
	// 创建配置文件
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建配置文件失败: %w", err)
	}
	defer file.Close()

	// 编码为JSON（带缩进）
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// Validate 验证配置的有效性
//
// 返回：
//   error - 如果配置无效，返回错误信息
func (c *Config) Validate() error {
	// 验证数据路径
	if c.Data.DailyDataPath == "" {
		return fmt.Errorf("日线数据路径不能为空")
	}
	if c.Data.M1DataPath == "" {
		return fmt.Errorf("分钟数据路径不能为空")
	}

	// 验证图像参数
	if c.Image.Width <= 0 {
		return fmt.Errorf("图像宽度必须大于0")
	}
	if c.Image.Height <= 0 {
		return fmt.Errorf("图像高度必须大于0")
	}

	// 验证计算参数
	if c.Calculation.Modulo <= 0 {
		return fmt.Errorf("取模基数必须大于0")
	}
	if c.Calculation.TradingMinutes <= 0 {
		return fmt.Errorf("交易分钟数必须大于0")
	}
	if c.Calculation.MorningEndMinute <= 0 || c.Calculation.MorningEndMinute >= c.Calculation.TradingMinutes {
		return fmt.Errorf("上午结束分钟数必须在有效范围内")
	}

	return nil
}
