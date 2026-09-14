package ana_dlt

import (
	"fmt"
	"sort"
	"strings"

	"github.com/before80/lot/db"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/models"
)

// CalDltNotTxFrontHmSlice
//
//	@Description:
func CalDltNotTxFrontHmSlice(hm7s, hm11s, hm15s, hm19s [][]string) (resSlice []string) {
	allFrontHmSlice := gen.Comb(gen.AllDltFrontHms, 5)
	resSlice = make([]string, 0)
	res2St := make(map[string]struct{})

	hm72St, hm112St, hm152St, hm192St := make(map[string]struct{}), make(map[string]struct{}), make(map[string]struct{}), make(map[string]struct{})
	for _, v := range hm7s {
		for _, iv := range v {
			hm72St[iv] = struct{}{}
		}
	}

	for _, v := range hm11s {
		for _, iv := range v {
			hm112St[iv] = struct{}{}
		}
	}

	for _, v := range hm15s {
		for _, iv := range v {
			hm152St[iv] = struct{}{}
		}
	}

	for _, v := range hm19s {
		for _, iv := range v {
			hm192St[iv] = struct{}{}
		}
	}

	for _, v := range allFrontHmSlice {
		_, ok1 := hm72St[v]
		_, ok2 := hm112St[v]
		_, ok3 := hm152St[v]
		_, ok4 := hm192St[v]
		_, ok5 := res2St[v]
		if !ok1 && !ok2 && !ok3 && !ok4 && !ok5 {
			resSlice = append(resSlice, v)
		}
	}
	return
}

func CalDltTxFrontHmSlice() (hm7s, hm11s, hm15s, hm19s [][]string) {
	// 找到类型为77777匹配历史开奖号码最多的前5种
	var moni77777s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "77777"}).Order("cs desc,id asc").Limit(5).Find(&moni77777s)
	// 找到类型为116666匹配历史开奖号码最多的前5种
	var moni116666s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "116666"}).Order("cs desc,id asc").Limit(5).Find(&moni116666s)
	// 找到类型为155555匹配历史开奖号码最多的前5种
	var moni155555s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "155555"}).Order("cs desc,id asc").Limit(5).Find(&moni155555s)
	// 找到类型为194444匹配历史开奖号码最多的前2种
	var moni194444s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "194444"}).Order("cs desc,id asc").Limit(2).Find(&moni194444s)

	//fmt.Println(moni77777s)
	//fmt.Println(moni116666s)
	//fmt.Println(moni155555s)
	//fmt.Println(moni194444s)
	hm7s = BuildFrontHmSlices(moni77777s, 5)
	hm11s = BuildFrontHmSlices(moni116666s, 5)
	hm15s = BuildFrontHmSlices(moni155555s, 5)
	hm19s = BuildFrontHmSlices(moni194444s, 2)
	return
}

var TxMap map[string]map[string][]string

// GetDltTxMap 获取大乐透前区号码的所有类型(T1到T19共19种类型)
func GetDltTxMap() map[string]map[string][]string {
	TxMap = getDltTxMap()
	return TxMap
}

