package main

import (
	"flag"
	"fmt"

	"os"
	"time"

	"gonum.org/v1/plot/plotter"
)

var h_head = 0
var g_src plotter.XYs

// 时间
var g_day_ *string = flag.String("day", "2020-02-04", "Use -day <2018-09-18>")
var g_day = ""

// 偏移1
var g_offset_start_ *int = flag.Int("py", 0, "Use -py 0")
var g_offset_start = 0

// 偏移 2
var g_offset_end_ *int = flag.Int("py2", 0, "Use -py2 0")
var g_offset_end = 0

// // 默认0 类型
// var g_ty *int = flag.Int("ty", 2, "Use -ty 0")

// 输出图片数量
var g_x *int = flag.Int("x", 99, "Use -x 0")
var imageNums = 0

// 局部
var g_local_ *int = flag.Int("jb", 0, "Use -jb 240")
var g_local = 0

// 高度
var g_gd_1 *int = flag.Int("gd", 200, "Use -gd 300")
var g_height int

// 子图的宽度
var g_ks_1 *int = flag.Int("ks", 400, "Use -ks 300")
var g_ks int

// 子图高度
var g_gs_1 *int = flag.Int("gs", 200, "Use -gs 100")
var g_gs int

// 子图 右侧局部图的线条数量
var g_local_vline_ *int = flag.Int("xx", 10, "Use -xx 20")
var g_local_vline int

// 大照片最多包括 子照片数量
var g_num_1 *int = flag.Int("gm", 40, "Use -gm 80")
var g_num int

// 宽度
var g_kd_1 *int = flag.Int("kd", 600, "Use -kd 600")
var g_width int

// 线条
var g_xt_1 *int = flag.Int("xt", 2, "Use -xt 1")
var g_xt int

var verticalImages = []string{}

var ads_all = ass{}
var g_map G_map
var g_map_m1 G_map

var g_dz map[string]int
var days []sed

// ----------------------------
func ready() {
	flag.Parse()
	os.RemoveAll("tmp")
	os.Mkdir("tmp", 0777)

	g_height = *g_gd_1
	g_width = *g_kd_1
	g_xt = *g_xt_1
	g_num = *g_num_1
	g_ks = *g_ks_1
	g_gs = *g_gs_1
	g_local_vline = *g_local_vline_
	g_day = *g_day_
	g_offset_start = *g_offset_start_
	g_offset_end = *g_offset_end_
	g_local = *g_local_
	imageNums = *g_x

	if g_height > 0 {
		h_head = g_height / 2
	}

	g_dz = make(map[string]int)
	g_map.init_data()
	g_map_m1.init_data()

	init_dz()
	load_data()

	fmt.Printf("历史数据:%d天\n", len(g_map.m_map))
	load_M1()
	fmt.Printf("M1数据:%d天\n", len(g_map_m1.m_map))
	load_day()
	fmt.Printf("日线数据:%d\n", len(days))
}

// 最终结果
var new_data = ass{}

