package ana_dlt

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/before80/lot/db"
	"github.com/before80/lot/excel"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
	"github.com/xuri/excelize/v2"
)

type BestCombSt struct {
	ColCount        int      // 组合的列数
	Tx              []string // 各个组合类型
	HisFuGaiDrawNum int      // 历史覆盖期数
	FrontHms        []string // 各个组合类型包含的前区号码的并集(需去重)
	FrontHmCount    int      // 各个组合类型包含的前区号码的并集(需去重)后的个数
	GaiLvTs         float64  // 概率提升多少倍
}

// DltDataToExcel 整理大乐透相关数据到Excel表中
func DltDataToExcel(prevRunMoni bool) {
	//var m1, m2 runtime.MemStats

	// 强制 GC，减少历史干扰（非常重要）
	//runtime.GC()
	//runtime.ReadMemStats(&m1)
	startTime := time.Now()
	UpdateMonisTable()
	if prevRunMoni {
		lg.InfoToFile(fmt.Sprintf("运行模拟数据中,请耐心等待...\n"))
		BatchMoni(3)
		lg.InfoToFile(fmt.Sprintf("运行模拟数据所需要的时间: %v\n", time.Now().Sub(startTime).Round(time.Second)))
	}

	// 找到类型为77777匹配历史开奖号码最多的前5种
	var moni77777s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "77777"}).Order("cs desc,id asc").Limit(5).Find(&moni77777s)
	// 找到类型为116666匹配历史开奖号码最多的前5种
	var moni116666s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "116666"}).Order("cs desc,id asc").Limit(5).Find(&moni116666s)
	// 找到类型为155555匹配历史开奖号码最多的前5种
	var moni155555s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "155555"}).Order("cs desc,id asc").Limit(5).Find(&moni155555s)
	// 找到类型为194444匹配历史开奖号码最多的前2种
	var moni194444s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "194444"}).Order("cs desc,id asc").Limit(2).Find(&moni194444s)

	//fmt.Println(moni77777s)
	//fmt.Println(moni116666s)
	//fmt.Println(moni155555s)
	//fmt.Println(moni194444s)
	hm7s := BuildFrontHmSlices(moni77777s, 5)
	hm11s := BuildFrontHmSlices(moni116666s, 5)
	hm15s := BuildFrontHmSlices(moni155555s, 5)
	hm19s := BuildFrontHmSlices(moni194444s, 2)
	otherTxHms := CalDltNotTxFrontHmSlice(hm7s, hm11s, hm15s, hm19s)
	//fmt.Println(len(hm7s[0]), len(hm7s[1]), len(hm7s[2]), len(hm7s[3]), len(hm7s[4]))
	//fmt.Println("-----------------")
	//fmt.Println(len(hm11s[0]), len(hm11s[1]), len(hm11s[2]), len(hm11s[3]), len(hm11s[4]))
	//fmt.Println("-----------------")
	//fmt.Println(len(hm15s[0]), len(hm15s[1]), len(hm15s[2]), len(hm15s[3]), len(hm15s[4]))
	//fmt.Println("-----------------")
	//fmt.Println(len(hm19s[0]), len(hm19s[1]))
	//fmt.Println("-----------------")
	//fmt.Printf("1 %v\n", time.Now().Sub(startTime).Round(time.Second))
	InitDlts()
	//fmt.Printf("2 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 计算重号累加
	drawNum2CHongHaoSt := CHongHaoLeiJia()
	dlts := DxDlts
	lenDlts := len(dlts)
	var dltExcelDatas []DltExcelData
	tr7 := make([][]int, 5)
	tr11 := make([][]int, 5)
	tr15 := make([][]int, 5)
	tr19 := make([][]int, 2)
	otherTCs := 0

	for xuHao, dlt := range dlts {
		rowNum := xuHao + 1
		frontHm := fmt.Sprintf("%s,%s,%s,%s,%s", dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5)
		hz := CalDltHz([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2})
		oe := CalDltOe([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2})
		t7 := make([]int, 5)
		for i, hm := range hm7s {
			if i >= len(t7) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t7[i] = 1
				tr7[i] = append(tr7[i], rowNum)
			}
		}

		t11 := make([]int, 5)
		for i, hm := range hm11s {
			if i >= len(t11) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t11[i] = 1
				tr11[i] = append(tr11[i], rowNum)
			}
		}

		t15 := make([]int, 5)
		for i, hm := range hm15s {
			if i >= len(t15) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t15[i] = 1
				tr15[i] = append(tr15[i], rowNum)
			}
		}

		t19 := make([]int, 2)
		for i, hm := range hm19s {
			if i >= len(t19) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t19[i] = 1
				tr19[i] = append(tr19[i], rowNum)
			}
		}

		otherT := 0
		if slices.Contains(otherTxHms, frontHm) {
			otherT = 1
			otherTCs++
		}

		dltExcelDatas = append(dltExcelDatas, DltExcelData{
			XuHao:               lenDlts - xuHao,
			DrawNum:             dlt.DrawNum,
			DrawTime:            dlt.DrawTime,
			EquipmentCount:      dlt.EquipmentCount,
			FrontHm:             frontHm,
			FullHm:              fmt.Sprintf("%s|%s,%s", frontHm, dlt.B1, dlt.B2),
			UnSortDrawResult:    dlt.UnSortDrawResult,
			PoolBalance:         dlt.PoolBalance,
			TotalSaleAmount:     dlt.TotalSaleAmount,
			StakeCount101:       dlt.StakeCount101,
			StakeAmount101:      dlt.StakeAmount101,
			StakeCount201:       dlt.StakeCount201,
			StakeAmount201:      dlt.StakeAmount201,
			StakeCount301:       dlt.StakeCount301,
			StakeAmount301:      dlt.StakeAmount301,
			StakeCount401:       dlt.StakeCount401,
			StakeAmount401:      dlt.StakeAmount401,
			Oe:                  oe,
			Hz:                  hz,
			Qzh:                 dlt.Qzh,
			NewAddCh4:           drawNum2CHongHaoSt[dlt.DrawNum].NewAddCh4,
			NewAddCh5:           drawNum2CHongHaoSt[dlt.DrawNum].NewAddCh5,
			NewAddCh6:           drawNum2CHongHaoSt[dlt.DrawNum].NewAddCh6,
			NewAddCh7:           drawNum2CHongHaoSt[dlt.DrawNum].NewAddCh7,
			DangQiTotalNewAddCh: drawNum2CHongHaoSt[dlt.DrawNum].DangQiTotalNewAddCh,
			LeiJiaCh:            drawNum2CHongHaoSt[dlt.DrawNum].LeiJiaCh,
			T71:                 t7[0], T72: t7[1], T73: t7[2], T74: t7[3], T75: t7[4],
			T111: t11[0], T112: t11[1], T113: t11[2], T114: t11[3], T115: t11[4],
			T151: t15[0], T152: t15[1], T153: t15[2], T154: t15[3], T155: t15[4],
			T191: t19[0], T192: t19[1],
			OtherT: otherT,
		})
	}

	rowCount := len(dltExcelDatas)
	data := map[string][]int{
		"T1":  tr7[0],
		"T2":  tr7[1],
		"T3":  tr7[2],
		"T4":  tr7[3],
		"T5":  tr7[4],
		"T6":  tr11[0],
		"T7":  tr11[1],
		"T8":  tr11[2],
		"T9":  tr11[3],
		"T10": tr11[4],
		"T11": tr15[0],
		"T12": tr15[1],
		"T13": tr15[2],
		"T14": tr15[3],
		"T15": tr15[4],
		"T16": tr19[0],
		"T17": tr19[1],
	}

	// 自动生成稳定的列名顺序
	colNames := make([]string, 0, len(data))
	for k := range data {
		colNames = append(colNames, k)
	}
	sort.Strings(colNames)

	// 构建列的 bitmask blocks
	blockCount := (rowCount + 63) / 64
	colsMasks := make([][]uint64, 0, len(colNames))

	for _, name := range colNames {
		blocks := make([]uint64, blockCount)
		for _, r := range data[name] {
			r-- // 0-based
			block := r / 64
			bit := uint(r % 64)
			blocks[block] |= 1 << bit
		}
		colsMasks = append(colsMasks, blocks)
	}

	//fmt.Printf("3 %v\n", time.Now().Sub(startTime).Round(time.Second))

	var hm7Combs []string
	for i, hms := range hm7s {
		if slices.Contains([]int{0, 3, 4}, i) {
			for _, hm := range hms {
				if !slices.Contains(hm7Combs, hm) {
					hm7Combs = append(hm7Combs, hm)
				}
			}
		}
	}

	//fmt.Println(len(hm7Combs))
	oeMs := make(map[string]int)
	for _, hm7Comb := range hm7Combs {
		hmCombs := strings.Split(hm7Comb, ",")
		oNum, eNum := 0, 0
		for _, f := range hmCombs {
			fi, _ := strconv.Atoi(f)
			if fi%2 == 0 {
				eNum++
			} else {
				oNum++
			}
		}
		oeStr := fmt.Sprintf("%d%d", oNum, eNum)
		if _, ok := oeMs[oeStr]; !ok {
			oeMs[oeStr] = 1
		} else {
			oeMs[oeStr]++
		}
	}
	//fmt.Println(oeMs)
	//fmt.Printf("4 %v\n", time.Now().Sub(startTime).Round(time.Second))
	saveDir := "./output/excel"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建目录失败: %v", err))
		return
	}

	fileName := filepath.Join(saveDir, fmt.Sprintf("大乐透截止至%s的数据分析", dlts[0].DrawTime))

	// TODO 判断文件是否已经存在, 若存在则删除后

	f, err := excel.CreateNewExcelFile(fileName)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建工作簿过程中出现错误：%v\n", err))
		return
	}

	defer func() {
		if err = f.Close(); err != nil {
			lg.ErrorToFile(fmt.Sprintf("关闭excel文件出现错误：%v", err))
		}
	}()

	defer func() {
		_ = f.Save()
	}()

	styleID, _ := f.NewStyle(&excelize.Style{
		//Font: &excelize.Font{Family: "Consolas", Size: 8},
		Font: &excelize.Font{Family: "monospace", Size: 8},
		Alignment: &excelize.Alignment{
			WrapText: false,
		},
	})

	greenStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#00FF00"},
			Pattern: 1,
		},
	})
	_ = greenStyle

	yellowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFF00"},
			Pattern: 1,
		},
	})
	_ = yellowStyle
	colorSlice := []int{greenStyle, yellowStyle}

	err = SetFirstSheetContent(f, colorSlice, styleID, "历史开奖", dltExcelDatas)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("5 %v\n", time.Now().Sub(startTime).Round(time.Second))

	typ2FrontHms := make(map[string][]string)
	combTyp2FrontHms := make(map[string][]string)
	tNum := 1
	for _, hms := range hm7s {
		rand.Shuffle(len(hms), func(i, j int) {
			hms[i], hms[j] = hms[j], hms[i]
		})
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}
	for _, hms := range hm11s {
		rand.Shuffle(len(hms), func(i, j int) {
			hms[i], hms[j] = hms[j], hms[i]
		})
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}
	for _, hms := range hm15s {
		rand.Shuffle(len(hms), func(i, j int) {
			hms[i], hms[j] = hms[j], hms[i]
		})
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}
	for _, hms := range hm19s {
		rand.Shuffle(len(hms), func(i, j int) {
			hms[i], hms[j] = hms[j], hms[i]
		})
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}

	//fmt.Println("列顺序:", colNames)

	bestCombSts := make([]*BestCombSt, 0)
	// 要选择 K 列
	for K := 1; K <= 17; K++ {
		bestCombos, bestCover := maxCoverageAll(colsMasks, colNames, K)

		//fmt.Println(K)
		//fmt.Printf("最大覆盖行数: %d 最佳组合数量: %d\n", bestCover, len(bestCombos))
		//fmt.Println("最佳组合如下:")

		for _, combo := range bestCombos {
			//fmt.Printf("\t%v\n", combo)
			frontHms := make([]string, 0, len(combo[0]))
			for _, tName := range combo {
				frontHms = append(frontHms, combTyp2FrontHms[tName]...)
			}
			totalFrontHms := gen.UniqueStrSlice(frontHms)
			//sort.Strings(totalFrontHms)
			bestCombSts = append(bestCombSts, &BestCombSt{
				ColCount:        K,
				Tx:              combo,
				HisFuGaiDrawNum: bestCover,
				FrontHms:        totalFrontHms,
				FrontHmCount:    len(totalFrontHms),
				GaiLvTs:         math.Trunc(float64(324632)/float64(len(totalFrontHms))*100) / 100,
			})
		}
	}

	err = SaveBestCombToExcelFile(f, colorSlice, styleID, "最佳组合", bestCombSts)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("6 %v\n", time.Now().Sub(startTime).Round(time.Second))

	typ2ae := make(map[string]models.DltMoni)
	tNum = 1
	for _, moni := range moni77777s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}
	for _, moni := range moni116666s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}
	for _, moni := range moni155555s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}
	for _, moni := range moni194444s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}

	typ2ae["OtherT"] = models.DltMoni{
		Cs:   otherTCs,
		Comb: len(otherTxHms),
	}

	err = SaveTypInfoToExcelFile(f, colorSlice, styleID, "相关类型", typ2ae)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("7 %v\n", time.Now().Sub(startTime).Round(time.Second))

	var dltBackHisData, dltOeHisData, dltHzHisData, dltQzhHisData, dltFrontOnlyOneHisData, dltBackOnlyOneHisData []DltHis
	var wg sync.WaitGroup
	hisTypCount := 6
	for i := 0; i < hisTypCount; i++ {
		wg.Add(1)
		if i == 0 {
			go func() {
				dltBackHisData = DltBackHis(&wg)
			}()

		}
		if i == 1 {
			go func() {
				dltOeHisData = DltOeHis(&wg)
			}()
		}
		if i == 2 {
			go func() {
				dltHzHisData = DltHzHis(&wg)
			}()
		}
		if i == 3 {
			go func() {
				dltQzhHisData = DltQzhHis(&wg)
			}()
		}
		if i == 4 {
			go func() {
				dltFrontOnlyOneHisData = DltFrontOnlyOneHis(&wg)
			}()
		}

		if i == 5 {
			go func() {
				dltBackOnlyOneHisData = DltBackOnlyOneHis(&wg)
			}()
		}
	}
	wg.Wait()
	//fmt.Printf("8 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 前区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "前区单号历史", dltFrontOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("9 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 后区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "后区单号历史", dltBackOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("10 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 后区历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "后区历史", dltBackHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("11 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 奇偶历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "奇偶历史", dltOeHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("12 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 和值历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "和值历史", dltHzHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("13 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 前中后历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "前中后历史", dltQzhHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("14 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 后区趋势1
	err = SaveBackQuShi1ToExcelFile(f, colorSlice, styleID, "后区趋势1", DltBackQuShi1())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("15 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 后区趋势2
	err = SaveBackQuShi2ToExcelFile(f, colorSlice, styleID, "后区趋势2", DltBackQuShi2())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("16 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 奇偶趋势
	err = SaveOeQuShiToExcelFile(f, colorSlice, styleID, "奇偶趋势", DltOeQuShi())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("17 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 和值趋势1
	err = SaveHzQuShi1ToExcelFile(f, colorSlice, styleID, "和值趋势1", DltHzQuShi1())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("18 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 和值趋势2
	err = SaveHzQuShi2ToExcelFile(f, colorSlice, styleID, "和值趋势2", DltHzQuShi2())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("19 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 前中后趋势
	err = SaveQzhQuShiToExcelFile(f, colorSlice, styleID, "前中后趋势", DltQzhQuShi())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("20 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 重号历史
	err = SaveCHongHaoToExcelFile(f, colorSlice, styleID, "重号历史", DltCHongHao())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("21 %v\n", time.Now().Sub(startTime).Round(time.Second))

	tNum = 1
	for tNum <= len(typ2FrontHms) {
		typName := fmt.Sprintf("类型T%d", tNum)
		_ = SaveTypHmsToExcelFile(f, colorSlice, styleID, typName, typ2FrontHms[typName])
		tNum++
	}
	//fmt.Printf("22 %v\n", time.Now().Sub(startTime).Round(time.Second))
	typName := fmt.Sprintf("OtherT")
	_ = SaveTypHmsToExcelFile(f, colorSlice, styleID, typName, otherTxHms)
	//fmt.Printf("23 %v\n", time.Now().Sub(startTime).Round(time.Second))

	//runtime.ReadMemStats(&m2)
	//
	//fmt.Printf("Alloc 变化: %d MB\n", (m2.Alloc-m1.Alloc)/1024/1024)
	//fmt.Printf("TotalAlloc 变化: %d MB\n", (m2.TotalAlloc-m1.TotalAlloc)/1024/1024)
	//fmt.Printf("Sys 变化: %d MB\n", (m2.Sys-m1.Sys)/1024/1024)
	//fmt.Printf("HeapAlloc 变化: %d MB\n", (m2.HeapAlloc-m1.HeapAlloc)/1024/1024)
}

// DltMoreDataToExcel 整理大乐透更多的相关数据到Excel表中(包括以设备号为基准的相关数据)
func DltMoreDataToExcel(prevRunMoni bool) {
	//var m1, m2 runtime.MemStats

	// 强制 GC，减少历史干扰（非常重要）
	//runtime.GC()
	//runtime.ReadMemStats(&m1)
	startTime := time.Now()
	UpdateMonisTable()
	if prevRunMoni {
		lg.InfoToFile(fmt.Sprintf("运行模拟数据中,请耐心等待...\n"))
		BatchMoni(3)
		lg.InfoToFile(fmt.Sprintf("运行模拟数据所需要的时间: %v\n", time.Now().Sub(startTime).Round(time.Second)))
	}

	// 找到类型为77777匹配历史开奖号码最多的前5种
	var moni77777s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "77777"}).Order("cs desc,id asc").Limit(5).Find(&moni77777s)
	// 找到类型为116666匹配历史开奖号码最多的前5种
	var moni116666s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "116666"}).Order("cs desc,id asc").Limit(5).Find(&moni116666s)
	// 找到类型为155555匹配历史开奖号码最多的前5种
	var moni155555s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "155555"}).Order("cs desc,id asc").Limit(5).Find(&moni155555s)
	// 找到类型为194444匹配历史开奖号码最多的前2种
	var moni194444s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "194444"}).Order("cs desc,id asc").Limit(2).Find(&moni194444s)

	//fmt.Println(moni77777s)
	//fmt.Println(moni116666s)
	//fmt.Println(moni155555s)
	//fmt.Println(moni194444s)
	hm7s := BuildFrontHmSlices(moni77777s, 5)
	hm11s := BuildFrontHmSlices(moni116666s, 5)
	hm15s := BuildFrontHmSlices(moni155555s, 5)
	hm19s := BuildFrontHmSlices(moni194444s, 2)
	otherTxHms := CalDltNotTxFrontHmSlice(hm7s, hm11s, hm15s, hm19s)
	//fmt.Println(len(hm7s[0]), len(hm7s[1]), len(hm7s[2]), len(hm7s[3]), len(hm7s[4]))
	//fmt.Println("-----------------")
	//fmt.Println(len(hm11s[0]), len(hm11s[1]), len(hm11s[2]), len(hm11s[3]), len(hm11s[4]))
	//fmt.Println("-----------------")
	//fmt.Println(len(hm15s[0]), len(hm15s[1]), len(hm15s[2]), len(hm15s[3]), len(hm15s[4]))
	//fmt.Println("-----------------")
	//fmt.Println(len(hm19s[0]), len(hm19s[1]))
	//fmt.Println("-----------------")
	//fmt.Printf("1 %v\n", time.Now().Sub(startTime).Round(time.Second))
	InitDlts()
	//fmt.Printf("2 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 计算重号累加
	drawNum2CHongHaoSt := CHongHaoLeiJia()
	dlts := DxDlts
	lenDlts := len(dlts)
	var dltExcelDatas []DltExcelData
	tr7 := make([][]int, 5)
	tr11 := make([][]int, 5)
	tr15 := make([][]int, 5)
	tr19 := make([][]int, 2)
	otherTCs := 0

	for xuHao, dlt := range dlts {
		rowNum := xuHao + 1
		frontHm := fmt.Sprintf("%s,%s,%s,%s,%s", dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5)
		hz := CalDltHz([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2})
		oe := CalDltOe([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2})
		t7 := make([]int, 5)
		for i, hm := range hm7s {
			if i >= len(t7) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t7[i] = 1
				tr7[i] = append(tr7[i], rowNum)
			}
		}

		t11 := make([]int, 5)
		for i, hm := range hm11s {
			if i >= len(t11) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t11[i] = 1
				tr11[i] = append(tr11[i], rowNum)
			}
		}

		t15 := make([]int, 5)
		for i, hm := range hm15s {
			if i >= len(t15) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t15[i] = 1
				tr15[i] = append(tr15[i], rowNum)
			}
		}

		t19 := make([]int, 2)
		for i, hm := range hm19s {
			if i >= len(t19) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t19[i] = 1
				tr19[i] = append(tr19[i], rowNum)
			}
		}

		otherT := 0
		if slices.Contains(otherTxHms, frontHm) {
			otherT = 1
			otherTCs++
		}

		dltExcelDatas = append(dltExcelDatas, DltExcelData{
			XuHao:               lenDlts - xuHao,
			DrawNum:             dlt.DrawNum,
			DrawTime:            dlt.DrawTime,
			EquipmentCount:      dlt.EquipmentCount,
			FrontHm:             frontHm,
			FullHm:              fmt.Sprintf("%s|%s,%s", frontHm, dlt.B1, dlt.B2),
			UnSortDrawResult:    dlt.UnSortDrawResult,
			PoolBalance:         dlt.PoolBalance,
			TotalSaleAmount:     dlt.TotalSaleAmount,
			StakeCount101:       dlt.StakeCount101,
			StakeAmount101:      dlt.StakeAmount101,
			StakeCount201:       dlt.StakeCount201,
			StakeAmount201:      dlt.StakeAmount201,
			StakeCount301:       dlt.StakeCount301,
			StakeAmount301:      dlt.StakeAmount301,
			StakeCount401:       dlt.StakeCount401,
			StakeAmount401:      dlt.StakeAmount401,
			Oe:                  oe,
			Hz:                  hz,
			Qzh:                 dlt.Qzh,
			NewAddCh4:           drawNum2CHongHaoSt[dlt.DrawNum].NewAddCh4,
			NewAddCh5:           drawNum2CHongHaoSt[dlt.DrawNum].NewAddCh5,
			NewAddCh6:           drawNum2CHongHaoSt[dlt.DrawNum].NewAddCh6,
			NewAddCh7:           drawNum2CHongHaoSt[dlt.DrawNum].NewAddCh7,
			DangQiTotalNewAddCh: drawNum2CHongHaoSt[dlt.DrawNum].DangQiTotalNewAddCh,
			LeiJiaCh:            drawNum2CHongHaoSt[dlt.DrawNum].LeiJiaCh,
			T71:                 t7[0], T72: t7[1], T73: t7[2], T74: t7[3], T75: t7[4],
			T111: t11[0], T112: t11[1], T113: t11[2], T114: t11[3], T115: t11[4],
			T151: t15[0], T152: t15[1], T153: t15[2], T154: t15[3], T155: t15[4],
			T191: t19[0], T192: t19[1],
			OtherT: otherT,
		})
	}

	rowCount := len(dltExcelDatas)
	data := map[string][]int{
		"T1":  tr7[0],
		"T2":  tr7[1],
		"T3":  tr7[2],
		"T4":  tr7[3],
		"T5":  tr7[4],
		"T6":  tr11[0],
		"T7":  tr11[1],
		"T8":  tr11[2],
		"T9":  tr11[3],
		"T10": tr11[4],
		"T11": tr15[0],
		"T12": tr15[1],
		"T13": tr15[2],
		"T14": tr15[3],
		"T15": tr15[4],
		"T16": tr19[0],
		"T17": tr19[1],
	}

	// 自动生成稳定的列名顺序
	colNames := make([]string, 0, len(data))
	for k := range data {
		colNames = append(colNames, k)
	}
	sort.Strings(colNames)

	// 构建列的 bitmask blocks
	blockCount := (rowCount + 63) / 64
	colsMasks := make([][]uint64, 0, len(colNames))

	for _, name := range colNames {
		blocks := make([]uint64, blockCount)
		for _, r := range data[name] {
			r-- // 0-based
			block := r / 64
			bit := uint(r % 64)
			blocks[block] |= 1 << bit
		}
		colsMasks = append(colsMasks, blocks)
	}

	//fmt.Printf("3 %v\n", time.Now().Sub(startTime).Round(time.Second))

	var hm7Combs []string
	for i, hms := range hm7s {
		if slices.Contains([]int{0, 3, 4}, i) {
			for _, hm := range hms {
				if !slices.Contains(hm7Combs, hm) {
					hm7Combs = append(hm7Combs, hm)
				}
			}
		}
	}

	//fmt.Println(len(hm7Combs))
	oeMs := make(map[string]int)
	for _, hm7Comb := range hm7Combs {
		hmCombs := strings.Split(hm7Comb, ",")
		oNum, eNum := 0, 0
		for _, f := range hmCombs {
			fi, _ := strconv.Atoi(f)
			if fi%2 == 0 {
				eNum++
			} else {
				oNum++
			}
		}
		oeStr := fmt.Sprintf("%d%d", oNum, eNum)
		if _, ok := oeMs[oeStr]; !ok {
			oeMs[oeStr] = 1
		} else {
			oeMs[oeStr]++
		}
	}
	//fmt.Println(oeMs)
	//fmt.Printf("4 %v\n", time.Now().Sub(startTime).Round(time.Second))
	saveDir := "./output/excel"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建目录失败: %v", err))
		return
	}

	fileName := filepath.Join(saveDir, fmt.Sprintf("大乐透截止至%s的数据分析_more", dlts[0].DrawTime))

	// TODO 判断文件是否已经存在, 若存在则删除后

	f, err := excel.CreateNewExcelFile(fileName)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建工作簿过程中出现错误：%v\n", err))
		return
	}

	defer func() {
		if err = f.Close(); err != nil {
			lg.ErrorToFile(fmt.Sprintf("关闭excel文件出现错误：%v", err))
		}
	}()

	defer func() {
		_ = f.Save()
	}()

	styleID, _ := f.NewStyle(&excelize.Style{
		//Font: &excelize.Font{Family: "Consolas", Size: 8},
		Font: &excelize.Font{Family: "monospace", Size: 8},
		Alignment: &excelize.Alignment{
			WrapText: false,
		},
	})

	greenStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#00FF00"},
			Pattern: 1,
		},
	})
	_ = greenStyle

	yellowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFF00"},
			Pattern: 1,
		},
	})
	_ = yellowStyle
	colorSlice := []int{greenStyle, yellowStyle}

	err = SetFirstSheetContent(f, colorSlice, styleID, "历史开奖", dltExcelDatas)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("5 %v\n", time.Now().Sub(startTime).Round(time.Second))

	typ2FrontHms := make(map[string][]string)
	combTyp2FrontHms := make(map[string][]string)
	tNum := 1
	for _, hms := range hm7s {
		rand.Shuffle(len(hms), func(i, j int) {
			hms[i], hms[j] = hms[j], hms[i]
		})
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}
	for _, hms := range hm11s {
		rand.Shuffle(len(hms), func(i, j int) {
			hms[i], hms[j] = hms[j], hms[i]
		})
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}
	for _, hms := range hm15s {
		rand.Shuffle(len(hms), func(i, j int) {
			hms[i], hms[j] = hms[j], hms[i]
		})
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}
	for _, hms := range hm19s {
		rand.Shuffle(len(hms), func(i, j int) {
			hms[i], hms[j] = hms[j], hms[i]
		})
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}

	//fmt.Println("列顺序:", colNames)

	bestCombSts := make([]*BestCombSt, 0)
	// 要选择 K 列
	for K := 1; K <= 17; K++ {
		bestCombos, bestCover := maxCoverageAll(colsMasks, colNames, K)

		//fmt.Println(K)
		//fmt.Printf("最大覆盖行数: %d 最佳组合数量: %d\n", bestCover, len(bestCombos))
		//fmt.Println("最佳组合如下:")

		for _, combo := range bestCombos {
			//fmt.Printf("\t%v\n", combo)
			frontHms := make([]string, 0, len(combo[0]))
			for _, tName := range combo {
				frontHms = append(frontHms, combTyp2FrontHms[tName]...)
			}
			totalFrontHms := gen.UniqueStrSlice(frontHms)
			//sort.Strings(totalFrontHms)
			bestCombSts = append(bestCombSts, &BestCombSt{
				ColCount:        K,
				Tx:              combo,
				HisFuGaiDrawNum: bestCover,
				FrontHms:        totalFrontHms,
				FrontHmCount:    len(totalFrontHms),
				GaiLvTs:         math.Trunc(float64(324632)/float64(len(totalFrontHms))*100) / 100,
			})
		}
	}

	err = SaveBestCombToExcelFile(f, colorSlice, styleID, "最佳组合", bestCombSts)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("6 %v\n", time.Now().Sub(startTime).Round(time.Second))

	typ2ae := make(map[string]models.DltMoni)
	tNum = 1
	for _, moni := range moni77777s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}
	for _, moni := range moni116666s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}
	for _, moni := range moni155555s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}
	for _, moni := range moni194444s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}

	typ2ae["OtherT"] = models.DltMoni{
		Cs:   otherTCs,
		Comb: len(otherTxHms),
	}

	err = SaveTypInfoToExcelFile(f, colorSlice, styleID, "相关类型", typ2ae)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("7 %v\n", time.Now().Sub(startTime).Round(time.Second))

	var dltBackHisData, dltOeHisData, dltHzHisData, dltQzhHisData, dltFrontOnlyOneHisData, dltBackOnlyOneHisData []DltHis
	var dltEq1FrontOnlyOneHisData, dltEq2FrontOnlyOneHisData, dltEq3FrontOnlyOneHisData []DltHis
	var dltEq1BackOnlyOneHisData, dltEq2BackOnlyOneHisData, dltEq3BackOnlyOneHisData []DltHis
	var dltEq1BackHisData, dltEq2BackHisData, dltEq3BackHisData []DltHis
	var dltEq1OeHisData, dltEq2OeHisData, dltEq3OeHisData []DltHis
	var dltEq1HzHisData, dltEq2HzHisData, dltEq3HzHisData []DltHis
	var dltEq1QzhHisData, dltEq2QzhHisData, dltEq3QzhHisData []DltHis

	var wg sync.WaitGroup
	hisTypCount := 24
	for i := 0; i < hisTypCount; i++ {
		wg.Add(1)
		if i == 0 {
			go func() {
				dltBackHisData = DltBackHis(&wg)
			}()
		}
		if i == 1 {
			go func() {
				dltOeHisData = DltOeHis(&wg)
			}()
		}
		if i == 2 {
			go func() {
				dltHzHisData = DltHzHis(&wg)
			}()
		}
		if i == 3 {
			go func() {
				dltQzhHisData = DltQzhHis(&wg)
			}()
		}
		if i == 4 {
			go func() {
				dltFrontOnlyOneHisData = DltFrontOnlyOneHis(&wg)
			}()
		}

		if i == 5 {
			go func() {
				dltBackOnlyOneHisData = DltBackOnlyOneHis(&wg)
			}()
		}

		if i == 6 {
			go func() {
				dltEq1FrontOnlyOneHisData = DltEqFrontOnlyOneHis(&wg, 1)
			}()
		}
		if i == 7 {
			go func() {
				dltEq2FrontOnlyOneHisData = DltEqFrontOnlyOneHis(&wg, 2)
			}()
		}

		if i == 8 {
			go func() {
				dltEq3FrontOnlyOneHisData = DltEqFrontOnlyOneHis(&wg, 3)
			}()
		}

		if i == 9 {
			go func() {
				dltEq1BackOnlyOneHisData = DltEqBackOnlyOneHis(&wg, 1)
			}()
		}

		if i == 10 {
			go func() {
				dltEq2BackOnlyOneHisData = DltEqBackOnlyOneHis(&wg, 2)
			}()
		}

		if i == 11 {
			go func() {
				dltEq3BackOnlyOneHisData = DltEqBackOnlyOneHis(&wg, 3)
			}()
		}

		if i == 12 {
			go func() {
				dltEq1BackHisData = DltEqBackHis(&wg, 1)
			}()
		}

		if i == 13 {
			go func() {
				dltEq2BackHisData = DltEqBackHis(&wg, 2)
			}()
		}

		if i == 14 {
			go func() {
				dltEq3BackHisData = DltEqBackHis(&wg, 3)
			}()
		}

		if i == 15 {
			go func() {
				dltEq1OeHisData = DltEqOeHis(&wg, 1)
			}()
		}

		if i == 16 {
			go func() {
				dltEq2OeHisData = DltEqOeHis(&wg, 2)
			}()
		}

		if i == 17 {
			go func() {
				dltEq3OeHisData = DltEqOeHis(&wg, 3)
			}()
		}

		if i == 18 {
			go func() {
				dltEq1HzHisData = DltEqHzHis(&wg, 1)
			}()
		}

		if i == 19 {
			go func() {
				dltEq2HzHisData = DltEqHzHis(&wg, 2)
			}()
		}

		if i == 20 {
			go func() {
				dltEq3HzHisData = DltEqHzHis(&wg, 3)
			}()
		}

		if i == 21 {
			go func() {
				dltEq1QzhHisData = DltEqQzhHis(&wg, 1)
			}()
		}

		if i == 22 {
			go func() {
				dltEq2QzhHisData = DltEqQzhHis(&wg, 2)
			}()
		}

		if i == 23 {
			go func() {
				dltEq3QzhHisData = DltEqQzhHis(&wg, 3)
			}()
		}
	}
	wg.Wait()
	//fmt.Printf("8 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 前区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "前区单号历史", dltFrontOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("9 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 关于设备号的历史都只考虑从11001期开始的历史数据
	// 设备1的前区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备1前区单号历史", dltEq1FrontOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备2的前区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备2前区单号历史", dltEq2FrontOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备3的前区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备3前区单号历史", dltEq3FrontOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 后区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "后区单号历史", dltBackOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("10 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 设备1的后区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备1后区单号历史", dltEq1BackOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备2的后区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备2后区单号历史", dltEq2BackOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备3的后区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备3后区单号历史", dltEq3BackOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}

	// 后区历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "后区历史", dltBackHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("11 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 设备1的后区历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备1后区历史", dltEq1BackHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备2的后区历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备2后区历史", dltEq2BackHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备3的后区历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备3后区历史", dltEq3BackHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}

	// 奇偶历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "奇偶历史", dltOeHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("12 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 设备1的奇偶历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备1奇偶历史", dltEq1OeHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备2的奇偶历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备2奇偶历史", dltEq2OeHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备3的奇偶历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备3奇偶历史", dltEq3OeHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 和值历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "和值历史", dltHzHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("13 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 设备1的和值历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备1和值历史", dltEq1HzHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备2的和值历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备2和值历史", dltEq2HzHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备3的和值历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备3和值历史", dltEq3HzHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}

	// 前中后历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "前中后历史", dltQzhHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("14 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 设备1的前中后历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备1前中后历史", dltEq1QzhHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备2的前中后历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备2前中后历史", dltEq2QzhHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备3的前中后历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "设备3前中后历史", dltEq3QzhHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}

	// 后区趋势1
	err = SaveBackQuShi1ToExcelFile(f, colorSlice, styleID, "后区趋势1", DltBackQuShi1())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("15 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 后区趋势2
	err = SaveBackQuShi2ToExcelFile(f, colorSlice, styleID, "后区趋势2", DltBackQuShi2())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("16 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 奇偶趋势
	err = SaveOeQuShiToExcelFile(f, colorSlice, styleID, "奇偶趋势", DltOeQuShi())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("17 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 和值趋势1
	err = SaveHzQuShi1ToExcelFile(f, colorSlice, styleID, "和值趋势1", DltHzQuShi1())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("18 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 和值趋势2
	err = SaveHzQuShi2ToExcelFile(f, colorSlice, styleID, "和值趋势2", DltHzQuShi2())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("19 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 前中后趋势
	err = SaveQzhQuShiToExcelFile(f, colorSlice, styleID, "前中后趋势", DltQzhQuShi())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("20 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 重号历史
	err = SaveCHongHaoToExcelFile(f, colorSlice, styleID, "重号历史", DltCHongHao())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	//fmt.Printf("21 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 设备1的重号历史
	err = SaveCHongHaoToExcelFile(f, colorSlice, styleID, "设备1重号历史", DltEqCHongHao(1))
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备2的重号历史
	err = SaveCHongHaoToExcelFile(f, colorSlice, styleID, "设备2重号历史", DltEqCHongHao(2))
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	// 设备3的重号历史
	err = SaveCHongHaoToExcelFile(f, colorSlice, styleID, "设备3重号历史", DltEqCHongHao(3))
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}

	tNum = 1
	for tNum <= len(typ2FrontHms) {
		typName := fmt.Sprintf("类型T%d", tNum)
		_ = SaveTypHmsToExcelFile(f, colorSlice, styleID, typName, typ2FrontHms[typName])
		tNum++
	}
	//fmt.Printf("22 %v\n", time.Now().Sub(startTime).Round(time.Second))
	typName := fmt.Sprintf("OtherT")
	_ = SaveTypHmsToExcelFile(f, colorSlice, styleID, typName, otherTxHms)
	//fmt.Printf("23 %v\n", time.Now().Sub(startTime).Round(time.Second))

	//runtime.ReadMemStats(&m2)
	//
	//fmt.Printf("Alloc 变化: %d MB\n", (m2.Alloc-m1.Alloc)/1024/1024)
	//fmt.Printf("TotalAlloc 变化: %d MB\n", (m2.TotalAlloc-m1.TotalAlloc)/1024/1024)
	//fmt.Printf("Sys 变化: %d MB\n", (m2.Sys-m1.Sys)/1024/1024)
	//fmt.Printf("HeapAlloc 变化: %d MB\n", (m2.HeapAlloc-m1.HeapAlloc)/1024/1024)
}

func DltFrontStatDataToExcel(crossDrawNumSli []int, flag string, eqNumCount int) {
	if len(DxDlts) == 0 {
		InitDlts()
	}

	saveDir := "./output/excel"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建目录失败: %v", err))
		return
	}
	var fileName string
	if eqNumCount == 0 {
		fileName = filepath.Join(saveDir, fmt.Sprintf("大乐透截止至%s的数据分析_前区跨期数统计_%s", DxDlts[0].DrawTime, flag))
	} else {
		fileName = filepath.Join(saveDir, fmt.Sprintf("设备%d_大乐透截止至%s的数据分析_前区跨期数统计_%s", eqNumCount, DxDlts[0].DrawTime, flag))
	}

	// TODO 判断文件是否已经存在, 若存在则删除后

	f, err := excel.CreateNewExcelFile(fileName)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建工作簿过程中出现错误：%v\n", err))
		return
	}

	defer func() {
		if err = f.Close(); err != nil {
			lg.ErrorToFile(fmt.Sprintf("关闭excel文件出现错误：%v", err))
		}
	}()

	defer func() {
		_ = f.Save()
	}()

	styleID, _ := f.NewStyle(&excelize.Style{
		//Font: &excelize.Font{Family: "Consolas", Size: 8},
		Font: &excelize.Font{Family: "monospace", Size: 8},
		Alignment: &excelize.Alignment{
			WrapText: false,
		},
	})

	greenStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#00FF00"},
			Pattern: 1,
		},
	})
	_ = greenStyle

	yellowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFF00"},
			Pattern: 1,
		},
	})
	_ = yellowStyle
	colorSlice := []int{greenStyle, yellowStyle}

	var statData []DltTSN
	for i, crossDrawNum := range crossDrawNumSli {
		if eqNumCount == 0 {
			statData = DltFrontCrossDrawNumStat(crossDrawNum)
		} else {
			statData = DltEqFrontCrossDrawNumStat(crossDrawNum, eqNumCount)
		}

		// 写入 excel 表格
		_ = SaveFrontStatToExcelFile(f, colorSlice, styleID, fmt.Sprintf("跨%d期", crossDrawNum), statData, i == 0)
	}
}

