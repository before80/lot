package ana_dlt

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/before80/lot/dbop"
)

// DltOeTable 大乐透奇偶的所有可能
var DltOeTable = []string{"07", "16", "25", "34", "43", "52", "61", "70"}

// CalDltOe 计算大乐透前后区号码的奇偶个数，前面的数字为奇数的个数，后面的数字为偶数的个数
//
//	@Description:
//	@param hmStrSlice
//	@return string
func CalDltOe(hmStrSlice []string) string {
	var oNum, eNum int
	for _, hmStr := range hmStrSlice {
		num, _ := strconv.Atoi(hmStr)
		if num%2 == 0 { // 偶数个数
			eNum++
		} else {
			oNum++
		}
	}
	return fmt.Sprintf("%d%d", oNum, eNum)
}

// CalDltOeFromStr 计算大乐透前后区号码的奇偶个数，前面的数字为奇数的个数，后面的数字为偶数的个数
//
//	@Description:
//	@param fullHmStr
//	@return string
func CalDltOeFromStr(fullHmStr string) string {
	fullHmStr = strings.Replace(fullHmStr, "|", ",", -1)
	fullHms := strings.Split(fullHmStr, ",")
	var oNum, eNum int
	for _, hmStr := range fullHms {
		num, _ := strconv.Atoi(hmStr)
		if num%2 == 0 { // 偶数个数
			eNum++
		} else {
			oNum++
		}
	}
	return fmt.Sprintf("%d%d", oNum, eNum)
}

// DltOeQuShi 奇偶趋势 (奇偶->奇偶)
//
//	@Description:
//	@return []KeyWithLength
func DltOeQuShi() []KeyWithLength {
	dlts, _ := dbop.ReadAllDlt(false)
	oeMs := make(map[string]int)
	for i, dlt := range dlts {
		if i == len(dlts)-1 {
			break
		}
		nextDlt := dlts[i+1]

		aeArrow := fmt.Sprintf("%s->%s", dlt.Oe, nextDlt.Oe)
		if _, ok := oeMs[aeArrow]; !ok {
			oeMs[aeArrow] = 1
		} else {
			oeMs[aeArrow]++
		}
	}

	// 对oeMs按照存放的值的大小进行排序
	kLens := make([]KeyWithLength, 0, len(oeMs))

	for aeArrow, v := range oeMs {
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

// DltOeHis 大乐透奇偶历史数据(按出现期数的大小,从大到小排序)
//
//	@Description:
//	@return res
func DltOeHis(wg *sync.WaitGroup) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2DltHis := make(map[string]*DltHis)
	// 初始化
	for k, v := range AllDltOe2Count {
		typ2DltHis[k] = &DltHis{Typ: k, AllCount: v}
	}
	lenDltHis := len(ZxDlts)

	for i, dlt := range ZxDlts {
		oe := dlt.Oe
		typ2DltHis[oe].Cs = typ2DltHis[oe].Cs + 1
		if lenDltHis-i <= 10 {
			typ2DltHis[oe].Last10 = typ2DltHis[oe].Last10 + 1
		}
		if lenDltHis-i <= 20 {
			typ2DltHis[oe].Last20 = typ2DltHis[oe].Last20 + 1
		}
		if lenDltHis-i <= 30 {
			typ2DltHis[oe].Last30 = typ2DltHis[oe].Last30 + 1
		}
		if lenDltHis-i <= 50 {
			typ2DltHis[oe].Last50 = typ2DltHis[oe].Last50 + 1
		}
		if lenDltHis-i <= 100 {
			typ2DltHis[oe].Last100 = typ2DltHis[oe].Last100 + 1
		}
		if lenDltHis-i <= 200 {
			typ2DltHis[oe].Last200 = typ2DltHis[oe].Last200 + 1
		}
		if lenDltHis-i <= 500 {
			typ2DltHis[oe].Last500 = typ2DltHis[oe].Last500 + 1
		}
		if lenDltHis-i <= 1000 {
			typ2DltHis[oe].Last1000 = typ2DltHis[oe].Last1000 + 1
		}
		if lenDltHis-i <= 1500 {
			typ2DltHis[oe].Last1500 = typ2DltHis[oe].Last1500 + 1
		}
		if lenDltHis-i <= 2000 {
			typ2DltHis[oe].Last2000 = typ2DltHis[oe].Last2000 + 1
		}
		if lenDltHis-i <= 2500 {
			typ2DltHis[oe].Last2500 = typ2DltHis[oe].Last2500 + 1
		}
		if lenDltHis-i <= 3500 {
			typ2DltHis[oe].Last3500 = typ2DltHis[oe].Last3500 + 1
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

func DltEqOeHis(wg *sync.WaitGroup, eqNumCount int) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2DltHis := make(map[string]*DltHis)
	// 初始化
	for k, v := range AllDltOe2Count {
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

		oe := dlt.Oe
		typ2DltHis[oe].Cs = typ2DltHis[oe].Cs + 1
		if lenDltHis-i <= 10 {
			typ2DltHis[oe].Last10 = typ2DltHis[oe].Last10 + 1
		}
		if lenDltHis-i <= 20 {
			typ2DltHis[oe].Last20 = typ2DltHis[oe].Last20 + 1
		}
		if lenDltHis-i <= 30 {
			typ2DltHis[oe].Last30 = typ2DltHis[oe].Last30 + 1
		}
		if lenDltHis-i <= 50 {
			typ2DltHis[oe].Last50 = typ2DltHis[oe].Last50 + 1
		}
		if lenDltHis-i <= 100 {
			typ2DltHis[oe].Last100 = typ2DltHis[oe].Last100 + 1
		}
		if lenDltHis-i <= 200 {
			typ2DltHis[oe].Last200 = typ2DltHis[oe].Last200 + 1
		}
		if lenDltHis-i <= 500 {
			typ2DltHis[oe].Last500 = typ2DltHis[oe].Last500 + 1
		}
		if lenDltHis-i <= 1000 {
			typ2DltHis[oe].Last1000 = typ2DltHis[oe].Last1000 + 1
		}
		if lenDltHis-i <= 1500 {
			typ2DltHis[oe].Last1500 = typ2DltHis[oe].Last1500 + 1
		}
		if lenDltHis-i <= 2000 {
			typ2DltHis[oe].Last2000 = typ2DltHis[oe].Last2000 + 1
		}
		if lenDltHis-i <= 2500 {
			typ2DltHis[oe].Last2500 = typ2DltHis[oe].Last2500 + 1
		}
		if lenDltHis-i <= 3500 {
			typ2DltHis[oe].Last3500 = typ2DltHis[oe].Last3500 + 1
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
