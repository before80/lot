package ana_ssq

import (
	"fmt"
	"math"
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

func SsqDataToExcel(prevRunMoni bool) {
	startTime := time.Now()
	if prevRunMoni {
		lg.InfoToFileAndStdOut(fmt.Sprintf("运行模拟数据中,请耐心等待...\n"))
		BatchMoni(true, 3)
		lg.InfoToFileAndStdOut(fmt.Sprintf("运行模拟数据所需要的时间: %v\n", time.Now().Sub(startTime).Round(time.Second)))
	}

	// 找到类型为666663匹配历史开奖号码最多的前5种
	var moni666663s []models.SsqMoni
	db.DB.Where(map[string]interface{}{"typ": "666663"}).Order("cs desc,id asc").Limit(5).Find(&moni666663s)
	// 找到类型为855555匹配历史开奖号码最多的前5种
	var moni855555s []models.SsqMoni
	db.DB.Where(map[string]interface{}{"typ": "855555"}).Order("cs desc,id asc").Limit(5).Find(&moni855555s)
	// 找到类型为1344444匹配历史开奖号码最多的前5种
	var moni1344444s []models.SsqMoni
	db.DB.Where(map[string]interface{}{"typ": "1344444"}).Order("cs desc,id asc").Limit(5).Find(&moni1344444s)
	// 找到类型为1833333匹配历史开奖号码最多的前2种
	var moni1833333s []models.SsqMoni
	db.DB.Where(map[string]interface{}{"typ": "1833333"}).Order("cs desc,id asc").Limit(2).Find(&moni1833333s)

	hm6s := BuildFrontHmSlices(moni666663s, 5)
	hm8s := BuildFrontHmSlices(moni855555s, 5)
	hm13s := BuildFrontHmSlices(moni1344444s, 5)
	hm18s := BuildFrontHmSlices(moni1833333s, 2)
	otherTxHms := CalSsqNotTxFrontHmSlice(hm6s, hm8s, hm13s, hm18s)

	fmt.Printf("1 %v\n", time.Now().Sub(startTime).Round(time.Second))
	InitSsqs()
	fmt.Printf("2 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 计算重号累加
	drawNum2CHongHaoSt := CHongHaoLeiJia()
	ssqs := DxSsqs
	lenSsqs := len(ssqs)
	var ssqExcelDatas []SsqExcelData
	tr6 := make([][]int, 5)
	tr8 := make([][]int, 5)
	tr13 := make([][]int, 5)
	tr18 := make([][]int, 2)
	otherTCs := 0

	for xuHao, ssq := range ssqs {
		rowNum := xuHao + 1
		frontHm := fmt.Sprintf("%s,%s,%s,%s,%s,%s", ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6)
		hz := CalSsqHz([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6, ssq.B1})
		oe := CalSsqOe([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6, ssq.B1})
		t6 := make([]int, 5)
		for i, hm := range hm6s {
			if i >= len(t6) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t6[i] = 1
				tr6[i] = append(tr6[i], rowNum)
			}
		}

		t8 := make([]int, 5)
		for i, hm := range hm8s {
			if i >= len(t8) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t8[i] = 1
				tr8[i] = append(tr8[i], rowNum)
			}
		}

		t13 := make([]int, 5)
		for i, hm := range hm13s {
			if i >= len(t13) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t13[i] = 1
				tr13[i] = append(tr13[i], rowNum)
			}
		}

		t18 := make([]int, 2)
		for i, hm := range hm18s {
			if i >= len(t18) {
				break
			}
			if slices.Contains(hm, frontHm) {
				t18[i] = 1
				tr18[i] = append(tr18[i], rowNum)
			}
		}

		otherT := 0
		if slices.Contains(otherTxHms, frontHm) {
			otherT = 1
			otherTCs++
		}

		ssqExcelDatas = append(ssqExcelDatas, SsqExcelData{
			XuHao:               lenSsqs - xuHao,
			DrawNum:             ssq.DrawNum,
			DrawTime:            ssq.DrawTime,
			FrontHm:             frontHm,
			FullHm:              fmt.Sprintf("%s|%s", frontHm, ssq.B1),
			Hz:                  hz,
			Oe:                  oe,
			Qzh:                 ssq.Qzh,
			NewAddCh4:           drawNum2CHongHaoSt[ssq.DrawNum].NewAddCh4,
			NewAddCh5:           drawNum2CHongHaoSt[ssq.DrawNum].NewAddCh5,
			NewAddCh6:           drawNum2CHongHaoSt[ssq.DrawNum].NewAddCh6,
			NewAddCh7:           drawNum2CHongHaoSt[ssq.DrawNum].NewAddCh7,
			DangQiTotalNewAddCh: drawNum2CHongHaoSt[ssq.DrawNum].DangQiTotalNewAddCh,
			LeiJiaCh:            drawNum2CHongHaoSt[ssq.DrawNum].LeiJiaCh,
			T61:                 t6[0], T62: t6[1], T63: t6[2], T64: t6[3], T65: t6[4],
			T81: t8[0], T82: t8[1], T83: t8[2], T84: t8[3], T85: t8[4],
			T131: t13[0], T132: t13[1], T133: t13[2], T134: t13[3], T135: t13[4],
			T181: t18[0], T182: t18[1],
			OtherT: otherT,
		})
	}

	rowCount := len(ssqExcelDatas)
	data := map[string][]int{
		"T1":  tr6[0],
		"T2":  tr6[1],
		"T3":  tr6[2],
		"T4":  tr6[3],
		"T5":  tr6[4],
		"T6":  tr8[0],
		"T7":  tr8[1],
		"T8":  tr8[2],
		"T9":  tr8[3],
		"T10": tr8[4],
		"T11": tr13[0],
		"T12": tr13[1],
		"T13": tr13[2],
		"T14": tr13[3],
		"T15": tr13[4],
		"T16": tr18[0],
		"T17": tr18[1],
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

	fmt.Printf("3 %v\n", time.Now().Sub(startTime).Round(time.Second))

	var hm6Combs []string
	for i, hms := range hm6s {
		if slices.Contains([]int{0, 3, 4}, i) {
			for _, hm := range hms {
				if !slices.Contains(hm6Combs, hm) {
					hm6Combs = append(hm6Combs, hm)
				}
			}
		}
	}

	//fmt.Println(len(hm6Combs))
	oeMs := make(map[string]int)
	for _, hm6Comb := range hm6Combs {
		hmCombs := strings.Split(hm6Comb, ",")
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
	fmt.Println(oeMs)
	fmt.Printf("4 %v\n", time.Now().Sub(startTime).Round(time.Second))

	f, err := excel.CreateNewExcelFile(fmt.Sprintf("双色球截止至%s的数据分析", ssqs[0].DrawTime))
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
	err = SetFirstSheetContent(f, colorSlice, styleID, "历史开奖", ssqExcelDatas)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("5 %v\n", time.Now().Sub(startTime).Round(time.Second))

	typ2FrontHms := make(map[string][]string)
	combTyp2FrontHms := make(map[string][]string)
	tNum := 1
	for _, hms := range hm6s {
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}
	for _, hms := range hm8s {
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}
	for _, hms := range hm13s {
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}
	for _, hms := range hm18s {
		typ2FrontHms[fmt.Sprintf("类型T%d", tNum)] = hms
		combTyp2FrontHms[fmt.Sprintf("T%d", tNum)] = hms
		tNum++
	}

	fmt.Println("列顺序:", colNames)

	bestCombSts := make([]*BestCombSt, 0)
	// 要选择 K 列
	for K := 1; K <= 17; K++ {
		bestCombos, bestCover := maxCoverageAll(colsMasks, colNames, K)

		fmt.Println(K)
		fmt.Printf("最大覆盖行数: %d 最佳组合数量: %d\n", bestCover, len(bestCombos))
		fmt.Println("最佳组合如下:")

		for _, combo := range bestCombos {
			fmt.Printf("\t%v\n", combo)
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
	fmt.Printf("6 %v\n", time.Now().Sub(startTime).Round(time.Second))

	typ2ae := make(map[string]models.SsqMoni)
	tNum = 1
	for _, moni := range moni666663s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}
	for _, moni := range moni855555s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}
	for _, moni := range moni1344444s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}
	for _, moni := range moni1833333s {
		typ2ae[fmt.Sprintf("类型T%d", tNum)] = moni
		tNum++
	}

	typ2ae["OtherT"] = models.SsqMoni{
		Cs:   otherTCs,
		Comb: len(otherTxHms),
	}

	err = SaveTypInfoToExcelFile(f, colorSlice, styleID, "相关类型", typ2ae)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("7 %v\n", time.Now().Sub(startTime).Round(time.Second))

	var ssqBackHisData, ssqOeHisData, ssqHzHisData, ssqQzhHisData, ssqFrontOnlyOneHisData []SsqHis
	var wg sync.WaitGroup
	hisTypCount := 5
	for i := 0; i < hisTypCount; i++ {
		wg.Add(1)
		if i == 0 {
			go func() {
				ssqBackHisData = SsqBackHis(&wg)
			}()

		}
		if i == 1 {
			go func() {
				ssqOeHisData = SsqOeHis(&wg)
			}()
		}
		if i == 2 {
			go func() {
				ssqHzHisData = SsqHzHis(&wg)
			}()
		}
		if i == 3 {
			go func() {
				ssqQzhHisData = SsqQzhHis(&wg)
			}()
		}
		if i == 4 {
			go func() {
				ssqFrontOnlyOneHisData = SsqFrontOnlyOneHis(&wg)
			}()
		}
	}
	wg.Wait()
	fmt.Printf("8 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 前区单号历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "前区单号历史", ssqFrontOnlyOneHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("9 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 后区历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "后区历史", ssqBackHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("10 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 奇偶历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "奇偶历史", ssqOeHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("11 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 和值历史
	err = SaveHisToExcelFile(f, colorSlice, styleID, "和值历史", ssqHzHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("12 %v\n", time.Now().Sub(startTime).Round(time.Second))

	err = SaveHisToExcelFile(f, colorSlice, styleID, "前中后历史", ssqQzhHisData)
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("13 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 后区趋势1
	err = SaveBackQuShi1ToExcelFile(f, colorSlice, styleID, "后区趋势", SsqBackQuShi1())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("14 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 奇偶趋势
	err = SaveOeQuShiToExcelFile(f, colorSlice, styleID, "奇偶趋势", SsqOeQuShi())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("15 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 和值趋势1
	err = SaveHzQuShi1ToExcelFile(f, colorSlice, styleID, "和值趋势1", SsqHzQuShi1())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("16 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 和值趋势2
	err = SaveHzQuShi2ToExcelFile(f, colorSlice, styleID, "和值趋势2", SsqHzQuShi2())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("17 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 前中后趋势
	err = SaveQzhQuShiToExcelFile(f, colorSlice, styleID, "前中后趋势", SsqQzhQuShi())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("18 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 重号历史
	err = SaveCHongHaoToExcelFile(f, colorSlice, styleID, "重号历史", SsqCHongHao())
	if err != nil {
		lg.InfoToFileAndStdOut(fmt.Sprintf("%v\n", err))
	}
	fmt.Printf("19 %v\n", time.Now().Sub(startTime).Round(time.Second))

	tNum = 1
	for tNum <= len(typ2FrontHms) {
		typName := fmt.Sprintf("类型T%d", tNum)
		_ = SaveTypHmsToExcelFile(f, colorSlice, styleID, typName, typ2FrontHms[typName])
		tNum++
	}
	fmt.Printf("20 %v\n", time.Now().Sub(startTime).Round(time.Second))
	typName := fmt.Sprintf("OtherT")
	_ = SaveTypHmsToExcelFile(f, colorSlice, styleID, typName, otherTxHms)
	fmt.Printf("21 %v\n", time.Now().Sub(startTime).Round(time.Second))

}

func SetFirstSheetContent(f *excelize.File, colorSlice []int, styleID int, sheetName string, data []SsqExcelData) (err error) {
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
		YSplit: 1,    // 冻结第一行
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
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 10)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "期号")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 10)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "开奖日期")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 12)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "号码")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 20)
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
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "前区前中后")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 10)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "新4重号数")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 12)
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
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "当新重")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)
	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "累重")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 12)

	colNum0++
	colZiMu = excel.GetColumnStr(colNum0)
	_ = f.SetCellValue(sheetName, fmt.Sprintf("%s1", colZiMu), "类型T1")
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 8)

	colWidth := float64(7)
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
	_ = f.SetColWidth(sheetName, colZiMu, colZiMu, 10)

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
			iData.T61, iData.T62, iData.T63, iData.T64, iData.T65,
			iData.T81, iData.T82, iData.T83, iData.T84, iData.T85,
			iData.T131, iData.T132, iData.T133, iData.T134, iData.T135,
			iData.T181, iData.T182,
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

func SaveTypInfoToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, typ2ae map[string]models.SsqMoni) (err error) {
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
			_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), fmt.Sprintf("%.2f", float64(1107568)/(float64(tempNum))))
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
			case "666663":
				tempNum = 6 * 6 * 6 * 6 * 6 * 3
			case "855555":
				tempNum = 8 * 5 * 5 * 5 * 5 * 5
			case "1344444":
				tempNum = 13 * 4 * 4 * 4 * 4 * 4
			case "1833333":
				tempNum = 18 * 3 * 3 * 3 * 3 * 3
			}

			_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), fmt.Sprintf("%d", tempNum))
			_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), fmt.Sprintf("%.2f", float64(1107568)/(float64(tempNum))))

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

	_ = f.SetColWidth(sheetName, "A", "A", 20)
	rowNum := 1
	for _, frontHm := range frontHms {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), frontHm)
		rowNum++
	}
	return
}

func SaveHisToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, ssqHisSlice []SsqHis) (err error) {
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
	for _, ssqHis := range ssqHisSlice {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), tNum)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), ssqHis.Typ)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), ssqHis.Cs)
		_ = f.SetCellStyle(sheetName, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("C%d", rowNum), greenStyle)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), ssqHis.AllCount)

		_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), ssqHis.Last10)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), ssqHis.Last20)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), ssqHis.Last30)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), ssqHis.Last50)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), ssqHis.Last100)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), ssqHis.Last200)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), ssqHis.Last500)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), ssqHis.Last1000)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), ssqHis.Last1500)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("N%d", rowNum), ssqHis.Last2000)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("O%d", rowNum), ssqHis.Last2500)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("P%d", rowNum), ssqHis.Last3500)

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

func SaveBackQuShi2ToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, quShiSlice []SsqBackChQuShi) (err error) {
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

func SaveHzQuShi2ToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, quShiSlice []SsqHzChQuShi) (err error) {
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

func SaveCHongHaoToExcelFile(f *excelize.File, colorSlice []int, styleID int, sheetName string, chSlice []SsqChSt) (err error) {
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
		for _, ssqInfo := range ch.SsqInfos {
			ziMu := excel.GetColumnStr(i)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", ziMu, rowNum), ssqInfo)
			i++
		}

		tNum++
		rowNum++
	}
	_ = f.AutoFilter(sheetName, "A1:E1", nil)
	return
}
