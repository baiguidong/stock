// main.go
package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type CustomTheme struct {
	fyne.Theme
}

func (c CustomTheme) ButtonColor() color.Color {
	return color.White
}

func (c CustomTheme) HoverColor() color.Color {
	return color.RGBA{230, 230, 230, 255} // 按钮 hover 时的灰色
}

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

// createSampleImage 创建一个示例 PNG，用来演示“输出原图/新图”等操作
func createSampleImage(path string, w, h int, bg color.RGBA, title string) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// 背景
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
	// 在左上角画一些条纹来模拟“线条”
	for i := 0; i < 10; i++ {
		for x := 0; x < w; x++ {
			y := (i*20 + 5) % h
			img.Set(x, (y+i)%h, color.RGBA{255, 255, 255, 30})
		}
	}
	// 写入文件
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return err
	}
	return nil
}

var select_day *widget.Entry
var select_status *widget.Label

func main() {

	a := app.New()
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
	countEntry.SetText("5")

	// 局图：宽 高 线
	localWidth := widget.NewEntry()
	localWidth.SetText("400")
	localHeight := widget.NewEntry()
	localHeight.SetText("200")
	localline := widget.NewEntry()
	localline.SetText("20")

	// 输出/智能按钮
	btnOut1 := widget.NewButton("输出原图", nil)

	// 统一按钮最小宽度（让界面看着整齐）
	minBtnWidth := float32(160)
	btnOut1.Resize(fyne.NewSize(minBtnWidth, 64))

	buttonRow1 := container.NewHBox(
		layout.NewSpacer(),
		btnOut1,
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
	generateAndSaveImage := func(title string) {
		p, _ := collectParams()
		g_msg <- "正在生成: "
		ui_start(p)
	}

	btnOut1.OnTapped = func() { generateAndSaveImage("输出原图") }

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
	)

	// 最后组合
	content := container.NewVBox(
		card,
		widget.NewSeparator(),
		buttonRow1,
		layout.NewSpacer(),
		select_status,
	)

	// 加边距
	padded := container.New(layout.NewVBoxLayout(), container.NewPadded(content))

	w.SetContent(padded)
	go Init()
	w.ShowAndRun()
}
