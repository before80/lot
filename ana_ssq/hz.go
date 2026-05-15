package ana_ssq

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// CalSsqHz 计算双色球前后区号码的和值
//
//	@Description:
//	@param lotteryDrawResults
//	@return total
func CalSsqHz(hmStrSlice []string) (total int) {
	for _, hmStr := range hmStrSlice {
		num, _ := strconv.Atoi(hmStr)
		total += num
	}
	return total
}

// CalSsqHzFromStr 计算双色球前后区号码的和值
//
//	@Description:
//	@param fullHmStr
//	@return total
func CalSsqHzFromStr(fullHmStr string) (total int) {
	fullHmStr = strings.Replace(fullHmStr, "|", ",", -1)
	fullHms := strings.Split(fullHmStr, ",")
	for _, fullHm := range fullHms {
		num, _ := strconv.Atoi(fullHm)
		total += num
	}
	return total
}

type MaxMin struct {
	Min int
	Max int
}

var HzABCDEMs = map[string]MaxMin{
	"A": MaxMin{Min: 22, Max: 57},
	"B": MaxMin{Min: 58, Max: 93},
	"C": MaxMin{Min: 94, Max: 129},
	"D": MaxMin{Min: 130, Max: 165},
	"E": MaxMin{Min: 166, Max: 199},
}

func SsqHzABCDE(hz int) string {
	if hz <= 57 {
		return "A"
	}
	if hz <= 93 {
		return "B"
	}
	if hz <= 129 {
		return "C"
	}
	if hz <= 165 {
		return "D"
	}
	if hz <= 199 {
		return "E"
	}
	return "E"
}

// SsqHzQuShi1 双色球和值趋势(和值->和值) 注意需要提前之前 InitSsqs()
//
//	@Description:
//	@return []KeyWithLength
func SsqHzQuShi1() []KeyWithLength {
	aeMs := make(map[string]int)
	for i, dlt := range ZxSsqs {
		if i == len(ZxSsqs)-1 {
			break
		}
		nextSsq := ZxSsqs[i+1]
		aeArrow := fmt.Sprintf("%d->%d", dlt.Hz, nextSsq.Hz)
		if _, ok := aeMs[aeArrow]; !ok {
			aeMs[aeArrow] = 1
		} else {
			aeMs[aeArrow]++
		}
	}

	kLens := make([]KeyWithLength, 0, len(aeMs))
	for aeArrow, v := range aeMs {
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

type SsqHzChQuShi struct {
	HzCh      string
	NextHzCh  string
	Cs        int
	Hz2NextHz []string
}

// SsqHzQuShi2 双色球和值趋势(和值Ch->和值Ch) 注意需要提前之前 InitSsqs()
//
//	@Description:
//	@return res
func SsqHzQuShi2() (res []SsqHzChQuShi) {
	hz2ChQuShiSt := make(map[string]*SsqHzChQuShi)
	for i, dlt := range ZxSsqs {
		if i == len(ZxSsqs)-1 {
			break
		}
		nextSsq := ZxSsqs[i+1]
		nextHzABCDE := SsqHzABCDE(nextSsq.Hz)
		hzABCDE := SsqHzABCDE(dlt.Hz)
		aeArrow := fmt.Sprintf("%s->%s", hzABCDE, nextHzABCDE)
		hz2NextHzStr := fmt.Sprintf("%d->%d", dlt.Hz, nextSsq.Hz)
		if _, ok := hz2ChQuShiSt[aeArrow]; !ok {
			hz2ChQuShiSt[aeArrow] = &SsqHzChQuShi{
				HzCh:      hzABCDE,
				NextHzCh:  nextHzABCDE,
				Cs:        1,
				Hz2NextHz: []string{hz2NextHzStr},
			}
		} else {
			hz2ChQuShiSt[aeArrow].Cs = hz2ChQuShiSt[aeArrow].Cs + 1
			if !slices.Contains(hz2ChQuShiSt[aeArrow].Hz2NextHz, hz2NextHzStr) {
				hz2ChQuShiSt[aeArrow].Hz2NextHz = append(hz2ChQuShiSt[aeArrow].Hz2NextHz, hz2NextHzStr)
			}
		}
	}

	// 对aeMs按照存放的值的大小进行排序
	kLens := make([]KeyWithLength, 0, len(hz2ChQuShiSt))
	for aeArrow, v := range hz2ChQuShiSt {
		kLens = append(kLens, KeyWithLength{
			Key:    aeArrow,
			Length: v.Cs,
		})
	}
	sort.Slice(kLens, func(i, j int) bool {
		if kLens[i].Length == kLens[j].Length {
			return kLens[i].Key > kLens[j].Key
		}
		return kLens[i].Length > kLens[j].Length
	})

	for _, v := range kLens {
		res = append(res, *hz2ChQuShiSt[v.Key])
	}
	return
}

// SsqHzHis 双色球和值历史数据(按出现期数的大小,从大到小排序)
//
//	@Description:
//	@return res
func SsqHzHis(wg *sync.WaitGroup) (res []SsqHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2SsqHis := make(map[string]*SsqHis)
	// 初始化
	for k, v := range AllSsqHz2Count {
		typ2SsqHis[k] = &SsqHis{Typ: k, AllCount: v}
	}
	lenSsqHis := len(ZxSsqs)

	for i, dlt := range ZxSsqs {
		hzStr := strconv.Itoa(dlt.Hz)
		typ2SsqHis[hzStr].Cs = typ2SsqHis[hzStr].Cs + 1
		if lenSsqHis-i <= 10 {
			typ2SsqHis[hzStr].Last10 = typ2SsqHis[hzStr].Last10 + 1
		}
		if lenSsqHis-i <= 20 {
			typ2SsqHis[hzStr].Last20 = typ2SsqHis[hzStr].Last20 + 1
		}
		if lenSsqHis-i <= 30 {
			typ2SsqHis[hzStr].Last30 = typ2SsqHis[hzStr].Last30 + 1
		}
		if lenSsqHis-i <= 50 {
			typ2SsqHis[hzStr].Last50 = typ2SsqHis[hzStr].Last50 + 1
		}
		if lenSsqHis-i <= 100 {
			typ2SsqHis[hzStr].Last100 = typ2SsqHis[hzStr].Last100 + 1
		}
		if lenSsqHis-i <= 200 {
			typ2SsqHis[hzStr].Last200 = typ2SsqHis[hzStr].Last200 + 1
		}
		if lenSsqHis-i <= 500 {
			typ2SsqHis[hzStr].Last500 = typ2SsqHis[hzStr].Last500 + 1
		}
		if lenSsqHis-i <= 1000 {
			typ2SsqHis[hzStr].Last1000 = typ2SsqHis[hzStr].Last1000 + 1
		}
		if lenSsqHis-i <= 1500 {
			typ2SsqHis[hzStr].Last1500 = typ2SsqHis[hzStr].Last1500 + 1
		}
		if lenSsqHis-i <= 2000 {
			typ2SsqHis[hzStr].Last2000 = typ2SsqHis[hzStr].Last2000 + 1
		}
		if lenSsqHis-i <= 2500 {
			typ2SsqHis[hzStr].Last2500 = typ2SsqHis[hzStr].Last2500 + 1
		}
		if lenSsqHis-i <= 3500 {
			typ2SsqHis[hzStr].Last3500 = typ2SsqHis[hzStr].Last3500 + 1
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
		return kLens[i].Length > kLens[j].Length
	})

	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2SsqHis[typ])
	}
	return
}
