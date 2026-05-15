package ana_dlt

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/before80/lot/cfg"
	"github.com/before80/lot/db"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
)

type XuHaoSt struct {
	Tx               []string
	Oes              []string
	HzMin            int
	HzMax            int
	RemoveHis        int
	FrontDanMaHms    []string
	FrontIncludeHms  []string
	BackIncludeHms   []string
	BackIncludeCombs []string
	FrontExcludeHms  []string
	BackExcludeHms   []string
	BackExcludeCombs []string
	QzhSlice         []string
	Ch4MustGetCount  int
}

func NewXuHaoSt() *XuHaoSt {
	// 选最佳组合
	txStr := cfg.Default.DltTxStr
	tx := strings.Split(txStr, ",")
	// 奇偶
	oes := cfg.Str2StrSliceWithSeparator(cfg.Default.DltOeStr, ",")
	// 和值
	hzMin := cfg.Default.DltHzMin
	hzMax := cfg.Default.DltHzMax

	// 移除历史开奖号码
	removeHis := cfg.Default.DltRemoveHis

	// 必须包含的号码
	//前区胆码
	frontDanMaHms := cfg.Str2StrSliceWithSeparator(cfg.Default.DltFrontDanMaHmStr, ",")
	// 区分前区包含,后区包含
	frontIncludeHms := cfg.Str2StrSliceWithSeparator(cfg.Default.DltFrontIncludeHmStr, ",")
	//fmt.Printf("len(frontIncludeHms)=%d\n", len(frontIncludeHms))
	backIncludeHms := cfg.Str2StrSliceWithSeparator(cfg.Default.DltBackIncludeHmStr, ",")
	//fmt.Printf("len(backIncludeHms)=%d\n", len(backIncludeHms))
	//fmt.Printf("cfg.Default.DltBackIncludeCombStr=%s\n", cfg.Default.DltBackIncludeCombStr)
	backIncludeCombs := cfg.Str2StrSliceWithSeparator(cfg.Default.DltBackIncludeCombStr, "|")
	//fmt.Printf("len(backIncludeCombs)=%d backIncludeCombs=%v\n", len(backIncludeCombs), backIncludeCombs)

	if len(frontDanMaHms) > 4 {
		frontDanMaHms = frontDanMaHms[0:4]
	}

	if len(backIncludeCombs) > 0 {
		newBackIncludeHmCombs := make([]string, len(backIncludeCombs))
		for _, comb := range backIncludeCombs {
			combs := strings.Split(comb, ",")
			a, _ := strconv.Atoi(combs[0])
			b, _ := strconv.Atoi(combs[1])
			if a > b {
				newBackIncludeHmCombs = append(newBackIncludeHmCombs, fmt.Sprintf("%02d,%02d", b, a))
			} else {
				newBackIncludeHmCombs = append(newBackIncludeHmCombs, fmt.Sprintf("%02d,%02d", a, b))
			}
		}
		//fmt.Printf("len(newBackIncludeHmCombs)=%d newBackIncludeHmCombs=%v\n", len(newBackIncludeHmCombs), newBackIncludeHmCombs)
		backIncludeCombs = nil
		backIncludeCombs = newBackIncludeHmCombs
	}

	// 必须排除的号码
	frontExcludeHms := cfg.Str2StrSliceWithSeparator(cfg.Default.DltFrontExcludeHmStr, ",")
	backExcludeHms := cfg.Str2StrSliceWithSeparator(cfg.Default.DltBackExcludeHmStr, ",")
	// 必须排除的后区组合号码
	backExcludeCombs := cfg.Str2StrSliceWithSeparator(cfg.Default.DltBackExcludeCombStr, "|")

	// 前中后号码个数分布
	qzhSlice := cfg.Str2StrSliceWithSeparator(cfg.Default.DltQzhStr, ",")

	// 新增4重号数必须大于等于的数值
	ch4MustGetCount := cfg.Default.DltCh4MustGETCount

	return &XuHaoSt{
		Tx:               tx,
		Oes:              oes,
		HzMin:            hzMin,
		HzMax:            hzMax,
		RemoveHis:        removeHis,
		FrontDanMaHms:    frontDanMaHms,
		FrontIncludeHms:  frontIncludeHms,
		BackIncludeHms:   backIncludeHms,
		BackIncludeCombs: backIncludeCombs,
		FrontExcludeHms:  frontExcludeHms,
		BackExcludeHms:   backExcludeHms,
		BackExcludeCombs: backExcludeCombs,
		QzhSlice:         qzhSlice,
		Ch4MustGetCount:  ch4MustGetCount,
	}
}

