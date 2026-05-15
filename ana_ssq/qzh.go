package ana_ssq

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/before80/lot/dbop"
)

// CalSsqQzh 计算大乐透全区5个号码中，前1到17，中18，后19到35，各自出现的个数
//
//	@Description:
//	@param lotteryDrawResults
//	@return string
func CalSsqQzh(hmStrSlice []string) string {
	var qNum, zNum, HNum int
	for i, hmStr := range hmStrSlice {
		if i < 5 {
			num, _ := strconv.Atoi(hmStr)
			if num < 17 {
				qNum++
			}
			if num == 17 {
				zNum++
			}
			if num > 17 {
				HNum++
			}
		}
	}
	return fmt.Sprintf("%d%d%d", qNum, zNum, HNum)
}

// SsqQzhQuShi
//
//	@Description:
//	@return []KeyWithLength
func SsqQzhQuShi() []KeyWithLength {
	ssqs, _ := dbop.ReadAllSsq(false)
	qzhMs := make(map[string]int)
	for i, ssq := range ssqs {
		if i == len(ssqs)-1 {
			break
		}
		nextSsq := ssqs[i+1]

		aeArrow := fmt.Sprintf("%s->%s", ssq.Qzh, nextSsq.Qzh)
		if _, ok := qzhMs[aeArrow]; !ok {
			qzhMs[aeArrow] = 1
		} else {
			qzhMs[aeArrow]++
		}
	}

	// 对oeMs按照存放的值的大小进行排序
	kLens := make([]KeyWithLength, 0, len(qzhMs))

	for aeArrow, v := range qzhMs {
		kLens = append(kLens, KeyWithLength{
			Key:    aeArrow,
			Length: v,
		})
	}
	sort.Slice(kLens, func(i, j int) bool {
		return kLens[i].Length > kLens[j].Length
	})
	return kLens
}

// SsqQzhHis
//
//	@Description:
//	@param wg
//	@return res
func SsqQzhHis(wg *sync.WaitGroup) (res []SsqHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2SsqHis := make(map[string]*SsqHis)

	// 初始化
	for k, v := range AllSsqQzh2Count {
		typ2SsqHis[k] = &SsqHis{Typ: k, AllCount: v}
	}

	lenSsqHis := len(ZxSsqs)

	for i, dlt := range ZxSsqs {
		qzhStr := dlt.Qzh

		typ2SsqHis[qzhStr].Cs = typ2SsqHis[qzhStr].Cs + 1
		if lenSsqHis-i <= 10 {
			typ2SsqHis[qzhStr].Last10 = typ2SsqHis[qzhStr].Last10 + 1
		}
		if lenSsqHis-i <= 20 {
			typ2SsqHis[qzhStr].Last20 = typ2SsqHis[qzhStr].Last20 + 1
		}
		if lenSsqHis-i <= 30 {
			typ2SsqHis[qzhStr].Last30 = typ2SsqHis[qzhStr].Last30 + 1
		}
		if lenSsqHis-i <= 50 {
			typ2SsqHis[qzhStr].Last50 = typ2SsqHis[qzhStr].Last50 + 1
		}
		if lenSsqHis-i <= 100 {
			typ2SsqHis[qzhStr].Last100 = typ2SsqHis[qzhStr].Last100 + 1
		}
		if lenSsqHis-i <= 200 {
			typ2SsqHis[qzhStr].Last200 = typ2SsqHis[qzhStr].Last200 + 1
		}
		if lenSsqHis-i <= 500 {
			typ2SsqHis[qzhStr].Last500 = typ2SsqHis[qzhStr].Last500 + 1
		}
		if lenSsqHis-i <= 1000 {
			typ2SsqHis[qzhStr].Last1000 = typ2SsqHis[qzhStr].Last1000 + 1
		}
		if lenSsqHis-i <= 1500 {
			typ2SsqHis[qzhStr].Last1500 = typ2SsqHis[qzhStr].Last1500 + 1
		}
		if lenSsqHis-i <= 2000 {
			typ2SsqHis[qzhStr].Last2000 = typ2SsqHis[qzhStr].Last2000 + 1
		}
		if lenSsqHis-i <= 2500 {
			typ2SsqHis[qzhStr].Last2500 = typ2SsqHis[qzhStr].Last2500 + 1
		}
		if lenSsqHis-i <= 3500 {
			typ2SsqHis[qzhStr].Last3500 = typ2SsqHis[qzhStr].Last3500 + 1
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
