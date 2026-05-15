## 版本1

```go
package ana

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/before80/lot/dbop"
)

type Issue [5]int // 前区 5 个号码

// 深拷贝 group
func cloneGroup(src map[int]int) map[int]int {
	dst := make(map[int]int, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// 计算命中期数
func countMatches(history []Issue, group map[int]int) int {
	hit := 0
	for _, issue := range history {
		seen := make(map[int]bool, 5)
		ok := true
		for _, n := range issue {
			g := group[n]
			if seen[g] {
				ok = false
				break
			}
			seen[g] = true
		}
		if ok {
			hit++
		}
	}
	return hit
}

// 随机生成有效分组（每组 7 个），使用传入 rng 保证随机性
func randomGroup(rng *rand.Rand) map[int]int {
	perm := rng.Perm(35)
	group := make(map[int]int, 35)
	for g := 0; g < 5; g++ {
		for i := 0; i < 7; i++ {
			group[perm[g*7+i]+1] = g
		}
	}
	return group
}

// 对已有分组尝试交换两个号码，若命中提高则接受（在副本 current 上操作）
// 返回改进后的分组和该副本的 hits
func improveGroup(history []Issue, group map[int]int, rng *rand.Rand, tries int) (map[int]int, int) {
	current := cloneGroup(group)
	bestHits := countMatches(history, current)

	for iter := 0; iter < tries; iter++ {
		a := rng.Intn(35) + 1
		b := rng.Intn(35) + 1
		if a == b {
			continue
		}
		if current[a] == current[b] {
			continue
		}

		// 交换并测试
		ga, gb := current[a], current[b]
		current[a], current[b] = gb, ga

		hits := countMatches(history, current)
		if hits > bestHits {
			bestHits = hits
			// 保持交换（current 已是交换后的状态），继续在此基础上寻找更多改进
		} else {
			// 回滚
			current[a], current[b] = ga, gb
		}
	}
	return current, bestHits
}

func Moni() {
	// 建立独立 RNG（确保每次运行差异）
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// ======= 在这里放入你的历史号码（2807 期） ========
	dlts, _ := dbop.ReadAllDlt(false)
	history := make([]Issue, 0, len(dlts))
	for _, dlt := range dlts {
		f1, _ := strconv.Atoi(dlt.F1)
		f2, _ := strconv.Atoi(dlt.F2)
		f3, _ := strconv.Atoi(dlt.F3)
		f4, _ := strconv.Atoi(dlt.F4)
		f5, _ := strconv.Atoi(dlt.F5)
		history = append(history, Issue{f1, f2, f3, f4, f5})
	}

	// 参数
	rounds := 200             // 主循环轮数
	triesPerImprove := 200000 // improveGroup 内部尝试次数
	noImproveLimit := 2       // 连续多少轮没改进就重启
	bestGroup := randomGroup(rng)
	bestHits := countMatches(history, bestGroup)

	fmt.Println("初始 hits:", bestHits)

	// 全局最佳
	globalBestHits := bestHits
	globalBestGroup := cloneGroup(bestGroup)

	noImproveCount := 0

	for round := 0; round < rounds; round++ {
		newGroup, newHits := improveGroup(history, bestGroup, rng, triesPerImprove)

		fmt.Printf("Round %3d: tryHits=%d, bestHits=%d, globalBest=%d\n",
			round+1, newHits, bestHits, globalBestHits)

		if newHits > bestHits {
			bestHits = newHits
			bestGroup = newGroup
			noImproveCount = 0

			if newHits > globalBestHits {
				globalBestHits = newHits
				globalBestGroup = cloneGroup(newGroup)
			}

			fmt.Printf("  >> Improvement! bestHits=%d globalBest=%d\n", bestHits, globalBestHits)

		} else {
			noImproveCount++
		}

		// 只重置“当前搜索状态”，不要覆盖全局最佳
		if noImproveCount >= noImproveLimit {
			fmt.Printf("  -- no improvement, restart search\n")

			bestGroup = randomGroup(rng)
			bestHits = countMatches(history, bestGroup)

			fmt.Printf("  -- restart bestHits=%d, globalBest=%d\n", bestHits, globalBestHits)
			noImproveCount = 0
		}
	}

	// 最终结果输出
	fmt.Println("最终最佳命中期数：", globalBestHits)

	groups := make([][]int, 5)
	for n := 1; n <= 35; n++ {
		g := globalBestGroup[n]
		groups[g] = append(groups[g], n)
	}
	for i := 0; i < 5; i++ {
		fmt.Printf("组 %c: %v\n", 'A'+i, groups[i])
	}
}

```

