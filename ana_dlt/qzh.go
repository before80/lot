package ana_dlt

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/before80/lot/dbop"
)

// CalDltQzh 计算大乐透全区5个号码中，前1到17，中18，后19到35，各自出现的个数
//
//	@Description:
//	@param lotteryDrawResults
//	@return string
func CalDltQzh(hmStrSlice []string) string {
	var qNum, zNum, HNum int
	for i, hmStr := range hmStrSlice {
		if i < 5 {
			num, _ := strconv.Atoi(hmStr)
			if num < 18 {
				qNum++
			}
			if num == 18 {
				zNum++
			}
			if num > 18 {
				HNum++
			}
		}
	}
	return fmt.Sprintf("%d%d%d", qNum, zNum, HNum)
}

// DltQzhQuShi
//
//	@Description:
//	@return []KeyWithLength
func DltQzhQuShi() []KeyWithLength {
	dlts, _ := dbop.ReadAllDlt(false)
	qzhMs := make(map[string]int)
	for i, dlt := range dlts {
		if i == len(dlts)-1 {
			break
		}
		nextDlt := dlts[i+1]

		aeArrow := fmt.Sprintf("%s->%s", dlt.Qzh, nextDlt.Qzh)
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

// DltQzhHis 大乐透前中后历史数据(按出现期数的大小,从大到小排序)
//
//	@Description:
//	@param wg
//	@return res
func DltQzhHis(wg *sync.WaitGroup) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2DltHis := make(map[string]*DltHis)

	// 初始化
	for k, v := range AllDltQzh2Count {
		typ2DltHis[k] = &DltHis{Typ: k, AllCount: v}
	}

	lenDltHis := len(ZxDlts)

	for i, dlt := range ZxDlts {
		qzhStr := dlt.Qzh

		typ2DltHis[qzhStr].Cs = typ2DltHis[qzhStr].Cs + 1
		if lenDltHis-i <= 10 {
			typ2DltHis[qzhStr].Last10 = typ2DltHis[qzhStr].Last10 + 1
		}
		if lenDltHis-i <= 20 {
			typ2DltHis[qzhStr].Last20 = typ2DltHis[qzhStr].Last20 + 1
		}
		if lenDltHis-i <= 30 {
			typ2DltHis[qzhStr].Last30 = typ2DltHis[qzhStr].Last30 + 1
		}
		if lenDltHis-i <= 50 {
			typ2DltHis[qzhStr].Last50 = typ2DltHis[qzhStr].Last50 + 1
		}
		if lenDltHis-i <= 100 {
			typ2DltHis[qzhStr].Last100 = typ2DltHis[qzhStr].Last100 + 1
		}
		if lenDltHis-i <= 200 {
			typ2DltHis[qzhStr].Last200 = typ2DltHis[qzhStr].Last200 + 1
		}
		if lenDltHis-i <= 500 {
			typ2DltHis[qzhStr].Last500 = typ2DltHis[qzhStr].Last500 + 1
		}
		if lenDltHis-i <= 1000 {
			typ2DltHis[qzhStr].Last1000 = typ2DltHis[qzhStr].Last1000 + 1
		}
		if lenDltHis-i <= 1500 {
			typ2DltHis[qzhStr].Last1500 = typ2DltHis[qzhStr].Last1500 + 1
		}
		if lenDltHis-i <= 2000 {
			typ2DltHis[qzhStr].Last2000 = typ2DltHis[qzhStr].Last2000 + 1
		}
		if lenDltHis-i <= 2500 {
			typ2DltHis[qzhStr].Last2500 = typ2DltHis[qzhStr].Last2500 + 1
		}
		if lenDltHis-i <= 3500 {
			typ2DltHis[qzhStr].Last3500 = typ2DltHis[qzhStr].Last3500 + 1
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
		return kLens[i].Length > kLens[j].Length
	})

	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2DltHis[typ])
	}
	return
}

// DltEqQzhHis 设备的大乐透前中后历史数据(按出现期数的大小,从大到小排序)
func DltEqQzhHis(wg *sync.WaitGroup, eqNumCount int) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2DltHis := make(map[string]*DltHis)

	// 初始化
	for k, v := range AllDltQzh2Count {
		typ2DltHis[k] = &DltHis{Typ: k, AllCount: v}
	}

	lenDltHis := 0
	for _, dlt := range ZxDlts {
		if dlt.DrawNum < "11001" {
			continue
		}
		if dlt.EquipmentCount != eqNumCount {
			continue
		}
		lenDltHis++
	}

	i := 0

	for _, dlt := range ZxDlts {
		if dlt.DrawNum < "11001" {
			continue
		}
		if dlt.EquipmentCount != eqNumCount {
			continue
		}

		qzhStr := dlt.Qzh

		typ2DltHis[qzhStr].Cs = typ2DltHis[qzhStr].Cs + 1
		if lenDltHis-i <= 10 {
			typ2DltHis[qzhStr].Last10 = typ2DltHis[qzhStr].Last10 + 1
		}
		if lenDltHis-i <= 20 {
			typ2DltHis[qzhStr].Last20 = typ2DltHis[qzhStr].Last20 + 1
		}
		if lenDltHis-i <= 30 {
			typ2DltHis[qzhStr].Last30 = typ2DltHis[qzhStr].Last30 + 1
		}
		if lenDltHis-i <= 50 {
			typ2DltHis[qzhStr].Last50 = typ2DltHis[qzhStr].Last50 + 1
		}
		if lenDltHis-i <= 100 {
			typ2DltHis[qzhStr].Last100 = typ2DltHis[qzhStr].Last100 + 1
		}
		if lenDltHis-i <= 200 {
			typ2DltHis[qzhStr].Last200 = typ2DltHis[qzhStr].Last200 + 1
		}
		if lenDltHis-i <= 500 {
			typ2DltHis[qzhStr].Last500 = typ2DltHis[qzhStr].Last500 + 1
		}
		if lenDltHis-i <= 1000 {
			typ2DltHis[qzhStr].Last1000 = typ2DltHis[qzhStr].Last1000 + 1
		}
		if lenDltHis-i <= 1500 {
			typ2DltHis[qzhStr].Last1500 = typ2DltHis[qzhStr].Last1500 + 1
		}
		if lenDltHis-i <= 2000 {
			typ2DltHis[qzhStr].Last2000 = typ2DltHis[qzhStr].Last2000 + 1
		}
		if lenDltHis-i <= 2500 {
			typ2DltHis[qzhStr].Last2500 = typ2DltHis[qzhStr].Last2500 + 1
		}
		if lenDltHis-i <= 3500 {
			typ2DltHis[qzhStr].Last3500 = typ2DltHis[qzhStr].Last3500 + 1
		}
		i++
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
		return kLens[i].Length > kLens[j].Length
	})

	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2DltHis[typ])
	}
	return
}