// getDltTxMap 获取大乐透前区号码的所有类型(T1到T19共19种类型)
func getDltTxMap() (txMap map[string]map[string][]string) {
	var moni77777s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "77777"}).Order("cs desc,id asc").Limit(5).Find(&moni77777s)
	// 找到类型为116666匹配历史开奖号码最多的前5种
	var moni116666s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "116666"}).Order("cs desc,id asc").Limit(5).Find(&moni116666s)
	// 找到类型为155555匹配历史开奖号码最多的前5种
	var moni155555s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "155555"}).Order("cs desc,id asc").Limit(5).Find(&moni155555s)
	// 找到类型为194444匹配历史开奖号码最多的前2种
	var moni194444s []models.DltMoni
	db.DB.Where(map[string]interface{}{"typ": "194444"}).Order("cs desc,id asc").Limit(2).Find(&moni194444s)

	for i, moni := range moni77777s {
		t := fmt.Sprintf("T%d", i+1)
		txMap[t] = make(map[string][]string)
		txMap[t]["A"] = strings.Split(moni.A, ",")
		txMap[t]["B"] = strings.Split(moni.B, ",")
		txMap[t]["C"] = strings.Split(moni.C, ",")
		txMap[t]["D"] = strings.Split(moni.D, ",")
		txMap[t]["E"] = strings.Split(moni.E, ",")
	}

	for i, moni := range moni116666s {
		t := fmt.Sprintf("T%d", i+5)
		txMap[t] = make(map[string][]string)
		txMap[t]["A"] = strings.Split(moni.A, ",")
		txMap[t]["B"] = strings.Split(moni.B, ",")
		txMap[t]["C"] = strings.Split(moni.C, ",")
		txMap[t]["D"] = strings.Split(moni.D, ",")
		txMap[t]["E"] = strings.Split(moni.E, ",")
	}

	for i, moni := range moni155555s {
		t := fmt.Sprintf("T%d", i+10)
		txMap[t] = make(map[string][]string)
		txMap[t]["A"] = strings.Split(moni.A, ",")
		txMap[t]["B"] = strings.Split(moni.B, ",")
		txMap[t]["C"] = strings.Split(moni.C, ",")
		txMap[t]["D"] = strings.Split(moni.D, ",")
		txMap[t]["E"] = strings.Split(moni.E, ",")
	}

	for i, moni := range moni194444s {
		t := fmt.Sprintf("T%d", i+15)
		txMap[t] = make(map[string][]string)
		txMap[t]["A"] = strings.Split(moni.A, ",")
		txMap[t]["B"] = strings.Split(moni.B, ",")
		txMap[t]["C"] = strings.Split(moni.C, ",")
		txMap[t]["D"] = strings.Split(moni.D, ",")
		txMap[t]["E"] = strings.Split(moni.E, ",")
	}

	return
}

var FrontHmStr2StMap = make(map[string]map[string]struct{})

func CalTxFrontHmStr2StMap() {
	hm7s, hm11s, hm15s, hm19s := CalDltTxFrontHmSlice()
	OtherTs := CalDltNotTxFrontHmSlice(hm7s, hm11s, hm15s, hm19s)

	for i := 0; i < len(hm7s); i++ {
		ti := fmt.Sprintf("T%d", i+1)
		FrontHmStr2StMap[ti] = make(map[string]struct{})
		for _, hmStr := range hm7s[i] {
			FrontHmStr2StMap[ti][hmStr] = struct{}{}
		}
	}

	for i := 5; i < len(hm11s); i++ {
		ti := fmt.Sprintf("T%d", i+1)
		FrontHmStr2StMap[ti] = make(map[string]struct{})
		for _, hmStr := range hm11s[i] {
			FrontHmStr2StMap[ti][hmStr] = struct{}{}
		}
	}

	for i := 10; i < len(hm15s); i++ {
		ti := fmt.Sprintf("T%d", i+1)
		FrontHmStr2StMap[ti] = make(map[string]struct{})
		for _, hmStr := range hm15s[i] {
			FrontHmStr2StMap[ti][hmStr] = struct{}{}
		}
	}

	for i := 15; i < len(hm19s); i++ {
		ti := fmt.Sprintf("T%d", i+1)
		FrontHmStr2StMap[ti] = make(map[string]struct{})
		for _, hmStr := range hm19s[i] {
			FrontHmStr2StMap[ti][hmStr] = struct{}{}
		}
	}

	ti := fmt.Sprintf("OtherT")
	FrontHmStr2StMap[ti] = make(map[string]struct{})
	for _, hmStr := range OtherTs {
		FrontHmStr2StMap[ti][hmStr] = struct{}{}
	}
}

func CalDltBelongToWhichTx(frontHms []string) (tx string) {
	sort.Strings(frontHms)

	return
}