## 版本2：
```go
package ana

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/before80/lot/dbop"
)

type Draw struct {
	Front []int
}

func loadHistory() []Draw {
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

func randomGroup() []int {
	g := make([]int, 35)
	nums := rand.Perm(35) // rand/v2 的 Perm 仍然可用

	groupSize := 7
	for gi := 0; gi < 5; gi++ {
		for j := 0; j < groupSize; j++ {
			idx := nums[gi*groupSize+j]
			g[idx] = gi
		}
	}

	return g
}

// 判断历史号码是否能被该分组方案构成（每组恰好选1个）
func calcHits(history []Draw, g []int) int {
	hits := 0

	for _, d := range history {
		usedGroup := map[int]bool{}
		ok := true

		for _, num := range d.Front {
			groupID := g[num-1]
			if usedGroup[groupID] {
				ok = false
				break
			}
			usedGroup[groupID] = true
		}
		if ok && len(usedGroup) == 5 {
			hits++
		}
	}

	return hits
}

// 建立邻域：随机交换两个号码的组别
func mutate(g []int) []int {
	ng := make([]int, 35)
	copy(ng, g)

	i := rand.IntN(35) // rand/v2 用 IntN
	j := rand.IntN(35)
	ng[i], ng[j] = ng[j], ng[i]
	return ng
}

// 模拟退火
func simulatedAnnealing(history []Draw) []int {
	T := 3000.0
	cooling := 0.998
	minTemp := 0.01

	current := randomGroup()
	currentScore := calcHits(history, current)

	best := make([]int, 35)
	copy(best, current)
	bestScore := currentScore

	for T > minTemp {
		next := mutate(current)
		nextScore := calcHits(history, next)

		// 若更好，接受
		if nextScore >= currentScore {
			current = next
			currentScore = nextScore
		} else {
			diff := float64(nextScore - currentScore)
			prob := rand.Float64() // rand/v2 仍然支持 Float64()

			if prob < exp(diff/T) {
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

func exp(x float64) float64 { return math.Exp(x) }

// 把 bestGroup 映射成 A,B,C,D,E 五组号码
func formatGroup(g []int) {
	groups := make([][]int, 5)
	for i := 0; i < 5; i++ {
		groups[i] = []int{}
	}

	for i := 0; i < 35; i++ {
		groupID := g[i]
		groups[groupID] = append(groups[groupID], i+1)
	}

	labels := []string{"A", "B", "C", "D", "E"}

	for i := 0; i < 5; i++ {
		fmt.Printf("%s: [", labels[i])
		for j, num := range groups[i] {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Printf(`"%02d"`, num)
		}
		fmt.Println("]")
	}
}

func Moni() {
	// rand/v2 不再支持 Seed，也不需要 Seed
	// 移除 rand.Seed(...)

	history := loadHistory()

	repeat := 1000
	var globalBest []int
	globalBestScore := -1

	for i := 1; i <= repeat; i++ {
		bestGroup := simulatedAnnealing(history)
		score := calcHits(history, bestGroup)

		if score > globalBestScore {
			globalBestScore = score
			globalBest = bestGroup
		}

		fmt.Printf("round=%10d，当前结果命中：%10d <- globalBestScore=%5d\n",
			i, score, globalBestScore)
	}

	fmt.Println("\n====== 最终全局最优结果 ======")
	fmt.Println("命中次数：", globalBestScore)
	fmt.Println("最终最佳组合（组别数组）：", globalBest)
	fmt.Println("格式化分组如下：")
	formatGroup(globalBest)
}


```



## 版本3（多线程，固定7个号）

```go
package ana

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"sync"

	"github.com/before80/lot/dbop"
)

type Draw struct {
	Front []int
}

func loadHistory() []Draw {
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

func randomGroup() []int {
	g := make([]int, 35)
	nums := rand.Perm(35) // rand/v2 的 Perm
	groupSize := 7

	for gi := 0; gi < 5; gi++ {
		for j := 0; j < groupSize; j++ {
			idx := nums[gi*groupSize+j]
			g[idx] = gi
		}
	}
	return g
}

// 判断历史号码是否能被该分组构成
func calcHits(history []Draw, g []int) int {
	hits := 0

	for _, d := range history {
		usedGroup := map[int]bool{}
		ok := true

		for _, num := range d.Front {
			groupID := g[num-1]
			if usedGroup[groupID] {
				ok = false
				break
			}
			usedGroup[groupID] = true
		}

		if ok && len(usedGroup) == 5 {
			hits++
		}
	}
	return hits
}

// 邻域：随机交换两个号码
func mutate(g []int) []int {
	ng := make([]int, 35)
	copy(ng, g)

	i := rand.IntN(35)
	j := rand.IntN(35)
	ng[i], ng[j] = ng[j], ng[i]
	return ng
}

func simulatedAnnealing(history []Draw) []int {
	T := 3000.0
	cooling := 0.998
	minTemp := 0.01

	current := randomGroup()
	currentScore := calcHits(history, current)

	best := make([]int, 35)
	copy(best, current)
	bestScore := currentScore

	for T > minTemp {
		next := mutate(current)
		nextScore := calcHits(history, next)

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
	groups := make([][]int, 5)
	for i := 0; i < 5; i++ {
		groups[i] = []int{}
	}

	for i := 0; i < 35; i++ {
		groupID := g[i]
		groups[groupID] = append(groups[groupID], i+1)
	}

	labels := []string{"A", "B", "C", "D", "E"}

	for i := 0; i < 5; i++ {
		fmt.Printf("%s: [", labels[i])
		for j, num := range groups[i] {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Printf(`"%02d"`, num)
		}
		fmt.Println("]")
	}
}

func Moni() {
	history := loadHistory()

	repeat := 1000
	workers := 30 // 多线程数量，可自由调整

	var globalBest []int
	globalBestScore := -1

	var mu sync.Mutex
	var wg sync.WaitGroup
	tasks := make(chan int, repeat)

	// 启动 worker
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for round := range tasks {
				bestGroup := simulatedAnnealing(history)
				score := calcHits(history, bestGroup)

				mu.Lock()
				if score > globalBestScore {
					globalBestScore = score
					globalBest = bestGroup
				}
				fmt.Printf("worker=%02d round=%4d，命中：%3d <- globalBestScore=%3d\n",
					workerID, round, score, globalBestScore)
				mu.Unlock()
			}
		}(w)
	}

	// 分发任务
	for i := 1; i <= repeat; i++ {
		tasks <- i
	}
	close(tasks)

	// 等待 goroutine 结束
	wg.Wait()

	fmt.Println("\n====== 最终全局最优结果 ======")
	fmt.Println("命中次数：", globalBestScore)
	fmt.Println("最终最佳组合：", globalBest)
	formatGroup(globalBest)
}

```





## 版本4（多线程，可配置）

