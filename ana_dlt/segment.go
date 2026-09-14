package ana_dlt

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// AllSegmentTypSli 所有段类型
var AllSegmentTypSli = []string{
	"11111",
	"2111",
	"221",
	"311",
	"32",
	"41",
	"5",
}

// // segments 是按历史出现次数分组的号码切片（次数从高到低）。
// // 每组内号码均为字符串，如 "01"、"13" 等，确保排序时字典序与数值序一致。
//
//	var segments = [][]string{
//		{"13"},                               // 出现12次
//		{"09"},                               // 出现10次
//		{"06", "21", "22", "26"},             // 出现8次
//		{"02", "03", "08", "10", "23"},       // 出现7次
//		{"04", "12", "14", "16", "18", "28"}, // 出现6次
//		{"01", "05", "07", "11", "27", "29", "32", "33", "34", "35"}, // 出现5次
//		{"24"},                   // 出现4次
//		{"17", "19", "30", "31"}, // 出现3次
//		{"15"},                   // 出现2次
//		{"20", "25"},             // 出现1次
//	}
//
// // Generate21111 根据规则2111（一个段出2个号，三个不同段各出1个号）生成所有5码组合。
// // 双号段必须来自号码数量 ≥ 2 的段，且所有4个被选段互不相同。
//
//	func Generate21111() [][]string {
//		var result [][]string
//
//		// 遍历所有可能作为“双号段”的段（号码个数 ≥ 2）
//		for dualIdx, dualSeg := range segments {
//			if len(dualSeg) < 2 {
//				continue
//			}
//
//			// 从双号段中任选2个号码的所有组合
//			dualCombos := combinations2(dualSeg)
//			for _, dualPair := range dualCombos {
//				// 构建剩余段列表（排除当前双号段）
//				var restSegs [][]string
//				for i, seg := range segments {
//					if i != dualIdx {
//						restSegs = append(restSegs, seg)
//					}
//				}
//
//				// 从剩余9个段中选3个段作为单号段（修正处）
//				segIndexCombos := combinationsK(len(restSegs), 3)
//				for _, idxs := range segIndexCombos {
//					var pickedSegs [][]string
//					for _, idx := range idxs {
//						pickedSegs = append(pickedSegs, restSegs[idx])
//					}
//					// 笛卡尔积：3个单号段各取1个号码
//					cartesian := cartesianProduct(pickedSegs)
//
//					for _, singleNums := range cartesian {
//						comb := make([]string, 0, 5)
//						comb = append(comb, dualPair...)   // 2个号码
//						comb = append(comb, singleNums...) // 3个号码
//						slices.Sort(comb)
//						result = append(result, comb)
//					}
//				}
//			}
//		}
//
//		return result
//	}
//
// // combinations2 从切片中任选2个元素的所有组合。
//
//	func combinations2(arr []string) [][]string {
//		var res [][]string
//		n := len(arr)
//		for i := 0; i < n-1; i++ {
//			for j := i + 1; j < n; j++ {
//				res = append(res, []string{arr[i], arr[j]})
//			}
//		}
//		return res
//	}
//
// // combinationsK 返回从 0..n-1 中选 k 个索引的所有组合。
//
//	func combinationsK(n, k int) [][]int {
//		var res [][]int
//		var comb []int
//		var backtrack func(start int)
//		backtrack = func(start int) {
//			if len(comb) == k {
//				tmp := make([]int, k)
//				copy(tmp, comb)
//				res = append(res, tmp)
//				return
//			}
//			for i := start; i < n; i++ {
//				comb = append(comb, i)
//				backtrack(i + 1)
//				comb = comb[:len(comb)-1]
//			}
//		}
//		backtrack(0)
//		return res
//	}
//
// // cartesianProduct 计算多个切片之间的笛卡尔积，每个切片选一个元素。
//
//	func cartesianProduct(sets [][]string) [][]string {
//		if len(sets) == 0 {
//			return [][]string{{}}
//		}
//		var res [][]string
//		var helper func(idx int, current []string)
//		helper = func(idx int, current []string) {
//			if idx == len(sets) {
//				tmp := make([]string, len(current))
//				copy(tmp, current)
//				res = append(res, tmp)
//				return
//			}
//			for _, val := range sets[idx] {
//				current = append(current, val)
//				helper(idx+1, current)
//				current = current[:len(current)-1]
//			}
//		}
//		helper(0, []string{})
//		return res
//	}
//

