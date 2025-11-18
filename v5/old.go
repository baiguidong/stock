package main

// import (
// 	"bufio"
// 	"bytes"
// 	"fmt"
// 	"image"
// 	"image/color"
// 	"image/draw"
// 	"image/jpeg"
// 	"io/ioutil"
// 	"math"
// 	"os"

// 	"github.com/disintegration/imaging"
// 	"github.com/golang/freetype"
// 	"gonum.org/v1/plot"
// 	"gonum.org/v1/plot/plotter"
// 	"gonum.org/v1/plot/text"
// 	"gonum.org/v1/plot/vg"
// )

// func test_hz_1(he int, bhe int, fname string) {
// 	text_py := fmt.Sprintf("    %d 平移", g_offset)
// 	text := fmt.Sprintf("    %d 合格", he)
// 	text1 := fmt.Sprintf("    %d 不合格", bhe)
// 	fontBytes, err := ioutil.ReadFile("FZBSJW.ttf")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	f, err := freetype.ParseFont(fontBytes)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fg, bg := image.Black, image.White
// 	rgba := image.NewRGBA(image.Rect(0, 0, 800, 1500))
// 	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
// 	c := freetype.NewContext()
// 	c.SetFont(f)
// 	c.SetFontSize(75)
// 	c.SetClip(rgba.Bounds())
// 	c.SetDst(rgba)
// 	c.SetSrc(fg)
// 	c.SetDPI(120)

// 	pt0 := freetype.Pt(0, 700)
// 	pt := freetype.Pt(0, 1000)
// 	pt1 := freetype.Pt(0, 1300)
// 	_, err = c.DrawString(text_py, pt0)
// 	_, err = c.DrawString(text, pt)
// 	_, err = c.DrawString(text1, pt1)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// Save that RGBA image to disk.
// 	outFile, err := os.Create(fname)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer outFile.Close()
// 	b := bufio.NewWriter(outFile)
// 	err = jpeg.Encode(b, rgba, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	err = b.Flush()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	s1, err := imaging.Open(fname)
// 	if err == nil {
// 		s1 = imaging.Resize(s1, width, high*3+h_head, imaging.Lanczos)
// 	}
// 	imaging.Save(s1, fname)
// }

// func test_hz_empty(fname string) {
// 	fg, bg := image.Black, image.White
// 	rgba := image.NewRGBA(image.Rect(0, 0, 10, 1500))
// 	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
// 	c := freetype.NewContext()
// 	c.SetClip(rgba.Bounds())
// 	c.SetDst(rgba)
// 	c.SetSrc(fg)
// 	c.SetDPI(120)

// 	// Save that RGBA image to disk.
// 	outFile, err := os.Create(fname)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer outFile.Close()
// 	b := bufio.NewWriter(outFile)
// 	err = jpeg.Encode(b, rgba, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	err = b.Flush()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	s1, err := imaging.Open(fname)
// 	if err == nil {
// 		s1 = imaging.Resize(s1, 10, high*3+h_head, imaging.Lanczos)
// 	}
// 	imaging.Save(s1, fname)
// }

// func merg_image_4(srcs []string, des string) {
// 	lens := len(srcs)
// 	ww := (width+2)*lens + a_head
// 	hh := high*4 + h_head + s_head

// 	dst := imaging.New(ww, hh, color.NRGBA{255, 255, 255, 255})
// 	for n, src := range srcs {
// 		s1, err := imaging.Open(src)
// 		if err == nil {
// 			dst = imaging.Paste(dst, s1, image.Pt(a_head/2+n*(width+2), 0))
// 		}
// 	}
// 	// Save the resulting image as JPEG.
// 	err := imaging.Save(dst, des)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// }

// // ----------------------------
// func osjl_jb(d1, d2 plotter.XYs) float64 {
// 	if len(d1) != len(d2) {
// 		return 0.0
// 	}
// 	num := 0.0
// 	for n, d := range d1 {
// 		num += (d.Y - d2[n].Y) * (d.Y - d2[n].Y)
// 	}
// 	xt := 0
// 	for n, _ := range d1 {
// 		if n > 0 && n%5 == 0 {
// 			z1 := d1[n].Y - d1[n-5].Y
// 			z2 := d2[n].Y - d2[n-5].Y
// 			if z1 > 0 && z2 > 0 {
// 				xt++
// 			} else if z1 < 0 && z2 < 0 {
// 				xt++
// 			}
// 		}
// 	}
// 	res := math.Sqrt(num)
// 	a := (float64)(xt * 100 / (len(d1) - 1))

