package ssq

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/before80/lot/ana_ssq"
	"github.com/before80/lot/bs"
	"github.com/before80/lot/cfg"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/defaults"
	"github.com/spf13/cobra"
)

func UpdateSsq(cmd *cobra.Command) {
	// 从命令行参数中获取ldn的值
	ldn, _ := cmd.Flags().GetString("ldn")
	lg.InfoToFileAndStdOut(fmt.Sprintf("ldn=%s\n", ldn))
	if ldn == "" {
		return
	}

	if len(ldn) != 7 {
		lg.InfoToFileAndStdOut(fmt.Sprintf("ldn长度必须为7位\n"))
		return
	}

	if ldn < "2003001" {
		lg.InfoToFileAndStdOut(fmt.Sprintf("ldn长度必须为7位\n"))
		return
	}

	startTime := time.Now()
	lastSsq := dbop.GetLastSsq()
	//db.DB.Last(&lastSsq)
	lg.InfoToFileAndStdOut(fmt.Sprintf("当前数据库中最新的一条记录为 %v \n", lastSsq))

	needUseMethod1 := false // 是否需要使用 https://datachart.500.com/ssq/history/history.shtml 获取开奖数据
	needUseMethod2 := false // 是否需要使用 https://www.cwl.gov.cn/ygkj/wqkjgg/ssq/ 取开奖数据
	useMethod1StartDrawNum := ""
	useMethod1EndDrawNum := ""

	useMethod2StartDrawNum := ""
	useMethod2EndDrawNum := ""

	if lastSsq.DrawNum == "" {
		needUseMethod1 = true

		if ldn < "2013001" {
			useMethod1StartDrawNum = "03001"
			useMethod1EndDrawNum = ldn[2:]
		}

		if ldn >= "2013001" {
			useMethod1StartDrawNum = "03001"
			useMethod1EndDrawNum = "12154"

			needUseMethod2 = true
			useMethod2StartDrawNum = "2013001"
			useMethod2EndDrawNum = ldn

		}
	}

	if lastSsq.DrawNum != "" {
		if lastSsq.DrawNum < "2013001" && ldn < "2013001" {
			needUseMethod1 = true
			useMethod1StartDrawNum = lastSsq.DrawNum[2:]
			useMethod1EndDrawNum = ldn[2:]
		}

		if lastSsq.DrawNum < "2013001" && ldn >= "2013001" {
			needUseMethod1 = true

			useMethod1StartDrawNum = lastSsq.DrawNum[2:]
			useMethod1EndDrawNum = "12154"

			needUseMethod2 = true
			useMethod2StartDrawNum = "2013001"
			useMethod2EndDrawNum = ldn
		}

		if lastSsq.DrawNum >= "2013001" && ldn >= "2013001" {
			needUseMethod2 = true
			useMethod2StartDrawNum = "2013001"
			useMethod2EndDrawNum = ldn
		}
	}

	lg.InfoToFileAndStdOut(fmt.Sprintf("needUseMethod1=%v  needUseMethod2=%v\n", needUseMethod1, needUseMethod2))
	if needUseMethod1 {
		lg.InfoToFileAndStdOut(fmt.Sprintf("useMethod1StartDrawNum=%s  useMethod1EndDrawNum=%s\n", useMethod1StartDrawNum, useMethod1EndDrawNum))
		ssqs, err := GetSomeSsqFromWeb1(lastSsq, useMethod1StartDrawNum, useMethod1EndDrawNum)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("从网页获取开奖数据出现错误：%v\n", err), 3)
			return
		}
		insertedRow, err := dbop.InsertSsqBatch(ssqs, 100)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("往数据表中插入数据出现错误：%v\n", err), 3)
			return
		} else {
			lg.InfoToFileAndStdOut(fmt.Sprintf("插入了 %d 条数据\n", insertedRow))
			lg.InfoToFileAndStdOut(fmt.Sprintf("程序运行时间：%.2f秒\n", time.Since(startTime).Seconds()))
		}
	}

	if needUseMethod2 {
		lg.InfoToFileAndStdOut(fmt.Sprintf("useMethod2StartDrawNum=%s useMethod2EndDrawNum=%s\n", useMethod2StartDrawNum, useMethod2EndDrawNum))
		if needUseMethod1 { // 注意这里需要重新获取最新的双色球记录
			lastSsq = dbop.GetLastSsq()
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("lastSsq=%v\n", lastSsq))

		ssqs, err := GetSomeSsqFromWeb2(lastSsq, useMethod2StartDrawNum, useMethod2EndDrawNum)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("从网页获取开奖数据出现错误：%v\n", err), 3)
			return
		}
		insertedRow, err := dbop.InsertSsqBatch(ssqs, 100)
		if err != nil {
			lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("往数据表中插入数据出现错误：%v\n", err), 3)
			return
		} else {
			lg.InfoToFileAndStdOut(fmt.Sprintf("插入了 %d 条数据\n", insertedRow))
			lg.InfoToFileAndStdOut(fmt.Sprintf("程序运行时间：%.2f秒\n", time.Since(startTime).Seconds()))
		}
	}
}

