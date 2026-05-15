package ana_dlt

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

	"github.com/before80/lot/cfg"
	"github.com/before80/lot/db"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
)

type Draw struct {
	Front []int
}

// LoadDltFrontHmHistory 获取所有大乐透开奖号码的前区号码
func LoadDltFrontHmHistory() []Draw {
	dlts, _ := dbop.ReadAllDlt(false)
	history := make([]Draw, 0, len(dlts))
	for _, dlt := range dlts {
		f1, _ := strconv.Atoi(dlt.F1)
		f2, _ := strconv.Atoi(dlt.F2)
		f3, _ := strconv.Atoi(dlt.F3)
		f4, _ := strconv.Atoi(dlt.F4)
		f5, _ := strconv.Atoi(dlt.F5)
		history = append(history, Draw{[]int{f1, f2, f3, f4, f5}})
	}
	return history
}

// randomGroup 随机生成分组 —— 使用 groupSizes 控制每组容量
func randomGroup(groupSizes []int) []int {
	g := make([]int, 35)
	nums := rand.Perm(35)

	// 验证 groupSizes 总和
	total := 0
	for _, size := range groupSizes {
		total += size
	}
	if total != 35 {
		panic("groupSizes 错误：五组数量之和必须为 35")
	}

	idx := 0
	for groupID := 0; groupID < 5; groupID++ {
		size := groupSizes[groupID]

		for j := 0; j < size; j++ {
			pos := nums[idx]
			g[pos] = groupID
			idx++
		}
	}

	return g
}

// 计算命中
//func CalcHits(history []Draw, g []int) int {
//	hits := 0
//
//	for _, d := range history {
//		usedGroup := map[int]bool{} // key：组编号（0~4，对应 A~E）; value：是否已被占用
//		ok := true                  // 本期是否仍然“合法”, 一旦发现 两个号码来自同一组, 立即标记 ok = false
//
//		for _, num := range d.Front {
//			groupID := g[num-1] // 找出该号码属于哪个组, 号码范围是 1~35,数组索引是 0~34,所以用 num-1
//			// 例如:
//			// num = 22
//			// groupID = g[21] = 3   // 说明 22 属于 D 组
//			if usedGroup[groupID] { // 判断是否“同组重复”, 如果这个组 已经被前面的号码用过, 则本期 不可能满足 “5 组各 1 个”, 直接终止本期检查（剪枝，提升性能）
//				ok = false
//				break
//			}
//			usedGroup[groupID] = true // 记录该组已被使用
//		}
//
//		// 本期是否计入命中
//		if ok && len(usedGroup) == 5 {
//			hits++
//		}
//	}
//	return hits
//}

