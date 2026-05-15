package ana_ssq

import (
	"fmt"
	"sort"
	"sync"

	"github.com/before80/lot/dbop"
	"github.com/before80/lot/gen"
)

// SsqBackHis 双色球后区历史数据(按出现期数的大小,从大到小排序)
//
//	@Description:
//	@return res
func SsqBackHis(wg *sync.WaitGroup) (res []SsqHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	typ2SsqHis := make(map[string]*SsqHis)
	backCombs := gen.AllSsqBackHms
	// 初始化 typ2SsqHis
	for _, backComb := range backCombs {
		typ2SsqHis[backComb] = &SsqHis{Typ: backComb, AllCount: 1107568}
	}

	ssqs, _ := dbop.ReadAllSsq(false)

	lenSsqHis := len(ssqs)
	for i, ssq := range ssqs {
		curBackComb := ssq.B1
		typ2SsqHis[curBackComb].Cs = typ2SsqHis[curBackComb].Cs + 1
		if lenSsqHis-i <= 10 {
			typ2SsqHis[curBackComb].Last10 = typ2SsqHis[curBackComb].Last10 + 1
		}
		if lenSsqHis-i <= 20 {
			typ2SsqHis[curBackComb].Last20 = typ2SsqHis[curBackComb].Last20 + 1
		}
		if lenSsqHis-i <= 30 {
			typ2SsqHis[curBackComb].Last30 = typ2SsqHis[curBackComb].Last30 + 1
		}
		if lenSsqHis-i <= 50 {
			typ2SsqHis[curBackComb].Last50 = typ2SsqHis[curBackComb].Last50 + 1
		}
		if lenSsqHis-i <= 100 {
			typ2SsqHis[curBackComb].Last100 = typ2SsqHis[curBackComb].Last100 + 1
		}
		if lenSsqHis-i <= 200 {
			typ2SsqHis[curBackComb].Last200 = typ2SsqHis[curBackComb].Last200 + 1
		}
		if lenSsqHis-i <= 500 {
			typ2SsqHis[curBackComb].Last500 = typ2SsqHis[curBackComb].Last500 + 1
		}
		if lenSsqHis-i <= 1000 {
			typ2SsqHis[curBackComb].Last1000 = typ2SsqHis[curBackComb].Last1000 + 1
		}
		if lenSsqHis-i <= 1500 {
			typ2SsqHis[curBackComb].Last1500 = typ2SsqHis[curBackComb].Last1500 + 1
		}
		if lenSsqHis-i <= 2000 {
			typ2SsqHis[curBackComb].Last2000 = typ2SsqHis[curBackComb].Last2000 + 1
		}
		if lenSsqHis-i <= 2500 {
			typ2SsqHis[curBackComb].Last2500 = typ2SsqHis[curBackComb].Last2500 + 1
		}
		if lenSsqHis-i <= 3500 {
			typ2SsqHis[curBackComb].Last3500 = typ2SsqHis[curBackComb].Last3500 + 1
		}
	}

	kLens := make([]KeyWithLength, 0, len(typ2SsqHis))
	for k, ssqHis := range typ2SsqHis {
		kLens = append(kLens, KeyWithLength{Key: k, Length: ssqHis.Cs})
	}
	// 对typ2SsqHis按照存放的Cs值的大小进行排序
	sort.Slice(kLens, func(i, j int) bool {
		return kLens[i].Length > kLens[j].Length
	})
	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2SsqHis[typ])
	}
	return
}

// SsqBackQuShi1 双色球后区趋势1(后区号码->后区号码)
//
//	@Description:
//	@return res
func SsqBackQuShi1() (res []KeyWithLength) {
	// 生成后区所有组合字符串的切片
	backCombs := gen.AllSsqBackHms
	quShi2Count := make(map[string]int)
	// 初始化
	for _, backComb1 := range backCombs {
		for _, backComb2 := range backCombs {
			quShi2Count[fmt.Sprintf("%s->%s", backComb1, backComb2)] = 0
		}
	}
	for i, dlt := range ZxSsqs {
		if i == len(ZxSsqs)-1 {
			break
		}
		nextSsq := ZxSsqs[i+1]
		quShiStr := fmt.Sprintf("%s->%s", dlt.B1, nextSsq.B1)
		quShi2Count[quShiStr] = quShi2Count[quShiStr] + 1
	}

	kLens := make([]KeyWithLength, 0, len(quShi2Count))
	for k, v := range quShi2Count {
		kLens = append(kLens, KeyWithLength{k, v})
	}

	sort.Slice(kLens, func(i, j int) bool {
		if kLens[i].Length == kLens[j].Length {
			return kLens[i].Key > kLens[j].Key
		}
		return kLens[i].Length > kLens[j].Length
	})

	return kLens
}