func DltNextFrontStatDataToExcel(crossDrawNumSli []int) {
	if len(DxDlts) == 0 {
		InitDlts()
	}

	saveDir := "./output/excel"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建目录失败: %v", err))
		return
	}
	fileName := filepath.Join(saveDir, fmt.Sprintf("大乐透截止至%s的数据分析_%s的下一期前区分析", DxDlts[0].DrawNum, DxDlts[0].DrawTime))

	f, err := excel.CreateNewExcelFile(fileName)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建工作簿过程中出现错误：%v\n", err))
		return
	}

	defer func() {
		if err = f.Close(); err != nil {
			lg.ErrorToFile(fmt.Sprintf("关闭excel文件出现错误：%v", err))
		}
	}()

	defer func() {
		_ = f.Save()
	}()

	styleID, _ := f.NewStyle(&excelize.Style{
		//Font: &excelize.Font{Family: "Consolas", Size: 8},
		Font: &excelize.Font{Family: "monospace", Size: 8},
		Alignment: &excelize.Alignment{
			WrapText: false,
		},
	})

	greenStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#00FF00"},
			Pattern: 1,
		},
	})
	_ = greenStyle

	yellowStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFF00"},
			Pattern: 1,
		},
	})
	_ = yellowStyle
	colorSlice := []int{greenStyle, yellowStyle}

	var coverStatData []DltCoverCrossTSN
	var statData []DltCrossTSN
	for i := 0; i <= 9 && i < len(DxDlts); i++ {
		if i == 0 {
			statData = DltFrontSpDrawNumCrossDrawNumStat(DxDlts[i].DrawNum, crossDrawNumSli, true, false)
			coverStatData = append(coverStatData, DltCoverCrossTSN{
				DrawNum:  DxDlts[i].DrawNum,
				CrossTsn: statData,
			})
			statData = DltFrontSpDrawNumCrossDrawNumStat(DxDlts[i].DrawNum, crossDrawNumSli, false, true)
		} else {
			statData = DltFrontSpDrawNumCrossDrawNumStat(DxDlts[i].DrawNum, crossDrawNumSli, false, true)
		}
		coverStatData = append(coverStatData, DltCoverCrossTSN{
			DrawNum:  DxDlts[i].DrawNum,
			CrossTsn: statData,
		})
	}

	// 写入 excel 表格
	_ = SaveNextFrontStatToExcelFile(f, colorSlice, styleID, fmt.Sprintf("%s的下一期", DxDlts[0].DrawNum), coverStatData, true)

}

func SaveNextFrontStatToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, dltCoverCrossTsnData []DltCoverCrossTSN, isFirst bool) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	_ = greenStyle
	//yellowStyle := colorSlice[1]
	// 创建自动换行样式
	wrapStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
			Vertical: "top",
		},
	})

	if isFirst {
		if err = f.SetSheetName("Sheet1", sheetName); err != nil {
			return fmt.Errorf("重命名Sheet1工作表为%s遇到错误：%v", sheetName, err)
		}
	} else {
		_, err = f.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
		}
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	for ik, dltCoverCrossTSN := range dltCoverCrossTsnData {
		if ik == 0 {
			tNum := 1
			colNum0 := 1 // 这里从1开始,因为 excel.GetColumnStr 会自动减1
			colZiMu := excel.GetColumnStr(colNum0)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "序号")
			_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)

			colNum0++
			colZiMu = excel.GetColumnStr(colNum0)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "跨几期")
			_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)

			colNum0++
			colZiMu = excel.GetColumnStr(colNum0)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "设备号")
			_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 9)

			colNum0++
			colZiMu = excel.GetColumnStr(colNum0)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "下一期待推测")
			_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 30)

			//colNum0++
			//for i := colNum0; i < 30; i++ {
			//	ziMu := excel.GetColumnStr(i)
			//	_ = f.SetColWidth(sheetName, ziMu, ziMu, 60)
			//}
			rowNum := 2
			_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)

			for _, tsn := range dltCoverCrossTSN.CrossTsn {
				_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), tsn.Cross)

				if tsn.EqCount == 0 {
					_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), "全")
				} else {
					_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), tsn.EqCount)
				}
				text := ""

				for _, d := range tsn.Tsn {
					if text == "" {
						text = fmt.Sprintf("%d*%v", d.Times, d.Sli)
					} else {
						text += fmt.Sprintf("\n%d*%v", d.Times, d.Sli)
					}
				}
				cellPos := fmt.Sprintf("D%d", rowNum)
				_ = f.SetCellValue(sheetName, cellPos, text)
				_ = f.SetCellStyle(sheetName, cellPos, cellPos, wrapStyle)
				tNum++
				rowNum++
			}
		} else {
			var ziMu string
			ziMu = excel.GetColumnStr(4 + (ik-1)*2 + 1)

			_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", ziMu), "设备号")
			_ = f.SetColWidth(sheetName, ziMu, ziMu, 9)

			ziMu = excel.GetColumnStr(4 + (ik-1)*2 + 2)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", ziMu), dltCoverCrossTSN.DrawNum)
			_ = f.SetColWidth(sheetName, ziMu, ziMu, 30)
			tNum := 1
			rowNum := 2
			_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)

			for _, tsn := range dltCoverCrossTSN.CrossTsn {
				ziMu = excel.GetColumnStr(4 + (ik-1)*2 + 1)

				if tsn.EqCount == 0 {
					_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", ziMu, rowNum), "全")
				} else {
					_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", ziMu, rowNum), tsn.EqCount)
				}
				ziMu = excel.GetColumnStr(4 + (ik-1)*2 + 2)
				var runs []excelize.RichTextRun
				runs = nil
				for _, d := range tsn.Tsn {
					if len(runs) == 0 {
						if len(d.NumStr) == 0 {
							runs = append(runs, excelize.RichTextRun{
								Text: fmt.Sprintf("%d*%v=>%v", d.Times, d.Sli, d.NumStr),
							})
						} else {
							runs = append(runs, excelize.RichTextRun{
								Text: fmt.Sprintf("%d*%v=>%v", d.Times, d.Sli, d.NumStr),
								Font: &excelize.Font{
									//Color: "00B050", // 设置为绿色 (十六进制 RRGGBB)
									Color: "FF0000", // 设置为红色 (十六进制 RRGGBB)
								},
							})
						}
					} else {
						runs = append(runs, excelize.RichTextRun{
							Text: "\n",
						})
						if len(d.NumStr) == 0 {
							runs = append(runs, excelize.RichTextRun{
								Text: fmt.Sprintf("%d*%v=>%v", d.Times, d.Sli, d.NumStr),
							})
						} else {
							runs = append(runs, excelize.RichTextRun{
								Text: fmt.Sprintf("%d*%v=>%v", d.Times, d.Sli, d.NumStr),
								Font: &excelize.Font{
									//Color: "00B050", // 设置为绿色 (十六进制 RRGGBB)
									Color: "FF0000", // 设置为红色 (十六进制 RRGGBB)
								},
							})
						}
					}
				}

				cellPos := fmt.Sprintf("%s%d", ziMu, rowNum)
				_ = f.SetCellRichText(sheetName, cellPos, runs)
				_ = f.SetCellStyle(sheetName, cellPos, cellPos, wrapStyle)

				tNum++
				rowNum++
			}
		}

	}

	_ = f.AutoFilter(sheetName, "A1:AZ1", nil)
	return
}

func SaveFrontStatToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, dltTsnData []DltTSN, isFirst bool) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	_ = greenStyle
	//yellowStyle := colorSlice[1]

	if isFirst {
		if err = f.SetSheetName("Sheet1", sheetName); err != nil {
			return fmt.Errorf("重命名Sheet1工作表为%s遇到错误：%v", sheetName, err)
		}
	} else {
		_, err = f.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
		}
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	colNum0 := 1 // 这里从1开始,因为 excel.GetColumnStr 会自动减1
	colZiMu := excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "序号")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "期号")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "开奖日期")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 9)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "顺序号码")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 18)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "来自段")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "共几段")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)

	colNum0++
	for i := colNum0; i < 30; i++ {
		ziMu := excel.GetColumnStr(i)
		_ = f.SetColWidth(sheetName, ziMu, ziMu, 60)
	}

	tNum := 1
	rowNum := 2
	colNum1 := 1
	for _, tsn := range dltTsnData {
		colNum1 = 1
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), tsn.Dlt.DrawNum)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), tsn.Dlt.DrawTime)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), fmt.Sprintf("%s,%s,%s,%s,%s|%s,%s", tsn.Dlt.F1, tsn.Dlt.F2, tsn.Dlt.F3, tsn.Dlt.F4, tsn.Dlt.F5, tsn.Dlt.B1, tsn.Dlt.B2))
		colNum1++
		// 来自段
		lzd := 0
		// 共几段
		gjd := len(tsn.Tsn)
		for _, d := range tsn.Tsn {
			if len(d.NumStr) > 0 {
				lzd++
			}
		}
		_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), lzd)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), gjd)

		i := 7
		//// []中间的号码
		//zjHm := ""
		//// => 后面的号码
		//hmHm := ""
		for _, d := range tsn.Tsn {
			ziMu := excel.GetColumnStr(i)
			colNum1++
			_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", ziMu, rowNum), fmt.Sprintf("%d*%v=>%v", d.Times, d.Sli, d.NumStr))

			if len(d.NumStr) > 0 {
				cellPos := fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum)
				_ = f.SetCellStyle(sheetName, cellPos, cellPos, greenStyle)
			}

			i++
		}

		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, "A1:AZ1", nil)
	return
}

