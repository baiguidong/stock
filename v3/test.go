package main

import (
	"flag"
	"fmt"

	"os"
	"sort"
	"time"

	"gonum.org/v1/plot/plotter"
)

//-----------------------------------
var high = 300
var width = 600
var h_head = high / 2
var s_head = width / 10
var a_head = width / 10

var g_src plotter.XYs

var g_day *string = flag.String("day", "2019-05-28", "Use -day <2018-09-18>")
var g_py *int = flag.Int("py", 0, "Use -py 0")

//默认0
var g_ty *int = flag.Int("ty", 0, "Use -ty 0")
var g_py2 *int = flag.Int("py2", 2, "Use -py2 0")
var g_x *int = flag.Int("x", 3, "Use -x 0")
var g_jb *int = flag.Int("jb", 60, "Use -jb 240")
var g_gd_1 *int = flag.Int("gd", 200, "Use -gd 300")

//子图的高宽
var g_ks_1 *int = flag.Int("ks", 400, "Use -ks 300")
var g_gs_1 *int = flag.Int("gs", 200, "Use -gs 100")
var g_xx_1 *int = flag.Int("xx", 20, "Use -xx 20")

//fazhi
var g_gz_1 *int = flag.Int("gz", 350, "Use -gz 300")
var g_num_1 *int = flag.Int("gm", 40, "Use -gm 80")
var g_kd_1 *int = flag.Int("kd", 600, "Use -kd 600")
var g_xt_1 *int = flag.Int("xt", 2, "Use -xt 1")

var g_gd int
var g_kd int
var g_ks int
var g_gs int
var g_xt int
var g_gz int
var g_num int
var g_xx int
var t_num = 30

var a1 = []string{}

var ads_all = ass{}
var g_map G_map
var g_map_m1 G_map

var g_dz map[string]int
var days []sed

var dy = ""
var py = 0
var py2 = 0
var jb = 0

//----------------------------
func ready() {
	flag.Parse()
	os.RemoveAll("tmp")
	os.Mkdir("tmp", 0777)

	g_gd = *g_gd_1
	g_kd = *g_kd_1
	g_xt = *g_xt_1
	g_gz = *g_gz_1
	g_num = *g_num_1
	g_ks = *g_ks_1
	g_gs = *g_gs_1
	g_xx = *g_xx_1

	if g_gd > 0 {
		high = g_gd
		h_head = high / 2
	}
	if g_kd > 0 {
		width = g_kd
		s_head = width / 2
		a_head = width / 2
	}
	dy = *g_day
	py = *g_py
	py2 = *g_py2
	jb = *g_jb

	g_dz = make(map[string]int)
	g_map.init_data()
	g_map_m1.init_data()

	init_dz()
	load_data()

	fmt.Println(fmt.Sprintf("历史数据:%d天", len(g_map.m_map)))
	load_M1()
	fmt.Println(fmt.Sprintf("M1数据:%d天", len(g_map_m1.m_map)))
	load_day()
	fmt.Println(fmt.Sprintf("日线数据:%d", len(days)))
}

//最终结果
var new_data = ass{}

