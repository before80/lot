package ana_dlt

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/before80/lot/db"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/models"
)

func UpdateDltHz() {
	dlts, _ := dbop.ReadAllDlt(false)
	var hz int
	var aeHz, oe, qzh string
	for _, dlt := range dlts {
		newCombs := []string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2}
		hz = CalDltHz(newCombs)
		aeHz = DltHzABCDE(hz)
		oe = CalDltOe(newCombs)
		qzh = CalDltQzh(newCombs)
		db.DB.Model(&dlt).Updates(map[string]interface{}{"oe": oe, "hz": hz, "ae_hz": aeHz, "qzh": qzh})
	}
	fmt.Println("UpdateDltHz done")
}

// CalDltHz 计算大乐透前后区号码的和值
//
//	@Description:
//	@param lotteryDrawResults
//	@return total
func CalDltHz(hmStrSlice []string) (total int) {
	for _, hmStr := range hmStrSlice {
		num, _ := strconv.Atoi(hmStr)
		total += num
	}
	return total
}

// CalDltHzFromStr 计算大乐透前后区号码的和值
//
//	@Description:
//	@param fullHmStr
//	@return total
func CalDltHzFromStr(fullHmStr string) (total int) {
	fullHmStr = strings.Replace(fullHmStr, "|", ",", -1)
	fullHms := strings.Split(fullHmStr, ",")
	for _, fullHm := range fullHms {
		num, _ := strconv.Atoi(fullHm)
		total += num
	}
	return total
}

func CalDltHz1(dlt models.Dlt) int {
	var total, num int
	num, _ = strconv.Atoi(dlt.F1)
	total += num
	num, _ = strconv.Atoi(dlt.F2)
	total += num
	num, _ = strconv.Atoi(dlt.F3)
	total += num
	num, _ = strconv.Atoi(dlt.F4)
	total += num
	num, _ = strconv.Atoi(dlt.F5)
	total += num
	num, _ = strconv.Atoi(dlt.B1)
	total += num
	num, _ = strconv.Atoi(dlt.B2)
	total += num
	return total
}

type MaxMin struct {
	Min int
	Max int
}

var HzABCDEMs = map[string]MaxMin{
	"A": MaxMin{Min: 18, Max: 52},
	"B": MaxMin{Min: 53, Max: 86},
	"C": MaxMin{Min: 87, Max: 120},
	"D": MaxMin{Min: 121, Max: 154},
	"E": MaxMin{Min: 155, Max: 188},
}

func DltHzABCDE(hz int) string {
	if hz <= 52 {
		return "A"
	}
	if hz <= 86 {
		return "B"
	}
	if hz <= 120 {
		return "C"
	}
	if hz <= 154 {
		return "D"
	}
	if hz <= 188 {
		return "E"
	}
	return "E"
}

// DltHzQuShi1 大乐透和值趋势(和值->和值) 注意需要提前之前 InitDlts()
//
//	@Description:
//	@return []KeyWithLength
func DltHzQuShi1() []KeyWithLength {
	aeMs := make(map[string]int)
	for i, dlt := range ZxDlts {
		if i == len(ZxDlts)-1 {
			break
		}
		nextDlt := ZxDlts[i+1]
		aeArrow := fmt.Sprintf("%d->%d", dlt.Hz, nextDlt.Hz)
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
		if kLens[i].Length == kLens[j].Length {
			return kLens[i].Key > kLens[j].Key
		}
		return kLens[i].Length > kLens[j].Length
	})
	return kLens
}

type DltHzChQuShi struct {
	HzCh      string
	NextHzCh  string
	Cs        int
	Hz2NextHz []string
}

