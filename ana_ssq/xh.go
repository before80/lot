package ana_ssq

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
	Tx              []string
	Oes             []string
	HzMin           int
	HzMax           int
	RemoveHis       int
	FrontIncludeHms []string
	BackIncludeHms  []string
	FrontExcludeHms []string
	BackExcludeHms  []string
	QzhSlice        []string
	Ch4MustGetCount int
}

func NewXuHaoSt() *XuHaoSt {
	// 选最佳组合
	txStr := cfg.Default.SsqTxStr
	tx := strings.Split(txStr, ",")
	// 奇偶
	oes := cfg.Str2StrSliceWithSeparator(cfg.Default.SsqOeStr, ",")
	// 和值
	hzMin := cfg.Default.SsqHzMin
	hzMax := cfg.Default.SsqHzMax

	// 移除历史开奖号码
	removeHis := cfg.Default.SsqRemoveHis

	// 必须包含的号码
	// 区分前区包含,后区包含
	frontIncludeHms := cfg.Str2StrSliceWithSeparator(cfg.Default.SsqFrontIncludeHmStr, ",")
	backIncludeHms := cfg.Str2StrSliceWithSeparator(cfg.Default.SsqBackIncludeHmStr, ",")
	// 必须排除的号码
	frontExcludeHms := cfg.Str2StrSliceWithSeparator(cfg.Default.SsqFrontExcludeHmStr, ",")
	backExcludeHms := cfg.Str2StrSliceWithSeparator(cfg.Default.SsqBackExcludeHmStr, ",")

	// 前中后号码个数分布
	qzhSlice := cfg.Str2StrSliceWithSeparator(cfg.Default.SsqQzhStr, ",")

	// 新增4重号数必须大于等于的数值
	ch4MustGetCount := cfg.Default.SsqCh4MustGETCount

	return &XuHaoSt{
		Tx:              tx,
		Oes:             oes,
		HzMin:           hzMin,
		HzMax:           hzMax,
		RemoveHis:       removeHis,
		FrontIncludeHms: frontIncludeHms,
		BackIncludeHms:  backIncludeHms,
		FrontExcludeHms: frontExcludeHms,
		BackExcludeHms:  backExcludeHms,
		QzhSlice:        qzhSlice,
		Ch4MustGetCount: ch4MustGetCount,
	}
}

