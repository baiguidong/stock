package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/shakinm/xlsReader/xls"
	"github.com/shakinm/xlsReader/xls/structure"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
)

func load_day() {
	if Isexist("data/日线.xls") {
		workbook, err := xls.OpenFile("data/日线.xls")
		if err != nil {
			log.Fatalf("打开 XLS 文件失败: %v", err)
		}

		// 遍历第一个 sheet（你可以遍历所有）
		sheet, err := workbook.GetSheet(0)
		if err != nil {
			log.Fatalf("获取 Sheet 失败: %v", err)
		}

		for i := 1; i <= int(sheet.GetNumberRows()); i++ {
			row, err := sheet.GetRow(i)
			if err != nil {
				continue
			}
			if row == nil {
				continue
			}
			day, err := row.GetCol(0)
			if err != nil {
				continue
			}
			if strings.Trim(day.GetString(), " ") == "" {
				break
			}
			day_tm, _ := parseExcelDateTime(day)

			open, err := row.GetCol(1)
			if err != nil {
				continue
			}
			high, err := row.GetCol(2)
			if err != nil {
				continue
			}
			low, err := row.GetCol(3)
			if err != nil {
				continue
			}
			close, err := row.GetCol(4)
			if err != nil {
				continue
			}
			ss := sed{dy: day_tm.Format("2006-01-02"), dk: open.GetFloat64(), dg: high.GetFloat64(), dd: low.GetFloat64(), ds: close.GetFloat64()}
			days = append(days, ss)
		}
		return
	}

	d, err := os.ReadFile("data/日线.txt")
	if err == nil {
		v2 := strings.Split(string(d), "\n")
		for _, v3 := range v2 {
			v3 = strings.Replace(v3, "\r", "", -1)
			v3 = strings.Replace(v3, " ", "", -1)
			vs := strings.Split(v3, "\t")
			if len(vs) >= 5 {
				a2, err := strconv.ParseFloat(vs[1], 10)
				if err != nil {
					continue
				}
				a3, err := strconv.ParseFloat(vs[2], 10)
				if err != nil {
					continue
				}
				a4, err := strconv.ParseFloat(vs[3], 10)
				if err != nil {
					continue
				}
				a5, err := strconv.ParseFloat(vs[4], 10)
				if err != nil {
					continue
				}
				ss := sed{dy: format_tm(vs[0]), dk: a2, dg: a3, dd: a4, ds: a5}
				days = append(days, ss)
			}
		}
	}
}
func get_val(k string) int {
	return g_dz[k]
}

func get_hes(f1 float64, f2 float64, f3 float64, f4 float64) (int, int, int, int) {
	return get_he(f1), get_he(f2), get_he(f3), get_he(f4)
}