// DltHzQuShi2 大乐透和值趋势(和值Ch->和值Ch) 注意需要提前之前 InitDlts()
//
//	@Description:
//	@return res
func DltHzQuShi2() (res []DltHzChQuShi) {
	hz2ChQuShiSt := make(map[string]*DltHzChQuShi)
	for i, dlt := range ZxDlts {
		if i == len(ZxDlts)-1 {
			break
		}
		nextDlt := ZxDlts[i+1]
		nextHzABCDE := DltHzABCDE(nextDlt.Hz)
		hzABCDE := DltHzABCDE(dlt.Hz)
		aeArrow := fmt.Sprintf("%s->%s", hzABCDE, nextHzABCDE)
		hz2NextHzStr := fmt.Sprintf("%d->%d", dlt.Hz, nextDlt.Hz)
		if _, ok := hz2ChQuShiSt[aeArrow]; !ok {
			hz2ChQuShiSt[aeArrow] = &DltHzChQuShi{
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

// DltHzHis 大乐透和值历史数据(按出现期数的大小,从大到小排序)
//
//	@Description:
//	@return res
func DltHzHis(wg *sync.WaitGroup) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2DltHis := make(map[string]*DltHis)
	// 初始化
	for k, v := range AllDltHz2Count {
		typ2DltHis[k] = &DltHis{Typ: k, AllCount: v}
	}
	lenDltHis := len(ZxDlts)

	for i, dlt := range ZxDlts {
		hzStr := strconv.Itoa(dlt.Hz)
		typ2DltHis[hzStr].Cs = typ2DltHis[hzStr].Cs + 1
		if lenDltHis-i <= 10 {
			typ2DltHis[hzStr].Last10 = typ2DltHis[hzStr].Last10 + 1
		}
		if lenDltHis-i <= 20 {
			typ2DltHis[hzStr].Last20 = typ2DltHis[hzStr].Last20 + 1
		}
		if lenDltHis-i <= 30 {
			typ2DltHis[hzStr].Last30 = typ2DltHis[hzStr].Last30 + 1
		}
		if lenDltHis-i <= 50 {
			typ2DltHis[hzStr].Last50 = typ2DltHis[hzStr].Last50 + 1
		}
		if lenDltHis-i <= 100 {
			typ2DltHis[hzStr].Last100 = typ2DltHis[hzStr].Last100 + 1
		}
		if lenDltHis-i <= 200 {
			typ2DltHis[hzStr].Last200 = typ2DltHis[hzStr].Last200 + 1
		}
		if lenDltHis-i <= 500 {
			typ2DltHis[hzStr].Last500 = typ2DltHis[hzStr].Last500 + 1
		}
		if lenDltHis-i <= 1000 {
			typ2DltHis[hzStr].Last1000 = typ2DltHis[hzStr].Last1000 + 1
		}
		if lenDltHis-i <= 1500 {
			typ2DltHis[hzStr].Last1500 = typ2DltHis[hzStr].Last1500 + 1
		}
		if lenDltHis-i <= 2000 {
			typ2DltHis[hzStr].Last2000 = typ2DltHis[hzStr].Last2000 + 1
		}
		if lenDltHis-i <= 2500 {
			typ2DltHis[hzStr].Last2500 = typ2DltHis[hzStr].Last2500 + 1
		}
		if lenDltHis-i <= 3500 {
			typ2DltHis[hzStr].Last3500 = typ2DltHis[hzStr].Last3500 + 1
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
		if kLens[i].Length == kLens[j].Length {
			return kLens[i].Key > kLens[j].Key
		}
		return kLens[i].Length > kLens[j].Length
	})

	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2DltHis[typ])
	}
	return
}

// DltEqHzHis 设备的大乐透和值历史数据(按出现期数的大小,从大到小排序)
func DltEqHzHis(wg *sync.WaitGroup, eqNumCount int) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2DltHis := make(map[string]*DltHis)
	// 初始化
	for k, v := range AllDltHz2Count {
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

		hzStr := strconv.Itoa(dlt.Hz)
		typ2DltHis[hzStr].Cs = typ2DltHis[hzStr].Cs + 1
		if lenDltHis-i <= 10 {
			typ2DltHis[hzStr].Last10 = typ2DltHis[hzStr].Last10 + 1
		}
		if lenDltHis-i <= 20 {
			typ2DltHis[hzStr].Last20 = typ2DltHis[hzStr].Last20 + 1
		}
		if lenDltHis-i <= 30 {
			typ2DltHis[hzStr].Last30 = typ2DltHis[hzStr].Last30 + 1
		}
		if lenDltHis-i <= 50 {
			typ2DltHis[hzStr].Last50 = typ2DltHis[hzStr].Last50 + 1
		}
		if lenDltHis-i <= 100 {
			typ2DltHis[hzStr].Last100 = typ2DltHis[hzStr].Last100 + 1
		}
		if lenDltHis-i <= 200 {
			typ2DltHis[hzStr].Last200 = typ2DltHis[hzStr].Last200 + 1
		}
		if lenDltHis-i <= 500 {
			typ2DltHis[hzStr].Last500 = typ2DltHis[hzStr].Last500 + 1
		}
		if lenDltHis-i <= 1000 {
			typ2DltHis[hzStr].Last1000 = typ2DltHis[hzStr].Last1000 + 1
		}
		if lenDltHis-i <= 1500 {
			typ2DltHis[hzStr].Last1500 = typ2DltHis[hzStr].Last1500 + 1
		}
		if lenDltHis-i <= 2000 {
			typ2DltHis[hzStr].Last2000 = typ2DltHis[hzStr].Last2000 + 1
		}
		if lenDltHis-i <= 2500 {
			typ2DltHis[hzStr].Last2500 = typ2DltHis[hzStr].Last2500 + 1
		}
		if lenDltHis-i <= 3500 {
			typ2DltHis[hzStr].Last3500 = typ2DltHis[hzStr].Last3500 + 1
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
		if kLens[i].Length == kLens[j].Length {
			return kLens[i].Key > kLens[j].Key
		}
		return kLens[i].Length > kLens[j].Length
	})

	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2DltHis[typ])
	}
	return
}
