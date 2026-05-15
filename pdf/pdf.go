package main

import (
	"fmt"
	"strconv"

	"github.com/before80/lot/ana_dlt"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/opdf"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// 设置文档保护 - 禁止编辑、复制、打印等
	pdf.SetProtection(
		gofpdf.CnProtectPrint,  // 允许打印
		"",                     // 用户密码（空字符串表示不需要密码打开）
		"owner_abcdefghijklmn", // 所有者密码（用于设置权限）
	)

	//pdf.AddUTF8Font("sourcehan", "L", "fonts/SourceHanSansSC-Light.ttf")
	pdf.AddUTF8Font("sourcehan", "B", "fonts/SourceHanSansSC-Bold.ttf")
	pdf.AddUTF8Font("sourcehan", "L", "fonts/SourceHanSansSC-Light.ttf")

	//// 设置文档保护
	//pdf.SetProtection(gofpdf.CnProtectCopy|gofpdf.CnProtectModify, "userpass", "ownerpass")

	// 设置自动分页
	pdf.SetAutoPageBreak(true, 20) // 20mm 的底部边距

	// 设置页脚函数 - 这会自动在每个页面调用
	pdf.SetFooterFunc(func() {
		// 定位到页面底部
		pdf.SetY(-15)
		pdf.SetFont("sourcehan", "B", 8)
		// 居中显示页码
		pdf.CellFormat(0, 10, "第 "+strconv.Itoa(pdf.PageNo())+" 页",
			"", 0, "C", false, 0, "")
	})

	//pdf.RegisterImageOptions("sk.png", gofpdf.ImageOptions{ReadDpi: true})
	//pdf.ImageOptions("sk.png",
	//	10, 10, // X=0左边开始放  Y=20
	//	200, 0,
	//	true, // 宽=页面宽度，高度=0（自动按比例）
	//	gofpdf.ImageOptions{
	//		ImageType: "PNG",
	//		ReadDpi:   true, // 读取 DPI，确保清晰度正确
	//	},
	//	0,
	//	"",
	//)

	//w, _ := pdf.GetImageInfo("sk.png").Width(), pdf.GetImageInfo("sk.png").Height()

	// 添加第一页
	pdf.AddPage()
	pdf.SetFont("sourcehan", "B", 30)
	pdf.Cell(0, opdf.CellHeight(30), "大乐透开奖分析")
	pdf.Bookmark("大乐透开奖分析", 0, 0)
	pdf.Ln(opdf.LineHeight(30))

	//pdf.AddPage()
	ana_dlt.CHongHaoForPDF(pdf, "sourcehan", 10, "重复开出号码", 1, [5]int{100, 0, 0, 0, 0}, false)
	pdf.AddPage()
	ana_dlt.DltBackQuShiForPDF(pdf, "sourcehan", 10, "后区趋势", 1)
	pdf.SetFont("sourcehan", "B", 30)
	pdf.Cell(0, opdf.CellHeight(30), "欢迎赞助，感谢您的赞助！")
	pdf.Bookmark("欢迎赞助，感谢您的赞助！", 1, 0)
	pdf.Ln(opdf.LineHeight(30))
	pdf.Image("sk.png", 5, 60, 200, 0, false, "", 0, "")
	//
	//pdf.SetFont("sourcehan", "B", 20)
	//pdf.Cell(0, 20, "第一章 概述")
	//pdf.Bookmark("第一章 概述", 1, 0)
	//pdf.Ln(20)

	//pdf.SetFont("sourcehan", "L", 11)
	//// 添加大量内容以触发自动分页
	//for i := 1; i <= 60; i++ {
	//	pdf.Cell(0, 11, "这是第 "+strconv.Itoa(i)+" 行内容，用于演示自动分页功能。")
	//	pdf.Ln(5) // 换行
	//}
	//
	//pdf.AddPageFormat("L", gofpdf.SizeType{Wd: 210, Ht: 297}) // "L" 表示横向
	//pdf.SetFont("sourcehan", "B", 20)
	//pdf.Cell(0, 20, "第二章 重号")
	//pdf.Bookmark("第二章 重号", 1, 0)
	//pdf.Ln(20)
	//
	//pdf.SetFont("sourcehan", "L", 11)
	//// 添加大量内容以触发自动分页
	//for i := 1; i <= 60; i++ {
	//	pdf.Cell(0, 11, "这是第 "+strconv.Itoa(i)+" 行内容，用于演示自动分页功能。")
	//	pdf.Ln(5) // 换行
	//}
	lastDlt := dbop.GetLastDlt()
	// 保存 PDF
	err := pdf.OutputFileAndClose(fmt.Sprintf("大乐透开奖分析—截止至%s的数据.pdf", lastDlt.DrawTime))
	if err != nil {
		panic(err)
	}
}
