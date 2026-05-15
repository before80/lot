package dlt

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/before80/lot/ana_dlt"
	"github.com/before80/lot/bs"
	"github.com/before80/lot/cfg"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/spf13/cobra"
)

//go:embed kjCommonFun.js
var kjCommonFunJs string

// UpdateDlt1 从官网获取最新的大乐透开奖数据并批量更新数据表
func UpdateDlt1() {
	startTime := time.Now()
	lastDlt := dbop.GetLastDlt()
	//db.DB.Last(&lastDlt)
	lg.InfoToFileAndStdOut(fmt.Sprintf("当前数据库中最新的一条记录为 %v \n", lastDlt))

	ldn := strconv.Itoa(time.Now().Year())[2:] + "156"

	dlts, err := GetSomeDltFromWeb(lastDlt, ldn)
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("从网页获取开奖数据出现错误：%v\n", err), 3)
		return
	}
	insertedRow, err := dbop.InsertDltBatch(dlts, 100)
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("往数据表中插入数据出现错误：%v\n", err), 3)
		lg.InfoToFileAndStdOut(fmt.Sprintf("程序运行时间：%.2f秒\n", time.Since(startTime).Seconds()))
		return
	} else {
		lg.InfoToFileAndStdOut(fmt.Sprintf("插入了 %d 条数据\n", insertedRow))
		lg.InfoToFileAndStdOut(fmt.Sprintf("程序运行时间：%.2f秒\n", time.Since(startTime).Seconds()))
	}
	lg.InfoToFileAndStdOut(fmt.Sprintf("已经处理完毕! \n等待10秒钟后自动关闭窗口\n"))
	time.Sleep(10 * time.Second)
}

// UpdateDlt 接收命令行参数并从官网获取最新的大乐透开奖数据并批量更新数据表,并且会下载对应开奖数据的PDF文档
func UpdateDlt(cmd *cobra.Command) {
	var err error
	// 从命令行参数中获取ldn的值
	ldn, err := cmd.Flags().GetString("ldn")
	lg.InfoToFileAndStdOut(fmt.Sprintf("ldn=%s\n", ldn))
	if ldn == "" {
		return
	}

	startTime := time.Now()
	lastDlt := dbop.GetLastDlt()
	//db.DB.Last(&lastDlt)
	lg.InfoToFileAndStdOut(fmt.Sprintf("当前数据库中最新的一条记录为 %v \n", lastDlt))
	dlts, err := GetSomeDltFromWeb(lastDlt, ldn)
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("从网页获取开奖数据出现错误：%v\n", err), 3)
		return
	}
	insertedRow, err := dbop.InsertDltBatch(dlts, 100)
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("往数据表中插入数据出现错误：%v\n", err), 3)
		return
	} else {
		lg.InfoToFileAndStdOut(fmt.Sprintf("插入了 %d 条数据\n", insertedRow))
		lg.InfoToFileAndStdOut(fmt.Sprintf("处理pdf中...\n"))
		if cfg.Default.CloseBrowser == 2 {
			//defaults.Show = true
			defaults.ResetWith("show=true")
		}

		browser, err1 := bs.GetBrowser(strconv.Itoa(1))
		if err1 != nil {
			lg.ErrorToFile(fmt.Sprintf("第%d次打开浏览器发生错误：%v\n", 0, err1))
			return
		}
		defer browser.Close()
		page := browser.MustPage("about:blank")
		dealPdfNum := 0
		for _, dlt := range dlts {
			if dlt.DrawNum >= "19081" {
				time.Sleep(6 * time.Second)
				page, err = DownloadPDF(browser, page, dlt.DrawPdfUrl)
				if err != nil {
					lg.ErrorToFile(fmt.Sprintf("下载 %s 所在PDF出现错误：%v\n", dlt.DrawPdfUrl, err))
				}
				dealPdfNum++
			}
		}

		lg.InfoToFileAndStdOut(fmt.Sprintf("成功下载了%d个pdf文件\n", dealPdfNum))
		lg.InfoToFileAndStdOut(fmt.Sprintf("程序运行时间：%.2f秒\n", time.Since(startTime).Seconds()))
	}
}