func ParseToSegments(dltTimesSliNumStrs []DltTimesSliNumStr) (segments [][]string) {
	for _, dltTimesSliNumStr := range dltTimesSliNumStrs {
		segments = append(segments, dltTimesSliNumStr.Sli)
	}
	return
}

func CalSegmentTyp(dltTimesSliNumStrs []DltTimesSliNumStr) string {
	var o1, t2, t3, f4, f5 int
	for _, dltTimesSliNumStr := range dltTimesSliNumStrs {
		switch len(dltTimesSliNumStr.NumStr) {
		case 1:
			o1++
		case 2:
			t2++
		case 3:
			t3++
		case 4:
			f4++
		case 5:
			f5++
		default:
		}
	}
	if o1 == 5 {
		return "11111"
	}

	if t2 == 1 && o1 == 3 {
		return "2111"
	}
	if t2 == 2 && o1 == 1 {
		return "221"
	}

	if t3 == 1 && o1 == 2 {
		return "311"
	}
	if t3 == 1 && t2 == 1 {
		return "32"
	}

	if f4 == 1 && o1 == 1 {
		return "41"
	}
	if f5 == 1 {
		return "5"
	}
	return "other"
}

// parseInput 解析输入的字符串，将其转换为按历史出现次数分段的号码组。
//
// 输入格式示例：
//
//	12*[13]=>[]
//	10*[09]=>[]
//	8*[06 21 22 26]=>[]
//	...
//
// 每一行代表一个“段”，段内号码用空格分隔。
//
// 参数:
//
//	input - 待解析的原始字符串
//
// 返回值:
//
//	一个二维字符串切片，每个内层切片是一个段的号码列表，号码已格式化为两位数字符串（例如 "01", "13"）。
func parseInput(input string) [][]string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var segments [][]string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		start := strings.Index(line, "[")
		end := strings.Index(line, "]")
		if start == -1 || end == -1 || end <= start {
			continue
		}
		numsStr := line[start+1 : end]
		numsStr = strings.TrimSpace(numsStr)
		var nums []string
		if numsStr != "" {
			parts := strings.Fields(numsStr)
			for _, p := range parts {
				n, err := strconv.Atoi(p)
				if err == nil {
					nums = append(nums, fmt.Sprintf("%02d", n))
				}
			}
		}
		segments = append(segments, nums)
	}
	return segments
}

// parseRule 将规则字符串解析为整数切片，每个数字代表需要从同一个段中选取的号码个数。
//
// 例如: "2111" 解析为 []int{2,1,1,1}，表示一个段取2个号码，三个不同的段各取1个号码。
//
// 参数:
//
//	rule - 由数字1-9组成的规则字符串
//
// 返回值:
//
//	规则对应的整数切片，以及可能的错误（如包含非数字字符）。
func parseRule(rule string) ([]int, error) {
	var pattern []int
	for _, ch := range rule {
		if ch < '1' || ch > '9' {
			return nil, fmt.Errorf("rule must contain digits 1-9, got %q", rule)
		}
		pattern = append(pattern, int(ch-'0'))
	}
	return pattern, nil
}

// GenerateCombinationsByInputStrAndRule 根据输入的段数据字符串和规则字符串，生成所有符合规则的5码组合。
//
// 规则说明: 例如规则 "2111" 表示从不同的4个段中选号，其中一个段选2个号码，其余三个段各选1个号码，
// 所有号码最终按从小到大排序。
//
// 参数:
//
//	input - 原始数据字符串（格式见 parseInput 说明）
//	rule  - 规则字符串，例如 "2111", "221", "32", "11111", "41", "5" 等
//
// 返回值:
//
//	所有可能的5个号码组合（已排序的字符串切片），以及可能的错误。
func GenerateCombinationsByInputStrAndRule(input string, rule string) ([][]string, error) {
	pattern, err := parseRule(rule)
	if err != nil {
		return nil, err
	}
	segments := parseInput(input)
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments parsed")
	}
	return generateFromSegments(segments, pattern), nil
}