// 转换 int  对 8 取余数
func get_he(f float64) int {
	s := fmt.Sprintf("%.2f", f)
	//fmt.Println(s)
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
func format_tm(s string) string {
	v := strings.Split(s, "/")
	des := ""
	if len(v) == 3 {
		v1, _ := strconv.ParseInt(v[0], 10, 64)
		v2, _ := strconv.ParseInt(v[1], 10, 10)
		v3, _ := strconv.ParseInt(v[2], 10, 10)

		des = fmt.Sprintf("%d-%02d-%02d", v1, v2, v3)
	}
	return des
}

type ass struct {
	data []sds
}
type ass_1 struct {
	data []sds_2
}

func (c ass_1) Len() int {
	return len(c.data)
}
func (c ass) Len() int {
	return len(c.data)
}

func (c ass) Swap(i, j int) {
	c.data[i], c.data[j] = c.data[j], c.data[i]
}

// func (c ass) Less(i, j int) bool {
// 	return c.data[i].osjl < c.data[j].osjl
// }

func get_m1(day string, offset int) (*seds, error) {
	if offset >= 240 {
		fmt.Println("ERR:偏移错误 ", offset)
		return nil, errors.New("ERR:偏移错误")
	}

	m1 := g_map_m1.m_map[day]
	if m1 == nil {
		fmt.Println("ERR:获取数据失败", day)
		return nil, errors.New("ERR:获取数据失败")
	}
	var before_m1 *seds
	if offset > 0 && offset < 240 {
		mm := 15
		beforeDay := day
		for mm > 1 {
			beforeDay = before_tm(beforeDay)
			before_m1 = g_map_m1.m_map[beforeDay]
			if before_m1 != nil {
				break
			}
			mm--
		}
		if before_m1 == nil {
			fmt.Println("ERR:获取数据失败——偏移数据")
			os.Exit(-1)
		}
	}

	em := &sed{}
	// rs 为前一天 240 + 当前日期的 240个点
	rs := []*sed{}
	i := 1
	for i < 241 {
		if before_m1 == nil {
			rs = append(rs, em)
		} else {
			vv := before_m1.vec[i]
			if vv != nil {
				rs = append(rs, vv)
			} else {
				rs = append(rs, em)
			}
		}
		i++
	}
	i = 1
	for i < 241 {
		vv := m1.vec[i]
		if vv != nil {
			rs = append(rs, vv)
		} else {
			rs = append(rs, em)
		}
		i++
	}
	m := &seds{}
	m.init_map()

	// 算出偏移后的 240 个点
	num := 240 - offset
	num_1 := 480 - offset
	xl := 1
	for num < num_1 {
		m.vec[xl] = rs[num]
		num++
		xl++
	}
	return m, nil
}

// 输出当前日期 偏移的点(没参与计算，x从0 开始)
func get_m1_py(dy string, py int) plotter.XYs {
	res := plotter.XYs{}
	data := g_map_m1.m_map[dy]

	num := 240 - py
	nn := num
	if data != nil {
		for num < 241 {
			point := data.vec[num]
			if point != nil {
				res = append(res,
					plotter.XY{X: float64(num - nn), Y: point.ds})
			} else {
				res = append(res,
					plotter.XY{X: float64(num - nn), Y: 0.0})
			}
			num++
		}
	}
	return res
}

// 指定日期 收盘价的 240 个点
func GetPointsByDay(day string) plotter.XYs {
	res := plotter.XYs{}
	data := g_map.m_map[day]
	if data != nil {
		n := 1
		for n <= 240 {
			point := data.vec[n]
			if point != nil {
				res = append(res,
					plotter.XY{X: float64(n), Y: point.ds})
			} else {
				res = append(res,
					plotter.XY{X: float64(n), Y: 0.0})
			}
			n++
		}
	}
	return res
}

// seds 转换 XYS 点
func CVPoints(s *seds) plotter.XYs {
	res := plotter.XYs{}
	if s != nil {
		n := 1
		for n <= 240 {
			point := s.vec[n]
			if point != nil {
				res = append(res,
					plotter.XY{X: float64(n), Y: point.ds})
			} else {
				res = append(res,
					plotter.XY{X: float64(n), Y: 0.0})
			}
			n++
		}
	}
	return res
}
func next_tm(s string) string {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err.Error()
	}

	if t.Weekday().String() == "Friday" {
		t = t.AddDate(0, 0, 3)
	} else if t.Weekday().String() == "Saturday" {
		t = t.AddDate(0, 0, 2)
	} else {
		t = t.AddDate(0, 0, 1)
	}
	return t.Format("2006-01-02")
}

// 获取前一个交易日
func before_tm(day string) string {
	t, err := time.Parse("2006-01-02", day)
	if err != nil {
		return err.Error()
	}

	if t.Weekday().String() == "Monday" {
		t = t.AddDate(0, 0, -3)
	} else {
		t = t.AddDate(0, 0, -1)
	}
	return t.Format("2006-01-02")
}

type sds struct {
	//dt1 前一天
	//dt2 后一天
	dt1, dt2 plotter.XYs
	// dy 当天 dy_1 前一天 dy_2 后一天
	dy, dy_1, dy_2 string
	zh, fa         int
	osjl           float64
	py             int
	sm             *seds
}
type sds_2 struct {
	dt1, dt2, dt3 plotter.XYs
	n1, n2, n3    int
	o1, o2        float64

	num int
}

