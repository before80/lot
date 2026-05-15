package ana_dlt

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/before80/lot/cfg"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
)

func GenAllComb(dlt models.Dlt) (result [][]string) {
	nums := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}
	//result = Gen1Comb(nums, result)
	//result = Gen2Comb(nums, result)
	//result = Gen3Comb(nums, result)
	//result = Gen4Comb(nums, result)
	// 直接预分配确切大小的切片
	result = make([][]string, 0, 30) // C(5,1)=5  +  C(5,2)=10 + C(5,3)=10 + C(5,4)=5

	// 对于固定小数据量，直接展开循环可能更快
	n := len(nums)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			result = append(result, []string{nums[i]})
		}
	}

	// 生成2组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			result = append(result, []string{nums[i], nums[j]})
		}
	}

	// 生成3组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				result = append(result, []string{nums[i], nums[j], nums[k]})
			}
		}
	}

	// 生成4组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				for l := k + 1; l < n; l++ {
					result = append(result, []string{nums[i], nums[j], nums[k], nums[l]})
				}
			}
		}
	}

	return result
}

func GenAllCombs(dlt models.Dlt) (result []string) {
	nums := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}
	//result = Gen1Combs(nums, result)
	//result = Gen2Combs(nums, result)
	//result = Gen3Combs(nums, result)
	//result = Gen4Combs(nums, result)
	// 直接预分配确切大小的切片
	result = make([]string, 0, 30) // C(5,1)=5  + C(5,2)=10 + C(5,3)=10 + C(5,4)=5

	// 对于固定小数据量，直接展开循环可能更快
	n := len(nums)

	// 生成1组合
	for i := 0; i < n; i++ {
		result = append(result, nums[i])
	}

	// 生成2组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			result = append(result, nums[i]+"|"+nums[j])
		}
	}

	// 生成3组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				result = append(result, nums[i]+"|"+nums[j]+"|"+nums[k])
			}
		}
	}

	// 生成4组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				for l := k + 1; l < n; l++ {
					result = append(result, nums[i]+"|"+nums[j]+"|"+nums[k]+"|"+nums[l])
				}
			}
		}
	}

	return result
}

func Gen1Comb(nums []string, result [][]string) [][]string {
	n := len(nums)

	for i := 0; i < n; i++ {
		result = append(result, []string{nums[i]})
	}

	return result
}

func Gen1Combs(nums []string, result []string) []string {
	//result = nums
	//return result
	n := len(nums)

	for i := 0; i < n; i++ {
		result = append(result, nums[i])
	}

	return result
}

func Gen2Comb(nums []string, result [][]string) [][]string {
	n := len(nums)

	// 双层循环组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			result = append(result, []string{nums[i], nums[j]})
		}
	}

	return result
}

func Gen2Combs(nums []string, result []string) []string {
	n := len(nums)

	// 双层循环组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			result = append(result, nums[i]+"|"+nums[j])
		}
	}

	return result
}

func Gen3Comb(nums []string, result [][]string) [][]string {
	n := len(nums)

	// 三层循环组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				result = append(result, []string{nums[i], nums[j], nums[k]})
			}
		}
	}

	return result
}

func Gen3Combs(nums []string, result []string) []string {
	n := len(nums)

	// 三层循环组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				result = append(result, nums[i]+"|"+nums[j]+"|"+nums[k])
			}
		}
	}

	return result
}

func Gen4Comb(nums []string, result [][]string) [][]string {
	n := len(nums)

	// 四层循环组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				for l := k + 1; l < n; l++ {
					result = append(result, []string{nums[i], nums[j], nums[k], nums[l]})
				}
			}
		}
	}

	return result
}

func Gen4Combs(nums []string, result []string) []string {
	n := len(nums)

	// 四层循环组合
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				for l := k + 1; l < n; l++ {
					result = append(result, nums[i]+"|"+nums[j]+"|"+nums[k]+"|"+nums[l])
				}
			}
		}
	}

	return result
}

// JudgeHaveCombs 判断是否存在组合
func JudgeHaveCombs(aCombs []string, bCombs []string) bool {
	if len(aCombs) == 0 || len(bCombs) == 0 {
		return false
	}

	// 使用map来记录第一个切片的元素
	set := make(map[string]bool, len(aCombs))
	for _, item := range aCombs {
		set[item] = true
	}

	// 检查第二个切片中的元素是否在map中
	for _, item := range bCombs {
		if set[item] {
			return true
		}
	}

	return false
}

// JudgeHaveComb 判断是否存在特定的组合
func JudgeHaveComb(newCombs []string, comb string) bool {
	if slices.Contains(newCombs, comb) {
		return true
	}
	return false
}

func FrontGuiLi(qiHao string) {
	allDlts, _ := dbop.ReadAllDlt(false)
	waitAnaDlt, _ := dbop.ReadDlt(qiHao)
	waitAnaDltCombs := GenAllCombs(waitAnaDlt)
	tj := make(map[string]int)
	needAnaNext := false
	lg.InfoToFile(fmt.Sprintf("waitAnaDlt=%v \n", waitAnaDlt))

	for _, waitAnaDltComb := range waitAnaDltCombs {
		tj = make(map[string]int)
		needAnaNext = false
		for _, dlt := range allDlts {
			itemCombs := GenAllCombs(dlt)
			if needAnaNext {
				// 统计下一期出现的组合的各自个数
				for _, item := range itemCombs {
					if _, ok := tj[item]; !ok {
						tj[item] = 1
					} else {
						tj[item] = tj[item] + 1
					}
				}
			}

			if JudgeHaveComb(itemCombs, waitAnaDltComb) {
				needAnaNext = true
			} else {
				needAnaNext = false
			}
		}
		tjs := sortMapByValueDesc(tj)
		//lg.InfoToFile(fmt.Sprintf("%s\n", waitAnaDltComb))
		lg.InfoToFile(fmt.Sprintf("%s -> tjs=%v\n", waitAnaDltComb, tjs))
	}

}

func sortMapByValueDesc(m map[string]int) []struct {
	Key   string
	Value int
} {
	// 创建结构体切片
	pairs := make([]struct {
		Key   string
		Value int
	}, 0, len(m))

	// 填充数据
	for k, v := range m {
		pairs = append(pairs, struct {
			Key   string
			Value int
		}{k, v})
	}

	// 按值从大到小排序
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	return pairs
}