func SsqXHao(xuHaoSt *XuHaoSt, zhuShu int, lastDrawNum string) {
	fmt.Printf("%#v", xuHaoSt)
	startTime := time.Now()
	InitSsqs()

	// 获取所有类型对应的 moni
	t2Monies := GenTx2Moni()

	// 获取当前选中类型的 moni 切片
	selectedMonies := make([]models.SsqMoni, 0, len(xuHaoSt.Tx))
	lTx := make([]string, 0, len(xuHaoSt.Tx))
	for _, t := range xuHaoSt.Tx {
		lt := strings.ToLower(t)
		lTx = append(lTx, lt)
		if lt != "othert" {
			selectedMonies = append(selectedMonies, t2Monies[t])
		}
	}

	// 获取选中类型对应的所有开奖号码
	var fullHmStrSlice []string
	fmt.Printf("1 selectedMonies: %#v\n", selectedMonies)
	fullHmStrSlice = GenFullHmStrSliceFromSelectedMonies(selectedMonies)
	fmt.Printf("2 selectedMonies: %#v\n", selectedMonies)
	if slices.Contains(lTx, "othert") {
		otherFullHmStrSlice := GenFullHmStrSliceFromFrontHmSlice(CalSsqNotTxFrontHmSlice(CalSsqTxFrontHmSlice()))
		fullHmStrSlice = append(fullHmStrSlice, otherFullHmStrSlice...)
	}

	fmt.Printf("1 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 获取历史已经开出的号码字符串切片,用于后续过滤使用
	hisFullHmStrSlice := GetSsqHisFullHmStrSlice(lastDrawNum)
	fmt.Printf("2 %v\n", time.Now().Sub(startTime).Round(time.Second))
	fmt.Printf("该组合的号码有%d注\n", len(fullHmStrSlice))
	// (是否)过滤掉历史已经开出的号码
	var step1FullHmStrSlice []string
	if xuHaoSt.RemoveHis == 1 {
		step1FullHmStrSlice = gen.DiffSlice(fullHmStrSlice, hisFullHmStrSlice)
		fmt.Printf("过滤掉历史已经开出的号码后还有%d注\n", len(step1FullHmStrSlice))
	} else {
		step1FullHmStrSlice = fullHmStrSlice
	}

	fmt.Printf("3 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 记得复制出多份,用于奇偶和和值过滤

	// 过滤掉非指定奇偶的号码 1
	// 过滤掉非指定奇偶的号码 2
	step2FullHmStrSlice := FilterFullHmStrSliceWithOes(step1FullHmStrSlice, xuHaoSt.Oes)
	fmt.Printf("过滤掉非指定奇偶的号码后还有%d注\n", len(step2FullHmStrSlice))
	fmt.Printf("4 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 过滤掉非指定范围内的和值号码
	step3FullHmStrSlice := FilterFullHmStrSliceWithHzRange(step2FullHmStrSlice, xuHaoSt.HzMin, xuHaoSt.HzMax)
	fmt.Printf("过滤掉非指定范围内的和值号码后还有%d注 hzMin=%d hzMax=%d\n", len(step3FullHmStrSlice), xuHaoSt.HzMin, xuHaoSt.HzMax)
	fmt.Printf("5 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 过滤掉不包含指定号码的号码(幸运号)
	step4FullHmStrSlice := FilterFullHmStrSliceWithIncludeHmSlice(step3FullHmStrSlice, xuHaoSt.FrontIncludeHms, xuHaoSt.BackIncludeHms)
	fmt.Printf("过滤掉不包含指定号码的号码后还有%d注 %v %v\n", len(step4FullHmStrSlice), xuHaoSt.FrontIncludeHms, xuHaoSt.BackIncludeHms)
	fmt.Printf("6 %v\n", time.Now().Sub(startTime).Round(time.Second))
	// 过滤掉包含指定号码的号码(认为当期不会开的号码)
	step5FullHmStrSlice := FilterFullHmStrSliceWithExcludeHmSlice(step4FullHmStrSlice, xuHaoSt.FrontExcludeHms, xuHaoSt.BackExcludeHms)
	fmt.Printf("过滤掉包含指定号码的号码后还有%d注\n", len(step5FullHmStrSlice))
	fmt.Printf("7 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 过滤掉非选定的前中后的号码
	step6FullHmStrSlice := FilterFullHmStrSliceWithQzh(step5FullHmStrSlice, xuHaoSt.QzhSlice)
	fmt.Printf("过滤掉非选定的前中后的号码后还有%d注\n", len(step6FullHmStrSlice))
	fmt.Printf("8 %v\n", time.Now().Sub(startTime).Round(time.Second))

	// 过滤掉新4重号数不大于指定数值的号码
	cms := SsqLastDrawNumNewAddCh4567(lastDrawNum)
	lastC4ms := cms[4]
	fmt.Printf("9 %v\n", time.Now().Sub(startTime).Round(time.Second))
	//lastC5ms := cms[5]
	//lastC6ms := cms[6]
	//lastC7ms := cms[7]
	step7FullHmStrSlice := FilterFullHmStrSliceWithNewAddCh4CountConcurrent(step6FullHmStrSlice, lastC4ms, xuHaoSt.Ch4MustGetCount, 30)
	fmt.Printf("过滤掉新4重号数不大于指定数值的号码后还有%d注\n", len(step7FullHmStrSlice))
	fmt.Printf("10 %v\n", time.Now().Sub(startTime).Round(time.Second))

	//随机选择的号码
	finalSelectFullHmStrSlice := make([]string, 0, zhuShu)
	tempNumSlice := make([]int, 0, zhuShu)
	if len(step7FullHmStrSlice) > 0 {
		for i := 0; i < zhuShu; i++ {
		LabelForContinue:
			time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
			tempNum := rand.IntN(len(step7FullHmStrSlice))
			if !slices.Contains(tempNumSlice, tempNum) {
				tempNumSlice = append(tempNumSlice, tempNum)
				finalSelectFullHmStrSlice = append(finalSelectFullHmStrSlice, step7FullHmStrSlice[tempNum])
			} else {
				if zhuShu <= len(step7FullHmStrSlice) {
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

// GenTx2Moni 生成 Tx对应的Moni
//
//	@Description:
//	@return t2Monis
func GenTx2Moni() (t2Monis map[string]models.SsqMoni) {
	t2Monis = make(map[string]models.SsqMoni, 18)
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

	if len(moni666663s) != 5 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到5个666663类型的行数据"), 1)
		return
	}
	if len(moni855555s) != 5 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到5个855555类型的行数据"), 1)
		return
	}
	if len(moni1344444s) != 5 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到5个1344444类型的行数据"), 1)
		return
	}

	if len(moni1833333s) != 2 {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("获取不到2个1833333类型的行数据"), 1)
		return
	}

	typNum := 1
	for _, moni := range moni666663s {
		t2Monis[fmt.Sprintf("T%d", typNum)] = moni
		typNum++
	}

	for _, moni := range moni855555s {
		t2Monis[fmt.Sprintf("T%d", typNum)] = moni
		typNum++
	}

	for _, moni := range moni1344444s {
		t2Monis[fmt.Sprintf("T%d", typNum)] = moni
		typNum++
	}

	for _, moni := range moni1833333s {
		t2Monis[fmt.Sprintf("T%d", typNum)] = moni
		typNum++
	}
	t2Monis["OtherT"] = models.SsqMoni{}

	return
}

// GenFullHmStrSliceFromSelectedMonies 从选定的组合中获取完整开奖号码字符串切片
//
//	@Description:
//	@param selectedMonies
//	@return fullHmStrSlice
func GenFullHmStrSliceFromSelectedMonies(selectedMonies []models.SsqMoni) (fullHmStrSlice []string) {
	frontHmStrSliceSli := BuildFrontHmSlices(selectedMonies, len(selectedMonies))
	var frontHmStrSlice []string
	for _, iFrontHmStrSlice := range frontHmStrSliceSli {
		frontHmStrSlice = append(frontHmStrSlice, iFrontHmStrSlice...)
	}
	// 加上后区号码
	backHmStrSlice := gen.AllSsqBackHms

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
	backHmStrSlice := gen.Comb(gen.AllSsqBackHms, 2)

	fullHmStrSlice = make([]string, 0, len(frontHmStrSlice)*66)
	for _, frontHmStr := range frontHmStrSlice {
		for _, backHmStr := range backHmStrSlice {
			fullHmStrSlice = append(fullHmStrSlice, frontHmStr+"|"+backHmStr)
		}
	}
	return
}

// GetSsqHisFullHmStrSlice 请预先执行 ana_ssq.InitSsqs
//
//	@Description:
//	@return hisFullHmStrSlice
func GetSsqHisFullHmStrSlice(lastDrawNum string) (hisFullHmStrSlice []string) {
	hisFullHmStrSlice = make([]string, 0, len(ZxSsqs))
	for _, ssq := range ZxSsqs {
		if lastDrawNum != "" && ssq.DrawNum > lastDrawNum {
			break
		}
		hisFullHmStrSlice = append(hisFullHmStrSlice, fmt.Sprintf("%s,%s,%s,%s,%s,%s|%s", ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6, ssq.B1))
	}
	return
}

func FilterFullHmStrSliceWithOes(fullHmStrSlice []string, oes []string) (resFullHmStrSlice []string) {
	for _, fullHmStr := range fullHmStrSlice {
		curOe := CalSsqOeFromStr(fullHmStr)
		if slices.Contains(oes, curOe) {
			resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
		}
	}
	return
}

func FilterFullHmStrSliceWithHzRange(fullHmStrSlice []string, hzMin, hzMax int) (resFullHmStrSlice []string) {
	for _, fullHmStr := range fullHmStrSlice {
		curHz := CalSsqHzFromStr(fullHmStr)
		if curHz >= hzMin && curHz <= hzMax {
			resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
		}
	}
	return
}

func FilterFullHmStrSliceWithIncludeHmSlice(fullHmStrSlice []string, frontIncludeHmSlice, backIncludeHmSlice []string) (resFullHmStrSlice []string) {
	if len(frontIncludeHmSlice) == 0 && len(backIncludeHmSlice) == 0 {
		fmt.Println("run here 1 ")
		return fullHmStrSlice
	}
	fmt.Println("run here 2 ")
	for _, fullHmStr := range fullHmStrSlice {
		// 区分前区包含,还是后区包含
		fbs := strings.Split(fullHmStr, "|")
		frontHmSlice := strings.Split(fbs[0], ",")
		backHmSlice := strings.Split(fbs[1], ",")

		if len(frontIncludeHmSlice) != 0 && len(backIncludeHmSlice) != 0 {
			if gen.HasIntersection(frontIncludeHmSlice, backHmSlice) && gen.HasIntersection(backIncludeHmSlice, backHmSlice) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
		} else if len(backIncludeHmSlice) != 0 {
			if gen.HasIntersection(backIncludeHmSlice, backHmSlice) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}

		} else if len(frontIncludeHmSlice) != 0 {
			if gen.HasIntersection(frontIncludeHmSlice, frontHmSlice) {
				resFullHmStrSlice = append(resFullHmStrSlice, fullHmStr)
			}
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

// SsqLastDrawNumNewAddCh4567 截止最新一期大乐透重号情况
//
//	@Description:
//	@return lastDrawNumNewAddCh4567Ms
func SsqLastDrawNumNewAddCh4567(lastDrawNum string) (lastDrawNumNewAddCh4567Ms map[int]map[string][]string) {
	ssqs, _ := dbop.ReadAllSsq(false)

	c7ms := make(map[string][]string)
	c6ms := make(map[string][]string)
	c5ms := make(map[string][]string)
	c4ms := make(map[string][]string)

	for _, ssq := range ssqs {
		if lastDrawNum != "" && ssq.DrawNum > lastDrawNum {
			break
		}
		var ic7s, ic6s, ic5s, ic4s []string
		ic7s, ic6s, ic5s, ic4s = nil, nil, nil, nil
		// 从开奖号码中生成6个组合号码
		ic7s = append(ic6s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5}, []string{ssq.F6, ssq.B1}, 5, 2)...)
		// 从开奖号码中生成6个组合号码
		ic6s = append(ic6s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5}, []string{ssq.F6, ssq.B1}, 5, 1)...)
		ic6s = append(ic6s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5}, []string{ssq.F6, ssq.B1}, 4, 2)...)
		// 从开奖号码中生成5个组合号码
		ic5s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5}, []string{ssq.F6, ssq.B1}, 5, 0)
		ic5s = append(ic5s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5}, []string{ssq.F6, ssq.B1}, 4, 1)...)
		ic5s = append(ic5s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5}, []string{ssq.F6, ssq.B1}, 3, 2)...)
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5}, []string{ssq.F6, ssq.B1}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5}, []string{ssq.F6, ssq.B1}, 3, 1)...)
		ic4s = append(ic4s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5}, []string{ssq.F6, ssq.B1}, 2, 2)...)

		for _, c7 := range ic7s {
			c7ms[c7] = append(c7ms[c7], fmt.Sprintf("%s,%s,%s,%s,%s,%s|%s", ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6, ssq.B1))
		}

		for _, c6 := range ic6s {
			c6ms[c6] = append(c6ms[c6], fmt.Sprintf("%s,%s,%s,%s,%s,%s|%s", ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6, ssq.B1))
		}

		for _, c5 := range ic5s {
			c5ms[c5] = append(c5ms[c5], fmt.Sprintf("%s,%s,%s,%s,%s,%s|%s", ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6, ssq.B1))
		}
		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], fmt.Sprintf("%s,%s,%s,%s,%s,%s|%s", ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6, ssq.B1))
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

			if hm < 17 {
				qNum++
			}
			if hm == 17 {
				zNum++
			}
			if hm > 17 {
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
