package ana_ssq

import (
	"sort"
	"sync"
)

// SsqFrontOnlyOneHis 双色球前区单个号码的历史
//
//	@Description:
//	@param wg
//	@return res
func SsqFrontOnlyOneHis(wg *sync.WaitGroup) (res []SsqHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2SsqHis := make(map[string]*SsqHis)

	// 初始化
	for k, v := range AllSsqOnlyOneFront2Count {
		typ2SsqHis[k] = &SsqHis{Typ: k, AllCount: v}
	}

	lenSsqHis := len(ZxSsqs)

	for i, ssq := range ZxSsqs {
		frontHms := []string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}

		for _, fHm := range frontHms {
			typ2SsqHis[fHm].Cs = typ2SsqHis[fHm].Cs + 1
			if lenSsqHis-i <= 10 {
				typ2SsqHis[fHm].Last10 = typ2SsqHis[fHm].Last10 + 1
			}
			if lenSsqHis-i <= 20 {
				typ2SsqHis[fHm].Last20 = typ2SsqHis[fHm].Last20 + 1
			}
			if lenSsqHis-i <= 30 {
				typ2SsqHis[fHm].Last30 = typ2SsqHis[fHm].Last30 + 1
			}
			if lenSsqHis-i <= 50 {
				typ2SsqHis[fHm].Last50 = typ2SsqHis[fHm].Last50 + 1
			}
			if lenSsqHis-i <= 100 {
				typ2SsqHis[fHm].Last100 = typ2SsqHis[fHm].Last100 + 1
			}
			if lenSsqHis-i <= 200 {
				typ2SsqHis[fHm].Last200 = typ2SsqHis[fHm].Last200 + 1
			}
			if lenSsqHis-i <= 500 {
				typ2SsqHis[fHm].Last500 = typ2SsqHis[fHm].Last500 + 1
			}
			if lenSsqHis-i <= 1000 {
				typ2SsqHis[fHm].Last1000 = typ2SsqHis[fHm].Last1000 + 1
			}
			if lenSsqHis-i <= 1500 {
				typ2SsqHis[fHm].Last1500 = typ2SsqHis[fHm].Last1500 + 1
			}
			if lenSsqHis-i <= 2000 {
				typ2SsqHis[fHm].Last2000 = typ2SsqHis[fHm].Last2000 + 1
			}
			if lenSsqHis-i <= 2500 {
				typ2SsqHis[fHm].Last2500 = typ2SsqHis[fHm].Last2500 + 1
			}
			if lenSsqHis-i <= 3500 {
				typ2SsqHis[fHm].Last3500 = typ2SsqHis[fHm].Last3500 + 1
			}
		}

	}

	kLens := make([]KeyWithLength, 0, len(typ2SsqHis))

	for k, ssqHis := range typ2SsqHis {
		kLens = append(kLens, KeyWithLength{
			Key:    k,
			Length: ssqHis.Cs,
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
