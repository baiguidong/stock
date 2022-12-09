package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"sort"
	"time"

	"github.com/disintegration/imaging"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var high = 300
var width = 600
var h_head = high / 2
var s_head = width / 10
var a_head = width / 10

func merg_image(srcs []string, des string) {
	lens := len(srcs)
	ww := (width+2)*lens + a_head
	hh := high*3 + h_head + s_head

	dst := imaging.New(ww, hh, color.NRGBA{255, 255, 255, 255})
	for n, src := range srcs {
		s1, err := imaging.Open(src)
		if err == nil {
			dst = imaging.Paste(dst, s1, image.Pt(a_head/2+n*(width+2), 0))
		}
	}
	// Save the resulting image as JPEG.
	err := imaging.Save(dst, des)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func merg_image_4(srcs []string, des string) {
	lens := len(srcs)
	ww := (width+2)*lens + a_head
	hh := high*4 + h_head + s_head

	dst := imaging.New(ww, hh, color.NRGBA{255, 255, 255, 255})
	for n, src := range srcs {
		s1, err := imaging.Open(src)
		if err == nil {
			dst = imaging.Paste(dst, s1, image.Pt(a_head/2+n*(width+2), 0))
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
	dst := imaging.New(width, high*(lens-1)+4*h_head, color.NRGBA{255, 255, 255, 255})
	for n, src := range srcs {
		s1, err := imaging.Open(src)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if n == 0 {
				s1 = imaging.Resize(s1, width, h_head, imaging.NearestNeighbor)
				dst = imaging.Paste(dst, s1, image.Pt(0, h_head))
			} else {
				s1 = imaging.Resize(s1, width, high, imaging.NearestNeighbor)
				if n == 1 {
					//n*400-200+4*n
					dst = imaging.Paste(dst, s1, image.Pt(0, 2*h_head+4*n))
				} else if n == 2 {
					dst = imaging.Paste(dst, s1, image.Pt(0, 2*h_head+4*n+high))
				} else if n == 3 {
					dst = imaging.Paste(dst, s1, image.Pt(0, 2*h_head+4*n+2*high))
				} else if n == 4 {
					dst = imaging.Paste(dst, s1, image.Pt(0, 2*h_head+4*n+3*high))
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
func get_p(dt plotter.XYs, name, title string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = title
	p.Title.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

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

	if err := p.Save(vg.Length(width), vg.Length(high), name); err != nil {
		panic(err)
	}
}

//基于另一个最大最小
func get_p_py(dt, dt1 plotter.XYs, name, title string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = title
	p.Title.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

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

	if err := p.Save(vg.Length(width), vg.Length(high), name); err != nil {
		panic(err)
	}
}
func get_mp(x plotter.XYs) float64 {
	mi := 999999.9
	mx := 0.0
	for _, p := range x {
		if p.Y > mx {
			mx = p.Y
		}
		if p.Y < mi {
			mi = p.Y
		}
	}
	return (mi + mx) / 2
}
func get_p_1(dt, dt1 plotter.XYs, name, title string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = title
	p.Title.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	p.HideX()

	p.HideY()

	p.HideAxes()

	p.BackgroundColor = color.RGBA{R: 0, A: 255}

	l, err := plotter.NewLine(dt)
	if err != nil {
		fmt.Println(err.Error())
	}
	l.LineStyle.Width = vg.Points(float64(g_xt))
	l.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	l2, err := plotter.NewLine(dt1)
	if err != nil {
		fmt.Println(err.Error())
	}
	l2.LineStyle.Width = vg.Points(float64(g_xt))
	l2.LineStyle.Color = color.RGBA{R: 0, G: 255, B: 0, A: 0}

	//-----------------------------
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
	p.Add(get_line(0, mi, mx), get_line(30, mi, mx),
		get_line(60, mi, mx),
		get_line(90, mi, mx),
		get_line(120, mi, mx),
		get_line(150, mi, mx),
		get_line(180, mi, mx),
		get_line(210, mi, mx), get_line(240, mi, mx), l, l2)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, name); err != nil {
		fmt.Println(err.Error())
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
func osjl_jb(d1, d2 plotter.XYs) float64 {
	if len(d1) != len(d2) {
		return 0.0
	}
	num := 0.0
	for n, d := range d1 {
		num += (d.Y - d2[n].Y) * (d.Y - d2[n].Y)
	}
	xt := 0
	for n, _ := range d1 {
		if n > 0 && n%5 == 0 {
			z1 := d1[n].Y - d1[n-5].Y
			z2 := d2[n].Y - d2[n-5].Y
			if z1 > 0 && z2 > 0 {
				xt++
			} else if z1 < 0 && z2 < 0 {
				xt++
			}
		}
	}
	res := math.Sqrt(num)
	a := (float64)(xt * 100 / (len(d1) - 1))

	//fmt.Println(">>", xt, "  ", a)
	res = res - res*a/100
	q := 0.0
	for n, _ := range d1 {
		if n > 0 && n%5 == 0 && n < (len(d1)-1) {
			if (d1[n].Y > d1[n+1].Y) && (d1[n].Y > d1[n-1].Y) {
				b := false
				if (d2[n].Y > d2[n+1].Y) && (d2[n].Y > d2[n-1].Y) {
					b = true
				}
				if n > 1 && n < (len(d1)-2) {
					if (d2[n-1].Y > d2[n].Y) && (d2[n-1].Y > d2[n-2].Y) {
						b = true
					}
					if (d2[n+1].Y > d2[n+2].Y) && (d2[n+1].Y > d2[n].Y) {
						b = true
					}
				}
				if n > 2 && n < (len(d1)-3) {
					if (d2[n-2].Y > d2[n-1].Y) && (d2[n-2].Y > d2[n-3].Y) {
						b = true
					}
					if (d2[n+2].Y > d2[n+3].Y) && (d2[n+2].Y > d2[n+1].Y) {
						b = true
					}
				}

				if b {
					q += 1.0
				}
			}

			if (d1[n].Y < d1[n+1].Y) && (d1[n].Y < d1[n-1].Y) {
				b := false
				if (d2[n].Y < d2[n+1].Y) && (d2[n].Y < d2[n-1].Y) {
					b = true
				}
				if n > 1 && n < (len(d1)-2) {
					if (d2[n-1].Y < d2[n].Y) && (d2[n-1].Y < d2[n-2].Y) {
						b = true
					}
					if (d2[n+1].Y < d2[n+2].Y) && (d2[n+1].Y < d2[n].Y) {
						b = true
					}
				}
				if n > 2 && n < (len(d1)-3) {
					if (d2[n-2].Y < d2[n-1].Y) && (d2[n-2].Y < d2[n-3].Y) {
						b = true
					}
					if (d2[n+2].Y < d2[n+3].Y) && (d2[n+2].Y < d2[n+1].Y) {
						b = true
					}
				}

				if b {
					q += 1.0
				}
			}
		}
	}
	ql := q * 100 / (float64)(len(d1))
	//fmt.Println("***", ql)
	res = res - res*ql/100
	return res
}
func osjl(d1, d2 plotter.XYs, jb int) float64 {
	if len(d1) != len(d2) {
		return 0.0
	}
	num := 0.0
	jj := 240 - jb
	for n, d := range d1 {
		if n >= jj {
			num += (d.Y - d2[n].Y) * (d.Y - d2[n].Y)
		}
	}
	return math.Sqrt(num)
}
func data_deal_jb(d1, d2 plotter.XYs, jb int) (plotter.XYs, plotter.XYs) {
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
			a := (p.Y - mi2) * jl_1 / jl_2
			d3 = append(d3, pt{p.X, Decimal(a + mi1)})
		}
	}
	for n1, p1 := range d1 {
		if n1 >= jj {
			d4 = append(d4, pt{p1.X, p1.Y})
		}
	}
	return d4, d3
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
			d3 = append(d3, pt{p.X, Decimal(a)})
		}
	}
	for n1, p1 := range d1 {
		if n1 >= jj {
			a := (p1.Y - mi1) * 100.0 / jl_1
			d4 = append(d4, pt{p1.X, Decimal(a)})
		}
	}
	return d4, d3
}
func data_deal(d1, d2 plotter.XYs) (plotter.XYs, plotter.XYs) {
	//最小值
	mi1 := 0.0
	//最大值
	mx1 := 0.0
	for n, p := range d1 {
		if n == 0 {
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
	mi2 := 0.0
	//最大值
	mx2 := 0.0
	for n, p := range d2 {
		if n == 0 {
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
	jl_1 := mx1 - mi1
	jl_2 := mx2 - mi2

	d3 := plotter.XYs{}
	for _, p := range d2 {
		a := (p.Y - mi2) * jl_1 / jl_2
		d3 = append(d3, pt{p.X, Decimal(a + mi1)})
	}
	return d1, d3
}
func data_base(d2 plotter.XYs) plotter.XYs {
	//最小值
	mi1 := 0.0
	//最大值
	mx1 := 120.0

	mi2 := 0.0
	//最大值
	mx2 := 0.0
	for n, p := range d2 {
		if n == 0 {
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
	jl_1 := mx1 - mi1
	jl_2 := mx2 - mi2

	d3 := plotter.XYs{}
	for _, p := range d2 {
		a := (p.Y - mi2) * jl_1 / jl_2
		d3 = append(d3, pt{p.X, Decimal(a + mi1)})
	}
	return d3
}

var g_src plotter.XYs

var g_day *string = flag.String("day", "2018-11-09", "Use -day <2018-09-18>")
var g_py *int = flag.Int("py", 120, "Use -py 0")
var g_jb *int = flag.Int("jb", 240, "Use -jb 240")
var g_ty *int = flag.Int("ty", 5, "Use -ty 1")
var g_gd_1 *int = flag.Int("gd", 200, "Use -gd 300")

//fazhi
var g_gz_1 *int = flag.Int("gz", 350, "Use -gz 300")
var g_num_1 *int = flag.Int("gm", 20, "Use -gm 80")
var g_kd_1 *int = flag.Int("kd", 400, "Use -kd 600")
var g_xt_1 *int = flag.Int("xt", 2, "Use -xt 1")

var g_gd int
var g_kd int
var g_xt int
var g_ty_i int
var g_gz int
var g_num int
var t_num = 30
var hgzs = 0
var bhgzs = 0
var res_vec1 = []string{}
var a1 = []string{}
var a2 = []string{}
var a3 = []string{}
var ads_all = ass{}
var g_map G_map
var g_map_m1 G_map
var g_dz map[string]int
var days []sed

//不合格
var ads_all_bhg = ass{}

//共振
var ads_all_gz = ass_1{}
var res_vec2 = []string{}

var name_p = ""

func main() {
	flag.Parse()
	os.RemoveAll("tmp")
	os.Mkdir("tmp", 0777)

	g_gd = *g_gd_1
	g_kd = *g_kd_1
	g_xt = *g_xt_1
	g_gz = *g_gz_1
	g_ty_i = *g_ty
	g_num = *g_num_1

	if g_gd > 0 {
		high = g_gd
		h_head = high / 2
	}
	if g_kd > 0 {
		width = g_kd
		s_head = width / 2
		a_head = width / 2
	}
	dy := *g_day
	py := *g_py
	jb := *g_jb

	g_dz = make(map[string]int)
	g_map.init_data()
	g_map_m1.init_data()

	init_dz()
	load_data()
	write_data()
	return

	fmt.Println(fmt.Sprintf("历史数据:%d天", len(g_map.m_map)))
	load_M1()
	fmt.Println(fmt.Sprintf("M1数据:%d天", len(g_map_m1.m_map)))
	load_day()
	fmt.Println(fmt.Sprintf("日线数据:%d", len(days)))

	m1 := get_m1(dy, py)
	k, g, d, s := get_hes(m1.deal())
	if len(m1.vec) != 240 {
		time.Sleep(1000 * 1000 * 1000 * 10)
		return
	}

	//算高低点
	m1.gd()
	m1.zd()

	g_src = pt_m1(m1)
	names := "tmp/" + dy + ".png"
	get_p(g_src, names, "")

	if py > 0 {
		g_src_1 := get_m1_py(dy, py)
		name_p = "tmp/" + dy + "_p.png"
		get_p_py(g_src, g_src_1, name_p, "")
	}

	all_num := 1

	//合格

	for _, ad := range days {
		k1, g1, d1, s1 := get_hes(ad.dk, ad.dg, ad.dd, ad.ds)
		if k1 == k && s1 == s {
			used := false

			if *g_ty == 5 {
				used = true
			} else {
				if d == d1 || g == g1 {
					used = true
				}
			}

			if used {
				//fmt.Println(".................", all_num)
				av := g_map.m_map[ad.dy]
				zh := 0
				fa := 0
				if av != nil {
					//					time.Sleep(1000 * 1000 * 1000 * 1)
					//					fmt.Println(ad.dy)
					av.zd()
					if len(av.vec) == 240 {
						g_map.m_map[ad.dy].gd()
						in := 1
						for in < 241 {
							m2 := m1.vec[in]
							av2 := av.vec[in]
							if m2 != nil && av2 != nil {
								if m2.g > 0 {
									if av2.g > 0 {
										zh++
									} else if av2.d > 0 {
										fa++
									}

								} else if m2.d > 0 {
									if av2.d > 0 {
										zh++
									} else if av2.g > 0 {
										fa++
									}
								}
							}
							in++
						}
					}
					av.z = zh
					av.f = fa
					dyn := next_tm(ad.dy)
					dt1 := randomPoints(ad.dy)
					dt2 := randomPoints(dyn)

					mm := 15
					for mm > 1 {
						if len(dt2) > 0 {
							break
						}
						dyn = next_tm(dyn)
						dt2 = randomPoints(dyn)
						mm--
					}
					if fa > zh {
						dt1 = resvr(dt1)
						dt2 = resvr(dt2)
					}

					sds_o := sds{}
					sds_o.dt1 = dt1
					sds_o.dt2 = dt2
					sds_o.dy = dy
					sds_o.dy_1 = ad.dy
					sds_o.dy_2 = dyn
					sds_o.zh = zh
					sds_o.fa = fa
					sds_o.sm = av

					hg := true

					m2 := dt_to_map(dt1)
					m2.zd()
					//中间震荡
					zdzt := false
					//涨
					if m1.isz() {
						if m2.isz() {
							hgzs++
						} else if m2.isd() {
							hg = false
							//fmt.Println(ad.dy, av.sg, av.sd, av.xg, av.xd)
							bhgzs++
						} else {
							hgzs++
						}
					} else if m1.isd() {
						//跌
						if m2.isd() {
							hgzs++
						} else if m2.isz() {
							hg = false
							bhgzs++
						} else {
							hgzs++
						}
					} else {
						zdzt = true
						hgzs++
					}
					if hg {
						if zdzt {
							if fa == zh {
								sds_new := sds{}
								sds_new = sds_o
								sds_new.dt1 = resvr(sds_new.dt1)
								sds_new.dt2 = resvr(sds_new.dt2)

								ads_all.data = append(ads_all.data, sds_o)
								ads_all.data = append(ads_all.data, sds_new)
								hgzs++
							} else {
								ads_all.data = append(ads_all.data, sds_o)
							}
						} else {
							ads_all.data = append(ads_all.data, sds_o)
						}
					} else {
						if fa == zh {
							sds_o.dt1 = resvr(sds_o.dt1)
							sds_o.dt2 = resvr(sds_o.dt2)
							ads_all.data = append(ads_all.data, sds_o)
							hgzs++
							bhgzs--
						} else {
							ads_all_bhg.data = append(ads_all_bhg.data, sds_o)
						}
					}
					all_num++
				} else {
					//fmt.Println(">>>>>>>", ad.dy)
				}
			}
		}
	}
	if *g_ty == 4 {

	} else {
		//合格
		ads_all_tmp := ass{}
		for n, rs := range ads_all.data {
			d1, d2 := data_deal_jb_2(g_src, rs.dt1, jb)
			rs.osjl = osjl_jb(d1, d2)

			ads_all.data[n].osjl = rs.osjl

			if g_ty_i == 5 {
				if rs.osjl < float64(g_gz) {
					ads_all_tmp.data = append(ads_all_tmp.data, ads_all.data[n])
				}
			}
		}
		if g_ty_i == 5 {
			ads_all = ads_all_tmp
		}
		sort.Sort(ads_all)

		//处理合格的照片  算共振

		for n, rs := range ads_all.data {
			sds_1_a := sds_2{}
			sds_1_a.num = 1
			sds_1_a.dt1 = rs.dt2
			sds_1_a.n1 = n
			for i := n + 1; i < len(ads_all.data); i++ {
				rs2 := ads_all.data[i]
				d1, d2 := data_deal_jb_2(rs.dt2, rs2.dt2, 240)
				f1 := osjl_jb(data_base(d1), data_base(d2))
				if f1 < float64(g_gz) {
					if sds_1_a.num == 1 {
						sds_1_a.dt2 = rs2.dt2
						sds_1_a.n2 = i
						sds_1_a.o1 = f1
						sds_1_a.num = 2
					} else if sds_1_a.num == 2 {
						sds_1_a.dt3 = rs2.dt2
						sds_1_a.n3 = i
						sds_1_a.o2 = f1
						sds_1_a.num = 3
						ads_all_gz.data = append(ads_all_gz.data, sds_1_a)
					} else {
						sds_1_a = sds_2{}
						sds_1_a.num = 1
						sds_1_a.dt1 = rs.dt2
						sds_1_a.n1 = n

						sds_1_a.dt2 = rs2.dt2
						sds_1_a.n2 = i
						sds_1_a.o1 = f1
						sds_1_a.num = 2
					}
				}
			}
			if sds_1_a.num == 2 {
				ads_all_gz.data = append(ads_all_gz.data, sds_1_a)
			}
		}

		fmt.Println(len(ads_all_gz.data))
		//不合格
		ads_all_tmp_bhg := ass{}
		for n, rs := range ads_all_bhg.data {
			d1, d2 := data_deal_jb_2(g_src, rs.dt1, jb)
			rs.osjl = osjl_jb(d1, d2)
			ads_all_bhg.data[n].osjl = rs.osjl

			if g_ty_i == 5 {
				if rs.osjl < float64(g_gz) {
					ads_all_tmp_bhg.data = append(ads_all_tmp_bhg.data, ads_all_bhg.data[n])
				}
			}
		}
		if g_ty_i == 5 {
			ads_all_bhg = ads_all_tmp_bhg
		}
		sort.Sort(ads_all_bhg)
	}

	//自动推荐
	ads_all_3 := ass{}
	if *g_ty == 4 {
		for key, val := range g_map.m_map {
			dyn := next_tm(key)
			dt1 := randomPoints_a(val)
			dt2 := randomPoints(dyn)

			mm := 15
			for mm > 1 {
				if len(dt2) > 0 {
					break
				}
				dyn = next_tm(dyn)
				dt2 = randomPoints(dyn)
				mm--
			}

			sds := sds{}
			sds.dt1 = dt1
			sds.dt2 = dt2
			sds.dy = dy
			sds.dy_1 = key
			sds.dy_2 = dyn
			ads_all_3.data = append(ads_all_3.data, sds)
		}
	}
	rnum := 0
	rnum_a := make(chan int, 5)

	fmt.Println("正在生成照片,请稍等...")
	if *g_ty == 4 {
		//智能挖掘
		for n, rs := range ads_all_3.data {
			d1, d2 := data_deal(g_src, rs.dt1)
			rs.osjl = osjl_jb(d1, d2)
			ads_all_3.data[n].osjl = rs.osjl
		}
		sort.Sort(ads_all_3)
		for n, rs := range ads_all_3.data {
			// if n == g_num/2 {
			// 	break
			// }
			if n == g_num {
				break
			}
			data_deal(rs.dt1, g_src)

			get_p(rs.dt1, "tmp/"+rs.dy_1+".png", "")

			ct := fmt.Sprintf("%d组(%.1f) %d正 %d反", n+1, rs.osjl, rs.zh, rs.fa)
			test_hz(ct, "        "+rs.dy_1, "tmp/"+rs.dy_1+"_h.png")

			get_p(rs.dt2, "tmp/"+rs.dy_2+".png", "")

			aaa := []string{"tmp/" + rs.dy_1 + "_h.png", "tmp/" + rs.dy_1 + ".png",
				"tmp/" + rs.dy + ".png", "tmp/" + rs.dy_2 + ".png"}
			fname := fmt.Sprintf("tmp/out_h_%d.png", n)

			if *g_py > 0 {
				aaa = append(aaa, name_p)
			}
			merg_image_h(aaa, fname)
			a3 = append(a3, fname)
		}
	} else {
		go f1(rnum_a)
		rnum++
		go f2(rnum_a)
		rnum++
		go f3(rnum_a)
		rnum++
	}
	//等待图片合并
	for a := 0; a < rnum; a++ {
		<-rnum_a
	}
	fmt.Println("正在合成成照片,请稍等...")

	if *g_ty == 4 {
		for _, s := range a3 {
			res_vec2 = append(res_vec2, s)
		}
	} else {
		if g_ty_i == 5 {
			test_hz_1(hgzs, bhgzs, "tmp/em.png")
			res_vec2 = append(res_vec2, fmt.Sprintf("tmp/em.png"))
		} else {
			if len(ads_all_bhg.data) > 0 {
				test_hz_1(hgzs, bhgzs, "tmp/em.png")
				res_vec2 = append(res_vec2, fmt.Sprintf("tmp/em.png"))
			}
		}

		for _, s := range a1 {
			res_vec2 = append(res_vec2, s)
		}
		if len(a2) > 0 {
			res_vec2 = append(res_vec2, fmt.Sprintf("tmp/em.png"))
		}
		for _, s := range a2 {
			res_vec2 = append(res_vec2, s)
		}
	}
	go f4(rnum_a)
	go f5(rnum_a)

	fmt.Println(hgzs, "", bhgzs)
	<-rnum_a
	<-rnum_a
	fmt.Println("完成")
	os.RemoveAll("tmp")
	return
}

//合成大照片 共振
func f4(i chan int) {
	num := 0
	if len(res_vec1) > 0 {
		tmp := []string{}
		for _, s := range res_vec1 {
			tmp = append(tmp, s)
			if len(tmp) >= g_num {
				name := fmt.Sprintf("共振分析_%d_a_%d.png", g_ty_i, num)
				os.Remove(name)
				merg_image(tmp, name)
				num++
				open_pic(name)
				tmp = []string{}
			}
		}
		if len(tmp) > 0 {
			name := fmt.Sprintf("共振分析_%d_a_%d.png", g_ty_i, num)
			os.Remove(name)
			merg_image(tmp, name)
			open_pic(name)
		}
	}
	i <- 1
}

//合成大照片 结果
func f5(i chan int) {
	num := 0
	if len(res_vec2) > 0 {
		tmp := []string{}
		for _, s := range res_vec2 {
			tmp = append(tmp, s)
			if len(tmp) >= g_num {
				name := fmt.Sprintf("比对结果_%d_a_%d.png", g_ty_i, num)
				os.Remove(name)
				if (g_ty_i == 6 || g_ty_i == 4 || g_ty_i == 5) && *g_py > 0 {
					merg_image_4(tmp, name)
				} else {
					merg_image(tmp, name)
				}
				num++
				open_pic(name)
				tmp = []string{}
			}
		}
		if len(tmp) > 0 {
			name := fmt.Sprintf("比对结果_%d_a_%d.png", g_ty_i, num)
			os.Remove(name)
			if (g_ty_i == 6 || g_ty_i == 4 || g_ty_i == 5) && *g_py > 0 {
				merg_image_4(tmp, name)
			} else {
				merg_image(tmp, name)
			}
			open_pic(name)
		}
	}
	i <- 1
}

//共振
func f1(i chan int) {
	for n, rs := range ads_all_gz.data {
		aaa := []string{}
		fm0 := fmt.Sprintf("tmp/%d_0.png", n)
		aaa = append(aaa, fm0)
		fm := fmt.Sprintf("tmp/%d_1.png", n)
		aaa = append(aaa, fm)
		get_p(rs.dt1, fm, "")
		fm = fmt.Sprintf("tmp/%d_2.png", n)
		get_p(rs.dt2, fm, "")
		aaa = append(aaa, fm)
		ct := ""
		ct2 := ""
		if rs.num == 3 {
			fm = fmt.Sprintf("tmp/%d_3.png", n)
			get_p(rs.dt3, fm, "")
			aaa = append(aaa, fm)
			ct = fmt.Sprintf("%d组 %d组(%.1f)", rs.n1+1, rs.n2+1, rs.o1)
			ct2 = fmt.Sprintf("%d组(%.1f)", rs.n3+1, rs.o2)
		} else {
			ct = fmt.Sprintf("%d组 %d组(%.1f)", rs.n1+1, rs.n2+1, rs.o1)
		}
		test_hz(ct, ct2, fm0)
		fname := fmt.Sprintf("tmp/out_hn_%d.png", n)
		merg_image_h(aaa, fname)
		res_vec1 = append(res_vec1, fname)
	}
	i <- 1
}

//合格
func f2(i chan int) {
	for n, rs := range ads_all.data {
		get_p(rs.dt1, "tmp/"+rs.dy_1+".png", "")

		ct := fmt.Sprintf("%d组(%.1f) %d正 %d反", n+1, rs.osjl, rs.zh, rs.fa)
		test_hz(ct, "        "+rs.dy_1, "tmp/"+rs.dy_1+"_h.png")

		get_p(rs.dt2, "tmp/"+rs.dy_2+".png", "")

		aaa := []string{"tmp/" + rs.dy_1 + "_h.png", "tmp/" + rs.dy_1 + ".png",
			"tmp/" + rs.dy + ".png", "tmp/" + rs.dy_2 + ".png"}

		fname := fmt.Sprintf("tmp/out_h_%d.png", n)
		if (g_ty_i == 6 || g_ty_i == 5) && *g_py > 0 {
			aaa = append(aaa, name_p)
		}
		merg_image_h(aaa, fname)
		a1 = append(a1, fname)
	}
	i <- 1
}

//不合格
func f3(i chan int) {
	for n, rs := range ads_all_bhg.data {
		get_p(rs.dt1, "tmp/"+rs.dy_1+".png", "")

		ct := fmt.Sprintf("%d组(%.1f) %d正 %d反", n+1, rs.osjl, rs.zh, rs.fa)
		test_hz(ct, "        "+rs.dy_1, "tmp/"+rs.dy_1+"_h.png")

		get_p(rs.dt2, "tmp/"+rs.dy_2+".png", "")

		aaa := []string{"tmp/" + rs.dy_1 + "_h.png", "tmp/" + rs.dy_1 + ".png",
			"tmp/" + rs.dy + ".png", "tmp/" + rs.dy_2 + ".png"}
		fname := fmt.Sprintf("tmp/out_h_b_%d.png", n)

		if (g_ty_i == 6 || g_ty_i == 5) && *g_py > 0 {
			aaa = append(aaa, name_p)
		}
		merg_image_h(aaa, fname)

		a2 = append(a2, fname)
	}
	i <- 1
}
