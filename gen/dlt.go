package gen

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var AllDltFrontHms = []string{
	"01", "02", "03", "04", "05", "06", "07",
	"08", "09", "10", "11", "12", "13", "14",
	"15", "16", "17", "18", "19", "20", "21",
	"22", "23", "24", "25", "26", "27", "28",
	"29", "30", "31", "32", "33", "34", "35",
}

var AllDltBackHms = []string{
	"01", "02", "03", "04", "05", "06",
	"07", "08", "09", "10", "11", "12",
}

// AllDltOes 所有大乐透奇偶可能
var AllDltOes = []string{"07", "16", "25", "34", "43", "52", "61", "70"}

// AllDltQzhs 所有大乐透前中后可能
var AllDltQzhs = []string{"500", "410", "401", "311", "302", "212", "203", "113", "104", "014", "005"}

// AllDltTxs 前区组合
var AllDltTxs = []string{
	"T1", "T2", "T3", "T4", "T5", "T6", "T7", "T8",
	"T9", "T10", "T11", "T12", "T13", "T14", "T15",
	"T16", "T17",
	"OtherT",
}

// GetDltBackQuShiHaoMasFromQuShi 获取大乐透指定后区两个号码（例如"01,02"）和趋势字符串（例如"sk"）的情况下，可能的下一期后区号码的组合的切片
//
//	@Description:
//	@param inputStr 类似'01,02' 这种格式的两个后区号码组合成的字符串
//	@param pattern 模式字符串，即 ks、sk等模式字符串
//	@return []string 下一期后区两个号码的组合字符串切片
//	@return error
func GetDltBackQuShiHaoMasFromQuShi(inputStr, pattern string) ([]string, error) {
	// 解析输入字符串
	parts := strings.Split(inputStr, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("输入字符串格式错误，应为 '01,02' 这样的格式")
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("第一个数字解析错误: %v", err)
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("第二个数字解析错误: %v", err)
	}

	if start > end {
		return nil, fmt.Errorf("第一个数字不能大于第二个数字")
	}

	var result []string

	switch pattern {
	case "ks":
		// 第一个数字往右扩，第二个数字往左扩
		for i := start + 1; i < end-1; i++ {
			for j := start + 2; j <= end-1; j++ {
				if i < j {
					result = append(result, fmt.Sprintf("%02d,%02d", i, j))
				}
			}
		}
	case "kq":
		// 第一个数字往右扩，第二个数字固定
		for i := start + 1; i < end; i++ {
			result = append(result, fmt.Sprintf("%02d,%02d", i, end))
		}
	case "kk":
		// 第一个数字往右扩，第二个数字往右扩
		for i := start + 1; i < 12; i++ {
			for j := end + 1; j <= 12; j++ {
				// 排除 kk_q的情况
				if i != end && i < j {
					result = append(result, fmt.Sprintf("%02d,%02d", i, j))
				}
			}
		}
	case "kk_q":
		// 第一个数字往右扩，第二个数字往左扩，但必须包含等于第一个数字的情况
		for j := end + 1; j <= 12; j++ {
			result = append(result, fmt.Sprintf("%02d,%02d", end, j))
		}
	case "ss":
		// 第一个数字往左扩，第二个数字往左扩
		for i := 1; i <= start-1; i++ {
			for j := 2; j <= end-1; j++ {
				// 排除 ss_q的情况
				if j != start && i < j {
					result = append(result, fmt.Sprintf("%02d,%02d", i, j))
				}
			}
		}
	case "ss_q":
		// 第一个数字往左扩，第二个数字往左扩，但必须包含等于第一个数字的情况
		for i := 1; i <= start-1; i++ {
			result = append(result, fmt.Sprintf("%02d,%02d", i, start))
		}
	case "sq":
		// 第一个数字往左扩，第二个数字固定
		for i := 1; i <= start-1; i++ {
			result = append(result, fmt.Sprintf("%02d,%02d", i, end))
		}
	case "sk":
		// 第一个数字往左扩，第二个数字往右扩
		for i := 1; i <= start-1; i++ {
			for j := end + 1; j <= 12; j++ {
				if i < j {
					result = append(result, fmt.Sprintf("%02d,%02d", i, j))
				}
			}
		}
	case "qs":
		// 第一个数字固定，第二个数字往左扩
		for j := start + 1; j <= end-1; j++ {
			result = append(result, fmt.Sprintf("%02d,%02d", start, j))
		}
	case "qq":
		// 两个数字都固定
		result = append(result, fmt.Sprintf("%02d,%02d", start, end))
	case "qk":
		// 第一个数字固定，第二个数字往右扩
		for j := end + 1; j <= 12; j++ {
			result = append(result, fmt.Sprintf("%02d,%02d", start, j))
		}
	default:
		return nil, fmt.Errorf("不支持的pattern: %s", pattern)
	}

	return result, nil
}

