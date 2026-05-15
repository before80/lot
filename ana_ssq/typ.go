package ana_ssq

import (
	"github.com/before80/lot/db"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/models"
)

// CalSsqNotTxFrontHmSlice
//
//	@Description:
func CalSsqNotTxFrontHmSlice(hm6s, hm8s, hm13s, hm18s [][]string) (resSlice []string) {
	allFrontHmSlice := gen.Comb(gen.AllSsqFrontHms, 6)
	resSlice = make([]string, 0)
	res2St := make(map[string]struct{})

	hm62St, hm82St, hm132St, hm182St := make(map[string]struct{}), make(map[string]struct{}), make(map[string]struct{}), make(map[string]struct{})
	for _, v := range hm6s {
		for _, iv := range v {
			hm62St[iv] = struct{}{}
		}
	}

	for _, v := range hm8s {
		for _, iv := range v {
			hm82St[iv] = struct{}{}
		}
	}

	for _, v := range hm13s {
		for _, iv := range v {
			hm132St[iv] = struct{}{}
		}
	}

	for _, v := range hm18s {
		for _, iv := range v {
			hm182St[iv] = struct{}{}
		}
	}

	for _, v := range allFrontHmSlice {
		_, ok1 := hm62St[v]
		_, ok2 := hm82St[v]
		_, ok3 := hm132St[v]
		_, ok4 := hm182St[v]
		_, ok5 := res2St[v]
		if !ok1 && !ok2 && !ok3 && !ok4 && !ok5 {
			resSlice = append(resSlice, v)
		}
	}
	return
}

func CalSsqTxFrontHmSlice() (hm6s, hm8s, hm13s, hm18s [][]string) {
	// 找到类型为666663匹配历史开奖号码最多的前5种
	var moni666663s []models.SsqMoni
	db.DB.Where(map[string]interface{}{"typ": "666663"}).Order("cs desc,id asc").Limit(5).Find(&moni666663s)
	// 找到类型为855555匹配历史开奖号码最多的前5种
	var moni855555s []models.SsqMoni
	db.DB.Where(map[string]interface{}{"typ": "855555"}).Order("cs desc,id asc").Limit(5).Find(&moni855555s)
	// 找到类型为1344444匹配历史开奖号码最多的前5种
	var moni1344444s []models.SsqMoni
	db.DB.Where(map[string]interface{}{"typ": "1344444"}).Order("cs desc,id asc").Limit(5).Find(&moni1344444s)
	// 找到类型为1833333匹配历史开奖号码最多的前2种
	var moni1833333s []models.SsqMoni
	db.DB.Where(map[string]interface{}{"typ": "1833333"}).Order("cs desc,id asc").Limit(2).Find(&moni1833333s)

	hm6s = BuildFrontHmSlices(moni666663s, 5)
	hm8s = BuildFrontHmSlices(moni855555s, 5)
	hm13s = BuildFrontHmSlices(moni1344444s, 5)
	hm18s = BuildFrontHmSlices(moni1833333s, 2)
	return
}
