package ana_dlt

import (
	"fmt"
	"slices"
	"sort"
	"strconv"

	"github.com/before80/lot/gen"
	"github.com/before80/lot/models"
)

func Stats(zxDlts []models.Dlt, t2MoniABCDEs map[string]map[string][]string) (
	eqHis map[string]map[string]int,
	txHis map[string]map[string]int,
	oeHis map[string]map[string]int,
	qzhHis map[string]map[string]int,
	frontDhHis map[string]map[string]int,
	backDhHis map[string]map[string]int,
	backCombHis map[string]map[string]int,
	quShi2St map[string]*DltBackChQuShi,
) {
	lenDltHis := len(zxDlts)

	// 初始化设备
	eqHis = make(map[string]map[string]int)
	for _, eq := range gen.AllDltEqs {
		eqHis[eq] = make(map[string]int)
		for _, last := range gen.LastHisSlice {
			eqHis[eq][last] = 0
		}
	}

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
		curEq := dlt.EquipmentCount
		curEqStr := ""
		if slices.Contains([]int{1, 2, 3}, curEq) {
			curEqStr = fmt.Sprintf("eq%d", curEq)
		} else {
			curEqStr = ""
		}

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

		for _, tempStr := range gen.LastHisSlice {
			tempNum, _ := strconv.Atoi(tempStr)
			if lenDltHis-i <= tempNum {
				if curEqStr != "" {
					eqHis[curEqStr][tempStr] = eqHis[curEqStr][tempStr] + 1
				}

				oeHis[curOe][tempStr] = oeHis[curOe][tempStr] + 1
				qzhHis[curQzh][tempStr] = qzhHis[curQzh][tempStr] + 1
				backCombHis[backComb][tempStr] = backCombHis[backComb][tempStr] + 1
				for _, tx := range curFrontBelongToTxs {
					txHis[tx][tempStr] = txHis[tx][tempStr] + 1
				}
			}
		}

		// 前区单号
		for _, fHm := range frontHms {
			for _, tempStr := range gen.LastHisSlice {
				tempNum, _ := strconv.Atoi(tempStr)
				if lenDltHis-i <= tempNum {
					frontDhHis[fHm][tempStr] = frontDhHis[fHm][tempStr] + 1
				}
			}

		}

		// 后区单号
		for _, bHm := range backHms {
			for _, tempStr := range gen.LastHisSlice {
				tempNum, _ := strconv.Atoi(tempStr)
				if lenDltHis-i <= tempNum {
					backDhHis[bHm][tempStr] = backDhHis[bHm][tempStr] + 1
				}
			}
		}
	}
	for k, v := range quShi2St {
		quShi2St[k].HadNotExistCombs = gen.DiffSlice(v.AllCombs, v.HadExistCombs)
	}
	return
}

