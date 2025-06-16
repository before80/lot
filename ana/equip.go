package ana

import (
	"github.com/before80/lot/models"
	"slices"
	"strconv"
	"strings"
)

type EqNumCount struct {
	EqNum   int       // 设备号
	Num     string    // 单个号码 / 其他组合（组合使用英文逗号隔开，若是前后区号码则使用|隔开）
	Count   int       // 出现次数
	Mon     int       // 周一出现次数
	Wed     int       // 周三出现次数
	Sat     int       // 周六出现次数
	His     []History // 出现的历史开奖结果
	Rate    float64   // 单个号码 / 其他组合 出现的比例
	MonRate float64   // 单个号码 / 其他组合 使用设备EqNum 在周一出现的比例
	WedRate float64   // 单个号码 / 其他组合 使用设备EqNum 在周三出现的比例
	SatRate float64   // 单个号码 / 其他组合 使用设备EqNum 在周六出现的比例
}

type EqAnyCount struct {
	Num     string    // 单个号码 / 其他组合（组合使用英文逗号隔开，若是前后区号码则使用|隔开）
	Count   int       // 出现次数
	Mon     int       // 周一出现次数
	Wed     int       // 周三出现次数
	Sat     int       // 周六出现次数
	His     []History // 出现的历史开奖结果
	Rate    float64   // 单个号码 / 其他组合 出现的比例
	MonRate float64   // 单个号码 / 其他组合 使用设备EqNum 在周一出现的比例
	WedRate float64   // 单个号码 / 其他组合 使用设备EqNum 在周三出现的比例
	SatRate float64   // 单个号码 / 其他组合 使用设备EqNum 在周六出现的比例
}

func EqNFMaxM(dlts []models.Dlt, n, m int) (counts []*EqNumCount) {
	num2Count := make(map[string]*EqNumCount)
	l := len(dlts)
	//useEqCount := CountGtEq11001(dlts)
	var ok bool
	for _, dlt := range dlts {
		if dlt.DrawNum < "110001" || dlt.EquipmentCount != n {
			continue
		}
		drawResult := dlt.UnSortDrawResult
		if drawResult == "" {
			drawResult = strings.Join([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}, " ")
		}
		week, _ := GetWeekday(dlt.DrawTime)
		fNums := GenComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, m)
		for _, fNum := range fNums {
			_, ok = num2Count[fNum]
			if ok {
				num2Count[fNum].Count = num2Count[fNum].Count + 1
				num2Count[fNum].His = append(num2Count[fNum].His, History{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult})
				switch week {
				case "Mon":
					num2Count[fNum].Mon += 1
				case "Wed":
					num2Count[fNum].Wed += 1
				case "Sat":
					num2Count[fNum].Sat += 1
				}
			} else {
				num2Count[fNum] = &EqNumCount{
					EqNum: n,
					Num:   fNum, Count: 1, Mon: 0, Wed: 0, Sat: 0,
					His: []History{
						{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult},
					},
				}

				switch week {
				case "Mon":
					num2Count[fNum].Mon = 1
				case "Wed":
					num2Count[fNum].Wed = 1
				case "Sat":
					num2Count[fNum].Sat = 1
				}
			}
		}

	}
	for _, v := range num2Count {
		v.Rate = float64(v.Count) / float64(l) * 100
		v.MonRate = float64(v.Mon) / float64(l) * 100
		v.WedRate = float64(v.Wed) / float64(l) * 100
		v.SatRate = float64(v.Sat) / float64(l) * 100
		counts = append(counts, v)
	}

	slices.SortFunc(counts, func(a, b *EqNumCount) int {
		return b.Count - a.Count
	})
	return
}

func EqNBMaxM(dlts []models.Dlt, n, m int) (counts []*EqNumCount) {
	num2Count := make(map[string]*EqNumCount)
	l := len(dlts)
	//useEqCount := CountGtEq11001(dlts)
	var ok bool
	for _, dlt := range dlts {
		if dlt.DrawNum < "110001" || dlt.EquipmentCount != n {
			continue
		}
		drawResult := dlt.UnSortDrawResult
		if drawResult == "" {
			drawResult = strings.Join([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}, " ")
		}
		week, _ := GetWeekday(dlt.DrawTime)
		fNums := GenComb([]string{dlt.B1, dlt.B2}, m)
		for _, fNum := range fNums {
			_, ok = num2Count[fNum]
			if ok {
				num2Count[fNum].Count = num2Count[fNum].Count + 1
				num2Count[fNum].His = append(num2Count[fNum].His, History{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult})
				switch week {
				case "Mon":
					num2Count[fNum].Mon += 1
				case "Wed":
					num2Count[fNum].Wed += 1
				case "Sat":
					num2Count[fNum].Sat += 1
				}
			} else {
				num2Count[fNum] = &EqNumCount{
					EqNum: n,
					Num:   fNum, Count: 1, Mon: 0, Wed: 0, Sat: 0,
					His: []History{
						{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult},
					},
				}

				switch week {
				case "Mon":
					num2Count[fNum].Mon = 1
				case "Wed":
					num2Count[fNum].Wed = 1
				case "Sat":
					num2Count[fNum].Sat = 1
				}
			}
		}

	}
	for _, v := range num2Count {
		v.Rate = float64(v.Count) / float64(l) * 100
		v.MonRate = float64(v.Mon) / float64(l) * 100
		v.WedRate = float64(v.Wed) / float64(l) * 100
		v.SatRate = float64(v.Sat) / float64(l) * 100
		counts = append(counts, v)
	}

	slices.SortFunc(counts, func(a, b *EqNumCount) int {
		return b.Count - a.Count
	})
	return
}