func SetFirstSheetContent(f *excelize.File, colorSlice []int, styleID int, sheetName string, data []DltExcelData) (err error) {
	defer func() {
		_ = f.Save()
	}()
	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]
	//greenStyle, _ := f.NewStyle(&excelize.Style{
	//	Fill: excelize.Fill{
	//		Type:    "pattern",
	//		Color:   []string{"#00FF00"},
	//		Pattern: 1,
	//	},
	//})

	if err = f.SetSheetName("Sheet1", sheetName); err != nil {
		return fmt.Errorf("重命名Sheet1工作表为%s遇到错误：%v", sheetName, err)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		XSplit: 3,    // 冻结前三列
		YSplit: 1,    // 冻结前1行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	t := reflect.TypeOf(data[0])

	// 如果是指针，先取 Elem
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	maxColNum := t.NumField() // t.NumField() + 1 加一列用于序号, 但这个结构体中已经有XuHao这个字段了,故不加1
	for colNum := 1; colNum <= maxColNum; colNum++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%s", excel.GetColumnStr(colNum)), styleID)
	}

	colNum0 := 1 // 这里从1开始,因为 excel.GetColumnStr 会自动减1
	colZiMu := excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "序号")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)

	// 1. 获取单元格当前样式的 ID
	oldStyleID, _ := f.GetCellStyle(sheetName, fmt.Sprintf("%s1", colZiMu))

	// 2. 通过样式 ID 获取完整的样式定义
	oldStyle, _ := f.GetStyle(oldStyleID)

	// 3. 修改数字格式（仅设置这一个属性）
	//oldStyle.NumFmt = 2                         // 内置数字格式，2 表示保留两位小数（带千位分隔符）
	oldStyle.NumFmt = 49                         // 使用自定义数字格式
	oldStyle.CustomNumFmt = &[]string{"0.00"}[0] // 设置格式为两位小数

	// 4. 根据修改后的样式定义，创建新的样式 ID
	fStyle, _ := f.NewStyle(oldStyle)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "期号")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "开奖日期")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 9)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "号码")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 18)

	// ---
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "顺序")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 18)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "奖池奖金")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 12)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "销售额")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 12)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "一基注")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "一基元")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "一追注")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "一追元")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "二基注")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "二基元")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "二追注")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "二追元")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "设备")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "奇偶")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "和值")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "前中后")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "4重")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "5重")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "6重")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "7重")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "总重")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "累重")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 6)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "类型T1")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	colWidth := float64(5)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T2")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T3")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T4")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T5")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T6")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T7")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T8")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T9")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T10")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T11")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T12")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T13")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T14")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T15")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T16")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "T17")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, colWidth)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "OtherT")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	maxColZiMu := colZiMu

	rowNum := 2
	colNum1 := 1
	for _, iData := range data {
		colNum1 = 1
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.XuHao)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.DrawNum)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.DrawTime)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.FullHm)
		// --
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.UnSortDrawResult)

		colNum1++
		_ = f.SetCellFloat(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.PoolBalance, 2, 64)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), fStyle)

		colNum1++
		_ = f.SetCellFloat(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.TotalSaleAmount, 2, 64)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), fStyle)

		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.StakeCount101)

		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.StakeAmount101)

		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.StakeCount201)

		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.StakeAmount201)

		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.StakeCount301)

		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.StakeAmount301)

		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.StakeCount401)

		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.StakeAmount401)
		//---

		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.EquipmentCount)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.Oe)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.Hz)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.Qzh)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.NewAddCh4)
		if iData.NewAddCh4 > 0 {
			cellPos := fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum)
			_ = f.SetCellStyle(sheetName, cellPos, cellPos, greenStyle)
		}

		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.NewAddCh5)
		if iData.NewAddCh5 > 0 {
			cellPos := fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum)
			_ = f.SetCellStyle(sheetName, cellPos, cellPos, greenStyle)
		}
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.NewAddCh6)
		if iData.NewAddCh6 > 0 {
			cellPos := fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum)
			_ = f.SetCellStyle(sheetName, cellPos, cellPos, greenStyle)
		}
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.NewAddCh7)
		if iData.NewAddCh7 > 0 {
			cellPos := fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum)
			_ = f.SetCellStyle(sheetName, cellPos, cellPos, greenStyle)
		}
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.DangQiTotalNewAddCh)
		if iData.DangQiTotalNewAddCh > 0 {
			cellPos := fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum)
			_ = f.SetCellStyle(sheetName, cellPos, cellPos, greenStyle)
		}
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.LeiJiaCh)

		tds := []int{
			iData.T71, iData.T72, iData.T73, iData.T74, iData.T75,
			iData.T111, iData.T112, iData.T113, iData.T114, iData.T115,
			iData.T151, iData.T152, iData.T153, iData.T154, iData.T155,
			iData.T191, iData.T192,
		}

		for iTd := 0; iTd <= len(tds)-1; iTd++ {
			colNum1++
			tempColZiMu := excel.GetColumnStr(colNum1)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", tempColZiMu, rowNum), tds[iTd])
			if tds[iTd] > 0 {
				_ = f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", tempColZiMu, rowNum), fmt.Sprintf("%s%d", tempColZiMu, rowNum), greenStyle)
			}
		}
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.OtherT)
		if iData.OtherT > 0 {
			cellPos := fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum)
			_ = f.SetCellStyle(sheetName, cellPos, cellPos, greenStyle)
		}

		rowNum++
	}

	_ = f.AutoFilter(sheetName, fmt.Sprintf("A1:%s1", maxColZiMu), nil)
	return
}

func SaveTypInfoToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, typ2ae map[string]models.DltMoni) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	_ = f.SetCellValue(sheetName, "A1", "序号")
	_ = f.SetCellValue(sheetName, "B1", "类型名称")
	_ = f.SetCellValue(sheetName, "C1", "所属类型")
	_ = f.SetCellValue(sheetName, "D1", "命中期数")
	_ = f.SetCellValue(sheetName, "E1", "前区号码组合数")
	_ = f.SetCellValue(sheetName, "F1", "概率提升倍数")

	_ = f.SetColWidth(sheetName, "A", "A", 5)
	_ = f.SetColWidth(sheetName, "B", "B", 10)
	_ = f.SetColWidth(sheetName, "C", "C", 10)
	_ = f.SetColWidth(sheetName, "D", "D", 10)
	_ = f.SetColWidth(sheetName, "E", "E", 15)
	_ = f.SetColWidth(sheetName, "F", "F", 12)

	_ = f.SetCellValue(sheetName, "G1", "A")
	_ = f.SetColWidth(sheetName, "G", "G", 58)

	_ = f.SetCellValue(sheetName, "H1", "B")
	_ = f.SetCellValue(sheetName, "I1", "C")
	_ = f.SetCellValue(sheetName, "J1", "D")
	_ = f.SetCellValue(sheetName, "K1", "E")

	for ziMu := 'H'; ziMu <= 'K'; ziMu++ {
		_ = f.SetColWidth(sheetName, fmt.Sprintf("%c", ziMu), fmt.Sprintf("%c", ziMu), 20)
	}

	tNum := 1
	rowNum := 2
	for tNum <= len(typ2ae) {
		if tNum == len(typ2ae) {
			iData := typ2ae["OtherT"]
			_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), "OtherT")
			_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), "OtherT")
			_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), iData.Cs)
			_ = f.SetCellStyle(sheetName, fmt.Sprintf("D%d", rowNum), fmt.Sprintf("D%d", rowNum), greenStyle)
			tempNum := iData.Comb
			_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), fmt.Sprintf("%d", tempNum))
			_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), fmt.Sprintf("%.2f", float64(324632)/(float64(tempNum))))
		} else {
			typName := fmt.Sprintf("类型T%d", tNum)
			iData := typ2ae[typName]
			_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), typName)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), iData.Typ)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), iData.Cs)
			_ = f.SetCellStyle(sheetName, fmt.Sprintf("D%d", rowNum), fmt.Sprintf("D%d", rowNum), greenStyle)
			tempNum := 1
			switch iData.Typ {
			case "77777":
				tempNum = 7 * 7 * 7 * 7 * 7
			case "116666":
				tempNum = 11 * 6 * 6 * 6 * 6
			case "155555":
				tempNum = 15 * 5 * 5 * 5 * 5
			case "194444":
				tempNum = 19 * 4 * 4 * 4 * 4
			}

			_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), fmt.Sprintf("%d", tempNum))
			_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), fmt.Sprintf("%.2f", float64(324632)/(float64(tempNum))))

			_ = f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), iData.A)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), iData.B)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), iData.C)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), iData.D)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), iData.E)
		}

		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, "A1:K1", nil)
	return
}

func SaveBestCombToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, bestCombs []*BestCombSt) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	t := reflect.TypeOf(bestCombs[0])

	// 如果是指针，先取 Elem
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	maxColNum := t.NumField() // t.NumField() + 1 加一列用于序号, 但这个结构体中已经有FrontHms 这个字段了且这个字段里没用,故刚好用于序号,因此不加1

	for colNum := 1; colNum <= maxColNum; colNum++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%s", excel.GetColumnStr(colNum)), styleID)
	}

	colNum0 := 1 // 这里从1开始,因为 excel.GetColumnStr 会自动减1
	colZiMu := excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "序号")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 5)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "K列")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 5)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "组合")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 58)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "历史覆盖期数")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 12)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "前区号码组合个数")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 20)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "概率提升倍数")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 16)
	maxColZiMu := colZiMu

	tNum := 1
	rowNum := 2
	colNum1 := 1
	for _, bestComb := range bestCombs {
		colNum1 = 1
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), tNum)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), bestComb.ColCount)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), strings.Join(bestComb.Tx, ","))
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), bestComb.HisFuGaiDrawNum)
		cellPos := fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum)
		_ = f.SetCellStyle(sheetName, cellPos, cellPos, greenStyle)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), bestComb.FrontHmCount)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), bestComb.GaiLvTs)

		tNum++
		rowNum++
	}

	_ = f.AutoFilter(sheetName, fmt.Sprintf("A1:%s1", maxColZiMu), nil)

	return
}

func SaveTypHmsToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, frontHms []string) (err error) {
	defer func() {
		_ = f.Save()
	}()

	//greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	colNum0 := 1 // 这里从1开始,因为 excel.GetColumnStr 会自动减1
	colZiMu := excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "号码")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 5)
	//colNum0++
	//colZiMu = excel.GetColumnStr(colNum0)
	//_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "奇偶")
	//_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 5)
	//colNum0++
	//colZiMu = excel.GetColumnStr(colNum0)
	//_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "和值")
	//_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 58)
	//colNum0++
	//colZiMu = excel.GetColumnStr(colNum0)
	//_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "新4重号")
	//_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 12)
	//maxColZiMu := colZiMu

	tNum := 1
	rowNum := 2
	colNum1 := 1
	for _, frontHm := range frontHms {
		colNum1 = 1
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), frontHm)

		tNum++
		rowNum++
	}

	//_ = f.AutoFilter(sheetName, fmt.Sprintf("A1:%s1", maxColZiMu), nil)

	return
}

func SaveHisToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, dltHisSlice []DltHis) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	_ = f.SetCellValue(sheetName, "A1", "序号")
	_ = f.SetCellValue(sheetName, "B1", "类型")
	_ = f.SetCellValue(sheetName, "C1", "出现期数")
	_ = f.SetCellValue(sheetName, "D1", "所有可能组合数")
	_ = f.SetColWidth(sheetName, "A", "A", 5)
	_ = f.SetColWidth(sheetName, "B", "B", 8)
	_ = f.SetColWidth(sheetName, "C", "C", 10)
	_ = f.SetColWidth(sheetName, "D", "D", 18)

	_ = f.SetCellValue(sheetName, "E1", "最新10期")
	_ = f.SetColWidth(sheetName, "E", "E", 10)
	_ = f.SetCellValue(sheetName, "F1", "20")
	_ = f.SetCellValue(sheetName, "G1", "30")
	_ = f.SetCellValue(sheetName, "H1", "50")
	_ = f.SetCellValue(sheetName, "I1", "100")
	_ = f.SetCellValue(sheetName, "J1", "200")
	_ = f.SetCellValue(sheetName, "K1", "500")
	_ = f.SetCellValue(sheetName, "L1", "1000")
	_ = f.SetCellValue(sheetName, "M1", "1500")
	_ = f.SetCellValue(sheetName, "N1", "2000")
	_ = f.SetCellValue(sheetName, "O1", "2500")
	_ = f.SetCellValue(sheetName, "P1", "3500")

	for ziMu := 'F'; ziMu <= 'P'; ziMu++ {
		_ = f.SetColWidth(sheetName, fmt.Sprintf("%c", ziMu), fmt.Sprintf("%c", ziMu), 8)
	}

	tNum := 1
	rowNum := 2
	for _, dltHis := range dltHisSlice {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), dltHis.Typ)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), dltHis.Cs)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("C%d", rowNum), greenStyle)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), dltHis.AllCount)

		_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), dltHis.Last10)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), dltHis.Last20)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), dltHis.Last30)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), dltHis.Last50)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), dltHis.Last100)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), dltHis.Last200)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), dltHis.Last500)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), dltHis.Last1000)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), dltHis.Last1500)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("N%d", rowNum), dltHis.Last2000)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("O%d", rowNum), dltHis.Last2500)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("P%d", rowNum), dltHis.Last3500)

		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, "A1:P1", nil)
	return
}