func GetSomeSsqFromWeb1(lastSsq models.Ssq, startTerm, endTerm string) (ssqs []models.Ssq, err error) {
	if lastSsq.DrawNum != endTerm {
		if lastSsq.DrawNum > endTerm {
			return ssqs, nil
		}

		ld := GetLotteryHistory1(startTerm, endTerm)
		var parseErr error
		if len(ld) > 0 {
			for _, v := range ld {
				// 过滤掉已经存在于数据库中的记录
				if lastSsq.DrawNum != "" && v.Code <= lastSsq.DrawNum {
					continue
				}
				newSsq := models.Ssq{}
				newSsq.DrawNum = v.Code

				khIndex := strings.Index(v.Date, "(")
				if khIndex != -1 {
					newSsq.DrawTime = v.Date[0:khIndex]
				} else {
					newSsq.DrawTime = v.Date
				}

				newSsq.Week = v.Week

				redHmStr := v.Red
				hmStrSlice := strings.Split(redHmStr, ",")
				hmStrSlice = append(hmStrSlice, v.Blue)
				if len(hmStrSlice) == 7 {
					newSsq.F1 = hmStrSlice[0]
					newSsq.F2 = hmStrSlice[1]
					newSsq.F3 = hmStrSlice[2]
					newSsq.F4 = hmStrSlice[3]
					newSsq.F5 = hmStrSlice[4]
					newSsq.F6 = hmStrSlice[5]
					newSsq.B1 = hmStrSlice[6]
					newSsq.Oe = ana_ssq.CalSsqOe(hmStrSlice)
					hz := ana_ssq.CalSsqHz(hmStrSlice)
					newSsq.Hz = hz
					newSsq.AeHz = ana_ssq.SsqHzABCDE(hz)
					newSsq.Qzh = ana_ssq.CalSsqQzh(hmStrSlice)
				}

				newSsq.PoolBalance, parseErr = strconv.ParseFloat(strings.ReplaceAll(v.PoolMoney, ",", ""), 64)
				if parseErr != nil {
					newSsq.PoolBalance = 0
				}
				newSsq.TotalSaleAmount, parseErr = strconv.ParseFloat(strings.ReplaceAll(v.Sales, ",", ""), 64)
				if parseErr != nil {
					newSsq.TotalSaleAmount = 0
				}
				newSsq.Content = v.Content

				newSsq.StakeCount1, newSsq.StakeAmount1 = GetStake(v.PrizeGrades, 1)
				newSsq.StakeCount2, newSsq.StakeAmount2 = GetStake(v.PrizeGrades, 2)
				newSsq.StakeCount3, newSsq.StakeAmount3 = GetStake(v.PrizeGrades, 3)
				newSsq.StakeCount4, newSsq.StakeAmount4 = GetStake(v.PrizeGrades, 4)
				newSsq.StakeCount5, newSsq.StakeAmount5 = GetStake(v.PrizeGrades, 5)
				newSsq.StakeCount6, newSsq.StakeAmount6 = GetStake(v.PrizeGrades, 6)
				newSsq.VideoUrl = v.VideoLink
				newSsq.DetailsUrl = v.DetailsLink
				// 追加到切片
				ssqs = append(ssqs, newSsq)
			}
		}
	}
	return
}