```go
package ana

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"sort"
	"strconv"
	"sync"

	"github.com/before80/lot/dbop"
	"github.com/before80/lot/models"
)

type Draw struct {
	Front []int
}

func loadHistory() []Draw {
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

// 随机生成分组 —— 使用 groupSizes 控制每组容量
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
func calcHits(history []Draw, g []int) int {
	hits := 0

	for _, d := range history {
		used := map[int]bool{}
		ok := true

		for _, num := range d.Front {
			grp := g[num-1]
			if used[grp] {
				ok = false
				break
			}
			used[grp] = true
		}

		if ok && len(used) == 5 {
			hits++
		}
	}
	return hits
}

// 邻域：随机交换两个号码
func mutate(g []int) []int {
	ng := make([]int, 35)
	copy(ng, g)

	i := rand.IntN(35)
	j := rand.IntN(35)
	ng[i], ng[j] = ng[j], ng[i]
	return ng
}

// 模拟退火 —— 完整支持 groupSizes
func simulatedAnnealing(history []Draw, groupSizes []int) []int {
	T := 3000.0
	cooling := 0.998
	minTemp := 0.01

	current := randomGroup(groupSizes)
	currentScore := calcHits(history, current)

	best := make([]int, 35)
	copy(best, current)
	bestScore := currentScore

	for T > minTemp {
		next := mutate(current)
		nextScore := calcHits(history, next)

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

func formatGroup1(g []int) map[string]string {
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

func joinAeMs(m map[string]string) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", m["A"], m["B"], m["C"], m["D"], m["E"])
}

// Moni —— 新版本传入 groupSizes
func Moni(groupSizes []int, spScore int, wg1 *sync.WaitGroup) {
	// 检查 groupSizes 长度
	if len(groupSizes) != 5 {
		panic("groupSizes 必须包含 5 个元素（A~E 各一个）")
	}
	var typ string
	for _, v := range groupSizes {
		typ += fmt.Sprintf("%d", v)
	}

	//
	var spScoreAeStrSlice []string

	history := loadHistory()

	repeat := 20000
	workers := 30

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
				bestGroup := simulatedAnnealing(history, groupSizes)
				score := calcHits(history, bestGroup)

				mu.Lock()
				if score > globalBestScore {
					globalBestScore = score
					globalBest = bestGroup
				}

				if score >= spScore {
					// aeMs 已经按照大小进行排序
					aeMs := formatGroup1(bestGroup)
					aeStr := joinAeMs(aeMs)

					if !slices.Contains(spScoreAeStrSlice, aeStr) {
						spScoreAeStrSlice = append(spScoreAeStrSlice, aeStr)
						dbop.UpdateOrInsertMoni(models.Moni{
							A:   aeMs["A"],
							B:   aeMs["B"],
							C:   aeMs["C"],
							D:   aeMs["D"],
							E:   aeMs["E"],
							Cs:  score,
							Typ: typ,
						})
					}
				}

				if round%1000 == 0 {
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

```



## 遗传学版本1：