// 	//fmt.Println(">>", xt, "  ", a)
// 	res = res - res*a/100
// 	q := 0.0
// 	for n, _ := range d1 {
// 		if n > 0 && n%5 == 0 && n < (len(d1)-1) {
// 			if (d1[n].Y > d1[n+1].Y) && (d1[n].Y > d1[n-1].Y) {
// 				b := false
// 				if (d2[n].Y > d2[n+1].Y) && (d2[n].Y > d2[n-1].Y) {
// 					b = true
// 				}
// 				if n > 1 && n < (len(d1)-2) {
// 					if (d2[n-1].Y > d2[n].Y) && (d2[n-1].Y > d2[n-2].Y) {
// 						b = true
// 					}
// 					if (d2[n+1].Y > d2[n+2].Y) && (d2[n+1].Y > d2[n].Y) {
// 						b = true
// 					}
// 				}
// 				if n > 2 && n < (len(d1)-3) {
// 					if (d2[n-2].Y > d2[n-1].Y) && (d2[n-2].Y > d2[n-3].Y) {
// 						b = true
// 					}
// 					if (d2[n+2].Y > d2[n+3].Y) && (d2[n+2].Y > d2[n+1].Y) {
// 						b = true
// 					}
// 				}

// 				if b {
// 					q += 1.0
// 				}
// 			}

// 			if (d1[n].Y < d1[n+1].Y) && (d1[n].Y < d1[n-1].Y) {
// 				b := false
// 				if (d2[n].Y < d2[n+1].Y) && (d2[n].Y < d2[n-1].Y) {
// 					b = true
// 				}
// 				if n > 1 && n < (len(d1)-2) {
// 					if (d2[n-1].Y < d2[n].Y) && (d2[n-1].Y < d2[n-2].Y) {
// 						b = true
// 					}
// 					if (d2[n+1].Y < d2[n+2].Y) && (d2[n+1].Y < d2[n].Y) {
// 						b = true
// 					}
// 				}
// 				if n > 2 && n < (len(d1)-3) {
// 					if (d2[n-2].Y < d2[n-1].Y) && (d2[n-2].Y < d2[n-3].Y) {
// 						b = true
// 					}
// 					if (d2[n+2].Y < d2[n+3].Y) && (d2[n+2].Y < d2[n+1].Y) {
// 						b = true
// 					}
// 				}

// 				if b {
// 					q += 1.0
// 				}
// 			}
// 		}
// 	}
// 	ql := q * 100 / (float64)(len(d1))
// 	//fmt.Println("***", ql)
// 	res = res - res*ql/100
// 	return res
// }
// func osjl(d1, d2 plotter.XYs, jb int) float64 {
// 	if len(d1) != len(d2) {
// 		return 0.0
// 	}
// 	num := 0.0
// 	jj := 240 - jb
// 	for n, d := range d1 {
// 		if n >= jj {
// 			num += (d.Y - d2[n].Y) * (d.Y - d2[n].Y)
// 		}
// 	}
// 	return math.Sqrt(num)
// }
// func data_deal_jb(d1, d2 plotter.XYs, jb int) (plotter.XYs, plotter.XYs) {
// 	jj := 240 - jb
// 	//最小值
// 	mi1 := 0.0
// 	//最大值
// 	mx1 := 0.0
// 	for n, p := range d1 {
// 		if n >= jj {
// 			if n == jj {
// 				mi1 = p.Y
// 				mx1 = p.Y
// 			}
// 			if p.Y < mi1 {
// 				mi1 = p.Y
// 			}
// 			if p.Y > mx1 {
// 				mx1 = p.Y
// 			}
// 		}

// 	}
// 	mi2 := 0.0
// 	//最大值
// 	mx2 := 0.0
// 	for n, p := range d2 {
// 		if n >= jj {
// 			if n == jj {
// 				mi2 = p.Y
// 				mx2 = p.Y
// 			}
// 			if p.Y < mi2 {
// 				mi2 = p.Y
// 			}
// 			if p.Y > mx2 {
// 				mx2 = p.Y
// 			}
// 		}
// 	}
// 	jl_1 := mx1 - mi1
// 	jl_2 := mx2 - mi2

// 	d3 := plotter.XYs{}
// 	d4 := plotter.XYs{}
// 	for n, p := range d2 {
// 		if n >= jj {
// 			a := (p.Y - mi2) * jl_1 / jl_2
// 			d3 = append(d3, plotter.XY{p.X, Decimal(a + mi1)})
// 		}
// 	}
// 	for n1, p1 := range d1 {
// 		if n1 >= jj {
// 			d4 = append(d4, plotter.XY{p1.X, p1.Y})
// 		}
// 	}
// 	return d4, d3
// }

