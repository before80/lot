package server

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/before80/lot/ana_dlt"
	"github.com/before80/lot/bs"
	"github.com/before80/lot/cfg"
	"github.com/before80/lot/db"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/dlt"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
	"github.com/go-rod/rod/lib/defaults"
	"golang.org/x/crypto/acme/autocert"
)

func StartServer() {
	mux := http.NewServeMux()
	//http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
	//	//io.WriteString(w, "Hello, world!\n")
	//	http.ServeFile(w, req, "./output/duidie_output.html")
	//})
	//
	//http.HandleFunc("/q", func(w http.ResponseWriter, req *http.Request) {
	//	http.ServeFile(w, req, "./requestHtml/q.html")
	//})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "欢迎访问 %s！", r.Host)
	})

	mux.HandleFunc("/qxh", func(w http.ResponseWriter, req *http.Request) {
		// ===== 1. 获取客户端 IP =====
		clientIP := getClientIP(req)

		// ===== 2. 写入 ips 表（不影响主流程）=====
		go func(ip string) {
			if ip == "" {
				return
			}

			err := db.DB.Create(&models.Ip{
				Ip:        ip,
				CreatedAt: time.Now(),
			}).Error

			if err != nil {
				lg.ErrorToFile(fmt.Sprintf("写入IP失败: %v", err))
			}
		}(clientIP)

		templatesDir := "./requestHtml"
		tmpl, err := template.New("qxh").ParseGlob(templatesDir + "/qxh.html")

		dlts, _ := dbop.ReadAllDlt(false)
		t2MoniABCDEs := ana_dlt.GenTx2MoniABCDE()
		txHis, oeHis, qzhHis, frontDhHis, backDhHis, backCombHis, quShi2St := ana_dlt.Stats(dlts, t2MoniABCDEs)

		lastDlt := dlts[len(dlts)-1]

		err = tmpl.ExecuteTemplate(w, "qxh.html", map[string]interface{}{
			"Dlts":         dlts,
			"LastBackComb": fmt.Sprintf("%s,%s", lastDlt.B1, lastDlt.B2),
			"T2MoniABCDEs": t2MoniABCDEs,
			"TxHis":        txHis,
			"OeHis":        oeHis,
			"QzhHis":       qzhHis,
			"FrontDhHis":   frontDhHis,
			"BackDhHis":    backDhHis,
			"BackCombHis":  backCombHis,
			"QuShi2St":     quShi2St,
			"LastHisSlice": gen.LastHisSlice,
			"LastDrawTime": dlts[len(dlts)-1].DrawTime,
		})

		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("执行模板错误: %v", err), 3)
			return
		}
	})

	mux.HandleFunc("/upde", func(w http.ResponseWriter, req *http.Request) {
		//fmt.Printf("run here \n")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if ana_dlt.GetDltKj() {
			lg.InfoToFile(fmt.Sprintf("访问/upde已执行完 GetDltKj"))
			ana_dlt.DltDataToExcel(true)
			lg.InfoToFile(fmt.Sprintf("访问/upde已执行完 ana_dlt.DltDataToExcel(true)"))
		}

		// 写入响应（Write 会自动设置 200 状态码）
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("已经执行"))
		if err != nil {
			lg.ErrorToFile(fmt.Sprintf("写入 HTTP 响应失败: %v", err))
		}
	})

	//http.HandleFunc("/moni", func(w http.ResponseWriter, req *http.Request) {
	//	templatesDir := "./requestHtml"
	//	tmpl, err := template.New("moni").ParseGlob(templatesDir + "/moni.html")
	//
	//	dlts, _ := dbop.ReadAllDlt(false)
	//
	//	err = tmpl.ExecuteTemplate(w, "moni.html", map[string]interface{}{
	//		"Dlts":         dlts,
	//		"T2MoniABCDEs": ana_dlt.GenTx2MoniABCDE(),
	//	})
	//	if err != nil {
	//		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("执行模板错误: %v", err), 3)
	//		return
	//	}
	//})

	//http.HandleFunc("/rxh", ana_dlt.GenerateNumbersHandler)

	//http.HandleFunc("/qmin", func(w http.ResponseWriter, req *http.Request) {
	//	http.ServeFile(w, req, "./requestHtml/qmin.html")
	//})
	//
	//http.HandleFunc("/r", handleRRequest)
	//http.HandleFunc("/rmin", handleRMinRequest)
	//http.HandleFunc("/his", handleHisRequest)
	//http.HandleFunc("/pdf", handlePdfRequest)   // 只用于临时下载pdf文件
	//http.HandleFunc("/g", handleGenRequest)     //
	//http.HandleFunc("/gl", handleGlRequest)     //
	//http.HandleFunc("/ch", handleChRequest)     //
	//http.HandleFunc("/bqs", handleBQsRequest)   //
	//http.HandleFunc("/zh", handleZhRequest)     //
	//http.HandleFunc("/vqs", handleVQsRequest)   //
	//http.HandleFunc("/aes", handleAesRequest)   //
	//http.HandleFunc("/naes", handleNAesRequest) //

	//_ = os.MkdirAll("download", os.ModePerm)
	//// 固定使用当前目录下的download目录
	//rootDir, err := filepath.Abs("./download")
	//if err != nil {
	//	lg.ErrorToFile(fmt.Sprintf("获取目录绝对路径失败: %v", err))
	//	return
	//}
	//
	//// 检查目录是否存在
	//if _, err = os.Stat(rootDir); os.IsNotExist(err) {
	//	lg.ErrorToFile(fmt.Sprintf("目录不存在: %s", rootDir))
	//	return
	//}
	//
	//// 确保目录可读
	//if _, err = os.Open(rootDir); err != nil {
	//	lg.ErrorToFile(fmt.Sprintf("无法访问目录: %s, 错误: %v", rootDir, err))
	//	return
	//}
	//
	//// 创建文件服务器
	//fileServer := http.FileServer(http.Dir(rootDir))

	// 自定义处理函数
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	// 记录请求日志
	//	lg.InfoToFileAndStdOut(fmt.Sprintf("请求: %s %s 来自 %s", r.Method, r.URL.Path, r.RemoteAddr))
	//
	//	// 安全检查：防止目录遍历攻击
	//	if strings.Contains(r.URL.Path, "..") {
	//		http.Error(w, "非法请求路径", http.StatusBadRequest)
	//		return
	//	}
	//
	//	// 确保只提供PDF文件
	//	if !strings.HasSuffix(strings.ToLower(r.URL.Path), ".pdf") {
	//		http.NotFound(w, r)
	//		return
	//	}
	//
	//	if strings.HasPrefix(r.URL.Path, "/download/") {
	//		_, file := filepath.Split(r.URL.Path)
	//		if file != "" { // 如果不是目录请求
	//			r.URL.Path = file
	//		}
	//	}
	//
	//	// 构建完整文件路径
	//	reqPath := filepath.Join(rootDir, filepath.Clean(r.URL.Path))
	//
	//	lg.InfoToFileAndStdOut(fmt.Sprintf("reqPath: %s\n", reqPath))
	//
	//	// 检查文件是否存在
	//	if _, err = os.Stat(reqPath); os.IsNotExist(err) {
	//		http.NotFound(w, r)
	//		return
	//	}
	//
	//	// 设置响应头，强制下载
	//	filename := filepath.Base(r.URL.Path)
	//	//w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	//	w.Header().Set("Content-Type", "application/pdf")
	//	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	//	// 启用跨域资源共享(CORS)，允许iframe嵌入
	//	w.Header().Set("Access-Control-Allow-Origin", "*")
	//	w.Header().Set("X-Content-Type-Options", "nosniff")
	//	//w.Header().Set("X-Frame-Options", "SAMEORIGIN") // 只能在同源的 iframe 中显示，但 X-Frame-Options: SAMEORIGIN 有一个例外：当被嵌入的是文件资源（如 PDF）时，某些浏览器会严格限制。
	//	//w.Header().Set("X-Frame-Options", "ALLOW-FROM *") // 允许所有域名的iframe嵌入（测试时使用）
	//	//w.Header().Set("Content-Security-Policy", "frame-ancestors 'self' https://lot.cn:8082;")
	//	w.Header().Set("Content-Security-Policy", "frame-ancestors *;") // CSP 策略优先级高于 X-Frame-Options
	//
	//	// 让文件服务器处理请求
	//	fileServer.ServeHTTP(w, r)
	//})

	//http.Handle("/download/", http.StripPrefix("/download/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	// 记录请求
	//	lg.InfoToFileAndStdOut(fmt.Sprintf("PDF请求: %s %s 来自 %s",
	//		r.Method, r.URL.Path, r.RemoteAddr))
	//
	//	// 安全检查
	//	if strings.Contains(r.URL.Path, "..") {
	//		http.Error(w, "非法请求路径", http.StatusBadRequest)
	//		return
	//	}
	//
	//	// 验证PDF扩展名
	//	if !strings.HasSuffix(strings.ToLower(r.URL.Path), ".pdf") {
	//		http.NotFound(w, r)
	//		return
	//	}
	//
	//	// 构建文件路径
	//	filePath := filepath.Join("download", filepath.Clean(r.URL.Path))
	//
	//	// 检查文件是否存在
	//	if _, err := os.Stat(filePath); os.IsNotExist(err) {
	//		http.NotFound(w, r)
	//		return
	//	}
	//
	//	// 设置关键响应头
	//	w.Header().Set("Content-Type", "application/pdf")
	//	w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(r.URL.Path))
	//	w.Header().Set("Access-Control-Allow-Origin", "*")
	//	w.Header().Set("X-Content-Type-Options", "nosniff")
	//
	//	// 关键修复: 使用 CSP 替代 X-Frame-Options
	//	// 允许从同源的 iframe 嵌入
	//	w.Header().Set("Content-Security-Policy", "frame-ancestors 'self'")
	//
	//	// 添加缓存控制头，避免浏览器重复验证证书
	//	w.Header().Set("Cache-Control", "public, max-age=3600")
	//
	//	// 直接提供文件
	//	http.ServeFile(w, r, filePath)
	//}),
	//))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// 处理Excel文件下载
	mux.HandleFunc("/excel/", excelDownloadHandler)

	if cfg.Default.EnvIsLocal == 1 {
		log.Fatal(http.ListenAndServeTLS(
			":8082",
			"lot.cn.crt",
			"lot.cn.key",
			mux,
		))
	} else {
		//startProductionServer()
		//cert, err := tls.LoadX509KeyPair(
		//	"be-better.top_bundle.crt", // 证书 + 中间链
		//	"be-better.top.key",        // 私钥
		//)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//
		//server := &http.Server{
		//	Addr: ":443",
		//	TLSConfig: &tls.Config{
		//		Certificates: []tls.Certificate{cert},
		//		MinVersion:   tls.VersionTLS12,
		//		NextProtos:   []string{"h2", "http/1.1"}, // 浏览器友好
		//	},
		//	Handler: nil,
		//}
		//
		//log.Fatal(server.ListenAndServeTLS("", ""))

		certManager := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("./autoCertCache"), // 证书缓存目录
			HostPolicy: autocert.HostWhitelist("be-better.top", "www.be-better.top"),
		}

		// 2026/01/25 17:39:05 http: TLS handshake error from 117.30.92.184:7991: acme/autocert: host "www.be-better.top" not configured in HostWhitelist

		// ===== HTTP 80：证书校验 + 跳转 HTTPS =====
		go func() {
			httpSrv := &http.Server{
				Addr:    ":80",
				Handler: certManager.HTTPHandler(nil),
			}
			log.Fatal(httpSrv.ListenAndServe())
		}()

		// ===== HTTPS 443：主服务 =====
		httpsSrv := &http.Server{
			Addr:    ":443",
			Handler: mux,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
				MinVersion:     tls.VersionTLS12,
				NextProtos:     []string{"h2", "http/1.1"},
			},
		}

		log.Fatal(httpsSrv.ListenAndServeTLS("", ""))
	}
}

