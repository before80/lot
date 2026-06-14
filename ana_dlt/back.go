package ana_dlt

import (
	"fmt"
	"slices"
	"sort"
	"sync"

	"github.com/before80/lot/dbop"
	"github.com/before80/lot/gen"
	"github.com/before80/lot/models"
	"github.com/before80/lot/opdf"
	"github.com/jung-kurt/gofpdf"
)

// DltBackQuShiChTable 后区号码可能的所有趋势
var DltBackQuShiChTable = []string{"ks", "kq", "kk", "kk_q", "ss", "ss_q", "sq", "sk", "qs", "qq", "qk"}

// DltBackQuShiStr 后区趋势
//
//	@Description:
//	@param nextNewDlt
//	@param b1
//	@param b2
//	@return string
func DltBackQuShiStr(nextNewDlt models.Dlt, b1, b2 string) string {
	if nextNewDlt.B1 < b1 && nextNewDlt.B2 < b2 {
		if nextNewDlt.B2 == b1 {
			// 1
			return "ss_q"
			//return "前缩后缩一重"
		}
		// 2
		return "ss"
		//return "前缩后缩"
	}

	if nextNewDlt.B1 < b1 && nextNewDlt.B2 == b2 {
		// 3
		return "sq"
		//return "前缩后重"
	}

	if nextNewDlt.B1 < b1 && nextNewDlt.B2 > b2 {
		// 4
		return "sk"
		//return "前缩后扩"
	}

	if nextNewDlt.B1 == b1 && nextNewDlt.B2 < b2 {
		// 5
		return "qs"
		//return "前重后缩"
	}

	if nextNewDlt.B1 == b1 && nextNewDlt.B2 == b2 {
		// 6
		return "qq"
		//return "一致"
	}

	if nextNewDlt.B1 == b1 && nextNewDlt.B2 > b2 {
		// 7
		return "qk"
		//return "前重后扩"
	}

	if nextNewDlt.B1 > b1 && nextNewDlt.B2 < b2 {
		// 8
		return "ks"
		//return "前扩后缩"
	}

	if nextNewDlt.B1 > b1 && nextNewDlt.B2 == b2 {
		// 9
		return "kq"
		//return "前扩后重"
	}

	if nextNewDlt.B1 > b1 && nextNewDlt.B2 > b2 {
		if nextNewDlt.B1 == b2 {
			// 10
			return "kk_q"
			//return "前扩后扩一同"
		}
		// 11
		return "kk"
		//return "前扩后扩"
	}
	return "kk"
	//return "前扩后扩"
}