func DStats(zxDlts []models.Dlt, t2MoniABCDEs map[string]map[string][]string) (
	eqHis map[string]map[string]int,
	txHis map[string]map[string]int,
	oeHis map[string]map[string]int,
	qzhHis map[string]map[string]int,
	frontDhHis map[string]map[string]int,
	backDhHis map[string]map[string]int,
	backCombHis map[string]map[string]int,
	quShi2St map[string]*DltBackChQuShi,
	txEqHis map[string]map[string]map[string]int,
	oeEqHis map[string]map[string]map[string]int,
	qzhEqHis map[string]map[string]map[string]int,
	frontDhEqHis map[string]map[string]map[string]int,
	backDhEqHis map[string]map[string]map[string]int,
	backCombEqHis map[string]map[string]map[string]int,
) {

	lenDltHis := len(zxDlts)

	// 初始化设备
	eqHis = make(map[string]map[string]int)
	txEqHis = make(map[string]map[string]map[string]int)
	oeEqHis = make(map[string]map[string]map[string]int)
	qzhEqHis = make(map[string]map[string]map[string]int)
	frontDhEqHis = make(map[string]map[string]map[string]int)
	backDhEqHis = make(map[string]map[string]map[string]int)
	backCombEqHis = make(map[string]map[string]map[string]int)
	for _, eq := range gen.AllDltEqs {
		eqHis[eq] = make(map[string]int)
		txEqHis[eq] = make(map[string]map[string]int)
		oeEqHis[eq] = make(map[string]map[string]int)
		qzhEqHis[eq] = make(map[string]map[string]int)
		frontDhEqHis[eq] = make(map[string]map[string]int)
		backDhEqHis[eq] = make(map[string]map[string]int)
		backCombEqHis[eq] = make(map[string]map[string]int)
		for _, last := range gen.LastHisSlice {
			eqHis[eq][last] = 0
		}
	}

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
		curEq := dlt.EquipmentCount
		curEqStr := ""
		if slices.Contains([]int{1, 2, 3}, curEq) {
			curEqStr = fmt.Sprintf("eq%d", curEq)
		} else {
			curEqStr = ""
		}

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

		for _, tempStr := range gen.LastHisSlice {
			tempNum, _ := strconv.Atoi(tempStr)
			if lenDltHis-i <= tempNum {
				if curEqStr != "" {
					eqHis[curEqStr][tempStr] = eqHis[curEqStr][tempStr] + 1
					if _, ok := oeEqHis[curEqStr][curOe]; !ok {
						oeEqHis[curEqStr][curOe] = make(map[string]int)
						oeEqHis[curEqStr][curOe][tempStr] = 1
					} else {
						oeEqHis[curEqStr][curOe][tempStr] = oeEqHis[curEqStr][curOe][tempStr] + 1
					}

					if _, ok := qzhEqHis[curEqStr][curQzh]; !ok {
						qzhEqHis[curEqStr][curQzh] = make(map[string]int)
						qzhEqHis[curEqStr][curQzh][tempStr] = 1
					} else {
						qzhEqHis[curEqStr][curQzh][tempStr] = qzhEqHis[curEqStr][curQzh][tempStr] + 1
					}

					if _, ok := backCombEqHis[curEqStr][backComb]; !ok {
						backCombEqHis[curEqStr][backComb] = make(map[string]int)
						backCombEqHis[curEqStr][backComb][tempStr] = 1
					} else {
						backCombEqHis[curEqStr][backComb][tempStr] = backCombEqHis[curEqStr][backComb][tempStr] + 1
					}

				}

				oeHis[curOe][tempStr] = oeHis[curOe][tempStr] + 1
				qzhHis[curQzh][tempStr] = qzhHis[curQzh][tempStr] + 1
				backCombHis[backComb][tempStr] = backCombHis[backComb][tempStr] + 1

				for _, tx := range curFrontBelongToTxs {
					txHis[tx][tempStr] = txHis[tx][tempStr] + 1

					if curEqStr != "" {
						if _, ok := txEqHis[curEqStr][tx]; !ok {
							txEqHis[curEqStr][tx] = make(map[string]int)
							txEqHis[curEqStr][tx][tempStr] = 1
						} else {
							txEqHis[curEqStr][tx][tempStr] = txEqHis[curEqStr][tx][tempStr] + 1
						}
					}
				}
			}
		}

		// 前区单号
		for _, fHm := range frontHms {
			for _, tempStr := range gen.LastHisSlice {
				tempNum, _ := strconv.Atoi(tempStr)
				if lenDltHis-i <= tempNum {
					frontDhHis[fHm][tempStr] = frontDhHis[fHm][tempStr] + 1
					if curEqStr != "" {
						if _, ok := frontDhEqHis[curEqStr][fHm]; !ok {
							frontDhEqHis[curEqStr][fHm] = make(map[string]int)
							frontDhEqHis[curEqStr][fHm][tempStr] = 1
						} else {
							frontDhEqHis[curEqStr][fHm][tempStr] = frontDhEqHis[curEqStr][fHm][tempStr] + 1
						}
					}
				}
			}
		}

		// 后区单号
		for _, bHm := range backHms {
			for _, tempStr := range gen.LastHisSlice {
				tempNum, _ := strconv.Atoi(tempStr)
				if lenDltHis-i <= tempNum {
					backDhHis[bHm][tempStr] = backDhHis[bHm][tempStr] + 1

					if curEqStr != "" {
						if _, ok := backDhEqHis[curEqStr][bHm]; !ok {
							backDhEqHis[curEqStr][bHm] = make(map[string]int)
							backDhEqHis[curEqStr][bHm][tempStr] = 1
						} else {
							backDhEqHis[curEqStr][bHm][tempStr] = backDhEqHis[curEqStr][bHm][tempStr] + 1
						}
					}
				}
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

// ---------------------------------------------------------------------------
// 分类函数（与前端 Worker 中 calDltLianHaoType / calDltFbChongHaoType / matchTeShuTypes 对齐）
// ---------------------------------------------------------------------------

// CalDltLianHaoType 判定一注前区的连号类型（互斥，一注只属一类）。
//
// 算法：
//  1. 将前区转为整数并升序；
//  2. 扫描相邻差为 1 的连续段，记录每段长度；
//  3. 只保留长度 >= 2 的段（真正的「连号段」），再按段数/长度映射到选项 key。
//
// 返回 "other" 表示不属于上述任一角标选项（例如更多复杂组合），统计时忽略不计。
func CalDltLianHaoType(front []string) string {
	nums := parseSortInts(front)
	if len(nums) == 0 {
		return "other"
	}

	// 连续段长度序列，例如前区 01,02,05,06,07 → runLens = [2, 3]
	runLens := make([]int, 0, len(nums))
	lenRun := 1
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[i-1]+1 {
			lenRun++
		} else {
			runLens = append(runLens, lenRun)
			lenRun = 1
		}
	}
	runLens = append(runLens, lenRun)

	// 只保留「连号段」（长度 >= 2）
	consec := make([]int, 0, len(runLens))
	for _, r := range runLens {
		if r >= 2 {
			consec = append(consec, r)
		}
	}
	if len(consec) == 0 {
		return "none"
	}

	// 单段连号
	if len(consec) == 1 {
		switch consec[0] {
		case 2:
			return "lh2"
		case 3:
			return "lh3"
		case 4:
			return "lh4"
		case 5:
			return "lh5"
		}
	}

	// 两段连号
	if len(consec) == 2 {
		a, b := consec[0], consec[1]
		if a > b {
			a, b = b, a
		}
		// 2+2 → 间隔 2 连号
		if a == 2 && b == 2 {
			return "lh2gap"
		}
		// 2+3 → 1 个 3 连 + 1 个 2 连
		if a == 2 && b == 3 {
			return "lh3_2"
		}
	}

	return "other"
}

// CalDltFbChongHaoType 判定前后区重号类型（互斥）。
// 比较方式与前端 getIntersection 一致：按号码字符串精确匹配（如 "03" 与 "03"）。
func CalDltFbChongHaoType(front, back []string) string {
	n := len(intersectStrings(front, back))
	switch n {
	case 0:
		return "none"
	case 1:
		return "dup1"
	case 2:
		return "dup2"
	default:
		// 正常大乐透后区只有 2 个号，理论上不会到这里
		return "other"
	}
}

// HasFrontArithGap 判断前区是否存在「n 码等差间隔」：
// 任取 n 个号升序后，相邻间隔全部相等，且公差 d∈[dMin,dMax]。
// 与前端 hasFrontArithGap 一致。
func HasFrontArithGap(front []string, n, dMin, dMax int) bool {
	nums := parseSortInts(front)
	if len(nums) < n || n < 2 {
		return false
	}
	chosen := make([]int, n)
	var dfs func(start, depth int) bool
	dfs = func(start, depth int) bool {
		if depth == n {
			d := chosen[1] - chosen[0]
			if d < dMin || d > dMax {
				return false
			}
			for i := 2; i < n; i++ {
				if chosen[i]-chosen[i-1] != d {
					return false
				}
			}
			return true
		}
		for i := start; i <= len(nums)-(n-depth); i++ {
			chosen[depth] = nums[i]
			if dfs(i+1, depth+1) {
				return true
			}
		}
		return false
	}
	return dfs(0, 0)
}

// HasFrontGap3_2_17 前区 3 号等差，公差 d∈[2,17]（如 01,03,05；01,18,35）。
func HasFrontGap3_2_17(front []string) bool {
	return HasFrontArithGap(front, 3, 2, 17)
}

// HasFrontGap4_2_11 前区 4 号等差，公差 d∈[2,11]（如 01,12,23,34）。
func HasFrontGap4_2_11(front []string) bool {
	return HasFrontArithGap(front, 4, 2, 11)
}

// HasFrontGap5_2_8 前区 5 号等差，公差 d∈[2,8]（如 01,09,17,25,33）。
func HasFrontGap5_2_8(front []string) bool {
	return HasFrontArithGap(front, 5, 2, 8)
}

// HasFbLianHaoBoth 判断前+后号码合集中，是否存在长度为 k 的连号窗口，
// 且该窗口内必须同时包含至少一个前区号、一个后区号。
//
// 算法：
//  1. 前区、后区分别建集合，合并去重后升序；
//  2. 在合集上切出极大连续段；
//  3. 在长度 >= k 的连续段上滑窗，检查每个长度为 k 的子段是否前后区都有贡献。
//
// 纯前区连号或纯后区连号都不算「前后区组成 k 连号」。
func HasFbLianHaoBoth(front, back []string, k int) bool {
	frontSet := make(map[int]struct{}, len(front))
	backSet := make(map[int]struct{}, len(back))
	allSet := make(map[int]struct{}, len(front)+len(back))
	for _, s := range front {
		n, _ := strconv.Atoi(s)
		frontSet[n] = struct{}{}
		allSet[n] = struct{}{}
	}
	for _, s := range back {
		n, _ := strconv.Atoi(s)
		backSet[n] = struct{}{}
		allSet[n] = struct{}{}
	}
	all := make([]int, 0, len(allSet))
	for n := range allSet {
		all = append(all, n)
	}
	sort.Ints(all)
	if len(all) < k {
		return false
	}

	// 遍历极大连续段，再在段内做长度为 k 的滑窗
	runStart := 0
	for i := 1; i <= len(all); i++ {
		if i == len(all) || all[i] != all[i-1]+1 {
			runLen := i - runStart
			if runLen >= k {
				for s := runStart; s+k <= i; s++ {
					hasF, hasB := false, false
					for t := s; t < s+k; t++ {
						n := all[t]
						if _, ok := frontSet[n]; ok {
							hasF = true
						}
						if _, ok := backSet[n]; ok {
							hasB = true
						}
					}
					if hasF && hasB {
						return true
					}
				}
			}
			runStart = i
		}
	}
	return false
}

// MatchDltTeShuTypes 返回一注命中的特殊号类型列表（可多选，供角标累加与前端筛选共用）。
//
// 规则：
//   - gap3 / gap4 / gap5：各自独立判断，可互相并存，也可与某一个 fbLhk 并存；
//   - fbLhk：从 k=7 往下到 2，取「最长且前后区都参与」的长度，只记一个（如最长为 5 只记 fbLh5，不记 2/3/4）；
//   - 若以上都不满足，返回 ["none"]。
//
// 角标统计时对每个命中类型各 +1。
func MatchDltTeShuTypes(front, back []string) []string {
	hit := make([]string, 0, 4)
	if HasFrontGap3_2_17(front) {
		hit = append(hit, "gap3")
	}
	if HasFrontGap4_2_11(front) {
		hit = append(hit, "gap4")
	}
	if HasFrontGap5_2_8(front) {
		hit = append(hit, "gap5")
	}

	// 从长到短找最大前后区共参与连号（含 2 连）
	maxFbLh := 0
	for k := 7; k >= 2; k-- {
		if HasFbLianHaoBoth(front, back, k) {
			maxFbLh = k
			break
		}
	}
	switch maxFbLh {
	case 2:
		hit = append(hit, "fbLh2")
	case 3:
		hit = append(hit, "fbLh3")
	case 4:
		hit = append(hit, "fbLh4")
	case 5:
		hit = append(hit, "fbLh5")
	case 6:
		hit = append(hit, "fbLh6")
	case 7:
		hit = append(hit, "fbLh7")
	}

	if len(hit) == 0 {
		hit = append(hit, "none")
	}
	return hit
}

// ---------------------------------------------------------------------------
// 历史频次统计（供模板注入 LianHaoHis / FbChongHaoHis / TeShuHis 及对应 *EqHis）
// ---------------------------------------------------------------------------

// DStatsLhFbTeShu 基于历史开奖列表，生成三组维度的全部设备统计与按设备统计。
//
// 参数 zxDlts：按期序排列的历史开奖（索引越大越新，与现有 DStats 约定一致）。
//
// 返回值说明：
//
//	lianHaoHis / fbChongHaoHis / teShuHis
//	  全部设备：选项 -> 最近N期 -> 出现次数
//	lianHaoEqHis / fbChongHaoEqHis / teShuEqHis
//	  按设备：eq1|eq2|eq3 -> 选项 -> 最近N期 -> 出现次数
//
// 计数窗口：对每期 i，若「从该期到最新」的跨度 len(zxDlts)-i <= N，
// 则该期计入对应 N 的桶（与 DStats 中 LastHisSlice 逻辑一致）。
//
// 注意：
//   - 连号/重号为互斥分类，一注只给一个 key +1；类型为 other 时不计入；
//   - 特殊号可多值，一注对 MatchDltTeShuTypes 返回的每个 key 各 +1；
//   - EquipmentCount 非 1/2/3 时不写入 *EqHis，但仍写入全部设备 *His。
func DStatsLhFbTeShu(zxDlts []models.Dlt) (
	lianHaoHis map[string]map[string]int,
	fbChongHaoHis map[string]map[string]int,
	teShuHis map[string]map[string]int,
	lianHaoEqHis map[string]map[string]map[string]int,
	fbChongHaoEqHis map[string]map[string]map[string]int,
	teShuEqHis map[string]map[string]map[string]int,
) {
	lenDltHis := len(zxDlts)

	// 先把全部选项 × 全部期数窗口初始化为 0，保证模板 range 时字段齐全
	lianHaoHis = initOptionHis(gen.AllDltLianHaos)
	fbChongHaoHis = initOptionHis(gen.AllDltFbChongHaos)
	teShuHis = initOptionHis(gen.AllDltTeShus)

	// 按设备外层只预建 eq1/eq2/eq3；内层选项在首次命中时惰性创建
	lianHaoEqHis = initEqOptionHis()
	fbChongHaoEqHis = initEqOptionHis()
	teShuEqHis = initEqOptionHis()

	for i, dlt := range zxDlts {
		frontHms := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}
		backHms := []string{dlt.B1, dlt.B2}

		// 本期分类结果
		curLh := CalDltLianHaoType(frontHms)
		curFb := CalDltFbChongHaoType(frontHms, backHms)
		curTeShus := MatchDltTeShuTypes(frontHms, backHms)

		// 设备 key：1/2/3 → "eq1"/"eq2"/"eq3"；其它空串表示不计入设备角标
		curEqStr := ""
		if slices.Contains([]int{1, 2, 3}, dlt.EquipmentCount) {
			curEqStr = fmt.Sprintf("eq%d", dlt.EquipmentCount)
		}

		for _, tempStr := range gen.LastHisSlice {
			tempNum, _ := strconv.Atoi(tempStr)
			// 不在「最近 tempNum 期」窗口内则跳过
			if lenDltHis-i > tempNum {
				continue
			}

			// 17. 前区连号：互斥；other 不在初始化 key 中，自动跳过
			if _, ok := lianHaoHis[curLh]; ok {
				lianHaoHis[curLh][tempStr]++
				if curEqStr != "" {
					incrEqOption(lianHaoEqHis, curEqStr, curLh, tempStr)
				}
			}

			// 18. 前后区重号：互斥
			if _, ok := fbChongHaoHis[curFb]; ok {
				fbChongHaoHis[curFb][tempStr]++
				if curEqStr != "" {
					incrEqOption(fbChongHaoEqHis, curEqStr, curFb, tempStr)
				}
			}

			// 19. 前后区特殊号：可多选，每个命中类型各 +1
			for _, ts := range curTeShus {
				if _, ok := teShuHis[ts]; !ok {
					continue
				}
				teShuHis[ts][tempStr]++
				if curEqStr != "" {
					incrEqOption(teShuEqHis, curEqStr, ts, tempStr)
				}
			}
		}
	}

	return
}

// ---------------------------------------------------------------------------
// 内部工具
// ---------------------------------------------------------------------------

// initOptionHis 初始化「选项 -> 最近N期 -> 0」的二维统计表。
func initOptionHis(options []string) map[string]map[string]int {
	his := make(map[string]map[string]int, len(options))
	for _, opt := range options {
		his[opt] = make(map[string]int, len(gen.LastHisSlice))
		for _, last := range gen.LastHisSlice {
			his[opt][last] = 0
		}
	}
	return his
}

// initEqOptionHis 初始化按设备外层 map（eq1/eq2/eq3）；内层选项按需填充。
func initEqOptionHis() map[string]map[string]map[string]int {
	eqHis := make(map[string]map[string]map[string]int, len(gen.AllDltEqs))
	for _, eq := range gen.AllDltEqs {
		eqHis[eq] = make(map[string]map[string]int)
	}
	return eqHis
}

// incrEqOption 给 *EqHis[设备][选项][期数] 累加 1（内层 map 惰性创建）。
func incrEqOption(eqHis map[string]map[string]map[string]int, eq, option, last string) {
	if _, ok := eqHis[eq]; !ok {
		eqHis[eq] = make(map[string]map[string]int)
	}
	if _, ok := eqHis[eq][option]; !ok {
		eqHis[eq][option] = make(map[string]int)
	}
	eqHis[eq][option][last]++
}

// parseSortInts 将号码字符串解析为 int 并升序（用于连号/间隔判断）。
func parseSortInts(ss []string) []int {
	nums := make([]int, 0, len(ss))
	for _, s := range ss {
		n, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		nums = append(nums, n)
	}
	sort.Ints(nums)
	return nums
}

// intersectStrings 求两串号码的交集（按字符串精确相等，保持 a 中出现顺序）。
func intersectStrings(a, b []string) []string {
	set := make(map[string]struct{}, len(b))
	for _, x := range b {
		set[x] = struct{}{}
	}
	out := make([]string, 0)
	for _, x := range a {
		if _, ok := set[x]; ok {
			out = append(out, x)
		}
	}
	return out
}

// DltHzStat 某一「最近 N 期」窗口内的和值统计。
type DltHzStat struct {
	Min int     // 最小值
	Max int     // 最大值
	Avg float64 // 平均值
}

// CalDltFullHz 一注开奖的完整和值（前5+后2）。
func CalDltFullHz(dlt models.Dlt) int {
	nums := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}
	sum := 0
	for _, s := range nums {
		n, _ := strconv.Atoi(s)
		sum += n
	}
	return sum
}

