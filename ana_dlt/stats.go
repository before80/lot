package ana_dlt

import (
	"fmt"
	"slices"

	"github.com/before80/lot/gen"
	"github.com/before80/lot/models"
)

func Stats(zxDlts []models.Dlt, t2MoniABCDEs map[string]map[string][]string) (
	txHis map[string]map[string]int,
	oeHis map[string]map[string]int,
	qzhHis map[string]map[string]int,
	frontDhHis map[string]map[string]int,
	backDhHis map[string]map[string]int,
	backCombHis map[string]map[string]int,
	quShi2St map[string]*DltBackChQuShi,
) {

	lenDltHis := len(zxDlts)
	// 初始化组合
	txHis = make(map[string]map[string]int)
	for _, tx := range gen.AllDltTxs {
		txHis[tx] = make(map[string]int)
		for _, last := range gen.LastHisSlice {
			txHis[tx][last] = 0
		}
	}

	// 初始化奇偶
	oeHis = make(map[string]map[string]int)
	for _, oe := range gen.AllDltOes {
		oeHis[oe] = make(map[string]int)
		for _, last := range gen.LastHisSlice {
			oeHis[oe][last] = 0
		}
	}

	// 初始化前中后
	qzhHis = make(map[string]map[string]int)
	for _, qzh := range gen.AllDltQzhs {
		qzhHis[qzh] = make(map[string]int)
		for _, last := range gen.LastHisSlice {
			qzhHis[qzh][last] = 0
		}
	}

	// 初始化前区单号
	frontDhHis = make(map[string]map[string]int)
	for _, frontHm := range gen.AllDltFrontHms {
		frontDhHis[frontHm] = make(map[string]int)
		for _, last := range gen.LastHisSlice {
			frontDhHis[frontHm][last] = 0
		}
	}

	// 初始化后区单号
	backDhHis = make(map[string]map[string]int)
	for _, backHm := range gen.AllDltBackHms {
		backDhHis[backHm] = make(map[string]int)
		for _, last := range gen.LastHisSlice {
			backDhHis[backHm][last] = 0
		}
	}

	// 初始化后区组合
	backCombHis = make(map[string]map[string]int)
	backCombSlice := gen.Comb(gen.AllDltBackHms, 2)
	for _, backComb := range backCombSlice {
		backCombHis[backComb] = make(map[string]int)
		for _, last := range gen.LastHisSlice {
			backCombHis[backComb][last] = 0
		}
	}

	lastDlt := zxDlts[len(zxDlts)-1]
	quShi2St = make(map[string]*DltBackChQuShi)
	for i, dlt := range zxDlts {
		frontHms := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}
		backHms := []string{dlt.B1, dlt.B2}
		backComb := fmt.Sprintf("%s,%s", dlt.B1, dlt.B2)
		fullHms := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}
		// 计算一注前区号码所属的Tx有哪些
		curFrontBelongToTxs := CalOneFrontHmBelongToTxs(frontHms, t2MoniABCDEs)
		//fmt.Printf("frontHms %v -> %v\n", frontHms, curFrontBelongToTxs)
		curOe := CalDltOe(fullHms)
		curQzh := CalDltQzh(frontHms)

		if dlt.B1 == lastDlt.B1 && dlt.B2 == lastDlt.B2 && i+1 < lenDltHis {
			nextDlt := zxDlts[i+1]
			qs := DltBackQuShiStr(nextDlt, dlt.B1, dlt.B2)
			nextBackComb := fmt.Sprintf("%s,%s", nextDlt.B1, nextDlt.B2)
			if _, ok := quShi2St[qs]; !ok {
				allCombs, _ := gen.GetDltBackQuShiHaoMasFromQuShi(backComb, qs)
				quShi2St[qs] = &DltBackChQuShi{
					HadExistCombs:    nil,
					HadNotExistCombs: nil,
					BackComb:         backComb,
					Qs:               qs,
					Cs:               0,
					AllCombs:         allCombs,
				}
			}
			quShi2St[qs].Cs = quShi2St[qs].Cs + 1
			if !slices.Contains(quShi2St[qs].HadExistCombs, nextBackComb) {
				quShi2St[qs].HadExistCombs = append(quShi2St[qs].HadExistCombs, nextBackComb)
			}
		}

		if lenDltHis-i <= 10 {
			oeHis[curOe]["10"] = oeHis[curOe]["10"] + 1
			qzhHis[curQzh]["10"] = qzhHis[curQzh]["10"] + 1
			backCombHis[backComb]["10"] = backCombHis[backComb]["10"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["10"] = txHis[tx]["10"] + 1
			}
		}
		if lenDltHis-i <= 20 {
			oeHis[curOe]["20"] = oeHis[curOe]["20"] + 1
			qzhHis[curQzh]["20"] = qzhHis[curQzh]["20"] + 1
			backCombHis[backComb]["20"] = backCombHis[backComb]["20"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["20"] = txHis[tx]["20"] + 1
			}
		}
		if lenDltHis-i <= 30 {
			oeHis[curOe]["30"] = oeHis[curOe]["30"] + 1
			qzhHis[curQzh]["30"] = qzhHis[curQzh]["30"] + 1
			backCombHis[backComb]["30"] = backCombHis[backComb]["30"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["30"] = txHis[tx]["30"] + 1
			}
		}
		if lenDltHis-i <= 50 {
			oeHis[curOe]["50"] = oeHis[curOe]["50"] + 1
			qzhHis[curQzh]["50"] = qzhHis[curQzh]["50"] + 1
			backCombHis[backComb]["50"] = backCombHis[backComb]["50"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["50"] = txHis[tx]["50"] + 1
			}
		}
		if lenDltHis-i <= 100 {
			oeHis[curOe]["100"] = oeHis[curOe]["100"] + 1
			qzhHis[curQzh]["100"] = qzhHis[curQzh]["100"] + 1
			backCombHis[backComb]["100"] = backCombHis[backComb]["100"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["100"] = txHis[tx]["100"] + 1
			}
		}
		if lenDltHis-i <= 200 {
			oeHis[curOe]["200"] = oeHis[curOe]["200"] + 1
			qzhHis[curQzh]["200"] = qzhHis[curQzh]["200"] + 1
			backCombHis[backComb]["200"] = backCombHis[backComb]["200"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["200"] = txHis[tx]["200"] + 1
			}
		}
		if lenDltHis-i <= 300 {
			oeHis[curOe]["300"] = oeHis[curOe]["300"] + 1
			qzhHis[curQzh]["300"] = qzhHis[curQzh]["300"] + 1
			backCombHis[backComb]["300"] = backCombHis[backComb]["300"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["300"] = txHis[tx]["300"] + 1
			}
		}

		if lenDltHis-i <= 500 {
			oeHis[curOe]["500"] = oeHis[curOe]["500"] + 1
			qzhHis[curQzh]["500"] = qzhHis[curQzh]["500"] + 1
			backCombHis[backComb]["500"] = backCombHis[backComb]["500"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["500"] = txHis[tx]["500"] + 1
			}
		}
		if lenDltHis-i <= 1000 {
			oeHis[curOe]["1000"] = oeHis[curOe]["1000"] + 1
			qzhHis[curQzh]["1000"] = qzhHis[curQzh]["1000"] + 1
			backCombHis[backComb]["1000"] = backCombHis[backComb]["1000"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["1000"] = txHis[tx]["1000"] + 1
			}
		}
		if lenDltHis-i <= 1500 {
			oeHis[curOe]["1500"] = oeHis[curOe]["1500"] + 1
			qzhHis[curQzh]["1500"] = qzhHis[curQzh]["1500"] + 1
			backCombHis[backComb]["1500"] = backCombHis[backComb]["1500"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["1500"] = txHis[tx]["1500"] + 1
			}
		}
		if lenDltHis-i <= 2000 {
			oeHis[curOe]["2000"] = oeHis[curOe]["2000"] + 1
			qzhHis[curQzh]["2000"] = qzhHis[curQzh]["2000"] + 1
			backCombHis[backComb]["2000"] = backCombHis[backComb]["2000"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["2000"] = txHis[tx]["2000"] + 1
			}
		}
		if lenDltHis-i <= 2500 {
			oeHis[curOe]["2500"] = oeHis[curOe]["2500"] + 1
			qzhHis[curQzh]["2500"] = qzhHis[curQzh]["2500"] + 1
			backCombHis[backComb]["2500"] = backCombHis[backComb]["2500"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["2500"] = txHis[tx]["2500"] + 1
			}
		}
		if lenDltHis-i <= 3500 {
			oeHis[curOe]["3500"] = oeHis[curOe]["3500"] + 1
			qzhHis[curQzh]["3500"] = qzhHis[curQzh]["3500"] + 1
			backCombHis[backComb]["3500"] = backCombHis[backComb]["3500"] + 1
			for _, tx := range curFrontBelongToTxs {
				txHis[tx]["3500"] = txHis[tx]["3500"] + 1
			}
		}

		// 前区单号
		for _, fHm := range frontHms {
			if lenDltHis-i <= 10 {
				frontDhHis[fHm]["10"] = frontDhHis[fHm]["10"] + 1
			}
			if lenDltHis-i <= 20 {
				frontDhHis[fHm]["20"] = frontDhHis[fHm]["20"] + 1
			}
			if lenDltHis-i <= 30 {
				frontDhHis[fHm]["30"] = frontDhHis[fHm]["30"] + 1
			}
			if lenDltHis-i <= 50 {
				frontDhHis[fHm]["50"] = frontDhHis[fHm]["50"] + 1
			}
			if lenDltHis-i <= 100 {
				frontDhHis[fHm]["100"] = frontDhHis[fHm]["100"] + 1
			}
			if lenDltHis-i <= 200 {
				frontDhHis[fHm]["200"] = frontDhHis[fHm]["200"] + 1
			}
			if lenDltHis-i <= 300 {
				frontDhHis[fHm]["300"] = frontDhHis[fHm]["300"] + 1
			}
			if lenDltHis-i <= 500 {
				frontDhHis[fHm]["500"] = frontDhHis[fHm]["500"] + 1
			}
			if lenDltHis-i <= 1000 {
				frontDhHis[fHm]["1000"] = frontDhHis[fHm]["1000"] + 1
			}
			if lenDltHis-i <= 1500 {
				frontDhHis[fHm]["1500"] = frontDhHis[fHm]["1500"] + 1
			}
			if lenDltHis-i <= 2000 {
				frontDhHis[fHm]["2000"] = frontDhHis[fHm]["2000"] + 1
			}
			if lenDltHis-i <= 2500 {
				frontDhHis[fHm]["2500"] = frontDhHis[fHm]["2500"] + 1
			}
			if lenDltHis-i <= 3500 {
				frontDhHis[fHm]["3500"] = frontDhHis[fHm]["3500"] + 1
			}
		}

		// 后区单号
		for _, bHm := range backHms {
			if lenDltHis-i <= 10 {
				backDhHis[bHm]["10"] = backDhHis[bHm]["10"] + 1
			}
			if lenDltHis-i <= 20 {
				backDhHis[bHm]["20"] = backDhHis[bHm]["20"] + 1
			}
			if lenDltHis-i <= 30 {
				backDhHis[bHm]["30"] = backDhHis[bHm]["30"] + 1
			}
			if lenDltHis-i <= 50 {
				backDhHis[bHm]["50"] = backDhHis[bHm]["50"] + 1
			}
			if lenDltHis-i <= 100 {
				backDhHis[bHm]["100"] = backDhHis[bHm]["100"] + 1
			}
			if lenDltHis-i <= 200 {
				backDhHis[bHm]["200"] = backDhHis[bHm]["200"] + 1
			}
			if lenDltHis-i <= 300 {
				backDhHis[bHm]["300"] = backDhHis[bHm]["300"] + 1
			}
			if lenDltHis-i <= 500 {
				backDhHis[bHm]["500"] = backDhHis[bHm]["500"] + 1
			}
			if lenDltHis-i <= 1000 {
				backDhHis[bHm]["1000"] = backDhHis[bHm]["1000"] + 1
			}
			if lenDltHis-i <= 1500 {
				backDhHis[bHm]["1500"] = backDhHis[bHm]["1500"] + 1
			}
			if lenDltHis-i <= 2000 {
				backDhHis[bHm]["2000"] = backDhHis[bHm]["2000"] + 1
			}
			if lenDltHis-i <= 2500 {
				backDhHis[bHm]["2500"] = backDhHis[bHm]["2500"] + 1
			}
			if lenDltHis-i <= 3500 {
				backDhHis[bHm]["3500"] = backDhHis[bHm]["3500"] + 1
			}
		}
	}
	for k, v := range quShi2St {
		quShi2St[k].HadNotExistCombs = gen.DiffSlice(v.AllCombs, v.HadExistCombs)
	}
	return
}

