package gen

import (
	"fmt"
	"sort"
	"strings"
)

var AllSsqFrontHms = []string{
	"01", "02", "03", "04", "05", "06", "07",
	"08", "09", "10", "11", "12", "13", "14",
	"15", "16", "17", "18", "19", "20", "21",
	"22", "23", "24", "25", "26", "27", "28",
	"29", "30", "31", "32", "33",
}

var AllSsqBackHms = []string{
	"01", "02", "03", "04", "05", "06", "07", "08",
	"09", "10", "11", "12", "13", "14", "15", "16",
}

// AllSsqOes 所有双色球奇偶可能
var AllSsqOes = []string{"07", "16", "25", "34", "43", "52", "61", "70"}

func SsqFrontHmStrSliceFromAeStr(a, b, c, d, e, f string) (result []string) {
	as := strings.Split(a, ",")
	bs := strings.Split(b, ",")
	cs := strings.Split(c, ",")
	ds := strings.Split(d, ",")
	es := strings.Split(e, ",")
	fs := strings.Split(f, ",")

	for _, ia := range as {
		for _, ib := range bs {
			for _, ic := range cs {
				for _, id := range ds {
					for _, ie := range es {
						for _, iif := range fs {
							s := []string{ia, ib, ic, id, ie, iif}
							sort.Strings(s)
							result = append(result, fmt.Sprintf("%s,%s,%s,%s,%s,%s", s[0], s[1], s[2], s[3], s[4], s[5]))
						}
					}
				}
			}
		}
	}

	return result
}

// SsqFullHmSliceFromAeStr 生成指定ae字符串可以生成的开奖全部7个号码
//
//	@Description:
//	@param a
//	@param b
//	@param c
//	@param d
//	@param e
//	@return result
func SsqFullHmSliceFromAeStr(a, b, c, d, e, f string) (result []string) {
	as := strings.Split(a, ",")
	bs := strings.Split(b, ",")
	cs := strings.Split(c, ",")
	ds := strings.Split(d, ",")
	es := strings.Split(e, ",")
	fs := strings.Split(f, ",")
	backCombs := AllSsqBackHms
	for _, ia := range as {
		for _, ib := range bs {
			for _, ic := range cs {
				for _, id := range ds {
					for _, ie := range es {
						for _, iif := range fs {
							s := []string{ia, ib, ic, id, ie, iif}
							sort.Strings(s)
							for _, comb := range backCombs {
								result = append(result, fmt.Sprintf("%s,%s,%s,%s,%s,%s|%s", s[0], s[1], s[2], s[3], s[4], s[5], comb))
							}
						}
					}
				}
			}
		}
	}

	return result
}