func open_pic(name string) {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("open", name)
		cmd.Start() // 不阻塞
	} else {
		cmd := fmt.Sprintf("%s rundll32.exe  C:/windows/system/shimgvw.dll", name)
		_, err := Runcmd(cmd)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func load_xls(v1 string) {
	v1_cache := strings.Replace(v1, "data/", "data/cache/", -1)
	// cache 文件
	v1_cache = strings.Replace(v1_cache, ".xls", ".csv", -1)
	if Isexist(v1_cache) {
		load_csv(v1_cache)
		return
	}

	workbook, err := xls.OpenFile(v1)
	if err != nil {
		log.Fatalf("打开 XLS 文件失败: %v", err)
	}
	// 遍历第一个 sheet（你可以遍历所有）
	sheet, err := workbook.GetSheet(0)
	if err != nil {
		log.Fatalf("获取 Sheet 失败: %v", err)
	}
	os.MkdirAll("data/cache", 0777)
	file, err := os.Create(v1_cache)
	if err != nil {
		fmt.Println("创建文件时出错:", err)
		return
	}
	defer file.Close() // 确保在程序结束前关闭文件

	for i := 1; i <= int(sheet.GetNumberRows()); i++ {
		if i%1000 == 0 {
			g_msg <- fmt.Sprintf("正在加载数据 %s-%d", v1, i)
		}
		row, err := sheet.GetRow(i)
		if err != nil {
			continue
		}

		day, err := row.GetCol(0)
		if err != nil {
			continue
		}
		if strings.Trim(day.GetString(), " ") == "" {
			break
		}
		day_str, _ := parseExcelDate(day)
		tm, err := row.GetCol(1)
		if err != nil {
			continue
		}
		tm_str := parseExcelTime(tm)

		open, err := row.GetCol(2)
		if err != nil {
			continue
		}
		high, err := row.GetCol(3)
		if err != nil {
			continue
		}
		low, err := row.GetCol(4)
		if err != nil {
			continue
		}
		close, err := row.GetCol(5)
		if err != nil {
			continue
		}
		day_str_csv := strings.ReplaceAll(day_str, "-", "/")
		file.WriteString(fmt.Sprintf("%s,%s,%g,%g,%g,%g\n", day_str_csv, tm_str, open.GetFloat64(), high.GetFloat64(), low.GetFloat64(), close.GetFloat64()))
		g_map.add_data(day_str, tm_str, open.GetFloat64(), high.GetFloat64(), low.GetFloat64(), close.GetFloat64())
	}
}
func load_csv(v1 string) {
	f, err := os.Open(v1)
	if err != nil {
		return
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	for {
		v3, err := buf.ReadString('\n')
		if err == nil {
			v3 = strings.Replace(v3, "\r", "", -1)
			v3 = strings.Replace(v3, "\n", "", -1)
			v3 = strings.Replace(v3, " ", "", -1)
			vs := strings.Split(v3, ",")
			if len(vs) >= 6 {
				a2, err := strconv.ParseFloat(vs[2], 10)
				if err != nil {
					continue
				}
				a3, err := strconv.ParseFloat(vs[3], 10)
				if err != nil {
					continue
				}
				a4, err := strconv.ParseFloat(vs[4], 10)
				if err != nil {
					continue
				}
				a5, err := strconv.ParseFloat(vs[5], 10)
				if err != nil {
					continue
				}
				g_map.add_data(format_tm(vs[0]), vs[1], a2, a3, a4, a5)
			}
		} else {
			break
		}
	}
}
func load_data() {
	v, err := List_file("data", "", 9999)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	for _, v1 := range v {
		g_msg <- fmt.Sprintf("正在加载数据 %s", v1)
		if v1 == "data/M1.txt" || v1 == "data/M1.xls" {
			continue
		}
		if v1 == "data/日线.txt" || v1 == "data/日线.xls" {
			continue
		}
		if strings.HasSuffix(v1, ".xls") {
			load_xls(v1)
			g_msg <- fmt.Sprintf("加载完成 %s", v1)
		} else {
			load_csv(v1)
			g_msg <- fmt.Sprintf("加载完成 %s", v1)
		}
	}
}

// 分钟数据
type sed struct {
	// 日期
	dy string
	// 时间
	tm string
	// 开 高 低 收
	dk, dg, dd, ds float64
	// 高点
	g int
	// 低点
	d int
}

// 一天的数据
type seds struct {
	// 240 个点数据
	vec map[int]*sed
	// 正 / 反 数量
	z, f int
	// 上午高点
	sg float64
	// 上午低点
	sd float64
	// 下午高点
	xg float64
	// 下午低点
	xd float64
}

func (m *seds) Add(s *sed) {
	m.vec[get_val(s.tm)] = s
}
func (m *seds) isz() bool {
	if m.xg > m.sg && m.xd > m.sd {
		return true
	}
	return false
}
func (m *seds) isd() bool {
	if m.sg > m.xg && m.sd > m.xd {
		return true
	}
	return false
}

// 获得上午/下午 高低点
func (m *seds) zd() {
	a := 1
	for a < 241 {
		aa := m.vec[a]
		if aa == nil {
			return
		}
		a++
	}
	i := 1
	m.sg = m.vec[1].ds
	m.sd = m.vec[1].ds

	m.xg = m.vec[121].ds
	m.xd = m.vec[121].ds
	for i < 241 {
		if i < 121 {
			if m.vec[i].ds > m.sg {
				m.sg = m.vec[i].ds
			}
			if m.vec[i].ds < m.sd {
				m.sd = m.vec[i].ds
			}
		} else {
			if m.vec[i].ds > m.xg {
				m.xg = m.vec[i].ds
			}
			if m.vec[i].ds < m.xd {
				m.xd = m.vec[i].ds
			}
		}
		i++
	}
}
func (m *seds) gd() (int, int) {
	a := 1
	for a < 241 {
		aa := m.vec[a]
		if aa == nil {
			return 0, 0
		}
		a++
	}
	i := 2
	g := 0
	d := 0
	if m.vec[1].ds > m.vec[2].ds {
		m.vec[1].g = 1
		g += 1
	} else {
		m.vec[1].d = 1
		d += 1
	}
	if m.vec[240].ds > m.vec[239].ds {
		m.vec[240].g = 1
		g += 1
	} else {
		m.vec[240].d = 1
		d += 1
	}
	for i < 240 {
		a := m.vec[i-1]
		b := m.vec[i]
		c := m.vec[i+1]
		if a == nil || b == nil || c == nil {
			i++
			continue
		}
		if b.ds > a.ds && b.ds > c.ds {
			b.g = 1
			g += 1
		}
		if b.ds < a.ds && b.ds < c.ds {
			b.d = 1
			d += 1
		}
		i++
	}

	return g, d
}

// 根据分钟 数据计算出 开盘价/高点/低点/收盘价
func (m *seds) deal() (float64, float64, float64, float64) {
	k := m.vec[1].dk
	s := m.vec[240].ds
	g := m.vec[1].dg
	d := m.vec[1].dd
	for _, v := range m.vec {
		if v.dg > g {
			g = v.dg
		}
		if v.dd < d {
			d = v.dd
		}
	}
	return k, g, d, s
}

func (m *seds) init_map() {
	m.vec = make(map[int]*sed)
}

type G_map struct {
	m_map map[string]*seds
}

func (g *G_map) init_data() {
	g.m_map = make(map[string]*seds)
}

func (g *G_map) add_data(k, tm string, dk, dg, dd, ds float64) {
	s := &sed{dy: k, tm: tm, dk: dk, dg: dg, dd: dd, ds: ds}
	if _, ok := g.m_map[k]; ok {
		g.m_map[k].Add(s)
	} else {
		mu := seds{}
		mu.init_map()
		mu.Add(s)

		g.m_map[k] = &mu
	}
}

func load_M1() string {
	lastday := ""
	if Isexist("data/M1.xls") {
		workbook, err := xls.OpenFile("data/M1.xls")
		if err != nil {
			log.Fatalf("打开 XLS 文件失败: %v", err)
		}
		// 遍历第一个 sheet（你可以遍历所有）
		sheet, err := workbook.GetSheet(0)
		if err != nil {
			log.Fatalf("获取 Sheet 失败: %v", err)
		}
		for i := 1; i <= int(sheet.GetNumberRows()); i++ {
			row, err := sheet.GetRow(i)
			if err != nil {
				continue
			}
			day, err := row.GetCol(0)
			if err != nil {
				continue
			}
			day_str := strings.Trim(day.GetString(), " ")
			if day_str == "" {
				break
			}
			day_tm, _ := parseExcelDateTime(day)

			open, err := row.GetCol(1)
			if err != nil {
				continue
			}
			high, err := row.GetCol(2)
			if err != nil {
				continue
			}
			low, err := row.GetCol(3)
			if err != nil {
				continue
			}
			close, err := row.GetCol(4)
			if err != nil {
				continue
			}

			g_map_m1.add_data(day_tm.Format("2006-01-02"), day_tm.Format("15:04"), open.GetFloat64(), high.GetFloat64(), low.GetFloat64(), close.GetFloat64())
			lastday = day_tm.Format("2006-01-02")
		}
		return lastday
	} else {
		d, err := os.ReadFile("data/M1.txt")
		if err == nil {
			v2 := strings.Split(string(d), "\n")
			for _, v3 := range v2 {
				v3 = strings.Replace(v3, "\r", "", -1)
				v3 = strings.Replace(v3, " ", "", -1)
				vs := strings.Split(v3, "\t")
				if len(vs) >= 5 {
					a2, err := strconv.ParseFloat(vs[1], 10)
					if err != nil {
						continue
					}
					a3, err := strconv.ParseFloat(vs[2], 10)
					if err != nil {
						continue
					}
					a4, err := strconv.ParseFloat(vs[3], 10)
					if err != nil {
						continue
					}
					a5, err := strconv.ParseFloat(vs[4], 10)
					if err != nil {
						continue
					}

					vs1 := strings.Split(vs[0], "-")
					if len(vs1) == 2 {
						g_map_m1.add_data(format_tm(vs1[0]), vs1[1], a2, a3, a4, a5)
					}
					lastday = format_tm(vs1[0])
				}
			}
		}
	}
	return lastday
}

func test_hz(text, text1, fname string) {
	fontBytes, err := ioutil.ReadFile("FZBSJW.ttf")
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, 800, 200))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetFont(f)
	c.SetFontSize(45)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetDPI(120)

	pt := freetype.Pt(0, 80)
	pt1 := freetype.Pt(0, 160)
	_, err = c.DrawString(text, pt)
	_, err = c.DrawString(text1, pt1)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Save that RGBA image to disk.
	outFile, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = jpeg.Encode(b, rgba, nil)
	if err != nil {
		fmt.Println(err)
	}
	err = b.Flush()
	if err != nil {
		fmt.Println(err)
	}

	s1, err := imaging.Open(fname)
	if err == nil {
		s1 = imaging.Resize(s1, g_width, h_head, imaging.Lanczos)
	}
	imaging.Save(s1, fname)

}