// GenerateCombinationsBySegmentsAndRule 根据输入的段数据二维切片字符串和规则字符串，生成所有符合规则的5码组合。
//
//	@Description:
//	@param segments 段数据二维切片字符串
//	@param rule 规则字符串
//	@return [][]string
//	@return error
func GenerateCombinationsBySegmentsAndRule(segments [][]string, rule string) ([][]string, error) {
	pattern, err := parseRule(rule)
	if err != nil {
		return nil, err
	}
	return generateFromSegments(segments, pattern), nil
}

// generateFromSegments 根据已经解析好的段数据和规则模式，生成所有可能的组合。
// 此函数为内部实现，不暴露给外部。
func generateFromSegments(segments [][]string, pattern []int) [][]string {
	var result [][]string

	// 对模式排序，用于生成所有唯一的排列（例如 2,1,1,1）
	sort.Ints(pattern)
	perms := uniquePermutations(pattern)

	k := len(pattern)
	// 从所有段中选取 k 个不同的段
	segIdxCombos := chooseIndices(len(segments), k)
	//fmt.Printf("segIdxCombos: %+v\n", segIdxCombos)

	for _, idxs := range segIdxCombos {
		// 当前选中的段
		selected := make([][]string, k)
		for i, idx := range idxs {
			selected[i] = segments[idx]
		}

		for _, perm := range perms {
			// 检查当前排列下每个段的号码数量是否足够
			valid := true
			for i, need := range perm {
				if len(selected[i]) < need {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}

			// 为每个段生成从该段中取 need 个号码的所有组合
			var segCombos [][][]string
			for i, need := range perm {
				combos := pickCombos(selected[i], need)
				segCombos = append(segCombos, combos)
			}

			// 对多个段的组合求笛卡尔积，拼接成一个完整的5个号码组合
			cartesianProductOfCombos(segCombos, func(combo []string) {
				sort.Strings(combo) // 整体排序，确保从小到大
				result = append(result, combo)
			})
		}
	}

	return result
}

// CountByInputStringAndRule 计算给定规则下所有可能的5码组合总数，不生成实际组合，速度更快。
//
// 参数:
//
//	input - 原始数据字符串
//	rule  - 规则字符串
//
// 返回值:
//
//	组合总数，以及可能的错误。
func CountByInputStringAndRule(input string, rule string) (int, error) {
	pattern, err := parseRule(rule)
	if err != nil {
		return 0, err
	}
	segments := parseInput(input)
	if len(segments) == 0 {
		return 0, fmt.Errorf("no segments parsed")
	}
	return countByRuleFromSegments(segments, pattern), nil
}

// CountBySegmentsAndRule 计算给定规则下所有可能的5码组合总数，不生成实际组合，速度更快。
//
//	@Description:
//	@param segments 段数据二维切片字符串
//	@param rule 规则字符串
//	@return int
//	@return error
func CountBySegmentsAndRule(segments [][]string, rule string) (int, error) {
	pattern, err := parseRule(rule)
	if err != nil {
		return 0, err
	}

	return countByRuleFromSegments(segments, pattern), nil
}

// countByRuleFromSegments 根据段数据和模式计算组合总数，使用数学组合数避免生成所有实例。
func countByRuleFromSegments(segments [][]string, pattern []int) int {
	sort.Ints(pattern)
	perms := uniquePermutations(pattern)

	k := len(pattern)
	segIdxCombos := chooseIndices(len(segments), k)

	total := 0
	for _, idxs := range segIdxCombos {
		selected := make([][]string, k)
		for i, idx := range idxs {
			selected[i] = segments[idx]
		}

		for _, perm := range perms {
			valid := true
			product := 1
			for i, need := range perm {
				m := len(selected[i])
				if m < need {
					valid = false
					break
				}
				// 从 m 个元素中选 need 个的组合数
				product *= nChooseK(m, need)
			}
			if valid {
				total += product
			}
		}
	}
	return total
}

// nChooseK 计算组合数 C(n, k)，即从 n 个元素中选取 k 个的组合方案数。
//
// 参数:
//
//	n - 总数
//	k - 选取数
//
// 返回值:
//
//	组合数；如果 k > n 返回 0。
func nChooseK(n, k int) int {
	if k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	// 利用对称性减少计算量
	if k > n-k {
		k = n - k
	}
	res := 1
	for i := 1; i <= k; i++ {
		res = res * (n - k + i) / i
	}
	return res
}

// uniquePermutations 返回已排序整数切片的唯一全排列。
// 例如输入 [1,1,2] 将返回 [[1,1,2],[1,2,1],[2,1,1]]。
//
// 参数:
//
//	arr - 已排序的整数切片
//
// 返回值:
//
//	一个二维切片，包含所有唯一的排列。
func uniquePermutations(arr []int) [][]int {
	var res [][]int
	n := len(arr)
	used := make([]bool, n)
	var cur []int
	var backtrack func()
	backtrack = func() {
		if len(cur) == n {
			tmp := make([]int, n)
			copy(tmp, cur)
			res = append(res, tmp)
			return
		}
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			// 跳过重复元素：当前元素与前一个相同且前一个未使用，则跳过以避免重复排列
			if i > 0 && arr[i] == arr[i-1] && !used[i-1] {
				continue
			}
			used[i] = true
			cur = append(cur, arr[i])
			backtrack()
			cur = cur[:len(cur)-1]
			used[i] = false
		}
	}
	sort.Ints(arr) // 确保输入有序
	backtrack()
	return res
}

