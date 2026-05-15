package ana_ssq

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/before80/lot/db"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
)

type Draw struct {
	Front []int
}

// LoadSsqFrontHmHistory
//
//	@Description:
//	@return []Draw
func LoadSsqFrontHmHistory() []Draw {
	ssqs, _ := dbop.ReadAllSsq(false)
	history := make([]Draw, 0, len(ssqs))
	for _, dlt := range ssqs {
		f1, _ := strconv.Atoi(dlt.F1)
		f2, _ := strconv.Atoi(dlt.F2)
		f3, _ := strconv.Atoi(dlt.F3)
		f4, _ := strconv.Atoi(dlt.F4)
		f5, _ := strconv.Atoi(dlt.F5)
		f6, _ := strconv.Atoi(dlt.F6)
		history = append(history, Draw{[]int{f1, f2, f3, f4, f5, f6}})
	}
	return history
}

// 随机生成分组 —— 使用 groupSizes 控制每组容量
func randomGroup(groupSizes []int) []int {
	g := make([]int, 33)
	nums := rand.Perm(33)

	// 验证 groupSizes 总和
	total := 0
	for _, size := range groupSizes {
		total += size
	}
	if total != 33 {
		panic("groupSizes 错误：五组数量之和必须为 33")
	}

	idx := 0
	for groupID := 0; groupID < 6; groupID++ {
		size := groupSizes[groupID]

		for j := 0; j < size; j++ {
			pos := nums[idx]
			g[pos] = groupID
			idx++
		}
	}
	//fmt.Println(g)
	return g
}

// CalcHits 计算命中（支持指定每组取几个）
func CalcHits(history []Draw, g []int, groupNeed []int) int {
	hits := 0

	for _, d := range history {
		// 统计本期各组命中数量
		groupCnt := make([]int, 6)
		ok := true

		for _, num := range d.Front {
			groupID := g[num-1]
			groupCnt[groupID]++

			// 剪枝：一旦超过需求，直接失败
			if groupCnt[groupID] > groupNeed[groupID] {
				ok = false
				break
			}
		}

		if !ok {
			continue
		}

		// 校验是否“刚好满足每组需求”
		for i := 0; i < 6; i++ {
			if groupCnt[i] != groupNeed[i] {
				ok = false
				break
			}
		}

		if ok {
			hits++
		}
	}

	return hits
}

// 邻域：随机交换两个号码
func mutate(g []int) []int {
	ng := make([]int, 33)
	copy(ng, g)

	i := rand.IntN(33)
	j := rand.IntN(33)
	ng[i], ng[j] = ng[j], ng[i]
	return ng
}

// 模拟退火 —— 完整支持 groupSizes
func simulatedAnnealing(history []Draw, groupSizes []int, groupNeed []int) []int {
	T := 3000.0
	cooling := 0.998
	minTemp := 0.01

	current := randomGroup(groupSizes)
	currentScore := CalcHits(history, current, groupNeed)

	best := make([]int, 33)
	copy(best, current)
	bestScore := currentScore

	for T > minTemp {
		next := mutate(current)
		nextScore := CalcHits(history, next, groupNeed)

		if nextScore >= currentScore {
			current = next
			currentScore = nextScore
		} else {
			diff := float64(nextScore - currentScore)
			prob := rand.Float64()
			if prob < math.Exp(diff/T) {
				current = next
				currentScore = nextScore
			}
		}

		if currentScore > bestScore {
			bestScore = currentScore
			copy(best, current)
		}

		T *= cooling
	}
	return best
}
func formatGroup(g []int) {
	// 先收集分组
	type grpInfo struct {
		nums     []int
		minValue int
	}

	groupList := make([]grpInfo, 6)
	for i := 0; i < 6; i++ {
		groupList[i] = grpInfo{nums: []int{}, minValue: 999}
	}

	for i := 0; i < 33; i++ {
		grp := g[i]
		num := i + 1
		groupList[grp].nums = append(groupList[grp].nums, num)
		if num < groupList[grp].minValue {
			groupList[grp].minValue = num
		}
	}

	// 组内排序（数字升序）
	for i := range groupList {
		sort.Ints(groupList[i].nums)
	}

	// 按你的规则排序：
	// 1. 数组长度降序
	// 2. 最小值升序
	sort.Slice(groupList, func(i, j int) bool {
		if len(groupList[i].nums) != len(groupList[j].nums) {
			return len(groupList[i].nums) > len(groupList[j].nums)
		}
		return groupList[i].minValue < groupList[j].minValue
	})

	labels := []string{"A", "B", "C", "D", "E", "F"}

	// 输出
	for i := 0; i < 6; i++ {
		fmt.Printf("%s: [", labels[i])
		for j, num := range groupList[i].nums {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Printf(`"%02d"`, num)
		}
		fmt.Println("]")
	}
}