// DltBackQuShiToStdOut 后区趋势
//
//	@Description:
func DltBackQuShiToStdOut() {
	dlts, _ := dbop.ReadAllDlt(false)
	// 生成后区所有组合字符串的切片
	backCombs := gen.Comb(gen.AllDltBackHms, 2)
	//bcs := make([][]string, len(backCombs))
	// 后区两个组合号为key，map[string]int为value（其中趋势字符组合为key，[]models.Dlt为value）
	bcMaps := make(map[string]map[string][]models.Dlt)
	// 后区两个组合号为key，map[string]int为value（其中趋势字符组合为key，出现次数为value）
	bcQMaps := make(map[string]map[string]int)
	// 后区两个组合号为key，[]string为value（其中元素为没有出现的趋势字符组合）
	bcNonQMaps := make(map[string][]string)

	for i, dlt := range dlts {
		if i+1 == len(dlts) {
			break
		}
		comb := dlt.B1 + "," + dlt.B2
		if _, ok1 := bcMaps[comb]; !ok1 {
			bcMaps[comb] = make(map[string][]models.Dlt)
		}
		qs := DltBackQuShiStr(dlts[i+1], dlt.B1, dlt.B2)
		bcMaps[comb][qs] = append(bcMaps[comb][qs], dlts[i+1])
	}

	for _, backComb := range backCombs {
		if _, ok2 := bcMaps[backComb]; ok2 {
			for _, chT := range DltBackQuShiChTable {
				if _, ok3 := bcMaps[backComb][chT]; ok3 {
					if _, ok4 := bcQMaps[backComb]; !ok4 {
						bcQMaps[backComb] = make(map[string]int)
					}
					bcQMaps[backComb][chT] = len(bcMaps[backComb][chT])
					sort.Slice(bcMaps[backComb][chT], func(i, j int) bool {
						return bcMaps[backComb][chT][i].DrawNum < bcMaps[backComb][chT][j].DrawNum
					})
					fmt.Printf("-> %s %s -----------\n", backComb, chT)
					for _, dlt := range bcMaps[backComb][chT] {
						fmt.Printf("%s %s %s %s %s %s %s %s %s\n", dlt.DrawNum, dlt.DrawTime, dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2)
					}
					fmt.Printf("---------------------------- \n")
				}
			}
		}
	}

	for _, backComb := range backCombs {
		canQss, _ := gen.GetDltAllQuShiFromSpBackComb(backComb)
		for _, chT := range canQss {
			if _, ok5 := bcQMaps[backComb][chT]; !ok5 {
				bcNonQMaps[backComb] = append(bcNonQMaps[backComb], chT)
			}
		}
	}
	var hmkgHaoMao []string
	var allNotAppearQsNum int
	for _, backComb := range backCombs {
		//fmt.Printf("%s 还未出现的趋势有: %v\n", backComb, bcNonQMaps[backComb])
		if len(bcNonQMaps[backComb]) > 0 {
			l := 0
			for _, v := range bcQMaps[backComb] {
				l = l + v
			}
			hmkgHaoMao = nil
			for _, chT := range bcNonQMaps[backComb] {
				items, _ := gen.GetDltBackQuShiHaoMasFromQuShi(backComb, chT)
				hmkgHaoMao = append(hmkgHaoMao, items...)
			}

			allQs, _ := gen.GetDltAllQuShiFromSpBackComb(backComb)

			sort.Strings(hmkgHaoMao)
			fmt.Printf("%s 历已开%d期数，所有可能出现的趋势有%d个（%q） 还未出现的趋势有（%d）: %v 即没开过的还有 %q \n", backComb, l, len(allQs), allQs, len(bcNonQMaps[backComb]), bcNonQMaps[backComb], hmkgHaoMao)
			allNotAppearQsNum += len(bcNonQMaps[backComb])
		} else {
			l := 0
			for _, v := range bcQMaps[backComb] {
				l = l + v
			}
			hmkgHaoMao = nil
			for _, chT := range bcNonQMaps[backComb] {
				items, _ := gen.GetDltBackQuShiHaoMasFromQuShi(backComb, chT)
				hmkgHaoMao = append(hmkgHaoMao, items...)
			}

			allQs, _ := gen.GetDltAllQuShiFromSpBackComb(backComb)
			fmt.Printf("%s 历已开%d期数，所有可能出现的趋势有%d个（%q） 还未出现的趋势有（%d）: %v 即没开过的还有 %q \n", backComb, l, len(allQs), allQs, len(bcNonQMaps[backComb]), bcNonQMaps[backComb], hmkgHaoMao)
		}
	}
	fmt.Printf("总共还有%d个趋势没有出现------------------------------\n", allNotAppearQsNum)
	fmt.Printf("-------------------------------\n")
	var allCanAppearQsNum int
	for _, backComb := range backCombs {
		allQs, _ := gen.GetDltAllQuShiFromSpBackComb(backComb)
		allCanAppearQsNum += len(allQs)
		fmt.Printf("%s %q\n", backComb, allQs)
	}
	fmt.Printf("总共可以出现%d个趋势------------------------------\n", allCanAppearQsNum)
}

// DltBackQuShi1 大乐透后区趋势1(后区号码->后区号码)
//
//	@Description:
//	@return res
func DltBackQuShi1() (res []KeyWithLength) {
	// 生成后区所有组合字符串的切片
	backCombs := gen.Comb(gen.AllDltBackHms, 2)
	quShi2Count := make(map[string]int)
	// 初始化
	for _, backComb1 := range backCombs {
		for _, backComb2 := range backCombs {
			quShi2Count[fmt.Sprintf("%s->%s", backComb1, backComb2)] = 0
		}
	}
	for i, dlt := range ZxDlts {
		if i == len(ZxDlts)-1 {
			break
		}
		nextDlt := ZxDlts[i+1]
		quShiStr := fmt.Sprintf("%s,%s->%s,%s", dlt.B1, dlt.B2, nextDlt.B1, nextDlt.B2)
		quShi2Count[quShiStr] = quShi2Count[quShiStr] + 1
	}

	kLens := make([]KeyWithLength, 0, len(quShi2Count))
	for k, v := range quShi2Count {
		kLens = append(kLens, KeyWithLength{k, v})
	}
	sort.Slice(kLens, func(i, j int) bool {
		if kLens[i].Length == kLens[j].Length {
			return kLens[i].Key > kLens[j].Key
		}
		return kLens[i].Length > kLens[j].Length
	})
	return kLens
}

type DltBackChQuShi struct {
	BackComb         string
	Qs               string
	Cs               int
	AllCombs         []string // 所有可能的组合
	HadExistCombs    []string // 已经出现的组合
	HadNotExistCombs []string // 未出现的组合
}