func DltXHaoForWeb(xuHaoSt *XuHaoSt, zhuShu int) []string {
	//startTime := time.Now()
	InitDlts()

	// 获取所有类型对应的 moni
	t2Monies := GenTx2Moni()
	//zjStr := "06,11,13,16,22|02,03"

	// 获取当前选中类型的 moni 切片
	selectedMonies := make([]models.DltMoni, 0, len(xuHaoSt.Tx))
	lTx := make([]string, len(xuHaoSt.Tx))
	//fmt.Printf("xuHaoSt.Tx=%v\n", xuHaoSt.Tx)
	for _, t := range xuHaoSt.Tx {
		lt := strings.ToLower(t)
		lTx = append(lTx, lt)
		if lt != "othert" && lt != "OtherT" {
			selectedMonies = append(selectedMonies, t2Monies[t])
		}
	}
	var m1, m2 runtime.MemStats

	// 强制 GC，减少历史干扰（非常重要）
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// 获取选中类型对应的所有开奖号码
	var fullHmStrSlice []string
	//fmt.Printf("1 selectedMonies: %#v\n", selectedMonies)
	fullHmStrSlice = GenFullHmStrSliceFromSelectedMonies(selectedMonies)

	runtime.ReadMemStats(&m2)

	fmt.Printf("Alloc 变化: %d MB\n", (m2.Alloc-m1.Alloc)/1024/1024)
	fmt.Printf("TotalAlloc 变化: %d MB\n", (m2.TotalAlloc-m1.TotalAlloc)/1024/1024)
	fmt.Printf("Sys 变化: %d MB\n", (m2.Sys-m1.Sys)/1024/1024)
	fmt.Printf("HeapAlloc 变化: %d MB\n", (m2.HeapAlloc-m1.HeapAlloc)/1024/1024)

	//fmt.Printf("2 selectedMonies: %#v\n", selectedMonies)
	if slices.Contains(lTx, "othert") || slices.Contains(lTx, "OtherT") {
		otherFullHmStrSlice := GenFullHmStrSliceFromFrontHmSlice(CalDltNotTxFrontHmSlice(CalDltTxFrontHmSlice()))
		fullHmStrSlice = append(fullHmStrSlice, otherFullHmStrSlice...)
	}

	//fmt.Printf("10 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(fullHmStrSlice, zjStr))

	// 获取历史已经开出的号码字符串切片,用于后续过滤使用
	hisFullHmStrSlice := GetDltHisFullHmStrSlice("")
	//fmt.Printf("20 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("该组合的号码有%d注\n", len(fullHmStrSlice))
	// (是否)过滤掉历史已经开出的号码
	var step10FullHmStrSlice []string
	if xuHaoSt.RemoveHis == 1 {
		step10FullHmStrSlice = gen.DiffSlice(fullHmStrSlice, hisFullHmStrSlice)
		//fmt.Printf("过滤掉历史已经开出的号码后还有%d注\n", len(step10FullHmStrSlice))
	} else {
		step10FullHmStrSlice = fullHmStrSlice
	}

	//fmt.Printf("30 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step10FullHmStrSlice, zjStr))
	// 记得复制出多份,用于奇偶和和值过滤

	// 过滤掉非指定奇偶的号码 1
	// 过滤掉非指定奇偶的号码 2
	step20FullHmStrSlice := FilterFullHmStrSliceWithOes(step10FullHmStrSlice, xuHaoSt.Oes)
	//fmt.Printf("仅保留指定奇偶的号码后还有%d注\n", len(step20FullHmStrSlice))
	//fmt.Printf("40 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step20FullHmStrSlice, zjStr))
	// 过滤掉非指定范围内的和值号码
	step30FullHmStrSlice := FilterFullHmStrSliceWithHzRange(step20FullHmStrSlice, xuHaoSt.HzMin, xuHaoSt.HzMax)
	//fmt.Printf("仅保留指定范围内的和值号码后还有%d注 hzMin=%d hzMax=%d\n", len(step30FullHmStrSlice), xuHaoSt.HzMin, xuHaoSt.HzMax)
	//fmt.Printf("50 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step30FullHmStrSlice, zjStr))

	// 过滤掉前区非胆码的号码
	step40FullHmStrSlice := FilterFullHmStrSliceWithIncludeDanMaHmSlice(step30FullHmStrSlice, xuHaoSt.FrontDanMaHms)
	//fmt.Printf("仅保留指定胆码号码后还有%d注\n", len(step40FullHmStrSlice))
	//fmt.Printf("60 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 过滤掉不包含指定号码的号码(幸运号)
	step50FullHmStrSlice := FilterFullHmStrSliceWithIncludeHmSlice(step40FullHmStrSlice, xuHaoSt.FrontIncludeHms, xuHaoSt.BackIncludeHms, xuHaoSt.BackIncludeCombs)
	//fmt.Printf("过滤掉不包含指定号码的号码后还有%d注 FrontIncludeHms=%v BackIncludeHms=%v BackIncludeCombs=%v\n", len(step50FullHmStrSlice), xuHaoSt.FrontIncludeHms, xuHaoSt.BackIncludeHms, xuHaoSt.BackIncludeCombs)
	//fmt.Printf("70 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step50FullHmStrSlice, zjStr))
	// 过滤掉包含指定号码的号码(认为当期不会开的号码)
	step60FullHmStrSlice := FilterFullHmStrSliceWithExcludeHmSlice(step50FullHmStrSlice, xuHaoSt.FrontExcludeHms, xuHaoSt.BackExcludeHms)
	//fmt.Printf("仅保留包含指定号码的号码后还有%d注\n", len(step60FullHmStrSlice))
	//fmt.Printf("80 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step5FullHmStrSlice, zjStr))
	step70FullHmStrSlice := FilterFullHmStrSliceWithBackExcludeHmCombSlice(step60FullHmStrSlice, xuHaoSt.BackExcludeCombs)
	//fmt.Printf("过滤掉包含指定后区组合号码的号码后还有%d注\n", len(step70FullHmStrSlice))
	//fmt.Printf("90 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step70FullHmStrSlice, zjStr))
	// 过滤掉非选定的前中后的号码
	step80FullHmStrSlice := FilterFullHmStrSliceWithQzh(step70FullHmStrSlice, xuHaoSt.QzhSlice)
	//fmt.Printf("仅保留选定的前中后的号码后还有%d注\n", len(step80FullHmStrSlice))
	//fmt.Printf("100 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step80FullHmStrSlice, zjStr))
	//for _, fullHmStr := range step80FullHmStrSlice {
	//	fmt.Println(fullHmStr)
	//}
	//fmt.Println("-----------------")
	// 过滤掉新4重号数不大于指定数值的号码
	cms := DltLastDrawNumNewAddCh4567("")
	lastC4ms := cms[4]
	//fmt.Printf("110 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//lastC5ms := cms[5]
	//lastC6ms := cms[6]
	//lastC7ms := cms[7]
	step90FullHmStrSlice := FilterFullHmStrSliceWithNewAddCh4CountConcurrent(step80FullHmStrSlice, lastC4ms, xuHaoSt.Ch4MustGetCount, 30)
	//fmt.Printf("仅保留新4重号数不大于指定数值的号码后还有%d注\n", len(step90FullHmStrSlice))
	//fmt.Printf("120 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step90FullHmStrSlice, zjStr))
	//for _, fullHmStr := range step90FullHmStrSlice {
	//	fmt.Println(fullHmStr)
	//}
	//
	//fmt.Println("-----------------")
	//随机选择的号码
	finalSelectFullHmStrSlice := make([]string, 0, zhuShu)
	tempNumSlice := make([]int, 0, zhuShu)
	if len(step90FullHmStrSlice) > 0 {
		for i := 0; i < zhuShu; i++ {
		LabelForContinue:
			time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
			tempNum := rand.IntN(len(step90FullHmStrSlice))
			if !slices.Contains(tempNumSlice, tempNum) {
				tempNumSlice = append(tempNumSlice, tempNum)
				finalSelectFullHmStrSlice = append(finalSelectFullHmStrSlice, step90FullHmStrSlice[tempNum])
			} else {
				if zhuShu < len(step90FullHmStrSlice) {
					goto LabelForContinue
				}
			}
		}
	}
	fmt.Printf("----------------------------")
	runtime.ReadMemStats(&m2)

	fmt.Printf("Alloc 变化: %d MB\n", (m2.Alloc-m1.Alloc)/1024/1024)
	fmt.Printf("TotalAlloc 变化: %d MB\n", (m2.TotalAlloc-m1.TotalAlloc)/1024/1024)
	fmt.Printf("Sys 变化: %d MB\n", (m2.Sys-m1.Sys)/1024/1024)
	fmt.Printf("HeapAlloc 变化: %d MB\n", (m2.HeapAlloc-m1.HeapAlloc)/1024/1024)

	return finalSelectFullHmStrSlice
	//for _, fullHmStr := range finalSelectFullHmStrSlice {
	//	fmt.Println(fullHmStr)
	//}
	//time.Sleep(9999 * time.Hour)
}