func AbcdeQs() {
	dlts, _ := dbop.ReadAllDlt(false)
	//
	aeMaps := make(map[string]map[string][]models.Dlt)
	allFrontAes := gen.GetDltAllFrontABCDE(5)
	// 5A、4A1B、3A2B等作为键，出现次数作为值
	cur5AeMaps := make(map[string]int)
	// 5A、3A1B、2A2B等作为键，出现次数作为值
	cur4AeMaps := make(map[string]int)
	// 3A、2A1B、1A2B等作为键，出现次数作为值
	cur3AeMaps := make(map[string]int)
	// 2A、1A1B等作为键，出现次数作为值
	cur2AeMaps := make(map[string]int)
	// 1A、1B等作为键，出现次数作为值
	cur1AeMaps := make(map[string]int)

	cur5QAeMaps := make(map[string]map[int]int)
	cur4QAeMaps := make(map[string]map[int]int)
	cur3QAeMaps := make(map[string]map[int]int)
	cur2QAeMaps := make(map[string]map[int]int)
	cur1QAeMaps := make(map[string]map[int]int)
	zXQs := []int{3000, 2000, 1000, 500, 200, 100, 50, 30, 20, 10}
	abcdeMap := gen.Rand77777ABCDEs(1)

	var dnf func(map[string]map[int]int, string, int, int)
	dnf = func(m map[string]map[int]int, curAe string, totalLen, index int) {
		var allKeys []int
		if totalLen-index <= 3000 {
			allKeys = append(allKeys, 3000)
		}
		if totalLen-index <= 2000 {
			allKeys = append(allKeys, 2000)
		}

		if totalLen-index <= 1000 {
			allKeys = append(allKeys, 1000)
		}

		if totalLen-index <= 500 {
			allKeys = append(allKeys, 500)
		}

		if totalLen-index <= 200 {
			allKeys = append(allKeys, 200)
		}
		if totalLen-index <= 100 {
			allKeys = append(allKeys, 100)
		}

		if totalLen-index <= 50 {
			allKeys = append(allKeys, 50)
		}

		if totalLen-index <= 30 {
			allKeys = append(allKeys, 30)
		}

		if totalLen-index <= 20 {
			allKeys = append(allKeys, 20)
		}

		if totalLen-index <= 10 {
			allKeys = append(allKeys, 10)
		}

		for _, key := range allKeys {
			if _, ok := m[curAe]; !ok {
				m[curAe] = make(map[int]int)
			}
			if _, ok := m[curAe][key]; !ok {
				m[curAe][key] = 1
			} else {
				m[curAe][key] += 1
			}
		}
	}

	allDltLen := len(dlts)
	for i, dlt := range dlts {
		//ae := gen.GetDltFront5ABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5})
		ae := gen.GetDltFrontABCDEStrFromCustomABCDEMap([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, abcdeMap)
		//ae4s := gen.GetDltFrontSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 4, true)
		//ae3s := gen.GetDltFrontSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 3, true)
		//ae2s := gen.GetDltFrontSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 2, true)
		//ae1s := gen.GetDltFrontSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 1, true)

		ae4s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 4, abcdeMap, true)
		ae3s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 3, abcdeMap, true)
		ae2s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 2, abcdeMap, true)
		ae1s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 1, abcdeMap, true)

		for _, iae := range ae4s {
			if _, ok := cur4AeMaps[iae]; !ok {
				cur4AeMaps[iae] = 1
			} else {
				cur4AeMaps[iae] += 1
			}
			dnf(cur4QAeMaps, iae, allDltLen, i)
		}

		for _, iae := range ae3s {
			if _, ok := cur3AeMaps[iae]; !ok {
				cur3AeMaps[iae] = 1
			} else {
				cur3AeMaps[iae] += 1
			}
			dnf(cur3QAeMaps, iae, allDltLen, i)
		}

		for _, iae := range ae2s {
			if _, ok := cur2AeMaps[iae]; !ok {
				cur2AeMaps[iae] = 1
			} else {
				cur2AeMaps[iae] += 1
			}
			dnf(cur2QAeMaps, iae, allDltLen, i)
		}

		for _, iae := range ae1s {
			if _, ok := cur1AeMaps[iae]; !ok {
				cur1AeMaps[iae] = 1
			} else {
				cur1AeMaps[iae] += 1
			}
			dnf(cur1QAeMaps, iae, allDltLen, i)
		}

		if _, ok := cur5AeMaps[ae]; !ok {
			cur5AeMaps[ae] = 1
		} else {
			cur5AeMaps[ae] += 1
		}
		dnf(cur5QAeMaps, ae, allDltLen, i)
		if i+1 == len(dlts) {
			break
		}

		//nae := gen.GetDltFront5ABCDE([]string{dlts[i+1].F1, dlts[i+1].F2, dlts[i+1].F3, dlts[i+1].F4, dlts[i+1].F5})
		nae := gen.GetDltFrontABCDEStrFromCustomABCDEMap([]string{dlts[i+1].F1, dlts[i+1].F2, dlts[i+1].F3, dlts[i+1].F4, dlts[i+1].F5}, abcdeMap)
		if _, ok1 := aeMaps[ae]; !ok1 {
			aeMaps[ae] = make(map[string][]models.Dlt)
		}
		aeMaps[ae][nae] = append(aeMaps[ae][nae], dlts[i+1])
	}

	// 把 map 转换为键值对切片
	type kv struct {
		Key   string
		Value int
	}
	var kv5s []kv
	var kv4s []kv
	var kv3s []kv
	var kv2s []kv
	var kv1s []kv
	var curAes []string
	for k, v := range cur5AeMaps {
		curAes = append(curAes, k)
		kv5s = append(kv5s, kv{k, v})
	}
	for k, v := range cur4AeMaps {
		kv4s = append(kv4s, kv{k, v})
	}
	for k, v := range cur3AeMaps {
		kv3s = append(kv3s, kv{k, v})
	}

	for k, v := range cur2AeMaps {
		kv2s = append(kv2s, kv{k, v})
	}
	for k, v := range cur1AeMaps {
		kv1s = append(kv1s, kv{k, v})
	}

	// 按 Value 从大到小排序
	sort.Slice(kv5s, func(i, j int) bool {
		return kv5s[i].Value > kv5s[j].Value
	})
	sort.Slice(kv4s, func(i, j int) bool {
		return kv4s[i].Value > kv4s[j].Value
	})
	sort.Slice(kv3s, func(i, j int) bool {
		return kv3s[i].Value > kv3s[j].Value
	})

	sort.Slice(kv2s, func(i, j int) bool {
		return kv2s[i].Value > kv2s[j].Value
	})

	sort.Slice(kv1s, func(i, j int) bool {
		return kv1s[i].Value > kv1s[j].Value
	})

	var curFrontNotAppearAes []string
	for ae, m := range aeMaps {
		for nae, v := range m {
			fmt.Printf("%10s -> %10s -> len=%3d\n", ae, nae, len(v))
		}
	}

	curFrontNotAppearAes = gen.DiffSlice(allFrontAes, curAes)
	fmt.Printf("当前未出现的前区ABCDE：%v\n", curFrontNotAppearAes)
	fmt.Printf("5个前区号码的出现情况------------------------------\n")
	for _, item := range kv5s {
		fmt.Printf("%s: %d\n", item.Key, item.Value)
		for _, q := range zXQs {
			if _, ok := cur5QAeMaps[item.Key][q]; ok {
				fmt.Printf("---最新%3d期: 出现%d期\n", q, cur5QAeMaps[item.Key][q])
			} else {
				fmt.Printf("---最新%3d期: 出现%d期\n", q, 0)
			}
		}
	}

	fmt.Printf("4个前区号码的出现情况------------------------------\n")
	for _, item := range kv4s {
		fmt.Printf("%s: %d\n", item.Key, item.Value)
		for _, q := range zXQs {
			if _, ok := cur4QAeMaps[item.Key][q]; ok {
				fmt.Printf("---最新%3d期: 出现%d期\n", q, cur4QAeMaps[item.Key][q])
			} else {
				fmt.Printf("---最新%3d期: 出现%d期\n", q, 0)
			}
		}
	}
	fmt.Printf("3个前区号码的出现情况------------------------------\n")
	for _, item := range kv3s {
		fmt.Printf("%s: %d\n", item.Key, item.Value)
		for _, q := range zXQs {
			if _, ok := cur3QAeMaps[item.Key][q]; ok {
				fmt.Printf("---最新%3d期: 出现%d期\n", q, cur3QAeMaps[item.Key][q])
			} else {
				fmt.Printf("---最新%3d期: 出现%d期\n", q, 0)
			}
		}
	}

	fmt.Printf("2个前区号码的出现情况------------------------------\n")
	for _, item := range kv2s {
		fmt.Printf("%s: %d\n", item.Key, item.Value)
		for _, q := range zXQs {
			if _, ok := cur2QAeMaps[item.Key][q]; ok {
				fmt.Printf("---最新%3d期: 出现%d期\n", q, cur2QAeMaps[item.Key][q])
			} else {
				fmt.Printf("---最新%3d期: 出现%d期\n", q, 0)
			}
		}
	}

	fmt.Printf("1个前区号码的出现情况------------------------------\n")
	for _, item := range kv1s {
		fmt.Printf("%s: %d\n", item.Key, item.Value)
		for _, q := range zXQs {
			if _, ok := cur1QAeMaps[item.Key][q]; ok {
				fmt.Printf("---最新%3d期: 出现%d期\n", q, cur1QAeMaps[item.Key][q])
			} else {
				fmt.Printf("---最新%3d期: 出现%d期\n", q, 0)
			}
		}
	}

}

