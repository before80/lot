package ana

import (
	"fmt"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type DrawNumCombNum struct {
	DrawNum   string
	CombCount int
	Combs     []string // 包括该期以及之前所有期号的组合
}

func DeDuplicateKeepOrder(strs []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(strs))
	for _, s := range strs {
		if _, exists := seen[s]; !exists {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}

	return result[:len(result)]
}

func CalAllComb() (combs []string) {
	var afs []string
	for i := 1; i <= 35; i++ {
		iStr := strconv.Itoa(i)
		if i < 10 {
			iStr = "0" + iStr
		}
		afs = append(afs, iStr)
	}
	var abs []string
	for i := 1; i <= 12; i++ {
		iStr := strconv.Itoa(i)
		if i < 10 {
			iStr = "0" + iStr
		}
		abs = append(abs, iStr)
	}
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 2; j++ {
			if (i == 0 && j == 0) || (i >= 4) {
				continue
			}
			combs = append(combs, GenCrossComb2(afs, abs, i, j)...)
		}
	}
	return
}

func CalCombCount() int {
	var combs []string
	var afs []string
	for i := 1; i <= 35; i++ {
		iStr := strconv.Itoa(i)
		if i < 10 {
			iStr = "0" + iStr
		}
		afs = append(afs, iStr)
	}
	var abs []string
	for i := 1; i <= 12; i++ {
		iStr := strconv.Itoa(i)
		if i < 10 {
			iStr = "0" + iStr
		}
		abs = append(abs, iStr)
	}
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 2; j++ {
			if (i == 0 && j == 0) || (i >= 4) {
				continue
			}
			combs = append(combs, GenCrossComb2(afs, abs, i, j)...)
		}
	}
	return len(combs)
}