// DltXHao
//
//	@Description:
//	@param xuHaoSt
//	@param zhuShu 选多少注
//	@param lastDrawNum 目的是为了验证
func DltXHao(xuHaoSt *XuHaoSt, zhuShu int, lastDrawNum string) {
	startTime := time.Now()
	InitDlts()

	// 获取所有类型对应的 moni
	t2Monies := GenTx2Moni()
	//zjStr := "06,11,13,16,22|02,03"

	// 获取当前选中类型的 moni 切片
	selectedMonies := make([]models.DltMoni, 0, len(xuHaoSt.Tx))
	lTx := make([]string, len(xuHaoSt.Tx))
	//fmt.Printf("xuHaoSt.Tx=%v\n", xuHaoSt.Tx)
	for _, t := range xuHaoSt.Tx {
		lt := strings.ToLower(t)
		lTx = append(lTx, lt)
		if lt != "othert" {
			selectedMonies = append(selectedMonies, t2Monies[t])
		}
	}

	// 获取选中类型对应的所有开奖号码
	var fullHmStrSlice []string
	//fmt.Printf("1 selectedMonies: %#v\n", selectedMonies)
	fullHmStrSlice = GenFullHmStrSliceFromSelectedMonies(selectedMonies)

	//fmt.Printf("2 selectedMonies: %#v\n", selectedMonies)
	if slices.Contains(lTx, "othert") {
		otherFullHmStrSlice := GenFullHmStrSliceFromFrontHmSlice(CalDltNotTxFrontHmSlice(CalDltTxFrontHmSlice()))
		fullHmStrSlice = append(fullHmStrSlice, otherFullHmStrSlice...)
	}

	fmt.Printf("10 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(fullHmStrSlice, zjStr))

	// 获取历史已经开出的号码字符串切片,用于后续过滤使用
	hisFullHmStrSlice := GetDltHisFullHmStrSlice(lastDrawNum)
	fmt.Printf("20 %v\n", time.Now().Sub(startTime).Round(time.Second))
	fmt.Printf("该组合的号码有%d注\n", len(fullHmStrSlice))
	// (是否)过滤掉历史已经开出的号码
	var step10FullHmStrSlice []string
	if xuHaoSt.RemoveHis == 1 {
		step10FullHmStrSlice = gen.DiffSlice(fullHmStrSlice, hisFullHmStrSlice)
		fmt.Printf("过滤掉历史已经开出的号码后还有%d注\n", len(step10FullHmStrSlice))
	} else {
		step10FullHmStrSlice = fullHmStrSlice
	}

	fmt.Printf("30 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step10FullHmStrSlice, zjStr))
	// 记得复制出多份,用于奇偶和和值过滤

	// 过滤掉非指定奇偶的号码 1
	// 过滤掉非指定奇偶的号码 2
	step20FullHmStrSlice := FilterFullHmStrSliceWithOes(step10FullHmStrSlice, xuHaoSt.Oes)
	fmt.Printf("仅保留指定奇偶的号码后还有%d注\n", len(step20FullHmStrSlice))
	fmt.Printf("40 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step20FullHmStrSlice, zjStr))
	// 过滤掉非指定范围内的和值号码
	step30FullHmStrSlice := FilterFullHmStrSliceWithHzRange(step20FullHmStrSlice, xuHaoSt.HzMin, xuHaoSt.HzMax)
	fmt.Printf("仅保留指定范围内的和值号码后还有%d注 hzMin=%d hzMax=%d\n", len(step30FullHmStrSlice), xuHaoSt.HzMin, xuHaoSt.HzMax)
	fmt.Printf("50 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step30FullHmStrSlice, zjStr))

	// 过滤掉前区非胆码的号码
	step40FullHmStrSlice := FilterFullHmStrSliceWithIncludeDanMaHmSlice(step30FullHmStrSlice, xuHaoSt.FrontDanMaHms)
	fmt.Printf("仅保留指定胆码号码后还有%d注\n", len(step40FullHmStrSlice))
	fmt.Printf("60 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 过滤掉不包含指定号码的号码(幸运号)
	step50FullHmStrSlice := FilterFullHmStrSliceWithIncludeHmSlice(step40FullHmStrSlice, xuHaoSt.FrontIncludeHms, xuHaoSt.BackIncludeHms, xuHaoSt.BackIncludeCombs)
	fmt.Printf("过滤掉不包含指定号码的号码后还有%d注 FrontIncludeHms=%v BackIncludeHms=%v BackIncludeCombs=%v\n", len(step50FullHmStrSlice), xuHaoSt.FrontIncludeHms, xuHaoSt.BackIncludeHms, xuHaoSt.BackIncludeCombs)
	fmt.Printf("70 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step50FullHmStrSlice, zjStr))
	// 过滤掉包含指定号码的号码(认为当期不会开的号码)
	step60FullHmStrSlice := FilterFullHmStrSliceWithExcludeHmSlice(step50FullHmStrSlice, xuHaoSt.FrontExcludeHms, xuHaoSt.BackExcludeHms)
	fmt.Printf("仅保留包含指定号码的号码后还有%d注\n", len(step60FullHmStrSlice))
	fmt.Printf("80 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step5FullHmStrSlice, zjStr))
	step70FullHmStrSlice := FilterFullHmStrSliceWithBackExcludeHmCombSlice(step60FullHmStrSlice, xuHaoSt.BackExcludeCombs)
	fmt.Printf("过滤掉包含指定后区组合号码的号码后还有%d注\n", len(step70FullHmStrSlice))
	fmt.Printf("90 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step70FullHmStrSlice, zjStr))
	// 过滤掉非选定的前中后的号码
	step80FullHmStrSlice := FilterFullHmStrSliceWithQzh(step70FullHmStrSlice, xuHaoSt.QzhSlice)
	fmt.Printf("仅保留选定的前中后的号码后还有%d注\n", len(step80FullHmStrSlice))
	fmt.Printf("100 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step80FullHmStrSlice, zjStr))
	//for _, fullHmStr := range step80FullHmStrSlice {
	//	fmt.Println(fullHmStr)
	//}
	//fmt.Println("-----------------")
	// 过滤掉新4重号数不大于指定数值的号码
	cms := DltLastDrawNumNewAddCh4567(lastDrawNum)
	lastC4ms := cms[4]
	fmt.Printf("110 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//lastC5ms := cms[5]
	//lastC6ms := cms[6]
	//lastC7ms := cms[7]
	step90FullHmStrSlice := FilterFullHmStrSliceWithNewAddCh4CountConcurrent(step80FullHmStrSlice, lastC4ms, xuHaoSt.Ch4MustGetCount, 30)
	fmt.Printf("仅保留新4重号数不大于指定数值的号码后还有%d注\n", len(step90FullHmStrSlice))
	fmt.Printf("120 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//fmt.Printf("exist = %t\n", slices.Contains(step90FullHmStrSlice, zjStr))
	//for _, fullHmStr := range step90FullHmStrSlice {
	//	fmt.Println(fullHmStr)
	//}
	//
	//fmt.Println("-----------------")
	//随机选择的号码
	finalSelectFullHmStrSlice := make([]string, 0, zhuShu)
	tempNumSlice := make([]int, 0, zhuShu)
	if len(step90FullHmStrSlice) > 0 {
		for i := 0; i < zhuShu; i++ {
		LabelForContinue:
			time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
			tempNum := rand.IntN(len(step90FullHmStrSlice))
			if !slices.Contains(tempNumSlice, tempNum) {
				tempNumSlice = append(tempNumSlice, tempNum)
				finalSelectFullHmStrSlice = append(finalSelectFullHmStrSlice, step90FullHmStrSlice[tempNum])
			} else {
				if zhuShu < len(step90FullHmStrSlice) {
					goto LabelForContinue
				}
			}
		}
	}

	for _, fullHmStr := range finalSelectFullHmStrSlice {
		fmt.Println(fullHmStr)
	}
	time.Sleep(9999 * time.Hour)
}

// GenTx2MoniABCDE
//
//	@Description:
//	@return t2MoniABCDEs
func GenTx2MoniABCDE() (t2MoniABCDEs map[string]map[string][]string) {
	t2MoniABCDEs = make(map[string]map[string][]string, 17)

	// 找到类型为77777匹配历史开奖号码最多的前5种
	var moni77777s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "77777", "method": "11111"}).Order("cs desc,id asc").Limit(5).Find(&moni77777s)
	// 找到类型为116666匹配历史开奖号码最多的前5种
	var moni116666s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "116666", "method": "11111"}).Order("cs desc,id asc").Limit(5).Find(&moni116666s)
	// 找到类型为155555匹配历史开奖号码最多的前5种
	var moni155555s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "155555", "method": "11111"}).Order("cs desc,id asc").Limit(5).Find(&moni155555s)
	// 找到类型为194444匹配历史开奖号码最多的前2种
	var moni194444s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "194444", "method": "11111"}).Order("cs desc,id asc").Limit(2).Find(&moni194444s)

	if len(moni77777s) != 5 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到5个77777类型的行数据"), 1)
		return
	}
	if len(moni116666s) != 5 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到5个116666类型的行数据"), 1)
		return
	}
	if len(moni155555s) != 5 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到5个155555类型的行数据"), 1)
		return
	}

	if len(moni194444s) != 2 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到2个194444类型的行数据"), 1)
		return
	}

	typNum := 1

	for _, moni := range moni77777s {
		if _, ok := t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]; !ok {
			t2MoniABCDEs[fmt.Sprintf("T%d", typNum)] = make(map[string][]string)
		}

		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["A"] = strings.Split(moni.A, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["B"] = strings.Split(moni.B, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["C"] = strings.Split(moni.C, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["D"] = strings.Split(moni.D, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["E"] = strings.Split(moni.E, ",")
		typNum++
	}

	for _, moni := range moni116666s {
		if _, ok := t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]; !ok {
			t2MoniABCDEs[fmt.Sprintf("T%d", typNum)] = make(map[string][]string)
		}

		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["A"] = strings.Split(moni.A, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["B"] = strings.Split(moni.B, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["C"] = strings.Split(moni.C, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["D"] = strings.Split(moni.D, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["E"] = strings.Split(moni.E, ",")
		typNum++
	}

	for _, moni := range moni155555s {
		if _, ok := t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]; !ok {
			t2MoniABCDEs[fmt.Sprintf("T%d", typNum)] = make(map[string][]string)
		}

		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["A"] = strings.Split(moni.A, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["B"] = strings.Split(moni.B, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["C"] = strings.Split(moni.C, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["D"] = strings.Split(moni.D, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["E"] = strings.Split(moni.E, ",")
		typNum++
	}

	for _, moni := range moni194444s {
		if _, ok := t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]; !ok {
			t2MoniABCDEs[fmt.Sprintf("T%d", typNum)] = make(map[string][]string)
		}

		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["A"] = strings.Split(moni.A, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["B"] = strings.Split(moni.B, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["C"] = strings.Split(moni.C, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["D"] = strings.Split(moni.D, ",")
		t2MoniABCDEs[fmt.Sprintf("T%d", typNum)]["E"] = strings.Split(moni.E, ",")
		typNum++
	}

	return
}

// GenTx2Moni 生成 Tx对应的Moni
//
//	@Description:
//	@return t2Monis
func GenTx2Moni() (t2Monis map[string]models.DltMoni) {
	t2Monis = make(map[string]models.DltMoni, 18)

	// 找到类型为77777匹配历史开奖号码最多的前5种
	var moni77777s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "77777", "method": "11111"}).Order("cs desc,id asc").Limit(5).Find(&moni77777s)
	// 找到类型为116666匹配历史开奖号码最多的前5种
	var moni116666s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "116666", "method": "11111"}).Order("cs desc,id asc").Limit(5).Find(&moni116666s)
	// 找到类型为155555匹配历史开奖号码最多的前5种
	var moni155555s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "155555", "method": "11111"}).Order("cs desc,id asc").Limit(5).Find(&moni155555s)
	// 找到类型为194444匹配历史开奖号码最多的前2种
	var moni194444s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "194444", "method": "11111"}).Order("cs desc,id asc").Limit(2).Find(&moni194444s)

	if len(moni77777s) != 5 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到5个77777类型的行数据"), 1)
		return
	}
	if len(moni116666s) != 5 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到5个116666类型的行数据"), 1)
		return
	}
	if len(moni155555s) != 5 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到5个155555类型的行数据"), 1)
		return
	}

	if len(moni194444s) != 2 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到2个194444类型的行数据"), 1)
		return
	}

	typNum := 1
	for _, moni := range moni77777s {
		t2Monis[fmt.Sprintf("T%d", typNum)] = moni
		typNum++
	}

	for _, moni := range moni116666s {
		t2Monis[fmt.Sprintf("T%d", typNum)] = moni
		typNum++
	}

	for _, moni := range moni155555s {
		t2Monis[fmt.Sprintf("T%d", typNum)] = moni
		typNum++
	}

	for _, moni := range moni194444s {
		t2Monis[fmt.Sprintf("T%d", typNum)] = moni
		typNum++
	}
	t2Monis["OtherT"] = models.DltMoni{}

	return
}