// DltBackQuShi2 大乐透后区趋势2(后区号码->后区skq)
//
//	@Description:
//	@return res
func DltBackQuShi2() (res []DltBackChQuShi) {
	quShi2St := make(map[string]*DltBackChQuShi)

	for i, dlt := range ZxDlts {
		if i == len(ZxDlts)-1 {
			break
		}
		nextDlt := ZxDlts[i+1]
		qs := DltBackQuShiStr(nextDlt, dlt.B1, dlt.B2)
		backComb := fmt.Sprintf("%s,%s", dlt.B1, dlt.B2)
		nextBackComb := fmt.Sprintf("%s,%s", nextDlt.B1, nextDlt.B2)
		typ := fmt.Sprintf("%s->%s", backComb, qs)
		if _, ok := quShi2St[typ]; !ok {
			allCombs, _ := gen.GetDltBackQuShiHaoMasFromQuShi(backComb, qs)
			quShi2St[typ] = &DltBackChQuShi{
				BackComb: backComb,
				Qs:       qs,
				Cs:       0,
				AllCombs: allCombs,
			}
		}
		quShi2St[typ].Cs = quShi2St[typ].Cs + 1
		if !slices.Contains(quShi2St[typ].HadExistCombs, nextBackComb) {
			quShi2St[typ].HadExistCombs = append(quShi2St[typ].HadExistCombs, nextBackComb)
		}
	}

	for k, v := range quShi2St {
		quShi2St[k].HadNotExistCombs = gen.DiffSlice(v.AllCombs, v.HadExistCombs)
	}

	kLens := make([]KeyWithLength, 0, len(quShi2St))
	for k, v := range quShi2St {
		kLens = append(kLens, KeyWithLength{k, v.Cs})
	}

	sort.Slice(kLens, func(i, j int) bool {
		if kLens[i].Length == kLens[j].Length {
			return kLens[i].Key > kLens[j].Key
		}
		return kLens[i].Length > kLens[j].Length
	})
	for _, v := range kLens {
		res = append(res, *quShi2St[v.Key])
	}
	return res
}

