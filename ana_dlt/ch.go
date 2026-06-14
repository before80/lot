package ana_dlt

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/before80/lot/dbop"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/models"
	"github.com/before80/lot/opdf"
	"github.com/jung-kurt/gofpdf"
)

type DltChSt struct {
	Typ      string
	Typ1     string
	Typ2     string
	Cs       int
	DltInfos []string // 期号->7个开奖号码
}

// DltCHongHao 重号分析
func DltCHongHao() []DltChSt {
	c7ms := make(map[string][]models.Dlt)
	c6ms := make(map[string][]models.Dlt)
	c5ms := make(map[string][]models.Dlt)
	c4ms := make(map[string][]models.Dlt)
	var c7s, c6s, c5s, c4s []string

	for _, dlt := range ZxDlts {
		var ic7s, ic6s, ic5s, ic4s []string
		ic7s, ic6s, ic5s, ic4s = nil, nil, nil, nil
		// 从开奖号码中生成7个组合号码
		ic7s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 2)
		// 从开奖号码中生成6个组合号码
		ic6s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 1)
		ic6s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 2)...)
		// 从开奖号码中生成5个组合号码
		ic5s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 0)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 1)...)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 2)...)
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 1)...)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 2, 2)...)

		for _, c7 := range ic7s {
			c7ms[c7] = append(c7ms[c7], dlt)
		}
		for _, c6 := range ic6s {
			c6ms[c6] = append(c6ms[c6], dlt)
		}
		for _, c5 := range ic5s {
			c5ms[c5] = append(c5ms[c5], dlt)
		}
		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], dlt)
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

	var dltChs, dltCh7s, dltCh6s, dltCh5s, dltCh4s []DltChSt

	var c7KHave2DrawNums, c6KHave2DrawNums, c5KHave2DrawNums, c4KHave2DrawNums []string
	for _, v := range c7KLens {
		typ := v.Key
		typ1 := "7重号"
		typ2 := "5+2"
		idlts := c7ms[v.Key]
		if len(idlts) > 1 {
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			for _, idlt := range idlts {
				c7KHave2DrawNums = append(c7KHave2DrawNums, idlt.DrawNum)
			}
			dltInfos := make([]string, 0, len(idlts))
			for _, idlt := range idlts {
				dltInfos = append(dltInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s|%s,%s", idlt.DrawNum, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
			}
			dltCh7s = append(dltCh7s, DltChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(idlts),
				DltInfos: dltInfos,
			})
		}
	}

	sort.Slice(dltCh7s, func(i, j int) bool { return dltCh7s[i].Cs > dltCh7s[j].Cs })
	dltChs = append(dltChs, dltCh7s...)

	for _, v := range c6KLens {
		typ := v.Key
		typ1 := "6重号"
		fN, bN := gen.JudgeFrontBackCountFromStr(typ)
		typ2 := fmt.Sprintf("%d+%d", fN, bN)

		idlts := c6ms[v.Key]
		if len(idlts) > 1 {
			//fmt.Printf("--> %v\n", idlts)
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			for _, idlt := range idlts {
				if !slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
					c6KHave2DrawNums = append(c6KHave2DrawNums, idlt.DrawNum)
				}
			}
			dltInfos := make([]string, 0, len(idlts))
			for _, idlt := range idlts {
				dltInfos = append(dltInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s|%s,%s", idlt.DrawNum, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
			}
			dltCh6s = append(dltCh6s, DltChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(idlts),
				DltInfos: dltInfos,
			})
		}
	}
	sort.Slice(dltCh6s, func(i, j int) bool { return dltCh6s[i].Cs > dltCh6s[j].Cs })
	dltChs = append(dltChs, dltCh6s...)

	for _, v := range c5KLens {
		typ := v.Key
		typ1 := "5重号"
		fN, bN := gen.JudgeFrontBackCountFromStr(typ)
		typ2 := fmt.Sprintf("%d+%d", fN, bN)
		idlts := c5ms[v.Key]
		if len(idlts) > 1 {
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			for _, idlt := range idlts {
				if !slices.Contains(c6KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
					c5KHave2DrawNums = append(c5KHave2DrawNums, idlt.DrawNum)
				}
			}
			dltInfos := make([]string, 0, len(idlts))
			for _, idlt := range idlts {
				dltInfos = append(dltInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s|%s,%s", idlt.DrawNum, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
			}
			dltCh5s = append(dltCh5s, DltChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(idlts),
				DltInfos: dltInfos,
			})
		}
	}

	sort.Slice(dltCh5s, func(i, j int) bool { return dltCh5s[i].Cs > dltCh5s[j].Cs })
	dltChs = append(dltChs, dltCh5s...)
	for _, v := range c4KLens {
		typ := v.Key
		typ1 := "4重号"
		fN, bN := gen.JudgeFrontBackCountFromStr(typ)
		typ2 := fmt.Sprintf("%d+%d", fN, bN)
		idlts := c4ms[v.Key]

		if len(idlts) > 1 {
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			for _, idlt := range idlts {
				if !slices.Contains(c5KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c6KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
					c4KHave2DrawNums = append(c4KHave2DrawNums, idlt.DrawNum)
				}
			}
			dltInfos := make([]string, 0, len(idlts))
			for _, idlt := range idlts {
				dltInfos = append(dltInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s|%s,%s", idlt.DrawNum, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
			}
			dltCh4s = append(dltCh4s, DltChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(idlts),
				DltInfos: dltInfos,
			})
		}
	}
	sort.Slice(dltCh4s, func(i, j int) bool { return dltCh4s[i].Cs > dltCh4s[j].Cs })
	dltChs = append(dltChs, dltCh4s...)
	return dltChs
}