```go
package ana

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/before80/lot/dbop"
)

type Draw struct {
	Front []int
}

func loadHistory() []Draw {
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

// =============================
// 配置：每组容量
// =============================
var groupSizes = []int{7, 7, 7, 7, 7} // 示例，可修改

type Individual struct {
	Gene  []int // 组别数组
	Score int
}

// =============================
// 初始化随机个体
// =============================
func newRandomIndividual(history []Draw) Individual {
	g := make([]int, 35)
	nums := rand.Perm(35)

	pos := 0
	for groupID, size := range groupSizes {
		for j := 0; j < size; j++ {
			g[nums[pos]] = groupID
			pos++
		}
	}

	return Individual{
		Gene:  g,
		Score: calcHits(history, g),
	}
}

// =============================
// 选择：锦标赛选择
// =============================
func tournamentSelect(pop []Individual, k int) Individual {
	best := pop[rand.IntN(len(pop))]
	for i := 1; i < k; i++ {
		ind := pop[rand.IntN(len(pop))]
		if ind.Score > best.Score {
			best = ind
		}
	}
	return best
}

// =============================
// 交叉（保组容量）：按片段交换 + 修复
// =============================
func crossover(p1, p2 Individual) Individual {
	child := make([]int, 35)

	// 单点交叉
	cut := rand.IntN(35)
	for i := 0; i < 35; i++ {
		if i < cut {
			child[i] = p1.Gene[i]
		} else {
			child[i] = p2.Gene[i]
		}
	}

	// 修复（关键步骤，保证 groupSizes 不变）
	child = fixGroupSizes(child)

	return Individual{
		Gene: child,
	}
}

// =============================
// 修复染色体：保证每组数量与 groupSizes 一致
// =============================
func fixGroupSizes(g []int) []int {
	count := make([]int, 5)
	for _, gid := range g {
		count[gid]++
	}

	// 找出缺少和多余的组别
	extra := []int{}
	need := []int{}

	for gid := range groupSizes {
		if count[gid] > groupSizes[gid] {
			for i := 0; i < count[gid]-groupSizes[gid]; i++ {
				extra = append(extra, gid)
			}
		} else if count[gid] < groupSizes[gid] {
			for i := 0; i < groupSizes[gid]-count[gid]; i++ {
				need = append(need, gid)
			}
		}
	}

	// 修复：把 extra 的位置换成 need
	ei := 0
	ni := 0
	for i := 0; i < 35 && ei < len(extra) && ni < len(need); i++ {
		if g[i] == extra[ei] {
			g[i] = need[ni]
			ei++
			ni++
		}
	}

	return g
}

// =============================
// 变异：交换两个不同组的号码
// =============================
func mutate(ind Individual, mutationRate float32) Individual {
	if rand.Float32() > mutationRate {
		return ind
	}

	g := make([]int, 35)
	copy(g, ind.Gene)

	var i, j int
	for {
		i = rand.IntN(35)
		j = rand.IntN(35)
		if g[i] != g[j] {
			break
		}
	}
	g[i], g[j] = g[j], g[i]

	ind.Gene = g
	return ind
}

// =============================
// GA 主流程
// =============================
func runGA(history []Draw) Individual {

	popSize := 1688               // 种群数量 -> 越大越稳，但越慢
	generations := 1000           // 进化代数 -> 越大越强烈寻找最优
	mutationRate := float32(0.05) // 5% 变异概率 -> 0.03~0.10
	tournamentSize := 7           // 3~7

	// 初始化种群
	pop := make([]Individual, popSize)
	for i := 0; i < popSize; i++ {
		pop[i] = newRandomIndividual(history)
	}

	// 全局最优
	best := pop[0]

	// GA 主循环
	for gen := 1; gen <= generations; gen++ {

		newPop := make([]Individual, 0, popSize)
		fmt.Printf("第%d代 ->", gen)
		// 精英保留策略
		for _, ind := range pop {
			if ind.Score > best.Score {
				best = ind
				fmt.Printf("\n🌟 第 %d 代出现更优：%d\n", gen, best.Score)
			}
		}
		newPop = append(newPop, best)

		// 产生新个体
		for len(newPop) < popSize {

			// 选择父母
			p1 := tournamentSelect(pop, tournamentSize)
			p2 := tournamentSelect(pop, tournamentSize)

			// 交叉
			child := crossover(p1, p2)

			// 变异
			child = mutate(child, mutationRate)

			// 计算适应度
			child.Score = calcHits(history, child.Gene)

			newPop = append(newPop, child)
		}

		pop = newPop
	}

	return best
}

// 判断历史号码是否能被该分组方案构成（每组恰好选1个）
func calcHits(history []Draw, g []int) int {
	hits := 0

	for _, d := range history {
		usedGroup := map[int]bool{}
		ok := true

		for _, num := range d.Front {
			groupID := g[num-1]
			if usedGroup[groupID] {
				ok = false
				break
			}
			usedGroup[groupID] = true
		}
		if ok && len(usedGroup) == 5 {
			hits++
		}
	}

	return hits
}

// 把 bestGroup 映射成 A,B,C,D,E 五组号码
func formatGroup(g []int) {
	groups := make([][]int, 5)
	for i := 0; i < 5; i++ {
		groups[i] = []int{}
	}

	for i := 0; i < 35; i++ {
		groupID := g[i]
		groups[groupID] = append(groups[groupID], i+1)
	}

	labels := []string{"A", "B", "C", "D", "E"}

	for i := 0; i < 5; i++ {
		fmt.Printf("%s: [", labels[i])
		for j, num := range groups[i] {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Printf(`"%02d"`, num)
		}
		fmt.Println("]")
	}
}

func Moni() {
	history := loadHistory()

	fmt.Println("开始遗传算法搜索…")

	best := runGA(history)

	fmt.Println("\n====== GA 最终最优结果 ======")
	fmt.Println("命中次数：", best.Score)
	fmt.Println("最终最佳组合（组别数组）：", best.Gene)
	fmt.Println("格式化分组如下：")
	formatGroup(best.Gene)
}


```

## 遗传学版本2(多线程)：