func DownloadPDF(browser *rod.Browser, page *rod.Page, url string) (*rod.Page, error) {
	lg.InfoToFileAndStdOut(fmt.Sprintf("url=%s\n", url))
	defer func() {
		if r := recover(); r != nil {
			lg.InfoToFileAndStdOut(fmt.Sprintf("err=%v\n", r))
		}
	}()
	_ = os.MkdirAll("download", os.ModePerm)
	// 从URL中提取文件名
	filename := filepath.Base(url)
	//lg.InfoToFileAndStdOut(fmt.Sprintf("3\n"))
	// 如果提取的文件名不包含.pdf后缀，则添加后缀
	if !hasPDFExtension(filename) {
		filename += ".pdf"
	}

	router := browser.HijackRequests()
	router.MustAdd("*.pdf", func(ctx *rod.Hijack) {
		//_ = ctx.LoadResponse(http.DefaultClient, true)
		//ctx.Response.SetBody(js.DltHistory)

		err1 := ctx.LoadResponse(http.DefaultClient, true)
		if err1 != nil {
			lg.InfoToFileAndStdOut(fmt.Sprintf("err1=%v\n", err1))
		}
		// 继续原始请求
		//ctx.MustLoadResponse()
		// 获取原始响应
		//time.Sleep(1 * time.Second)

		originalResponse := ctx.Response.Body()
		lg.InfoToFileAndStdOut(fmt.Sprintf("%s len=%d\n", filename, len(originalResponse)))
		//lg.InfoToFileAndStdOut(fmt.Sprintf("%s len=%d %v\n", filename, len(originalResponse), originalResponse))
		_ = os.WriteFile(filepath.Join("download", filename), []byte(originalResponse), 0644)
		lg.InfoToFileAndStdOut(fmt.Sprintf("已保存%s\n", filename))
	})
	defer router.Stop()
	go router.Run()

	page.MustNavigate(url).MustWaitLoad()
	time.Sleep(6 * time.Second)

	//	pageUrlObj, _ := page.Eval(`() => {
	//	return window.PageUrl;
	//}`)
	//
	//	lg.InfoToFileAndStdOut(fmt.Sprintf("pageUrl=%s\n", pageUrlObj.Value))

	//time.Sleep(1000 * time.Second)

	return page, nil
}