// GetDltAllQuShiFromSpBackComb 获取大乐透指定后区两个号码的所有可能趋势
//
//	@Description:
//	@param inputStr 类似'01,02' 这种格式的两个后区号码组合成的字符串
//	@return []string 下一期后区两个号码的所有可能趋势字符串切片，例如 []string{"ks","qs"}
//	@return error
func GetDltAllQuShiFromSpBackComb(inputStr string) ([]string, error) {
	// 解析输入字符串
	parts := strings.Split(inputStr, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("输入字符串格式错误，应为 '01,02' 这样的格式")
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("第一个数字解析错误: %v", err)
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("第二个数字解析错误: %v", err)
	}

	if start > end {
		return nil, fmt.Errorf("第一个数字不能大于第二个数字")
	}

	var patterns []string

	// 分析模式

	// ks: 第一个数字往右扩，第二个数字往左扩
	if end >= start+3 {
		patterns = append(patterns, "ks")
	}

	// kq: 第一个数字往右扩，第二个数字固定
	if start < 11 && end >= start+2 {
		patterns = append(patterns, "kq")
	}

	// kk: 第一个数字往右扩，第二个数字往右扩
	// kk_q: 在kk的基础上，第一个数字往右扩时包含等于end的情况
	if start < 11 && end < 12 {
		patterns = append(patterns, "kk")
		patterns = append(patterns, "kk_q")
	}

	// ss: 第一个数字往左扩，第二个数字往左扩
	if start > 1 && end > 3 && end >= (start+2) {
		patterns = append(patterns, "ss")
	}

	if start > 1 {
		patterns = append(patterns, "ss_q")
	}

	// sk: 第一个数字往左扩，第二个数字往右扩
	// sk_q: 在sk的基础上，第二个数字往右扩时包含等于start的情况
	if start > 1 && end < 12 {
		patterns = append(patterns, "sk")
	}

	// sq: 第一个数字往左扩，第二个数字固定
	if start > 1 {
		patterns = append(patterns, "sq")
	}

	// qs: 第一个数字固定，第二个数字往左扩
	if end >= start+2 {
		patterns = append(patterns, "qs")
	}

	// qk: 第一个数字固定，第二个数字往右扩
	if end < 12 {
		patterns = append(patterns, "qk")
	}

	// qq: 两个数字都固定
	patterns = append(patterns, "qq")

	return patterns, nil
}

// Shuffle01To35 生成 "01"~"35" 的切片并随机乱序
//
//	@Description:
//	@param seedHi
//	@param seedLo
//	@return []string
func Shuffle01To35(seedHi, seedLo uint64) []string {
	// 1. 生成 "01" 到 "35" 的切片
	nums := make([]string, 35)
	//for i := 1; i <= 35; i++ {
	//	nums[i-1] = fmt.Sprintf("%02d", i)
	//}

	//for i := 35; i >= 1; i-- {
	//	nums[i-1] = fmt.Sprintf("%02d", i)
	//}
	nums = []string{"02", "33", "06", "35", "01", "03", "27", "15", "32", "34", "25", "26", "28", "08", "12", "21", "20", "29", "19", "09", "24", "30", "10", "22", "23", "13", "17", "18", "11", "14", "16", "04", "07", "05", "31"}

	//_ = seedHi
	//_ = seedLo
	r := rand.New(rand.NewPCG(seedHi, seedLo))
	num := r.IntN(5) + 1
	// 3. 洗牌算法（Fisher–Yates 洗牌）
	for i := len(nums) - 1; i > 0; i-- {
		if i%num == 0 {
			j := r.IntN(i + 1)
			nums[i], nums[j] = nums[j], nums[i]
		}
	}
	//for i := len(nums) - 1; i > 0; i-- {
	//	j := rand.IntN(i + 1)
	//	nums[i], nums[j] = nums[j], nums[i]
	//}

	return nums
}

