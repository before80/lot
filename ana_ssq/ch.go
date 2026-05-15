package ana_ssq

import (
	"fmt"
	"slices"
	"sort"

	"github.com/before80/lot/dbop"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/models"
	"github.com/before80/lot/opdf"
	"github.com/jung-kurt/gofpdf"
)

type SsqChSt struct {
	Typ      string
	Typ1     string
	Typ2     string
	Cs       int
	SsqInfos []string // 期号->7个开奖号码
}

// SsqCHongHao 重号分析
func SsqCHongHao() []SsqChSt {
	c7ms := make(map[string][]models.Ssq)
	c6ms := make(map[string][]models.Ssq)
	c5ms := make(map[string][]models.Ssq)
	c4ms := make(map[string][]models.Ssq)
	var c7s, c6s, c5s, c4s []string

	for _, ssq := range ZxSsqs {
		var ic7s, ic6s, ic5s, ic4s []string
		ic7s, ic6s, ic5s, ic4s = nil, nil, nil, nil
		// 从开奖号码中生成7个组合号码
		ic7s = append(ic6s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 6, 1)...)
		// 从开奖号码中生成6个组合号码
		ic6s = append(ic6s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 6, 0)...)
		ic6s = append(ic6s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 5, 1)...)
		// 从开奖号码中生成5个组合号码
		ic5s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 5, 0)
		ic5s = append(ic5s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 4, 1)...)
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 3, 1)...)

		for _, c7 := range ic7s {
			c7ms[c7] = append(c7ms[c7], ssq)
		}
		for _, c6 := range ic6s {
			c6ms[c6] = append(c6ms[c6], ssq)
		}
		for _, c5 := range ic5s {
			c5ms[c5] = append(c5ms[c5], ssq)
		}
		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], ssq)
		}

		c7s = append(c7s, ic7s...)
		c6s = append(c6s, ic6s...)
		c5s = append(c5s, ic5s...)
		c4s = append(c4s, ic4s...)
	}

	c7KLens := make([]KeyWithLength, 0, len(c7ms))
	c6KLens := make([]KeyWithLength, 0, len(c6ms))
	c5KLens := make([]KeyWithLength, 0, len(c5ms))
	c4KLens := make([]KeyWithLength, 0, len(c4ms))

	for qhComb, v := range c7ms {
		c7KLens = append(c7KLens, KeyWithLength{
			Key:    qhComb,
			Length: len(v),
		})
	}

	for qhComb, v := range c6ms {
		c6KLens = append(c6KLens, KeyWithLength{
			Key:    qhComb,
			Length: len(v),
		})
	}

	for qhComb, v := range c5ms {
		c5KLens = append(c5KLens, KeyWithLength{
			Key:    qhComb,
			Length: len(v),
		})
	}

	for qhComb, v := range c4ms {
		c4KLens = append(c4KLens, KeyWithLength{
			Key:    qhComb,
			Length: len(v),
		})
	}

	sort.Slice(c7KLens, func(i, j int) bool {
		return c7KLens[i].Length > c7KLens[j].Length
	})
	sort.Slice(c6KLens, func(i, j int) bool {
		return c6KLens[i].Length > c6KLens[j].Length
	})
	sort.Slice(c5KLens, func(i, j int) bool {
		return c5KLens[i].Length > c5KLens[j].Length
	})
	sort.Slice(c4KLens, func(i, j int) bool {
		return c4KLens[i].Length > c4KLens[j].Length
	})

	var ssqChs, ssqCh7s, ssqCh6s, ssqCh5s, ssqCh4s []SsqChSt

	var c7KHave2DrawNums, c6KHave2DrawNums, c5KHave2DrawNums, c4KHave2DrawNums []string
	for _, v := range c7KLens {
		typ := v.Key
		typ1 := "7重号"
		typ2 := "5+2"
		issqs := c7ms[v.Key]
		if len(issqs) > 1 {
			sort.Slice(issqs, func(i, j int) bool {
				return issqs[i].DrawNum > issqs[j].DrawNum
			})
			for _, issq := range issqs {
				c7KHave2DrawNums = append(c7KHave2DrawNums, issq.DrawNum)
			}
			ssqInfos := make([]string, 0, len(issqs))
			for _, issq := range issqs {
				ssqInfos = append(ssqInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s,%s|%s", issq.DrawNum, issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.F6, issq.B1))
			}
			ssqCh7s = append(ssqCh7s, SsqChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(issqs),
				SsqInfos: ssqInfos,
			})
		}
	}

	sort.Slice(ssqCh7s, func(i, j int) bool { return ssqCh7s[i].Cs > ssqCh7s[j].Cs })
	ssqChs = append(ssqChs, ssqCh7s...)

	for _, v := range c6KLens {
		typ := v.Key
		typ1 := "6重号"
		fN, bN := gen.JudgeFrontBackCountFromStr(typ)
		typ2 := fmt.Sprintf("%d+%d", fN, bN)

		issqs := c6ms[v.Key]
		if len(issqs) > 1 {
			sort.Slice(issqs, func(i, j int) bool {
				return issqs[i].DrawNum > issqs[j].DrawNum
			})
			for _, issq := range issqs {
				if !slices.Contains(c7KHave2DrawNums, issq.DrawNum) {
					c6KHave2DrawNums = append(c6KHave2DrawNums, issq.DrawNum)
				}
			}
			if len(issqs) > 2 {
				ssqInfos := make([]string, 0, len(issqs))
				for _, issq := range issqs {
					ssqInfos = append(ssqInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s,%s|%s", issq.DrawNum, issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.F6, issq.B1))
				}
				ssqCh6s = append(ssqCh6s, SsqChSt{
					Typ:      typ,
					Typ1:     typ1,
					Typ2:     typ2,
					Cs:       len(issqs),
					SsqInfos: ssqInfos,
				})
			}
		}
	}
	sort.Slice(ssqCh6s, func(i, j int) bool { return ssqCh6s[i].Cs > ssqCh6s[j].Cs })
	ssqChs = append(ssqChs, ssqCh6s...)

	for _, v := range c5KLens {
		typ := v.Key
		typ1 := "5重号"
		fN, bN := gen.JudgeFrontBackCountFromStr(typ)
		typ2 := fmt.Sprintf("%d+%d", fN, bN)
		issqs := c5ms[v.Key]
		if len(issqs) > 1 {
			sort.Slice(issqs, func(i, j int) bool {
				return issqs[i].DrawNum > issqs[j].DrawNum
			})
			for _, issq := range issqs {
				if !slices.Contains(c6KHave2DrawNums, issq.DrawNum) && !slices.Contains(c7KHave2DrawNums, issq.DrawNum) {
					c5KHave2DrawNums = append(c5KHave2DrawNums, issq.DrawNum)
				}
			}
			ssqInfos := make([]string, 0, len(issqs))
			for _, issq := range issqs {
				ssqInfos = append(ssqInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s,%s|%s", issq.DrawNum, issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.F6, issq.B1))
			}
			ssqCh5s = append(ssqCh5s, SsqChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(issqs),
				SsqInfos: ssqInfos,
			})
		}
	}

	sort.Slice(ssqCh5s, func(i, j int) bool { return ssqCh5s[i].Cs > ssqCh5s[j].Cs })
	ssqChs = append(ssqChs, ssqCh5s...)
	for _, v := range c4KLens {
		typ := v.Key
		typ1 := "4重号"
		fN, bN := gen.JudgeFrontBackCountFromStr(typ)
		typ2 := fmt.Sprintf("%d+%d", fN, bN)
		issqs := c4ms[v.Key]

		if len(issqs) > 1 {
			sort.Slice(issqs, func(i, j int) bool {
				return issqs[i].DrawNum > issqs[j].DrawNum
			})
			for _, issq := range issqs {
				if !slices.Contains(c5KHave2DrawNums, issq.DrawNum) && !slices.Contains(c6KHave2DrawNums, issq.DrawNum) && !slices.Contains(c7KHave2DrawNums, issq.DrawNum) {
					c4KHave2DrawNums = append(c4KHave2DrawNums, issq.DrawNum)
				}
			}
			ssqInfos := make([]string, 0, len(issqs))
			for _, issq := range issqs {
				ssqInfos = append(ssqInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s,%s|%s", issq.DrawNum, issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.F6, issq.B1))
			}
			ssqCh4s = append(ssqCh4s, SsqChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(issqs),
				SsqInfos: ssqInfos,
			})
		}
	}
	sort.Slice(ssqCh4s, func(i, j int) bool { return ssqCh4s[i].Cs > ssqCh4s[j].Cs })
	ssqChs = append(ssqChs, ssqCh4s...)
	return ssqChs
}

