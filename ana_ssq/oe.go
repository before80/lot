package ana_ssq

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/before80/lot/dbop"
)

// CalSsqOe 计算双色球前后区号码的奇偶个数，前面的数字为奇数的个数，后面的数字为偶数的个数
//
//	@Description:
//	@param hmStrSlice
//	@return string
func CalSsqOe(hmStrSlice []string) string {
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

// CalSsqOeFromStr 计算双色球前后区号码的奇偶个数，前面的数字为奇数的个数，后面的数字为偶数的个数
//
//	@Description:
//	@param fullHmStr
//	@return string
func CalSsqOeFromStr(fullHmStr string) string {
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

// SsqOeQuShi 奇偶趋势 (奇偶->奇偶)
//
//	@Description:
//	@return []KeyWithLength
func SsqOeQuShi() []KeyWithLength {
	dlts, _ := dbop.ReadAllSsq(false)
	oeMs := make(map[string]int)
	for i, dlt := range dlts {
		if i == len(dlts)-1 {
			break
		}
		nextSsq := dlts[i+1]

		aeArrow := fmt.Sprintf("%s->%s", dlt.Oe, nextSsq.Oe)
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

// SsqOeHis 大乐透奇偶历史数据(按出现期数的大小,从大到小排序)
//
//	@Description:
//	@return res
func SsqOeHis(wg *sync.WaitGroup) (res []SsqHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2SsqHis := make(map[string]*SsqHis)
	// 初始化
	for k, v := range AllSsqOe2Count {
		typ2SsqHis[k] = &SsqHis{Typ: k, AllCount: v}
	}
	lenSsqHis := len(ZxSsqs)

	for i, dlt := range ZxSsqs {
		oe := dlt.Oe
		typ2SsqHis[oe].Cs = typ2SsqHis[oe].Cs + 1
		if lenSsqHis-i <= 10 {
			typ2SsqHis[oe].Last10 = typ2SsqHis[oe].Last10 + 1
		}
		if lenSsqHis-i <= 20 {
			typ2SsqHis[oe].Last20 = typ2SsqHis[oe].Last20 + 1
		}
		if lenSsqHis-i <= 30 {
			typ2SsqHis[oe].Last30 = typ2SsqHis[oe].Last30 + 1
		}
		if lenSsqHis-i <= 50 {
			typ2SsqHis[oe].Last50 = typ2SsqHis[oe].Last50 + 1
		}
		if lenSsqHis-i <= 100 {
			typ2SsqHis[oe].Last100 = typ2SsqHis[oe].Last100 + 1
		}
		if lenSsqHis-i <= 200 {
			typ2SsqHis[oe].Last200 = typ2SsqHis[oe].Last200 + 1
		}
		if lenSsqHis-i <= 500 {
			typ2SsqHis[oe].Last500 = typ2SsqHis[oe].Last500 + 1
		}
		if lenSsqHis-i <= 1000 {
			typ2SsqHis[oe].Last1000 = typ2SsqHis[oe].Last1000 + 1
		}
		if lenSsqHis-i <= 1500 {
			typ2SsqHis[oe].Last1500 = typ2SsqHis[oe].Last1500 + 1
		}
		if lenSsqHis-i <= 2000 {
			typ2SsqHis[oe].Last2000 = typ2SsqHis[oe].Last2000 + 1
		}
		if lenSsqHis-i <= 2500 {
			typ2SsqHis[oe].Last2500 = typ2SsqHis[oe].Last2500 + 1
		}
		if lenSsqHis-i <= 3500 {
			typ2SsqHis[oe].Last3500 = typ2SsqHis[oe].Last3500 + 1
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