// DltBackQuShiForPDF 用于在pdf文档中生成后区趋势
//
//	@Description:
//	@param pdf
//	@param fontName 在pdf文档中使用的字体
//	@param chapterFontSize 在pdf文档中使用的字体大小
//	@param chapter 将生成的内容在pdf中的标题名称
//	@param level 将生成的内容在pdf中的标题层级
func DltBackQuShiForPDF(pdf *gofpdf.Fpdf, fontName string, chapterFontSize float64, chapter string, level int) {
	pdf.SetFont(fontName, "B", chapterFontSize)
	pdf.Cell(0, opdf.CellHeight(chapterFontSize), chapter)
	pdf.Bookmark(chapter, level, 0)
	pdf.Ln(opdf.LineHeight(chapterFontSize))

	pdf.AddPage()
	pdf.SetFont(fontName, "B", chapterFontSize-1)
	pdf.Cell(0, opdf.CellHeight(chapterFontSize-1), "s k q的解释")
	pdf.Bookmark("s k q的解释", level+1, 0)
	pdf.Ln(opdf.LineHeight(chapterFontSize - 1))
	pdf.Image("suokuo.png", 5, 30, 200, 0, false, "", 0, "")
	pdf.AddPage()
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
	// 生成后区所有组合字符串的切片
	backCombs := gen.Comb(gen.AllDltBackHms, 2)
	//bcs := make([][]string, len(backCombs))
	// 后区两个组合号为key，map[string]int为value（其中趋势字符组合为key，[]models.Dlt为value）
	bcMaps := make(map[string]map[string][]models.Dlt)
	// 后区两个组合号为key，map[string]int为value（其中趋势字符组合为key，出现次数为value）
	bcQMaps := make(map[string]map[string]int)
	// 后区两个组合号为key，[]string为value（其中元素为没有出现的趋势字符组合）
	bcNonQMaps := make(map[string][]string)
	comb2ExistCombSlice := make(map[string][]string, 0)

	for i, dlt := range dlts {
		if i+1 == len(dlts) {
			break
		}
		comb := dlt.B1 + "," + dlt.B2
		if _, ok1 := bcMaps[comb]; !ok1 {
			bcMaps[comb] = make(map[string][]models.Dlt)
		}
		qs := DltBackQuShiStr(dlts[i+1], dlt.B1, dlt.B2)
		nextComb := dlts[i+1].B1 + "," + dlts[i+1].B2
		if _, ok2 := comb2ExistCombSlice[comb]; !ok2 {
			comb2ExistCombSlice[comb] = append(comb2ExistCombSlice[comb], nextComb)
		} else {
			if !slices.Contains(comb2ExistCombSlice[comb], nextComb) {
				comb2ExistCombSlice[comb] = append(comb2ExistCombSlice[comb], nextComb)
			}
		}
		bcMaps[comb][qs] = append(bcMaps[comb][qs], dlts[i+1])
	}

	for _, backComb := range backCombs {
		if _, ok2 := bcMaps[backComb]; ok2 {
			for _, chT := range DltBackQuShiChTable {
				if _, ok3 := bcMaps[backComb][chT]; ok3 {
					if _, ok4 := bcQMaps[backComb]; !ok4 {
						bcQMaps[backComb] = make(map[string]int)
					}
					bcQMaps[backComb][chT] = len(bcMaps[backComb][chT])
					sort.Slice(bcMaps[backComb][chT], func(i, j int) bool {
						return bcMaps[backComb][chT][i].DrawNum < bcMaps[backComb][chT][j].DrawNum
					})
				}
			}
		}
	}

	for _, backComb := range backCombs {
		canQss, _ := gen.GetDltAllQuShiFromSpBackComb(backComb)
		for _, chT := range canQss {
			if _, ok5 := bcQMaps[backComb][chT]; !ok5 {
				bcNonQMaps[backComb] = append(bcNonQMaps[backComb], chT)
			}
		}
	}
	var hmkgHaoMao []string
	var allNotAppearQsNum int
	for _, backComb := range backCombs {
		//fmt.Printf("%s 还未出现的趋势有: %v\n", backComb, bcNonQMaps[backComb])
		if len(bcNonQMaps[backComb]) > 0 {
			l := 0
			for _, v := range bcQMaps[backComb] {
				l = l + v
			}
			hmkgHaoMao = nil
			for _, chT := range bcNonQMaps[backComb] {
				items, _ := gen.GetDltBackQuShiHaoMasFromQuShi(backComb, chT)
				hmkgHaoMao = append(hmkgHaoMao, items...)
			}

			allQs, _ := gen.GetDltAllQuShiFromSpBackComb(backComb)
			sort.Strings(hmkgHaoMao)
			str := fmt.Sprintf("%s 的历史统计", backComb)
			pdf.SetFont(fontName, "B", chapterFontSize-1)
			pdf.Cell(0, opdf.CellHeight(chapterFontSize-1), str)
			pdf.Bookmark(str, level+1, 0)
			pdf.Ln(opdf.LineHeight(chapterFontSize - 1)) // 换行
			pdf.SetFont(fontName, "L", chapterFontSize-2)

			pdf.Cell(0, opdf.CellHeight(chapterFontSize-2), fmt.Sprintf("%s 历史已开%d期数", backComb, l))
			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
			pdf.MultiCell(0, opdf.CellHeight(chapterFontSize-2), fmt.Sprintf("所有可能出现的趋势有%d个，即 %s", len(allQs), opdf.SliceToStr(allQs, ",")), "0", "L", false)
			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
			pdf.MultiCell(0, opdf.CellHeight(chapterFontSize-2), fmt.Sprintf("还未出现的趋势有%d个： %s，这些趋势对应号码（还没有出现的）：%s", len(bcNonQMaps[backComb]), opdf.SliceToStr(bcNonQMaps[backComb], ","), opdf.SliceToStr(hmkgHaoMao, "|")), "0", "L", false)
			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
			notAppearCombs := gen.DiffSlice(backCombs, comb2ExistCombSlice[backComb])
			pdf.MultiCell(0, opdf.CellHeight(chapterFontSize-2), fmt.Sprintf("已出现的号码有%d个： %s", len(comb2ExistCombSlice[backComb]), opdf.SliceToStr(comb2ExistCombSlice[backComb], "|")), "0", "L", false)
			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
			pdf.MultiCell(0, opdf.CellHeight(chapterFontSize-2), fmt.Sprintf("综上，还没有出现的号码有%d个： %s", len(notAppearCombs), opdf.SliceToStr(notAppearCombs, "|")), "0", "L", false)

			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
			allNotAppearQsNum += len(bcNonQMaps[backComb])
		} else {
			l := 0
			for _, v := range bcQMaps[backComb] {
				l = l + v
			}
			hmkgHaoMao = nil
			for _, chT := range bcNonQMaps[backComb] {
				items, _ := gen.GetDltBackQuShiHaoMasFromQuShi(backComb, chT)
				hmkgHaoMao = append(hmkgHaoMao, items...)
			}

			allQs, _ := gen.GetDltAllQuShiFromSpBackComb(backComb)
			str := fmt.Sprintf("%s 的历史统计", backComb)
			pdf.SetFont(fontName, "B", chapterFontSize-1)
			pdf.Cell(0, opdf.CellHeight(chapterFontSize-1), str)
			pdf.Bookmark(str, level+1, 0)
			pdf.Ln(opdf.LineHeight(chapterFontSize - 1)) // 换行
			pdf.SetFont(fontName, "L", chapterFontSize-2)

			pdf.Cell(0, opdf.CellHeight(chapterFontSize-2), fmt.Sprintf("%s 历史已开%d期数", backComb, l))
			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
			pdf.MultiCell(0, opdf.CellHeight(chapterFontSize-2), fmt.Sprintf("所有可能出现的趋势有%d个，即 %s", len(allQs), opdf.SliceToStr(allQs, ",")), "0", "L", false)
			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
			pdf.MultiCell(0, opdf.CellHeight(chapterFontSize-2), fmt.Sprintf("还未出现的趋势有%d个： %s，这些趋势对应号码（还没有出现的）：%s", len(bcNonQMaps[backComb]), opdf.SliceToStr(bcNonQMaps[backComb], ","), opdf.SliceToStr(hmkgHaoMao, "|")), "0", "L", false)
			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
			notAppearCombs := gen.DiffSlice(backCombs, comb2ExistCombSlice[backComb])
			pdf.MultiCell(0, opdf.CellHeight(chapterFontSize-2), fmt.Sprintf("已出现的号码有%d个： %s", len(comb2ExistCombSlice[backComb]), opdf.SliceToStr(comb2ExistCombSlice[backComb], "|")), "0", "L", false)
			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
			pdf.MultiCell(0, opdf.CellHeight(chapterFontSize-2), fmt.Sprintf("综上，还没有出现的号码有%d个： %s", len(notAppearCombs), opdf.SliceToStr(notAppearCombs, "|")), "0", "L", false)

			pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
		}
	}

	for _, backComb := range backCombs {
		canExistChTs, _ := gen.GetDltAllQuShiFromSpBackComb(backComb)
		if _, ok2 := bcMaps[backComb]; ok2 {
			for _, chT := range DltBackQuShiChTable {
				if slices.Contains(canExistChTs, chT) {
					str := fmt.Sprintf("%s -> %s （%d）", backComb, chT, len(bcMaps[backComb][chT]))
					pdf.SetFont(fontName, "B", chapterFontSize-2)
					pdf.Cell(0, opdf.CellHeight(chapterFontSize-2), str)
					pdf.Bookmark(str, level+1, 0)
					pdf.Ln(opdf.LineHeight(chapterFontSize - 2)) // 换行
					pdf.SetFont(fontName, "L", chapterFontSize-3)

					//for _, dlt := range bcMaps[backComb][chT] {
					//	fmt.Printf("%s %s %s %s %s %s %s %s %s\n", dlt.DrawNum, dlt.DrawTime, dlt.F1, dlt.F2, dlt.F3, dlt.F4, dlt.F5, dlt.B1, dlt.B2)
					//}

					var rows [][]string
					for _, idlt := range bcMaps[backComb][chT] {
						rows = append(rows, []string{idlt.DrawNum, idlt.DrawTime, fmt.Sprintf("%s %s %s %s %s | %s %s", idlt.F1, idlt.F2, idlt.F3, idlt.F4, idlt.F5, idlt.B1, idlt.B2)})
					}
					dnf(rows)
				}
			}
		}
	}

	//
	//fmt.Printf("总共还有%d个趋势没有出现------------------------------\n", allNotAppearQsNum)
	//fmt.Printf("-------------------------------\n")
	//var allCanAppearQsNum int
	//for _, backComb := range backCombs {
	//	allQs, _ := gen.GetDltAllQuShiFromSpBackComb(backComb)
	//	allCanAppearQsNum += len(allQs)
	//	fmt.Printf("%s %q\n", backComb, allQs)
	//}
	//fmt.Printf("总共可以出现%d个趋势------------------------------\n", allCanAppearQsNum)
}