func DownloadPDF1(url string) error {
	// 创建HTTP客户端
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头，模拟浏览器访问
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 从URL中提取文件名
	filename := filepath.Base(url)

	// 如果提取的文件名不包含.pdf后缀，则添加后缀
	if !hasPDFExtension(filename) {
		filename += ".pdf"
	}

	_ = os.MkdirAll("download", os.ModePerm)
	// 创建文件
	out, err := os.Create("download/" + filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	// 下载文件内容并写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	//fmt.Printf("文件已成功下载并保存为: %s\n", filename)
	return nil
}

// 检查文件名是否有.pdf后缀
func hasPDFExtension(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".pdf" || ext == ".PDF"
}

func judgeCrossYearGt(startTerm, endTerm string, gtNum int) bool {
	startYear, _ := strconv.Atoi(startTerm[0:2])
	endYear, _ := strconv.Atoi(endTerm[0:2])
	if endYear-startYear > gtNum {
		return true
	}
	return false
}

func GenCrossTermYears(startTerm, endTerm string, gtNum int) (crossTermYears [][2]string) {
	curEndTerm := startTerm
	curStartTerm := startTerm
	var curEndTermNum int
	if judgeCrossYearGt(startTerm, endTerm, gtNum) {
	LabelForContinue0:
		curStartYear, _ := strconv.Atoi(curEndTerm[0:2])
		curEndTerm = strconv.Itoa(curStartYear+gtNum) + "001"
		curEndTermNum, _ = strconv.Atoi(curEndTerm)
		if curEndTerm <= endTerm {
			crossTermYears = append(crossTermYears, [][2]string{
				{curStartTerm, curEndTerm},
			}...)
			curEndTermNum += 1
			curStartTerm = strconv.Itoa(curEndTermNum)
			// 补0
			if len(curStartTerm) < len(startTerm) {
				curStartTerm = "0" + curStartTerm
			}
		} else {
			crossTermYears = append(crossTermYears, [][2]string{
				{curStartTerm, endTerm},
			}...)
		}

		if curEndTerm <= endTerm {
			if judgeCrossYearGt(curEndTerm, endTerm, gtNum) {
				goto LabelForContinue0
			} else {
				crossTermYears = append(crossTermYears, [][2]string{
					{curStartTerm, endTerm},
				}...)
			}
		}
	} else {
		crossTermYears = [][2]string{
			{startTerm, endTerm},
		}
	}
	return
}

// GetLotteryHistory 从官网中获取开奖数据(已经整理好的数据)
func GetLotteryHistory(intervalSecond int, startTerm, endTerm string) (lotteryHistory []ListItem) {
	crossTermYears := GenCrossTermYears(startTerm, endTerm, 5)
	for _, crossTermYear := range crossTermYears {
		// 组合返回的结果
		iListItems := getLotteryHistory3(crossTermYear[0], crossTermYear[1])
		fmt.Printf("%s~%s len=%d\n", crossTermYear[0], crossTermYear[1], len(iListItems))
		lotteryHistory = append(lotteryHistory, iListItems...)
		time.Sleep(time.Duration(intervalSecond) * time.Second)
	}

	// 排序
	slices.SortFunc(lotteryHistory, func(a, b ListItem) int {
		//aTerm, _ := strconv.Atoi(a.LotteryDrawNum)
		//bTerm, _ := strconv.Atoi(a.LotteryDrawNum)
		if a.LotteryDrawNum > a.LotteryDrawNum {
			return 1
		} else if a.LotteryDrawNum < a.LotteryDrawNum {
			return -1
		}
		return 0
	})

	return
}

func GetLotteryHistory0(gameNo, pageSize, pageNo, intervalSecond int, startTerm, endTerm string) (lotteryHistory []ListItem) {
	crossTermYears := GenCrossTermYears(startTerm, endTerm, 5)
	for _, crossTermYear := range crossTermYears {
		// 组合返回的结果
		curPageNo := pageNo
		lotterySt := getLotteryHistory0(gameNo, pageSize, pageNo, crossTermYear[0], crossTermYear[1])
		lotteryHistory = append(lotteryHistory, lotterySt.Value.List...)
		if curPageNo <= lotterySt.Value.Pages {
		LabelForContinue1:
			curPageNo += 1
			for curPageNo <= lotterySt.Value.Pages {
				time.Sleep(time.Duration(intervalSecond) * time.Second)
				lotterySt = getLotteryHistory0(gameNo, pageSize, curPageNo, crossTermYear[0], crossTermYear[1])
				lotteryHistory = append(lotteryHistory, lotterySt.Value.List...)
				goto LabelForContinue1
			}
		}
	}

	// 排序
	slices.SortFunc(lotteryHistory, func(a, b ListItem) int {
		//aTerm, _ := strconv.Atoi(a.LotteryDrawNum)
		//bTerm, _ := strconv.Atoi(a.LotteryDrawNum)
		if a.LotteryDrawNum > a.LotteryDrawNum {
			return 1
		} else if a.LotteryDrawNum < a.LotteryDrawNum {
			return -1
		}
		return 0
	})

	return
}

func getLotteryHistory0(gameNo, pageSize, pageNo int, startTerm, endTerm string) (lotteryData LotterySt) {
	if cfg.Default.UseGoRodToGetLottery == 1 {
		return getLotteryHistory1(gameNo, pageSize, pageNo, startTerm, endTerm)
	} else {
		return getLotteryHistory2(gameNo, pageSize, pageNo, startTerm, endTerm)
	}
}

func GetLotteryHistory2(startTerm, endTerm string) (lotteryHistory []ListItem) {
	lotteryHistory = getLotteryHistory3(startTerm, endTerm)

	// 排序
	slices.SortFunc(lotteryHistory, func(a, b ListItem) int {
		if a.LotteryDrawNum > b.LotteryDrawNum {
			return 1
		} else if a.LotteryDrawNum < b.LotteryDrawNum {
			return -1
		}
		return 0
	})

	return
}

// getLotteryHistory3 从官网中获取相关开奖数据(未整理的数据)
func getLotteryHistory3(startTerm, endTerm string) (lotteryDataSlice []ListItem) {
	lg.InfoToFileAndStdOut(fmt.Sprintf("%s~%s 初始化浏览器中... \n", startTerm, endTerm))
	if cfg.Default.CloseBrowser == 2 {
		//defaults.Show = true
		defaults.ResetWith("show=true")
	}

	lg.InfoToFileAndStdOut(fmt.Sprintf("%s~%s 初始化浏览器中... \n", startTerm, endTerm))
	sNum, _ := strconv.Atoi(startTerm)
	eNum, _ := strconv.Atoi(endTerm)

	browser, err1 := bs.GetBrowser(strconv.Itoa(0))
	if err1 != nil {
		lg.ErrorToFile(fmt.Sprintf("第%d次打开浏览器发生错误：%v\n", 0, err1))
		return
	}
	lg.InfoToFileAndStdOut(fmt.Sprintf("1 \n"))

	var dataMu sync.Mutex

	chNextPageNo := make(chan int)
	var closeOnce sync.Once
	hadExistDrawNum := make(map[string]struct{})

	appendFunc := func(data []ListItem) {
		for _, v := range data {
			if v.LotteryDrawNum > endTerm {
				continue
			}
			if _, ok := hadExistDrawNum[v.LotteryDrawNum]; !ok {
				hadExistDrawNum[v.LotteryDrawNum] = struct{}{}
				lotteryDataSlice = append(lotteryDataSlice, v)
			}
		}
	}

	safeClose := func() {
		closeOnce.Do(func() {
			close(chNextPageNo)
		})
	}

	// https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=85&provinceId=0&pageSize=30&isVerify=1&pageNo=1&startTerm=07001&endTerm=12002
	router := browser.HijackRequests()
	router.MustAdd("https://static.sporttery.cn/res_1_0/jcw/default/kj/kjCommonFun.js", func(ctx *rod.Hijack) {
		_ = ctx.LoadResponse(http.DefaultClient, true)
		ctx.Response.SetBody(kjCommonFunJs)
	})

	router.MustAdd("https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=85*", func(ctx *rod.Hijack) {
		url := ctx.Request.URL().RawPath
		ctx.MustLoadResponse()
		// 获取原始响应
		originalResponse := ctx.Response.Body()
		lotteryData := LotterySt{}

		err := json.Unmarshal([]byte(originalResponse), &lotteryData)
		if err != nil {
			close(chNextPageNo)
			panic(fmt.Sprintf("解析%s返回结果，遇到错误：%v\n", url, err))
		}

		// ===== 线程安全写入 =====
		dataMu.Lock()
		appendFunc(lotteryData.Value.List)
		dataMu.Unlock()
		// ======================
		lg.InfoToFileAndStdOut(fmt.Sprintf("lotteryData.Value.PageNo = %d lotteryData.Value.Pages = %d\n", lotteryData.Value.PageNo, lotteryData.Value.Pages))
		if lotteryData.Value.PageNo >= lotteryData.Value.Pages {
			lg.InfoToFileAndStdOut(fmt.Sprintf("发现页面PageNo >= Pages\n"))
			// 已是最后一页
			safeClose()
			return
		}

		lg.InfoToFileAndStdOut(fmt.Sprintf("未发现页面PageNo >= Pages\n"))
		nextPageNo := lotteryData.Value.PageNo + 1

		// 非阻塞写，避免 channel 已关闭或满导致死锁
		select {
		case chNextPageNo <- nextPageNo:
		default:
		}

	})
	// 先设置好要监测的网址中的相关http请求,再打开相关网页,这样才能监测到
	go router.Run()
	lg.InfoToFileAndStdOut(fmt.Sprintf("2 \n"))
	url := "https://www.lottery.gov.cn/kj/kjlb.html?dlt"

	// 再打开相关网页
	page := browser.MustPage(url)
	lg.InfoToFileAndStdOut(fmt.Sprintf("3 \n"))
	page.MustWaitLoad()
	page.MustSetViewport(cfg.Default.BrowserWidth, cfg.Default.BrowserHeight, 1, false)

	iframe := page.MustElement("#iFrame1").MustFrame()
	iframe.Eval(`() => {
			document.querySelector(".g-history .u-zdy ul").style.display = 'block';
		}`)
	time.Sleep(1 * time.Second)
	iframe.Eval(fmt.Sprintf(`() => {
			document.getElementById("bterm").value = '%s';
			document.getElementById("eterm").value = '%s';
			document.getElementById("searchBrn").click();
		}`, startTerm, endTerm))
	lg.InfoToFileAndStdOut(fmt.Sprintf("4 \n"))
	time.Sleep(3 * time.Second)
	lg.InfoToFile(fmt.Sprintf("在大乐透历史开奖页面已点击查询"))

	for nextPageNo := range chNextPageNo {
		if nextPageNo > 0 {
			// 目前设置一个分页大小为100, 以此为依据,防止首次打开页面,卡在这里
			if eNum-sNum > 100 {
				lg.InfoToFileAndStdOut(fmt.Sprintf("nextPageNo=%d > 0\n", nextPageNo))
			LabelForContinue:
				els, _ := iframe.Elements(fmt.Sprintf("li[onclick=\"kjCommonFun.goNextPage(%d)\"]", nextPageNo))
				if len(els) < 1 {
					time.Sleep(1 * time.Second)
					goto LabelForContinue
				}

				iframe.MustElement(fmt.Sprintf("li[onclick=\"kjCommonFun.goNextPage(%d)\"]", nextPageNo)).MustClick()
			}
		}
	}
	lg.InfoToFileAndStdOut(fmt.Sprintf("5 \n"))
	time.Sleep(2 * time.Second)
	lg.InfoToFileAndStdOut(fmt.Sprintf("6 \n"))
	// 关闭打开的浏览器
	_ = browser.Close()
	lg.InfoToFileAndStdOut(fmt.Sprintf("7 \n"))
	return
}

func getLotteryHistory1(gameNo, pageSize, pageNo int, startTerm, endTerm string) (lotteryData LotterySt) {
	lg.InfoToFileAndStdOut(fmt.Sprintf("初始化浏览器中... \n"))
	if cfg.Default.CloseBrowser == 2 {
		//defaults.Show = true
		defaults.ResetWith("show=true")
	}

	browser, err1 := bs.GetBrowser(strconv.Itoa(0))
	if err1 != nil {
		lg.ErrorToFile(fmt.Sprintf("第%d次打开浏览器发生错误：%v\n", 0, err1))
		return
	}

	router := browser.HijackRequests()
	router.MustAdd("https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry*", func(ctx *rod.Hijack) {
		url := ctx.Request.URL().RawPath
		//_ = ctx.LoadResponse(http.DefaultClient, true)
		//ctx.Response.SetBody("")
		// 继续原始请求
		ctx.MustLoadResponse()
		// 获取原始响应
		originalResponse := ctx.Response.Body()

		err := json.Unmarshal([]byte(originalResponse), &lotteryData)
		if err != nil {
			panic(fmt.Sprintf("解析%s返回结果，遇到错误：%v\n", url, err))
		}
	})

	go router.Run()

	url := "https://www.lottery.gov.cn/kj/kjlb.html?dlt"
	page := browser.MustPage(url)
	//if err != nil {
	//	lg.ErrorToFile(fmt.Sprintf("打不开大乐透历史开奖页面：%v", err))
	//	return
	//}
	page.MustWaitLoad()
	page.MustSetViewport(cfg.Default.BrowserWidth, cfg.Default.BrowserHeight, 1, false)

	iframe := page.MustElement("#iFrame1").MustFrame()
	iframe.Eval(`() => {
			document.querySelector(".g-history .u-zdy ul").style.display = 'block';
		}`)
	time.Sleep(2 * time.Second)
	iframe.Eval(fmt.Sprintf(`() => {
			document.getElementById("bterm").value = '%s';
			document.getElementById("eterm").value = '%s';
			document.getElementById("searchBrn").click();
		}`, startTerm, endTerm))

	//iframe.MustElement("#bterm").MustSelectAllText().MustInput("").MustInput(startTerm)
	//lg.InfoToFile(fmt.Sprintf("已填入开始期号"))
	//iframe.MustElement("#eterm").MustSelectAllText().MustInput("").MustInput(endTerm)
	//lg.InfoToFile(fmt.Sprintf("已填入结束期号"))
	//
	//iframe.MustElement("#searchBrn").MustClick()
	lg.InfoToFile(fmt.Sprintf("在大乐透历史开奖页面已点击查询"))
	time.Sleep(4 * time.Second)

	// 关闭打开的浏览器
	_ = browser.Close()

	fmt.Printf("lotteryData=%v\n", lotteryData)
	return
}

// getLotteryHistory2 使用go标准库中的http包获取官网的开奖数据(已经被官网封锁,不能再使用)
func getLotteryHistory2(gameNo, pageSize, pageNo int, startTerm, endTerm string) (lotteryData LotterySt) {
	client := &http.Client{
		Timeout: 20 * time.Second,
		//Transport: &http.Transport{
		//	MaxIdleConns: 100,
		//},
	}

	url := fmt.Sprintf("https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=%d&provinceId=0&pageSize=%d&isVerify=1&pageNo=%d&startTerm=%s&endTerm=%s", gameNo, pageSize, pageNo, startTerm, endTerm)
	// 创建HTTP GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Sprintf("创建请求失败: %v\n", err))
	}

	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	////req.Header.Set("Access-Control-Request-Method", "GET")
	//req.Header.Set("Access-Control-Request-Headers", "x-dev-fp,x-full-ref,x-sensors-id")
	//////req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	//req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	//req.Header.Set("Priority", "u=1, i")
	//req.Header.Set("Origin", "https://static.sporttery.cn")
	//req.Header.Set("Referer", "https://static.sporttery.cn/")
	//req.Header.Set("Sec-Ch-Ua-Platform", "Windows")
	//req.Header.Set("Sec-Fetch-Dest", "")
	//req.Header.Set("Sec-Fetch-Mode", "cors")
	//req.Header.Set("Sec-Fetch-Site", "same-site")
	//req.Header.Set("x-dev-fp", "8e26f46ee3eceeac0b38f819dcce5fc7")
	//req.Header.Set("x-full-ref", "//www.lottery.gov.cn/kj/kjlb.html?dlt")
	//req.Header.Set("x-sensors-id", "196bd9a41b7842-064d65a09f91dfc-26011f51-3686400-196bd9a41b9136b")

	//req.Header.Set(":authority", "webapi.sporttery.cn")
	//req.Header.Set(":method", "OPTIONS")
	//req.Header.Set(":path", fmt.Sprintf("/gateway/lottery/getHistoryPageListV1.qry?gameNo=85&provinceId=0&pageSize=100&isVerify=1&pageNo=%d", pageNo))
	//req.Header.Set(":scheme", "https")
	//req.Header.Set(":authority", "webapi.sporttery.cn")
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("打不开%s，遇到错误：%v\n", url, err))
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		// 处理非200响应（如401、500等）
		panic(fmt.Sprintf("打开%s遇到响应状态码：%d\n", url, resp.StatusCode))
	}

	body, _ := io.ReadAll(resp.Body)
	//fmt.Printf("%v\n", string(body))

	err = json.Unmarshal(body, &lotteryData)
	if err != nil {
		panic(fmt.Sprintf("解析%s返回结果，遇到错误：%v\n", url, err))
	}
	fmt.Printf("已完成第%d页（总页数为%d） %s~%s\n", pageNo, lotteryData.Value.Pages, startTerm, endTerm)
	return
}

