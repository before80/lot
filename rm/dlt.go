package rm

import (
	"slices"
	"strconv"
	"strings"

	"github.com/before80/lot/cnt"
)

// DltFullHmWhenItNotInSpOes 移除不需要的奇偶组合号码
//
//	@Description:
//	@param allFullHms
//	@param oes
//	@return afterRmFullHms
func DltFullHmWhenItNotInSpOes(allFullHms []string, oes []string) (afterRmFullHms []string) {
	for _, fullHm := range allFullHms {
		curOe := cnt.DltFullHmContainOe(fullHm)
		if !slices.Contains(oes, curOe) {
			afterRmFullHms = append(afterRmFullHms, fullHm)
		}
	}
	return
}

// DltFullHmWhenItInSpHms 保留含有指定号码的号码
//
//	@Description:
//	@param allFullHms
//	@param spHms
//	@return afterRmFullHms
func DltFullHmWhenItInSpHms(allFullHms []string, spHms []string) (afterRmFullHms []string) {
	for _, fullHm := range allFullHms {
		hms := strings.Split(strings.Replace(fullHm, "|", ",", -1), ",")
		if cnt.HasIntersection(hms, spHms) {
			afterRmFullHms = append(afterRmFullHms, fullHm)
		}
	}
	return
}

// DltFullHmWhenItNotInSpHms 移除含有指定号码的号码
//
//	@Description:
//	@param allFullHms
//	@param spHms
//	@return afterRmFullHms
func DltFullHmWhenItNotInSpHms(allFullHms []string, spHms []string) (afterRmFullHms []string) {
	for _, fullHm := range allFullHms {
		hms := strings.Split(strings.Replace(fullHm, "|", ",", -1), ",")
		if !cnt.HasIntersection(hms, spHms) {
			afterRmFullHms = append(afterRmFullHms, fullHm)
		}
	}
	return
}

// DltFullHmWhenItNotInSpBackHms 移除不在指定后区号码的号码
//
//	@Description:
//	@param allFullHms
//	@param spBackHms
//	@return afterRmFullHms
func DltFullHmWhenItNotInSpBackHms(allFullHms []string, spBackHms []string) (afterRmFullHms []string) {
	for _, fullHm := range allFullHms {
		hmStr := strings.Split(fullHm, "|")
		if len(hmStr) != 2 {
			continue
		}
		backHms := strings.Split(hmStr[1], ",")
		if !cnt.HasIntersection(spBackHms, backHms) {
			afterRmFullHms = append(afterRmFullHms, fullHm)
		}
	}
	return
}

// DltFullHmWhenItInSpBackHms 保留在指定后区号码的号码
//
//	@Description:
//	@param allFullHms
//	@param spBackHms
//	@return afterRmFullHms
func DltFullHmWhenItInSpBackHms(allFullHms []string, spBackHms []string) (afterRmFullHms []string) {
	for _, fullHm := range allFullHms {
		hmStr := strings.Split(fullHm, "|")
		if len(hmStr) != 2 {
			continue
		}
		backHms := strings.Split(hmStr[1], ",")
		if cnt.HasIntersection(spBackHms, backHms) {
			afterRmFullHms = append(afterRmFullHms, fullHm)
		}
	}
	return
}

// DltFullHmWhenItNotInHzRange 移除不在和值范围内的号码
//
//	@Description:
//	@param allFullHms
//	@param HzRangeMin
//	@param HzRangeMax
//	@return afterRmFullHms
func DltFullHmWhenItNotInHzRange(allFullHms []string, HzRangeMin, HzRangeMax int) (afterRmFullHms []string) {
	for _, fullHm := range allFullHms {
		hms := strings.Split(strings.Replace(fullHm, "|", ",", -1), ",")
		curTotal := 0
		for _, hm := range hms {
			hmNum, _ := strconv.Atoi(hm)
			curTotal += hmNum
		}
		if curTotal >= HzRangeMin && curTotal <= HzRangeMax {
			afterRmFullHms = append(afterRmFullHms, fullHm)
		}
	}
	return
}

// DltFullHmWhenItInHisFrontHms 移除含有历史前区号码的号码
//
//	@Description:
//	@param allFullHms
//	@param hisFrontHms
//	@return afterRmFullHms
func DltFullHmWhenItInHisFrontHms(allFullHms []string, hisFrontHms []string) (afterRmFullHms []string) {
	for _, fullHm := range allFullHms {
		hmStr := strings.Split(fullHm, "|")
		if len(hmStr) != 2 {
			continue
		}
		frontHms := hmStr[0]
		if !slices.Contains(hisFrontHms, frontHms) {
			afterRmFullHms = append(afterRmFullHms, fullHm)
		}
	}
	return
}
