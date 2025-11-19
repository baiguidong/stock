// main.go
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Params holds all the UI parameters (可扩展)
type Params struct {
	Date        string `json:"date"`
	Offset1     int    `json:"offset1"`
	Offset2     int    `json:"offset2"`
	Local       int    `json:"local"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	LineWidth   int    `json:"line_width"`
	Count       int    `json:"count"`
	LocalWidth  int    `json:"local_width"`
	LocalHeight int    `json:"local_height"`
	LocalLine   int    `json:"local_line"`
	LastUpdated string `json:"last_updated"`
}

var select_day *widget.Entry
var select_status *widget.Label

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("比对工具")
	w.Resize(fyne.NewSize(520, 420))
	// ---- 状态条 ----
	select_status = widget.NewLabel("")
	// 时间
	select_day = widget.NewEntry()
	select_day.SetText(time.Now().Format("2006-01-02"))

	// 偏移（两个）
	offsetStart := widget.NewEntry()
	offsetStart.SetPlaceHolder("开始")
	offsetStart.SetText("0")

	// 偏移 截止（只能输入数字）
	offsetEnd := widget.NewEntry()
	offsetEnd.SetPlaceHolder("截止")
	offsetEnd.SetText("0")

	offsetRow := container.NewVBox(
		container.NewHBox(
			offsetStart,
			widget.NewLabel(" - "),
			offsetEnd,
		),
	)
	// 局部（下拉）
	localSelect := widget.NewSelect([]string{"0", "30", "60", "120", "240"}, func(s string) {})
	localSelect.SetSelected("0")

	// 宽 / 高 / 线条 / 张数
	widthEntry := widget.NewEntry()
	widthEntry.SetText("600")
	heightEntry := widget.NewEntry()
	heightEntry.SetText("200")
	lineEntry := widget.NewEntry()
	lineEntry.SetText("3")
	countEntry := widget.NewEntry()
	countEntry.SetText("100")

	// 局图：宽 高 线
	localWidth := widget.NewEntry()
	localWidth.SetText("400")
	localHeight := widget.NewEntry()
	localHeight.SetText("200")
	localline := widget.NewEntry()
	localline.SetText("20")

	// 输出/智能按钮
	btnOut1 := widget.NewButton("全数据计算", nil)
	btnOut1.Importance = widget.HighImportance
	btnOut2 := widget.NewButton("多组计算", nil)
	btnOut2.Importance = widget.HighImportance

	// 统一按钮最小宽度（让界面看着整齐）
	minBtnWidth := float32(160)
	btnOut1.Resize(fyne.NewSize(minBtnWidth, 64))
	btnOut2.Resize(fyne.NewSize(minBtnWidth, 64))

	buttonRow1 := container.NewHBox(
		layout.NewSpacer(),
		btnOut1,
		layout.NewSpacer(),
		btnOut2,
		layout.NewSpacer(),
	)

	// ---- 帮助函数 ----
	collectParams := func() (*Params, error) {
		// convert entries
		p := &Params{}
		p.Date = select_day.Text
		if v, err := strconv.Atoi(offsetStart.Text); err == nil {
			p.Offset1 = v
		}
		if v, err := strconv.Atoi(offsetEnd.Text); err == nil {
			p.Offset2 = v
		}
		if v, err := strconv.Atoi(localSelect.Selected); err == nil {
			p.Local = v
		}
		if v, err := strconv.Atoi(widthEntry.Text); err == nil {
			p.Width = v
		}
		if v, err := strconv.Atoi(heightEntry.Text); err == nil {
			p.Height = v
		}
		if v, err := strconv.Atoi(lineEntry.Text); err == nil {
			p.LineWidth = v
		}
		if v, err := strconv.Atoi(countEntry.Text); err == nil {
			p.Count = v
		}
		if v, err := strconv.Atoi(localWidth.Text); err == nil {
			p.LocalWidth = v
		}
		if v, err := strconv.Atoi(localHeight.Text); err == nil {
			p.LocalHeight = v
		}
		if v, err := strconv.Atoi(localline.Text); err == nil {
			p.LocalLine = v
		}
		p.LastUpdated = time.Now().Format(time.RFC3339)
		return p, nil
	}

	// 生成图片并保存（用于输出按钮）
	generateAndSaveImage := func(signal int) {
		p, _ := collectParams()
		if signal == 0 {
			p.Offset2 = 0
		}
		g_msg <- "正在生成: "
		ui_start(p)
	}

	btnOut1.OnTapped = func() { generateAndSaveImage(0) }
	btnOut2.OnTapped = func() { generateAndSaveImage(1) }

	// ---- 美化布局 ----
	// 左标签列宽（固定）与右控件排列更像 WinForm
	leftCol := container.NewVBox(
		widget.NewLabelWithStyle("时间", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("偏移", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("局部", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("宽度", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("高度", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("线条", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("张数", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("局图", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	// 右列控件（按行）
	rightCol := container.NewVBox(
		select_day,
		offsetRow,
		localSelect,
		widthEntry,
		heightEntry,
		lineEntry,
		countEntry,
		container.NewHBox(widget.NewLabel("宽"), localWidth, widget.NewLabel("高"), localHeight, widget.NewLabel("线"), localline),
	)

	// 表单两列组合
	form := container.NewHBox(
		layout.NewSpacer(),
		container.NewVBox(leftCol),
		layout.NewSpacer(),
		container.NewVBox(rightCol),
		layout.NewSpacer(),
	)

	// 统一控件卡片（类似 GroupBox）
	card := container.NewVBox(
		form,
		widget.NewSeparator(),
		layout.NewSpacer(),
	)

	// 最后组合
	content := container.NewVBox(
		card,
		buttonRow1,
		layout.NewSpacer(),
		select_status,
	)

	// 加边距
	padded := container.New(layout.NewVBoxLayout(), container.NewPadded(content))

	w.SetContent(padded)
	go Init_UI()
	w.ShowAndRun()
}

func Init_UI() {
	go func() {
		for {
			value, ok := <-g_msg
			if !ok {
				break // channel 已关闭，退出循环
			}
			if select_status != nil {
				fyne.Do(func() {
					{
						select_status.SetText(value)
					}
				})
			}
		}
	}()
	verticalImages = []string{}
	load_data()
	g_msg <- fmt.Sprintf("历史数据:%d天\n", len(g_map.m_map))

	lastDay := load_M1()
	fmt.Println(lastDay)
	g_msg <- fmt.Sprintf("M1数据:%d天\n", len(g_map_m1.m_map))
	load_day()
	g_msg <- fmt.Sprintf("日线数据:%d\n", len(days))
	fyne.Do(func() {
		select_day.SetText(lastDay)
	})
}
func ui_start(param *Params) {
	os.RemoveAll("tmp")
	os.Mkdir("tmp", 0777)
	new_data = ass{}

	g_height = param.Height
	g_width = param.Width
	g_xt = param.LineWidth
	g_num = 40
	g_ks = param.LocalWidth
	g_gs = param.LocalHeight
	g_local_vline = param.LocalLine
	g_day = param.Date
	g_offset_start = param.Offset1
	g_offset_end = param.Offset2
	g_local = param.Local
	imageNums = param.Count

	if g_height > 0 {
		h_head = g_height / 2
	}

	if g_offset_end <= g_offset_start {
		g_offset_end = g_offset_start
	}
	verticalImages = []string{}
	// g_offset_start -> g_offset_end 所有offset
	for i := g_offset_start; i <= g_offset_end; i++ {
		// 计算逻辑
		run_offset(i)
	}

	g_msg <- "正在生成照片,请稍等..."
	generate_vimages()
	g_msg <- "正在合成成照片,请稍等..."
	generate_result()
	g_msg <- "完成"
	os.RemoveAll("tmp")
}