func EqNFBMaxM2K(dlts []models.Dlt, eqN, fM, bK int) (counts []*EqNumCount) {
	num2Count := make(map[string]*EqNumCount)
	l := len(dlts)
	//useEqCount := CountGtEq11001(dlts)
	var ok bool
	for _, dlt := range dlts {
		if dlt.DrawNum < "110001" || dlt.EquipmentCount != eqN {
			continue
		}
		drawResult := dlt.UnSortDrawResult
		if drawResult == "" {
			drawResult = strings.Join([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}, " ")
		}
		week, _ := GetWeekday(dlt.DrawTime)
		fNums := GenCrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, fM, bK)
		for _, fNum := range fNums {
			_, ok = num2Count[fNum]
			if ok {
				num2Count[fNum].Count = num2Count[fNum].Count + 1
				num2Count[fNum].His = append(num2Count[fNum].His, History{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult})
				switch week {
				case "Mon":
					num2Count[fNum].Mon += 1
				case "Wed":
					num2Count[fNum].Wed += 1
				case "Sat":
					num2Count[fNum].Sat += 1
				}
			} else {
				num2Count[fNum] = &EqNumCount{
					EqNum: eqN,
					Num:   fNum, Count: 1, Mon: 0, Wed: 0, Sat: 0,
					His: []History{
						{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult},
					},
				}

				switch week {
				case "Mon":
					num2Count[fNum].Mon = 1
				case "Wed":
					num2Count[fNum].Wed = 1
				case "Sat":
					num2Count[fNum].Sat = 1
				}
			}
		}

	}
	for _, v := range num2Count {
		v.Rate = float64(v.Count) / float64(l) * 100
		v.MonRate = float64(v.Mon) / float64(l) * 100
		v.WedRate = float64(v.Wed) / float64(l) * 100
		v.SatRate = float64(v.Sat) / float64(l) * 100
		counts = append(counts, v)
	}

	slices.SortFunc(counts, func(a, b *EqNumCount) int {
		return b.Count - a.Count
	})
	return
}