func GToAeMs(g []int) map[string]string {
	// 先收集分组
	type grpInfo struct {
		nums     []int
		minValue int
	}

	groupList := make([]grpInfo, 6)
	for i := 0; i < 6; i++ {
		groupList[i] = grpInfo{nums: []int{}, minValue: 999}
	}

	for i := 0; i < 33; i++ {
		grp := g[i]
		num := i + 1
		groupList[grp].nums = append(groupList[grp].nums, num)
		if num < groupList[grp].minValue {
			groupList[grp].minValue = num
		}
	}

	// 组内排序（数字升序）
	for i := range groupList {
		sort.Ints(groupList[i].nums)
	}

	// 按你的规则排序：
	// 1. 数组长度降序
	// 2. 最小值升序
	sort.Slice(groupList, func(i, j int) bool {
		if len(groupList[i].nums) != len(groupList[j].nums) {
			return len(groupList[i].nums) > len(groupList[j].nums)
		}
		return groupList[i].minValue < groupList[j].minValue
	})

	labels := []string{"A", "B", "C", "D", "E", "F"}

	aeMs := make(map[string]string)
	var str string
	// 输出
	for i := 0; i < 6; i++ {
		str = ""
		for j, num := range groupList[i].nums {
			if j > 0 {
				str += fmt.Sprint(",")
			}
			str += fmt.Sprintf(`%02d`, num)
		}
		aeMs[labels[i]] = str
	}

	return aeMs
}

func AeMsToG(aeMs map[string]string) []int {
	gLen := 0
	for _, frontHmStr := range aeMs {
		gLen += len(strings.Split(frontHmStr, ","))
	}
	g := make([]int, gLen)
	for ae, frontHmStr := range aeMs {
		hmStrs := strings.Split(frontHmStr, ",")
		var hms []int
		for _, hmStr := range hmStrs {
			h, _ := strconv.Atoi(hmStr)
			hms = append(hms, h)
		}
		switch ae {
		case "A":
			for _, hm := range hms {
				g[hm-1] = 0
			}
		case "B":
			for _, hm := range hms {
				g[hm-1] = 1
			}
		case "C":
			for _, hm := range hms {
				g[hm-1] = 2
			}
		case "D":
			for _, hm := range hms {
				g[hm-1] = 3
			}
		case "E":
			for _, hm := range hms {
				g[hm-1] = 4
			}
		case "F":
			for _, hm := range hms {
				g[hm-1] = 5
			}
		}
	}
	return g
}

func joinAeMs(m map[string]string) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s", m["A"], m["B"], m["C"], m["D"], m["E"], m["F"])
}

// Moni
//
//	@Description:
//	@param groupSizes
//	@param groupNeed
//	@param spScore
//	@param repeat
//	@param workers
//	@param wg1
func Moni(groupSizes, groupNeed []int, spScore, repeat, workers int, wg1 *sync.WaitGroup) {
	// 检查 groupSizes 长度
	if len(groupSizes) != 6 {
		panic("groupSizes 必须包含 6 个元素（A~F 各一个）")
	}
	if len(groupNeed) != 6 {
		panic("groupNeed 必须包含 6 个元素（A~F 各一个）")
	}
	var typ string
	for _, v := range groupSizes {
		typ += fmt.Sprintf("%d", v)
	}
	var groupNeedStr string
	for _, v := range groupNeed {
		groupNeedStr += fmt.Sprintf("%d", v)
	}
	methodStr := groupNeedStr

	//
	var spScoreAeStrSlice []string

	history := LoadSsqFrontHmHistory()

	var globalBest []int
	globalBestScore := -1

	var mu sync.Mutex
	var wg sync.WaitGroup
	tasks := make(chan int, repeat)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for round := range tasks {
				bestGroup := simulatedAnnealing(history, groupSizes, groupNeed)
				score := CalcHits(history, bestGroup, groupNeed)

				mu.Lock()
				if score > globalBestScore {
					globalBestScore = score
					globalBest = bestGroup
				}

				if score >= spScore {
					// aeMs 已经按照大小进行排序
					aeMs := GToAeMs(bestGroup)
					aeStr := joinAeMs(aeMs)

					if !slices.Contains(spScoreAeStrSlice, aeStr) {
						spScoreAeStrSlice = append(spScoreAeStrSlice, aeStr)
						dbop.UpdateOrInsertSsqMoni(models.SsqMoni{
							A:      aeMs["A"],
							B:      aeMs["B"],
							C:      aeMs["C"],
							D:      aeMs["D"],
							E:      aeMs["E"],
							F:      aeMs["F"],
							Cs:     score,
							Typ:    typ,
							Method: methodStr,
						}, groupNeedStr)
					}
				}

				if round%500 == 0 {
					fmt.Printf("typ=%10s worker=%02d round=%4d，命中：%3d <- globalBestScore=%3d\n", typ,
						workerID, round, score, globalBestScore)
				}
				mu.Unlock()
			}
		}(w)
	}

	for i := 1; i <= repeat; i++ {
		tasks <- i
	}
	close(tasks)

	wg.Wait()

	fmt.Printf("\n 类型 %s ====== 最终全局最优结果 ======\n", typ)
	fmt.Printf("类型 %s 命中次数：%d\n", typ, globalBestScore)
	fmt.Printf("类型 %s 最终最佳组合：%v\n", typ, globalBest)
	formatGroup(globalBest)
	wg1.Done()
}