func run_py(pyn int) {
	ads_all = ass{}
	m1 := get_m1(dy, pyn)
	k, g, d, s := get_hes(m1.deal())
	if len(m1.vec) != 240 {
		time.Sleep(1000 * 1000 * 1000 * 10)
		return
	}

	//算高低点
	m1.gd()
	m1.zd()

	g_src = pt_m1(m1)
	names := fmt.Sprintf("tmp/%s_%d.png", dy, pyn)
	get_p(g_src, names, "")

	if jb >= 30 {
		names_120 := fmt.Sprintf("tmp/%s_%d_120.png", dy, pyn)
		get_p_120(g_src, names_120, "")
	}
	if *g_ty == 2 {
		if pyn > 0 {
			g_src_1 := get_m1_py(dy, pyn)
			name_p := fmt.Sprintf("tmp/%d_p.png", pyn)
			get_p_py(g_src, g_src_1, name_p, "")
		}
	}

	var hgzs = 0
	var bhgzs = 0

	all_num := 1

	//合格
	for _, ad := range days {
		k1, g1, d1, s1 := get_hes(ad.dk, ad.dg, ad.dd, ad.ds)
		if k1 == k && s1 == s {
			used := false

			if d == d1 || g == g1 {
				used = true
			} else {
				//fmt.Println(fmt.Sprintf("%s 不符合:%d %d %d %d", ad.dy, k1, g1, d1, s1))
			}

			if used {
				//fmt.Println(fmt.Sprintf("%s 符合:%d %d %d %d", ad.dy, k1, g1, d1, s1))
				av := g_map.m_map[ad.dy]
				zh := 0
				fa := 0
				if av != nil {
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
						}
					}
					all_num++
				} else {
					//fmt.Println(">>>>>>>", ad.dy)
				}
			}
		}
	}
	if hgzs > bhgzs {
		for n, rs := range ads_all.data {
			d1, d2 := data_deal_jb_2(g_src, rs.dt1, jb)
			rs.osjl = sim_jb(d1, d2)
			ads_all.data[n].osjl = rs.osjl
			ads_all.data[n].py = pyn
		}
		sort.Sort(ads_all)

		for n, rs := range ads_all.data {
			if n < *g_x {
				new_data.data = append(new_data.data, rs)
			}
		}
	}

}
func main() {
	ready()

	if py2 <= py {
		py2 = py
	}
	if py2 >= 0 && py2 >= py {
		for i := py; i <= py2; i++ {
			run_py(i)
		}
	}
	//sort.Sort(new_data)
	fmt.Println("正在生成照片,请稍等...")
	//处理 new_data.data 结果放到 a1  数组
	f2()
	fmt.Println("正在合成成照片,请稍等...")
	//处理a1数据，生成大图
	f5()
	fmt.Println("完成")

	os.RemoveAll("tmp")
	return
}

//合成大照片 结果
func f5() {
	num := 0
	tmp := []string{}
	for _, s := range a1 {
		tmp = append(tmp, s)
		if len(tmp) >= g_num {
			name := fmt.Sprintf(get_home_dir()+"/结果_%d.png", num)
			os.Remove(name)

			merg_image(tmp, name)

			num++
			open_pic(name)
			tmp = []string{}
		}
	}
	if len(tmp) > 0 {
		name := fmt.Sprintf(get_home_dir()+"/结果_%d.png", num)
		os.Remove(name)

		merg_image(tmp, name)

		open_pic(name)
	}
}

//合格---合成竖向的照片
func f2() {
	for n, rs := range new_data.data {
		get_p(rs.dt1, "tmp/"+rs.dy_1+".png", "")

		names_120 := ""
		if jb >= 30 {
			names_120 = fmt.Sprintf("tmp/%s_120.png", rs.dy_1)
			get_p_120(rs.dt1, names_120, "")
		}

		ct := fmt.Sprintf("%d组(%.1f) %d正 %d反", n+1, rs.osjl, rs.zh, rs.fa)
		ct2 := fmt.Sprintf("偏移(%d) %s", rs.py, rs.dy_1)
		test_hz(ct, ct2, "tmp/"+rs.dy_1+"_h.png")

		get_p(rs.dt2, "tmp/"+rs.dy_2+".png", "")

		aaa := []string{"tmp/" + rs.dy_1 + "_h.png", "tmp/" + rs.dy_1 + ".png",
			fmt.Sprintf("tmp/%s_%d.png", rs.dy, rs.py), "tmp/" + rs.dy_2 + ".png"}

		fname := fmt.Sprintf("tmp/out_h_%d.png", n)

		if n > 0 && rs.py != new_data.data[n-1].py {
			a1 = append(a1, "empty.png")
		}

		if *g_ty == 2 {
			f_py := fmt.Sprintf("tmp/%d_p.png", rs.py)
			if Isexist(f_py) {
				aaa = append(aaa, f_py)
			}
		}

		merg_image_h(aaa, fname)
		a1 = append(a1, fname)

		if jb >= 30 {
			fname_120 := fmt.Sprintf("tmp/out_h120_%s_%d.png", rs.dy_1, n)
			aaa_120 := []string{}
			aaa_120 = append(aaa_120, names_120)
			names_120_2 := fmt.Sprintf("tmp/%s_%d_120.png", rs.dy, rs.py)
			aaa_120 = append(aaa_120, names_120_2)

			merg_image_h_120(aaa_120, fname_120)
			a1 = append(a1, fname_120)
		}
	}
}