```go
package ana

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/before80/lot/dbop"
)

type Draw struct {
	Front []int
}

func loadHistory() []Draw {
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

// =============================
// 配置：每组容量
// =============================
// var groupSizes = []int{7, 7, 7, 7, 7} // 示例，可修改
var groupSizes = []int{11, 6, 6, 6, 6} // 示例，可修改

type Individual struct {
	Gene  []int // 组别数组
	Score int
}

// =============================
// 初始化随机个体
// =============================
func newRandomIndividual(history []Draw) Individual {
	g := make([]int, 35)
	nums := rand.Perm(35)

	pos := 0
	for groupID, size := range groupSizes {
		for j := 0; j < size; j++ {
			g[nums[pos]] = groupID
			pos++
		}
	}

	return Individual{
		Gene:  g,
		Score: calcHits(history, g),
	}
}

// =============================
// 选择：锦标赛选择
// =============================
func tournamentSelect(pop []Individual, k int) Individual {
	best := pop[rand.IntN(len(pop))]
	for i := 1; i < k; i++ {
		ind := pop[rand.IntN(len(pop))]
		if ind.Score > best.Score {
			best = ind
		}
	}
	return best
}

// =============================
// 交叉（保组容量）：按片段交换 + 修复
// =============================
func crossover(p1, p2 Individual) Individual {
	child := make([]int, 35)

	// 单点交叉
	cut := rand.IntN(35)
	for i := 0; i < 35; i++ {
		if i < cut {
			child[i] = p1.Gene[i]
		} else {
			child[i] = p2.Gene[i]
		}
	}

	// 修复（关键步骤，保证 groupSizes 不变）
	child = fixGroupSizes(child)

	return Individual{
		Gene: child,
	}
}

// =============================
// 修复染色体：保证每组数量与 groupSizes 一致
// =============================
func fixGroupSizes(g []int) []int {
	count := make([]int, 5)
	for _, gid := range g {
		count[gid]++
	}

	// 找出缺少和多余的组别
	extra := []int{}
	need := []int{}

	for gid := range groupSizes {
		if count[gid] > groupSizes[gid] {
			for i := 0; i < count[gid]-groupSizes[gid]; i++ {
				extra = append(extra, gid)
			}
		} else if count[gid] < groupSizes[gid] {
			for i := 0; i < groupSizes[gid]-count[gid]; i++ {
				need = append(need, gid)
			}
		}
	}

	// 修复：把 extra 的位置换成 need
	ei := 0
	ni := 0
	for i := 0; i < 35 && ei < len(extra) && ni < len(need); i++ {
		if g[i] == extra[ei] {
			g[i] = need[ni]
			ei++
			ni++
		}
	}

	return g
}

// =============================
// 变异：交换两个不同组的号码
// =============================
func mutate(ind Individual, mutationRate float32) Individual {
	if rand.Float32() > mutationRate {
		return ind
	}

	g := make([]int, 35)
	copy(g, ind.Gene)

	var i, j int
	for {
		i = rand.IntN(35)
		j = rand.IntN(35)
		if g[i] != g[j] {
			break
		}
	}
	g[i], g[j] = g[j], g[i]

	ind.Gene = g
	return ind
}

// 并行计算所有个体的适应度
func evaluatePopulationParallel(pop []Individual, history []Draw, workers int) {

	type job struct {
		idx  int
		gene []int
	}

	jobs := make(chan job, len(pop))
	results := make(chan struct{}, len(pop))

	// worker
	for w := 0; w < workers; w++ {
		go func() {
			for j := range jobs {
				pop[j.idx].Score = calcHits(history, j.gene)
				results <- struct{}{}
			}
		}()
	}

	// 投递任务
	for i := range pop {
		jobs <- job{
			idx:  i,
			gene: pop[i].Gene,
		}
	}
	close(jobs)

	// 等待全部完成
	for i := 0; i < len(pop); i++ {
		<-results
	}
}

// =============================
// GA 主流程
// =============================
func runGA(history []Draw) Individual {

	popSize := 588                // 种群数量 -> 越大越稳，但越慢
	generations := 5000           // 进化代数 -> 越大越强烈寻找最优
	mutationRate := float32(0.09) // 5% 变异概率 -> 0.03~0.10
	tournamentSize := 3           // 3~7

	// 启用 CPU 核心数（根据你的机器自动设置）
	workers := 20 // 你也可以 runtime.NumCPU()

	// 初始化种群
	pop := make([]Individual, popSize)
	for i := 0; i < popSize; i++ {
		pop[i] = newRandomIndividual(history)
	}

	// 并行计算初始适应度
	evaluatePopulationParallel(pop, history, workers)

	// 全局最优
	best := pop[0]

	// GA 主循环
	for gen := 1; gen <= generations; gen++ {

		if gen%500 == 0 {
			fmt.Printf("第%d代 ->", gen)
		}

		newPop := make([]Individual, 0, popSize)

		// 精英保留
		for _, ind := range pop {
			if ind.Score > best.Score {
				best = ind
				fmt.Printf("\n🌟 第 %d 代出现更优：%d\n", gen, best.Score)
			}
		}
		newPop = append(newPop, best)

		// 产生其他个体
		for len(newPop) < popSize {

			p1 := tournamentSelect(pop, tournamentSize)
			p2 := tournamentSelect(pop, tournamentSize)

			child := crossover(p1, p2)
			child = mutate(child, mutationRate)

			newPop = append(newPop, child)
		}

		// 并行计算新种群的 fitness
		evaluatePopulationParallel(newPop, history, workers)

		pop = newPop
	}

	return best
}

// 判断历史号码是否能被该分组方案构成（每组恰好选1个）
func calcHits(history []Draw, g []int) int {
	hits := 0

	for _, d := range history {
		usedGroup := map[int]bool{}
		ok := true

		for _, num := range d.Front {
			groupID := g[num-1]
			if usedGroup[groupID] {
				ok = false
				break
			}
			usedGroup[groupID] = true
		}
		if ok && len(usedGroup) == 5 {
			hits++
		}
	}

	return hits
}

// 把 bestGroup 映射成 A,B,C,D,E 五组号码
func formatGroup(g []int) {
	groups := make([][]int, 5)
	for i := 0; i < 5; i++ {
		groups[i] = []int{}
	}

	for i := 0; i < 35; i++ {
		groupID := g[i]
		groups[groupID] = append(groups[groupID], i+1)
	}

	labels := []string{"A", "B", "C", "D", "E"}

	for i := 0; i < 5; i++ {
		fmt.Printf("%s: [", labels[i])
		for j, num := range groups[i] {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Printf(`"%02d"`, num)
		}
		fmt.Println("]")
	}
}

func Moni() {
	history := loadHistory()

	fmt.Println("开始遗传算法搜索…")

	best := runGA(history)

	fmt.Println("\n====== GA 最终最优结果 ======")
	fmt.Println("命中次数：", best.Score)
	fmt.Println("最终最佳组合（组别数组）：", best.Gene)
	fmt.Println("格式化分组如下：")
	formatGroup(best.Gene)
}


```

