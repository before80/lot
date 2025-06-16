package main

import (
	"fmt"
	"github.com/before80/lot/ana"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/lg"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func main() {

	//var err error
	//startTime := time.Now()
	//lastDlt := dbop.GetLastDlt()
	//db.DB.Last(&lastDlt)
	//lg.InfoToFileAndStdOut(fmt.Sprintf("当前数据库中最新的一条记录为 %v \n", lastDlt))
	//dlts, err := sport.GetSomeDltFromWeb(lastDlt, "25065")
	//if err != nil {
	//	lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("从网页获取开奖数据出现错误：%v\n", err), 3)
	//	return
	//}
	//insertedRow, err := dbop.InsertDltBatch(dlts, 100)
	//if err != nil {
	//	lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("往数据表中插入数据出现错误：%v\n", err), 3)
	//	return
	//} else {
	//	lg.InfoToFileAndStdOut(fmt.Sprintf("插入了 %d 条数据\n", insertedRow))
	//	lg.InfoToFileAndStdOut(fmt.Sprintf("程序运行时间：%.2f秒\n", time.Since(startTime).Seconds()))
	//}

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

	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		//io.WriteString(w, "Hello, world!\n")
		http.ServeFile(w, req, "./output/duidie_output.html")
	})

	http.HandleFunc("/q", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "./requestHtml/q.html")
	})

	http.HandleFunc("/qmin", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "./requestHtml/qmin.html")
	})

	http.HandleFunc("/r", handleRRequest)
	http.HandleFunc("/rmin", handleRMinRequest)

	log.Fatal(http.ListenAndServe(":8080", nil))

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

// 处理/r路径的请求
func handleRRequest(w http.ResponseWriter, r *http.Request) {
	// 设置响应头为JSON格式
	//w.Header().Set("Content-Type", "application/json")

	// 解析URL参数
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "解析请求参数失败", http.StatusBadRequest)
		return
	}

	// 从URL参数中获取表单数据
	startDrawNum := r.FormValue("startDrawNum")
	eqNumStr := r.FormValue("eqNum")
	intervalStr := r.FormValue("interval")
	typeStr := r.FormValue("type")
	numbers := r.Form["numbers"] // 获取所有名为"numbers"的参数值
	lg.InfoToFile(fmt.Sprintf("startDrawNum: %s, eqNum: %s, interval: %s, type: %s, numbers: %v\n", startDrawNum, eqNumStr, intervalStr, typeStr, numbers))

	// 将字符串类型的参数转换为数值类型
	eqNum, err := strconv.Atoi(eqNumStr)
	if err != nil {
		eqNum = 0 // 默认值
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		interval = 20 // 默认值
	}

	typeVal, err := strconv.Atoi(typeStr)
	if err != nil {
		typeVal = 0 // 默认值
	}

	// 验证起始期号格式
	if len(startDrawNum) != 5 {
		http.Error(w, "起始期号必须为5位数字", http.StatusBadRequest)
		return
	}

	// 将数据转换为JSON格式并返回
	//response := fmt.Sprintf(`{
	//	"startDrawNum": "%s",
	//	"eqNum": %d,
	//	"interval": %d,
	//	"type": %d,
	//	"numbers": %q
	//}`, startDrawNum, eqNum, interval, typeVal, numbers)

	dlts, err := dbop.ReadDltGETDrawNum(startDrawNum)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("从数据表中读取数据出现错误：%v\n", err))
		return
	}
	var drawNums []string
	var ys [][]string
	if eqNum == 0 {
		var counts1 []*ana.EqAnyCount

		for _, legend := range numbers {
			counts1, drawNums = ana.EqAnyFBIntervalDrawNumSpNum(dlts, interval, typeVal, legend, true)
			var y []string
			for _, c := range counts1 {
				y = append(y, strconv.Itoa(c.Count))
			}
			ys = append(ys, y)
		}
	} else {
		var counts2 []*ana.EqNumCount
		for _, legend := range numbers {
			counts2, drawNums = ana.EqNFBIntervalDrawNumSpNum(dlts, eqNum, interval, typeVal, legend, true)
			var y []string
			for _, c := range counts2 {
				y = append(y, strconv.Itoa(c.Count))
			}

			ys = append(ys, y)
		}
	}

	type Data struct {
		Title    string
		EqNum    int
		FB       string
		Type     string
		Xuhao    int
		Interval int
		Legends  []string
		Ys       [][]string
		X        []string
	}
	fms := make(template.FuncMap)
	fms["join"] = strings.Join
	templatesDir := "./templates/http"

	// 解析目录中所有匹配的模板文件
	tmpl, err := template.New("qt").Funcs(fms).ParseGlob(templatesDir + "/*.html")
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("解析模板错误: %v", err))
		return
	}

	var title, fb string
	switch typeVal {
	case 1:
		fb = "前区"
	case 2:
		fb = "后区"
	case 3:
		fb = "前后区"

	}

	switch eqNum {
	case 1:
		title = "设备1-" + fb + "-间隔" + strconv.Itoa(interval) + "-" + "递增"
	case 2:
		title = "设备2-" + fb + "-间隔" + strconv.Itoa(interval) + "-" + "递增"
	case 3:
		title = "设备3-" + fb + "-间隔" + strconv.Itoa(interval) + "-" + "递增"
	case 0:
		title = "任意-" + fb + "-间隔" + strconv.Itoa(interval) + "-" + "递增"
	}

	err = tmpl.ExecuteTemplate(w, "qt.html", Data{
		Title: title, EqNum: eqNum, FB: fb, Type: "递增", Xuhao: 1, Interval: interval, Legends: numbers, Ys: ys, X: drawNums,
	})
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("执行模板错误: %v", err), 3)
		return
	}
}