type MoniDataSt struct {
	GroupSize []int
	GroupNeed []int
	SpScore   int
	Repeat    int
	Workers   int
}

// BatchMoni 批量模拟
//
//	@Description:
//	@param index 最大可以设置为 4
func BatchMoni(prevUpdateMoniTable bool, index int) {
	if prevUpdateMoniTable {
		UpdateMonisTable()
	}
	ds := []MoniDataSt{
		{GroupSize: []int{6, 6, 6, 6, 6, 3}, GroupNeed: []int{1, 1, 1, 1, 1, 1}, SpScore: 155, Repeat: 20000, Workers: 8},
		{GroupSize: []int{8, 5, 5, 5, 5, 5}, GroupNeed: []int{1, 1, 1, 1, 1, 1}, SpScore: 158, Repeat: 20000, Workers: 8},
		{GroupSize: []int{8, 7, 6, 5, 4, 3}, GroupNeed: []int{1, 1, 1, 1, 1, 1}, SpScore: 138, Repeat: 20000, Workers: 8},
		{GroupSize: []int{13, 4, 4, 4, 4, 4}, GroupNeed: []int{1, 1, 1, 1, 1, 1}, SpScore: 104, Repeat: 20000, Workers: 8},
		{GroupSize: []int{18, 3, 3, 3, 3, 3}, GroupNeed: []int{1, 1, 1, 1, 1, 1}, SpScore: 51, Repeat: 20000, Workers: 8},
	}
	var wg sync.WaitGroup
	for i, d := range ds {
		if index != -1 {
			if i <= index {
				wg.Add(1)
				go Moni(d.GroupSize, d.GroupNeed, d.SpScore, d.Repeat, d.Workers, &wg)
			}
		}
	}
	wg.Wait()
	lg.InfoToFileAndStdOut("执行完成!")
}

func UpdateMonisTable() {
	var list []models.SsqMoni

	err := db.DB.Find(&list).Error
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("query monis failed: %v", err), 1)
		return
	}

	his := LoadSsqFrontHmHistory()

	for _, moni := range list {
		groupNeedStr := strings.Split(strings.Trim(moni.Method, " "), "")
		groupNeed := make([]int, len(groupNeedStr))
		for i, v := range groupNeedStr {
			groupNeed[i], _ = strconv.Atoi(v)
		}
		aeMs := make(map[string]string, 5)
		aeMs["A"] = moni.A
		aeMs["B"] = moni.B
		aeMs["C"] = moni.C
		aeMs["D"] = moni.D
		aeMs["E"] = moni.E
		aeMs["F"] = moni.F
		hit := CalcHits(his, AeMsToG(aeMs), groupNeed)
		if hit != moni.Cs {
			fmt.Printf("hit=%d cs=%d\n", hit, moni.Cs)
			err = db.DB.Model(&models.SsqMoni{}).
				Where("id = ?", moni.ID).
				Updates(models.SsqMoni{UpdatedAt: time.Now(), Cs: hit}).Error
			if err != nil {
				lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("update failed, id=%d, err=%v", moni.ID, err), 1)
			}
		}
	}
}