// GenFullHmStrSliceFromSelectedMonies 从选定的组合中获取完整开奖号码字符串切片
//
//	@Description:
//	@param selectedMonies
//	@return fullHmStrSlice
func GenFullHmStrSliceFromSelectedMonies(selectedMonies []models.DltMoni) (fullHmStrSlice []string) {
	frontHmStrSliceSli := BuildFrontHmSlices(selectedMonies, len(selectedMonies))
	var frontHmStrSlice []string
	for _, iFrontHmStrSlice := range frontHmStrSliceSli {
		frontHmStrSlice = append(frontHmStrSlice, iFrontHmStrSlice...)
	}
	// 加上后区号码
	backHmStrSlice := gen.Comb(gen.AllDltBackHms, 2)

	fullHmStrSlice = make([]string, 0, len(frontHmStrSlice)*66)
	for _, frontHmStr := range frontHmStrSlice {
		for _, backHmStr := range backHmStrSlice {
			fullHmStrSlice = append(fullHmStrSlice, frontHmStr+"|"+backHmStr)
		}
	}
	return
}

// GenFullHmStrSliceFromFrontHmSlice 从给定前区号码中获取完整开奖号码字符串切片
//
//	@Description:
//	@param frontHmStrSlice
//	@return fullHmStrSlice
func GenFullHmStrSliceFromFrontHmSlice(frontHmStrSlice []string) (fullHmStrSlice []string) {
	// 加上后区号码
	backHmStrSlice := gen.Comb(gen.AllDltBackHms, 2)

	fullHmStrSlice = make([]string, 0, len(frontHmStrSlice)*66)
	for _, frontHmStr := range frontHmStrSlice {
		for _, backHmStr := range backHmStrSlice {
			fullHmStrSlice = append(fullHmStrSlice, frontHmStr+"|"+backHmStr)
		}
	}
	return
}

