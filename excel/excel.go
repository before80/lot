package excel

import (
	"fmt"

	"github.com/before80/lot/lg"
	"github.com/xuri/excelize/v2"
)

func CreateNewExcelFile(fileName string) (f *excelize.File, err error) {
	newExcelFileNameWithSuffix := fileName + ".xlsx"
	lg.InfoToFile(fmt.Sprintf("正在创建数据表：%s\n", newExcelFileNameWithSuffix))

	f = excelize.NewFile()

	// 根据指定路径保存文件
	if err = f.SaveAs(newExcelFileNameWithSuffix); err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建后保存excel文件出现错误：%v", err))
		return nil, err
	}

	lg.InfoToFile(fmt.Sprintf("已创建数据表：%s\n", newExcelFileNameWithSuffix))

	return f, err
}

var column = []string{
	"A", "B", "C", "D", "E", "F",
	"G", "H", "I", "J", "K", "L",
	"M", "N", "O", "P", "Q", "R",
	"S", "T", "U", "V", "W", "X", "Y", "Z",
}

// GetColumnStr 获取单元格所在字母编号
func GetColumnStr(index int) (columnStr string) {
	index = index - 1
	for index >= 0 {
		columnStr = column[index%26] + columnStr // 取模并拼接当前字母
		index = index/26 - 1                     // 更新索引，注意减 1
	}
	return
}

func GetColumnStr2(index int) (columnStr string) {
	if index <= 25 {
		columnStr = column[index]
	} else if index <= 676 {
		first := index / 26
		second := index % 26
		columnStr = column[first-1] + column[second]
	} else if index <= 17576 {
		first := index / 676
		second := index % 676 / 26
		third := index % 676 % 26
		columnStr = column[first-1] + column[second-1] + column[third]
	} else if index <= 456976 {
		first := index / 17576
		second := index % 17576 / 676
		third := index % 17576 % 676 / 26
		fourth := index % 17576 % 676 % 26
		columnStr = column[first-1] + column[second-1] + column[third-1] + column[fourth]
	}
	return
}