func startLocalServer() {
	lg.InfoToFile("本地环境：使用固定证书文件启动HTTPS...")
	log.Fatal(http.ListenAndServeTLS(
		":8082",
		"lot.cn.crt",
		"lot.cn.key",
		nil,
	))
}

// excelDownloadHandler Excel文件下载处理器
func excelDownloadHandler(w http.ResponseWriter, r *http.Request) {
	// ===== 1. 获取客户端 IP =====
	clientIP := getClientIP(r)

	// ===== 2. 写入 ips 表（不影响主流程）=====
	go func(ip string) {
		if ip == "" {
			return
		}

		err := db.DB.Create(&models.Ip{
			Ip:        ip,
			CreatedAt: time.Now(),
		}).Error

		if err != nil {
			lg.ErrorToFile(fmt.Sprintf("写入IP失败: %v", err))
		}
	}(clientIP)

	// 获取请求的文件名
	filename := filepath.Base(r.URL.Path)

	// 安全验证：确保是Excel文件
	ext := filepath.Ext(filename)

	onlyRiQi := filename[:len(filename)-len(ext)]
	lg.InfoToFileAndStdOut(fmt.Sprintf("要下载的文件的日期是%q\n", onlyRiQi))

	if ext != ".xlsx" {
		http.Error(w, "无效的文件类型", http.StatusBadRequest)
		return
	}

	localFilename := fmt.Sprintf("大乐透截止至%s的数据分析.xlsx", onlyRiQi)

	// 构建文件路径（根据你的实际存储位置调整）
	filePath := filepath.Join("./output/excel", localFilename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 尝试使用今天的日期生成文件名
		//today := time.Now().Format("20060102")
		//filename = today + ".xlsx"
		//filePath = filepath.Join("./output/excel", filename)

		// 如果今天的文件也不存在，返回404
		if _, err = os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "文件不存在", http.StatusNotFound)
			return
		}
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "无法打开文件", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "无法获取文件信息", http.StatusInternalServerError)
		return
	}

	// 设置HTTP头
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", localFilename))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// 将文件内容写入响应
	http.ServeContent(w, r, filename, fileInfo.ModTime(), file)
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
		var counts1 []*ana_dlt.EqAnyCount

		for _, legend := range numbers {
			counts1, drawNums = ana_dlt.EqAnyFBIntervalDrawNumSpNum(dlts, interval, typeVal, legend, true)
			var y []string
			for _, c := range counts1 {
				y = append(y, strconv.Itoa(c.Count))
			}
			ys = append(ys, y)
		}
	} else {
		var counts2 []*ana_dlt.EqNumCount
		for _, legend := range numbers {
			counts2, drawNums = ana_dlt.EqNFBIntervalDrawNumSpNum(dlts, eqNum, interval, typeVal, legend, true)
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
	tmpl, err := template.New("qt").Funcs(fms).ParseGlob(templatesDir + "/qt*.html")
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
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("从数据表中读取数据出现错误：%v\n", err), 1)
		return
	}
	drawNum2Type2NotExistCount := ana_dlt.MinComb(dlts, eqNum, accumulateStartDrawNum)
	//fmt.Printf("drawNum2Type2NotExistCount=%v\n\n\n", drawNum2Type2NotExistCount)

	var drawNums []string
	for drawNum, _ := range drawNum2Type2NotExistCount {
		drawNums = append(drawNums, drawNum)
	}
	slices.Sort(drawNums)
	var ys [][]string
	for _, comb := range legends {
		//fmt.Printf("comb: %s\n\n\n", comb)
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
	fms["split"] = strings.Split
	templatesDir := "./templates/http"

	// 解析目录中所有匹配的模板文件
	tmpl, err := template.New("qmint").Funcs(fms).ParseGlob(templatesDir + "/qmint*.html")
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("解析模板错误: %v", err), 1)
		return
	}
	lg.InfoToFileAndStdOut(fmt.Sprintf("len(remainCombs): %d\n", len(remainCombs)))
	title := "组合递减折线图"
	err = tmpl.ExecuteTemplate(w, "qmint.html", Data{
		Title: title, EqNum: eqNum, Xuhao: 1, Interval: 10, Legends: legends, Ys: ys, X: drawNums, RemainCombs: remainCombs,
	})
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("执行模板错误: %v", err), 1)
		return
	}
}