func CalOneFrontHmBelongToTxs(frontHms []string, t2MoniABCDEs map[string]map[string][]string) (res []string) {
	found := false

	for tx, moni := range t2MoniABCDEs {
		curFound := CheckElementsFromDiffABCDE(frontHms, [][]string{moni["A"], moni["B"], moni["C"], moni["D"], moni["E"]})

		if curFound {
			res = append(res, tx)
		}
		if curFound && !found {
			found = true
		}
	}
	if !found {
		return []string{"OtherT"}
	}
	return res
}

func CheckElementsFromDiffABCDE(frontHms []string, moniABCDE [][]string) bool {
	if len(moniABCDE) != 5 {
		return false
	}

	// 将 moniABCDE 转为 map
	targetSets := make([]map[string]struct{}, 5)
	for i := 0; i < 5; i++ {
		set := make(map[string]struct{}, len(moniABCDE[i]))
		for _, v := range moniABCDE[i] {
			set[v] = struct{}{}
		}
		targetSets[i] = set
	}

	// 标记数组是否已被使用
	usedArrays := make([]bool, 5)

	// 遍历 frontHms
	for _, element := range frontHms {
		found := false

		// 为当前元素寻找一个未使用且包含该元素的数组
		for i := 0; i < 5; i++ {
			if !usedArrays[i] {
				if _, ok := targetSets[i][element]; ok {
					usedArrays[i] = true
					found = true
					break
				}
			}
		}

		// 当前元素无法匹配任何未使用数组
		if !found {
			return false
		}
	}

	// 验证所有数组都已被使用
	for _, used := range usedArrays {
		if !used {
			return false
		}
	}

	return true
}