//var Aae = []string{"01", "10", "15", "18", "20", "30", "35"}
//var Bae = []string{"02", "12", "16", "22", "26", "27", "34"}
//var Cae = []string{"03", "04", "11", "17", "24", "29", "33"}
//var Dae = []string{"05", "06", "07", "19", "25", "31", "32"}
//var Eae = []string{"08", "09", "13", "14", "21", "28", "23"}

// Rand77777ABCDEs 将01到35这35个号码，分成ABCDE 5组，每组各7个号码 共有 7*7*7*7*7 = 16807种可能
//
//	@Description:
//	@param seedHi
//	@return abcdeMap
func Rand77777ABCDEs(seedHi uint64) (abcdeMap map[string][]string) {
	all := Shuffle01To35(seedHi, 1)
	abcdeMap = make(map[string][]string)
	a, b, c, d, e := all[0:7], all[7:14], all[14:21], all[21:28], all[28:35]
	sort.Strings(a)
	sort.Strings(b)
	sort.Strings(c)
	sort.Strings(d)
	sort.Strings(e)
	abcdeMap["A"] = a
	abcdeMap["B"] = b
	abcdeMap["C"] = c
	abcdeMap["D"] = d
	abcdeMap["E"] = e
	return
}

// Rand116666ABCDEs 将01到35这35个号码，分成ABCDE 5组，A组11个，其他组都是6个，共有 11*6*6*6*6 = 14256种可能
//
//	@Description:
//	@param seedHi
//	@return abcdeMap
func Rand116666ABCDEs(seedHi uint64) (abcdeMap map[string][]string) {
	all := Shuffle01To35(seedHi, 1)
	abcdeMap = make(map[string][]string)
	a, b, c, d, e := all[0:11], all[11:17], all[17:23], all[23:29], all[29:35]
	sort.Strings(a)
	sort.Strings(b)
	sort.Strings(c)
	sort.Strings(d)
	sort.Strings(e)
	abcdeMap["A"] = a
	abcdeMap["B"] = b
	abcdeMap["C"] = c
	abcdeMap["D"] = d
	abcdeMap["E"] = e
	return
}

// Rand155555ABCDEs 将01到35这35个号码，分成ABCDE 5组，A组15个，其他组都是5个，共有 15*5*5*5*5 = 9375种可能
//
//	@Description:
//	@param seedHi
//	@return abcdeMap
func Rand155555ABCDEs(seedHi uint64) (abcdeMap map[string][]string) {
	all := Shuffle01To35(seedHi, 1)
	abcdeMap = make(map[string][]string)
	a, b, c, d, e := all[0:15], all[15:20], all[20:25], all[25:30], all[30:35]
	sort.Strings(a)
	sort.Strings(b)
	sort.Strings(c)
	sort.Strings(d)
	sort.Strings(e)
	abcdeMap["A"] = a
	abcdeMap["B"] = b
	abcdeMap["C"] = c
	abcdeMap["D"] = d
	abcdeMap["E"] = e
	return
}

// Rand194444ABCDEs 将01到35这35个号码，分成ABCDE 5组，A组19个，其他组都是4个，共有 19*4*4*4*4 = 4864种可能
//
//	@Description:
//	@param seedHi
//	@return abcdeMap
func Rand194444ABCDEs(seedHi uint64) (abcdeMap map[string][]string) {
	all := Shuffle01To35(seedHi, 1)
	abcdeMap = make(map[string][]string)
	a, b, c, d, e := all[0:19], all[19:23], all[23:27], all[27:31], all[31:35]
	sort.Strings(a)
	sort.Strings(b)
	sort.Strings(c)
	sort.Strings(d)
	sort.Strings(e)
	abcdeMap["A"] = a
	abcdeMap["B"] = b
	abcdeMap["C"] = c
	abcdeMap["D"] = d
	abcdeMap["E"] = e
	return
}

// Rand233333ABCDEs 将01到35这35个号码，分成ABCDE 5组，A组23个，其他组都是3个，共有 23*3*3*3*3 = 1863种可能、
//
//	@Description:
//	@param seedHi
//	@return abcdeMap
func Rand233333ABCDEs(seedHi uint64) (abcdeMap map[string][]string) {
	all := Shuffle01To35(seedHi, 1)
	abcdeMap = make(map[string][]string)
	a, b, c, d, e := all[0:23], all[23:26], all[26:29], all[29:32], all[32:35]
	sort.Strings(a)
	sort.Strings(b)
	sort.Strings(c)
	sort.Strings(d)
	sort.Strings(e)
	abcdeMap["A"] = a
	abcdeMap["B"] = b
	abcdeMap["C"] = c
	abcdeMap["D"] = d
	abcdeMap["E"] = e
	return
}