// GetStake 从切片中获取对应奖级的中奖个数和奖金
func GetStake(pls []PrizeLevelListItem2, sort int) (stakeCount int, stakeAmount int) {
	var parseErr error
	for _, pl := range pls {
		if pl.Sort == sort {
			stakeCount, parseErr = strconv.Atoi(strings.ReplaceAll(pl.StakeCount, ",", ""))
			if parseErr != nil || stakeCount < 0 {
				stakeCount = 0
			}
			stakeAmount, parseErr = strconv.Atoi(pl.StakeAmountFormat)
			if parseErr != nil || stakeAmount < 0 {
				stakeAmount = 0
			}
			return
		}
	}
	return
}

// GetSomeDltFromWeb 从官网所在网页中获取开奖数据(已整理好的数据)
func GetSomeDltFromWeb(lastDlt models.Dlt, endTerm string) (dlts []models.Dlt, err error) {
	if lastDlt.DrawNum != endTerm {
		startDrawNum := "07001"
		if lastDlt.DrawNum != "" {
			startDrawNum = lastDlt.DrawNum
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("startDrawNum=%s endTerm=%s\n", startDrawNum, endTerm))
		ld := GetLotteryHistory(5, startDrawNum, endTerm)
		var parseErr error
		if len(ld) > 0 {
			// 按照 期号,从小到大
			sort.Slice(ld, func(i, j int) bool {
				return ld[i].LotteryDrawNum < ld[j].LotteryDrawNum
			})
			for _, v := range ld {
				if lastDlt.DrawNum != "" && v.LotteryDrawNum <= lastDlt.DrawNum {
					continue
				}
				newDlt := models.Dlt{}
				newDlt.DrawNum = v.LotteryDrawNum
				newDlt.DrawTime = v.LotteryDrawTime
				newDlt.EquipmentCount = v.LotteryEquipmentCount
				newDlt.DrawPdfUrl = v.DrawPdfUrl
				newDlt.UnSortDrawResult = string(v.LotteryUnSortDrawResult)
				hmStrSlice := strings.Split(v.LotteryDrawResult, " ")
				if len(hmStrSlice) == 7 {
					newDlt.F1 = hmStrSlice[0]
					newDlt.F2 = hmStrSlice[1]
					newDlt.F3 = hmStrSlice[2]
					newDlt.F4 = hmStrSlice[3]
					newDlt.F5 = hmStrSlice[4]
					newDlt.B1 = hmStrSlice[5]
					newDlt.B2 = hmStrSlice[6]
					newDlt.Oe = ana_dlt.CalDltOe(hmStrSlice)
					hz := ana_dlt.CalDltHz(hmStrSlice)
					newDlt.Hz = hz
					newDlt.AeHz = ana_dlt.DltHzABCDE(hz)
					newDlt.Qzh = ana_dlt.CalDltQzh(hmStrSlice)
				}

				newDlt.PoolBalance, parseErr = strconv.ParseFloat(strings.ReplaceAll(v.PoolBalanceAfterDraw, ",", ""), 64)
				if parseErr != nil {
					newDlt.PoolBalance = 0
				}
				newDlt.TotalSaleAmount, parseErr = strconv.ParseFloat(strings.ReplaceAll(v.TotalSaleAmount, ",", ""), 64)
				if parseErr != nil {
					newDlt.TotalSaleAmount = 0
				}

				newDlt.StakeCount60, newDlt.StakeAmount60 = GetStake(v.PrizeLevelList, 60)
				newDlt.StakeCount80, newDlt.StakeAmount80 = GetStake(v.PrizeLevelList, 80)
				newDlt.StakeCount100, newDlt.StakeAmount100 = GetStake(v.PrizeLevelList, 100)
				newDlt.StakeCount101, newDlt.StakeAmount101 = GetStake(v.PrizeLevelList, 101)
				newDlt.StakeCount102, newDlt.StakeAmount102 = GetStake(v.PrizeLevelList, 102)
				newDlt.StakeCount201, newDlt.StakeAmount201 = GetStake(v.PrizeLevelList, 201)
				newDlt.StakeCount202, newDlt.StakeAmount202 = GetStake(v.PrizeLevelList, 202)
				newDlt.StakeCount301, newDlt.StakeAmount301 = GetStake(v.PrizeLevelList, 301)
				newDlt.StakeCount302, newDlt.StakeAmount302 = GetStake(v.PrizeLevelList, 302)
				newDlt.StakeCount401, newDlt.StakeAmount401 = GetStake(v.PrizeLevelList, 401)
				newDlt.StakeCount402, newDlt.StakeAmount402 = GetStake(v.PrizeLevelList, 402)
				newDlt.StakeCount501, newDlt.StakeAmount501 = GetStake(v.PrizeLevelList, 501)
				newDlt.StakeCount601, newDlt.StakeAmount601 = GetStake(v.PrizeLevelList, 601)
				newDlt.StakeCount701, newDlt.StakeAmount701 = GetStake(v.PrizeLevelList, 701)
				newDlt.StakeCount801, newDlt.StakeAmount801 = GetStake(v.PrizeLevelList, 801)
				newDlt.StakeCount901, newDlt.StakeAmount901 = GetStake(v.PrizeLevelList, 901)
				newDlt.StakeCount1001, newDlt.StakeAmount1001 = GetStake(v.PrizeLevelList, 1001)
				newDlt.StakeCount1101, newDlt.StakeAmount1101 = GetStake(v.PrizeLevelList, 1101)

				// 追加到切片
				dlts = append(dlts, newDlt)
			}
		}
	}
	return
}
