package cnt

import (
	"fmt"
	"strconv"
	"strings"
)

// DltFullHmContainOe 计算一注大乐透号码中奇偶个数
//
//	@Description:
//	@param fullHm
//	@return string
func DltFullHmContainOe(fullHm string) string {
	hms := strings.Split(strings.Replace(fullHm, "|", ",", -1), ",")
	var oNum, eNum int
	for _, hmStr := range hms {
		hm, _ := strconv.Atoi(hmStr)
		if hm%2 == 0 {
			eNum++
		} else {
			oNum++
		}
	}
	return fmt.Sprintf("%d%d", oNum, eNum)
}

// HasIntersection 计算两个切片是否存在交集
//
//	@Description:
//	@param a
//	@param b
//	@return bool
func HasIntersection(a, b []string) bool {
	set := make(map[string]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := set[v]; ok {
			return true
		}
	}
	return false
}