func WriteFBCombToFile() {
	startTime := time.Now()
	var combs []string
	var afs []string
	for i := 1; i <= 35; i++ {
		iStr := strconv.Itoa(i)
		if i < 10 {
			iStr = "0" + iStr
		}
		afs = append(afs, iStr)
	}
	var abs []string
	for i := 1; i <= 12; i++ {
		iStr := strconv.Itoa(i)
		if i < 10 {
			iStr = "0" + iStr
		}
		abs = append(abs, iStr)
	}

	var f0b1, f0b2, f1b0, f1b1, f1b2, f2b0, f2b1, f2b2, f3b0, f3b1, f3b2 []string
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 2; j++ {
			if (i == 0 && j == 0) || (i >= 4) {
				continue
			}
			combs = append(combs, GenCrossComb2(afs, abs, i, j)...)
			if i == 0 && j == 1 {
				f0b1 = append(f0b1, GenCrossComb2(afs, abs, i, j)...)
			}
			if i == 0 && j == 2 {
				f0b2 = append(f0b2, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 1 && j == 0 {
				f1b0 = append(f1b0, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 1 && j == 1 {
				f1b1 = append(f1b1, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 1 && j == 2 {
				f1b2 = append(f1b2, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 2 && j == 0 {
				f2b0 = append(f2b0, GenCrossComb2(afs, abs, i, j)...)
			}
			if i == 2 && j == 1 {
				f2b1 = append(f2b1, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 2 && j == 2 {
				f2b2 = append(f2b2, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 3 && j == 0 {
				f3b0 = append(f3b0, GenCrossComb2(afs, abs, i, j)...)
			}
			if i == 3 && j == 1 {
				f3b1 = append(f3b1, GenCrossComb2(afs, abs, i, j)...)
			}
			if i == 3 && j == 2 {
				f3b2 = append(f3b2, GenCrossComb2(afs, abs, i, j)...)
			}
		}
	}

	lg.InfoToFileAndStdOut(fmt.Sprintf("用时1：%v\n", time.Since(startTime)))
	dstDir := "fbcombs"
	_ = os.MkdirAll(dstDir, os.ModePerm)
	f0b1F, _ := os.OpenFile(dstDir+"/f0b1.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f0b1F.Close()
	_, _ = f0b1F.WriteString(strings.Join(f0b1, "\n"))
	f0b2F, _ := os.OpenFile(dstDir+"/f0b2.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f0b2F.Close()
	_, _ = f0b2F.WriteString(strings.Join(f0b2, "\n"))
	f1b0F, _ := os.OpenFile(dstDir+"/f1b0.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f1b0F.Close()
	_, _ = f1b0F.WriteString(strings.Join(f1b0, "\n"))
	f1b1F, _ := os.OpenFile(dstDir+"/f1b1.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f1b1F.Close()
	_, _ = f1b1F.WriteString(strings.Join(f1b1, "\n"))
	f1b2F, _ := os.OpenFile(dstDir+"/f1b2.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f1b2F.Close()
	_, _ = f1b2F.WriteString(strings.Join(f1b2, "\n"))
	f2b0F, _ := os.OpenFile(dstDir+"/f2b0.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f2b0F.Close()
	_, _ = f2b0F.WriteString(strings.Join(f2b0, "\n"))
	f2b1F, _ := os.OpenFile(dstDir+"/f2b1.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f2b1F.Close()
	_, _ = f2b1F.WriteString(strings.Join(f2b1, "\n"))
	f2b2F, _ := os.OpenFile(dstDir+"/f2b2.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f2b2F.Close()
	_, _ = f2b2F.WriteString(strings.Join(f2b2, "\n"))
	f3b0F, _ := os.OpenFile(dstDir+"/f3b0.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f3b0F.Close()
	_, _ = f3b0F.WriteString(strings.Join(f3b0, "\n"))
	f3b1F, _ := os.OpenFile(dstDir+"/f3b1.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f3b1F.Close()
	_, _ = f3b1F.WriteString(strings.Join(f3b1, "\n"))
	f3b2F, _ := os.OpenFile(dstDir+"/f3b2.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f3b2F.Close()
	_, _ = f3b2F.WriteString(strings.Join(f3b2, "\n"))

	allCombF, _ := os.OpenFile(dstDir+"/allcomb.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer allCombF.Close()
	_, _ = allCombF.WriteString(strings.Join(combs, "\n"))
	lg.InfoToFileAndStdOut(fmt.Sprintf("用时2：%v\n", time.Since(startTime)))
}

func MaxComb(dlts []models.Dlt) (counts []*NumCount, existCombs []string, drawNum2CombNums []DrawNumCombNum) {
	num2Count := make(map[string]*NumCount)
	useEqCount := CountGtEq11001(dlts)
	tempNum := 0
	preCombs := make([]string, 0, 37)
	for _, dlt := range dlts {
		drawNum := dlt.DrawNum
		fs := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}
		bs := []string{dlt.B1, dlt.B2}
		itemCombs := make([]string, 0, 37)
		for i := 0; i <= 5; i++ {
			for j := 0; j <= 2; j++ {
				if (i == 0 && j == 0) || (i >= 4) {
					continue
				}
				itemCombs = append(itemCombs, GenCrossComb2(fs, bs, i, j)...)
			}
		}
		if tempNum == 0 {
			tempNum++
			preCombs = itemCombs
			// , Combs: itemCombs
			drawNum2CombNums = append(drawNum2CombNums, DrawNumCombNum{DrawNum: drawNum, CombCount: len(itemCombs)})
		} else {
			preCombs = append(preCombs, itemCombs...)
			//Combs:= DeDuplicateKeepOrder(curCombs)
			preCombs = DeDuplicateKeepOrder(preCombs)
			drawNum2CombNums = append(drawNum2CombNums, DrawNumCombNum{DrawNum: drawNum, CombCount: len(preCombs)})
		}

		drawResult := dlt.UnSortDrawResult
		if drawResult == "" {
			drawResult = strings.Join([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}, " ")
		}
		week, _ := GetWeekday(dlt.DrawTime)

		var ok bool
		for _, num := range itemCombs {
			_, ok = num2Count[num]
			if ok {
				num2Count[num].Count = num2Count[num].Count + 1
				num2Count[num].His = append(num2Count[num].His, History{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult})
				switch week {
				case "Mon":
					num2Count[num].Mon += 1
				case "Wed":
					num2Count[num].Wed += 1
				case "Sat":
					num2Count[num].Sat += 1
				}
				switch dlt.EquipmentCount {
				case 1:
					num2Count[num].EqCount += 1
					num2Count[num].Eq1Count += 1
				case 2:
					num2Count[num].EqCount += 1
					num2Count[num].Eq2Count += 1
				case 3:
					num2Count[num].EqCount += 1
					num2Count[num].Eq3Count += 1
				}
			} else {
				num2Count[num] = &NumCount{
					Num: num, Count: 1, Mon: 0, Wed: 0, Sat: 0,
					His: []History{
						{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult},
					},
				}

				switch week {
				case "Mon":
					num2Count[num].Mon = 1
				case "Wed":
					num2Count[num].Wed = 1
				case "Sat":
					num2Count[num].Sat = 1
				}
				switch dlt.EquipmentCount {
				case 1:
					num2Count[num].EqCount = 1
					num2Count[num].Eq1Count = 1
				case 2:
					num2Count[num].EqCount = 1
					num2Count[num].Eq2Count = 1
				case 3:
					num2Count[num].EqCount = 1
					num2Count[num].Eq3Count = 1
				}
			}
		}
	}
	l := len(dlts)

	for _, v := range num2Count {
		existCombs = append(existCombs, v.Num)
		v.Rate = float64(v.Count) / float64(l) * 100
		v.MonRate = float64(v.Mon) / float64(useEqCount) * 100
		v.WedRate = float64(v.Wed) / float64(useEqCount) * 100
		v.SatRate = float64(v.Sat) / float64(useEqCount) * 100
		if useEqCount > 0 {
			v.Eq1Rate = float64(v.Eq1Count) / float64(useEqCount) * 100
			v.Eq2Rate = float64(v.Eq2Count) / float64(useEqCount) * 100
			v.Eq3Rate = float64(v.Eq3Count) / float64(useEqCount) * 100
		}
		counts = append(counts, v)
	}

	slices.SortFunc(counts, func(a, b *NumCount) int {
		return b.Count - a.Count
	})
	return
}

func MaxCombCount(dlts []models.Dlt) (drawNum2CombNums []DrawNumCombNum) {
	tempNum := 0
	preCombs := make([]string, 0, 37)
	for _, dlt := range dlts {
		drawNum := dlt.DrawNum
		fs := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}
		bs := []string{dlt.B1, dlt.B2}
		itemCombs := make([]string, 0, 37)
		for i := 0; i <= 5; i++ {
			for j := 0; j <= 2; j++ {
				if (i == 0 && j == 0) || (i >= 4) {
					continue
				}
				itemCombs = append(itemCombs, GenCrossComb2(fs, bs, i, j)...)
			}
		}
		if tempNum == 0 {
			tempNum++
			preCombs = itemCombs
			// , Combs: itemCombs
			drawNum2CombNums = append(drawNum2CombNums, DrawNumCombNum{DrawNum: drawNum, CombCount: len(itemCombs)})
		} else {
			preCombs = append(preCombs, itemCombs...)
			//Combs:= DeDuplicateKeepOrder(curCombs)
			preCombs = DeDuplicateKeepOrder(preCombs)
			drawNum2CombNums = append(drawNum2CombNums, DrawNumCombNum{DrawNum: drawNum, CombCount: len(preCombs)})
		}
	}

	return
}

type FBCount struct {
	FCount int
	BCount int
}

func GenFBCounts() (fbCounts []FBCount) {
	fbCounts = []FBCount{
		{0, 1},
		{0, 2},
		{1, 0},
		{1, 1},
		{1, 2},
		{2, 0},
		{2, 1},
		{2, 2},
		{3, 0},
		{3, 1},
		{3, 2},
	}
	return
}

type FbCombSt struct {
	FbCount int
	CombMap map[string]uint8
}

func GenFbCombMap(mapVal uint8) (fbType2FbCombSts map[string]*FbCombSt) {
	startTime := time.Now()
	var combs []string
	var afs []string
	for i := 1; i <= 35; i++ {
		iStr := strconv.Itoa(i)
		if i < 10 {
			iStr = "0" + iStr
		}
		afs = append(afs, iStr)
	}
	var abs []string
	for i := 1; i <= 12; i++ {
		iStr := strconv.Itoa(i)
		if i < 10 {
			iStr = "0" + iStr
		}
		abs = append(abs, iStr)
	}

	var f0b1, f0b2, f1b0, f1b1, f1b2, f2b0, f2b1, f2b2, f3b0, f3b1, f3b2 []string
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 2; j++ {
			if (i == 0 && j == 0) || (i >= 4) {
				continue
			}
			combs = append(combs, GenCrossComb2(afs, abs, i, j)...)
			if i == 0 && j == 1 {
				f0b1 = append(f0b1, GenCrossComb2(afs, abs, i, j)...)
			}
			if i == 0 && j == 2 {
				f0b2 = append(f0b2, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 1 && j == 0 {
				f1b0 = append(f1b0, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 1 && j == 1 {
				f1b1 = append(f1b1, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 1 && j == 2 {
				f1b2 = append(f1b2, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 2 && j == 0 {
				f2b0 = append(f2b0, GenCrossComb2(afs, abs, i, j)...)
			}
			if i == 2 && j == 1 {
				f2b1 = append(f2b1, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 2 && j == 2 {
				f2b2 = append(f2b2, GenCrossComb2(afs, abs, i, j)...)
			}

			if i == 3 && j == 0 {
				f3b0 = append(f3b0, GenCrossComb2(afs, abs, i, j)...)
			}
			if i == 3 && j == 1 {
				f3b1 = append(f3b1, GenCrossComb2(afs, abs, i, j)...)
			}
			if i == 3 && j == 2 {
				f3b2 = append(f3b2, GenCrossComb2(afs, abs, i, j)...)
			}
		}
	}

	lg.InfoToFileAndStdOut(fmt.Sprintf("用时1：%v\n", time.Since(startTime)))
	fbType2FbCombSts = make(map[string]*FbCombSt)
	allMap := make(map[string]uint8)
	for _, comb := range combs {
		allMap[comb] = 1
	}
	fbType2FbCombSts["all"] = &FbCombSt{FbCount: len(allMap), CombMap: allMap}

	f0b1Map := make(map[string]uint8)
	for _, comb := range f0b1 {
		f0b1Map[comb] = mapVal
	}
	fbType2FbCombSts["0|1"] = &FbCombSt{FbCount: len(f0b1Map), CombMap: f0b1Map}

	f0b2Map := make(map[string]uint8)
	for _, comb := range f0b2 {
		f0b2Map[comb] = mapVal
	}
	fbType2FbCombSts["0|2"] = &FbCombSt{FbCount: len(f0b2Map), CombMap: f0b2Map}

	f1b0Map := make(map[string]uint8)
	for _, comb := range f1b0 {
		f1b0Map[comb] = mapVal
	}
	fbType2FbCombSts["1|0"] = &FbCombSt{FbCount: len(f1b0Map), CombMap: f1b0Map}

	f1b1Map := make(map[string]uint8)
	for _, comb := range f1b1 {
		f1b1Map[comb] = mapVal
	}
	fbType2FbCombSts["1|1"] = &FbCombSt{FbCount: len(f1b1Map), CombMap: f1b1Map}

	f1b2Map := make(map[string]uint8)
	for _, comb := range f1b2 {
		f1b2Map[comb] = mapVal
	}
	fbType2FbCombSts["1|2"] = &FbCombSt{FbCount: len(f1b2Map), CombMap: f1b2Map}

	f2b0Map := make(map[string]uint8)
	for _, comb := range f2b0 {
		f2b0Map[comb] = mapVal
	}
	fbType2FbCombSts["2|0"] = &FbCombSt{FbCount: len(f2b0Map), CombMap: f2b0Map}

	f2b1Map := make(map[string]uint8)
	for _, comb := range f2b1 {
		f2b1Map[comb] = mapVal
	}
	fbType2FbCombSts["2|1"] = &FbCombSt{FbCount: len(f2b1Map), CombMap: f2b1Map}

	f2b2Map := make(map[string]uint8)
	for _, comb := range f2b2 {
		f2b2Map[comb] = mapVal
	}
	fbType2FbCombSts["2|2"] = &FbCombSt{FbCount: len(f2b2Map), CombMap: f2b2Map}

	f3b0Map := make(map[string]uint8)
	for _, comb := range f3b0 {
		f3b0Map[comb] = mapVal
	}
	fbType2FbCombSts["3|0"] = &FbCombSt{FbCount: len(f3b0Map), CombMap: f3b0Map}

	f3b1Map := make(map[string]uint8)
	for _, comb := range f3b1 {
		f3b1Map[comb] = mapVal
	}
	fbType2FbCombSts["3|1"] = &FbCombSt{FbCount: len(f3b1Map), CombMap: f3b1Map}

	f3b2Map := make(map[string]uint8)
	for _, comb := range f3b2 {
		f3b2Map[comb] = mapVal
	}
	fbType2FbCombSts["3|2"] = &FbCombSt{FbCount: len(f3b2Map), CombMap: f3b2Map}

	lg.InfoToFileAndStdOut(fmt.Sprintf("用时2：%v\n", time.Since(startTime)))
	return
}

type DrawNumNotExistComb struct {
	NoExistCombCount int
	NoExistCombs     []string
}

func MinComb(dlts []models.Dlt, eqN int, accumulateStartDrawNum string) (drawNum2Type2NotExistCount map[string]map[string]*DrawNumNotExistComb) {
	fbCounts := GenFBCounts()
	fbType2FbCombSts := GenFbCombMap(2)
	drawNum2Type2NotExistCount = make(map[string]map[string]*DrawNumNotExistComb)

	var lastDrawNum string
	for _, dlt := range dlts {
		if eqN == 0 {
			if dlt.DrawNum < accumulateStartDrawNum {
				continue
			}
		} else {
			if dlt.DrawNum < accumulateStartDrawNum || dlt.EquipmentCount != eqN {
				continue
			}
		}
		lastDrawNum = dlt.DrawNum

		fs := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}
		bs := []string{dlt.B1, dlt.B2}

		for _, fbCount := range fbCounts {
			itemCombs := GenCrossComb2(fs, bs, fbCount.FCount, fbCount.BCount)
			typeStr := fmt.Sprintf("%d|%d", fbCount.FCount, fbCount.BCount)
			if _, ok := drawNum2Type2NotExistCount[dlt.DrawNum]; !ok {
				drawNum2Type2NotExistCount[dlt.DrawNum] = make(map[string]*DrawNumNotExistComb)
			}

			if _, ok := drawNum2Type2NotExistCount[dlt.DrawNum][typeStr]; !ok {
				drawNum2Type2NotExistCount[dlt.DrawNum][typeStr] = &DrawNumNotExistComb{NoExistCombCount: fbType2FbCombSts[typeStr].FbCount}
			}

			if _, ok := drawNum2Type2NotExistCount[dlt.DrawNum]["all"]; !ok {
				drawNum2Type2NotExistCount[dlt.DrawNum]["all"] = &DrawNumNotExistComb{NoExistCombCount: fbType2FbCombSts["all"].FbCount}
			}

			for _, itemComb := range itemCombs {
				if v, ok := fbType2FbCombSts[typeStr].CombMap[itemComb]; ok && v == 2 {
					fbType2FbCombSts[typeStr].CombMap[itemComb] = 1
					fbType2FbCombSts["all"].CombMap[itemComb] = 1
					fbType2FbCombSts[typeStr].FbCount -= 1
					fbType2FbCombSts["all"].FbCount -= 1
					if _, ok := drawNum2Type2NotExistCount[dlt.DrawNum]; !ok {
						drawNum2Type2NotExistCount[dlt.DrawNum] = make(map[string]*DrawNumNotExistComb)
					}
					drawNum2Type2NotExistCount[dlt.DrawNum][typeStr] = &DrawNumNotExistComb{NoExistCombCount: fbType2FbCombSts[typeStr].FbCount}
					drawNum2Type2NotExistCount[dlt.DrawNum]["all"] = &DrawNumNotExistComb{NoExistCombCount: fbType2FbCombSts["all"].FbCount}
				}
			}
		}
	}

	for _, fbCount := range fbCounts {
		typeStr := fmt.Sprintf("%d|%d", fbCount.FCount, fbCount.BCount)
		if drawNum2Type2NotExistCount[lastDrawNum][typeStr].NoExistCombCount != 0 {
			for comb, combCount := range fbType2FbCombSts[typeStr].CombMap {
				if combCount == 2 {
					drawNum2Type2NotExistCount[lastDrawNum][typeStr].NoExistCombs = append(drawNum2Type2NotExistCount[lastDrawNum][typeStr].NoExistCombs, comb)
				}
			}
		}
	}
	fmt.Println(drawNum2Type2NotExistCount[lastDrawNum])
	return
}