func EqNFBIntervalDrawNumSpNum(dlts []models.Dlt, eqN, intervalDrawNum, fbType int, spNum string, isIncrease bool) (counts []*EqNumCount, drawNums []string) {
	var fSpNums []string
	var bSpNums []string
	maxDrawNum := dlts[len(dlts)-1].DrawNum
	keyStr2DrawNum := make(map[string]string)
	useEqCount := CountGtEq11001(dlts)
	if fbType == 1 {
		fSpNums = strings.Split(spNum, ",")
	} else if fbType == 2 {
		bSpNums = strings.Split(spNum, ",")
	} else if fbType == 3 {
		spNums := strings.Split(spNum, "|")
		if spNums[0] == "" {
			fSpNums = []string{}
		} else {
			fSpNums = strings.Split(spNums[0], ",")
		}

		if spNums[1] == "" {
			bSpNums = []string{}
		} else {
			bSpNums = strings.Split(spNums[1], ",")
		}
	}

	num2Count := make(map[string]*EqNumCount)
	tempNum := 0
	//l := len(dlts)
	curKeyStr := strconv.Itoa(intervalDrawNum)
	prevKeyStr := strconv.Itoa(intervalDrawNum)
	var ok bool
	var legend string
	newSpNum := spNum
	if strings.HasPrefix(spNum, "|") {
		newSpNum = strings.ReplaceAll(spNum, "|", "")
	}
	for _, dlt := range dlts {
		if dlt.DrawNum < "11001" {
			continue
		}
		if tempNum == 0 || tempNum%intervalDrawNum == 0 {
			drawNumI, _ := strconv.Atoi(dlt.DrawNum)
			legendI := drawNumI + intervalDrawNum - 1
			legend = strconv.Itoa(legendI)
			if len(legend) == 4 {
				legend = "0" + legend
			}
			if legend > maxDrawNum {
				legend = maxDrawNum
			}
		}

		if tempNum == 0 {
			tempSt := EqNumCount{}
			tempSt.EqNum = eqN
			tempSt.Num = spNum
			tempSt.Count = 0
			num2Count[prevKeyStr] = &tempSt
			keyStr2DrawNum[prevKeyStr] = legend
		}

		if tempNum != 0 && tempNum%intervalDrawNum == 0 {
			rNum := tempNum / intervalDrawNum
			prevKeyStr = strconv.Itoa(intervalDrawNum * rNum)
			curKeyStr = strconv.Itoa(intervalDrawNum * (rNum + 1))

			_, ok = num2Count[curKeyStr]
			if !ok {
				tempSt := EqNumCount{}
				tempSt.EqNum = num2Count[prevKeyStr].EqNum
				tempSt.Num = num2Count[prevKeyStr].Num
				if isIncrease {
					tempSt.Count = num2Count[prevKeyStr].Count
					tempSt.Mon = num2Count[prevKeyStr].Mon
					tempSt.Wed = num2Count[prevKeyStr].Wed
					tempSt.Sat = num2Count[prevKeyStr].Sat
					tempSt.His = num2Count[prevKeyStr].His[:]
				}

				num2Count[curKeyStr] = &tempSt
				keyStr2DrawNum[curKeyStr] = legend
			}
		}

		tempNum++
		if dlt.EquipmentCount != eqN {
			continue
		}

		drawResult := dlt.UnSortDrawResult
		if drawResult == "" {
			drawResult = strings.Join([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}, " ")
		}
		week, _ := GetWeekday(dlt.DrawTime)
		nums := GenCrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, len(fSpNums), len(bSpNums))
		for _, num := range nums {
			if num != newSpNum {
				continue
			}

			_, ok = num2Count[curKeyStr]
			if ok {
				num2Count[curKeyStr].Count = num2Count[curKeyStr].Count + 1
				num2Count[curKeyStr].His = append(num2Count[curKeyStr].His, History{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult})
				switch week {
				case "Mon":
					num2Count[curKeyStr].Mon += 1
				case "Wed":
					num2Count[curKeyStr].Wed += 1
				case "Sat":
					num2Count[curKeyStr].Sat += 1
				}
			} else {
				num2Count[curKeyStr] = &EqNumCount{
					EqNum: eqN,
					Num:   num, Count: 1, Mon: 0, Wed: 0, Sat: 0,
					His: []History{
						{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult},
					},
				}

				switch week {
				case "Mon":
					num2Count[curKeyStr].Mon = 1
				case "Wed":
					num2Count[curKeyStr].Wed = 1
				case "Sat":
					num2Count[curKeyStr].Sat = 1
				}
			}
		}

	}
	var ks []string
	for k, _ := range num2Count {
		ks = append(ks, k)
	}

	slices.SortFunc(ks, func(a, b string) int {
		aInt, _ := strconv.Atoi(a)
		bInt, _ := strconv.Atoi(b)
		return aInt - bInt
	})
	counts = append(counts, &EqNumCount{})
	for _, k := range ks {
		v := num2Count[k]
		//lg.InfoToFileAndStdOut(fmt.Sprintf("k: %v\n", k))
		v.Rate = float64(v.Count) / float64(useEqCount) * 100
		v.MonRate = float64(v.Mon) / float64(useEqCount) * 100
		v.WedRate = float64(v.Wed) / float64(useEqCount) * 100
		v.SatRate = float64(v.Sat) / float64(useEqCount) * 100
		counts = append(counts, v)
	}
	drawNums = append(drawNums, "0")
	for _, v := range keyStr2DrawNum {
		drawNums = append(drawNums, v)
	}
	slices.Sort(drawNums)
	//slices.SortFunc(counts, func(a, b *EqNumCount) int {
	//	return b.Count - a.Count
	//})
	return
}