// DltEqCHongHao 设备的重号分析
func DltEqCHongHao(eqNumCount int) []DltChSt {
	c7ms := make(map[string][]models.Dlt)
	c6ms := make(map[string][]models.Dlt)
	c5ms := make(map[string][]models.Dlt)
	c4ms := make(map[string][]models.Dlt)
	var c7s, c6s, c5s, c4s []string

	for _, dlt := range ZxDlts {
		if dlt.DrawNum < "11001" {
			continue
		}
		if dlt.EquipmentCount != eqNumCount {
			continue
		}

		var ic7s, ic6s, ic5s, ic4s []string
		ic7s, ic6s, ic5s, ic4s = nil, nil, nil, nil
		// 从开奖号码中生成7个组合号码
		ic7s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 2)
		// 从开奖号码中生成6个组合号码
		ic6s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 1)
		ic6s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 2)...)
		// 从开奖号码中生成5个组合号码
		ic5s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 0)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 1)...)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 2)...)
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 1)...)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 2, 2)...)

		for _, c7 := range ic7s {
			c7ms[c7] = append(c7ms[c7], dlt)
		}
		for _, c6 := range ic6s {
			c6ms[c6] = append(c6ms[c6], dlt)
		}
		for _, c5 := range ic5s {
			c5ms[c5] = append(c5ms[c5], dlt)
		}
		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], dlt)
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

	var dltChs, dltCh7s, dltCh6s, dltCh5s, dltCh4s []DltChSt

	var c7KHave2DrawNums, c6KHave2DrawNums, c5KHave2DrawNums, c4KHave2DrawNums []string
	for _, v := range c7KLens {
		typ := v.Key
		typ1 := "7重号"
		typ2 := "5+2"
		idlts := c7ms[v.Key]
		if len(idlts) > 1 {
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			for _, idlt := range idlts {
				c7KHave2DrawNums = append(c7KHave2DrawNums, idlt.DrawNum)
			}
			dltInfos := make([]string, 0, len(idlts))
			for _, idlt := range idlts {
				dltInfos = append(dltInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s|%s,%s", idlt.DrawNum, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
			}
			dltCh7s = append(dltCh7s, DltChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(idlts),
				DltInfos: dltInfos,
			})
		}
	}

	sort.Slice(dltCh7s, func(i, j int) bool { return dltCh7s[i].Cs > dltCh7s[j].Cs })
	dltChs = append(dltChs, dltCh7s...)

	for _, v := range c6KLens {
		typ := v.Key
		typ1 := "6重号"
		fN, bN := gen.JudgeFrontBackCountFromStr(typ)
		typ2 := fmt.Sprintf("%d+%d", fN, bN)

		idlts := c6ms[v.Key]
		if len(idlts) > 1 {
			//fmt.Printf("--> %v\n", idlts)
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			for _, idlt := range idlts {
				if !slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
					c6KHave2DrawNums = append(c6KHave2DrawNums, idlt.DrawNum)
				}
			}
			dltInfos := make([]string, 0, len(idlts))
			for _, idlt := range idlts {
				dltInfos = append(dltInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s|%s,%s", idlt.DrawNum, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
			}
			dltCh6s = append(dltCh6s, DltChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(idlts),
				DltInfos: dltInfos,
			})
		}
	}
	sort.Slice(dltCh6s, func(i, j int) bool { return dltCh6s[i].Cs > dltCh6s[j].Cs })
	dltChs = append(dltChs, dltCh6s...)

	for _, v := range c5KLens {
		typ := v.Key
		typ1 := "5重号"
		fN, bN := gen.JudgeFrontBackCountFromStr(typ)
		typ2 := fmt.Sprintf("%d+%d", fN, bN)
		idlts := c5ms[v.Key]
		if len(idlts) > 1 {
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			for _, idlt := range idlts {
				if !slices.Contains(c6KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
					c5KHave2DrawNums = append(c5KHave2DrawNums, idlt.DrawNum)
				}
			}
			dltInfos := make([]string, 0, len(idlts))
			for _, idlt := range idlts {
				dltInfos = append(dltInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s|%s,%s", idlt.DrawNum, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
			}
			dltCh5s = append(dltCh5s, DltChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(idlts),
				DltInfos: dltInfos,
			})
		}
	}

	sort.Slice(dltCh5s, func(i, j int) bool { return dltCh5s[i].Cs > dltCh5s[j].Cs })
	dltChs = append(dltChs, dltCh5s...)
	for _, v := range c4KLens {
		typ := v.Key
		typ1 := "4重号"
		fN, bN := gen.JudgeFrontBackCountFromStr(typ)
		typ2 := fmt.Sprintf("%d+%d", fN, bN)
		idlts := c4ms[v.Key]

		if len(idlts) > 1 {
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			for _, idlt := range idlts {
				if !slices.Contains(c5KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c6KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
					c4KHave2DrawNums = append(c4KHave2DrawNums, idlt.DrawNum)
				}
			}
			dltInfos := make([]string, 0, len(idlts))
			for _, idlt := range idlts {
				dltInfos = append(dltInfos, fmt.Sprintf("%s->%s,%s,%s,%s,%s|%s,%s", idlt.DrawNum, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
			}
			dltCh4s = append(dltCh4s, DltChSt{
				Typ:      typ,
				Typ1:     typ1,
				Typ2:     typ2,
				Cs:       len(idlts),
				DltInfos: dltInfos,
			})
		}
	}
	sort.Slice(dltCh4s, func(i, j int) bool { return dltCh4s[i].Cs > dltCh4s[j].Cs })
	dltChs = append(dltChs, dltCh4s...)
	return dltChs
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

	dlts, _ := dbop.ReadAllDlt(false)
	// key为前后区字符串组合，value为models.Dlt切片
	c7ms := make(map[string][]models.Dlt)
	c6ms := make(map[string][]models.Dlt)
	c5ms := make(map[string][]models.Dlt)
	c4ms := make(map[string][]models.Dlt)
	c3ms := make(map[string][]models.Dlt)
	c2ms := make(map[string][]models.Dlt)
	// 存放所有开奖结果前后区各种的字符串组合
	var c7s, c6s, c5s, c4s, c3s, c2s []string

	for _, dlt := range dlts {
		// 存放当前这注开奖结果前后区各种的字符串组合
		var ic7s, ic6s, ic5s, ic4s, ic3s, ic2s []string
		ic7s, ic6s, ic5s, ic4s, ic3s, ic2s = nil, nil, nil, nil, nil, nil
		// 从开奖号码中生成7个组合号码
		ic7s = append(ic7s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 2)...)
		// 从开奖号码中生成6个组合号码
		ic6s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 1)...)
		ic6s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 2)...)
		// 从开奖号码中生成5个组合号码
		ic5s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 0)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 1)...)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 2)...)
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 1)...)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 2, 2)...)
		// 从开奖号码中生成3个组合号码
		ic3s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 0)
		ic3s = append(ic3s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 2, 1)...)
		ic3s = append(ic3s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 1, 2)...)
		// 从开奖号码中生成2个组合号码
		ic2s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 2, 0)
		ic2s = append(ic2s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 1, 1)...)
		ic2s = append(ic2s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 0, 2)...)
		for _, c7 := range ic7s {
			c7ms[c7] = append(c7ms[c7], dlt)
		}
		for _, c6 := range ic6s {
			c6ms[c6] = append(c6ms[c6], dlt)
		}
		for _, c5 := range ic5s {
			c5ms[c5] = append(c5ms[c5], dlt)
		}
		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], dlt)
		}
		for _, c3 := range ic3s {
			c3ms[c3] = append(c3ms[c3], dlt)
		}
		for _, c2 := range ic2s {
			c2ms[c2] = append(c2ms[c2], dlt)
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
		idlts := c7ms[v.Key]
		if len(idlts) > 1 {
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			var iidlts []models.Dlt
			for _, idlt := range idlts {
				c7KHave2DrawNums = append(c7KHave2DrawNums, idlt.DrawNum)
				iidlts = append(iidlts, idlt)
			}
			if len(iidlts) > 0 {
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
				for _, idlt := range iidlts {
					rows = append(rows, []string{idlt.DrawNum, idlt.DrawTime, fmt.Sprintf("%s %s %s %s %s | %s %s", idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2)})
				}
				dnf(rows)
				//for _, idlt := range iidlts {
				//	dnf([][]string{})
				//	pdf.Cell(0, opdf.CellHeight(chapterFontSize-3), fmt.Sprintf("%s %s -> %s %s %s %s %s | %s %s", idlt.DrawNum, idlt.DrawTime, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
				//	pdf.Ln(opdf.LineHeight(chapterFontSize - 3)) // 换行
				//}
			}
		}

	}
	tempNum = 0
	for k, v := range c6KLens {
		idlts := c6ms[v.Key]
		if limitLen65432[0] != 0 && k > limitLen65432[0] {
			break
		}
		//if k > 10 {
		//	break
		//}
		if len(idlts) > 1 {
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			var iidlts []models.Dlt
			for _, idlt := range idlts {
				if autoFilterMore {
					if !slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
						c6KHave2DrawNums = append(c6KHave2DrawNums, idlt.DrawNum)
						iidlts = append(iidlts, idlt)
					}
				} else {
					iidlts = append(iidlts, idlt)
				}

			}
			if len(iidlts) > 0 {
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
				//for _, idlt := range iidlts {
				//	pdf.Cell(0, opdf.CellHeight(chapterFontSize-3), fmt.Sprintf("%s %s -> %s %s %s %s %s | %s %s", idlt.DrawNum, idlt.DrawTime, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
				//	pdf.Ln(opdf.LineHeight(chapterFontSize - 3)) // 换行
				//}
				var rows [][]string
				for _, idlt := range iidlts {
					rows = append(rows, []string{idlt.DrawNum, idlt.DrawTime, fmt.Sprintf("%s %s %s %s %s | %s %s", idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2)})
				}
				dnf(rows)
			}
		}
	}
	tempNum = 0
	for k, v := range c5KLens {
		idlts := c5ms[v.Key]
		if limitLen65432[1] != 0 && k > limitLen65432[1] {
			break
		}
		//if k > 10 {
		//	break
		//}
		if len(idlts) > 1 {
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			var iidlts []models.Dlt
			for _, idlt := range idlts {
				if autoFilterMore {
					if !slices.Contains(c6KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
						c5KHave2DrawNums = append(c5KHave2DrawNums, idlt.DrawNum)
						iidlts = append(iidlts, idlt)
					}
				} else {
					iidlts = append(iidlts, idlt)
				}

			}
			if len(iidlts) > 0 {
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
				//for _, idlt := range iidlts {
				//	pdf.Cell(0, opdf.CellHeight(chapterFontSize-3), fmt.Sprintf("%s %s -> %s %s %s %s %s | %s %s", idlt.DrawNum, idlt.DrawTime, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
				//	pdf.Ln(opdf.LineHeight(chapterFontSize - 3)) // 换行
				//}
				var rows [][]string
				for _, idlt := range iidlts {
					rows = append(rows, []string{idlt.DrawNum, idlt.DrawTime, fmt.Sprintf("%s %s %s %s %s | %s %s", idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2)})
				}
				dnf(rows)
			}
		}
	}
	tempNum = 0
	for k, v := range c4KLens {
		idlts := c4ms[v.Key]
		if limitLen65432[2] != 0 && k > limitLen65432[2] {
			break
		}
		//if k > 10 {
		//	break
		//}
		if len(idlts) > 1 {
			sort.Slice(idlts, func(i, j int) bool {
				return idlts[i].DrawNum > idlts[j].DrawNum
			})
			var iidlts []models.Dlt
			for _, idlt := range idlts {
				if autoFilterMore {
					if !slices.Contains(c5KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c6KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
						c4KHave2DrawNums = append(c4KHave2DrawNums, idlt.DrawNum)
						iidlts = append(iidlts, idlt)
					}
				} else {
					iidlts = append(iidlts, idlt)
				}
			}
			if len(iidlts) > 0 {
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
				//for _, idlt := range iidlts {
				//	pdf.Cell(0, opdf.CellHeight(chapterFontSize-3), fmt.Sprintf("%s %s -> %s %s %s %s %s | %s %s", idlt.DrawNum, idlt.DrawTime, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
				//	pdf.Ln(opdf.LineHeight(chapterFontSize - 3)) // 换行
				//}
				var rows [][]string
				for _, idlt := range iidlts {
					rows = append(rows, []string{idlt.DrawNum, idlt.DrawTime, fmt.Sprintf("%s %s %s %s %s | %s %s", idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2)})
				}
				dnf(rows)
			}
		}
	}
	//tempNum = 0
	//for k, v := range c3KLens {
	//	idlts := c3ms[v.Key]
	//	if limitLen65432[3] != 0 && k > limitLen65432[3] {
	//		break
	//	}
	//	//if k > 10 {
	//	//	break
	//	//}
	//
	//	if len(idlts) > 1 {
	//		sort.Slice(idlts, func(i, j int) bool {
	//			return idlts[i].DrawNum > idlts[j].DrawNum
	//		})
	//		var iidlts []models.Dlt
	//		for _, idlt := range idlts {
	//			if autoFilterMore {
	//				if !slices.Contains(c4KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c5KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c6KHave2DrawNums, idlt.DrawNum) && !slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
	//					c3KHave2DrawNums = append(c3KHave2DrawNums, idlt.DrawNum)
	//					iidlts = append(iidlts, idlt)
	//				}
	//			} else {
	//				iidlts = append(iidlts, idlt)
	//			}
	//
	//		}
	//		if len(iidlts) > 0 {
	//			if tempNum == 0 {
	//				tempNum++
	//				str := fmt.Sprintf("重复3个号码")
	//				pdf.SetFont(fontName, "B", chapterFontSize-1)
	//				pdf.Cell(0, opdf.CellHeight(chapterFontSize-1), str)
	//				pdf.Bookmark(str, level+1, 0)
	//				pdf.Ln(opdf.LineHeight(chapterFontSize - 1)) // 换行
	//			}
	//			str := fmt.Sprintf("%s -> 已出现%d期", v.Key, v.Length)
	//			pdf.SetFont(fontName, "B", chapterFontSize-2)
	//			pdf.Cell(0, opdf.CellHeight(chapterFontSize-2), str)
	//			//pdf.Bookmark(str, level+2, 0)
	//			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
	//			pdf.SetFont(fontName, "L", chapterFontSize-3)
	//			//for _, idlt := range iidlts {
	//			//	pdf.Cell(0, opdf.CellHeight(chapterFontSize-3), fmt.Sprintf("%s %s -> %s %s %s %s %s | %s %s", idlt.DrawNum, idlt.DrawTime, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
	//			//	pdf.Ln(opdf.LineHeight(chapterFontSize - 3)) // 换行
	//			//}
	//			var rows [][]string
	//			for _, idlt := range iidlts {
	//				rows = append(rows, []string{idlt.DrawNum, idlt.DrawTime, fmt.Sprintf("%s %s %s %s %s | %s %s", idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2)})
	//			}
	//			dnf(rows)
	//		}
	//	}
	//}
	//tempNum = 0
	//for k, v := range c2KLens {
	//	idlts := c2ms[v.Key]
	//	//if limitLen65432[4] == 0  {
	//	//
	//	//}
	//	if limitLen65432[4] != 0 && k > limitLen65432[4] {
	//		break
	//	}
	//	//if k > 10 {
	//	//	break
	//	//}
	//	if len(idlts) > 1 {
	//		sort.Slice(idlts, func(i, j int) bool {
	//			return idlts[i].DrawNum > idlts[j].DrawNum
	//		})
	//		var iidlts []models.Dlt
	//		for _, idlt := range idlts {
	//			if autoFilterMore {
	//				if !slices.Contains(c3KHave2DrawNums, idlt.DrawNum) &&
	//					!slices.Contains(c4KHave2DrawNums, idlt.DrawNum) &&
	//					!slices.Contains(c5KHave2DrawNums, idlt.DrawNum) &&
	//					!slices.Contains(c6KHave2DrawNums, idlt.DrawNum) &&
	//					!slices.Contains(c7KHave2DrawNums, idlt.DrawNum) {
	//					iidlts = append(iidlts, idlt)
	//				}
	//			} else {
	//				iidlts = append(iidlts, idlt)
	//			}
	//		}
	//		if len(iidlts) > 0 {
	//
	//			if tempNum == 0 {
	//				tempNum++
	//				str := fmt.Sprintf("重复2个号码")
	//				pdf.SetFont(fontName, "B", chapterFontSize-1)
	//				pdf.Cell(0, opdf.CellHeight(chapterFontSize-1), str)
	//				pdf.Bookmark(str, level+1, 0)
	//				pdf.Ln(opdf.LineHeight(chapterFontSize - 1)) // 换行
	//			}
	//			str := fmt.Sprintf("%s -> 已出现%d期", v.Key, v.Length)
	//			pdf.SetFont(fontName, "B", chapterFontSize-2)
	//			pdf.Cell(0, opdf.CellHeight(chapterFontSize-2), str)
	//			//pdf.Bookmark(str, level+2, 0)
	//			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
	//			pdf.SetFont(fontName, "L", chapterFontSize-3)
	//			//for _, idlt := range iidlts {
	//			//	pdf.Cell(0, opdf.CellHeight(chapterFontSize-3), fmt.Sprintf("%s %s -> %s %s %s %s %s | %s %s", idlt.DrawNum, idlt.DrawTime, idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2))
	//			//	pdf.Ln(opdf.LineHeight(chapterFontSize - 3)) // 换行
	//			//}
	//			var rows [][]string
	//			for _, idlt := range iidlts {
	//				rows = append(rows, []string{idlt.DrawNum, idlt.DrawTime, fmt.Sprintf("%s %s %s %s %s | %s %s", idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2)})
	//			}
	//			dnf(rows)
	//		}
	//	}
	//}
}

