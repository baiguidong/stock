# 代码重构指南

本文档说明如何将现有代码改进为更规范、更易维护的版本。

---

## 改进概览

### 当前代码问题

| 问题类型 | 具体表现 | 影响 |
|---------|---------|------|
| **命名不清** | 变量使用拼音缩写 (dy, sed, seds) | 可读性差 |
| **注释缺失** | 大部分函数无注释 | 难以理解 |
| **硬编码** | 路径、参数直接写在代码中 | 灵活性差 |
| **结构混乱** | 所有逻辑在两个大文件中 | 难以维护 |
| **错误处理** | 缺少完善的错误处理 | 容易崩溃 |

### 改进目标

- ✅ 清晰的命名规范
- ✅ 完善的代码注释
- ✅ 配置文件化
- ✅ 模块化设计
- ✅ 错误处理机制
- ✅ 日志系统

---

## 命名规范改进

### 变量命名对照表

#### 当前命名 → 建议命名

```go
// 数据结构
sed     → StockData           // 股票数据
seds    → StockDataSet        // 股票数据集
sds     → MatchResult         // 匹配结果

// 字段名
dy      → date                // 日期
tm      → time                // 时间
dk      → open                // 开盘价
dg      → high                // 最高价
dd      → low                 // 最低价
ds      → close               // 收盘价
g       → isHighPoint         // 是否高点
d       → isLowPoint          // 是否低点

// 集合字段
sg      → morningHigh         // 上午最高
sd      → morningLow          // 上午最低
xg      → afternoonHigh       // 下午最高
xd      → afternoonLow        // 下午最低
zh      → positiveCount       // 正值数量
fa      → negativeCount       // 反值数量

// 全局变量
g_map   → historicalDataMap   // 历史数据映射
g_map_m1→ minuteDataMap       // 分钟数据映射
g_dz    → timeIndexMap        // 时间索引映射
g_py    → offsetValue         // 偏移值
g_day   → targetDate          // 目标日期
```

### 改进示例

**改进前:**
```go
type sed struct {
    dy, tm         string
    dk, dg, dd, ds float64
    g, d           int
}
```

**改进后:**
```go
// StockData 股票数据结构
// 包含单根K线的完整信息:日期、时间、开高低收价格以及高低点标记
type StockData struct {
    Date  string  // 交易日期,格式:YYYY-MM-DD
    Time  string  // 交易时间,格式:HH:MM

    Open  float64 // 开盘价
    High  float64 // 最高价
    Low   float64 // 最低价
    Close float64 // 收盘价

    IsHighPoint bool // 是否为高点(比前后都高)
    IsLowPoint  bool // 是否为低点(比前后都低)
}
```

---

## 代码注释规范

### 函数注释模板

```go
// CalculateRemainder 计算价格特征余数
//
// 功能说明:
//   将价格的每一位数字相加后除以8取余数,作为该价格的特征值
//
// 参数:
//   price float64 - 输入价格,如 3324.62
//
// 返回值:
//   int - 特征余数,范围[1-8]
//
// 特殊处理:
//   - 余数为0时返回8
//   - 相加结果<8时返回8
//
// 示例:
//   CalculateRemainder(3324.62) // 返回: 4
//   计算过程: 3+3+2+4+6+2=20, 20÷8=2余4
func CalculateRemainder(price float64) int {
    // 格式化为两位小数字符串
    priceStr := fmt.Sprintf("%.2f", price)

    // 累加所有数字(跳过小数点)
    sum := int64(0)
    for i := 0; i < len(priceStr); i++ {
        char := priceStr[i : i+1]
        if char != "." {
            digit, _ := strconv.ParseInt(char, 10, 10)
            sum += digit
        }
    }

    // 对8取余
    remainder := sum % 8

    // 特殊情况处理
    if remainder == 0 {
        remainder = 8
    }

    return int(remainder)
}
```

### 改进对照

**改进前:**
```go
//转换 int  对 8 取余数
func get_he(f float64) int {
    s := fmt.Sprintf("%.2f", f)
    i := 0
    num := int64(0)
    for i < len(s) {
        s1 := s[i : i+1]
        if s1 != "." {
            n, _ := strconv.ParseInt(s1, 10, 10)
            num += n
        }
        i++
    }
    num = num % 8
    if num == 0 {
        num = 8
    }
    return int(num)
}
```