遗传学版本：加强版单线程
```go
package ana

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/before80/lot/dbop"
)

type Draw struct {
	Front []int
}

const (
	GroupCount   = 5
	NumsPerGroup = 7
	TotalNums    = 35
)

func loadHistory() []Draw {
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

// =============================
// 配置：每组容量
// =============================
// var groupSizes = []int{7, 7, 7, 7, 7} // 示例，可修改
var groupSizes = []int{11, 6, 6, 6, 6} // 示例，可修改

// Individual 一个个体 = 一个 35 长度数组，每个值属于 {0..4}
type Individual struct {
	Gene  [TotalNums]int // 组别数组
	Score int
}

// =============================
// 初始化随机个体
// =============================
func newValidIndividual() Individual {
	var ind Individual

	// 将 0..34 随机排列
	nums := rand.Perm(TotalNums)

	// 前 7 个 → 组 0
	// 下 7 个 → 组 1
	// ...
	for g := 0; g < GroupCount; g++ {
		for i := 0; i < NumsPerGroup; i++ {
			idx := nums[g*NumsPerGroup+i]
			ind.Gene[idx] = g
		}
	}

	return ind
}

func fitness(ind *Individual, history []Draw) int {
	count := 0

	for _, d := range history {
		found := 0
		groupUsed := make(map[int]bool)

		for _, num := range d.Front {
			g := ind.Gene[num-1]

			if groupUsed[g] {
				found = -99
				break
			}
			groupUsed[g] = true
			found++
		}

		if found == 5 {
			count++
		}
	}

	ind.Score = count
	return count
}

func localImprove(ind *Individual, history []Draw) {
	bestScore := ind.Score

	for i := 0; i < 3; i++ {
		// 随机两个号码交换
		a := rand.IntN(TotalNums)
		b := rand.IntN(TotalNums)

		ind.Gene[a], ind.Gene[b] = ind.Gene[b], ind.Gene[a]

		newScore := fitness(ind, history)

		if newScore <= bestScore {
			// 恢复
			ind.Gene[a], ind.Gene[b] = ind.Gene[b], ind.Gene[a]
			ind.Score = bestScore
		} else {
			bestScore = newScore
		}
	}
}

func crossover(p1, p2 Individual) Individual {
	var child Individual

	// 每个 group 取一半的号码来自 p1，一半来自 p2
	count := [GroupCount]int{}

	for num := 0; num < TotalNums; num++ {

		g := p1.Gene[num]
		if count[g] < NumsPerGroup/2 {
			child.Gene[num] = g
			count[g]++
		} else {
			// 来自 p2
			g2 := p2.Gene[num]
			child.Gene[num] = g2
			count[g2]++
		}
	}

	// 修正（如果因冲突导致某组不满）
	fixGroups(&child)

	return child
}

func fixGroups(ind *Individual) {
	groupCount := [GroupCount]int{}
	for n := 0; n < TotalNums; n++ {
		groupCount[ind.Gene[n]]++
	}

	// 不满的组
	var need []int
	var extra []int

	for g := 0; g < GroupCount; g++ {
		if groupCount[g] < NumsPerGroup {
			for i := 0; i < NumsPerGroup-groupCount[g]; i++ {
				need = append(need, g)
			}
		} else if groupCount[g] > NumsPerGroup {
			for i := 0; i < groupCount[g]-NumsPerGroup; i++ {
				extra = append(extra, g)
			}
		}
	}

	ei := 0
	for n := 0; n < TotalNums && ei < len(extra); n++ {
		if ind.Gene[n] == extra[ei] {
			ind.Gene[n] = need[ei]
			ei++
		}
	}
}

func mutate(ind *Individual, rate float64) Individual {
	for n := 0; n < TotalNums; n++ {
		if rand.Float64() < rate {
			oldG := ind.Gene[n]
			ng := rand.IntN(GroupCount)
			if ng != oldG {
				ind.Gene[n] = ng
			}
		}
	}

	fixGroups(ind)
	return *ind
}
func runGA(history []Draw) Individual {

	islands := 8       // 6
	popSize := 800     // 每岛大小 600
	totalGen := 6      // 6000
	migrationGap := 60 // 每隔 120 代迁移
	mutationRate := 0.08
	maxStagnation := 12

	// 初始化岛屿
	pops := make([][]Individual, islands)
	bests := make([]Individual, islands)

	for i := 0; i < islands; i++ {
		pops[i] = make([]Individual, popSize)
		for j := 0; j < popSize; j++ {
			pops[i][j] = newValidIndividual()
			fitness(&pops[i][j], history)
		}
		bests[i] = pops[i][0]
	}

	// 全局最优
	globalBest := bests[0]
	stagnation := 0

	for gen := 1; gen <= totalGen; gen++ {

		for isl := 0; isl < islands; isl++ {

			// 演化一代
			newPop := make([]Individual, popSize)

			// 精英保留
			for _, ind := range pops[isl] {
				if ind.Score > bests[isl].Score {
					bests[isl] = ind
				}
			}
			newPop[0] = bests[isl]

			// 产生新个体
			for k := 1; k < popSize; k++ {

				p1 := pops[isl][rand.IntN(popSize)]
				p2 := pops[isl][rand.IntN(popSize)]

				c := crossover(p1, p2)
				c = mutate(&c, mutationRate)

				fitness(&c, history)
				localImprove(&c, history)

				newPop[k] = c
			}

			pops[isl] = newPop
		}

		// 更新全局最优
		improved := false
		for isl := 0; isl < islands; isl++ {
			if bests[isl].Score > globalBest.Score {
				globalBest = bests[isl]
				improved = true
				stagnation = 0
				fmt.Printf("🔥 第 %d 代出现更优：%d\n", gen, globalBest.Score)
			}
		}

		if !improved {
			stagnation++
		}

		// 停滞 → 自适应处理
		if stagnation >= maxStagnation {
			fmt.Println("⚠️ 停滞，触发 Mutation Explosion")

			mutationRate *= 1.25
			if mutationRate > 0.35 {
				mutationRate = 0.35
			}

			// 注入新个体
			for isl := 0; isl < islands; isl++ {
				for j := 1; j < popSize; j++ {
					if rand.Float64() < 0.10 {
						pops[isl][j] = newValidIndividual()
						fitness(&pops[isl][j], history)
					}
				}
			}

			stagnation = 0
		}

		// 岛屿迁移
		if gen%migrationGap == 0 {
			fmt.Printf("🏝 岛屿迁移（第 %d 代）\n", gen)
			for isl := 0; isl < islands; isl++ {
				next := (isl + 1) % islands
				// 把最佳传给下一岛
				pops[next][rand.IntN(popSize)] = bests[isl]
			}
		}
	}

	return globalBest
}

// 把 bestGroup 映射成 A,B,C,D,E 五组号码
func formatGroup(g [35]int) {
	groups := make([][]int, 5)
	for i := 0; i < 5; i++ {
		groups[i] = []int{}
	}

	for i := 0; i < 35; i++ {
		groupID := g[i]
		groups[groupID] = append(groups[groupID], i+1)
	}

	labels := []string{"A", "B", "C", "D", "E"}

	for i := 0; i < 5; i++ {
		fmt.Printf("%s: [", labels[i])
		for j, num := range groups[i] {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Printf(`"%02d"`, num)
		}
		fmt.Println("]")
	}
}

func Moni() {
	history := loadHistory()

	fmt.Println("开始遗传算法搜索…")

	best := runGA(history)

	fmt.Println("\n====== GA 最终最优结果 ======")
	fmt.Println("命中次数：", best.Score)
	fmt.Println("最终最佳组合（组别数组）：", best.Gene)
	fmt.Println("格式化分组如下：")
	formatGroup(best.Gene)
}


```
遗传学版本：加强版多线程