func NewAbcdeQs(abcdeTyp string) {
	dlts, _ := dbop.ReadAllDlt(false)
	allDltLen := len(dlts)
	allFrontAes := gen.GetDltAllFrontABCDE(5)
	abcdeMap := make(map[string][]string)
	seedHi := uint64(1)
	switch abcdeTyp {
	case "77777":
		abcdeMap = gen.Rand77777ABCDEs(seedHi)
	case "215432":
		abcdeMap = gen.Rand215432ABCDEs(seedHi)
	case "224432":
		abcdeMap = gen.Rand224432ABCDEs(seedHi)
	case "224441":
		abcdeMap = gen.Rand224441ABCDEs(seedHi)
	case "253322":
		abcdeMap = gen.Rand253322ABCDEs(seedHi)
	case "272222":
		abcdeMap = gen.Rand272222ABCDEs(seedHi)
	default:
		abcdeMap = gen.Rand77777ABCDEs(seedHi)
	}

	err := os.MkdirAll("results", 0777)
	if err != nil {
		panic(fmt.Sprintf("无法创建%s目录：%v\n", "results", err))
	}

	fileName := fmt.Sprintf("%s_%s%s", abcdeTyp, time.Now().Format("2006_0102_1504"), ".txt")
	resFile, err := os.OpenFile(filepath.Join("results", fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	resWt := bufio.NewWriter(resFile)

	var dnf func(map[string]map[int]int, string, int, int)
	dnf = func(m map[string]map[int]int, curAe string, totalLen, index int) {
		var allKeys []int
		if totalLen-index <= 3000 {
			allKeys = append(allKeys, 3000)
		}
		if totalLen-index <= 2000 {
			allKeys = append(allKeys, 2000)
		}

		if totalLen-index <= 1000 {
			allKeys = append(allKeys, 1000)
		}

		if totalLen-index <= 500 {
			allKeys = append(allKeys, 500)
		}

		if totalLen-index <= 200 {
			allKeys = append(allKeys, 200)
		}
		if totalLen-index <= 100 {
			allKeys = append(allKeys, 100)
		}

		if totalLen-index <= 50 {
			allKeys = append(allKeys, 50)
		}

		if totalLen-index <= 30 {
			allKeys = append(allKeys, 30)
		}

		if totalLen-index <= 20 {
			allKeys = append(allKeys, 20)
		}

		if totalLen-index <= 10 {
			allKeys = append(allKeys, 10)
		}

		for _, key := range allKeys {
			if _, ok := m[curAe]; !ok {
				m[curAe] = make(map[int]int)
			}
			if _, ok := m[curAe][key]; !ok {
				m[curAe][key] = 1
			} else {
				m[curAe][key] += 1
			}
		}
	}

	//aeMaps := make(map[string]map[string][]models.Dlt)
	cur5AeMaps := make(map[string]int)
	//cur4AeMaps := make(map[string]int)
	//cur3AeMaps := make(map[string]int)
	//cur2AeMaps := make(map[string]int)
	//cur1AeMaps := make(map[string]int)

	cur5QAeMaps := make(map[string]map[int]int)
	//cur4QAeMaps := make(map[string]map[int]int)
	//cur3QAeMaps := make(map[string]map[int]int)
	//cur2QAeMaps := make(map[string]map[int]int)
	//cur1QAeMaps := make(map[string]map[int]int)
	zXQs := []int{3000, 2000, 1000, 500, 200, 100, 50, 30, 20, 10}

	// 把 map 转换为键值对切片
	type kv struct {
		Key   string
		Value int
	}
	var kv5s []kv
	//var kv4s []kv
	//var kv3s []kv
	//var kv2s []kv
	//var kv1s []kv
	var curAes []string
	var curFrontNotAppearAes []string
	var curMax int

LabelForContinue:
	//aeMaps = make(map[string]map[string][]models.Dlt)
	cur5AeMaps = make(map[string]int)
	//cur4AeMaps = make(map[string]int)
	//cur3AeMaps = make(map[string]int)
	//cur2AeMaps = make(map[string]int)
	//cur1AeMaps = make(map[string]int)

	cur5QAeMaps = make(map[string]map[int]int)
	//cur4QAeMaps = make(map[string]map[int]int)
	//cur3QAeMaps = make(map[string]map[int]int)
	//cur2QAeMaps = make(map[string]map[int]int)
	//cur1QAeMaps = make(map[string]map[int]int)

	for i, dlt := range dlts {
		ae := gen.GetDltFrontABCDEStrFromCustomABCDEMap([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, abcdeMap)
		//ae4s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 4, abcdeMap, true)
		//ae3s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 3, abcdeMap, true)
		//ae2s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 2, abcdeMap, true)
		//ae1s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 1, abcdeMap, true)

		//for _, iae := range ae4s {
		//	if _, ok := cur4AeMaps[iae]; !ok {
		//		cur4AeMaps[iae] = 1
		//	} else {
		//		cur4AeMaps[iae] += 1
		//	}
		//	dnf(cur4QAeMaps, iae, allDltLen, i)
		//}
		//
		//for _, iae := range ae3s {
		//	if _, ok := cur3AeMaps[iae]; !ok {
		//		cur3AeMaps[iae] = 1
		//	} else {
		//		cur3AeMaps[iae] += 1
		//	}
		//	dnf(cur3QAeMaps, iae, allDltLen, i)
		//}
		//
		//for _, iae := range ae2s {
		//	if _, ok := cur2AeMaps[iae]; !ok {
		//		cur2AeMaps[iae] = 1
		//	} else {
		//		cur2AeMaps[iae] += 1
		//	}
		//	dnf(cur2QAeMaps, iae, allDltLen, i)
		//}
		//
		//for _, iae := range ae1s {
		//	if _, ok := cur1AeMaps[iae]; !ok {
		//		cur1AeMaps[iae] = 1
		//	} else {
		//		cur1AeMaps[iae] += 1
		//	}
		//	dnf(cur1QAeMaps, iae, allDltLen, i)
		//}

		if _, ok := cur5AeMaps[ae]; !ok {
			cur5AeMaps[ae] = 1
		} else {
			cur5AeMaps[ae] += 1
		}
		dnf(cur5QAeMaps, ae, allDltLen, i)
		if i+1 == len(dlts) {
			break
		}

		//nae := gen.GetDltFrontABCDEStrFromCustomABCDEMap([]string{dlts[i+1].F1, dlts[i+1].F2, dlts[i+1].F3, dlts[i+1].F4, dlts[i+1].F5}, abcdeMap)
		//if _, ok1 := aeMaps[ae]; !ok1 {
		//	aeMaps[ae] = make(map[string][]models.Dlt)
		//}
		//aeMaps[ae][nae] = append(aeMaps[ae][nae], dlts[i+1])
	}
	kv5s = nil
	//kv4s = nil
	//kv3s = nil
	//kv2s = nil
	//kv1s = nil
	curAes = nil
	curFrontNotAppearAes = nil

	for k, v := range cur5AeMaps {
		curAes = append(curAes, k)
		kv5s = append(kv5s, kv{k, v})
	}
	//for k, v := range cur4AeMaps {
	//	kv4s = append(kv4s, kv{k, v})
	//}
	//for k, v := range cur3AeMaps {
	//	kv3s = append(kv3s, kv{k, v})
	//}
	//
	//for k, v := range cur2AeMaps {
	//	kv2s = append(kv2s, kv{k, v})
	//}
	//for k, v := range cur1AeMaps {
	//	kv1s = append(kv1s, kv{k, v})
	//}

	// 按 Value 从大到小排序
	sort.Slice(kv5s, func(i, j int) bool {
		return kv5s[i].Value > kv5s[j].Value
	})
	//sort.Slice(kv4s, func(i, j int) bool {
	//	return kv4s[i].Value > kv4s[j].Value
	//})
	//sort.Slice(kv3s, func(i, j int) bool {
	//	return kv3s[i].Value > kv3s[j].Value
	//})
	//
	//sort.Slice(kv2s, func(i, j int) bool {
	//	return kv2s[i].Value > kv2s[j].Value
	//})
	//
	//sort.Slice(kv1s, func(i, j int) bool {
	//	return kv1s[i].Value > kv1s[j].Value
	//})

	curFrontNotAppearAes = gen.DiffSlice(allFrontAes, curAes)
	for _, item := range kv5s {
		if !judgeHavaABCDE(item.Key) {
			continue
		}

		if item.Value > curMax {
			curMax = item.Value
			_, _ = resWt.WriteString(fmt.Sprintf("类型：%s\n", abcdeTyp))
			_, _ = resWt.WriteString(fmt.Sprintf("A=%q\nB=%q\nC=%q\nD=%q\nE=%q\n", abcdeMap["A"], abcdeMap["B"], abcdeMap["C"], abcdeMap["D"], abcdeMap["E"]))
			_, _ = resWt.WriteString(fmt.Sprintf("当前未出现的ABCDE：%v\n", curFrontNotAppearAes))
			_, _ = resWt.WriteString(fmt.Sprintf("%s: %d\n", item.Key, item.Value))

			for _, q := range zXQs {
				if _, ok := cur5QAeMaps[item.Key][q]; ok {
					_, _ = resWt.WriteString(fmt.Sprintf("---最新%4d期出现%4d期\n", q, cur5QAeMaps[item.Key][q]))
				} else {
					_, _ = resWt.WriteString(fmt.Sprintf("---最新%4d期出现%4d期\n", q, 0))
				}
			}
			_ = resWt.Flush()
		}
	}

	goto LabelForContinue
}

//var JX35 = 35 * 34 * 33 * 32 * 31 * 30 * 29 * 28 * 27 * 26 * 25 * 24 * 23 * 22 * 21 * 20 * 19 * 18 * 17 * 16 * 15 * 14 * 13 * 12 * 11 * 10 * 9 * 8 * 7 * 6 * 5 * 4 * 3 * 2

func NewAbcdeQsWithWg(abcdeTyp string, wg *sync.WaitGroup) {
	dlts, _ := dbop.ReadAllDlt(false)
	allDltLen := len(dlts)
	allFrontAes := gen.GetDltAllFrontABCDE(5)
	abcdeMap := make(map[string][]string)
	var ANum, BNum, CNum, DNum, ENum int
	var hisMaxNum int
	switch abcdeTyp {
	case "77777":
		{
			ANum, BNum, CNum, DNum, ENum = 7, 7, 7, 7, 7
			hisMaxNum = cfg.Default.AeHisMax77777
		}
	case "116666":
		{
			ANum, BNum, CNum, DNum, ENum = 11, 6, 6, 6, 6
			hisMaxNum = cfg.Default.AeHisMax116666
		}
	case "155555":
		{
			ANum, BNum, CNum, DNum, ENum = 15, 5, 5, 5, 5
			hisMaxNum = cfg.Default.AeHisMax155555
		}
	case "194444":
		{
			ANum, BNum, CNum, DNum, ENum = 19, 4, 4, 4, 4
			hisMaxNum = cfg.Default.AeHisMax194444
		}
	case "215432":
		{
			ANum, BNum, CNum, DNum, ENum = 21, 5, 4, 3, 2
			hisMaxNum = cfg.Default.AeHisMax215432
		}
	case "224432":
		{
			ANum, BNum, CNum, DNum, ENum = 22, 4, 4, 3, 2
			hisMaxNum = cfg.Default.AeHisMax224432
		}
	case "224441":
		{
			ANum, BNum, CNum, DNum, ENum = 22, 4, 4, 4, 1
			hisMaxNum = cfg.Default.AeHisMax224441
		}
	case "233333":
		{
			ANum, BNum, CNum, DNum, ENum = 25, 3, 3, 3, 3
			hisMaxNum = cfg.Default.AeHisMax233333
		}
	case "253322":
		{
			ANum, BNum, CNum, DNum, ENum = 25, 3, 3, 2, 2
			hisMaxNum = cfg.Default.AeHisMax253322
		}
	case "272222":
		{
			ANum, BNum, CNum, DNum, ENum = 27, 2, 2, 2, 2
			hisMaxNum = cfg.Default.AeHisMax272222
		}
	default:
		{
			ANum, BNum, CNum, DNum, ENum = 7, 7, 7, 7, 7
			hisMaxNum = 223
		}
	}

	err := os.MkdirAll("results", 0777)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("无法创建%s目录：%v\n", "results", err))
		wg.Done()
		return
	}

	fileName := fmt.Sprintf("%s_%s_hismax_%d%s", abcdeTyp, time.Now().Format("2006_0102_1504"), hisMaxNum, ".txt")
	resFile, err := os.OpenFile(filepath.Join("results", fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("无法打开%s文件：%v\n", fileName, err))
		wg.Done()
		return
	}
	resWt := bufio.NewWriter(resFile)

	var dnf func(map[string]map[int]int, string, int, int)
	dnf = func(m map[string]map[int]int, curAe string, totalLen, index int) {
		var allKeys []int
		if totalLen-index <= 3000 {
			allKeys = append(allKeys, 3000)
		}
		if totalLen-index <= 2000 {
			allKeys = append(allKeys, 2000)
		}

		if totalLen-index <= 1000 {
			allKeys = append(allKeys, 1000)
		}

		if totalLen-index <= 500 {
			allKeys = append(allKeys, 500)
		}

		if totalLen-index <= 200 {
			allKeys = append(allKeys, 200)
		}
		if totalLen-index <= 100 {
			allKeys = append(allKeys, 100)
		}

		if totalLen-index <= 50 {
			allKeys = append(allKeys, 50)
		}

		if totalLen-index <= 30 {
			allKeys = append(allKeys, 30)
		}

		if totalLen-index <= 20 {
			allKeys = append(allKeys, 20)
		}

		if totalLen-index <= 10 {
			allKeys = append(allKeys, 10)
		}

		for _, key := range allKeys {
			if _, ok := m[curAe]; !ok {
				m[curAe] = make(map[int]int)
			}
			if _, ok := m[curAe][key]; !ok {
				m[curAe][key] = 1
			} else {
				m[curAe][key] += 1
			}
		}
	}

	//aeMaps := make(map[string]map[string][]models.Dlt)
	cur5AeMaps := make(map[string]int)
	//cur4AeMaps := make(map[string]int)
	//cur3AeMaps := make(map[string]int)
	//cur2AeMaps := make(map[string]int)
	//cur1AeMaps := make(map[string]int)

	cur5QAeMaps := make(map[string]map[int]int)
	//cur4QAeMaps := make(map[string]map[int]int)
	//cur3QAeMaps := make(map[string]map[int]int)
	//cur2QAeMaps := make(map[string]map[int]int)
	//cur1QAeMaps := make(map[string]map[int]int)
	zXQs := []int{3000, 2000, 1000, 500, 200, 100, 50, 30, 20, 10}
	var curDontExistAes []string
	// 把 map 转换为键值对切片
	type kv struct {
		Key   string
		Value int
	}
	var kv5s []kv
	//var kv4s []kv
	//var kv3s []kv
	//var kv2s []kv
	//var kv1s []kv
	var curAes []string
	var curFrontNotAppearAes []string
	var curMax int
	seedHi := uint64(0)
	tempNum := 0
LabelForContinue:
	if tempNum == 0 {
		tempNum = 1
		var tempAbcdes [][]string
		switch abcdeTyp {
		case "77777":
			tempAbcdes = cfg.Default.AeMax77777
		case "116666":
			tempAbcdes = cfg.Default.AeMax116666
		case "155555":
			tempAbcdes = cfg.Default.AeMax155555
		case "194444":
			tempAbcdes = cfg.Default.AeMax194444
		case "215432":
			tempAbcdes = cfg.Default.AeMax215432
		case "224432":
			tempAbcdes = cfg.Default.AeMax224432
		case "224441":
			tempAbcdes = cfg.Default.AeMax224441
		case "233333":
			tempAbcdes = cfg.Default.AeMax233333
		case "253322":
			tempAbcdes = cfg.Default.AeMax253322
		case "272222":
			tempAbcdes = cfg.Default.AeMax272222
		default:
			tempAbcdes = cfg.Default.AeMax77777
		}

		abcdeMap["A"] = tempAbcdes[0]
		abcdeMap["B"] = tempAbcdes[1]
		abcdeMap["C"] = tempAbcdes[2]
		abcdeMap["D"] = tempAbcdes[3]
		abcdeMap["E"] = tempAbcdes[4]
	} else {
		switch abcdeTyp {
		case "77777":
			abcdeMap = gen.Rand77777ABCDEs(seedHi)
		case "116666":
			abcdeMap = gen.Rand116666ABCDEs(seedHi)
		case "155555":
			abcdeMap = gen.Rand155555ABCDEs(seedHi)
		case "194444":
			abcdeMap = gen.Rand194444ABCDEs(seedHi)
		case "215432":
			abcdeMap = gen.Rand215432ABCDEs(seedHi)
		case "224432":
			abcdeMap = gen.Rand224432ABCDEs(seedHi)
		case "224441":
			abcdeMap = gen.Rand224441ABCDEs(seedHi)
		case "233333":
			abcdeMap = gen.Rand233333ABCDEs(seedHi)
		case "253322":
			abcdeMap = gen.Rand253322ABCDEs(seedHi)
		case "272222":
			abcdeMap = gen.Rand272222ABCDEs(seedHi)
		default:
			abcdeMap = gen.Rand77777ABCDEs(seedHi)
		}
	}

	//aeMaps = make(map[string]map[string][]models.Dlt)
	cur5AeMaps = make(map[string]int)
	//cur4AeMaps = make(map[string]int)
	//cur3AeMaps = make(map[string]int)
	//cur2AeMaps = make(map[string]int)
	//cur1AeMaps = make(map[string]int)

	cur5QAeMaps = make(map[string]map[int]int)
	//cur4QAeMaps = make(map[string]map[int]int)
	//cur3QAeMaps = make(map[string]map[int]int)
	//cur2QAeMaps = make(map[string]map[int]int)
	//cur1QAeMaps = make(map[string]map[int]int)

	for i, dlt := range dlts {
		ae := gen.GetDltFrontABCDEStrFromCustomABCDEMap([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, abcdeMap)
		//ae4s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 4, abcdeMap, true)
		//ae3s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 3, abcdeMap, true)
		//ae2s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 2, abcdeMap, true)
		//ae1s := gen.GetDltFrontCustomSpNumABCDE([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, 1, abcdeMap, true)

		//for _, iae := range ae4s {
		//	if _, ok := cur4AeMaps[iae]; !ok {
		//		cur4AeMaps[iae] = 1
		//	} else {
		//		cur4AeMaps[iae] += 1
		//	}
		//	dnf(cur4QAeMaps, iae, allDltLen, i)
		//}
		//
		//for _, iae := range ae3s {
		//	if _, ok := cur3AeMaps[iae]; !ok {
		//		cur3AeMaps[iae] = 1
		//	} else {
		//		cur3AeMaps[iae] += 1
		//	}
		//	dnf(cur3QAeMaps, iae, allDltLen, i)
		//}
		//
		//for _, iae := range ae2s {
		//	if _, ok := cur2AeMaps[iae]; !ok {
		//		cur2AeMaps[iae] = 1
		//	} else {
		//		cur2AeMaps[iae] += 1
		//	}
		//	dnf(cur2QAeMaps, iae, allDltLen, i)
		//}
		//
		//for _, iae := range ae1s {
		//	if _, ok := cur1AeMaps[iae]; !ok {
		//		cur1AeMaps[iae] = 1
		//	} else {
		//		cur1AeMaps[iae] += 1
		//	}
		//	dnf(cur1QAeMaps, iae, allDltLen, i)
		//}

		if _, ok := cur5AeMaps[ae]; !ok {
			cur5AeMaps[ae] = 1
		} else {
			cur5AeMaps[ae] += 1
		}
		dnf(cur5QAeMaps, ae, allDltLen, i)
		if i+1 == len(dlts) {
			break
		}

		//nae := gen.GetDltFrontABCDEStrFromCustomABCDEMap([]string{dlts[i+1].F1, dlts[i+1].F2, dlts[i+1].F3, dlts[i+1].F4, dlts[i+1].F5}, abcdeMap)
		//if _, ok1 := aeMaps[ae]; !ok1 {
		//	aeMaps[ae] = make(map[string][]models.Dlt)
		//}
		//aeMaps[ae][nae] = append(aeMaps[ae][nae], dlts[i+1])
	}
	kv5s = nil
	//kv4s = nil
	//kv3s = nil
	//kv2s = nil
	//kv1s = nil
	curAes = nil
	curFrontNotAppearAes = nil

	for k, v := range cur5AeMaps {
		curAes = append(curAes, k)
		kv5s = append(kv5s, kv{k, v})
	}
	//for k, v := range cur4AeMaps {
	//	kv4s = append(kv4s, kv{k, v})
	//}
	//for k, v := range cur3AeMaps {
	//	kv3s = append(kv3s, kv{k, v})
	//}
	//
	//for k, v := range cur2AeMaps {
	//	kv2s = append(kv2s, kv{k, v})
	//}
	//for k, v := range cur1AeMaps {
	//	kv1s = append(kv1s, kv{k, v})
	//}

	// 按 Value 从大到小排序
	sort.Slice(kv5s, func(i, j int) bool {
		return kv5s[i].Value > kv5s[j].Value
	})
	//sort.Slice(kv4s, func(i, j int) bool {
	//	return kv4s[i].Value > kv4s[j].Value
	//})
	//sort.Slice(kv3s, func(i, j int) bool {
	//	return kv3s[i].Value > kv3s[j].Value
	//})
	//
	//sort.Slice(kv2s, func(i, j int) bool {
	//	return kv2s[i].Value > kv2s[j].Value
	//})
	//
	//sort.Slice(kv1s, func(i, j int) bool {
	//	return kv1s[i].Value > kv1s[j].Value
	//})

	curFrontNotAppearAes = gen.DiffSlice(allFrontAes, curAes)
	curDontExistAes = nil
	for _, v := range curFrontNotAppearAes {
		curABCDENums := getABCDENum(v)
		if curABCDENums[0] > ANum ||
			curABCDENums[1] > BNum ||
			curABCDENums[2] > CNum ||
			curABCDENums[3] > DNum ||
			curABCDENums[4] > ENum {
			curDontExistAes = append(curDontExistAes, v)
		}
	}
	curFrontNotAppearAes = gen.DiffSlice(curFrontNotAppearAes, curDontExistAes)
	for _, item := range kv5s {
		//if !judgeHavaBCDE(item.Key) {
		//	continue
		//}
		if !judgeHavaABCDE(item.Key) {
			continue
		}

		if item.Value >= curMax ||
			//(abcdeTyp == "77777" && item.Value >= 200) ||
			//(abcdeTyp == "116666" && item.Value >= 200) ||
			//(abcdeTyp == "155555" && item.Value >= 200) ||
			//(abcdeTyp == "194444" && item.Value >= 200) ||
			//(abcdeTyp == "215432" && item.Value >= 330) ||
			//(abcdeTyp == "224432" && item.Value >= 340) ||
			//(abcdeTyp == "224441" && item.Value >= 330) ||
			//(abcdeTyp == "233333" && item.Value >= 390) ||
			//(abcdeTyp == "253322" && item.Value >= 570) ||
			//(abcdeTyp == "272222" && item.Value >= 830
			(abcdeTyp == "77777" && item.Value >= 200) ||
			(abcdeTyp == "116666" && item.Value >= 190) ||
			(abcdeTyp == "155555" && item.Value >= 200) ||
			(abcdeTyp == "194444" && item.Value >= 200) ||
			(abcdeTyp == "215432" && item.Value >= 330) ||
			(abcdeTyp == "224432" && item.Value >= 340) ||
			(abcdeTyp == "224441" && item.Value >= 330) ||
			(abcdeTyp == "233333" && item.Value >= 390) ||
			(abcdeTyp == "253322" && item.Value >= 570) ||
			(abcdeTyp == "272222" && item.Value >= 830) {
			if item.Value >= curMax {
				curMax = item.Value
			}

			_, _ = resWt.WriteString(fmt.Sprintf("类型：%s <- 当前最大值：%d\n", abcdeTyp, curMax))
			_, _ = resWt.WriteString(fmt.Sprintf("A=%q\nB=%q\nC=%q\nD=%q\nE=%q\n", abcdeMap["A"], abcdeMap["B"], abcdeMap["C"], abcdeMap["D"], abcdeMap["E"]))
			_, _ = resWt.WriteString(fmt.Sprintf("当前未出现的ABCDE：%v\n", curFrontNotAppearAes))
			_, _ = resWt.WriteString(fmt.Sprintf("%s: %d<-\n", item.Key, item.Value))

			for _, q := range zXQs {
				if _, ok := cur5QAeMaps[item.Key][q]; ok {
					_, _ = resWt.WriteString(fmt.Sprintf("---最新%4d期出现%4d期\n", q, cur5QAeMaps[item.Key][q]))
				} else {
					_, _ = resWt.WriteString(fmt.Sprintf("---最新%4d期出现%4d期\n", q, 0))
				}
			}
			_ = resWt.Flush()
		}
	}
	if seedHi >= 18446744073709551615 {
		wg.Done()
	} else {
		seedHi++
		if seedHi%1000000 == 0 {
			lg.InfoToFileAndStdOut(fmt.Sprintf("%10s 已运行 %20d 次\n", abcdeTyp, seedHi))
		}
		goto LabelForContinue
	}
}

func judgeHavaABCDE(str string) bool {
	if strings.Contains(str, "A") &&
		strings.Contains(str, "B") &&
		strings.Contains(str, "C") &&
		strings.Contains(str, "D") &&
		strings.Contains(str, "E") {
		return true
	} else {
		return false
	}
}

func judgeHavaBCDE(str string) bool {
	if strings.Contains(str, "B") &&
		strings.Contains(str, "C") &&
		strings.Contains(str, "D") &&
		strings.Contains(str, "E") {
		return true
	} else {
		return false
	}
}

func getABCDENum(str string) (res []int) {
	res = append(res, getNum(str, "A"))
	res = append(res, getNum(str, "B"))
	res = append(res, getNum(str, "C"))
	res = append(res, getNum(str, "D"))
	res = append(res, getNum(str, "E"))
	return res
}

// getNum 注意只能用于英文和数字组成的字符串，不能用于中文等字符串中
func getNum(str, spStr string) (n int) {
	if strings.Contains(str, spStr) {
		first := strings.Split(str, spStr)[0]
		n, _ = strconv.Atoi(first[len(first)-1:])
		return n
	}
	return n
}

func ValidateABCDE(abcdeTyp string) {
	dlts, _ := dbop.ReadAllDlt(false)
	allDltLen := len(dlts)
	allFrontAes := gen.GetDltAllFrontABCDE(5)
	abcdeMap := make(map[string][]string)

	err := os.MkdirAll("validate", 0777)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("无法创建%s目录：%v\n", "validate", err))
		return
	}

	fileName := fmt.Sprintf("%s_%s%s", abcdeTyp, time.Now().Format("2006_0102_1504"), ".txt")
	resFile, err := os.OpenFile(filepath.Join("validate", fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("无法打开%s文件：%v\n", fileName, err))
		return
	}
	resWt := bufio.NewWriter(resFile)

	var dnf func(map[string]map[int]int, string, int, int)
	dnf = func(m map[string]map[int]int, curAe string, totalLen, index int) {
		var allKeys []int
		if totalLen-index <= 3000 {
			allKeys = append(allKeys, 3000)
		}
		if totalLen-index <= 2000 {
			allKeys = append(allKeys, 2000)
		}

		if totalLen-index <= 1000 {
			allKeys = append(allKeys, 1000)
		}

		if totalLen-index <= 500 {
			allKeys = append(allKeys, 500)
		}

		if totalLen-index <= 200 {
			allKeys = append(allKeys, 200)
		}
		if totalLen-index <= 100 {
			allKeys = append(allKeys, 100)
		}

		if totalLen-index <= 50 {
			allKeys = append(allKeys, 50)
		}

		if totalLen-index <= 30 {
			allKeys = append(allKeys, 30)
		}

		if totalLen-index <= 20 {
			allKeys = append(allKeys, 20)
		}

		if totalLen-index <= 10 {
			allKeys = append(allKeys, 10)
		}

		for _, key := range allKeys {
			if _, ok := m[curAe]; !ok {
				m[curAe] = make(map[int]int)
			}
			if _, ok := m[curAe][key]; !ok {
				m[curAe][key] = 1
			} else {
				m[curAe][key] += 1
			}
		}
	}

	//aeMaps := make(map[string]map[string][]models.Dlt)
	cur5AeMaps := make(map[string]int)
	//cur4AeMaps := make(map[string]int)
	//cur3AeMaps := make(map[string]int)
	//cur2AeMaps := make(map[string]int)
	//cur1AeMaps := make(map[string]int)

	cur5QAeMaps := make(map[string]map[int]int)
	//cur4QAeMaps := make(map[string]map[int]int)
	//cur3QAeMaps := make(map[string]map[int]int)
	//cur2QAeMaps := make(map[string]map[int]int)
	//cur1QAeMaps := make(map[string]map[int]int)
	zXQs := []int{3000, 2000, 1000, 500, 200, 100, 50, 30, 20, 10}

	// 把 map 转换为键值对切片
	type kv struct {
		Key   string
		Value int
	}
	var kv5s []kv
	//var kv4s []kv
	//var kv3s []kv
	//var kv2s []kv
	//var kv1s []kv
	var curAes []string
	var curFrontNotAppearAes []string
	//var curMax int
	abcdeMap["A"] = []string{"09", "13", "17", "20", "27", "28", "32"}
	abcdeMap["B"] = []string{"02", "05", "11", "14", "24", "30", "33"}
	abcdeMap["C"] = []string{"04", "07", "10", "16", "21", "29", "34"}
	abcdeMap["D"] = []string{"08", "12", "18", "19", "22", "23", "35"}
	abcdeMap["E"] = []string{"01", "03", "06", "15", "25", "26", "31"}

	//aeMaps = make(map[string]map[string][]models.Dlt)
	cur5AeMaps = make(map[string]int)
	//cur4AeMaps = make(map[string]int)
	//cur3AeMaps = make(map[string]int)
	//cur2AeMaps = make(map[string]int)
	//cur1AeMaps = make(map[string]int)

	cur5QAeMaps = make(map[string]map[int]int)
	//cur4QAeMaps = make(map[string]map[int]int)
	//cur3QAeMaps = make(map[string]map[int]int)
	//cur2QAeMaps = make(map[string]map[int]int)
	//cur1QAeMaps = make(map[string]map[int]int)

	for i, dlt := range dlts {
		ae := gen.GetDltFrontABCDEStrFromCustomABCDEMap([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, abcdeMap)

		if _, ok := cur5AeMaps[ae]; !ok {
			cur5AeMaps[ae] = 1
		} else {
			cur5AeMaps[ae] += 1
		}
		dnf(cur5QAeMaps, ae, allDltLen, i)
		if i+1 == len(dlts) {
			break
		}

	}
	kv5s = nil
	curAes = nil
	curFrontNotAppearAes = nil

	for k, v := range cur5AeMaps {
		curAes = append(curAes, k)
		kv5s = append(kv5s, kv{k, v})
	}

	// 按 Value 从大到小排序
	sort.Slice(kv5s, func(i, j int) bool {
		return kv5s[i].Value > kv5s[j].Value
	})

	curFrontNotAppearAes = gen.DiffSlice(allFrontAes, curAes)
	for _, item := range kv5s {
		if item.Value > 50 {
			_, _ = resWt.WriteString(fmt.Sprintf("类型：%s\n", abcdeTyp))
			_, _ = resWt.WriteString(fmt.Sprintf("A=%q\nB=%q\nC=%q\nD=%q\nE=%q\n", abcdeMap["A"], abcdeMap["B"], abcdeMap["C"], abcdeMap["D"], abcdeMap["E"]))
			_, _ = resWt.WriteString(fmt.Sprintf("当前未出现的ABCDE：%v\n", curFrontNotAppearAes))
			_, _ = resWt.WriteString(fmt.Sprintf("%s: %d\n", item.Key, item.Value))

			for _, q := range zXQs {
				if _, ok := cur5QAeMaps[item.Key][q]; ok {
					_, _ = resWt.WriteString(fmt.Sprintf("---最新%4d期出现%4d期\n", q, cur5QAeMaps[item.Key][q]))
				} else {
					_, _ = resWt.WriteString(fmt.Sprintf("---最新%4d期出现%4d期\n", q, 0))
				}
			}
			_ = resWt.Flush()
		}
	}
}

func DltFrontNextOneFCount() {
	dlts, _ := dbop.ReadAllDlt(false)
	fMs := make(map[string]map[string]int, 35)
	for i, dlt := range dlts {
		if i+1 == len(dlts) {
			break
		}
		nextDlt := dlts[i+1]
		fs := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}
		nfs := []string{nextDlt.F1, nextDlt.F2, nextDlt.F3, nextDlt.F4, nextDlt.F5}

		for _, f := range fs {
			if _, ok := fMs[f]; !ok {
				fMs[f] = make(map[string]int)
			}
			for _, nf := range nfs {
				if _, ok := fMs[f][nf]; !ok {
					fMs[f][nf] = 1
				} else {
					fMs[f][nf] += 1
				}
			}
		}
	}

	for i := 1; i <= 35; i++ {
		f := fmt.Sprintf("%02d", i)
		if _, ok := fMs[f]; ok {
			kLens := make([]KeyWithLength, 0, len(fMs[f]))
			for nf, v := range fMs[f] {
				kLens = append(kLens, KeyWithLength{
					Key:    nf,
					Length: v,
				})
			}
			sort.Slice(kLens, func(i, j int) bool {
				return kLens[i].Length > kLens[j].Length
			})
			fmt.Printf("%s------\n", f)
			for _, kLen := range kLens {
				fmt.Printf("%s %d\n", kLen.Key, kLen.Length)
			}
			fmt.Printf("%s------------\n", f)
		}
	}
}