// CHongHaoUpdateTable 更新 typ 表
//
//	@Description:
func CHongHaoUpdateTable() {
	dlts, _ := dbop.ReadAllDlt(false)
	// key为前后区号码的字符串组合，value为models.Dlt切片
	c7ms := make(map[string][]models.Dlt)
	c6ms := make(map[string][]models.Dlt)
	c5ms := make(map[string][]models.Dlt)
	c4ms := make(map[string][]models.Dlt)

	// 找到最新的 dlt_id
	lastTyp := dbop.ReadLastTyp()

	for _, dlt := range dlts {
		// 存放当前这注开奖结果前后区各种的字符串组合
		var ic7s, ic6s, ic5s, ic4s []string
		ic7s, ic6s, ic5s, ic4s = nil, nil, nil, nil
		// 从开奖号码中生成7个组合号码
		ic7s = append(ic7s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 2)...)
		// 从开奖号码中生成6个组合号码
		ic6s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 1)...)
		ic6s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 2)...)
		// 从开奖号码中生成5个组合号码
		ic5s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 0)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 1)...)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 2)...)
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 1)...)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 2, 2)...)

		for _, c7 := range ic7s {
			c7ms[c7] = append(c7ms[c7], dlt)
			if dlt.ID > lastTyp.DltID && len(c7ms[c7]) > 1 {
				for _, curDlt := range c7ms[c7] {
					if curDlt.ID != dlt.ID {
						//fmt.Printf("c7\n")
						dbop.InsertTyp(models.Typ{
							DltID:     dlt.ID,
							PrevDltId: curDlt.ID,
							Typ_1:     gen.CalDltTyp1(c7),
							Typ_2:     gen.CalDltTyp2(c7),
							Hm:        c7,
						})
					}
				}
			}
		}
		for _, c6 := range ic6s {
			c6ms[c6] = append(c6ms[c6], dlt)
			if dlt.ID > lastTyp.DltID && len(c6ms[c6]) > 1 {
				for _, curDlt := range c6ms[c6] {
					if curDlt.ID != dlt.ID {
						//fmt.Printf("c6\n")
						dbop.InsertTyp(models.Typ{
							DltID:     dlt.ID,
							PrevDltId: curDlt.ID,
							Typ_1:     gen.CalDltTyp1(c6),
							Typ_2:     gen.CalDltTyp2(c6),
							Hm:        c6,
						})
					}
				}
			}
		}
		for _, c5 := range ic5s {
			c5ms[c5] = append(c5ms[c5], dlt)
			if dlt.ID > lastTyp.DltID && len(c5ms[c5]) > 1 {
				for _, curDlt := range c5ms[c5] {
					if curDlt.ID != dlt.ID {
						//fmt.Printf("c5\n")
						dbop.InsertTyp(models.Typ{
							DltID:     dlt.ID,
							PrevDltId: curDlt.ID,
							Typ_1:     gen.CalDltTyp1(c5),
							Typ_2:     gen.CalDltTyp2(c5),
							Hm:        c5,
						})
					}
				}
			}
		}
		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], dlt)
			if dlt.ID > lastTyp.DltID && len(c4ms[c4]) > 1 {
				for _, curDlt := range c4ms[c4] {
					if curDlt.ID != dlt.ID {
						//fmt.Printf("c4\n")
						dbop.InsertTyp(models.Typ{
							DltID:     dlt.ID,
							PrevDltId: curDlt.ID,
							Typ_1:     gen.CalDltTyp1(c4),
							Typ_2:     gen.CalDltTyp2(c4),
							Hm:        c4,
						})
					}
				}
			}
		}
		//fmt.Printf("已完成第%s\n", dlt.DrawNum)
	}
	fmt.Printf("已处理完成\n")
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

