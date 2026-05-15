package ana_ssq

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
func BuildFrontHmSlices(src []models.SsqMoni, size int) [][]string {
	out := make([][]string, size)
	for i, m := range src {
		if i >= size {
			break
		}
		out[i] = gen.SsqFrontHmStrSliceFromAeStr(m.A, m.B, m.C, m.D, m.E, m.F)
	}
	return out
}

type SsqExcelData struct {
	XuHao               int
	DrawNum             string
	DrawTime            string
	FrontHm             string
	FullHm              string
	Oe                  string
	Hz                  int
	Qzh                 string
	NewAddCh4           int // 当期新增4重号数
	NewAddCh5           int // 当期新增5重号数
	NewAddCh6           int // 当期新增6重号数
	NewAddCh7           int // 当期新增7重号数
	DangQiTotalNewAddCh int // 当期总的新增重号数
	LeiJiaCh            int // 累加重号数
	T61                 int
	T62                 int
	T63                 int
	T64                 int
	T65                 int
	T81                 int
	T82                 int
	T83                 int
	T84                 int
	T85                 int
	T131                int
	T132                int
	T133                int
	T134                int
	T135                int
	T181                int
	T182                int
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