func CHongHaoForPDF(pdf *gofpdf.Fpdf, fontName string, chapterFontSize float64, chapter string, level int, limitLen65432 [5]int, autoFilterMore bool) {
	pdf.SetFont(fontName, "B", chapterFontSize)
	pdf.Cell(0, opdf.CellHeight(chapterFontSize), chapter)
	pdf.Bookmark(chapter, level, 0)
	pdf.Ln(opdf.LineHeight(chapterFontSize))
	tableHeader := []string{"期号", "日期", "号码"}
	colWidth := []float64{10, 16, 26}
	colHeight := opdf.LineHeight(chapterFontSize - 2)

	var dnf = func(rows [][]string) {
		pdf.SetFont(fontName, "L", chapterFontSize-3)
		// --- 表头 ---
		for i, str := range tableHeader {
			pdf.CellFormat(colWidth[i], colHeight, str, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)

		// --- 表内容 ---
		for _, row := range rows {
			for i, txt := range row {
				pdf.CellFormat(colWidth[i], colHeight, txt, "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
		}
		pdf.Ln(2)
	}

	ssqs, _ := dbop.ReadAllSsq(false)
	// key为前后区字符串组合，value为models.Ssq切片
	c7ms := make(map[string][]models.Ssq)
	c6ms := make(map[string][]models.Ssq)
	c5ms := make(map[string][]models.Ssq)
	c4ms := make(map[string][]models.Ssq)
	c3ms := make(map[string][]models.Ssq)
	c2ms := make(map[string][]models.Ssq)
	// 存放所有开奖结果前后区各种的字符串组合
	var c7s, c6s, c5s, c4s, c3s, c2s []string

	for _, ssq := range ssqs {
		// 存放当前这注开奖结果前后区各种的字符串组合
		var ic7s, ic6s, ic5s, ic4s, ic3s, ic2s []string
		ic7s, ic6s, ic5s, ic4s, ic3s, ic2s = nil, nil, nil, nil, nil, nil
		// 从开奖号码中生成7个组合号码
		ic7s = append(ic7s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 5, 2)...)
		// 从开奖号码中生成6个组合号码
		ic6s = append(ic6s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 5, 1)...)
		ic6s = append(ic6s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 4, 2)...)
		// 从开奖号码中生成5个组合号码
		ic5s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 5, 0)
		ic5s = append(ic5s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 4, 1)...)
		ic5s = append(ic5s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 3, 2)...)
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 3, 1)...)
		ic4s = append(ic4s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 2, 2)...)
		// 从开奖号码中生成3个组合号码
		ic3s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 3, 0)
		ic3s = append(ic3s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 2, 1)...)
		ic3s = append(ic3s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 1, 2)...)
		// 从开奖号码中生成2个组合号码
		ic2s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 2, 0)
		ic2s = append(ic2s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 1, 1)...)
		ic2s = append(ic2s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 0, 2)...)
		for _, c7 := range ic7s {
			c7ms[c7] = append(c7ms[c7], ssq)
		}
		for _, c6 := range ic6s {
			c6ms[c6] = append(c6ms[c6], ssq)
		}
		for _, c5 := range ic5s {
			c5ms[c5] = append(c5ms[c5], ssq)
		}
		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], ssq)
		}
		for _, c3 := range ic3s {
			c3ms[c3] = append(c3ms[c3], ssq)
		}
		for _, c2 := range ic2s {
			c2ms[c2] = append(c2ms[c2], ssq)
		}
		c7s = append(c7s, ic7s...)
		c6s = append(c6s, ic6s...)
		c5s = append(c5s, ic5s...)
		c4s = append(c4s, ic4s...)
		c3s = append(c3s, ic3s...)
		c2s = append(c2s, ic2s...)
	}

	c7KLens := make([]KeyWithLength, 0, len(c7ms))
	c6KLens := make([]KeyWithLength, 0, len(c6ms))
	c5KLens := make([]KeyWithLength, 0, len(c5ms))
	c4KLens := make([]KeyWithLength, 0, len(c4ms))
	c3KLens := make([]KeyWithLength, 0, len(c3ms))
	c2KLens := make([]KeyWithLength, 0, len(c2ms))

	for qhComb, v := range c7ms {
		c7KLens = append(c7KLens, KeyWithLength{
			Key:    qhComb,
			Length: len(v),
		})
	}

	for qhComb, v := range c6ms {
		c6KLens = append(c6KLens, KeyWithLength{
			Key:    qhComb,
			Length: len(v),
		})
	}

	for qhComb, v := range c5ms {
		c5KLens = append(c5KLens, KeyWithLength{
			Key:    qhComb,
			Length: len(v),
		})
	}

	for qhComb, v := range c4ms {
		c4KLens = append(c4KLens, KeyWithLength{
			Key:    qhComb,
			Length: len(v),
		})
	}

	for qhComb, v := range c3ms {
		c3KLens = append(c3KLens, KeyWithLength{
			Key:    qhComb,
			Length: len(v),
		})
	}

	for qhComb, v := range c2ms {
		c2KLens = append(c2KLens, KeyWithLength{
			Key:    qhComb,
			Length: len(v),
		})
	}

	sort.Slice(c7KLens, func(i, j int) bool {
		return c7KLens[i].Length > c7KLens[j].Length
	})
	sort.Slice(c6KLens, func(i, j int) bool {
		return c6KLens[i].Length > c6KLens[j].Length
	})
	sort.Slice(c5KLens, func(i, j int) bool {
		return c5KLens[i].Length > c5KLens[j].Length
	})
	sort.Slice(c4KLens, func(i, j int) bool {
		return c4KLens[i].Length > c4KLens[j].Length
	})
	sort.Slice(c3KLens, func(i, j int) bool {
		return c3KLens[i].Length > c3KLens[j].Length
	})
	sort.Slice(c2KLens, func(i, j int) bool {
		return c2KLens[i].Length > c2KLens[j].Length
	})

	var c7KHave2DrawNums, c6KHave2DrawNums, c5KHave2DrawNums, c4KHave2DrawNums, c3KHave2DrawNums []string
	_ = c3KHave2DrawNums
	tempNum := 0
	for _, v := range c7KLens {
		issqs := c7ms[v.Key]
		if len(issqs) > 1 {
			sort.Slice(issqs, func(i, j int) bool {
				return issqs[i].DrawNum > issqs[j].DrawNum
			})
			var iissqs []models.Ssq
			for _, issq := range issqs {
				c7KHave2DrawNums = append(c7KHave2DrawNums, issq.DrawNum)
				iissqs = append(iissqs, issq)
			}
			if len(iissqs) > 0 {
				if tempNum == 0 {
					tempNum++
					str := fmt.Sprintf("重复7个号码")
					pdf.SetFont(fontName, "B", chapterFontSize-1)
					pdf.Cell(0, opdf.CellHeight(chapterFontSize-1), str)
					pdf.Bookmark(str, level+1, 0)
					pdf.Ln(opdf.LineHeight(chapterFontSize - 1)) // 换行
				}
				str := fmt.Sprintf("%s -> 已出现%d期", v.Key, v.Length)
				pdf.SetFont(fontName, "B", chapterFontSize-2)
				pdf.Cell(0, opdf.CellHeight(chapterFontSize-2), str)
				//pdf.Bookmark(str, level+2, 0)
				pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
				pdf.SetFont(fontName, "L", chapterFontSize-3)
				var rows [][]string
				for _, issq := range iissqs {
					rows = append(rows, []string{issq.DrawNum, issq.DrawTime, fmt.Sprintf("%s %s %s %s %s %s | %s", issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.F6, issq.B1)})
				}
				dnf(rows)
				//for _, issq := range iissqs {
				//	dnf([][]string{})
				//	pdf.Cell(0, opdf.CellHeight(chapterFontSize-3), fmt.Sprintf("%s %s -> %s %s %s %s %s | %s %s", issq.DrawNum, issq.DrawTime, issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.B1, issq.B2))
				//	pdf.Ln(opdf.LineHeight(chapterFontSize - 3)) // 换行
				//}
			}
		}

	}
	tempNum = 0
	for k, v := range c6KLens {
		issqs := c6ms[v.Key]
		if limitLen65432[0] != 0 && k > limitLen65432[0] {
			break
		}
		//if k > 10 {
		//	break
		//}
		if len(issqs) > 1 {
			sort.Slice(issqs, func(i, j int) bool {
				return issqs[i].DrawNum > issqs[j].DrawNum
			})
			var iissqs []models.Ssq
			for _, issq := range issqs {
				if autoFilterMore {
					if !slices.Contains(c7KHave2DrawNums, issq.DrawNum) {
						c6KHave2DrawNums = append(c6KHave2DrawNums, issq.DrawNum)
						iissqs = append(iissqs, issq)
					}
				} else {
					iissqs = append(iissqs, issq)
				}

			}
			if len(iissqs) > 0 {
				if tempNum == 0 {
					tempNum++
					str := fmt.Sprintf("重复6个号码")
					pdf.SetFont(fontName, "B", chapterFontSize-1)
					pdf.Cell(0, opdf.CellHeight(chapterFontSize-1), str)
					pdf.Bookmark(str, level+1, 0)
					pdf.Ln(opdf.LineHeight(chapterFontSize - 1)) // 换行
				}
				str := fmt.Sprintf("%s -> 已出现%d期", v.Key, v.Length)
				pdf.SetFont(fontName, "B", chapterFontSize-2)
				pdf.Cell(0, opdf.CellHeight(chapterFontSize-2), str)
				//pdf.Bookmark(str, level+2, 0)
				pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
				pdf.SetFont(fontName, "L", chapterFontSize-3)
				//for _, issq := range iissqs {
				//	pdf.Cell(0, opdf.CellHeight(chapterFontSize-3), fmt.Sprintf("%s %s -> %s %s %s %s %s | %s %s", issq.DrawNum, issq.DrawTime, issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.B1, issq.B2))
				//	pdf.Ln(opdf.LineHeight(chapterFontSize - 3)) // 换行
				//}
				var rows [][]string
				for _, issq := range iissqs {
					rows = append(rows, []string{issq.DrawNum, issq.DrawTime, fmt.Sprintf("%s %s %s %s %s %s | %s", issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.F6, issq.B1)})
				}
				dnf(rows)
			}
		}
	}
	tempNum = 0
	for k, v := range c5KLens {
		issqs := c5ms[v.Key]
		if limitLen65432[1] != 0 && k > limitLen65432[1] {
			break
		}
		//if k > 10 {
		//	break
		//}
		if len(issqs) > 1 {
			sort.Slice(issqs, func(i, j int) bool {
				return issqs[i].DrawNum > issqs[j].DrawNum
			})
			var iissqs []models.Ssq
			for _, issq := range issqs {
				if autoFilterMore {
					if !slices.Contains(c6KHave2DrawNums, issq.DrawNum) && !slices.Contains(c7KHave2DrawNums, issq.DrawNum) {
						c5KHave2DrawNums = append(c5KHave2DrawNums, issq.DrawNum)
						iissqs = append(iissqs, issq)
					}
				} else {
					iissqs = append(iissqs, issq)
				}

			}
			if len(iissqs) > 0 {
				if tempNum == 0 {
					tempNum++
					str := fmt.Sprintf("重复5个号码")
					pdf.SetFont(fontName, "B", chapterFontSize-1)
					pdf.Cell(0, opdf.CellHeight(chapterFontSize-1), str)
					pdf.Bookmark(str, level+1, 0)
					pdf.Ln(opdf.LineHeight(chapterFontSize - 1)) // 换行
				}
				str := fmt.Sprintf("%s -> 已出现%d期", v.Key, v.Length)
				pdf.SetFont(fontName, "B", chapterFontSize-2)
				pdf.Cell(0, opdf.CellHeight(chapterFontSize-2), str)
				//pdf.Bookmark(str, level+2, 0)
				pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
				pdf.SetFont(fontName, "L", chapterFontSize-3)
				//for _, issq := range iissqs {
				//	pdf.Cell(0, opdf.CellHeight(chapterFontSize-3), fmt.Sprintf("%s %s -> %s %s %s %s %s | %s %s", issq.DrawNum, issq.DrawTime, issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.B1, issq.B2))
				//	pdf.Ln(opdf.LineHeight(chapterFontSize - 3)) // 换行
				//}
				var rows [][]string
				for _, issq := range iissqs {
					rows = append(rows, []string{issq.DrawNum, issq.DrawTime, fmt.Sprintf("%s %s %s %s %s %s | %s", issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.F6, issq.B1)})
				}
				dnf(rows)
			}
		}
	}
	tempNum = 0
	for k, v := range c4KLens {
		issqs := c4ms[v.Key]
		if limitLen65432[2] != 0 && k > limitLen65432[2] {
			break
		}
		//if k > 10 {
		//	break
		//}
		if len(issqs) > 1 {
			sort.Slice(issqs, func(i, j int) bool {
				return issqs[i].DrawNum > issqs[j].DrawNum
			})
			var iissqs []models.Ssq
			for _, issq := range issqs {
				if autoFilterMore {
					if !slices.Contains(c5KHave2DrawNums, issq.DrawNum) && !slices.Contains(c6KHave2DrawNums, issq.DrawNum) && !slices.Contains(c7KHave2DrawNums, issq.DrawNum) {
						c4KHave2DrawNums = append(c4KHave2DrawNums, issq.DrawNum)
						iissqs = append(iissqs, issq)
					}
				} else {
					iissqs = append(iissqs, issq)
				}
			}
			if len(iissqs) > 0 {
				if tempNum == 0 {
					tempNum++
					str := fmt.Sprintf("重复4个号码")
					pdf.SetFont(fontName, "B", chapterFontSize-1)
					pdf.Cell(0, opdf.CellHeight(chapterFontSize-1), str)
					pdf.Bookmark(str, level+1, 0)
					pdf.Ln(opdf.LineHeight(chapterFontSize - 1)) // 换行
				}
				str := fmt.Sprintf("%s -> 已出现%d期", v.Key, v.Length)
				pdf.SetFont(fontName, "B", chapterFontSize-2)
				pdf.Cell(0, opdf.CellHeight(chapterFontSize-2), str)
				//pdf.Bookmark(str, level+2, 0)
				pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
				pdf.SetFont(fontName, "L", chapterFontSize-3)
				//for _, issq := range iissqs {
				//	pdf.Cell(0, opdf.CellHeight(chapterFontSize-3), fmt.Sprintf("%s %s -> %s %s %s %s %s | %s %s", issq.DrawNum, issq.DrawTime, issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.B1, issq.B2))
				//	pdf.Ln(opdf.LineHeight(chapterFontSize - 3)) // 换行
				//}
				var rows [][]string
				for _, issq := range iissqs {
					rows = append(rows, []string{issq.DrawNum, issq.DrawTime, fmt.Sprintf("%s %s %s %s %s %s | %s", issq.F1, issq.F2, issq.F3, issq.F4, issq.F5, issq.F6, issq.B1)})
				}
				dnf(rows)
			}
		}
	}
}