// GetDltHisFullHmStrSlice 请预先执行 ana_dlt.InitDlts
//
//	@Description:
//	@return hisFullHmStrSlice
func GetDltHisFullHmStrSlice(lastDrawNum string) (hisFullHmStrSlice []string) {
	hisFullHmStrSlice = make([]string, 0, len(ZxDlts))
	for _, dlt := range ZxDlts {
		if lastDrawNum != "" && dlt.DrawNum > lastDrawNum {
			break
		}
		hisFullHmStrSlice = append(hisFullHmStrSlice, fmt.Sprintf("%s,%s,%s,%s,%s|%s,%s", dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2))
	}
	return
}

func FilterFullHmStrSliceWithOes(fullHmStrSlice []string, oes []string) (resFullHmStrSlice []string) {
	for _, fullHmStr := range fullHmStrSlice {
		curOe := CalDltOeFromStr(fullHmStr)
		if slices.Contains(oes, curOe) {
			resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
		}
	}
	return
}

func FilterFullHmStrSliceWithHzRange(fullHmStrSlice []string, hzMin, hzMax int) (resFullHmStrSlice []string) {
	for _, fullHmStr := range fullHmStrSlice {
		curHz := CalDltHzFromStr(fullHmStr)
		if curHz >= hzMin && curHz <= hzMax {
			resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
		}
	}
	return
}