// Rand215432ABCDEs 将01到35这35个号码，分成ABCDE 5组，A组21个，B组5个，C组4个，D组3个，E组2个，共有 21*5*4*3*2 = 2520种可能
//
//	@Description:
//	@param seedHi
//	@return abcdeMap
func Rand215432ABCDEs(seedHi uint64) (abcdeMap map[string][]string) {
	all := Shuffle01To35(seedHi, 1)
	abcdeMap = make(map[string][]string)

	a, b, c, d, e := all[0:21], all[21:26], all[26:30], all[30:33], all[33:35]
	sort.Strings(a)
	sort.Strings(b)
	sort.Strings(c)
	sort.Strings(d)
	sort.Strings(e)
	abcdeMap["A"] = a
	abcdeMap["B"] = b
	abcdeMap["C"] = c
	abcdeMap["D"] = d
	abcdeMap["E"] = e
	return
}

// Rand224432ABCDEs 将01到35这35个号码，分成ABCDE 5组，A组22个，B组4个，C组4个，D组3个，E组2个，共有 22*4*4*3*2 = 2212种可能
//
//	@Description:
//	@param seedHi
//	@return abcdeMap
func Rand224432ABCDEs(seedHi uint64) (abcdeMap map[string][]string) {
	all := Shuffle01To35(seedHi, 1)
	abcdeMap = make(map[string][]string)

	a, b, c, d, e := all[0:22], all[22:26], all[26:30], all[30:33], all[33:35]
	sort.Strings(a)
	sort.Strings(b)
	sort.Strings(c)
	sort.Strings(d)
	sort.Strings(e)
	abcdeMap["A"] = a
	abcdeMap["B"] = b
	abcdeMap["C"] = c
	abcdeMap["D"] = d
	abcdeMap["E"] = e
	return
}

// Rand224441ABCDEs 将01到35这35个号码，分成ABCDE 5组，A组22个，B组4个，C组4个，D组4个，E组1个，共有 22*4*4*4*1 = 1408种可能
//
//	@Description:
//	@param seedHi
//	@return abcdeMap
func Rand224441ABCDEs(seedHi uint64) (abcdeMap map[string][]string) {
	all := Shuffle01To35(seedHi, 1)
	abcdeMap = make(map[string][]string)

	a, b, c, d, e := all[0:22], all[22:26], all[26:30], all[30:34], all[34:35]
	sort.Strings(a)
	sort.Strings(b)
	sort.Strings(c)
	sort.Strings(d)
	sort.Strings(e)
	abcdeMap["A"] = a
	abcdeMap["B"] = b
	abcdeMap["C"] = c
	abcdeMap["D"] = d
	abcdeMap["E"] = e
	return
}

// Rand253322ABCDEs 将01到35这35个号码，分成ABCDE 5组，A组25个，B组3个，C组3个，D组2个，E组1个，共有 25*3*3*2*2 = 900种可能
//
//	@Description:
//	@param seedHi
//	@return abcdeMap
func Rand253322ABCDEs(seedHi uint64) (abcdeMap map[string][]string) {
	all := Shuffle01To35(seedHi, 1)
	abcdeMap = make(map[string][]string)

	a, b, c, d, e := all[0:25], all[25:28], all[28:31], all[31:33], all[33:35]
	sort.Strings(a)
	sort.Strings(b)
	sort.Strings(c)
	sort.Strings(d)
	sort.Strings(e)
	abcdeMap["A"] = a
	abcdeMap["B"] = b
	abcdeMap["C"] = c
	abcdeMap["D"] = d
	abcdeMap["E"] = e
	return
}

// Rand272222ABCDEs 将01到35这35个号码，分成ABCDE 5组，A组27个，B组2个，C组2个，D组2个，E组2个，共有 27*2*2*2*2 = 432种可能
//
//	@Description:
//	@param seedHi
//	@return abcdeMap
func Rand272222ABCDEs(seedHi uint64) (abcdeMap map[string][]string) {
	all := Shuffle01To35(seedHi, 1)
	abcdeMap = make(map[string][]string)

	a, b, c, d, e := all[0:27], all[27:29], all[29:31], all[31:33], all[33:35]
	sort.Strings(a)
	sort.Strings(b)
	sort.Strings(c)
	sort.Strings(d)
	sort.Strings(e)
	abcdeMap["A"] = a
	abcdeMap["B"] = b
	abcdeMap["C"] = c
	abcdeMap["D"] = d
	abcdeMap["E"] = e
	return
}

