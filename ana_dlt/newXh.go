package ana_dlt

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"time"

	"github.com/before80/lot/dbop"
	"github.com/before80/lot/gen"
)

func DltXHaoForWeb1(xuHaoSt *XuHaoSt, zhuShu int) (finalSelectFullHmStrSlice []string) {
	allDltFrontHms := ShuffleAllDltFrontHms()
	allDltBackHms := ShuffleAllDltBackHms()

	dlts, _ := dbop.ReadAllDlt(false)

	c4ms := make(map[string][]string)
	for _, dlt := range dlts {
		var ic4s []string
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 1)...)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 2, 2)...)

		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], fmt.Sprintf("%s,%s,%s,%s,%s|%s,%s", dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2))
		}
	}

	var hisFullHmStrSlice = make([]string, 0, len(dlts))
	for _, dlt := range dlts {
		hisFullHmStrSlice = append(hisFullHmStrSlice, fmt.Sprintf("%s,%s,%s,%s,%s|%s,%s", dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2))
	}

	// 获取所有类型对应的 moni
	t2MoniABCDEs := GenTx2MoniABCDE()

	finalCombCount := 0
	finalCombMaxCount := 6666
	finalCombForRandSlice := make([]string, 0, finalCombMaxCount)
	for a := 0; a <= 30; a++ {
		for b := a + 1; b <= 31; b++ {
			for c := b + 1; c <= 32; c++ {
				for d := c + 1; d <= 33; d++ {
				LabelForE:
					for e := d + 1; e <= 34; e++ {
						front := []string{
							allDltFrontHms[a],
							allDltFrontHms[b],
							allDltFrontHms[c],
							allDltFrontHms[d],
							allDltFrontHms[e],
						}

						slices.Sort(front)
						frontStr := strings.Join(front, ",")

						// 计算前区和值
						frontHz := CalDltHz(front)
						if frontHz < (xuHaoSt.HzMin - 23) {
							continue LabelForE
						}

						danMaLen := len(xuHaoSt.FrontDanMaHms)

						// 判断是否存在指定胆码
						if danMaLen != 0 {
							if len(gen.SliceIntersection(front, xuHaoSt.FrontDanMaHms)) != danMaLen {
								continue LabelForE
							}
						}

						// 判断前区是否包括必须的号码
						if len(xuHaoSt.FrontIncludeHms) > 0 && len(gen.SliceIntersection(front, xuHaoSt.FrontIncludeHms)) == 0 {
							continue LabelForE
						}

						// 判断前区是否包括必须排除的号码
						if len(xuHaoSt.FrontExcludeHms) > 0 && len(gen.SliceIntersection(front, xuHaoSt.FrontExcludeHms)) != 0 {
							continue LabelForE
						}

						// 判断前区前中后是否在选定的值中
						if len(xuHaoSt.QzhSlice) > 0 && !slices.Contains(xuHaoSt.QzhSlice, CalDltQzh(front)) {
							continue LabelForE
						}

						// 判断是否是各种类型 特别是 OtherT
						if slices.Contains(xuHaoSt.Tx, "OtherT") {
							isT1, isT2, isT3, isT4, isT5, isT6, isT7, isT8, isT9, isT10, isT11, isT12, isT13, isT14, isT15, isT16, isT17, isOtherT := false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false
							// 是否在当前选中的非OtherT中
							isExistSelected1Dao17T := false

						LabelForJudge01:
							for _, tx := range xuHaoSt.Tx {
								if tx != "OtherT" {
									if slices.Contains(t2MoniABCDEs[tx]["A"], front[0]) &&
										slices.Contains(t2MoniABCDEs[tx]["B"], front[1]) &&
										slices.Contains(t2MoniABCDEs[tx]["C"], front[2]) &&
										slices.Contains(t2MoniABCDEs[tx]["D"], front[3]) &&
										slices.Contains(t2MoniABCDEs[tx]["E"], front[4]) {
										isExistSelected1Dao17T = true
										break LabelForJudge01
									}
								}
							}

							if !isExistSelected1Dao17T {
								for tx, moniABCDEms := range t2MoniABCDEs {
									if slices.Contains(moniABCDEms["A"], front[0]) &&
										slices.Contains(moniABCDEms["B"], front[1]) &&
										slices.Contains(moniABCDEms["C"], front[2]) &&
										slices.Contains(moniABCDEms["D"], front[3]) &&
										slices.Contains(moniABCDEms["E"], front[4]) {
										switch tx {
										case "T1":
											isT1 = true
										case "T2":
											isT2 = true
										case "T3":
											isT3 = true
										case "T4":
											isT4 = true
										case "T5":
											isT5 = true
										case "T6":
											isT6 = true
										case "T7":
											isT7 = true
										case "T8":
											isT8 = true
										case "T9":
											isT9 = true
										case "T10":
											isT10 = true
										case "T11":
											isT11 = true
										case "T12":
											isT12 = true
										case "T13":
											isT13 = true
										case "T14":
											isT14 = true
										case "T15":
											isT15 = true
										case "T16":
											isT16 = true
										case "T17":
											isT17 = true
										}
									}
								}
								if !(isT1 || isT2 || isT3 || isT4 || isT5 || isT6 || isT7 || isT8 || isT9 || isT10 || isT11 || isT12 || isT13 || isT14 || isT15 || isT16 || isT17) {
									isOtherT = true
								}
							}

						LabelForJudge03:
							for _, tx := range xuHaoSt.Tx {
								switch tx {
								case "T1":
									if isT1 {
										break LabelForJudge03
									}
								case "T2":
									if isT2 {
										break LabelForJudge03
									}
								case "T3":
									if isT3 {
										break LabelForJudge03
									}
								case "T4":
									if isT4 {
										break LabelForJudge03
									}
								case "T5":
									if isT5 {
										break LabelForJudge03
									}
								case "T6":
									if isT6 {
										break LabelForJudge03
									}
								case "T7":
									if isT7 {
										break LabelForJudge03
									}
								case "T8":
									if isT8 {
										break LabelForJudge03
									}
								case "T9":
									if isT9 {
										break LabelForJudge03
									}
								case "T10":
									if isT10 {
										break LabelForJudge03
									}
								case "T11":
									if isT11 {
										break LabelForJudge03
									}
								case "T12":
									if isT12 {
										break LabelForJudge03
									}
								case "T13":
									if isT13 {
										break LabelForJudge03
									}
								case "T14":
									if isT14 {
										break LabelForJudge03
									}
								case "T15":
									if isT15 {
										break LabelForJudge03
									}
								case "T16":
									if isT16 {
										break LabelForJudge03
									}
								case "T17":
									if isT17 {
										break LabelForJudge03
									}
								case "OtherT":
									if isOtherT {
										break LabelForJudge03
									}
								}
							}
						} else { // 选择的组合不包含OtherT的情况
						LabelForJudge02:
							for _, tx := range xuHaoSt.Tx {
								if slices.Contains(t2MoniABCDEs[tx]["A"], front[0]) &&
									slices.Contains(t2MoniABCDEs[tx]["B"], front[1]) &&
									slices.Contains(t2MoniABCDEs[tx]["C"], front[2]) &&
									slices.Contains(t2MoniABCDEs[tx]["D"], front[3]) &&
									slices.Contains(t2MoniABCDEs[tx]["E"], front[4]) {
									// 跳出当前循环,继续后区的生成
									break LabelForJudge02
								}
							}
						}

						// 后区组合生成（2个数字从12个中选取）
						for x := 0; x <= 10; x++ {
						LabelForY:
							for y := x + 1; y <= 11; y++ {
								back := []string{allDltBackHms[x], allDltBackHms[y]}
								slices.Sort(back)

								backStr := strings.Join(back, ",")

								// 判断后区是否包含必须存在的号码
								if len(xuHaoSt.BackIncludeHms) > 0 && len(gen.SliceIntersection(xuHaoSt.BackIncludeHms, back)) == 0 {
									continue LabelForY
								}

								// 判断后区是否包括必须排除的号码
								if len(xuHaoSt.BackExcludeHms) > 0 && len(gen.SliceIntersection(xuHaoSt.BackExcludeHms, back)) != 0 {
									continue LabelForY
								}

								// 判断后区是否是必须包含的组合
								if len(xuHaoSt.BackIncludeCombs) > 0 && !slices.Contains(xuHaoSt.BackIncludeCombs, backStr) {
									continue LabelForY
								}

								// 判断后区是否是必须排除的组合
								if len(xuHaoSt.BackExcludeCombs) > 0 && slices.Contains(xuHaoSt.BackExcludeCombs, backStr) {
									continue LabelForY
								}

								// 判断是否需要移除历史
								fullStr := frontStr + "|" + backStr
								if xuHaoSt.RemoveHis == 1 && slices.Contains(hisFullHmStrSlice, fullStr) {
									continue LabelForY
								}

								// 计算奇偶
								full := append(front, back...)
								oe := CalDltOe(full)
								// 判断奇偶是否在选定的范围内
								if !slices.Contains(xuHaoSt.Oes, oe) {
									continue LabelForY
								}

								// 计算和值
								totalHz := frontHz + CalDltHz(back)
								// 判断和值是否在选定的范围内
								if totalHz < xuHaoSt.HzMin || totalHz > xuHaoSt.HzMax {
									continue LabelForY
								}

								// 计算新4重号
								if xuHaoSt.Ch4MustGetCount > 0 {
									ic4s := make([]string, 0)
									// 从开奖号码中生成4个组合号码
									ic4s = gen.CrossComb(front, back, 4, 0)
									ic4s = append(ic4s, gen.CrossComb(front, back, 3, 1)...)
									ic4s = append(ic4s, gen.CrossComb(front, back, 2, 2)...)

									addCh4Count := 0

									// 严格按 len(v) > 1 的规则计算增量
									for _, ic4 := range ic4s {
										n := len(c4ms[ic4])
										if n == 1 {
											addCh4Count += 2
										} else if n >= 2 {
											addCh4Count += 1
										}
									}

									if addCh4Count < xuHaoSt.Ch4MustGetCount {
										continue LabelForY
									}
								}

								if finalCombCount < finalCombMaxCount {
									finalCombForRandSlice = append(finalCombForRandSlice, fullStr)
								} else {
									//goto LabelForReturn
									// 超过 finalCombMaxCount 则随机存放
									finalCombForRandSlice[rand.IntN(finalCombMaxCount)] = fullStr
								}
								finalCombCount++
							}
						}
					}
				}
			}
		}
	}

	//LabelForReturn:
	finalSelectFullHmStrSlice = make([]string, 0, zhuShu)
	tempNumSlice := make([]int, 0, zhuShu)
	if len(finalCombForRandSlice) > 0 {
		if len(finalCombForRandSlice) <= zhuShu {
			finalSelectFullHmStrSlice = append(finalSelectFullHmStrSlice, finalCombForRandSlice...)
		} else {
			for i := 0; i < zhuShu; i++ {
			LabelForContinue:
				time.Sleep(time.Duration(rand.IntN(1)) * time.Nanosecond)
				tempNum := rand.IntN(len(finalCombForRandSlice))
				if !slices.Contains(tempNumSlice, tempNum) {
					tempNumSlice = append(tempNumSlice, tempNum)
					finalSelectFullHmStrSlice = append(finalSelectFullHmStrSlice, finalCombForRandSlice[tempNum])
				} else {
					if len(tempNumSlice) < zhuShu {
						goto LabelForContinue
					}
				}
			}
		}
	}

	return
}

// ShuffleAllDltFrontHms 乱序大乐透前区的35个号码
//
//	@Description:
//	@return allDltFrontHms
func ShuffleAllDltFrontHms() (allDltFrontHms []string) {
	allDltFrontHms = gen.AllDltFrontHms
	rand.Shuffle(len(allDltFrontHms), func(i, j int) {
		allDltFrontHms[i], allDltFrontHms[j] = allDltFrontHms[j], allDltFrontHms[i]
	})
	return
}

func ShuffleAllDltBackHms() (allDltBackHms []string) {
	allDltBackHms = gen.AllDltBackHms
	rand.Shuffle(len(allDltBackHms), func(i, j int) {
		allDltBackHms[i], allDltBackHms[j] = allDltBackHms[j], allDltBackHms[i]
	})
	return
}