// chooseIndices 返回从 0 到 n-1 中选取 k 个索引的所有组合。
//
// 参数:
//
//	n - 索引范围 [0, n)
//	k - 选取的索引个数
//
// 返回值:
//
//	一个二维整数切片，每个内层切片是长度为 k 的索引组合，并按升序排列。
func chooseIndices(n, k int) [][]int {
	var res [][]int
	var comb []int
	var backtrack func(start int)
	backtrack = func(start int) {
		if len(comb) == k {
			tmp := make([]int, k)
			copy(tmp, comb)
			res = append(res, tmp)
			return
		}
		for i := start; i < n; i++ {
			comb = append(comb, i)
			backtrack(i + 1)
			comb = comb[:len(comb)-1]
		}
	}
	backtrack(0)
	return res
}

// pickCombos 从字符串切片 arr 中选取 k 个元素的所有组合，每个组合内部元素已按字符串升序排列。
//
// 参数:
//
//	arr - 待选取的字符串切片（要求格式化为两位数字符串，如 "01"）
//	k   - 要选取的元素个数
//
// 返回值:
//
//	所有可能的组合，每个组合为 []string。
func pickCombos(arr []string, k int) [][]string {
	if k == 0 {
		return [][]string{{}}
	}
	// 先复制并排序，确保生成的组合有序
	sorted := make([]string, len(arr))
	copy(sorted, arr)
	sort.Strings(sorted)

	var res [][]string
	n := len(sorted)
	var comb []string
	var backtrack func(start int)
	backtrack = func(start int) {
		if len(comb) == k {
			tmp := make([]string, k)
			copy(tmp, comb)
			res = append(res, tmp)
			return
		}
		for i := start; i < n; i++ {
			comb = append(comb, sorted[i])
			backtrack(i + 1)
			comb = comb[:len(comb)-1]
		}
	}
	backtrack(0)
	return res
}

// cartesianProductOfCombos 对多个组合集合求笛卡尔积，每当产生一个完整组合时调用回调函数 emit。
//
// 参数:
//
//	sets - 每个元素是一个段的全部组合（[][]string）
//	emit - 回调函数，接收一个完整的 []string 组合（由各段组合拼接而成）
func cartesianProductOfCombos(sets [][][]string, emit func([]string)) {
	if len(sets) == 0 {
		emit([]string{})
		return
	}
	var helper func(idx int, current []string)
	helper = func(idx int, current []string) {
		if idx == len(sets) {
			emit(current)
			return
		}
		for _, combo := range sets[idx] {
			// 创建新切片，长度 = len(current) + len(combo)
			newCur := make([]string, len(current)+len(combo))
			copy(newCur, current)
			copy(newCur[len(current):], combo)
			helper(idx+1, newCur)
		}
	}
	helper(0, []string{})
}