**改进后:**
```go
// CalculateRemainder 计算价格特征余数
//
// 算法说明:
//   将价格数字逐位相加后除以8,取余数作为特征值
//   这是本系统的核心特征提取算法
//
// 参数:
//   price - 股票价格(保留两位小数)
//
// 返回:
//   特征值,范围[1,8],用于历史匹配
//
// 示例:
//   3324.62 → 3+3+2+4+6+2=20 → 20%8=4 → 返回4
//   3328.08 → 3+3+2+8+0+8=24 → 24%8=0 → 返回8(整除按8处理)
//   3300.01 → 3+3+0+0+0+1=7 → 7<8 → 返回8(小于8按8处理)
func CalculateRemainder(price float64) int {
    // 转换为固定格式的字符串(保留两位小数)
    priceStr := fmt.Sprintf("%.2f", price)

    // 累加所有数字位(忽略小数点)
    digitSum := int64(0)
    for _, char := range priceStr {
        if char != '.' {
            digit := int64(char - '0') // 字符转数字
            digitSum += digit
        }
    }

    // 对8取模获得余数
    remainder := digitSum % 8

    // 特殊情况:余数为0时按8处理
    if remainder == 0 {
        remainder = 8
    }

    return int(remainder)
}
```

---

## 模块化重构

### 建议文件结构

```
stock_predict/
├── main.go              # 程序入口
├── config/
│   ├── config.go        # 配置加载
│   └── config.json      # 配置文件
├── model/
│   ├── stock.go         # 数据结构定义
│   └── result.go        # 结果结构
├── data/
│   ├── loader.go        # 数据加载
│   ├── daily.go         # 日线数据
│   └── minute.go        # 分钟数据
├── algorithm/
│   ├── feature.go       # 特征计算
│   ├── match.go         # 历史匹配
│   ├── highlow.go       # 高低点识别
│   ├── trend.go         # 涨跌判断
│   └── similarity.go    # 相似度计算
├── image/
│   ├── generator.go     # 图像生成
│   └── merger.go        # 图像合并
├── utils/
│   ├── date.go          # 日期工具
│   ├── logger.go        # 日志工具
│   └── file.go          # 文件工具
└── docs/
    ├── README.md
    ├── API.md
    └── USAGE.md
```

### 模块拆分示例

#### 1. 配置模块 (config/config.go)

```go
package config

import (
    "encoding/json"
    "os"
)

// Config 系统配置
type Config struct {
    Data struct {
        DailyDataPath    string `json:"daily_data_path"`
        M1DataPath       string `json:"m1_data_path"`
        HistoryDataDir   string `json:"history_data_dir"`
        OutputDir        string `json:"output_dir"`
    } `json:"data"`

    Image struct {
        Width            int    `json:"width"`
        Height           int    `json:"height"`
        LineWidth        int    `json:"line_width"`
    } `json:"image"`

    Calculation struct {
        Modulo           int    `json:"modulo"`
        TradingMinutes   int    `json:"trading_minutes"`
        MorningEndMinute int    `json:"morning_end_minute"`
    } `json:"calculation"`
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    config := &Config{}
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(config); err != nil {
        return nil, err
    }

    return config, nil
}
```

#### 2. 数据模型 (model/stock.go)

```go
package model

// StockData 单根K线数据
type StockData struct {
    Date  string  // 日期
    Time  string  // 时间
    Open  float64 // 开盘价
    High  float64 // 最高价
    Low   float64 // 最低价
    Close float64 // 收盘价

    IsHighPoint bool // 高点标记
    IsLowPoint  bool // 低点标记
}

// StockDataSet K线数据集合
type StockDataSet struct {
    Data map[int]*StockData // 分钟索引 -> 数据

    MorningHigh    float64 // 上午最高价
    MorningLow     float64 // 上午最低价
    AfternoonHigh  float64 // 下午最高价
    AfternoonLow   float64 // 下午最低价

    PositiveCount  int     // 正向匹配数
    NegativeCount  int     // 反向匹配数
}

// TrendType 走势类型
type TrendType int

const (
    TrendUp    TrendType = 1  // 上涨
    TrendDown  TrendType = 2  // 下跌
    TrendFlat  TrendType = 3  // 震荡
)

// GetTrend 判断走势类型
func (s *StockDataSet) GetTrend() TrendType {
    if s.AfternoonHigh > s.MorningHigh && s.AfternoonLow > s.MorningLow {
        return TrendUp
    }
    if s.AfternoonHigh < s.MorningHigh && s.AfternoonLow < s.MorningLow {
        return TrendDown
    }
    return TrendFlat
}
```

