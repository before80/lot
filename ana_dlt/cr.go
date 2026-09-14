package ana_dlt

import (
	"sync"

	"github.com/before80/lot/models"
)

//
//type EqCrCs struct {
//	Eq  int // 作为分析的那一期的设备号
//	Cr  int // 跨多少期进行统计(不包含作为分析的那一期, 也就是只能是作为分析的那一期的之前历史是的 Cr 期)
//	Typ int // 类型: 11111,2111,311,221,32,41,5 对应 Cs 字段的注释说明
//	Cs  int // 作为分析的那一期的前区单号出现在 5个不同段中的次数(1->1->1->1->1 对应 Typ 中的 11111 类型) 或
//	// 作为分析的那一期的前区单号出现在 4个不同段中的次数(2->1->1->1 对应 Typ 中的 2111 类型) 或
//	// 作为分析的那一期的前区单号出现在 3个不同段中的次数(3->1->1 [对应 Typ 中的 311 类型]以及 2->2->1 [对应 Typ 中的 221 类型]) 或
//	// 作为分析的那一期的前区单号出现在 2个不同段中的次数(3->2 [对应 Typ 中的 32 类型]以及 4->1 [对应 Typ 中的 41 类型]) 或
//	// 作为分析的那一期的前区单号出现在 1个相同段中的次数[对应 Typ 中的 1 类型]
//}
//
//type TypCs struct {
//	Typ int // 类型: 11111,2111,311,221,32,41,5 对应 Cs 字段的注释说明
//	Cs  int // 作为分析的那一期的前区单号出现在 5个不同段中的次数(1->1->1->1->1 对应 Typ 中的 11111 类型) 或
//	// 作为分析的那一期的前区单号出现在 4个不同段中的次数(2->1->1->1 对应 Typ 中的 2111 类型) 或
//	// 作为分析的那一期的前区单号出现在 3个不同段中的次数(3->1->1 [对应 Typ 中的 311 类型]以及 2->2->1 [对应 Typ 中的 221 类型]) 或
//	// 作为分析的那一期的前区单号出现在 2个不同段中的次数(3->2 [对应 Typ 中的 32 类型]以及 4->1 [对应 Typ 中的 41 类型]) 或
//	// 作为分析的那一期的前区单号出现在 1个相同段中的次数[对应 Typ 中的 1 类型]
//}

func CalDuanTyp() (allETMap map[int]map[int]map[int]map[string]int) {
	if len(DxDlts) == 0 {
		InitDlts()
	}
	//InitDrawNum2Dlt()
	InitEqDxDlt()

	allETMap = make(map[int]map[int]map[int]map[string]int)
	allETMap[1] = make(map[int]map[int]map[string]int)
	allETMap[2] = make(map[int]map[int]map[string]int)
	allETMap[3] = make(map[int]map[int]map[string]int)
	var wg sync.WaitGroup
	var lc sync.Mutex
	for i := 1; i <= 3; i++ {
		if i == 1 {
			wg.Add(1)
			go func(eqCount int) {
				eTMap := calDuanT(eqCount)
				lc.Lock()
				allETMap[eqCount] = eTMap
				lc.Unlock()
				wg.Done()
			}(1)
		}
		if i == 2 {
			wg.Add(1)
			go func(eqCount int) {
				eTMap := calDuanT(eqCount)
				lc.Lock()
				allETMap[eqCount] = eTMap
				lc.Unlock()
				wg.Done()
			}(2)
		}
		if i == 3 {
			wg.Add(1)
			go func(eqCount int) {
				eTMap := calDuanT(eqCount)
				lc.Lock()
				allETMap[eqCount] = eTMap
				lc.Unlock()
				wg.Done()
			}(3)
		}
	}

	wg.Wait()
	return
}

var CrossDrawNumSli = []int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
}

// LastestStatTyp 用于统计的最近多少期
var LastestStatTyp = []int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 25, 30, 40, 50, 60, 70, 80, 90, 100, 120, 150, 200, 300, 500, 600, 700, 800, 900, 1000, 1200,
}

// StartCrossDrawNum 起始的跨期数
var StartCrossDrawNum = 2

// EndCrossDrawNum 结束的跨期数
var EndCrossDrawNum = 300

func calDuanT(eqCount int) map[int]map[int]map[string]int {
	eTMap := make(map[int]map[int]map[string]int)

	for i := StartCrossDrawNum; i <= EndCrossDrawNum; i++ {
		eTMap[i] = make(map[int]map[string]int)
		ts := CalCommonCrossDrawNumAndEqCount(i, eqCount)

		// 需要分最近多少期进行统计
		for _, l := range LastestStatTyp {
			eTMap[i][l] = make(map[string]int)
			eTMap[i][l] = StatTypWithLastest(ts, l)
		}
	}
	return eTMap
}

// CalCommonCrossDrawNumAndEqCount 计算相同跨期和相同设备号下前区号码的统计数据
//
//	@Description:
//	@param crossDrawNum
//	@param eqCount
//	@return res
func CalCommonCrossDrawNumAndEqCount(crossDrawNum, eqCount int) (res []DltCrossTSN) {
	dlts := make([]models.Dlt, 0)
	switch eqCount {
	case 1:
		dlts = Eq1DxDlts
	case 2:
		dlts = Eq2DxDlts
	case 3:
		dlts = Eq3DxDlts
	default:
		return
	}

	for _, dlt := range dlts {
		curFrontHms := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}
		eDlt := FindNextCommonEqNumDltByDxEqDlt(dlt.DrawNum, dlt.EquipmentCount, crossDrawNum, false)

		eDltTimesSliNumStrSli := dealWaitStatDlt(eDlt, curFrontHms, true)
		res = append(res, DltCrossTSN{
			Cross:   crossDrawNum,
			EqCount: dlt.EquipmentCount,
			Tsn:     eDltTimesSliNumStrSli,
		})
	}
	return
}