// GetDltFrontABCDEStrFromCustomABCDEMap 根据自定义的ABCDE map，计算出 fronts中abcdeMap["A"]、abcdeMap["B"]、abcdeMap["C"]、abcdeMap["D"]、abcdeMap["E"]各自出现的号码个数
//
//	@Description:
//	@param fronts 前区5个号码组成的字符串切片
//	@param abcdeMap
//	@return result 返回类似 5A、4A1B等结果
func GetDltFrontABCDEStrFromCustomABCDEMap(fronts []string, abcdeMap map[string][]string) (result string) {
	var Ai, Bi, Ci, Di, Ei int
	for _, front := range fronts {
		//num, _ := strconv.Atoi(front)
		if slices.Contains(abcdeMap["A"], front) {
			Ai += 1
			continue
		}
		if slices.Contains(abcdeMap["B"], front) {
			Bi += 1
			continue
		}

		if slices.Contains(abcdeMap["C"], front) {
			Ci += 1
			continue
		}

		if slices.Contains(abcdeMap["D"], front) {
			Di += 1
			continue
		}

		if slices.Contains(abcdeMap["E"], front) {
			Ei += 1
			continue
		}
	}

	if Ai > 0 {
		result += fmt.Sprintf("%dA", Ai)
	}
	if Bi > 0 {
		result += fmt.Sprintf("%dB", Bi)
	}

	if Ci > 0 {
		result += fmt.Sprintf("%dC", Ci)
	}

	if Di > 0 {
		result += fmt.Sprintf("%dD", Di)
	}

	if Ei > 0 {
		result += fmt.Sprintf("%dE", Ei)
	}

	return result
}

// GetDltFrontABCDEStrFromStandardABCDE 根据标准的ABCDE map（即1-7为A，8-14为B以此类推），计算出inputs中map["A"]、map["B"]、map["C"]、map["D"]、map["E"]各自出现的号码个数
//
//	@Description:
//	@param fronts 前区5个号码组成的字符串切片
//	@return result
func GetDltFrontABCDEStrFromStandardABCDE(fronts []string) (result string) {
	var Ai, Bi, Ci, Di, Ei int
	for _, front := range fronts {
		num, _ := strconv.Atoi(front)
		if num <= 7 {
			Ai += 1
			continue
		}
		if num <= 14 {
			Bi += 1
			continue
		}

		if num <= 21 {
			Ci += 1
			continue
		}

		if num <= 28 {
			Di += 1
			continue
		}

		if num <= 35 {
			Ei += 1
			continue
		}
	}

	if Ai > 0 {
		result += fmt.Sprintf("%dA", Ai)
	}
	if Bi > 0 {
		result += fmt.Sprintf("%dB", Bi)
	}

	if Ci > 0 {
		result += fmt.Sprintf("%dC", Ci)
	}

	if Di > 0 {
		result += fmt.Sprintf("%dD", Di)
	}

	if Ei > 0 {
		result += fmt.Sprintf("%dE", Ei)
	}

	return result
}

// GetDltFrontCustomSpNumABCDE 获取 inputs中的5个号码，
//
//	@Description:
//	@param inputs
//	@param spNum
//	@param abcdeMap
//	@param removeRepeat
//	@return result
func GetDltFrontCustomSpNumABCDE(inputs []string, spNum int, abcdeMap map[string][]string, removeRepeat bool) (result []string) {
	var Ai, Bi, Ci, Di, Ei int
	for _, input := range inputs {
		//num, _ := strconv.Atoi(input)
		if slices.Contains(abcdeMap["A"], input) {
			Ai += 1
			continue
		}
		if slices.Contains(abcdeMap["B"], input) {
			Bi += 1
			continue
		}

		if slices.Contains(abcdeMap["C"], input) {
			Ci += 1
			continue
		}

		if slices.Contains(abcdeMap["D"], input) {
			Di += 1
			continue
		}

		if slices.Contains(abcdeMap["E"], input) {
			Ei += 1
			continue
		}
	}

	var letters []rune
	if Ai > 0 {
		for i := 0; i < Ai; i++ {
			letters = append(letters, 'A')
		}
	}
	if Bi > 0 {
		for i := 0; i < Bi; i++ {
			letters = append(letters, 'B')
		}
	}

	if Ci > 0 {
		for i := 0; i < Ci; i++ {
			letters = append(letters, 'C')
		}
	}

	if Di > 0 {
		for i := 0; i < Di; i++ {
			letters = append(letters, 'D')
		}
	}

	if Ei > 0 {
		for i := 0; i < Ei; i++ {
			letters = append(letters, 'E')
		}
	}

	return GetDltAllFrontABCDEWithLettersAndNotAny(letters, spNum, removeRepeat)
}