func SaveBackQuShi1ToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, quShiSlice []KeyWithLength) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	_ = f.SetCellValue(sheetName, "A1", "序号")
	_ = f.SetCellValue(sheetName, "B1", "类型")
	_ = f.SetCellValue(sheetName, "C1", "出现期数")

	_ = f.SetColWidth(sheetName, "A", "A", 5)
	_ = f.SetColWidth(sheetName, "B", "B", 12)
	_ = f.SetColWidth(sheetName, "C", "C", 10)

	tNum := 1
	rowNum := 2
	for _, quShi := range quShiSlice {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), quShi.Key)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), quShi.Length)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("C%d", rowNum), greenStyle)

		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, "A1:C1", nil)
	return
}

func SaveBackQuShi2ToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, quShiSlice []DltBackChQuShi) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	_ = f.SetCellValue(sheetName, "A1", "序号")
	_ = f.SetCellValue(sheetName, "B1", "类型")
	_ = f.SetCellValue(sheetName, "C1", "出现期数")
	_ = f.SetCellValue(sheetName, "D1", "所组合数")
	_ = f.SetCellValue(sheetName, "E1", "已组合数")
	_ = f.SetCellValue(sheetName, "F1", "未组合数")
	_ = f.SetCellValue(sheetName, "G1", "所有可能趋势组合")
	_ = f.SetCellValue(sheetName, "H1", "已经出现趋势组合")
	_ = f.SetCellValue(sheetName, "I1", "未出现趋势组合")

	_ = f.SetColWidth(sheetName, "A", "A", 5)
	_ = f.SetColWidth(sheetName, "B", "B", 12)
	_ = f.SetColWidth(sheetName, "C", "C", 10)
	_ = f.SetColWidth(sheetName, "D", "D", 10)
	_ = f.SetColWidth(sheetName, "E", "E", 10)
	_ = f.SetColWidth(sheetName, "F", "F", 10)
	_ = f.SetColWidth(sheetName, "G", "G", 242)
	_ = f.SetColWidth(sheetName, "H", "H", 200)
	_ = f.SetColWidth(sheetName, "I", "I", 200)

	tNum := 1
	rowNum := 2
	for _, qs := range quShiSlice {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), fmt.Sprintf("%s->%s", qs.BackComb, qs.Qs))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), qs.Cs)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("C%d", rowNum), greenStyle)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), len(qs.AllCombs))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), len(qs.HadExistCombs))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), len(qs.HadNotExistCombs))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), strings.Join(qs.AllCombs, "|"))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), strings.Join(qs.HadExistCombs, "|"))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), strings.Join(qs.HadNotExistCombs, "|"))
		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, "A1:F1", nil)
	return
}

func SaveOeQuShiToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, quShiSlice []KeyWithLength) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	_ = f.SetCellValue(sheetName, "A1", "序号")
	_ = f.SetCellValue(sheetName, "B1", "类型")
	_ = f.SetCellValue(sheetName, "C1", "出现期数")

	_ = f.SetColWidth(sheetName, "A", "A", 5)
	_ = f.SetColWidth(sheetName, "B", "B", 12)
	_ = f.SetColWidth(sheetName, "C", "C", 10)

	tNum := 1
	rowNum := 2
	for _, quShi := range quShiSlice {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), quShi.Key)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), quShi.Length)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("C%d", rowNum), greenStyle)

		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, "A1:C1", nil)
	return
}

func SaveHzQuShi1ToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, quShiSlice []KeyWithLength) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	_ = f.SetCellValue(sheetName, "A1", "序号")
	_ = f.SetCellValue(sheetName, "B1", "类型")
	_ = f.SetCellValue(sheetName, "C1", "出现期数")

	_ = f.SetColWidth(sheetName, "A", "A", 5)
	_ = f.SetColWidth(sheetName, "B", "B", 12)
	_ = f.SetColWidth(sheetName, "C", "C", 10)

	tNum := 1
	rowNum := 2
	for _, quShi := range quShiSlice {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), quShi.Key)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), quShi.Length)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("C%d", rowNum), greenStyle)

		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, "A1:C1", nil)
	return
}

func SaveHzQuShi2ToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, quShiSlice []DltHzChQuShi) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	_ = f.SetCellValue(sheetName, "A1", "序号")
	_ = f.SetCellValue(sheetName, "B1", "类型")
	_ = f.SetCellValue(sheetName, "C1", "类型解释")
	_ = f.SetCellValue(sheetName, "D1", "出现期数")
	_ = f.SetCellValue(sheetName, "E1", "已经出现和值趋势")

	_ = f.SetColWidth(sheetName, "A", "A", 5)
	_ = f.SetColWidth(sheetName, "B", "B", 12)
	_ = f.SetColWidth(sheetName, "C", "C", 20)
	_ = f.SetColWidth(sheetName, "D", "D", 10)
	_ = f.SetColWidth(sheetName, "E", "E", 255)

	tNum := 1
	rowNum := 2
	for _, qs := range quShiSlice {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), fmt.Sprintf("%s->%s", qs.HzCh, qs.NextHzCh))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("[%d~%d]到[%d~%d]", HzABCDEMs[qs.HzCh].Min, HzABCDEMs[qs.HzCh].Max, HzABCDEMs[qs.NextHzCh].Min, HzABCDEMs[qs.NextHzCh].Max))
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), qs.Cs)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("D%d", rowNum), fmt.Sprintf("D%d", rowNum), greenStyle)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), strings.Join(qs.Hz2NextHz, "|"))
		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, "A1:D1", nil)
	return
}

func SaveQzhQuShiToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, qzhQsSlice []KeyWithLength) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	colNum0 := 1 // 这里从1开始,因为 excel.GetColumnStr 会自动减1
	colZiMu := excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "序号")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "类型")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 10)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "出现期数")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 10)

	maxColNum := colNum0
	maxColZiMu := colZiMu

	for colNum := 1; colNum <= maxColNum; colNum++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%s", excel.GetColumnStr(colNum)), styleID)
	}

	tNum := 1
	rowNum := 2
	colNum1 := 1
	for _, iData := range qzhQsSlice {
		colNum1 = 1
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), tNum)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.Key)
		colNum1++
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), iData.Length)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), fmt.Sprintf("%s%d", excel.GetColumnStr(colNum1), rowNum), greenStyle)

		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, fmt.Sprintf("A1:%s1", maxColZiMu), nil)
	return
}

func SaveCHongHaoToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, chSlice []DltChSt) (err error) {
	defer func() {
		_ = f.Save()
	}()

	greenStyle := colorSlice[0]
	//yellowStyle := colorSlice[1]

	_, err = f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("新建工作表%s出现错误：%v\n", sheetName, err)
	}

	for ch := 'A'; ch <= 'Z'; ch++ {
		_ = f.SetColStyle(sheetName, fmt.Sprintf("%c", ch), styleID)
	}

	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze: true, // 启用冻结窗口
		YSplit: 1,    // 冻结第一行
		//TopLeftCell: "A2",         // 冻结后左上角的单元格
		//ActivePane:  "bottomLeft", // 冻结后活动区域
	})

	_ = f.SetCellValue(sheetName, "A1", "序号")
	_ = f.SetCellValue(sheetName, "B1", "类型")
	_ = f.SetCellValue(sheetName, "C1", "类型1")
	_ = f.SetCellValue(sheetName, "D1", "类型2")
	_ = f.SetCellValue(sheetName, "E1", "出现期数")
	tempNum := 1
	for i := 6; i < 600; i++ {
		ziMu := excel.GetColumnStr(i)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", ziMu), fmt.Sprintf("历史%d", tempNum))
		tempNum++
	}

	_ = f.SetColWidth(sheetName, "A", "A", 5)
	_ = f.SetColWidth(sheetName, "B", "B", 18)
	_ = f.SetColWidth(sheetName, "C", "C", 10)
	_ = f.SetColWidth(sheetName, "D", "D", 10)
	_ = f.SetColWidth(sheetName, "E", "E", 10)

	for i := 6; i < 600; i++ {
		ziMu := excel.GetColumnStr(i)
		_ = f.SetColWidth(sheetName, ziMu, ziMu, 25)
	}

	tNum := 1
	rowNum := 2
	for _, ch := range chSlice {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), ch.Typ)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), ch.Typ1)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), ch.Typ2)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), ch.Cs)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("E%d", rowNum), fmt.Sprintf("E%d", rowNum), greenStyle)
		i := 6
		for _, dltInfo := range ch.DltInfos {
			ziMu := excel.GetColumnStr(i)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", ziMu, rowNum), dltInfo)
			i++
		}

		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, "A1:E1", nil)
	return
}