// StatTyp
//
//	@Description:
//	@param ts 所有同设备的
//	@return typMap
func StatTyp(ts []DltCrossTSN) (typMap map[string]int) {
	typMap = make(map[string]int)

	for _, t := range ts {
		// 判断是什么类型
		typ := judgeTyp(t.Tsn)
		if _, ok := typMap[typ]; !ok {
			typMap[typ] = 1
		} else {
			typMap[typ]++
		}
	}

	// TODO 降序排序或升序排序

	return
}

func StatTypWithLastest(ts []DltCrossTSN, lastest int) (typMap map[string]int) {
	typMap = make(map[string]int)
	count := 0
	for _, t := range ts {
		// 判断是什么类型
		typ := judgeTyp(t.Tsn)
		if _, ok := typMap[typ]; !ok {
			typMap[typ] = 1
		} else {
			typMap[typ]++
		}
		count++
		if count >= lastest {
			break
		}
	}

	// TODO 降序排序或升序排序

	return
}

// AllDuanTyp 所有段类型
var AllDuanTyp = map[string]string{
	"T11111": "T11111",
	"T2111":  "T2111",
	"T221":   "T221",
	"T311":   "T311",
	"T32":    "T32",
	"T41":    "T41",
	"T5":     "T5",
}

// judgeTyp 判断是什么类型
//
//	@Description:
//	@param dltTimesSliNumStrSli
//	@return string
func judgeTyp(dltTimesSliNumStrSli []DltTimesSliNumStr) string {
	var l int
	var i1s, i2s, i3s, i4s, i5s int
	var existNum int
	totalLen := 0
	for _, dltTimesSliNumStr := range dltTimesSliNumStrSli {
		l = len(dltTimesSliNumStr.NumStr)
		if l > 0 {
			existNum += 1
			totalLen += l
		}
		switch l {
		case 1:
			i1s += 1
		case 2:
			i2s += 1
		case 3:
			i3s += 1
		case 4:
			i4s += 1
		case 5:
			i5s += 1
		default:
		}
		if totalLen == 5 {
			break
		}
		if existNum == 5 {
			break
		}
	}

	if i1s == 5 {
		return AllDuanTyp["T11111"]
	}

	if i1s == 3 && i2s == 1 {
		return AllDuanTyp["T2111"]
	}

	if i1s == 1 && i2s == 2 {
		return AllDuanTyp["T221"]
	}

	if i1s == 2 && i3s == 1 {
		return AllDuanTyp["T311"]
	}

	if i2s == 1 && i3s == 1 {
		return AllDuanTyp["T32"]
	}

	if i1s == 1 && i4s == 1 {
		return AllDuanTyp["T41"]
	}

	if i5s == 1 {
		return AllDuanTyp["T5"]
	}

	return "NotKnown"
}

// FindNextCommonEqNumDltByDxEqDlt
//
//	@Description:
//	@param curDrawNum
//	@param eqCount
//	@param crossDrawNum
//	@param includeCurDrawNum
//	@return nextDlts
func FindNextCommonEqNumDltByDxEqDlt(curDrawNum string, eqCount int, crossDrawNum int, includeCurDrawNum bool) (nextDlts []models.Dlt) {
	dlts := make([]models.Dlt, 0)
	switch eqCount {
	case 1:
		dlts = Eq1DxDlts
	case 2:
		dlts = Eq2DxDlts
	case 3:
		dlts = Eq3DxDlts
	default:
		return
	}
	count := 0
	for _, dlt := range dlts {
		if includeCurDrawNum {
			if dlt.DrawNum > curDrawNum {
				continue
			}
		} else {
			if dlt.DrawNum >= curDrawNum {
				continue
			}
		}

		nextDlts = append(nextDlts, dlt)
		count++
		if count >= crossDrawNum {
			break
		}
	}

	return
}

func CalFrontStatData(startDrawNum string, startCrossDrawNum, endCrossDrawNum, anaDrawNum int) {

	var curDlt models.Dlt
	if startDrawNum == "" {
		curDlt = DxDlts[0]
	} else {
		curDlt = DrawNum2Dlt[startDrawNum]
	}

	var coverStatData []DltCoverCrossTSN
	var statData []DltCrossTSN
	var j int
	for i := 0; i < len(DxDlts); i++ {
		if DxDlts[i].DrawNum > curDlt.DrawNum {
			continue
		}
		j++

		statData = DltFrontSpDrawNumCrossDrawNumStat(DxDlts[i].DrawNum, startCrossDrawNum, endCrossDrawNum, false, true)
		coverStatData = append(coverStatData, DltCoverCrossTSN{
			DrawNum:  DxDlts[i].DrawNum,
			CrossTsn: statData,
		})

		if j >= anaDrawNum || DxDlts[i].EquipmentCount == 0 {
			break
		}
	}

}