// GetDltFrontSpNumABCDE 从前区号码中计算出spNum个ABCDE
func GetDltFrontSpNumABCDE(inputs []string, spNum int, removeRepeat bool) (result []string) {
	var Ai, Bi, Ci, Di, Ei int
	for _, input := range inputs {
		num, _ := strconv.Atoi(input)
		if num <= 7 {
			Ai += 1
			continue
		}
		if num <= 14 {
			Bi += 1
			continue
		}

		if num <= 21 {
			Ci += 1
			continue
		}

		if num <= 28 {
			Di += 1
			continue
		}

		if num <= 35 {
			Ei += 1
			continue
		}
	}

	var letters []rune
	if Ai > 0 {
		for i := 0; i < Ai; i++ {
			letters = append(letters, 'A')
		}
	}
	if Bi > 0 {
		for i := 0; i < Bi; i++ {
			letters = append(letters, 'B')
		}
	}

	if Ci > 0 {
		for i := 0; i < Ci; i++ {
			letters = append(letters, 'C')
		}
	}

	if Di > 0 {
		for i := 0; i < Di; i++ {
			letters = append(letters, 'D')
		}
	}

	if Ei > 0 {
		for i := 0; i < Ei; i++ {
			letters = append(letters, 'E')
		}
	}

	return GetDltAllFrontABCDEWithLettersAndNotAny(letters, spNum, removeRepeat)
}

// GetDltAllFrontABCDE 生成从 A,B,C,D,E 中取出 n 个的所有可能情况
//
//	@Description:
//	@param n
//	@return []string 若n=5，则返回的数据类似 []string{"5A","4A1B",“3A1B1C”...}
func GetDltAllFrontABCDE(n int) []string {
	letters := []rune{'A', 'B', 'C', 'D', 'E'}
	var results []string
	counts := make([]int, len(letters))

	var dfs func(pos int, remaining int)
	dfs = func(pos int, remaining int) {
		// 如果已经分配完所有字母
		if pos == len(letters) {
			if remaining == 0 {
				results = append(results, formatResult(letters, counts))
			}
			return
		}

		// 从剩余数的最大值往下枚举（保证“多的”优先出现）
		for i := remaining; i >= 0; i-- {
			counts[pos] = i
			dfs(pos+1, remaining-i)
		}
	}

	dfs(0, n)
	return results
}

// formatResult 把计数结果转为字符串形式，例如 "2A1C2E"
//
//	@Description:
//	@param letters
//	@param counts
//	@return string
func formatResult(letters []rune, counts []int) string {
	result := ""
	for i, c := range counts {
		// 在 ASCII 码表中 大写英文字母 排在 0~9数字的后面
		if c > 0 {
			result += fmt.Sprintf("%d%c", c, letters[i])
		}
	}
	return result
}

// GetDltAllFrontABCDEWithLettersAndNotAny 生成指定长度的所有字母计数组合，如 3A、2A1B、1A1B1C
//
//	@Description:
//	@param letters
//	@param spNum
//	@param removeRepeat
//	@return []string
func GetDltAllFrontABCDEWithLettersAndNotAny(letters []rune, spNum int, removeRepeat bool) []string {
	// 1. 统计每个字母最多能出现几次
	freq := make(map[rune]int)
	for _, ch := range letters {
		freq[ch]++
	}

	// 2. 提取并排序字母（保证输出有序）
	var chars []rune
	for ch := range freq {
		chars = append(chars, ch)
	}
	sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })

	// 3. 用回溯生成所有符合总数为 spNum 的组合
	var results []string
	var backtrack func(idx, remain int, current []string)

	backtrack = func(idx, remain int, current []string) {
		// 如果正好凑够 spNum 个
		if remain == 0 {
			results = append(results, joinStrSliceParts(current))
			return
		}

		// 遍历可选字母
		for i := idx; i < len(chars); i++ {
			ch := chars[i]
			maxCount := freq[ch]
			for cnt := 1; cnt <= maxCount && cnt <= remain; cnt++ {
				part := fmt.Sprintf("%d%c", cnt, ch)
				backtrack(i+1, remain-cnt, append(current, part))
			}
		}
	}

	backtrack(0, spNum, nil)
	if removeRepeat {
		return UniqueStrSlice(results)
	} else {
		return results
	}
}