// Zhs
//
//	@Description:
func Zhs() {
	// 生成后区所有组合字符串的切片
	backCombs := gen.Comb(gen.AllDltBackHms, 2)
	for _, backComb := range backCombs {
		for _, chT := range DltBackQuShiChTable {
			res, _ := gen.GetDltBackQuShiHaoMasFromQuShi(backComb, chT)
			fmt.Printf("%s -> %6s -> len=%2d -> %q\n", backComb, chT, len(res), res)
		}
	}
}

func deepCopyNestedMap(src map[string]map[string]int) map[string]map[string]int {
	dst := make(map[string]map[string]int, len(src))
	for k, v := range src {
		inner := make(map[string]int, len(v))
		for ik, iv := range v {
			inner[ik] = iv
		}
		dst[k] = inner
	}
	return dst
}

func sliceRemoveFirst(s []string, target string) []string {
	for i, v := range s {
		if v == target {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s // 没找到则返回原切片
}

func ValidateDltBackQuShi() {
	// 只分析最新的多少期
	dlts, _ := dbop.ReadAllDlt(false)
	// 生成后区所有组合字符串的切片
	backCombs := gen.Comb(gen.AllDltBackHms, 2)

	backComb2Qs := make(map[string][]string)
	var startAllCanQs []string // 元素类似：01,02->kk
	var itemQs []string
	for _, backComb := range backCombs {
		itemQs = nil
		itemQs, _ = gen.GetDltAllQuShiFromSpBackComb(backComb)
		backComb2Qs[backComb] = itemQs
		for _, chT := range itemQs {
			startAllCanQs = append(startAllCanQs, backComb+"->"+chT)
		}
	}
	drawNum2CurNotAppearQs := make(map[string][]string)
	curNotAppearQs := make([]string, len(startAllCanQs))
	copy(curNotAppearQs, startAllCanQs)
	var allDrawNums []string
	for i, dlt := range dlts {
		if i+1 == len(dlts) {
			break
		}
		comb := dlt.B1 + "," + dlt.B2
		qs := DltBackQuShiStr(dlts[i+1], dlt.B1, dlt.B2)
		combQs := comb + "->" + qs
		curNotAppearQs = sliceRemoveFirst(curNotAppearQs, combQs)
		drawNum2CurNotAppearQs[dlt.DrawNum] = curNotAppearQs
		allDrawNums = append(allDrawNums, dlt.DrawNum)
	}
	var lastDrawNum string

	for i, drawNum := range allDrawNums {
		if i+1 == len(allDrawNums) {
			lastDrawNum = drawNum
		}
		fmt.Printf("%s -> len=%d\n", drawNum, len(drawNum2CurNotAppearQs[drawNum]))
	}
	fmt.Printf("---------------------------- \n")
	for _, iCurNotAppearQs := range drawNum2CurNotAppearQs[lastDrawNum] {
		fmt.Printf("%s\n", iCurNotAppearQs)
	}

}

// DltBackHis 大乐透后区历史数据(按出现期数的大小,从大到小排序)
//
//	@Description:
//	@return res
func DltBackHis(wg *sync.WaitGroup) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	typ2DltHis := make(map[string]*DltHis)
	backCombs := gen.Comb(gen.AllDltBackHms, 2)
	// 初始化 typ2DltHis
	for _, backComb := range backCombs {
		typ2DltHis[backComb] = &DltHis{Typ: backComb, AllCount: 324632}
	}
	lenDltHis := len(ZxDlts)
	for i, dlt := range ZxDlts {
		curBackComb := dlt.B1 + "," + dlt.B2
		typ2DltHis[curBackComb].Cs = typ2DltHis[curBackComb].Cs + 1
		if lenDltHis-i <= 10 {
			typ2DltHis[curBackComb].Last10 = typ2DltHis[curBackComb].Last10 + 1
		}
		if lenDltHis-i <= 20 {
			typ2DltHis[curBackComb].Last20 = typ2DltHis[curBackComb].Last20 + 1
		}
		if lenDltHis-i <= 30 {
			typ2DltHis[curBackComb].Last30 = typ2DltHis[curBackComb].Last30 + 1
		}
		if lenDltHis-i <= 50 {
			typ2DltHis[curBackComb].Last50 = typ2DltHis[curBackComb].Last50 + 1
		}
		if lenDltHis-i <= 100 {
			typ2DltHis[curBackComb].Last100 = typ2DltHis[curBackComb].Last100 + 1
		}
		if lenDltHis-i <= 200 {
			typ2DltHis[curBackComb].Last200 = typ2DltHis[curBackComb].Last200 + 1
		}
		if lenDltHis-i <= 500 {
			typ2DltHis[curBackComb].Last500 = typ2DltHis[curBackComb].Last500 + 1
		}
		if lenDltHis-i <= 1000 {
			typ2DltHis[curBackComb].Last1000 = typ2DltHis[curBackComb].Last1000 + 1
		}
		if lenDltHis-i <= 1500 {
			typ2DltHis[curBackComb].Last1500 = typ2DltHis[curBackComb].Last1500 + 1
		}
		if lenDltHis-i <= 2000 {
			typ2DltHis[curBackComb].Last2000 = typ2DltHis[curBackComb].Last2000 + 1
		}
		if lenDltHis-i <= 2500 {
			typ2DltHis[curBackComb].Last2500 = typ2DltHis[curBackComb].Last2500 + 1
		}
		if lenDltHis-i <= 3500 {
			typ2DltHis[curBackComb].Last3500 = typ2DltHis[curBackComb].Last3500 + 1
		}
	}

	kLens := make([]KeyWithLength, 0, len(typ2DltHis))
	for k, dltHis := range typ2DltHis {
		kLens = append(kLens, KeyWithLength{Key: k, Length: dltHis.Cs})
	}
	// 对typ2DltHis按照存放的Cs值的大小进行排序
	sort.Slice(kLens, func(i, j int) bool {
		return kLens[i].Length > kLens[j].Length
	})
	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2DltHis[typ])
	}
	return
}