func EqAnyFBIntervalDrawNumSpNum(dlts []models.Dlt, intervalDrawNum, fbType int, spNum string, isIncrease bool) (counts []*EqAnyCount, drawNums []string) {
	var fSpNums []string
	var bSpNums []string
	keyStr2DrawNum := make(map[string]string)
	maxDrawNum := dlts[len(dlts)-1].DrawNum
	l := len(dlts)
	//useEqCount := CountGtEq11001(dlts)
	if fbType == 1 {
		fSpNums = strings.Split(spNum, ",")
	} else if fbType == 2 {
		bSpNums = strings.Split(spNum, ",")
	} else if fbType == 3 {
		spNums := strings.Split(spNum, "|")
		if spNums[0] == "" {
			fSpNums = []string{}
		} else {
			fSpNums = strings.Split(spNums[0], ",")
		}

		if spNums[1] == "" {
			bSpNums = []string{}
		} else {
			bSpNums = strings.Split(spNums[1], ",")
		}
	}
	num2Count := make(map[string]*EqAnyCount)
	tempNum := 0
	//l := len(dlts)
	curKeyStr := strconv.Itoa(intervalDrawNum)
	prevKeyStr := strconv.Itoa(intervalDrawNum)
	var ok bool
	var legend string
	newSpNum := spNum
	if strings.HasPrefix(spNum, "|") {
		newSpNum = strings.ReplaceAll(spNum, "|", "")
	}
	for _, dlt := range dlts {
		//if dlt.DrawNum < "11001" {
		//	continue
		//}
		if tempNum == 0 || tempNum%intervalDrawNum == 0 {
			drawNumI, _ := strconv.Atoi(dlt.DrawNum)
			legendI := drawNumI + intervalDrawNum - 1
			legend = strconv.Itoa(legendI)
			if len(legend) == 4 {
				legend = "0" + legend
			}
			if legend > maxDrawNum {
				legend = maxDrawNum
			}
		}

		if tempNum == 0 {
			tempSt := EqAnyCount{}
			tempSt.Num = spNum
			tempSt.Count = 0
			num2Count[prevKeyStr] = &tempSt
			keyStr2DrawNum[prevKeyStr] = legend
		}

		if tempNum != 0 && tempNum%intervalDrawNum == 0 {
			rNum := tempNum / intervalDrawNum
			prevKeyStr = strconv.Itoa(intervalDrawNum * rNum)
			curKeyStr = strconv.Itoa(intervalDrawNum * (rNum + 1))

			_, ok = num2Count[curKeyStr]
			if !ok {
				tempSt := EqAnyCount{}
				tempSt.Num = num2Count[prevKeyStr].Num
				if isIncrease {
					tempSt.Count = num2Count[prevKeyStr].Count
					tempSt.Mon = num2Count[prevKeyStr].Mon
					tempSt.Wed = num2Count[prevKeyStr].Wed
					tempSt.Sat = num2Count[prevKeyStr].Sat
					tempSt.His = num2Count[prevKeyStr].His[:]
				}

				num2Count[curKeyStr] = &tempSt
				keyStr2DrawNum[curKeyStr] = legend
			}
		}

		tempNum++

		drawResult := dlt.UnSortDrawResult
		if drawResult == "" {
			drawResult = strings.Join([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}, " ")
		}
		week, _ := GetWeekday(dlt.DrawTime)
		nums := GenCrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, len(fSpNums), len(bSpNums))
		//lg.InfoToFileAndStdOut(fmt.Sprintf("len(fSpNums)=%d len(bSpNums)=%d  nums: %v\n", len(fSpNums), len(bSpNums), nums))
		for _, num := range nums {
			if num != newSpNum {
				continue
			}

			_, ok = num2Count[curKeyStr]
			if ok {
				num2Count[curKeyStr].Count = num2Count[curKeyStr].Count + 1
				num2Count[curKeyStr].His = append(num2Count[curKeyStr].His, History{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult})
				switch week {
				case "Mon":
					num2Count[curKeyStr].Mon += 1
				case "Wed":
					num2Count[curKeyStr].Wed += 1
				case "Sat":
					num2Count[curKeyStr].Sat += 1
				}
			} else {
				num2Count[curKeyStr] = &EqAnyCount{
					Num: num, Count: 1, Mon: 0, Wed: 0, Sat: 0,
					His: []History{
						{DrawNum: dlt.DrawNum, DrawTime: dlt.DrawTime, DrawResult: drawResult},
					},
				}

				switch week {
				case "Mon":
					num2Count[curKeyStr].Mon = 1
				case "Wed":
					num2Count[curKeyStr].Wed = 1
				case "Sat":
					num2Count[curKeyStr].Sat = 1
				}
			}
		}

	}
	var ks []string
	for k, _ := range num2Count {
		ks = append(ks, k)
	}

	slices.SortFunc(ks, func(a, b string) int {
		aInt, _ := strconv.Atoi(a)
		bInt, _ := strconv.Atoi(b)
		return aInt - bInt
	})
	counts = append(counts, &EqAnyCount{})
	for _, k := range ks {
		v := num2Count[k]
		//lg.InfoToFileAndStdOut(fmt.Sprintf("k: %v\n", k))
		v.Rate = float64(v.Count) / float64(l) * 100
		v.MonRate = float64(v.Mon) / float64(l) * 100
		v.WedRate = float64(v.Wed) / float64(l) * 100
		v.SatRate = float64(v.Sat) / float64(l) * 100
		counts = append(counts, v)
	}

	drawNums = append(drawNums, "0")
	for _, v := range keyStr2DrawNum {
		drawNums = append(drawNums, v)
	}
	slices.Sort(drawNums)
	//slices.SortFunc(counts, func(a, b *EqNumCount) int {
	//	return b.Count - a.Count
	//})
	return
}