// CalcHits 计算命中（支持指定每组取几个）
func CalcHits(history []Draw, g []int, groupNeed []int) int {
	hits := 0

	for _, d := range history {
		// 统计本期各组命中数量
		groupCnt := make([]int, 5)
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
		for i := 0; i < 5; i++ {
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

// mutate 邻域：随机交换两个号码
func mutate(g []int) []int {
	ng := make([]int, 35)
	copy(ng, g)

	i := rand.IntN(35)
	j := rand.IntN(35)
	ng[i], ng[j] = ng[j], ng[i]
	return ng
}

// simulatedAnnealing 模拟退火 —— 完整支持 groupSizes
func simulatedAnnealing(history []Draw, groupSizes []int, groupNeed []int) []int {
	T := 3000.0
	cooling := 0.998
	minTemp := 0.01

	current := randomGroup(groupSizes)
	currentScore := CalcHits(history, current, groupNeed)

	best := make([]int, 35)
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

	groupList := make([]grpInfo, 5)
	for i := 0; i < 5; i++ {
		groupList[i] = grpInfo{nums: []int{}, minValue: 999}
	}

	for i := 0; i < 35; i++ {
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

	labels := []string{"A", "B", "C", "D", "E"}

	// 输出
	for i := 0; i < 5; i++ {
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

	groupList := make([]grpInfo, 5)
	for i := 0; i < 5; i++ {
		groupList[i] = grpInfo{nums: []int{}, minValue: 999}
	}

	for i := 0; i < 35; i++ {
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

	labels := []string{"A", "B", "C", "D", "E"}

	aeMs := make(map[string]string)
	var str string
	// 输出
	for i := 0; i < 5; i++ {
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
		}
	}
	return g
}

func joinAeMs(m map[string]string) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", m["A"], m["B"], m["C"], m["D"], m["E"])
}

// Moni 模拟大乐透开奖
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
	if len(groupSizes) != 5 {
		panic("groupSizes 必须包含 5 个元素（A~E 各一个）")
	}
	if len(groupNeed) != 5 {
		panic("groupNeed 必须包含 5 个元素（A~E 各一个）")
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

	history := LoadDltFrontHmHistory()

	var globalBest []int
	_ = globalBest
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
						dbop.UpdateOrInsertDltMoni(models.DltMoni{
							A:      aeMs["A"],
							B:      aeMs["B"],
							C:      aeMs["C"],
							D:      aeMs["D"],
							E:      aeMs["E"],
							Cs:     score,
							Typ:    typ,
							Method: methodStr,
						}, groupNeedStr)
					}
				}

				if round%1000 == 0 {
					//fmt.Printf("typ=%10s worker=%02d round=%4d，命中：%3d <- globalBestScore=%3d\n", typ, workerID, round, score, globalBestScore)
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

	//fmt.Printf("\n 类型 %s ====== 最终全局最优结果 ======\n", typ)
	//fmt.Printf("类型 %s 命中次数：%d\n", typ, globalBestScore)
	//fmt.Printf("类型 %s 最终最佳组合：%v\n", typ, globalBest)
	//formatGroup(globalBest)
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
//	@param index 最大可以设置为 9
func BatchMoni(index int) {
	var ds []MoniDataSt

	if cfg.Default.EnvIsLocal == 1 {
		ds = []MoniDataSt{
			{GroupSize: []int{7, 7, 7, 7, 7}, GroupNeed: []int{1, 1, 1, 1, 1}, SpScore: 258, Repeat: 20000, Workers: 8},
			{GroupSize: []int{11, 6, 6, 6, 6}, GroupNeed: []int{1, 1, 1, 1, 1}, SpScore: 233, Repeat: 20000, Workers: 8},
			{GroupSize: []int{15, 5, 5, 5, 5}, GroupNeed: []int{1, 1, 1, 1, 1}, SpScore: 172, Repeat: 20000, Workers: 8},
			{GroupSize: []int{19, 4, 4, 4, 4}, GroupNeed: []int{1, 1, 1, 1, 1}, SpScore: 112, Repeat: 20000, Workers: 8},
			{GroupSize: []int{23, 3, 3, 3, 3}, GroupNeed: []int{1, 1, 1, 1, 1}, SpScore: 59, Repeat: 20000, Workers: 4},
			{GroupSize: []int{23, 3, 3, 3, 3}, GroupNeed: []int{5, 0, 0, 0, 0}, SpScore: 420, Repeat: 20000, Workers: 4},
			{GroupSize: []int{24, 4, 3, 2, 2}, GroupNeed: []int{5, 0, 0, 0, 0}, SpScore: 507, Repeat: 20000, Workers: 4},
			{GroupSize: []int{25, 4, 3, 2, 1}, GroupNeed: []int{5, 0, 0, 0, 0}, SpScore: 620, Repeat: 20000, Workers: 4},
			{GroupSize: []int{26, 3, 2, 2, 2}, GroupNeed: []int{5, 0, 0, 0, 0}, SpScore: 726, Repeat: 20000, Workers: 4},
			{GroupSize: []int{27, 2, 2, 2, 2}, GroupNeed: []int{5, 0, 0, 0, 0}, SpScore: 860, Repeat: 20000, Workers: 4},
		}
	} else {
		ds = []MoniDataSt{
			{GroupSize: []int{7, 7, 7, 7, 7}, GroupNeed: []int{1, 1, 1, 1, 1}, SpScore: 258, Repeat: 200, Workers: 1},
			{GroupSize: []int{11, 6, 6, 6, 6}, GroupNeed: []int{1, 1, 1, 1, 1}, SpScore: 233, Repeat: 200, Workers: 1},
			{GroupSize: []int{15, 5, 5, 5, 5}, GroupNeed: []int{1, 1, 1, 1, 1}, SpScore: 172, Repeat: 200, Workers: 1},
			{GroupSize: []int{19, 4, 4, 4, 4}, GroupNeed: []int{1, 1, 1, 1, 1}, SpScore: 112, Repeat: 200, Workers: 1},
			{GroupSize: []int{23, 3, 3, 3, 3}, GroupNeed: []int{1, 1, 1, 1, 1}, SpScore: 59, Repeat: 200, Workers: 1},
			{GroupSize: []int{23, 3, 3, 3, 3}, GroupNeed: []int{5, 0, 0, 0, 0}, SpScore: 420, Repeat: 200, Workers: 1},
			{GroupSize: []int{24, 4, 3, 2, 2}, GroupNeed: []int{5, 0, 0, 0, 0}, SpScore: 507, Repeat: 200, Workers: 1},
			{GroupSize: []int{25, 4, 3, 2, 1}, GroupNeed: []int{5, 0, 0, 0, 0}, SpScore: 620, Repeat: 200, Workers: 1},
			{GroupSize: []int{26, 3, 2, 2, 2}, GroupNeed: []int{5, 0, 0, 0, 0}, SpScore: 726, Repeat: 200, Workers: 1},
			{GroupSize: []int{27, 2, 2, 2, 2}, GroupNeed: []int{5, 0, 0, 0, 0}, SpScore: 860, Repeat: 200, Workers: 1},
		}
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
	lg.InfoToFile("模拟执行完成!")
}

// UpdateMonisTable
func UpdateMonisTable() {
	var list []models.DltMoni

	err := db.DB.Find(&list).Error
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("query monis failed: %v", err), 1)
		return
	}

	his := LoadDltFrontHmHistory()

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
		hit := CalcHits(his, AeMsToG(aeMs), groupNeed)
		if hit != moni.Cs {
			//fmt.Printf("hit=%d cs=%d\n", hit, moni.Cs)
			err = db.DB.Model(&models.DltMoni{}).
				Where("id = ?", moni.ID).
				Updates(models.DltMoni{UpdatedAt: time.Now(), Cs: hit}).Error
			if err != nil {
				lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("update failed, id=%d, err=%v", moni.ID, err), 1)
			}
		}
	}
}
