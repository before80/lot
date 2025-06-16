package ana

import (
	"fmt"
	"github.com/before80/lot/models"
	"slices"
	"strings"
	"time"
	_ "time/tzdata"
)

//type WeekCount struct {
//	Mon int
//	Wed int
//	Sat int
//}

type History struct {
	DrawNum    string // 期号
	DrawTime   string // 开奖时间
	DrawResult string // 历史开奖结果
}

type NumCount struct {
	Num      string    // 单个号码 / 其他组合（组合使用英文逗号隔开，若是前后区号码则使用|隔开）
	Count    int       // 出现次数
	Mon      int       // 周一出现次数
	Wed      int       // 周三出现次数
	Sat      int       // 周六出现次数
	His      []History // 出现的历史开奖结果
	Rate     float64   // 出现的比例
	MonRate  float64   // 周一出现的比例
	WedRate  float64   // 周三出现的比例
	SatRate  float64   // 周六出现的比例
	EqCount  int       // 使用所有设备（任意，包括没有设备号的）总数
	Eq1Count int       // 使用设备1的次数
	Eq2Count int       // 使用设备2的次数
	Eq3Count int       // 使用设备3的次数
	Eq1Rate  float64   // 使用设备1的比例
	Eq2Rate  float64   // 使用设备2的比例
	Eq3Rate  float64   // 使用设备3的比例
}

func CountGtEq11001(dlts []models.Dlt) (count int) {
	for _, dlt := range dlts {
		if dlt.DrawNum >= "11001" {
			count += 1
		}
	}
	return
}

// GenComb 生成指定切片中k个元素的所有组合
func GenComb(strs []string, k int) []string {
	var result []string
	n := len(strs)

	// 如果k小于1或大于切片长度，无法形成有效组合
	if k < 1 || k > n {
		return result
	}

	// 创建临时数组存储当前组合
	current := make([]string, k)

	// 递归生成组合
	var backtrack func(start, depth int)
	backtrack = func(start, depth int) {
		// 如果已经选择了k个元素，将当前组合添加到结果中
		if depth == k {
			result = append(result, strings.Join(current[:k], ","))
			return
		}

		// 从start开始选择元素
		for i := start; i < n; i++ {
			current[depth] = strs[i]
			// 递归选择下一个元素，注意下一个start是i+1，避免重复选择
			backtrack(i+1, depth+1)
		}
	}

	// 开始递归
	backtrack(0, 0)

	return result
}

// GenCrossComb 生成两个切片的交叉组合
func GenCrossComb(a, b []string, c, d int) []string {
	// 从a切片中获取c个元素的所有组合
	combAs := GenComb(a, c)

	// 从b切片中获取d个元素的所有组合
	combBs := GenComb(b, d)

	// 如果任一组合为空，直接返回空结果
	if len(combAs) == 0 && c == 0 {
		return combBs
	}

	if len(combBs) == 0 && d == 0 {
		return combAs
	}

	// 组合所有可能的结果
	var result []string
	for _, combA := range combAs {
		for _, combB := range combBs {
			// 将两个组合用"|"连接
			result = append(result, combA+"|"+combB)
		}
	}

	return result
}

func GenCrossComb2(a, b []string, c, d int) []string {
	// 从a切片中获取c个元素的所有组合
	combAs := GenComb(a, c)

	// 从b切片中获取d个元素的所有组合
	combBs := GenComb(b, d)

	// 如果任一组合为空，直接返回空结果
	if len(combAs) == 0 && c == 0 {
		newCombBs := make([]string, 0)
		for _, combB := range combBs {
			newCombBs = append(newCombBs, "|"+combB)
		}
		return newCombBs
	}

	if len(combBs) == 0 && d == 0 {
		newCombAs := make([]string, 0)
		for _, combA := range combAs {
			newCombAs = append(newCombAs, combA+"|")
		}
		return newCombAs
	}

	// 组合所有可能的结果
	var result []string
	for _, combA := range combAs {
		for _, combB := range combBs {
			// 将两个组合用"|"连接
			result = append(result, combA+"|"+combB)
		}
	}

	//fmt.Printf("result: %v\n\n\n", result)
	return result
}

// GetWeekday 获取指定日期的星期几(英文单词前三个字母)
func GetWeekday(dateStr string) (string, error) {
	// 创建北京时间时区（东八区，UTC+8）
	beijing, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 备选方案：手动创建UTC+8时区
		beijing = time.FixedZone("CST", 8*60*60)
	}

	// 在北京时间时区解析日期
	date, err := time.ParseInLocation(time.DateOnly, dateStr, beijing)
	if err != nil {
		return "", fmt.Errorf("解析日期失败: %v", err)
	}

	// 获取星期几并转换为中文
	weekday := date.Weekday().String()
	//weekdayCN := map[string]string{
	//	"Sunday":    "星期日",
	//	"Monday":    "星期一",
	//	"Tuesday":   "星期二",
	//	"Wednesday": "星期三",
	//	"Thursday":  "星期四",
	//	"Friday":    "星期五",
	//	"Saturday":  "星期六",
	//}

	//return weekdayCN[weekday], nil
	return weekday[:3], nil
}

// FMaxOne 分析前区单个号码
func FMaxOne(dlts []models.Dlt) (counts []*NumCount) {
	//return FBMaxN2M(dlts, 1, 0)
	return FMaxN(dlts, 1)
}