func GetSomeSsqFromWeb2(lastSsq models.Ssq, startTerm, endTerm string) (ssqs []models.Ssq, err error) {
	if lastSsq.DrawNum != endTerm {
		//startDrawNum := "2013001"
		//if lastSsq.DrawNum != "" && lastSsq.DrawNum >= startDrawNum {
		//	startDrawNum = lastSsq.DrawNum
		//}

		if lastSsq.DrawNum > endTerm {
			return ssqs, nil
		}

		ld := GetLotteryHistory2(startTerm, endTerm)
		//fmt.Printf("%v\n", len(ld))
		//fmt.Printf("%v\n", ld)
		//fmt.Printf("%v\n", GenCrossTermYears("07001", "25053", 5))
		var parseErr error
		if len(ld) > 0 {
			for _, v := range ld {
				// 过滤掉已经存在于数据库中的记录
				if lastSsq.DrawNum != "" && v.Code <= lastSsq.DrawNum {
					continue
				}
				newSsq := models.Ssq{}
				newSsq.DrawNum = v.Code

				khIndex := strings.Index(v.Date, "(")
				if khIndex != -1 {
					newSsq.DrawTime = v.Date[0:khIndex]
				} else {
					newSsq.DrawTime = v.Date
				}

				newSsq.Week = v.Week

				redHmStr := v.Red
				hmStrSlice := strings.Split(redHmStr, ",")
				hmStrSlice = append(hmStrSlice, v.Blue)
				if len(hmStrSlice) == 7 {
					newSsq.F1 = hmStrSlice[0]
					newSsq.F2 = hmStrSlice[1]
					newSsq.F3 = hmStrSlice[2]
					newSsq.F4 = hmStrSlice[3]
					newSsq.F5 = hmStrSlice[4]
					newSsq.F6 = hmStrSlice[5]
					newSsq.B1 = hmStrSlice[6]
					newSsq.Oe = ana_ssq.CalSsqOe(hmStrSlice)
					hz := ana_ssq.CalSsqHz(hmStrSlice)
					newSsq.Hz = hz
					newSsq.AeHz = ana_ssq.SsqHzABCDE(hz)
					newSsq.Qzh = ana_ssq.CalSsqQzh(hmStrSlice)
				}

				newSsq.PoolBalance, parseErr = strconv.ParseFloat(strings.ReplaceAll(v.PoolMoney, ",", ""), 64)
				if parseErr != nil {
					newSsq.PoolBalance = 0
				}
				newSsq.TotalSaleAmount, parseErr = strconv.ParseFloat(strings.ReplaceAll(v.Sales, ",", ""), 64)
				if parseErr != nil {
					newSsq.TotalSaleAmount = 0
				}
				newSsq.Content = v.Content

				newSsq.StakeCount1, newSsq.StakeAmount1 = GetStake(v.PrizeGrades, 1)
				newSsq.StakeCount2, newSsq.StakeAmount2 = GetStake(v.PrizeGrades, 2)
				newSsq.StakeCount3, newSsq.StakeAmount3 = GetStake(v.PrizeGrades, 3)
				newSsq.StakeCount4, newSsq.StakeAmount4 = GetStake(v.PrizeGrades, 4)
				newSsq.StakeCount5, newSsq.StakeAmount5 = GetStake(v.PrizeGrades, 5)
				newSsq.StakeCount6, newSsq.StakeAmount6 = GetStake(v.PrizeGrades, 6)
				newSsq.VideoUrl = v.VideoLink
				newSsq.DetailsUrl = v.DetailsLink
				// 追加到切片
				ssqs = append(ssqs, newSsq)
			}
		}
	}
	return
}