func FilterFullHmStrSliceWithIncludeDanMaHmSlice(fullHmStrSlice []string, frontIncludeDanMaHmSlice []string) (resFullHmStrSlice []string) {
	if len(frontIncludeDanMaHmSlice) == 0 {
		//fmt.Println("run here 1 ")
		return fullHmStrSlice
	}
	//fmt.Println("run here 2 ")
	danMaLen := len(frontIncludeDanMaHmSlice)
	for _, fullHmStr := range fullHmStrSlice {
		// 区分前区包含,还是后区包含
		fbs := strings.Split(fullHmStr, "|")
		frontHmSlice := strings.Split(fbs[0], ",")

		if len(gen.SliceIntersection(frontHmSlice, frontIncludeDanMaHmSlice)) == danMaLen {
			//fmt.Println("8")
			resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
		}
	}
	return
}

func FilterFullHmStrSliceWithIncludeHmSlice(fullHmStrSlice []string, frontIncludeHmSlice, backIncludeHmSlice []string, backIncludeCombSlice []string) (resFullHmStrSlice []string) {
	if len(frontIncludeHmSlice) == 0 && len(backIncludeHmSlice) == 0 && len(backIncludeCombSlice) == 0 {
		//fmt.Println("run here 1 ")
		return fullHmStrSlice
	}
	//fmt.Println("run here 2 ")
	for _, fullHmStr := range fullHmStrSlice {
		// 区分前区包含,还是后区包含
		fbs := strings.Split(fullHmStr, "|")
		frontHmSlice := strings.Split(fbs[0], ",")
		backHmSlice := strings.Split(fbs[1], ",")
		backCombStr := fbs[1]

		if len(frontIncludeHmSlice) != 0 && len(backIncludeHmSlice) != 0 && len(backIncludeCombSlice) != 0 {
			//fmt.Println("1")
			if gen.HasIntersection(frontIncludeHmSlice, frontHmSlice) && (gen.HasIntersection(backIncludeHmSlice, backHmSlice) || gen.HasIntersection(backIncludeCombSlice, []string{backCombStr})) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		} else if len(frontIncludeHmSlice) != 0 && len(backIncludeHmSlice) != 0 && len(backIncludeCombSlice) == 0 {
			//fmt.Println("2")
			if gen.HasIntersection(frontIncludeHmSlice, frontHmSlice) && gen.HasIntersection(backIncludeHmSlice, backHmSlice) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		} else if len(frontIncludeHmSlice) != 0 && len(backIncludeHmSlice) == 0 && len(backIncludeCombSlice) != 0 {
			//fmt.Println("3")
			if gen.HasIntersection(frontIncludeHmSlice, frontHmSlice) && gen.HasIntersection(backIncludeCombSlice, []string{backCombStr}) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		} else if len(frontIncludeHmSlice) == 0 && len(backIncludeHmSlice) != 0 && len(backIncludeCombSlice) != 0 {
			//fmt.Println("4")
			if gen.HasIntersection(backIncludeHmSlice, backHmSlice) || gen.HasIntersection(backIncludeCombSlice, []string{backCombStr}) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		} else if len(frontIncludeHmSlice) != 0 && len(backIncludeHmSlice) == 0 && len(backIncludeCombSlice) == 0 {
			//fmt.Println("5")
			if gen.HasIntersection(backIncludeHmSlice, frontHmSlice) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		} else if len(frontIncludeHmSlice) == 0 && len(backIncludeHmSlice) != 0 && len(backIncludeCombSlice) == 0 {
			//fmt.Println("6")
			if gen.HasIntersection(backIncludeHmSlice, backHmSlice) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		} else if len(frontIncludeHmSlice) == 0 && len(backIncludeHmSlice) == 0 && len(backIncludeCombSlice) != 0 {
			//fmt.Println("7")
			if gen.HasIntersection(backIncludeCombSlice, []string{backCombStr}) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		} else if len(frontIncludeHmSlice) == 0 && len(backIncludeHmSlice) == 0 && len(backIncludeCombSlice) == 0 {
			//fmt.Println("8")
			resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
		}
	}
	return
}