// FMaxTwo 分析前区两个号码
func FMaxTwo(dlts []models.Dlt) (counts []*NumCount) {
	return FMaxN(dlts, 2)
}

// FMaxThree 分析前区三个号码
func FMaxThree(dlts []models.Dlt) (counts []*NumCount) {
	return FMaxN(dlts, 3)
}

// FMaxFour 分析前区四个号码
func FMaxFour(dlts []models.Dlt) (counts []*NumCount) {
	return FMaxN(dlts, 4)
}

// FMaxFive 分析前区五个号码
func FMaxFive(dlts []models.Dlt) (counts []*NumCount) {
	return FMaxN(dlts, 5)
}

// FBMaxOne2One 分析前区单个号码和后区单个号码
func FBMaxOne2One(dlts []models.Dlt) (counts []*NumCount) {
	return FBMaxN2M(dlts, 1, 1)
}

// FBMaxOne2Two 分析前区单个号码和后区两个号码
func FBMaxOne2Two(dlts []models.Dlt) (counts []*NumCount) {
	return FBMaxN2M(dlts, 1, 2)
}

// FBMaxTwo2One 分析前区两个号码和后区单个号码
func FBMaxTwo2One(dlts []models.Dlt) (counts []*NumCount) {
	return FBMaxN2M(dlts, 2, 1)
}

// FBMaxTwo2Two 分析前区两个号码和后区两个号码
func FBMaxTwo2Two(dlts []models.Dlt) (counts []*NumCount) {
	return FBMaxN2M(dlts, 2, 2)
}

// FBMaxThree2One 分析前区三个号码和后区单个号码
func FBMaxThree2One(dlts []models.Dlt) (counts []*NumCount) {
	return FBMaxN2M(dlts, 3, 1)
}

// FBMaxThree2Two 分析前区三个号码和后区两个号码
func FBMaxThree2Two(dlts []models.Dlt) (counts []*NumCount) {
	return FBMaxN2M(dlts, 3, 2)
}

// FBMaxFour2One 分析前区四个号码和后区单个号码
func FBMaxFour2One(dlts []models.Dlt) (counts []*NumCount) {
	return FBMaxN2M(dlts, 4, 1)
}

// FBMaxFour2Two 分析前区四个号码和后区两个号码
func FBMaxFour2Two(dlts []models.Dlt) (counts []*NumCount) {
	return FBMaxN2M(dlts, 4, 2)
}

// BMaxOne 分析后区单个号码
func BMaxOne(dlts []models.Dlt) (counts []*NumCount) {
	return BMaxN(dlts, 1)
}

// BMaxTwo 分析后区2个号码
func BMaxTwo(dlts []models.Dlt) (counts []*NumCount) {
	return BMaxN(dlts, 2)
}

// FMaxN 分析前区n个号码
func FMaxN(dlts []models.Dlt, fN int) (counts []*NumCount) {
	num2Count := make(map[string]*NumCount)
	l := len(dlts)
	useEqCount := CountGtEq11001(dlts)
	var ok bool
	for _, dlt := range dlts {
		drawResult := dlt.UnSortDrawResult
		if drawResult == "" {
			drawResult = strings.Join([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}, " ")
		}
		week, _ := GetWeekday(dlt.DrawTime)
		nums := GenComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, fN)
		for _, num := range nums {
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
	for _, v := range num2Count {
		v.Rate = float64(v.Count) / float64(l) * 100
		v.MonRate = float64(v.Mon) / float64(l) * 100
		v.WedRate = float64(v.Wed) / float64(l) * 100
		v.SatRate = float64(v.Sat) / float64(l) * 100
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

// FBMaxN2M 分析前区n个号码和后区m个号码
func FBMaxN2M(dlts []models.Dlt, fN, bM int) (counts []*NumCount) {
	num2Count := make(map[string]*NumCount)
	l := len(dlts)
	useEqCount := CountGtEq11001(dlts)
	var ok bool
	for _, dlt := range dlts {
		drawResult := dlt.UnSortDrawResult
		if drawResult == "" {
			drawResult = strings.Join([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}, " ")
		}
		week, _ := GetWeekday(dlt.DrawTime)
		nums := GenCrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, fN, bM)
		for _, num := range nums {
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
	for _, v := range num2Count {
		v.Rate = float64(v.Count) / float64(l) * 100
		v.MonRate = float64(v.Mon) / float64(l) * 100
		v.WedRate = float64(v.Wed) / float64(l) * 100
		v.SatRate = float64(v.Sat) / float64(l) * 100
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

// BMaxN 分析后区n个号码
func BMaxN(dlts []models.Dlt, bN int) (counts []*NumCount) {
	num2Count := make(map[string]*NumCount)
	l := len(dlts)
	useEqCount := CountGtEq11001(dlts)
	var ok bool
	for _, dlt := range dlts {
		drawResult := dlt.UnSortDrawResult
		if drawResult == "" {
			drawResult = strings.Join([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}, " ")
		}
		week, _ := GetWeekday(dlt.DrawTime)
		nums := GenComb([]string{dlt.B1, dlt.B2}, bN)
		for _, num := range nums {
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
	for _, v := range num2Count {
		v.Rate = float64(v.Count) / float64(l) * 100
		v.MonRate = float64(v.Mon) / float64(l) * 100
		v.WedRate = float64(v.Wed) / float64(l) * 100
		v.SatRate = float64(v.Sat) / float64(l) * 100
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