// DltEqBackHis 设备的大乐透后区历史数据(按出现期数的大小,从大到小排序)
func DltEqBackHis(wg *sync.WaitGroup, eqNumCount int) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	typ2DltHis := make(map[string]*DltHis)
	backCombs := gen.Comb(gen.AllDltBackHms, 2)
	// 初始化 typ2DltHis
	for _, backComb := range backCombs {
		typ2DltHis[backComb] = &DltHis{Typ: backComb, AllCount: 324632}
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

		curBackComb := dlt.B1 + "," + dlt.B2
		typ2DltHis[curBackComb].Cs = typ2DltHis[curBackComb].Cs + 1
		if lenDltHis-i <= 10 {
			typ2DltHis[curBackComb].Last10 = typ2DltHis[curBackComb].Last10 + 1
		}
		if lenDltHis-i <= 20 {
			typ2DltHis[curBackComb].Last20 = typ2DltHis[curBackComb].Last20 + 1
		}
		if lenDltHis-i <= 30 {
			typ2DltHis[curBackComb].Last30 = typ2DltHis[curBackComb].Last30 + 1
		}
		if lenDltHis-i <= 50 {
			typ2DltHis[curBackComb].Last50 = typ2DltHis[curBackComb].Last50 + 1
		}
		if lenDltHis-i <= 100 {
			typ2DltHis[curBackComb].Last100 = typ2DltHis[curBackComb].Last100 + 1
		}
		if lenDltHis-i <= 200 {
			typ2DltHis[curBackComb].Last200 = typ2DltHis[curBackComb].Last200 + 1
		}
		if lenDltHis-i <= 500 {
			typ2DltHis[curBackComb].Last500 = typ2DltHis[curBackComb].Last500 + 1
		}
		if lenDltHis-i <= 1000 {
			typ2DltHis[curBackComb].Last1000 = typ2DltHis[curBackComb].Last1000 + 1
		}
		if lenDltHis-i <= 1500 {
			typ2DltHis[curBackComb].Last1500 = typ2DltHis[curBackComb].Last1500 + 1
		}
		if lenDltHis-i <= 2000 {
			typ2DltHis[curBackComb].Last2000 = typ2DltHis[curBackComb].Last2000 + 1
		}
		if lenDltHis-i <= 2500 {
			typ2DltHis[curBackComb].Last2500 = typ2DltHis[curBackComb].Last2500 + 1
		}
		if lenDltHis-i <= 3500 {
			typ2DltHis[curBackComb].Last3500 = typ2DltHis[curBackComb].Last3500 + 1
		}
		i++
	}

	kLens := make([]KeyWithLength, 0, len(typ2DltHis))
	for k, dltHis := range typ2DltHis {
		kLens = append(kLens, KeyWithLength{Key: k, Length: dltHis.Cs})
	}
	// 对typ2DltHis按照存放的Cs值的大小进行排序
	sort.Slice(kLens, func(i, j int) bool {
		return kLens[i].Length > kLens[j].Length
	})
	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2DltHis[typ])
	}
	return
}