func GetStake(pgs []PrizeGrades, typ int) (stakeCount int, stakeAmount int) {
	var parseErr error
	for _, pg := range pgs {
		if pg.Type == typ {
			stakeCount, parseErr = strconv.Atoi(strings.ReplaceAll(pg.TypeNum, ",", ""))
			if parseErr != nil || stakeCount < 0 {
				stakeCount = 0
			}
			//s := "8333333（含加奖3333333）"
			khIndex := strings.Index(pg.TypeMoney, "（")
			if khIndex != -1 {
				stakeAmount, parseErr = strconv.Atoi(pg.TypeMoney[0:khIndex])
				if parseErr != nil || stakeAmount < 0 {
					stakeAmount = 0
				}
			} else {
				stakeAmount, parseErr = strconv.Atoi(pg.TypeMoney)
				if parseErr != nil || stakeAmount < 0 {
					stakeAmount = 0
				}
			}

			return
		}
	}
	return
}

func GetLotteryHistory1(startTerm, endTerm string) (lotteryHistory []ListItem) {
	lotteryHistory = getLotteryHistory1(startTerm, endTerm)

	// 排序
	slices.SortFunc(lotteryHistory, func(a, b ListItem) int {
		aTerm, _ := strconv.Atoi(a.Code)
		bTerm, _ := strconv.Atoi(b.Code)
		if aTerm > bTerm {
			return 1
		} else if aTerm < bTerm {
			return -1
		}
		return 0
	})

	return
}

func getLotteryHistory1(startTerm, endTerm string) (lotteryDataSlice []ListItem) {
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

	url := "https://datachart.500.com/ssq/history/history.shtml"
	page := browser.MustPage(url)

	page.MustWaitLoad()
	page.MustSetViewport(cfg.Default.BrowserWidth, cfg.Default.BrowserHeight, 1, false)

	//ch := make(chan int)

	//hijackUrl := fmt.Sprintf("https://datachart.500.com/ssq/history/newinc/history.php?start=%s&end=%s", startTerm, endTerm)
	//router := browser.HijackRequests()
	//router.MustAdd(hijackUrl, func(ctx *rod.Hijack) {
	//	ctx.MustLoadResponse()
	//	ch <- 1
	//})

	//go router.Run()

	page.Eval(fmt.Sprintf(`() => {
			document.querySelector("#start").value = '%s';
			document.querySelector("#end").value = '%s';
			document.querySelector("img[onclick^='getChartdata']").click();
		}`, startTerm, endTerm))

	lg.InfoToFile(fmt.Sprintf("在双色球历史开奖页面已点击查询"))

	lg.InfoToFileAndStdOut(fmt.Sprintf("等待6秒前\n"))
	time.Sleep(6 * time.Second)
	lg.InfoToFileAndStdOut(fmt.Sprintf("等待6秒后\n"))
	trs := page.MustElements("#tdata > tr")
	if len(trs) > 0 {
		for _, tr := range trs {
			lotteryData := ListItem{}
			tds := tr.MustElements("td")
			red := ""
			var pg []PrizeGrades
			pg = nil
			ipg1, ipg2 := PrizeGrades{}, PrizeGrades{}
			for i, td := range tds {
				if i == 0 {
					lotteryData.Code = "20" + td.MustText()
				}
				if i >= 1 && i <= 6 {
					if red != "" {
						red = red + "," + td.MustText()
					} else {
						red = td.MustText()
					}
					if i == 6 {
						lotteryData.Red = red
					}
				}
				if i == 7 {
					lotteryData.Blue = td.MustText()
				}
				if i == 9 {
					lotteryData.PoolMoney = td.MustText()
				}
				if i == 10 {
					ipg1.Type = 1
					ipg1.TypeNum = td.MustText()
				}
				if i == 11 {
					ipg1.TypeMoney = td.MustText()
				}
				if i == 12 {
					ipg2.Type = 2
					ipg2.TypeNum = td.MustText()
				}
				if i == 13 {
					ipg2.TypeMoney = td.MustText()
				}

				if i == 14 {
					lotteryData.Sales = td.MustText()
				}
				if i == 15 {
					lotteryData.Date = td.MustText()
					lotteryData.Week, _ = gen.GetWeekdayCN(td.MustText())
				}
			}
			pg = append(pg, ipg1, ipg2)
			lotteryData.PrizeGrades = pg
			lotteryDataSlice = append(lotteryDataSlice, lotteryData)
		}
	}

	// 关闭打开的浏览器
	_ = browser.Close()
	return
}

