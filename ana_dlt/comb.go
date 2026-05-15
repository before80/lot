package ana_dlt

import (
	"github.com/before80/lot/gen"
	"github.com/before80/lot/models"
)

// BuildFrontHmSlices 构建前区号码
//
//	@Description:
//	@param src
//	@param size
//	@return [][]string
func BuildFrontHmSlices(src []models.DltMoni, size int) [][]string {
	out := make([][]string, size)
	for i, m := range src {
		if i >= size {
			break
		}
		out[i] = gen.DltFrontHmStrSliceFromAeStr(m.A, m.B, m.C, m.D, m.E)
	}
	return out
}

type DltExcelData struct {
	XuHao               int
	DrawNum             string
	DrawTime            string
	EquipmentCount      int
	FrontHm             string
	FullHm              string
	UnSortDrawResult    string
	PoolBalance         float64 // 奖池奖金(元)
	TotalSaleAmount     float64 // 当期总销售额(元)
	StakeCount101       int     // 一等奖注数
	StakeAmount101      int     // 一等奖奖金
	StakeCount201       int     // 一等奖追加注数
	StakeAmount201      int     // 一等奖追加奖金
	StakeCount301       int     // 二等奖注数
	StakeAmount301      int     // 二等奖奖金
	StakeCount401       int     // 二等奖追加注数
	StakeAmount401      int     // 二等奖追加奖金
	Oe                  string
	Hz                  int
	Qzh                 string
	NewAddCh4           int // 当期新增4重号数
	NewAddCh5           int // 当期新增5重号数
	NewAddCh6           int // 当期新增6重号数
	NewAddCh7           int // 当期新增7重号数
	DangQiTotalNewAddCh int // 当期总的新增重号数
	LeiJiaCh            int // 累加重号数
	T71                 int
	T72                 int
	T73                 int
	T74                 int
	T75                 int
	T111                int
	T112                int
	T113                int
	T114                int
	T115                int
	T151                int
	T152                int
	T153                int
	T154                int
	T155                int
	T191                int
	T192                int
	OtherT              int
}

// OR 两个 block-slices（假设长度相同）
func orBlocks(a, b []uint64) {
	for i := 0; i < len(a); i++ {
		a[i] |= b[i]
	}
}

// 统计 blocks 中的 1 数量
func countBitsBlocks(blocks []uint64) int {
	cnt := 0
	for _, x := range blocks {
		for x > 0 {
			x &= x - 1
			cnt++
		}
	}
	return cnt
}

// 从 n 个元素中选 k 个，返回索引组合列表
func combinations(n, k int) [][]int {
	var res [][]int
	comb := make([]int, k)

	var dfs func(start, depth int)
	dfs = func(start, depth int) {
		if depth == k {
			tmp := make([]int, k)
			copy(tmp, comb)
			res = append(res, tmp)
			return
		}
		for i := start; i <= n-(k-depth); i++ {
			comb[depth] = i
			dfs(i+1, depth+1)
		}
	}
	if k > 0 && k <= n {
		dfs(0, 0)
	}
	return res
}

// 返回多个最佳组合
func maxCoverageAll(colsMasks [][]uint64, colNames []string, K int) ([][]string, int) {
	N := len(colsMasks)
	if K <= 0 || K > N {
		return nil, 0
	}

	combs := combinations(N, K)
	blockCount := len(colsMasks[0])

	bestCover := -1
	var bestCombos [][]string

	for _, comb := range combs {
		acc := make([]uint64, blockCount)
		for _, idx := range comb {
			orBlocks(acc, colsMasks[idx])
		}
		covered := countBitsBlocks(acc)

		if covered > bestCover {
			bestCover = covered
			bestCombos = bestCombos[:0] // 清空旧的
			tmp := make([]string, 0, K)
			for _, idx := range comb {
				tmp = append(tmp, colNames[idx])
			}
			bestCombos = append(bestCombos, tmp)
		} else if covered == bestCover {
			tmp := make([]string, 0, K)
			for _, idx := range comb {
				tmp = append(tmp, colNames[idx])
			}
			bestCombos = append(bestCombos, tmp)
		}
	}

	return bestCombos, bestCover
}