//func CalCHLj1(cms map[string][]models.Dlt) (lj int) {
//	for _, v := range cms {
//		if len(v) > 1 {
//			lj += len(v)
//		}
//	}
//	return
//}

func CalCHLj(prevCms map[string][]string, iCms []string) (lj int) {
	for _, v := range iCms {
		if iv, ok := prevCms[v]; ok {
			if len(iv) == 1 {
				lj += 2
			} else if len(iv) >= 2 {
				lj += 1
			}
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
	dlts, _ := dbop.ReadAllDlt(false)
	drawNum2CHongHaoSt = make(map[string]CHongHaoSt, len(dlts)-1)

	c7ms := make(map[string][]string)
	c6ms := make(map[string][]string)
	c5ms := make(map[string][]string)
	c4ms := make(map[string][]string)

	var c7SLLj, c6SLLj, c5SLLj, c4SLLj int
	var ic7SLLj, ic6SLLj, ic5SLLj, ic4SLLj int

	for _, dlt := range dlts {
		var ic7s, ic6s, ic5s, ic4s []string
		ic7s, ic6s, ic5s, ic4s = nil, nil, nil, nil
		ic7SLLj, ic6SLLj, ic5SLLj, ic4SLLj = 0, 0, 0, 0
		fullStr := strings.Join([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, ",") + "|" + strings.Join([]string{dlt.B1, dlt.B2}, ",")
		// 从开奖号码中生成7个组合号码
		ic7s = append(ic7s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 2)...)
		// 从开奖号码中生成6个组合号码
		ic6s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 1)...)
		ic6s = append(ic6s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 2)...)
		// 从开奖号码中生成5个组合号码
		ic5s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 5, 0)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 1)...)
		ic5s = append(ic5s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 2)...)
		// 从开奖号码中生成4个组合号码
		ic4s = gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 4, 0)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 3, 1)...)
		ic4s = append(ic4s, gen.CrossComb([]string{dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5}, []string{dlt.B1, dlt.B2}, 2, 2)...)

		ic7SLLj = CalCHLj(c7ms, ic7s)
		c7SLLj += ic7SLLj
		for _, c7 := range ic7s {
			c7ms[c7] = append(c7ms[c7], fullStr)
		}

		ic6SLLj = CalCHLj(c6ms, ic6s)
		c6SLLj += ic6SLLj
		for _, c6 := range ic6s {
			c6ms[c6] = append(c6ms[c6], fullStr)
		}

		ic5SLLj = CalCHLj(c5ms, ic5s)
		c5SLLj += ic5SLLj
		for _, c5 := range ic5s {
			c5ms[c5] = append(c5ms[c5], fullStr)
		}
		ic4SLLj = CalCHLj(c4ms, ic4s)
		c4SLLj += ic4SLLj
		for _, c4 := range ic4s {
			c4ms[c4] = append(c4ms[c4], fullStr)
		}

		drawNum2CHongHaoSt[dlt.DrawNum] = CHongHaoSt{
			DrawNum:             dlt.DrawNum,
			NewAddCh4:           ic4SLLj,
			NewAddCh5:           ic5SLLj,
			NewAddCh6:           ic6SLLj,
			NewAddCh7:           ic7SLLj,
			DangQiTotalNewAddCh: ic4SLLj + ic5SLLj + ic6SLLj + ic7SLLj,
			LeiJiaCh:            c4SLLj + c5SLLj + c6SLLj + c7SLLj,
		}
	}
	return
}