type SsqBackChQuShi struct {
	BackComb         string
	Qs               string
	Cs               int
	AllCombs         []string // 所有可能的组合
	HadExistCombs    []string // 已经出现的组合
	HadNotExistCombs []string // 未出现的组合
}

// SsqBackOnlyOneHis 双色球后区单个号码的历史
//
//	@Description:
//	@param wg
//	@return res
func SsqBackOnlyOneHis(wg *sync.WaitGroup) (res []SsqHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2SsqHis := make(map[string]*SsqHis)

	// 初始化
	for k, v := range AllSsqOnlyOneBack2Count {
		typ2SsqHis[k] = &SsqHis{Typ: k, AllCount: v}
	}
	ssqs, _ := dbop.ReadAllSsq(false)
	lenSsqHis := len(ssqs)

	for i, dlt := range ssqs {
		backHms := []string{dlt.B1}

		for _, bHm := range backHms {
			typ2SsqHis[bHm].Cs = typ2SsqHis[bHm].Cs + 1
			if lenSsqHis-i <= 10 {
				typ2SsqHis[bHm].Last10 = typ2SsqHis[bHm].Last10 + 1
			}
			if lenSsqHis-i <= 20 {
				typ2SsqHis[bHm].Last20 = typ2SsqHis[bHm].Last20 + 1
			}
			if lenSsqHis-i <= 30 {
				typ2SsqHis[bHm].Last30 = typ2SsqHis[bHm].Last30 + 1
			}
			if lenSsqHis-i <= 50 {
				typ2SsqHis[bHm].Last50 = typ2SsqHis[bHm].Last50 + 1
			}
			if lenSsqHis-i <= 100 {
				typ2SsqHis[bHm].Last100 = typ2SsqHis[bHm].Last100 + 1
			}
			if lenSsqHis-i <= 200 {
				typ2SsqHis[bHm].Last200 = typ2SsqHis[bHm].Last200 + 1
			}
			if lenSsqHis-i <= 500 {
				typ2SsqHis[bHm].Last500 = typ2SsqHis[bHm].Last500 + 1
			}
			if lenSsqHis-i <= 1000 {
				typ2SsqHis[bHm].Last1000 = typ2SsqHis[bHm].Last1000 + 1
			}
			if lenSsqHis-i <= 1500 {
				typ2SsqHis[bHm].Last1500 = typ2SsqHis[bHm].Last1500 + 1
			}
			if lenSsqHis-i <= 2000 {
				typ2SsqHis[bHm].Last2000 = typ2SsqHis[bHm].Last2000 + 1
			}
			if lenSsqHis-i <= 2500 {
				typ2SsqHis[bHm].Last2500 = typ2SsqHis[bHm].Last2500 + 1
			}
			if lenSsqHis-i <= 3500 {
				typ2SsqHis[bHm].Last3500 = typ2SsqHis[bHm].Last3500 + 1
			}
		}

	}

	kLens := make([]KeyWithLength, 0, len(typ2SsqHis))

	for k, dltHis := range typ2SsqHis {
		kLens = append(kLens, KeyWithLength{
			Key:    k,
			Length: dltHis.Cs,
		})
	}
	// 对typ2SsqHis按照存放的Cs值的大小进行排序
	sort.Slice(kLens, func(i, j int) bool {
		if kLens[i].Length == kLens[j].Length {
			return kLens[i].Key > kLens[j].Key
		}
		return kLens[i].Length > kLens[j].Length
	})

	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2SsqHis[typ])
	}

	return
}