type CHongHaoSt struct {
	DrawNum             string
	NewAddCh4           int // 当期新增4重号数
	NewAddCh5           int // 当期新增5重号数
	NewAddCh6           int // 当期新增6重号数
	NewAddCh7           int // 当期新增7重号数
	DangQiTotalNewAddCh int // 当期总的新增重号数
	LeiJiaCh            int // 累计重号数
}

func CalCHLj(cms map[string][]models.Ssq) (lj int) {
	for _, v := range cms {
		if len(v) > 1 {
			lj += len(v)
		}
	}
	return
}

func CalCHLjFromStrSlice(cms map[string][]string) (lj int) {
	for _, v := range cms {
		if len(v) > 1 {
			lj += len(v)
		}
	}
	return
}

// CHongHaoLeiJia 重号累加
//
//	@Description:
//	@return drawNum2CHongHaoSt
func CHongHaoLeiJia() (drawNum2CHongHaoSt map[string]CHongHaoSt) {
	ssqs, _ := dbop.ReadAllSsq(false)
	drawNum2CHongHaoSt = make(map[string]CHongHaoSt, len(ssqs)-1)

	c7ms := make(map[string][]models.Ssq)
	c6ms := make(map[string][]models.Ssq)
	c5ms := make(map[string][]models.Ssq)
	c4ms := make(map[string][]models.Ssq)

	var c7SLLj, c6SLLj, c5SLLj, c4SLLj int
	var c7SLLjPrev, c6SLLjPrev, c5SLLjPrev, c4SLLjPrev int

	for _, ssq := range ssqs {
		var ic7s, ic6s, ic5s, ic4s []string
		ic7s, ic6s, ic5s, ic4s = nil, nil, nil, nil
		// 从开奖号码中生成7个组合号码
		ic7s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 6, 1)
		// 从开奖号码中生成6个组合号码
		ic6s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 6, 0)
		ic6s = append(ic6s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 5, 1)...)
		// 从开奖号码中生成5个组合号码
		ic5s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 5, 0)
		ic5s = append(ic5s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 4, 1)...)
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{ssq.F1, ssq.F2, ssq.F3, ssq.F4, ssq.F5, ssq.F6}, []string{ssq.B1}, 3, 1)...)

		c7SLLjPrev = c7SLLj
		c6SLLjPrev = c6SLLj
		c5SLLjPrev = c5SLLj
		c4SLLjPrev = c4SLLj
		for _, c7 := range ic7s {
			c7ms[c7] = append(c7ms[c7], ssq)
		}
		c7SLLj = CalCHLj(c7ms)

		for _, c6 := range ic6s {
			c6ms[c6] = append(c6ms[c6], ssq)
		}
		c6SLLj = CalCHLj(c6ms)

		for _, c5 := range ic5s {
			c5ms[c5] = append(c5ms[c5], ssq)
		}
		c5SLLj = CalCHLj(c5ms)
		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], ssq)
		}
		c4SLLj = CalCHLj(c4ms)
		addCh4, addCh5, addCh6, addCh7 := c4SLLj-c4SLLjPrev, c5SLLj-c5SLLjPrev, c6SLLj-c6SLLjPrev, c7SLLj-c7SLLjPrev

		//fmt.Printf("DrawNum= %s c4SLLj=%d c5SLLj=%d c6SLLj=%d c7SLLj=%d addCh4=%d addCh5=%d addCh6=%d addCh7=%d LeiJiaCh=%d\n", ssq.DrawNum, c4SLLj, c5SLLj, c6SLLj, c7SLLj, addCh4, addCh5, addCh6, addCh7, c4SLLj+c5SLLj+c6SLLj+c7SLLj)

		drawNum2CHongHaoSt[ssq.DrawNum] = CHongHaoSt{
			DrawNum:             ssq.DrawNum,
			NewAddCh4:           addCh4,
			NewAddCh5:           addCh5,
			NewAddCh6:           addCh6,
			NewAddCh7:           addCh7,
			DangQiTotalNewAddCh: addCh4 + addCh5 + addCh6 + addCh7,
			LeiJiaCh:            c4SLLj + c5SLLj + c6SLLj + c7SLLj,
		}
	}
	return
}