func merg_image(srcs []string, des string) {
	lens := len(srcs)

	emptynum := 0
	for _, src := range srcs {
		if src == "empty.png" {
			emptynum++
		}
	}
	lens = lens - emptynum
	ww := (g_width+2)*lens + g_width/2 + 10*emptynum
	// if jb >= 30 {
	// 	ww = (g_width+2)*lens/2 + (g_ks+2)*lens/2 + a_head + 10*emptynum
	// }
	//ww := (g_width+2)*lens + a_head + 10*emptynum
	hh := g_height*3 + h_head + g_width

	dst := imaging.New(ww, hh, color.NRGBA{255, 255, 255, 255})
	a := 0
	for n, src := range srcs {
		if src == "empty.png" {
			a++
		} else {
			s1, err := imaging.Open(src)
			if err == nil {
				// if jb >= 30 {
				// 	dst = imaging.Paste(dst, s1, image.Pt(a_head/2+(n-a)*((width+g_ks)/2+2)+a*10, 0))
				// } else {
				dst = imaging.Paste(dst, s1, image.Pt(g_width/4+(n-a)*(g_width+2)+a*10, 0))
				//}
			}
		}
	}
	// Save the resulting image as JPEG.
	err := imaging.Save(dst, des)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func merg_image_h(srcs []string, des string) {
	lens := len(srcs)
	dst := imaging.New(g_width, g_height*(lens-1)+4*h_head, color.NRGBA{255, 255, 255, 255})
	for n, src := range srcs {
		s1, err := imaging.Open(src)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if n == 0 {
				s1 = imaging.Resize(s1, g_width, h_head, imaging.Linear)
				dst = imaging.Paste(dst, s1, image.Pt(0, h_head))
			} else {
				s1 = imaging.Resize(s1, g_width, g_height, imaging.Linear)
				switch n {
				case 1:
					dst = imaging.Paste(dst, s1, image.Pt(0, 2*h_head+4*n))
				case 2:
					dst = imaging.Paste(dst, s1, image.Pt(0, 2*h_head+4*n+g_height))
				case 3:
					dst = imaging.Paste(dst, s1, image.Pt(0, 2*h_head+4*n+2*g_height))
				case 4:
					dst = imaging.Paste(dst, s1, image.Pt(0, 2*h_head+4*n+3*g_height))
				}
			}
		}
	}

	// Save the resulting image as JPEG.
	err := imaging.Save(dst, des)
	if err != nil {
		fmt.Println(err.Error())
	}
}
func get_home_dir() string {
	u, err := user.Current()
	if err == nil {
		return u.HomeDir //+ "/临时结果"
	} else {
		return ""
	}
}