func handleHisRequest(w http.ResponseWriter, r *http.Request) {
	fms := make(template.FuncMap)
	fms["join"] = strings.Join
	fms["split"] = strings.Split
	fms["oe"] = func(f1, f2, f3, f4, f5, b1, b2 string) string {
		oddNum := 0
		evenNum := 0
		if1, _ := strconv.Atoi(f1)
		if if1%2 == 0 {
			evenNum++
		} else {
			oddNum++
		}
		if2, _ := strconv.Atoi(f2)
		if if2%2 == 0 {
			evenNum++
		} else {
			oddNum++
		}
		if3, _ := strconv.Atoi(f3)
		if if3%2 == 0 {
			evenNum++
		} else {
			oddNum++
		}
		if4, _ := strconv.Atoi(f4)
		if if4%2 == 0 {
			evenNum++
		} else {
			oddNum++
		}
		if5, _ := strconv.Atoi(f5)
		if if5%2 == 0 {
			evenNum++
		} else {
			oddNum++
		}
		ib1, _ := strconv.Atoi(b1)
		if ib1%2 == 0 {
			evenNum++
		} else {
			oddNum++
		}
		ib2, _ := strconv.Atoi(b2)
		if ib2%2 == 0 {
			evenNum++
		} else {
			oddNum++
		}
		return fmt.Sprintf("%d奇%d偶", oddNum, evenNum)
	}
	fms["newSum"] = func(f1, f2, f3, f4, f5, b1, b2 string) int {
		total := 0
		if1, _ := strconv.Atoi(f1)
		total += if1
		if2, _ := strconv.Atoi(f2)
		total += if2
		if3, _ := strconv.Atoi(f3)
		total += if3
		if4, _ := strconv.Atoi(f4)
		total += if4
		if5, _ := strconv.Atoi(f5)
		total += if5
		ib1, _ := strconv.Atoi(b1)
		total += ib1
		ib2, _ := strconv.Atoi(b2)
		total += ib2
		return total
	}
	templatesDir := "./templates/http"

	// 解析目录中所有匹配的模板文件
	tmpl, err := template.New("his").Funcs(fms).ParseGlob(templatesDir + "/dltHis.html")
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("解析模板错误: %v", err))
		return
	}

	//allDlts, err := dbop.ReadDltGETDrawNum("25051")
	allDlts, err := dbop.ReadAllDlt(true)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("从数据表中读取数据出现错误：%v\n", err))
		return
	}

	for i, d := range allDlts {
		//if i < 10 {
		//	lg.InfoToFileAndStdOut(fmt.Sprintf("d.DrawNum: %v\n", d.DrawNum))
		//}
		if strings.Contains(d.DrawPdfUrl, ".pdf") {
			allDlts[i].DrawPdfUrl = filepath.Base(d.DrawPdfUrl)
		}
	}

	err = tmpl.ExecuteTemplate(w, "dltHis.html", allDlts)
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("执行模板错误: %v", err), 3)
		return
	}
}

func handlePdfRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	//dlts, err := dbop.ReadDltGETDrawNum("19081")
	// 从 URL 查询参数中获取 "dn" 的值
	dn := r.URL.Query().Get("dn")
	if dn == "" || len(dn) != 5 {
		http.Error(w, "缺少dn参数，dn即最新的一期期号", http.StatusBadRequest)
		return
	}
	dlts, err := dbop.ReadDltGETDrawNum(dn)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("从数据表中读取数据出现错误：%v\n", err))
	}
	var urls []string
	for _, dlt := range dlts {
		urls = append(urls, dlt.DrawPdfUrl)
	}

	if cfg.Default.CloseBrowser == 2 {
		//defaults.Show = true
		defaults.ResetWith("show=true")
	}

	browser, err1 := bs.GetBrowser(strconv.Itoa(0))
	if err1 != nil {
		lg.ErrorToFile(fmt.Sprintf("第%d次打开浏览器发生错误：%v\n", 0, err1))
		return
	}
	defer browser.Close()
	//fmt.Printf("urls: %v\n", urls)
	page := browser.MustPage("about:blank")
	for _, url := range urls {
		time.Sleep(6 * time.Second)
		page, err = dlt.DownloadPDF(browser, page, url)
		if err != nil {
			lg.ErrorToFile(fmt.Sprintf("下载%s 所在的PDF文件失败: %v\n", url, err))
		}
	}
}

