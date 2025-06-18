package server

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

func StartServer() {
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

// 处理/rmin路径的请求
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
