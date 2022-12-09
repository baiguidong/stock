package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"gonum.org/v1/plot/plotter"
)

func load_day() {
	d, err := ioutil.ReadFile("data/日线.txt")
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

//转换 int  对 8 取余数
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
func (c ass_1) Clear() {
	c.data = append([]sds_2{})
}
func (c ass) Swap(i, j int) {
	c.data[i], c.data[j] = c.data[j], c.data[i]
}
func (c ass) Less(i, j int) bool {
	return c.data[i].osjl < c.data[j].osjl
}

func get_m1(dy string, py int) *seds {
	m1 := g_map_m1.m_map[dy]
	if m1 == nil {
		fmt.Println("ERR:获取数据失败")
		os.Exit(-1)
	}
	var m2 *seds
	if py > 0 && py < 240 {
		mm := 15
		dyn := dy
		for mm > 1 {
			dyn = after_tm(dyn)
			m2 = g_map_m1.m_map[dyn]
			if m2 != nil {
				break
			}
			mm--
		}
		if m2 == nil {
			fmt.Println("ERR:获取数据失败_1")
			time.Sleep(1000 * 1000 * 1000 * 8)
			os.Exit(-1)
		}
	}
	if py >= 240 {
		fmt.Println("ERR:偏移错误")
		time.Sleep(1000 * 1000 * 1000 * 8)
		os.Exit(-1)
	}

	em := &sed{}
	rs := []*sed{}
	i := 1
	for i < 241 {
		if m2 == nil {
			rs = append(rs, em)
		} else {
			vv := m2.vec[i]
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

	num := 240 - py
	num_1 := 480 - py
	xl := 1
	for num < num_1 {
		m.vec[xl] = rs[num]
		num++
		xl++
	}
	return m
}

func get_m1_py(dy string, py int) plotter.XYs {
	res := plotter.XYs{}
	a := g_map_m1.m_map[dy]

	num := 240 - py
	nn := num
	if a != nil {
		for num < 241 {
			aa := a.vec[num]
			if aa != nil {
				res = append(res,
					pt{X: float64(num - nn), Y: aa.ds})
			} else {
				res = append(res,
					pt{X: float64(num - nn), Y: 0.0})
			}
			num++
		}
	} else {
		fmt.Println(dy)
	}
	return res
}
func randomPoints(dy string) plotter.XYs {
	res := plotter.XYs{}
	a := g_map.m_map[dy]
	if a != nil {
		n := 1
		for n <= 240 {
			aa := a.vec[n]
			if aa != nil {
				res = append(res,
					pt{X: float64(n), Y: aa.ds})
			} else {
				//fmt.Println(">>>>", dy)
				res = append(res,
					pt{X: float64(n), Y: 0.0})
			}
			n++
		}
	} else {
		//fmt.Println(">>", dy)
	}
	return res
}
func randomPoints_m(dy string) plotter.XYs {
	res := plotter.XYs{}
	//fmt.Println(dy)
	a := g_map_m1.m_map[dy]
	if a != nil {
		n := 1
		for n <= 240 {
			aa := a.vec[n]
			if aa != nil {
				res = append(res,
					pt{X: float64(n), Y: aa.ds})
			} else {
				res = append(res,
					pt{X: float64(n), Y: 0.0})
			}
			n++
		}
	}

	return res
}

func pt_m1(s *seds) plotter.XYs {
	res := plotter.XYs{}
	if s != nil {
		n := 1
		for n <= 240 {
			aa := s.vec[n]
			if aa != nil {
				res = append(res,
					pt{X: float64(n), Y: aa.ds})
			} else {
				res = append(res,
					pt{X: float64(n), Y: 0.0})
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
func after_tm(s string) string {
	t, err := time.Parse("2006-01-02", s)
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
	dt1, dt2       plotter.XYs
	dy, dy_1, dy_2 string
	zh, fa         int
	osjl           float64
	sm             *seds
}
type sds_2 struct {
	dt1, dt2, dt3 plotter.XYs
	n1, n2, n3    int
	o1, o2        float64

	num int
}

func open_pic(name string) {
	cmd := fmt.Sprintf("%s rundll32.exe  C:/windows/system/shimgvw.dll", name)
	_, err := Runcmd(cmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}
func Runcmd(cmd string) ([]byte, error) {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", cmd)
	} else {
		c = exec.Command("/bin/sh", "-c", cmd)
	}
	return c.Output()
}
func write_data() {
	buf := bytes.NewBufferString("")
	for k, v := range g_map.m_map {
		buf.WriteString(k)
		for i := 1; i <= 240; i++ {
			buf.WriteString(",")
			buf.WriteString(fmt.Sprintf("%.2f", v.vec[i].ds))
		}
		buf.WriteString("\n")
	}
	ioutil.WriteFile("a.bcp", buf.Bytes(), 0777)
}
func load_data() {
	v, err := List_file("data", "", 9999)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	for _, v1 := range v {
		if v1 == "data/M1.txt" {
			//fmt.Println(v1)
			continue
		}
		if v1 == "data/日线.txt" {
			//fmt.Println(v1)
			continue
		}
		fmt.Println(v1)
		f, err := os.Open(v1)
		if err != nil {
			continue
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
					//fmt.Println("aa")
					//					if vs[0] == "2016/11/11" {
					//						fmt.Println("aaaaaaaaaa")
					//					}
					g_map.add_data(format_tm(vs[0]), vs[1], a2, a3, a4, a5)
					//fmt.Println(format_tm(vs[0]), vs[1], a2, a3, a4, a5)
				} else {
					//					fmt.Println(vs[0], len(vs))
				}
			} else {
				//fmt.Println(err.Error())
				break
			}
		}
	}
}
func List_file(path string, reg string, max int) ([]string, error) {
	filelist := []string{}

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return filelist, err
	}
	sep := "/"
	for _, fi := range dir {
		if fi.IsDir() {
			s, err := List_file(path+sep+fi.Name(), reg, max)
			if err == nil {
				for _, a := range s {
					filelist = append(filelist, a)
					if len(filelist) >= max {
						return filelist, nil
					}
				}
			}
		} else {
			if len(reg) > 0 {
				m, _ := regexp.MatchString(reg, fi.Name())
				if m {
					filelist = append(filelist, path+sep+fi.Name())
				}
			} else {
				if strings.HasSuffix(fi.Name(), ".tmp") || strings.HasSuffix(fi.Name(), ".temp") || strings.HasSuffix(fi.Name(), ".dealling") {
				} else {
					filelist = append(filelist, path+sep+fi.Name())
				}
			}
			if len(filelist) >= max {
				return filelist, nil
			}
		}
	}
	return filelist, nil
}

type sed struct {
	dy, tm         string
	dk, dg, dd, ds float64
	g, d           int
}

type seds struct {
	vec            map[int]*sed
	z, f           int
	sg, sd, xg, xd float64
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
func (m *seds) iszd() bool {
	if m.isz() || m.isd() {
		return false
	}
	return true
}

//判断涨跌
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
func (m *seds) deal() (float64, float64, float64, float64) {
	//fmt.Println(m.vec[1])
	k := m.vec[1].dk
	s := m.vec[240].ds
	g := m.vec[1].dg
	d := m.vec[1].dd
	//fmt.Println(k, g, d, s)
	for _, v := range m.vec {
		if v.dg > g {
			g = v.dg
		}
		if v.dd < d {
			d = v.dd
		}
	}
	//fmt.Println(k, g, d, s)
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
func Isexist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

type pt struct{ X, Y float64 }

func randomPoints_a(a *seds) plotter.XYs {
	res := plotter.XYs{}
	if a != nil {
		n := 1
		for n <= 240 {
			aa := a.vec[n]
			if aa != nil {
				res = append(res,
					pt{X: float64(n), Y: aa.ds})
			} else {
				res = append(res,
					pt{X: float64(n), Y: 0.0})
			}
			n++
		}
	}
	return res
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
func (g *G_map) print() {
	fmt.Println(fmt.Sprintf("size:%d", len(g.m_map)))

	for key, val := range g.m_map {
		fmt.Println(fmt.Sprintf("key %s size:%d", key, len(val.vec)))
		for key1, val1 := range val.vec {
			fmt.Println(fmt.Sprintf("%d %s [%.2f,%.2f,%.2f,%.2f]",
				key1, val1.tm, val1.dk, val1.dg, val1.dd, val1.ds))
		}
	}
}
func init_dz() {
	g_dz = map[string]int{
		"9:31":  1,
		"9:32":  2,
		"9:33":  3,
		"9:34":  4,
		"9:35":  5,
		"9:36":  6,
		"9:37":  7,
		"9:38":  8,
		"9:39":  9,
		"9:40":  10,
		"9:41":  11,
		"9:42":  12,
		"9:43":  13,
		"9:44":  14,
		"9:45":  15,
		"9:46":  16,
		"9:47":  17,
		"9:48":  18,
		"9:49":  19,
		"9:50":  20,
		"9:51":  21,
		"9:52":  22,
		"9:53":  23,
		"9:54":  24,
		"9:55":  25,
		"9:56":  26,
		"9:57":  27,
		"9:58":  28,
		"9:59":  29,
		"09:31": 1,
		"09:32": 2,
		"09:33": 3,
		"09:34": 4,
		"09:35": 5,
		"09:36": 6,
		"09:37": 7,
		"09:38": 8,
		"09:39": 9,
		"09:40": 10,
		"09:41": 11,
		"09:42": 12,
		"09:43": 13,
		"09:44": 14,
		"09:45": 15,
		"09:46": 16,
		"09:47": 17,
		"09:48": 18,
		"09:49": 19,
		"09:50": 20,
		"09:51": 21,
		"09:52": 22,
		"09:53": 23,
		"09:54": 24,
		"09:55": 25,
		"09:56": 26,
		"09:57": 27,
		"09:58": 28,
		"09:59": 29,
		"10:00": 30,
		"10:01": 31,
		"10:02": 32,
		"10:03": 33,
		"10:04": 34,
		"10:05": 35,
		"10:06": 36,
		"10:07": 37,
		"10:08": 38,
		"10:09": 39,
		"10:10": 40,
		"10:11": 41,
		"10:12": 42,
		"10:13": 43,
		"10:14": 44,
		"10:15": 45,
		"10:16": 46,
		"10:17": 47,
		"10:18": 48,
		"10:19": 49,
		"10:20": 50,
		"10:21": 51,
		"10:22": 52,
		"10:23": 53,
		"10:24": 54,
		"10:25": 55,
		"10:26": 56,
		"10:27": 57,
		"10:28": 58,
		"10:29": 59,
		"10:30": 60,
		"10:31": 61,
		"10:32": 62,
		"10:33": 63,
		"10:34": 64,
		"10:35": 65,
		"10:36": 66,
		"10:37": 67,
		"10:38": 68,
		"10:39": 69,
		"10:40": 70,
		"10:41": 71,
		"10:42": 72,
		"10:43": 73,
		"10:44": 74,
		"10:45": 75,
		"10:46": 76,
		"10:47": 77,
		"10:48": 78,
		"10:49": 79,
		"10:50": 80,
		"10:51": 81,
		"10:52": 82,
		"10:53": 83,
		"10:54": 84,
		"10:55": 85,
		"10:56": 86,
		"10:57": 87,
		"10:58": 88,
		"10:59": 89,
		"11:00": 90,
		"11:01": 91,
		"11:02": 92,
		"11:03": 93,
		"11:04": 94,
		"11:05": 95,
		"11:06": 96,
		"11:07": 97,
		"11:08": 98,
		"11:09": 99,
		"11:10": 100,
		"11:11": 101,
		"11:12": 102,
		"11:13": 103,
		"11:14": 104,
		"11:15": 105,
		"11:16": 106,
		"11:17": 107,
		"11:18": 108,
		"11:19": 109,
		"11:20": 110,
		"11:21": 111,
		"11:22": 112,
		"11:23": 113,
		"11:24": 114,
		"11:25": 115,
		"11:26": 116,
		"11:27": 117,
		"11:28": 118,
		"11:29": 119,
		"11:30": 120,
		"11:31": 121,
		"13:00": 120,
		"13:01": 121,
		"13:02": 122,
		"13:03": 123,
		"13:04": 124,
		"13:05": 125,
		"13:06": 126,
		"13:07": 127,
		"13:08": 128,
		"13:09": 129,
		"13:10": 130,
		"13:11": 131,
		"13:12": 132,
		"13:13": 133,
		"13:14": 134,
		"13:15": 135,
		"13:16": 136,
		"13:17": 137,
		"13:18": 138,
		"13:19": 139,
		"13:20": 140,
		"13:21": 141,
		"13:22": 142,
		"13:23": 143,
		"13:24": 144,
		"13:25": 145,
		"13:26": 146,
		"13:27": 147,
		"13:28": 148,
		"13:29": 149,
		"13:30": 150,
		"13:31": 151,
		"13:32": 152,
		"13:33": 153,
		"13:34": 154,
		"13:35": 155,
		"13:36": 156,
		"13:37": 157,
		"13:38": 158,
		"13:39": 159,
		"13:40": 160,
		"13:41": 161,
		"13:42": 162,
		"13:43": 163,
		"13:44": 164,
		"13:45": 165,
		"13:46": 166,
		"13:47": 167,
		"13:48": 168,
		"13:49": 169,
		"13:50": 170,
		"13:51": 171,
		"13:52": 172,
		"13:53": 173,
		"13:54": 174,
		"13:55": 175,
		"13:56": 176,
		"13:57": 177,
		"13:58": 178,
		"13:59": 179,
		"14:00": 180,
		"14:01": 181,
		"14:02": 182,
		"14:03": 183,
		"14:04": 184,
		"14:05": 185,
		"14:06": 186,
		"14:07": 187,
		"14:08": 188,
		"14:09": 189,
		"14:10": 190,
		"14:11": 191,
		"14:12": 192,
		"14:13": 193,
		"14:14": 194,
		"14:15": 195,
		"14:16": 196,
		"14:17": 197,
		"14:18": 198,
		"14:19": 199,
		"14:20": 200,
		"14:21": 201,
		"14:22": 202,
		"14:23": 203,
		"14:24": 204,
		"14:25": 205,
		"14:26": 206,
		"14:27": 207,
		"14:28": 208,
		"14:29": 209,
		"14:30": 210,
		"14:31": 211,
		"14:32": 212,
		"14:33": 213,
		"14:34": 214,
		"14:35": 215,
		"14:36": 216,
		"14:37": 217,
		"14:38": 218,
		"14:39": 219,
		"14:40": 220,
		"14:41": 221,
		"14:42": 222,
		"14:43": 223,
		"14:44": 224,
		"14:45": 225,
		"14:46": 226,
		"14:47": 227,
		"14:48": 228,
		"14:49": 229,
		"14:50": 230,
		"14:51": 231,
		"14:52": 232,
		"14:53": 233,
		"14:54": 234,
		"14:55": 235,
		"14:56": 236,
		"14:57": 237,
		"14:58": 238,
		"14:59": 239,
		"15:00": 240,
		"15:01": 240,
	}
}
func load_M1() {
	d, err := ioutil.ReadFile("data/M1.txt")
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

			}
		}

	}
}
func Decimal(value float64) float64 {
	val := int64(value * 100)
	if val/10 == 0 {
		val += 1
	}
	return float64(val) / 100
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
		s1 = imaging.Resize(s1, width, h_head, imaging.Lanczos)
	}
	imaging.Save(s1, fname)

}
func test_hz_1(he int, bhe int, fname string) {
	text_py := fmt.Sprintf("    %d 平移", *g_py)
	text := fmt.Sprintf("    %d 合格", he)
	text1 := fmt.Sprintf("    %d 不合格", bhe)
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
	rgba := image.NewRGBA(image.Rect(0, 0, 800, 1500))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetFont(f)
	c.SetFontSize(75)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetDPI(120)

	pt0 := freetype.Pt(0, 700)
	pt := freetype.Pt(0, 1000)
	pt1 := freetype.Pt(0, 1300)
	_, err = c.DrawString(text_py, pt0)
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
		s1 = imaging.Resize(s1, width, high*3+h_head, imaging.Lanczos)
	}
	imaging.Save(s1, fname)
}