// DltBackOnlyOneHis 大乐透后区单个号码的历史
//
//	@Description:
//	@param wg
//	@return res
func DltBackOnlyOneHis(wg *sync.WaitGroup) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2DltHis := make(map[string]*DltHis)

	// 初始化
	for k, v := range AllDltOnlyOneBack2Count {
		typ2DltHis[k] = &DltHis{Typ: k, AllCount: v}
	}

	lenDltHis := len(ZxDlts)

	for i, dlt := range ZxDlts {
		backHms := []string{dlt.B1, dlt.B2}

		for _, bHm := range backHms {
			typ2DltHis[bHm].Cs = typ2DltHis[bHm].Cs + 1
			if lenDltHis-i <= 10 {
				typ2DltHis[bHm].Last10 = typ2DltHis[bHm].Last10 + 1
			}
			if lenDltHis-i <= 20 {
				typ2DltHis[bHm].Last20 = typ2DltHis[bHm].Last20 + 1
			}
			if lenDltHis-i <= 30 {
				typ2DltHis[bHm].Last30 = typ2DltHis[bHm].Last30 + 1
			}
			if lenDltHis-i <= 50 {
				typ2DltHis[bHm].Last50 = typ2DltHis[bHm].Last50 + 1
			}
			if lenDltHis-i <= 100 {
				typ2DltHis[bHm].Last100 = typ2DltHis[bHm].Last100 + 1
			}
			if lenDltHis-i <= 200 {
				typ2DltHis[bHm].Last200 = typ2DltHis[bHm].Last200 + 1
			}
			if lenDltHis-i <= 500 {
				typ2DltHis[bHm].Last500 = typ2DltHis[bHm].Last500 + 1
			}
			if lenDltHis-i <= 1000 {
				typ2DltHis[bHm].Last1000 = typ2DltHis[bHm].Last1000 + 1
			}
			if lenDltHis-i <= 1500 {
				typ2DltHis[bHm].Last1500 = typ2DltHis[bHm].Last1500 + 1
			}
			if lenDltHis-i <= 2000 {
				typ2DltHis[bHm].Last2000 = typ2DltHis[bHm].Last2000 + 1
			}
			if lenDltHis-i <= 2500 {
				typ2DltHis[bHm].Last2500 = typ2DltHis[bHm].Last2500 + 1
			}
			if lenDltHis-i <= 3500 {
				typ2DltHis[bHm].Last3500 = typ2DltHis[bHm].Last3500 + 1
			}
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
		return kLens[i].Length > kLens[j].Length
	})

	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2DltHis[typ])
	}

	return
}