func run_offset(offset int) {
	ads_all = ass{}
	m1 := get_m1(g_day, offset)
	k, g, d, s := get_hes(m1.deal())
	if len(m1.vec) != 240 {
		time.Sleep(1000 * 1000 * 1000 * 10)
		return
	}

	//算高低点
	m1.gd()
	m1.zd()

	g_src = CVPoints(m1)
	get_p(g_src, fmt.Sprintf("tmp/%s_%d.png", g_day, offset), "")

	if g_local > 0 {
		get_p_120(g_src, fmt.Sprintf("tmp/%s_%d_120.png", g_day, offset), "")
	}

	if offset > 0 {
		g_src_1 := get_m1_py(g_day, offset)
		get_p_py(g_src, g_src_1, fmt.Sprintf("tmp/%d_p.png", offset), "")
	}
	// 合格的图片数量
	var hgzs = 0
	// 不合格的图片数量
	var bhgzs = 0

	all_num := 1

	//合格
	for _, ad := range days {
		k1, g1, d1, s1 := get_hes(ad.dk, ad.dg, ad.dd, ad.ds)
		if k1 == k && s1 == s && (d == d1 || g == g1) {
			//fmt.Println(fmt.Sprintf("%s 符合:%d %d %d %d", ad.dy, k1, g1, d1, s1))
			av := g_map.m_map[ad.dy]
			zh := 0 // 正
			fa := 0 // 反
			if av != nil {
				av.zd()
				if len(av.vec) != 240 {
					fmt.Println("数据错误:", ad.dy, k1, g1, d1, s1)
					continue
				}

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

				av.z = zh
				av.f = fa
				dyn := next_tm(ad.dy)
				dt1 := GetPointsByDay(ad.dy)
				dt2 := GetPointsByDay(dyn)

				mm := 15
				for mm > 1 {
					if len(dt2) > 0 {
						break
					}
					dyn = next_tm(dyn)
					dt2 = GetPointsByDay(dyn)
					mm--
				}
				if fa > zh {
					dt1 = resvr(dt1)
					dt2 = resvr(dt2)
				}

				sds_o := sds{}
				sds_o.dt1 = dt1
				sds_o.dt2 = dt2
				sds_o.dy = g_day
				sds_o.dy_1 = ad.dy
				sds_o.dy_2 = dyn
				sds_o.zh = zh
				sds_o.fa = fa
				sds_o.sm = av
				sds_o.py = offset

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
							sds_new := sds_o
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
			}
		}
	}
	if hgzs > bhgzs {
		// sort.Sort(ads_all)
		for n, rs := range ads_all.data {
			if n < imageNums {
				new_data.data = append(new_data.data, rs)
			}
		}
	}
}
func main() {
	ready()
	if g_offset_end <= g_offset_start {
		g_offset_end = g_offset_start
	}

	// g_offset_start -> g_offset_end 所有offset
	for i := g_offset_start; i <= g_offset_end; i++ {
		// 计算逻辑
		run_offset(i)
	}

	fmt.Println("正在生成照片,请稍等...")
	generate_vimages()
	fmt.Println("正在合成成照片,请稍等...")
	generate_result()
	fmt.Println("完成")
	os.RemoveAll("tmp")
}

// 合成大照片 结果
func generate_result() {
	num := 0
	tmp := []string{}
	for _, s := range verticalImages {
		tmp = append(tmp, s)
		if len(tmp) >= g_num+1 {
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

// 合格---合成竖向的照片
// 返回 竖向照片数组
func generate_vimages() {
	for n, rs := range new_data.data {
		get_p(rs.dt1, "tmp/"+rs.dy_1+".png", "")

		ct := fmt.Sprintf("%d组 %d正 %d反", n+1, rs.zh, rs.fa)
		ct2 := fmt.Sprintf("偏移(%d) %s", rs.py, rs.dy_1)
		test_hz(ct, ct2, "tmp/"+rs.dy_1+"_h.png")

		get_p(rs.dt2, "tmp/"+rs.dy_2+".png", "")

		vimages := []string{"tmp/" + rs.dy_1 + "_h.png", "tmp/" + rs.dy_1 + ".png",
			fmt.Sprintf("tmp/%s_%d.png", rs.dy, rs.py), "tmp/" + rs.dy_2 + ".png"}

		fname := fmt.Sprintf("tmp/out_h_%d.png", n)

		// 和前一个点的偏移不相同,增加空白
		if n > 0 && rs.py != new_data.data[n-1].py {
			verticalImages = append(verticalImages, "empty.png")
		}
		f_py := fmt.Sprintf("tmp/%d_p.png", rs.py)
		if Isexist(f_py) {
			vimages = append(vimages, f_py)
		}
		merg_image_h(vimages, fname)
		verticalImages = append(verticalImages, fname)

		if g_local > 0 {
			// 前一天 局部
			fname_before := fmt.Sprintf("tmp/%s_120.png", rs.dy_1)
			get_p_120(rs.dt1, fname_before, "")
			// 当天 局部
			fname_offset := fmt.Sprintf("tmp/%s_%d_120.png", rs.dy, rs.py)
			fname := fmt.Sprintf("tmp/out_h120_%s_%d.png", rs.dy_1, n)
			merg_image_offset([]string{fname_before, fname_offset}, fname)
			verticalImages = append(verticalImages, fname)
		}
	}
}
