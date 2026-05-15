package gen

import (
	"fmt"
	"strings"
	"time"
)

// LastHisSlice 最近多少期
var LastHisSlice = []string{"10", "20", "30", "50", "100", "200", "300", "500", "1500", "2000", "2500", "3500"}

// Comb 生成指定切片中k个元素的所有组合
func Comb(strs []string, k int) []string {
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

// CrossComb 生成两个切片的拼接组合
func CrossComb(a, b []string, c, d int) []string {
	// 从a切片中获取c个元素的所有组合
	combAs := Comb(a, c)

	// 从b切片中获取d个元素的所有组合
	combBs := Comb(b, d)

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

func CrossComb2(a, b []string, c, d int) []string {
	// 从a切片中获取c个元素的所有组合
	combAs := Comb(a, c)

	// 从b切片中获取d个元素的所有组合
	combBs := Comb(b, d)

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

	// 获取星期几，并转换为三个英文字母表示的星期几
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

// DiffSlice 返回在 slice a 中存在、但在 slice b 中不存在的元素
func DiffSlice(a, b []string) []string {
	// 以下运行速度太慢,故注释掉
	//if len(b) == 0 {
	//	return a
	//}
	//
	//// 遍历 a，筛选出不在 b 中的元素
	//var diff []string
	//for _, v := range a {
	//	if !slices.Contains(b, v) {
	//		diff = append(diff, v)
	//	}
	//}
	//
	//return diff

	if len(a) == 0 {
		return nil
	}
	if len(b) == 0 {
		return a
	}

	// 将 b 转为 set
	bSet := make(map[string]struct{}, len(b))
	for _, v := range b {
		bSet[v] = struct{}{}
	}

	// 预分配容量，减少扩容
	diff := make([]string, 0, len(a))
	for _, v := range a {
		if _, ok := bSet[v]; !ok {
			diff = append(diff, v)
		}
	}

	return diff
}

// GetWeekdayCN 获取日期是星期几
//
//	@Description:
//	@param dateStr
//	@return string
//	@return error
func GetWeekdayCN(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	weekMap := []string{"日", "一", "二", "三", "四", "五", "六"}
	return weekMap[t.Weekday()], nil
}

func JudgeFrontBackCountFromStr(inputStr string) (fN, bN int) {
	sxIndex := strings.Index(inputStr, "|")
	if sxIndex == -1 {
		fN = len(strings.Split(inputStr, ","))
	} else {
		s := strings.Split(inputStr, "|")
		fN = len(strings.Split(s[0], ","))
		s2 := strings.Split(s[1], ",")
		if len(s2[0]) == 0 {
			bN = 0
		} else {
			bN = len(s2)
		}
	}
	return
}

// HasIntersection 判断是否存在交集
//
//	@Description:
//	@param a
//	@param b
//	@return bool
func HasIntersection(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	set := make(map[string]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := set[v]; ok {
			return true
		}
	}

	return false
}

func SliceIntersection(a, b []string) (res []string) {
	if len(a) == 0 || len(b) == 0 {
		return []string{}
	}

	set := make(map[string]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := set[v]; ok {
			res = append(res, v)
		}
	}

	return
}