type hzAcc struct {
	sum   int
	count int
	min   int
	max   int
	init  bool
}

func (a *hzAcc) add(hz int) {
	if a == nil {
		return
	}
	if !a.init {
		a.min, a.max, a.init = hz, hz, true
	} else {
		if hz < a.min {
			a.min = hz
		}
		if hz > a.max {
			a.max = hz
		}
	}
	a.sum += hz
	a.count++
}

func (a *hzAcc) toStat() *DltHzStat {
	if a == nil || a.count == 0 {
		return &DltHzStat{Min: 0, Max: 0, Avg: 0}
	}
	return &DltHzStat{
		Min: a.min,
		Max: a.max,
		Avg: float64(a.sum) / float64(a.count),
	}
}

// DStatsHz 按 LastHisSlice 各期数窗口，统计全设备与 eq1/eq2/eq3 的和值 min/max/avg。
//
// 约定：zxDlts 下标越大越新。
//   - 全设备：窗口为「时间轴上最近 N 期」（与其它 *His 一致：len-i <= N）
//   - 单设备：从新到旧只收集该设备的开奖，凑满 N 期后再算 min/max/avg
//     （不足 N 期则用已有期数统计）
func DStatsHz(zxDlts []models.Dlt) (
	hzHis map[string]*DltHzStat,
	hzEqHis map[string]map[string]*DltHzStat,
) {
	lenDltHis := len(zxDlts)
	hzHis = make(map[string]*DltHzStat, len(gen.LastHisSlice))
	hzEqHis = make(map[string]map[string]*DltHzStat, len(gen.AllDltEqs))
	for _, eq := range gen.AllDltEqs {
		hzEqHis[eq] = make(map[string]*DltHzStat, len(gen.LastHisSlice))
	}

	// ---------- 1) 全设备：最近 N 期（时间轴窗口）----------
	allAcc := make(map[string]*hzAcc, len(gen.LastHisSlice))
	for _, last := range gen.LastHisSlice {
		allAcc[last] = &hzAcc{}
	}
	for i, dlt := range zxDlts {
		hz := CalDltFullHz(dlt)
		for _, tempStr := range gen.LastHisSlice {
			tempNum, _ := strconv.Atoi(tempStr)
			if lenDltHis-i > tempNum {
				continue
			}
			allAcc[tempStr].add(hz)
		}
	}
	for _, last := range gen.LastHisSlice {
		hzHis[last] = allAcc[last].toStat()
	}

	// ---------- 2) 单设备：该设备自身最近 N 期 ----------
	// 先按设备收集从新到旧的和值序列
	eqHzSeries := make(map[string][]int, len(gen.AllDltEqs))
	for _, eq := range gen.AllDltEqs {
		eqHzSeries[eq] = make([]int, 0)
	}
	for i := lenDltHis - 1; i >= 0; i-- {
		dlt := zxDlts[i]
		if !slices.Contains([]int{1, 2, 3}, dlt.EquipmentCount) {
			continue
		}
		eq := fmt.Sprintf("eq%d", dlt.EquipmentCount)
		eqHzSeries[eq] = append(eqHzSeries[eq], CalDltFullHz(dlt))
	}

	for _, eq := range gen.AllDltEqs {
		series := eqHzSeries[eq] // 下标 0 为该设备最新一期
		for _, tempStr := range gen.LastHisSlice {
			n, _ := strconv.Atoi(tempStr)
			acc := &hzAcc{}
			limit := n
			if limit > len(series) {
				limit = len(series)
			}
			for j := 0; j < limit; j++ {
				acc.add(series[j])
			}
			hzEqHis[eq][tempStr] = acc.toStat()
		}
	}

	return
}