#### 3. 特征计算 (algorithm/feature.go)

```go
package algorithm

import (
    "fmt"
    "strconv"
)

// FeatureExtractor 特征提取器
type FeatureExtractor struct {
    Modulo int // 取模基数,默认8
}

// NewFeatureExtractor 创建特征提取器
func NewFeatureExtractor(modulo int) *FeatureExtractor {
    return &FeatureExtractor{Modulo: modulo}
}

// ExtractFromPrice 从价格提取特征
func (fe *FeatureExtractor) ExtractFromPrice(price float64) int {
    // 转为字符串
    priceStr := fmt.Sprintf("%.2f", price)

    // 累加数字
    sum := int64(0)
    for _, char := range priceStr {
        if char >= '0' && char <= '9' {
            digit := int64(char - '0')
            sum += digit
        }
    }

    // 取模
    remainder := sum % int64(fe.Modulo)

    // 特殊处理
    if remainder == 0 {
        remainder = int64(fe.Modulo)
    }

    return int(remainder)
}

// ExtractFromOHLC 从OHLC提取特征组
func (fe *FeatureExtractor) ExtractFromOHLC(open, high, low, close float64) (int, int, int, int) {
    return fe.ExtractFromPrice(open),
           fe.ExtractFromPrice(high),
           fe.ExtractFromPrice(low),
           fe.ExtractFromPrice(close)
}
```

#### 4. 数据加载 (data/loader.go)

```go
package data

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"

    "stock_predict/model"
)

// DailyDataLoader 日线数据加载器
type DailyDataLoader struct {
    FilePath string
}

// Load 加载日线数据
func (l *DailyDataLoader) Load() ([]model.StockData, error) {
    file, err := os.Open(l.FilePath)
    if err != nil {
        return nil, fmt.Errorf("打开文件失败: %w", err)
    }
    defer file.Close()

    var data []model.StockData
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        fields := strings.Split(line, "\t")

        if len(fields) < 5 {
            continue
        }

        // 解析数据
        date := formatDate(fields[0])
        open, _ := strconv.ParseFloat(fields[1], 64)
        high, _ := strconv.ParseFloat(fields[2], 64)
        low, _ := strconv.ParseFloat(fields[3], 64)
        close, _ := strconv.ParseFloat(fields[4], 64)

        data = append(data, model.StockData{
            Date:  date,
            Open:  open,
            High:  high,
            Low:   low,
            Close: close,
        })
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("读取文件失败: %w", err)
    }

    return data, nil
}

// formatDate 格式化日期
func formatDate(dateStr string) string {
    // 将 2018/11/1 转为 2018-11-01
    parts := strings.Split(dateStr, "/")
    if len(parts) != 3 {
        return dateStr
    }

    year, _ := strconv.Atoi(parts[0])
    month, _ := strconv.Atoi(parts[1])
    day, _ := strconv.Atoi(parts[2])

    return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}
```

---

## 错误处理改进

### 改进前

```go
d, err := ioutil.ReadFile("data/日线.txt")
if err == nil {
    // 处理数据
}
```

### 改进后

```go
// LoadDailyData 加载日线数据,包含完善的错误处理
func LoadDailyData(path string) ([]StockData, error) {
    // 检查文件是否存在
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, fmt.Errorf("数据文件不存在: %s", path)
    }

    // 读取文件
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("读取文件失败: %w", err)
    }

    // 检查文件是否为空
    if len(data) == 0 {
        return nil, fmt.Errorf("数据文件为空: %s", path)
    }

    // 解析数据
    lines := strings.Split(string(data), "\n")
    var stockData []StockData

    for lineNum, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

        // 解析单行
        stock, err := parseStockLine(line)
        if err != nil {
            // 记录警告但继续处理
            log.Printf("警告: 第%d行数据解析失败: %v", lineNum+1, err)
            continue
        }

        stockData = append(stockData, stock)
    }

    // 检查是否有有效数据
    if len(stockData) == 0 {
        return nil, fmt.Errorf("未能解析出有效数据")
    }

    log.Printf("成功加载 %d 条日线数据", len(stockData))
    return stockData, nil
}
```

