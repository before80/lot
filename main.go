package main

import "github.com/before80/lot/cmd"

func main() {
	cmd.Execute()
	//allDlts, err := dbop.ReadAllDlt()
	//if err != nil {
	//	lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("从数据表中读取数据出现错误：%v\n", err), 3)
	//	return
	//}
	//fmt.Printf("读取到 %d 条数据\n", len(allDlts))
	//drawNum2Type2NotExistCount := ana.MinComb(allDlts, 0, "07001")
	//fmt.Printf("drawNum2Type2NotExistCount=%v\n", drawNum2Type2NotExistCount)

	////var counts []*ana.NumCount
	//var counts1 []*ana.EqNumCount
	//var counts2 []*ana.EqAnyCount
	//counts = ana.FMaxOne(allDlts)
	//for _, c := range counts {
	//	fmt.Printf("count: %v\n\n\n", c)
	//}

	//counts = ana.BMaxOne(allDlts)
	//for _, c := range counts {
	//	fmt.Printf("count: %v\n\n\n", c)
	//}

	//counts = ana.BMaxTwo(allDlts)
	//for _, c := range counts {
	//	fmt.Printf("count: %v\n\n\n", c)
	//}

	//counts = ana.FMaxTwo(allDlts)
	//for k, c := range counts {
	//	if k > 10 {
	//		break
	//	}
	//	fmt.Printf("count: %v\n\n\n", c)
	//}

	//counts = ana.FMaxThree(allDlts)
	//for k, c := range counts {
	//	if k > 10 {
	//		break
	//	}
	//	fmt.Printf("count: %v\n\n\n", c)
	//}

	//counts = ana.FMaxFour(allDlts)
	//for k, c := range counts {
	//	if k > 20 {
	//		break
	//	}
	//	fmt.Printf("count: %v\n\n\n", c)
	//}
	//var drawNums []string
	//var ys [][][]string
	//var legends1 []string
	//var legends2 []string
	//for i := 1; i <= 35; i++ {
	//	iStr := strconv.Itoa(i)
	//	if len(iStr) == 1 {
	//		iStr = "0" + iStr
	//	}
	//	legends1 = append(legends1, iStr)
	//}
	//for i := 1; i <= 12; i++ {
	//	iStr := strconv.Itoa(i)
	//	if len(iStr) == 1 {
	//		iStr = "0" + iStr
	//	}
	//	legends2 = append(legends2, iStr)
	//}
	//for _, legend := range legends1 {
	//	counts1, drawNums = ana.EqNFBIntervalDrawNumSpNum(allDlts, 1, 50, 1, legend, true)
	//	var y []string
	//	for _, c := range counts1 {
	//		y = append(y, strconv.Itoa(c.Count))
	//	}
	//	if len(ys) == 0 {
	//		ys = append(ys, [][]string{})
	//	}
	//	ys[0] = append(ys[0], y)
	//}
	//
	//for _, legend := range legends1 {
	//	counts1, drawNums = ana.EqNFBIntervalDrawNumSpNum(allDlts, 2, 50, 1, legend, true)
	//	var y []string
	//	for _, c := range counts1 {
	//		y = append(y, strconv.Itoa(c.Count))
	//	}
	//	if len(ys) == 1 {
	//		ys = append(ys, [][]string{})
	//	}
	//	ys[1] = append(ys[1], y)
	//}
	//
	//for _, legend := range legends1 {
	//	counts1, drawNums = ana.EqNFBIntervalDrawNumSpNum(allDlts, 3, 50, 1, legend, true)
	//	var y []string
	//	for _, c := range counts1 {
	//		y = append(y, strconv.Itoa(c.Count))
	//	}
	//	if len(ys) == 2 {
	//		ys = append(ys, [][]string{})
	//	}
	//	ys[2] = append(ys[2], y)
	//}
	//
	//for _, legend := range legends1 {
	//	counts2, drawNums = ana.EqAnyFBIntervalDrawNumSpNum(allDlts, 50, 1, legend, true)
	//	var y []string
	//	for _, c := range counts2 {
	//		y = append(y, strconv.Itoa(c.Count))
	//	}
	//	if len(ys) == 3 {
	//		ys = append(ys, [][]string{})
	//	}
	//	ys[3] = append(ys[3], y)
	//}
	//
	//for _, legend := range legends2 {
	//	counts1, drawNums = ana.EqNFBIntervalDrawNumSpNum(allDlts, 1, 50, 1, legend, true)
	//	var y []string
	//	for _, c := range counts1 {
	//		y = append(y, strconv.Itoa(c.Count))
	//	}
	//	if len(ys) == 4 {
	//		ys = append(ys, [][]string{})
	//	}
	//	ys[4] = append(ys[4], y)
	//}
	//
	//for _, legend := range legends2 {
	//	counts1, drawNums = ana.EqNFBIntervalDrawNumSpNum(allDlts, 2, 50, 1, legend, true)
	//	var y []string
	//	for _, c := range counts1 {
	//		y = append(y, strconv.Itoa(c.Count))
	//	}
	//	if len(ys) == 5 {
	//		ys = append(ys, [][]string{})
	//	}
	//	ys[5] = append(ys[5], y)
	//}
	//
	//for _, legend := range legends2 {
	//	counts1, drawNums = ana.EqNFBIntervalDrawNumSpNum(allDlts, 3, 50, 1, legend, true)
	//	var y []string
	//	for _, c := range counts1 {
	//		y = append(y, strconv.Itoa(c.Count))
	//	}
	//	if len(ys) == 6 {
	//		ys = append(ys, [][]string{})
	//	}
	//	ys[6] = append(ys[6], y)
	//}
	//
	//for _, legend := range legends2 {
	//	counts2, drawNums = ana.EqAnyFBIntervalDrawNumSpNum(allDlts, 50, 1, legend, true)
	//	var y []string
	//	for _, c := range counts2 {
	//		y = append(y, strconv.Itoa(c.Count))
	//	}
	//	if len(ys) == 7 {
	//		ys = append(ys, [][]string{})
	//	}
	//	ys[7] = append(ys[7], y)
	//}
	////lg.InfoToFileAndStdOut(fmt.Sprintf("ys1: %v\n\n\n", ys1))
	//
	////fmt.Printf("comb: %v\n\n\n", ana.GenComb([]string{"01", "02", "03", "04", "05"}, 3))
	//
	//templatesDir := "./templates"
	//fms := make(template.FuncMap)
	//fms["join"] = strings.Join
	//
	//// 解析目录中所有匹配的模板文件
	//tmpl, err := template.New("t").Funcs(fms).ParseGlob(templatesDir + "/*.html")
	//if err != nil {
	//	log.Fatalf("解析模板错误: %v", err)
	//}
	//_ = os.Mkdir("output", os.ModePerm)
	//file, err := os.Create("output/duidie_output.html")
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer file.Close()
	//
	//type Data struct {
	//	EqNum    string
	//	FB       string
	//	Type     string
	//	Xuhao    int
	//	Interval int
	//	Legends  []string
	//	Ys       [][]string
	//	X        []string
	//}
	//
	//_ = tmpl.ExecuteTemplate(file, "t1.html", []Data{
	//	{EqNum: "1", FB: "前区", Type: "递增", Xuhao: 1, Interval: 50, Legends: legends1, Ys: ys[0], X: drawNums},
	//	{EqNum: "2", FB: "前区", Type: "递增", Xuhao: 2, Interval: 50, Legends: legends1, Ys: ys[1], X: drawNums},
	//	{EqNum: "3", FB: "前区", Type: "递增", Xuhao: 3, Interval: 50, Legends: legends1, Ys: ys[2], X: drawNums},
	//	{EqNum: "任意", FB: "前区", Type: "递增", Xuhao: 4, Interval: 50, Legends: legends1, Ys: ys[3], X: drawNums},
	//	{EqNum: "1", FB: "后区", Type: "递增", Xuhao: 5, Interval: 50, Legends: legends2, Ys: ys[4], X: drawNums},
	//	{EqNum: "2", FB: "后区", Type: "递增", Xuhao: 6, Interval: 50, Legends: legends2, Ys: ys[5], X: drawNums},
	//	{EqNum: "3", FB: "后区", Type: "递增", Xuhao: 7, Interval: 50, Legends: legends2, Ys: ys[6], X: drawNums},
	//	{EqNum: "任意", FB: "后区", Type: "递增", Xuhao: 8, Interval: 50, Legends: legends2, Ys: ys[7], X: drawNums},
	//})

	//fmt.Printf("combCount: %d\n\n\n", ana.CalCombCount())
	//counts, existCombs, drawNum2CombNums := ana.MaxComb(allDlts)
	//_ = drawNum2CombNums
	//_ = counts
	//_ = existCombs
	//fmt.Printf("existCombs: %d\n", len(existCombs))
	//for i, d := range drawNum2CombNums {
	//	if i < len(drawNum2CombNums)-10 {
	//		continue
	//	}
	//	fmt.Printf("drawNum: %s, combCount: %d, combs: %v\n", d.DrawNum, d.CombCount, d.Combs)
	//}

	//for i, count := range counts {
	//	if i > 10 {
	//		break
	//	}
	//
	//	// 前区两个号
	//
	//	fmt.Printf("Num: %s Count: %d\n", count.Num, count.Count)
	//}
	//
	//allCombs := ana.CalAllComb()
	//var notExistCombs []string
	//for _, comb := range allCombs {
	//	if !slices.Contains(existCombs, comb) {
	//		notExistCombs = append(notExistCombs, comb)
	//	}
	//}
	//fmt.Printf("notExistCombs: %d\n", len(notExistCombs))
	//var notExistsF2B0Combs []string
	//var notExistsF1B1Combs []string
	//var notExistsF2B1Combs []string
	//var notExistsF2B2Combs []string
	//var notExistsF3B0Combs []string
	//var notExistsF3B1Combs []string
	//var notExistsF3B2Combs []string
	//for _, comb := range notExistCombs {
	//	//suffixHasSx := strings.HasSuffix(comb, "|")
	//	//containSx := strings.Contains(comb, "|")
	//	commaCount := strings.Count(comb, ",")
	//	sxIndex := strings.Index(comb, "|")
	//	lenComb := len(comb)
	//	if sxIndex == 5 && commaCount == 1 && lenComb == 5 {
	//		notExistsF2B0Combs = append(notExistsF2B0Combs, comb)
	//	} else if sxIndex == 2 && commaCount == 0 && lenComb == 5 {
	//		notExistsF1B1Combs = append(notExistsF1B1Combs, comb)
	//	} else if sxIndex == 5 && commaCount == 1 && lenComb == 8 {
	//		notExistsF2B1Combs = append(notExistsF2B1Combs, comb)
	//	} else if sxIndex == 5 && commaCount == 2 && len(comb) == 11 {
	//		notExistsF2B2Combs = append(notExistsF2B2Combs, comb)
	//	} else if sxIndex == 8 && commaCount == 2 && len(comb) == 9 {
	//		notExistsF3B0Combs = append(notExistsF3B0Combs, comb)
	//	} else if sxIndex == 8 && commaCount == 2 && len(comb) == 11 {
	//		notExistsF3B1Combs = append(notExistsF3B1Combs, comb)
	//	} else if sxIndex == 8 && commaCount == 3 && len(comb) == 14 {
	//		notExistsF3B2Combs = append(notExistsF3B2Combs, comb)
	//	}
	//}
	//
	//fmt.Printf("notExistsF2B0Combs count: %d\n", len(notExistsF2B0Combs))
	//fmt.Printf("notExistsF2B0Combs %#v\n", notExistsF2B0Combs)
	////for _, comb := range notExistsF2B0Combs {
	////	fmt.Printf("notExistsF2B0Combs: %s\n", comb)
	////}
	//fmt.Printf("notExistsF1B1Combs count: %d\n", len(notExistsF1B1Combs))
	//fmt.Printf("notExistsF1B1Combs %#v\n", notExistsF1B1Combs)
	//
	//fmt.Printf("notExistsF2B1Combs count: %d\n", len(notExistsF2B1Combs))
	//fmt.Printf("notExistsF2B1Combs %#v\n", notExistsF2B1Combs)
	//
	//fmt.Printf("notExistsF2B2Combs count: %d\n", len(notExistsF2B2Combs))
	////fmt.Printf("notExistsF2B2Combs %#v\n", notExistsF2B2Combs)
	//
	//fmt.Printf("notExistsF3B0Combs count: %d\n", len(notExistsF3B0Combs))
	////fmt.Printf("notExistsF3B0Combs %#v\n", notExistsF3B0Combs)
	//
	//fmt.Printf("notExistsF3B1Combs count: %d\n", len(notExistsF3B1Combs))
	////fmt.Printf("notExistsF3B1Combs %#v\n", notExistsF3B1Combs)
	//
	//fmt.Printf("notExistsF3B2Combs count: %d\n", len(notExistsF3B2Combs))
	////fmt.Printf("notExistsF3B2Combs %#v\n", notExistsF3B2Combs)

	//for _, comb := range notExistsF1B1Combs {
	//	fmt.Printf("notExistsF1B1Combs: %s\n", comb)
	//}
}