```go
package ana

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"sync"

	"github.com/before80/lot/dbop"
)

type Draw struct {
	Front []int
}

const (
	GroupCount = 5
	TotalNums  = 35
)

// ========== 数据加载 ==========
func loadHistory() []Draw {
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

// ========== 个体定义 ==========
type Individual struct {
	Gene  [TotalNums]int // 每个索引 i 表示号码 i+1 属于哪个组 0..4
	Score int
}

// ========== 工具：校验 sizes ==========
func validateSizes(sizes []int) {
	if len(sizes) != GroupCount {
		panic("sizes length must be 5")
	}
	sum := 0
	for _, v := range sizes {
		sum += v
	}
	if sum != TotalNums {
		panic(fmt.Sprintf("sizes sum must be %d, got %d", TotalNums, sum))
	}
}

// ========== 初始化合法个体（按 sizes 分配） ==========
func newValidIndividual(sizes []int) Individual {
	validateSizes(sizes)

	var ind Individual
	nums := rand.Perm(TotalNums) // 0..34 随机排列
	pos := 0
	for g := 0; g < GroupCount; g++ {
		for i := 0; i < sizes[g]; i++ {
			ind.Gene[nums[pos]] = g
			pos++
		}
	}
	return ind
}

// ========== 适应度计算（不依赖 sizes） ==========
func fitness(ind *Individual, history []Draw) int {
	count := 0
	for _, d := range history {
		groupUsed := [GroupCount]bool{}
		ok := true
		for _, num := range d.Front {
			g := ind.Gene[num-1]
			if groupUsed[g] {
				ok = false
				break
			}
			groupUsed[g] = true
		}
		if ok {
			// 所有5个号码来自不同组 -> 命中
			count++
		}
	}
	ind.Score = count
	return count
}

// ========== 局部改进（Lamarckian） ==========
func localImprove(ind *Individual, history []Draw) {
	bestScore := ind.Score
	// 做少量随机交换尝试
	for i := 0; i < 3; i++ {
		a := rand.IntN(TotalNums)
		b := rand.IntN(TotalNums)
		ind.Gene[a], ind.Gene[b] = ind.Gene[b], ind.Gene[a]

		newScore := fitness(ind, history)
		if newScore <= bestScore {
			// 回滚
			ind.Gene[a], ind.Gene[b] = ind.Gene[b], ind.Gene[a]
			ind.Score = bestScore
		} else {
			bestScore = newScore
		}
	}
}

// ========== 修复组大小（依据 sizes） ==========
func fixGroups(ind *Individual, sizes []int) {
	// 统计当前各组人数
	groupCount := [GroupCount]int{}
	for n := 0; n < TotalNums; n++ {
		groupCount[ind.Gene[n]]++
	}

	// 计算缺和多
	var need []int  // 需要哪些组（每个元素是组id，需要一个额外名额）
	var extra []int // 哪些组多了（每个元素表示一个要移走的位置对应的组id）

	for g := 0; g < GroupCount; g++ {
		if groupCount[g] < sizes[g] {
			for i := 0; i < sizes[g]-groupCount[g]; i++ {
				need = append(need, g)
			}
		} else if groupCount[g] > sizes[g] {
			for i := 0; i < groupCount[g]-sizes[g]; i++ {
				extra = append(extra, g)
			}
		}
	}

	// 把出现 extra 的位置替换为 need
	ei := 0
	ni := 0
	for n := 0; n < TotalNums && ei < len(extra) && ni < len(need); n++ {
		if ind.Gene[n] == extra[ei] {
			ind.Gene[n] = need[ni]
			ei++
			ni++
		}
	}
}

// ========== 交叉（使用 sizes） ==========
func crossover(p1, p2 Individual, sizes []int) Individual {
	var child Individual
	// 尝试把每个组从 p1 取 sizes[g]/2 个，剩余从 p2 填
	count := [GroupCount]int{}
	// 第一步：先试从 p1 拿
	for n := 0; n < TotalNums; n++ {
		g1 := p1.Gene[n]
		if count[g1] < sizes[g1]/2 {
			child.Gene[n] = g1
			count[g1]++
		} else {
			// 暂时标记为 -1，用 p2 填充后修正
			child.Gene[n] = -1
		}
	}
	// 第二步：用 p2 填充剩余位置
	for n := 0; n < TotalNums; n++ {
		if child.Gene[n] == -1 {
			child.Gene[n] = p2.Gene[n]
			count[child.Gene[n]]++
		}
	}
	// 修复以确保严格满足 sizes
	fixGroups(&child, sizes)
	return child
}

// ========== 变异（按概率尝试单个位置改变组），然后修复 ==========
func mutate(ind *Individual, rate float64, sizes []int) Individual {
	for n := 0; n < TotalNums; n++ {
		if rand.Float64() < rate {
			oldG := ind.Gene[n]
			ng := rand.IntN(GroupCount)
			if ng != oldG {
				ind.Gene[n] = ng
			}
		}
	}
	// 修复组大小
	fixGroups(ind, sizes)
	return *ind
}

// ========== 并行评估（线程池） ==========
func evaluatePopulationParallel(pop []Individual, history []Draw, workers int) {
	n := len(pop)
	jobs := make(chan int, n)
	wg := sync.WaitGroup{}
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				fitness(&pop[idx], history)
			}
		}()
	}
	for i := 0; i < n; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
}

// ========== 多岛并行 GA（sizes 作为参数传入） ==========
func runGA(history []Draw, sizes []int) Individual {
	// 参数（你可以按需调整）
	islands := 80
	popSize := 8866
	generations := 1000
	migrationGap := 6
	mutationRate := 0.05
	maxStagnation := 12
	workers := 16

	// 校验 sizes
	validateSizes(sizes)

	// 初始化岛屿
	pops := make([][]Individual, islands)
	bests := make([]Individual, islands)
	for i := 0; i < islands; i++ {
		pops[i] = make([]Individual, popSize)
		for j := 0; j < popSize; j++ {
			pops[i][j] = newValidIndividual(sizes)
		}
		// 并行评估初始适应度
		evaluatePopulationParallel(pops[i], history, workers)
		bests[i] = pops[i][0]
	}

	globalBest := bests[0]
	stagnation := 0
	var globalMu sync.Mutex

	for gen := 1; gen <= generations; gen++ {
		// ---- 每个岛并行演化 ----
		wg := sync.WaitGroup{}
		wg.Add(islands)

		for isl := 0; isl < islands; isl++ {
			go func(island int) {
				defer wg.Done()
				pop := pops[island]
				newPop := make([]Individual, popSize)

				// 精英保留：找到本岛最优
				best := bests[island]
				for _, ind := range pop {
					if ind.Score > best.Score {
						best = ind
					}
				}
				newPop[0] = best
				bests[island] = best

				// 并行生成新个体（内部使用 worker 池）
				jobs := make(chan int, popSize)
				subwg := sync.WaitGroup{}
				subwg.Add(workers)

				for w := 0; w < workers; w++ {
					go func() {
						defer subwg.Done()
						for k := range jobs {
							if k == 0 {
								continue
							}
							// 父代随机选择（可改为 tournament 以更好）
							p1 := pop[rand.IntN(popSize)]
							p2 := pop[rand.IntN(popSize)]

							c := crossover(p1, p2, sizes)
							c = mutate(&c, mutationRate, sizes)
							fitness(&c, history)
							localImprove(&c, history)

							newPop[k] = c
						}
					}()
				}

				for k := 1; k < popSize; k++ {
					jobs <- k
				}
				close(jobs)
				subwg.Wait()

				pops[island] = newPop
			}(isl)
		}

		wg.Wait()

		// ---- 更新全局最优 ----
		improved := false
		globalMu.Lock()
		for isl := 0; isl < islands; isl++ {
			if bests[isl].Score > globalBest.Score {
				globalBest = bests[isl]
				improved = true
				stagnation = 0
				fmt.Printf("🔥 第 %d 代出现更优：%d\n", gen, globalBest.Score)
			}
		}
		globalMu.Unlock()

		if !improved {
			stagnation++
		}

		// 停滞处理：自适应变异并注入新个体
		if stagnation >= maxStagnation {
			fmt.Println("⚠️ 停滞 -> Mutation Explosion")
			mutationRate *= 1.25
			if mutationRate > 0.35 {
				mutationRate = 0.35
			}
			for isl := 0; isl < islands; isl++ {
				for j := 1; j < popSize; j++ {
					if rand.Float64() < 0.10 {
						pops[isl][j] = newValidIndividual(sizes)
						fitness(&pops[isl][j], history)
					}
				}
			}
			stagnation = 0
		}

		// 岛屿迁移
		if gen%migrationGap == 0 {
			fmt.Printf("🏝 岛屿迁移（第 %d 代）\n", gen)
			for isl := 0; isl < islands; isl++ {
				next := (isl + 1) % islands
				pops[next][rand.IntN(popSize)] = bests[isl]
			}
		}
	}

	return globalBest
}

// ========== 映射并格式化输出 ==========
func formatGroup(g [TotalNums]int) {
	groups := make([][]int, GroupCount)
	for i := 0; i < GroupCount; i++ {
		groups[i] = []int{}
	}
	for i := 0; i < TotalNums; i++ {
		groupID := g[i]
		groups[groupID] = append(groups[groupID], i+1)
	}
	labels := []string{"A", "B", "C", "D", "E"}
	for i := 0; i < GroupCount; i++ {
		fmt.Printf("%s: [", labels[i])
		for j, num := range groups[i] {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Printf(`"%02d"`, num)
		}
		fmt.Println("]")
	}
}

// Moni 接受 sizes 参数（线程安全） ==========
func Moni(sizes []int) {
	// sizes 必须合法
	validateSizes(sizes)

	history := loadHistory()

	fmt.Println("开始遗传算法搜索… sizes =", sizes)

	best := runGA(history, sizes)

	fmt.Println("\n====== GA 最终最优结果 ======")
	fmt.Println("命中次数：", best.Score)
	fmt.Println("最终最佳组合（组别数组）：", best.Gene)
	fmt.Println("格式化分组如下：")
	formatGroup(best.Gene)
}


```