func merg_image_offset(srcs []string, des string) {
	dst := imaging.New(g_ks, g_height*(3)+4*h_head, color.NRGBA{255, 255, 255, 255})
	for n, src := range srcs {
		s1, err := imaging.Open(src)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if n == 0 {
				s1 = imaging.Resize(s1, g_ks, g_gs, imaging.Linear)
				dst = imaging.Paste(dst, s1, image.Pt(0, 2*h_head+4+g_height/2))
			} else {
				s1 = imaging.Resize(s1, g_ks, g_gs, imaging.Linear)
				if n == 1 {
					//n*400-200+4*n
					dst = imaging.Paste(dst, s1, image.Pt(0, 2*h_head+8+g_gs+g_height/2))
				}
			}
		}
	}

	// Save the resulting image as JPEG.
	err := imaging.Save(dst, des)
	if err != nil {
		fmt.Println(err.Error())
	}
}
func get_line(w, mi, mx float64) *plotter.Line {
	c, _ := plotter.NewLine(plotter.XYs{
		{w, mi},
		{w, mx},
	})

	if w == 120.0 || w == 0.0 || w == 240.0 {
		c.LineStyle.Width = vg.Points(2)
		c.LineStyle.Color = color.RGBA{R: 255, A: 0, G: 0, B: 30}
	} else {
		c.LineStyle.Width = vg.Points(1)
		c.LineStyle.Color = color.RGBA{R: 255, A: 0}
		c.LineStyle.Dashes = []vg.Length{2, 2}
	}
	return c
}
func get_line_new(w, mi, mx float64) *plotter.Line {
	c, _ := plotter.NewLine(plotter.XYs{
		{w, mi},
		{w, mx},
	})

	c.LineStyle.Width = vg.Points(1)
	c.LineStyle.Color = color.RGBA{R: 255, A: 0, G: 0, B: 30}
	c.LineStyle.Dashes = []vg.Length{1, 1}
	return c
}