---

## 日志系统

### 添加日志模块 (utils/logger.go)

```go
package utils

import (
    "log"
    "os"
    "path/filepath"
)

// Logger 日志记录器
type Logger struct {
    infoLogger  *log.Logger
    errorLogger *log.Logger
    debugLogger *log.Logger
}

// NewLogger 创建日志记录器
func NewLogger(logDir string) (*Logger, error) {
    // 确保日志目录存在
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    // 创建日志文件
    logFile := filepath.Join(logDir, "app.log")
    file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }

    return &Logger{
        infoLogger:  log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
        errorLogger: log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
        debugLogger: log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
    }, nil
}

// Info 记录信息日志
func (l *Logger) Info(format string, v ...interface{}) {
    l.infoLogger.Printf(format, v...)
}

// Error 记录错误日志
func (l *Logger) Error(format string, v ...interface{}) {
    l.errorLogger.Printf(format, v...)
}

// Debug 记录调试日志
func (l *Logger) Debug(format string, v ...interface{}) {
    l.debugLogger.Printf(format, v...)
}
```

---

## 单元测试

### 测试示例 (algorithm/feature_test.go)

```go
package algorithm

import "testing"

func TestFeatureExtractor_ExtractFromPrice(t *testing.T) {
    extractor := NewFeatureExtractor(8)

    tests := []struct {
        name     string
        price    float64
        expected int
    }{
        {"普通情况", 3324.62, 4},
        {"整除情况", 3328.08, 8},
        {"小于基数", 3300.01, 8},
        {"最高价", 3352.08, 5},
        {"最低价", 3312.44, 1},
        {"收盘价", 3333.02, 6},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := extractor.ExtractFromPrice(tt.price)
            if result != tt.expected {
                t.Errorf("ExtractFromPrice(%f) = %d, 期望 %d",
                    tt.price, result, tt.expected)
            }
        })
    }
}
```

---

## 重构步骤建议

### 第一阶段:基础重构(不改变功能)

1. **添加注释**
   - 为所有函数添加文档注释
   - 为复杂逻辑添加行内注释

2. **重命名变量**
   - 使用清晰的英文命名
   - 遵循Go命名规范

3. **提取配置**
   - 创建config.json
   - 移除硬编码值

### 第二阶段:结构重构

4. **拆分文件**
   - 按功能划分模块
   - 创建包结构

5. **错误处理**
   - 添加完善的错误检查
   - 返回有意义的错误信息

6. **日志系统**
   - 添加日志记录
   - 便于调试和监控

### 第三阶段:功能增强

7. **单元测试**
   - 为核心算法编写测试
   - 确保重构不破坏功能

8. **性能优化**
   - 识别性能瓶颈
   - 优化关键路径

9. **文档完善**
   - 更新README
   - 编写API文档

---

## 迁移检查清单

重构完成后,检查以下项目:

- [ ] 所有函数都有文档注释
- [ ] 变量命名清晰易懂
- [ ] 配置项都在配置文件中
- [ ] 有完善的错误处理
- [ ] 有日志记录功能
- [ ] 代码按模块组织
- [ ] 通过所有测试用例
- [ ] 文档已更新
- [ ] 性能没有明显下降
- [ ] 功能与原代码一致

---

## 建议优先级

### 高优先级(必须完成)

1. 添加关键函数注释
2. 重命名核心变量
3. 添加错误处理
4. 创建配置文件

### 中优先级(强烈建议)

5. 拆分模块
6. 添加日志系统
7. 编写测试用例

### 低优先级(可选)

8. 性能优化
9. 功能扩展
10. UI改进

---

**注意:** 重构应逐步进行,每次改动后都要测试确保功能正常!

---

**最后更新:** 2025-11-11