// UniqueStrSlice 去除字符串切片中的重复项
//
//	@Description:
//	@param slice
//	@return []string
func UniqueStrSlice(slice []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}

// joinStrSliceParts 拼接 ["1A", "2B"] => "1A2B"
//
//	@Description:
//	@param parts
//	@return res
func joinStrSliceParts(parts []string) (res string) {
	for _, p := range parts {
		res += p
	}
	return res
}

// CalDltTyp1 计算给的号码中的类型1，类型有：7,6,5,4,3,2,
//
//	@Description:
//	@param input
//	@return uint8
func CalDltTyp1(input string) uint8 {
	re := regexp.MustCompile(`\|$`)
	input = re.ReplaceAllString(input, "")
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "|", ",")
	inputs := strings.Split(input, ",")
	return uint8(len(inputs))
}

// CalDltTyp2 计算给的号码中的类型2，类型有：70,60,51,42,50,41,32,40,31,22,
//
//	@Description:
//	@param input
//	@return uint8
func CalDltTyp2(input string) uint8 {
	sxIndex := strings.Index(input, "|")
	if sxIndex == -1 {
		input = strings.ReplaceAll(input, " ", "")
		inputs := strings.Split(input, ",")
		str := fmt.Sprintf("%d0", len(inputs))
		num, _ := strconv.Atoi(str)
		return uint8(num)
	} else {
		input = strings.ReplaceAll(input, " ", "")
		inputs := strings.Split(input, "|")
		//if len(inputs) != 2 {
		//	return 0
		//}
		fronts := strings.Split(inputs[0], ",")
		backs := strings.Split(inputs[1], ",")
		str := fmt.Sprintf("%d%d", len(fronts), len(backs))
		num, _ := strconv.Atoi(str)
		return uint8(num)
	}
}

// DltFrontHmStrSliceFromAeStr 生成指定ae字符串可以生成的开奖前区号码 (完整的前区号码,使用逗号隔开的字符串)
//
//	@Description:
//	@param a
//	@param b
//	@param c
//	@param d
//	@param e
//	@return result
func DltFrontHmStrSliceFromAeStr(a, b, c, d, e string) (result []string) {
	as := strings.Split(a, ",")
	bs := strings.Split(b, ",")
	cs := strings.Split(c, ",")
	ds := strings.Split(d, ",")
	es := strings.Split(e, ",")

	for _, ia := range as {
		for _, ib := range bs {
			for _, ic := range cs {
				for _, id := range ds {
					for _, ie := range es {
						s := []string{ia, ib, ic, id, ie}
						sort.Strings(s)
						result = append(result, fmt.Sprintf("%s,%s,%s,%s,%s", s[0], s[1], s[2], s[3], s[4]))
					}
				}
			}
		}
	}

	return result
}

// DltFullHmSliceFromAeStr 生成指定ae字符串可以生成的开奖全部7个号码
//
//	@Description:
//	@param a
//	@param b
//	@param c
//	@param d
//	@param e
//	@return result
func DltFullHmSliceFromAeStr(a, b, c, d, e string) (result []string) {
	as := strings.Split(a, ",")
	bs := strings.Split(b, ",")
	cs := strings.Split(c, ",")
	ds := strings.Split(d, ",")
	es := strings.Split(e, ",")
	backCombs := Comb(AllDltBackHms, 2)
	for _, ia := range as {
		for _, ib := range bs {
			for _, ic := range cs {
				for _, id := range ds {
					for _, ie := range es {
						s := []string{ia, ib, ic, id, ie}
						sort.Strings(s)
						for _, comb := range backCombs {
							result = append(result, fmt.Sprintf("%s,%s,%s,%s,%s|%s", s[0], s[1], s[2], s[3], s[4], comb))
						}
					}
				}
			}
		}
	}

	return result
}