func handleGenRequest(w http.ResponseWriter, r *http.Request) {
	dlts, err := dbop.ReadDltGETDrawNum("25120")
	lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", dlts[0]))
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("从数据表中读取数据出现错误：%v\n", err))
	}

	//GenAllComb(dlts[0])
	//fmt.Printf("%v\n", GenAllComb(dlts[0]))
	lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", ana_dlt.GenAllComb(dlts[0])))
}

func handleGlRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("hhh")
	ana_dlt.FrontGuiLi("25121")
}
func handleChRequest(w http.ResponseWriter, r *http.Request) {

}

func handleBQsRequest(w http.ResponseWriter, r *http.Request) {
	res := ana_dlt.DltBackQuShi1()
	fmt.Println(res)
}
func handleZhRequest(w http.ResponseWriter, r *http.Request) {
	ana_dlt.Zhs()
}

func handleVQsRequest(w http.ResponseWriter, r *http.Request) {
	ana_dlt.ValidateDltBackQuShi()
}
func handleAesRequest(w http.ResponseWriter, r *http.Request) {
	ana_dlt.AbcdeQs()
}

func handleNAesRequest(w http.ResponseWriter, r *http.Request) {
	ana_dlt.NewAbcdeQs("77777")
}

// 辅助函数：打印目录内容（调试用）
func printDirContents(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("  [目录] %s/\n", file.Name())
		} else if strings.HasSuffix(strings.ToLower(file.Name()), ".pdf") {
			fmt.Printf("  [PDF]  %s\n", file.Name())
		} else {
			fmt.Printf("  [其他] %s\n", file.Name())
		}
	}
	return nil
}

func getClientIP(r *http.Request) string {
	// 1. 反向代理 / CDN 常用
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// 可能是多个 IP，取第一个
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}

	// 2. Nginx / Apache
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return strings.TrimSpace(ip)
	}

	// 3. 直连
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}

	return r.RemoteAddr
}