// 生成 image
func get_p(dt plotter.XYs, name, title string) {
	p := plot.New()
	p.Title.Text = title
	p.Title.TextStyle = text.Style{
		Color: color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}

	p.HideX()

	p.X.Padding = 0
	p.Y.Padding = 0
	p.HideY()

	p.HideAxes()

	p.BackgroundColor = color.RGBA{R: 0, A: 255}

	l, err := plotter.NewLine(dt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	l.LineStyle.Width = vg.Points(float64(g_xt))
	l.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	mx := 0.0
	mi := 0.0
	for _, p := range dt {
		if mi == 0.0 {
			mi = p.Y
		}
		if p.Y > mx {
			mx = p.Y
		}
		if p.Y < mi {
			mi = p.Y
		}
	}
	p.Add(
		get_line(0, mi, mx),
		get_line(30, mi, mx),
		get_line(60, mi, mx),
		get_line(90, mi, mx),
		get_line(120, mi, mx),
		get_line(150, mi, mx),
		get_line(180, mi, mx),
		get_line(210, mi, mx),
		get_line(240, mi, mx), l)

	if err := p.Save(vg.Length(g_width), vg.Length(g_height), name); err != nil {
		panic(err)
	}
}

// 生成 局部的图
func get_p_120(dt plotter.XYs, name, title string) {
	p := plot.New()
	p.Title.Text = title
	p.Title.TextStyle = text.Style{
		Color: color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}

	p.HideX()

	p.X.Padding = 0
	p.Y.Padding = 0
	p.HideY()

	p.HideAxes()

	p.BackgroundColor = color.RGBA{R: 0, A: 255}
	dt_new := plotter.XYs{}
	for n, x := range dt {
		if n >= 240-g_local {
			x.X = x.X - 240.0 + float64(g_local)
			dt_new = append(dt_new, x)
		}
	}
	l, err := plotter.NewLine(dt_new)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	l.LineStyle.Width = vg.Points(float64(g_xt))
	l.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	mx := 0.0
	mi := 0.0
	for _, p := range dt_new {
		if mi == 0.0 {
			mi = p.Y
		}
		if p.Y > mx {
			mx = p.Y
		}
		if p.Y < mi {
			mi = p.Y
		}
	}
	// 局部 g_local 点
	// g_local_vline 竖向线条数量
	xx := g_local / g_local_vline
	for a := 0; a <= g_local; a++ {
		if a%xx == 0 {
			// 画中间红色竖线
			p.Add(get_line_new(float64(a), mi, mx))
		}
	}
	p.Add(l)
	if err := p.Save(vg.Length(g_ks), vg.Length(g_gs), name); err != nil {
		panic(err)
	}
}

// 基于另一个最大最小
func get_p_py(dt, dt1 plotter.XYs, name, title string) {
	p := plot.New()

	p.Title.Text = title
	p.Title.TextStyle = text.Style{
		Color: color.RGBA{R: 255, G: 255, B: 255, A: 255},
	}

	p.HideX()

	p.X.Padding = 0
	p.Y.Padding = 0
	p.HideY()

	p.HideAxes()

	p.BackgroundColor = color.RGBA{R: 0, A: 255}

	l, err := plotter.NewLine(dt1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	l.LineStyle.Width = vg.Points(float64(g_xt))
	l.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	mx := 0.0
	mi := 0.0
	for _, p := range dt {
		if mi == 0.0 {
			mi = p.Y
		}
		if p.Y > mx {
			mx = p.Y
		}
		if p.Y < mi {
			mi = p.Y
		}
	}
	for _, p := range dt1 {
		if mi == 0.0 {
			mi = p.Y
		}
		if p.Y > mx {
			mx = p.Y
		}
		if p.Y < mi {
			mi = p.Y
		}
	}
	p.Add(
		get_line(0, mi, mx),
		get_line(30, mi, mx),
		get_line(60, mi, mx),
		get_line(90, mi, mx),
		get_line(120, mi, mx),
		get_line(150, mi, mx),
		get_line(180, mi, mx),
		get_line(210, mi, mx),
		get_line(240, mi, mx), l)

	if err := p.Save(vg.Length(g_width), vg.Length(g_height), name); err != nil {
		panic(err)
	}
}

func dt_to_map(dt plotter.XYs) *seds {
	m := &seds{}
	m.init_map()
	for _, p := range dt {
		s := &sed{}
		s.ds = p.Y
		m.vec[int(p.X)] = s
	}
	return m
}
func resvr(dt plotter.XYs) plotter.XYs {
	mi := 999999.9
	mx := 0.0
	res := plotter.XYs{}
	for _, p := range dt {
		if p.Y > mx {
			mx = p.Y
		}
		if p.Y < mi {
			mi = p.Y
		}
	}
	mp := (mi + mx) / 2

	for _, p := range dt {
		if p.Y > mp {
			p.Y = p.Y - 2*(p.Y-mp)
		} else if p.Y < mp {
			p.Y = p.Y + 2*(mp-p.Y)
		}
		res = append(res, p)
	}
	return res
}

// -----------------------计算相似度
func sim_jb(d1, d2 plotter.XYs) float64 {
	return sim_jb_ss(d1, d2, 5) + sim_jb_ss(d1, d2, 10) + sim_jb_ss(d1, d2, 15) + sim_jb_ss(d1, d2, 30)
}
func sim_jb_ss(d1, d2 plotter.XYs, ss int) float64 {
	if len(d1) != len(d2) {
		return 0.0
	}
	d1_d := []float64{}
	d2_d := []float64{}
	res := 0.0

	for n1, d1_1 := range d1 {
		if n1 > 0 && (n1+1)%ss == 0 {
			r := 0.0
			mm := n1 - ss
			if mm < 0 {
				mm = 0
			}
			if d1_1.Y != d1[mm].Y && d1[mm].Y > 0 {
				a1 := d1_1.Y - d1[mm].Y
				r = a1 / d1[mm].Y
			}
			d1_d = append(d1_d, r)
		}
	}
	for n2, d2_1 := range d2 {
		if n2 > 0 && (n2+1)%ss == 0 {
			r := 0.0
			mm := n2 - ss
			if mm < 0 {
				mm = 0
			}
			if d2_1.Y != d2[mm].Y && d2[mm].Y > 0 {
				a2 := d2_1.Y - d2[mm].Y
				//fmt.Println("@@", a2)
				r = a2 / d2[mm].Y
			}
			d2_d = append(d2_d, r)
		}
	}

	for a, _ := range d1_d {
		//fmt.Println(d1_d[a], "==", d2_d[a])
		if d1_d[a] > 0 && d2_d[a] > 0 {
			res -= float64(ss / 5.0)
		} else if d1_d[a] < 0 && d2_d[a] < 0 {
			res -= float64(ss / 5.0)
		} else {
			res += float64(ss / 5.0)
		}
	}
	return res
}

func data_deal_jb_2(d1, d2 plotter.XYs, jb int) (plotter.XYs, plotter.XYs) {
	jj := 240 - jb
	//最小值

	mi1 := 0.0
	//最大值
	mx1 := 0.0

	for n, p := range d1 {
		if n >= jj {
			if n == jj {
				mi1 = p.Y
				mx1 = p.Y
			}
			if p.Y < mi1 {
				mi1 = p.Y
			}
			if p.Y > mx1 {
				mx1 = p.Y
			}
		}

	}

	mi2 := 0.0
	//最大值
	mx2 := 0.0
	for n, p := range d2 {
		if n >= jj {
			if n == jj {
				mi2 = p.Y
				mx2 = p.Y
			}
			if p.Y < mi2 {
				mi2 = p.Y
			}
			if p.Y > mx2 {
				mx2 = p.Y
			}
		}
	}
	jl_1 := mx1 - mi1
	jl_2 := mx2 - mi2

	d3 := plotter.XYs{}
	d4 := plotter.XYs{}
	for n, p := range d2 {
		if n >= jj {
			a := (p.Y - mi2) * 100.0 / jl_2
			d3 = append(d3, plotter.XY{p.X, Decimal(a)})
		}
	}
	for n1, p1 := range d1 {
		if n1 >= jj {
			a := (p1.Y - mi1) * 100.0 / jl_1
			d4 = append(d4, plotter.XY{p1.X, Decimal(a)})
		}
	}
	return d4, d3
}

// func data_base(d2 plotter.XYs) plotter.XYs {
// 	//最小值
// 	mi1 := 0.0
// 	//最大值
// 	mx1 := 120.0

// 	mi2 := 0.0
// 	//最大值
// 	mx2 := 0.0
// 	for n, p := range d2 {
// 		if n == 0 {
// 			mi2 = p.Y
// 			mx2 = p.Y
// 		}
// 		if p.Y < mi2 {
// 			mi2 = p.Y
// 		}
// 		if p.Y > mx2 {
// 			mx2 = p.Y
// 		}
// 	}
// 	jl_1 := mx1 - mi1
// 	jl_2 := mx2 - mi2

// 	d3 := plotter.XYs{}
// 	for _, p := range d2 {
// 		a := (p.Y - mi2) * jl_1 / jl_2
// 		d3 = append(d3, plotter.XY{p.X, Decimal(a + mi1)})
// 	}
// 	return d3
// }

var excelBase = time.Date(1899, 12, 30, 0, 0, 0, 0, time.Local)

func parseExcelTime(cell structure.CellData) string {
	// timeLayouts := []string{
	// 	"15:04",
	// 	"15:04:05",
	// }
	// for _, layout := range timeLayouts {
	// 	if t, err := time.ParseInLocation(layout, cell.GetString(), time.Local); err == nil {
	// 		return t.Format("15:04")
	// 	}
	// }
	f := cell.GetFloat64()
	seconds := math.Round(f * 86400)
	h := int(seconds) / 3600
	m := (int(seconds) % 3600) / 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

func parseExcelDateTime(cell structure.CellData) (time.Time, error) {
	s := cell.GetString()
	s = strings.Trim(s, " ")
	layouts := []string{
		"2006/01/02-15:04",
		"2006/01/02",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}
	f := cell.GetFloat64()
	seconds := math.Round(f * 86400)
	if f > 0 {
		d := excelBase.Add(time.Duration(seconds * float64(time.Second)))
		return d, nil
	}
	return time.Time{}, fmt.Errorf("无法解析日期/时间: %q", s)
}
func parseExcelDate(cell structure.CellData) (string, error) {
	s := cell.GetString()
	// s = strings.Trim(s, " ")
	// layouts := []string{
	// 	"2006/1/2",
	// 	"2006-1-2",
	// 	"1/2/2006",
	// 	"1-2-2006",
	// 	"2006/01/02",
	// 	"2006-01-02",
	// }
	// for _, layout := range layouts {
	// 	if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
	// 		return t.Format("2006-01-02"), nil
	// 	}
	// }
	f := cell.GetFloat64()
	seconds := math.Round(f * 86400)
	if f > 0 {
		d := excelBase.Add(time.Duration(seconds * float64(time.Second)))
		return d.Format("2006-01-02"), nil
	}
	return "", fmt.Errorf("无法解析日期/时间: %q", s)
}