// func data_deal(d1, d2 plotter.XYs) (plotter.XYs, plotter.XYs) {
// 	//最小值
// 	mi1 := 0.0
// 	//最大值
// 	mx1 := 0.0
// 	for n, p := range d1 {
// 		if n == 0 {
// 			mi1 = p.Y
// 			mx1 = p.Y
// 		}
// 		if p.Y < mi1 {
// 			mi1 = p.Y
// 		}
// 		if p.Y > mx1 {
// 			mx1 = p.Y
// 		}
// 	}
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
// 	return d1, d3
// }

// func randomPoints_m(dy string) plotter.XYs {
// 	res := plotter.XYs{}
// 	//fmt.Println(dy)
// 	a := g_map_m1.m_map[dy]
// 	if a != nil {
// 		n := 1
// 		for n <= 240 {
// 			aa := a.vec[n]
// 			if aa != nil {
// 				res = append(res,
// 					plotter.XY{X: float64(n), Y: aa.ds})
// 			} else {
// 				res = append(res,
// 					plotter.XY{X: float64(n), Y: 0.0})
// 			}
// 			n++
// 		}
// 	}

// 	return res
// }

// func randomPoints_a(a *seds) plotter.XYs {
// 	res := plotter.XYs{}
// 	if a != nil {
// 		n := 1
// 		for n <= 240 {
// 			aa := a.vec[n]
// 			if aa != nil {
// 				res = append(res,
// 					plotter.XY{X: float64(n), Y: aa.ds})
// 			} else {
// 				res = append(res,
// 					plotter.XY{X: float64(n), Y: 0.0})
// 			}
// 			n++
// 		}
// 	}
// 	return res
// }

// func get_mp(x plotter.XYs) float64 {
// 	mi := 999999.9
// 	mx := 0.0
// 	for _, p := range x {
// 		if p.Y > mx {
// 			mx = p.Y
// 		}
// 		if p.Y < mi {
// 			mi = p.Y
// 		}
// 	}
// 	return (mi + mx) / 2
// }
// func get_p_1(dt, dt1 plotter.XYs, name, title string) {
// 	p := plot.New()
// 	p.Title.Text = title
// 	p.Title.TextStyle = text.Style{
// 		Color: color.RGBA{R: 255, G: 255, B: 255, A: 255},
// 	}

// 	p.HideX()

// 	p.HideY()

// 	p.HideAxes()

// 	p.BackgroundColor = color.RGBA{R: 0, A: 255}

// 	l, err := plotter.NewLine(dt)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	l.LineStyle.Width = vg.Points(float64(g_xt))
// 	l.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

// 	l2, err := plotter.NewLine(dt1)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	l2.LineStyle.Width = vg.Points(float64(g_xt))
// 	l2.LineStyle.Color = color.RGBA{R: 0, G: 255, B: 0, A: 0}

// 	//-----------------------------
// 	mx := 0.0
// 	mi := 0.0
// 	for _, p := range dt {
// 		if mi == 0.0 {
// 			mi = p.Y
// 		}
// 		if p.Y > mx {
// 			mx = p.Y
// 		}
// 		if p.Y < mi {
// 			mi = p.Y
// 		}
// 	}
// 	p.Add(get_line(0, mi, mx), get_line(30, mi, mx),
// 		get_line(60, mi, mx),
// 		get_line(90, mi, mx),
// 		get_line(120, mi, mx),
// 		get_line(150, mi, mx),
// 		get_line(180, mi, mx),
// 		get_line(210, mi, mx), get_line(240, mi, mx), l, l2)

// 	if err := p.Save(8*vg.Inch, 4*vg.Inch, name); err != nil {
// 		fmt.Println(err.Error())
// 	}
// }

// func write_data() {
// 	buf := bytes.NewBufferString("")
// 	for k, v := range g_map.m_map {
// 		buf.WriteString(k)
// 		//data_base
// 		ps := plotter.XYs{}
// 		for i := 1; i <= 240; i++ {
// 			//buf.WriteString(",")
// 			//buf.WriteString(fmt.Sprintf("%.2f", v.vec[i].ds))

// 			p := plotter.XY{}
// 			p.X = float64(i)
// 			p.Y = float64(v.vec[i].ds)
// 			ps = append(ps, p)
// 		}
// 		ps_new := data_base(ps)
// 		for _, a := range ps_new {
// 			buf.WriteString(",")
// 			buf.WriteString(fmt.Sprintf("%.2f", a.Y))
// 		}
// 		buf.WriteString("\n")
// 	}
// 	ioutil.WriteFile("b.bcp", buf.Bytes(), 0777)
// }