func DltFrontNextTwoFCount() {
	dlts, _ := dbop.ReadAllDlt(false)
	fMs := make(map[string]map[string]int, 35)
	for i, dlt := range dlts {
		if i+1 == len(dlts) {
			break
		}
		nextDlt := dlts[i+1]
		fs := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}
		nfs := []string{nextDlt.F1, nextDlt.F2, nextDlt.F3, nextDlt.F4, nextDlt.F5}
		nextCombs := gen.Comb(nfs, 2)

		for _, f := range fs {
			if _, ok := fMs[f]; !ok {
				fMs[f] = make(map[string]int)
			}
			for _, nComb := range nextCombs {
				if _, ok := fMs[f][nComb]; !ok {
					fMs[f][nComb] = 1
				} else {
					fMs[f][nComb] += 1
				}
			}
		}
	}

	for i := 1; i <= 35; i++ {
		f := fmt.Sprintf("%02d", i)
		if _, ok := fMs[f]; ok {
			kLens := make([]KeyWithLength, 0, len(fMs[f]))
			for nf, v := range fMs[f] {
				kLens = append(kLens, KeyWithLength{
					Key:    nf,
					Length: v,
				})
			}
			sort.Slice(kLens, func(i, j int) bool {
				return kLens[i].Length > kLens[j].Length
			})
			fmt.Printf("%s------\n", f)
			for _, kLen := range kLens {
				fmt.Printf("%s %d\n", kLen.Key, kLen.Length)
			}
			fmt.Printf("%s------------\n", f)
		}
	}
}