func handleRMinRequest(w http.ResponseWriter, r *http.Request) {
	// 设置响应头为JSON格式
	//w.Header().Set("Content-Type", "application/json")

	// 解析URL参数
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "解析请求参数失败", http.StatusBadRequest)
		return
	}

	// 从URL参数中获取表单数据
	accumulateStartDrawNum := r.FormValue("accumulateStartDrawNum")
	startDrawNum := r.FormValue("startDrawNum")
	eqNumStr := r.FormValue("eqNum")
	legends := r.Form["combinations"] // 获取所有名为"numbers"的参数值

	// 将字符串类型的参数转换为数值类型
	eqNum, err := strconv.Atoi(eqNumStr)
	if err != nil {
		eqNum = 0 // 默认值
	}

	lg.InfoToFile(fmt.Sprintf("startDrawNum: %s, eqNum: %d,  legends: %v\n", startDrawNum, eqNum, legends))

	// 验证起始期号格式
	if len(accumulateStartDrawNum) != 5 {
		http.Error(w, "从某期数开始累计必须为5位数字", http.StatusBadRequest)
		return
	}

	if len(startDrawNum) != 5 {
		http.Error(w, "起始期号必须为5位数字", http.StatusBadRequest)
		return
	}

	dlts, err := dbop.ReadDltGETDrawNum(startDrawNum)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("从数据表中读取数据出现错误：%v\n", err))
		return
	}
	drawNum2Type2NotExistCount := ana.MinComb(dlts, eqNum, accumulateStartDrawNum)
	//fmt.Printf("drawNum2Type2NotExistCount=%v\n\n\n", drawNum2Type2NotExistCount)

	var drawNums []string
	for drawNum, _ := range drawNum2Type2NotExistCount {
		drawNums = append(drawNums, drawNum)
	}
	slices.Sort(drawNums)
	var ys [][]string
	for _, comb := range legends {
		fmt.Printf("comb: %s\n\n\n", comb)
		var y []string
		for _, drawNum := range drawNums {
			y = append(y, strconv.Itoa(drawNum2Type2NotExistCount[drawNum][comb].NoExistCombCount))
		}
		ys = append(ys, y)
	}

	lastDrawNum := drawNums[len(drawNums)-1]
	remainCombs := make(map[string][]string)

	for _, comb := range legends {
		//if !(comb == "all" || comb == "2|2" || comb == "3|2") {
		if !(comb == "all") {
			remainCombs[comb] = append(remainCombs[comb], drawNum2Type2NotExistCount[lastDrawNum][comb].NoExistCombs...)
		}
	}
	type Data struct {
		Title       string
		EqNum       int
		Xuhao       int
		Interval    int
		Legends     []string
		Ys          [][]string
		X           []string
		RemainCombs map[string][]string
	}
	fms := make(template.FuncMap)
	fms["join"] = strings.Join
	templatesDir := "./templates/http"

	// 解析目录中所有匹配的模板文件
	tmpl, err := template.New("qmint").Funcs(fms).ParseGlob(templatesDir + "/*.html")
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("解析模板错误: %v", err))
		return
	}
	title := "组合递减折线图"
	err = tmpl.ExecuteTemplate(w, "qmint.html", Data{
		Title: title, EqNum: eqNum, Xuhao: 1, Interval: 10, Legends: legends, Ys: ys, X: drawNums, RemainCombs: remainCombs,
	})
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("执行模板错误: %v", err), 3)
		return
	}
}