func FilterFullHmStrSliceWithExcludeHmSlice(fullHmStrSlice []string, frontExcludeHmSlice, backExcludeHmSlice []string) (resFullHmStrSlice []string) {
	if len(frontExcludeHmSlice) == 0 && len(backExcludeHmSlice) == 0 {
		return fullHmStrSlice
	}

	for _, fullHmStr := range fullHmStrSlice {
		// 区分前区包含,还是后区包含
		fbs := strings.Split(fullHmStr, "|")
		frontHmSlice := strings.Split(fbs[0], ",")
		backHmSlice := strings.Split(fbs[1], ",")

		if len(frontExcludeHmSlice) != 0 && len(backExcludeHmSlice) != 0 {
			if !gen.HasIntersection(backExcludeHmSlice, backHmSlice) && !gen.HasIntersection(frontExcludeHmSlice, frontHmSlice) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		} else if len(backExcludeHmSlice) != 0 {
			if !gen.HasIntersection(backExcludeHmSlice, backHmSlice) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		} else if len(frontExcludeHmSlice) != 0 {
			if !gen.HasIntersection(frontExcludeHmSlice, frontHmSlice) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		}
	}
	return
}

func FilterFullHmStrSliceWithBackExcludeHmCombSlice(fullHmStrSlice []string, backExcludeHmCombSlice []string) (resFullHmStrSlice []string) {
	if len(backExcludeHmCombSlice) == 0 {
		return fullHmStrSlice
	}

	for _, fullHmStr := range fullHmStrSlice {
		// 区分前区包含,还是后区包含
		fbs := strings.Split(fullHmStr, "|")
		//frontHmSlice := strings.Split(fbs[0], ",")
		if !gen.HasIntersection(backExcludeHmCombSlice, []string{fbs[1]}) {
			resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
		}
	}
	return
}

// DltLastDrawNumNewAddCh4567 截止最新一期大乐透重号情况
//
//	@Description:
//	@return lastDrawNumNewAddCh4567Ms
func DltLastDrawNumNewAddCh4567(lastDrawNum string) (lastDrawNumNewAddCh4567Ms map[int]map[string][]string) {
	dlts, _ := dbop.ReadAllDlt(false)

	c7ms := make(map[string][]string)
	c6ms := make(map[string][]string)
	c5ms := make(map[string][]string)
	c4ms := make(map[string][]string)

	for _, dlt := range dlts {
		if lastDrawNum != "" && dlt.DrawNum > lastDrawNum {
			break
		}
		var ic7s, ic6s, ic5s, ic4s []string
		ic7s, ic6s, ic5s, ic4s = nil, nil, nil, nil
		// 从开奖号码中生成6个组合号码
		ic7s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 2)...)
		// 从开奖号码中生成6个组合号码
		ic6s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 1)...)
		ic6s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 2)...)
		// 从开奖号码中生成5个组合号码
		ic5s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 0)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 1)...)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 2)...)
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 1)...)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 2, 2)...)

		for _, c7 := range ic7s {
			c7ms[c7] = append(c7ms[c7], fmt.Sprintf("%s,%s,%s,%s,%s|%s,%s", dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2))
		}

		for _, c6 := range ic6s {
			c6ms[c6] = append(c6ms[c6], fmt.Sprintf("%s,%s,%s,%s,%s|%s,%s", dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2))
		}

		for _, c5 := range ic5s {
			c5ms[c5] = append(c5ms[c5], fmt.Sprintf("%s,%s,%s,%s,%s|%s,%s", dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2))
		}
		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], fmt.Sprintf("%s,%s,%s,%s,%s|%s,%s", dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2))
		}
	}
	lastDrawNumNewAddCh4567Ms = make(map[int]map[string][]string, 4)

	lastDrawNumNewAddCh4567Ms[4] = c4ms
	lastDrawNumNewAddCh4567Ms[5] = c5ms
	lastDrawNumNewAddCh4567Ms[6] = c6ms
	lastDrawNumNewAddCh4567Ms[7] = c7ms
	return
}