// DltFrontOnlyOneHis 大乐透前区单个号码的历史
//
//	@Description:
//	@param wg
//	@return res
func DltFrontOnlyOneHis(wg *sync.WaitGroup) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2DltHis := make(map[string]*DltHis)

	// 初始化
	for k, v := range AllDltOnlyOneFront2Count {
		typ2DltHis[k] = &DltHis{Typ: k, AllCount: v}
	}

	lenDltHis := len(ZxDlts)

	for i, dlt := range ZxDlts {
		frontHms := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}

		for _, fHm := range frontHms {
			typ2DltHis[fHm].Cs = typ2DltHis[fHm].Cs + 1
			if lenDltHis-i <= 10 {
				typ2DltHis[fHm].Last10 = typ2DltHis[fHm].Last10 + 1
			}
			if lenDltHis-i <= 20 {
				typ2DltHis[fHm].Last20 = typ2DltHis[fHm].Last20 + 1
			}
			if lenDltHis-i <= 30 {
				typ2DltHis[fHm].Last30 = typ2DltHis[fHm].Last30 + 1
			}
			if lenDltHis-i <= 50 {
				typ2DltHis[fHm].Last50 = typ2DltHis[fHm].Last50 + 1
			}
			if lenDltHis-i <= 100 {
				typ2DltHis[fHm].Last100 = typ2DltHis[fHm].Last100 + 1
			}
			if lenDltHis-i <= 200 {
				typ2DltHis[fHm].Last200 = typ2DltHis[fHm].Last200 + 1
			}
			if lenDltHis-i <= 500 {
				typ2DltHis[fHm].Last500 = typ2DltHis[fHm].Last500 + 1
			}
			if lenDltHis-i <= 1000 {
				typ2DltHis[fHm].Last1000 = typ2DltHis[fHm].Last1000 + 1
			}
			if lenDltHis-i <= 1500 {
				typ2DltHis[fHm].Last1500 = typ2DltHis[fHm].Last1500 + 1
			}
			if lenDltHis-i <= 2000 {
				typ2DltHis[fHm].Last2000 = typ2DltHis[fHm].Last2000 + 1
			}
			if lenDltHis-i <= 2500 {
				typ2DltHis[fHm].Last2500 = typ2DltHis[fHm].Last2500 + 1
			}
			if lenDltHis-i <= 3500 {
				typ2DltHis[fHm].Last3500 = typ2DltHis[fHm].Last3500 + 1
			}
		}

	}

	kLens := make([]KeyWithLength, 0, len(typ2DltHis))

	for k, dltHis := range typ2DltHis {
		kLens = append(kLens, KeyWithLength{
			Key:    k,
			Length: dltHis.Cs,
		})
	}
	// 对typ2DltHis按照存放的Cs值的大小进行排序
	sort.Slice(kLens, func(i, j int) bool {
		if kLens[i].Length == kLens[j].Length {
			return kLens[i].Key > kLens[j].Key
		}
		return kLens[i].Length > kLens[j].Length
	})

	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2DltHis[typ])
	}

	return
}