// DltEqBackOnlyOneHis 设备的大乐透后区单个号码的历史
func DltEqBackOnlyOneHis(wg *sync.WaitGroup, eqNumCount int) (res []DltHis) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	typ2DltHis := make(map[string]*DltHis)

	// 初始化
	for k, v := range AllDltOnlyOneBack2Count {
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

		backHms := []string{dlt.B1, dlt.B2}

		for _, bHm := range backHms {
			typ2DltHis[bHm].Cs = typ2DltHis[bHm].Cs + 1
			if lenDltHis-i <= 10 {
				typ2DltHis[bHm].Last10 = typ2DltHis[bHm].Last10 + 1
			}
			if lenDltHis-i <= 20 {
				typ2DltHis[bHm].Last20 = typ2DltHis[bHm].Last20 + 1
			}
			if lenDltHis-i <= 30 {
				typ2DltHis[bHm].Last30 = typ2DltHis[bHm].Last30 + 1
			}
			if lenDltHis-i <= 50 {
				typ2DltHis[bHm].Last50 = typ2DltHis[bHm].Last50 + 1
			}
			if lenDltHis-i <= 100 {
				typ2DltHis[bHm].Last100 = typ2DltHis[bHm].Last100 + 1
			}
			if lenDltHis-i <= 200 {
				typ2DltHis[bHm].Last200 = typ2DltHis[bHm].Last200 + 1
			}
			if lenDltHis-i <= 500 {
				typ2DltHis[bHm].Last500 = typ2DltHis[bHm].Last500 + 1
			}
			if lenDltHis-i <= 1000 {
				typ2DltHis[bHm].Last1000 = typ2DltHis[bHm].Last1000 + 1
			}
			if lenDltHis-i <= 1500 {
				typ2DltHis[bHm].Last1500 = typ2DltHis[bHm].Last1500 + 1
			}
			if lenDltHis-i <= 2000 {
				typ2DltHis[bHm].Last2000 = typ2DltHis[bHm].Last2000 + 1
			}
			if lenDltHis-i <= 2500 {
				typ2DltHis[bHm].Last2500 = typ2DltHis[bHm].Last2500 + 1
			}
			if lenDltHis-i <= 3500 {
				typ2DltHis[bHm].Last3500 = typ2DltHis[bHm].Last3500 + 1
			}
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
		return kLens[i].Length > kLens[j].Length
	})

	for _, kvLen := range kLens {
		typ := kvLen.Key
		res = append(res, *typ2DltHis[typ])
	}

	return
}