// CalNewAddCh4Count 该号码作为下一期开奖号码情况下,计算新增4重号数
//
//	@Description:
//	@param fullHmStr
//	@param lastC4ms
//	@return addCh4Count
func CalNewAddCh4Count(fullHmStr string, lastC4ms map[string][]string) (addCh4Count int) {
	//fullHmStrSlice := strings.Split(strings.Replace(fullHmStr, "|", ",", -1), ",")
	// 高效切分，避免 Replace + Split
	fullHmStrSlice := strings.FieldsFunc(fullHmStr, func(r rune) bool {
		return r == '|' || r == ','
	})

	//neLastC4ms := make(map[string][]string, len(lastC4ms))
	//for k, v := range lastC4ms {
	//	neLastC4ms[k] = v
	//}

	//c4SLLjPrev := CalCHLjFromStrSlice(neLastC4ms)
	//fmt.Printf("len(neLastC4ms)= %v\n", len(neLastC4ms))
	ic4s := make([]string, 0)
	// 从开奖号码中生成4个组合号码
	ic4s = gen.CrossComb([]string{fullHmStrSlice[0], fullHmStrSlice[1], fullHmStrSlice[2], fullHmStrSlice[3], fullHmStrSlice[4]}, []string{fullHmStrSlice[5], fullHmStrSlice[6]}, 4, 0)
	ic4s = append(ic4s, gen.CrossComb([]string{fullHmStrSlice[0], fullHmStrSlice[1], fullHmStrSlice[2], fullHmStrSlice[3], fullHmStrSlice[4]}, []string{fullHmStrSlice[5], fullHmStrSlice[6]}, 3, 1)...)
	ic4s = append(ic4s, gen.CrossComb([]string{fullHmStrSlice[0], fullHmStrSlice[1], fullHmStrSlice[2], fullHmStrSlice[3], fullHmStrSlice[4]}, []string{fullHmStrSlice[5], fullHmStrSlice[6]}, 2, 2)...)
	//for _, c4 := range ic4s {
	//	neLastC4ms[c4] = append(neLastC4ms[c4], fullHmStr)
	//}
	//curC4SLLj := CalCHLjFromStrSlice(neLastC4ms)
	//return curC4SLLj - c4SLLjPrev

	// 严格按 len(v) > 1 的规则计算增量
	for _, c4 := range ic4s {
		n := len(lastC4ms[c4])
		if n == 1 {
			addCh4Count += 2
		} else if n >= 2 {
			addCh4Count += 1
		}
	}
	return
}

func FilterFullHmStrSliceWithNewAddCh4Count(fullHmStrSlice []string, lastC4ms map[string][]string, mustGETNewAddCh4Count int) (resFullHmStrSlice []string) {
	for _, fullHmStr := range fullHmStrSlice {
		if CalNewAddCh4Count(fullHmStr, lastC4ms) >= mustGETNewAddCh4Count {
			resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
		}
		//fmt.Printf("剩%d\n", len(fullHmStrSlice)-i)
	}
	return
}

func FilterFullHmStrSliceWithQzh(fullHmStrSlice []string, qzhSlice []string) (resFullHmStrSlice []string) {
	if len(qzhSlice) == 0 {
		return fullHmStrSlice
	}
	for _, fullHmStr := range fullHmStrSlice {
		frontHmStrSlice := strings.Split(strings.Split(fullHmStr, "|")[0], ",")
		qNum, zNum, hNum := 0, 0, 0
		for _, frontHmStr := range frontHmStrSlice {
			hm, _ := strconv.Atoi(frontHmStr)

			if hm < 18 {
				qNum++
			}
			if hm == 18 {
				zNum++
			}
			if hm > 18 {
				hNum++
			}
		}
		qzhStr := fmt.Sprintf("%d%d%d", qNum, zNum, hNum)
		if slices.Contains(qzhSlice, qzhStr) {
			resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
		}
	}
	return
}

func FilterFullHmStrSliceWithNewAddCh4CountConcurrent(
	fullHmStrSlice []string,
	lastC4ms map[string][]string,
	mustGETNewAddCh4Count int,
	workerNum int,
) (resFullHmStrSlice []string) {

	if workerNum <= 0 {
		workerNum = runtime.NumCPU()
	}

	jobs := make(chan string)
	results := make(chan string)

	var wg sync.WaitGroup

	// 启动 worker
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fullHmStr := range jobs {
				if CalNewAddCh4Count(fullHmStr, lastC4ms) >= mustGETNewAddCh4Count {
					results <- fullHmStr
				}
			}
		}()
	}

	// 投递任务
	go func() {
		for _, fullHmStr := range fullHmStrSlice {
			jobs <- fullHmStr
		}
		close(jobs)
	}()

	// 关闭 results
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	for r := range results {
		resFullHmStrSlice = append(resFullHmStrSlice, r)
	}

	return
}