func GetLotteryHistory2(startTerm, endTerm string) (lotteryHistory []ListItem) {
	lotteryHistory = getLotteryHistory2(startTerm, endTerm)

	// 排序
	slices.SortFunc(lotteryHistory, func(a, b ListItem) int {
		aTerm, _ := strconv.Atoi(a.Code)
		bTerm, _ := strconv.Atoi(b.Code)
		if aTerm > bTerm {
			return 1
		} else if aTerm < bTerm {
			return -1
		}
		return 0
	})

	return
}

func getLotteryHistory2(startTerm, endTerm string) (lotteryDataSlice []ListItem) {
	lg.InfoToFileAndStdOut(fmt.Sprintf("%s~%s 初始化浏览器中... \n", startTerm, endTerm))
	if cfg.Default.CloseBrowser == 2 {
		//defaults.Show = true
		defaults.ResetWith("show=true")
	}

	browser, err1 := bs.GetBrowser(strconv.Itoa(0))
	if err1 != nil {
		lg.ErrorToFile(fmt.Sprintf("第%d次打开浏览器发生错误：%v\n", 0, err1))
		return
	}

	var dataMu sync.Mutex

	chNextPageNo := make(chan int, 1)
	var closeOnce sync.Once
	hadExistDrawNum := make(map[string]struct{})

	appendFunc := func(data []ListItem) {
		for _, v := range data {
			if v.Code > endTerm {
				continue
			}
			if _, ok := hadExistDrawNum[v.Code]; !ok {
				hadExistDrawNum[v.Code] = struct{}{}
				lotteryDataSlice = append(lotteryDataSlice, v)
			}
		}
	}

	safeClose := func() {
		closeOnce.Do(func() {
			close(chNextPageNo)
		})
	}

	// https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=ssq&issueCount=&issueStart=&issueEnd=&dayStart=&dayEnd=&pageNo=1&pageSize=30&week=&systemType=PC
	router := browser.HijackRequests()
	router.MustAdd("https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice*", func(ctx *rod.Hijack) {
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
		appendFunc(lotteryData.Result)
		dataMu.Unlock()
		// ======================

		if lotteryData.PageNo >= lotteryData.PageNum {
			// 已是最后一页
			safeClose()
			return
		}

		nextPageNo := lotteryData.PageNo + 1

		// 非阻塞写，避免 channel 已关闭或满导致死锁
		select {
		case chNextPageNo <- nextPageNo:
		default:
		}

	})

	go router.Run()

	url := "https://www.cwl.gov.cn/ygkj/wqkjgg/ssq/"
	page := browser.MustPage(url)

	page.MustWaitLoad()
	page.MustSetViewport(cfg.Default.BrowserWidth, cfg.Default.BrowserHeight, 1, false)

	page.Eval(`() => { handleShowCustom(); }`)
	//page.MustElement(`div[onclick="handleShowCustom()"]`).MustClick()
	time.Sleep(100 * time.Millisecond)
	page.MustElement("div.item-content")

	page.Eval(fmt.Sprintf(`() => {
			document.querySelector(".issue-start").value = '%s';
			document.querySelector(".issue-end").value = '%s';
			document.querySelector(".custom-query").click();
		}`, startTerm, endTerm))

	lg.InfoToFile(fmt.Sprintf("在双色球历史开奖页面已点击查询"))

	for nextPageNo := range chNextPageNo {
	LabelForContinue:
		els, _ := page.Elements(fmt.Sprintf("a[data-page='%d']", nextPageNo))
		if len(els) < 1 {
			time.Sleep(1 * time.Second)
			goto LabelForContinue
		}
		page.MustElement(fmt.Sprintf("a[data-page='%d']", nextPageNo)).MustClick()
	}

	time.Sleep(3 * time.Second)

	// 关闭打开的浏览器
	_ = browser.Close()
	return
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
		curStartYear, _ := strconv.Atoi(curEndTerm[0:4])
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
