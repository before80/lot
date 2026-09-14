package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/before80/lot/ana_dlt"
)

func main() {

	// (从500.com网页中)获取大乐透开奖数据, 并批量插入到数据表中
	//ana_dlt.GetDltKj()

	// | ------ 以下只能在本地电脑执行, 服务器不能使用go-rod打开浏览器-----------|
	// 从官网获取最新的大乐透开奖数据并批量更新数据表
	//dlt.UpdateDlt1()
	// | ------ 以上只能在本地电脑执行, 服务器不能使用go-rod打开浏览器-----------|

	// 不预先执行模拟,再整理大乐透相关数据到Excel表中
	//ana_dlt.DltDataToExcel(false)

	// 预先执行模拟,再整理大乐透相关数据到Excel表中
	// ana_dlt.DltDataToExcel(true)

	// 不预先执行模拟,再整理大乐透更多的相关数据到Excel表中
	// go build -o dlt_excel_more_without_moni_new.exe .
	//ana_dlt.DltMoreDataToExcel(false)

	// 预先执行模拟,再整理大乐透更多的相关数据到Excel表中
	// go build -o dlt_excel_more_with_moni_new.exe .
	ana_dlt.DltMoreDataToExcel(true)

	// 前区跨期数统计
	//ana_dlt.DltFrontStatDataToExcel([]int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30}, "2_30", 0)
	//ana_dlt.DltFrontStatDataToExcel([]int{10, 15, 20, 30, 40, 50, 60, 70, 80, 90, 100}, "10_100", 0)
	//
	//// 某个设备的前区跨期数统计
	//ana_dlt.DltFrontStatDataToExcel([]int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30}, "2_30", 1)
	//ana_dlt.DltFrontStatDataToExcel([]int{10, 15, 20, 30, 40, 50, 60, 70, 80, 90, 100}, "10_100", 1)
	//
	//ana_dlt.DltFrontStatDataToExcel([]int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30}, "2_30", 2)
	//ana_dlt.DltFrontStatDataToExcel([]int{10, 15, 20, 30, 40, 50, 60, 70, 80, 90, 100}, "10_100", 2)
	//
	//ana_dlt.DltFrontStatDataToExcel([]int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30}, "2_30", 3)
	//ana_dlt.DltFrontStatDataToExcel([]int{10, 15, 20, 30, 40, 50, 60, 70, 80, 90, 100}, "10_100", 3)
	//
	//ana_dlt.DltNextFrontStatDataToExcel("", []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	//	21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
	//	41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60,
	//	61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 78, 79, 80,
	//	81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100})

	// go build -o dlt_excel_next_spDrawNum_front_stat.exe .
	//ana_dlt.DltNextFrontStatDataToExcel(GetSpDrawNumString(), 2, 200, 50)

	// 计算 某个模式下的某种规则, 与 指定 Tx 类型的前区号码的交集
	//
	//	input := `6*[20]
	//5*[03 09 15 26]
	//4*[04 28 33]
	//3*[10 17 19 21 23 24 29 35]
	//2*[01 02 05 07 14 16 18 25 27 31 32]
	//1*[06 11 13 22 30 34]
	//0*[08 12]
	//`
	//
	//	ana_dlt.CalTxFrontHmStr2StMap()
	//
	//	res, _ := ana_dlt.GenerateCombinationsByInputStrAndRule(input, "2111")
	//	var jjs [][]string
	//	for _, r := range res {
	//		hms := fmt.Sprintf("%s,%s,%s,%s,%s", r[0], r[1], r[2], r[3], r[4])
	//		if _, ok := ana_dlt.FrontHmStr2StMap["OtherT"][hms]; ok {
	//			jjs = append(jjs, r)
	//		}
	//	}
	//
	//	fmt.Printf("res=%d\n", len(res))
	//	fmt.Printf("ana_dlt.FrontHmStr2StMap[\"OtherT\"]=%d\n", len(ana_dlt.FrontHmStr2StMap["OtherT"]))
	//	fmt.Printf("result=%d\n", len(jjs))

	//res := ana_dlt.Generate21111()
	//fmt.Println(len(res))
	//fmt.Println(res)

	//	input := `12*[13]=>[]
	//10*[09]=>[]
	//8*[06 21 22 26]=>[]
	//7*[02 03 08 10 23]=>[]
	//6*[04 12 14 16 18 28]=>[]
	//5*[01 05 07 11 27 29 32 33 34 35]=>[]
	//4*[24]=>[]
	//3*[17 19 30 31]=>[]
	//2*[15]=>[]
	//1*[20 25]=>[]`
	//
	//	for _, rule := range []string{"2111", "221", "32", "311", "11111", "41", "5"} {
	//		// 生成组合
	//		result, err := ana_dlt.GenerateCombinationsByInputStrAndRule(input, rule)
	//		if err != nil {
	//			fmt.Printf("Rule %s generate error: %v\n", rule, err)
	//			continue
	//		}
	//		// 计数
	//		count, err := ana_dlt.CountByInputStringAndRule(input, rule)
	//		if err != nil {
	//			fmt.Printf("Rule %s count error: %v\n", rule, err)
	//			continue
	//		}
	//		fmt.Printf("Rule %s: count = %d, generated = %d\n", rule, count, len(result))
	//		// 打印前2个示例
	//		for i := 0; i < 2 && i < len(result); i++ {
	//			fmt.Println("  ", result[i])
	//		}
	//		fmt.Println()
	//	}

	//ana_dlt.DltFrontDuanStatDataToExcel()

	//ana_dlt.UpdateMonisTable()
	//ana_dlt.GenAllDltsInfoForGoCode()

	//lastDrawNum := GetString()
	//if lastDrawNum == "" {
	//	dlt := dbop.GetLastDlt()
	//	lastDrawNum = dlt.DrawNum
	//}
	//
	//ana_dlt.DltXHao(ana_dlt.NewXuHaoSt(), GetNum(), lastDrawNum)

	//ana_dlt.BatchMoni(true, 9)
	//ana_dlt.UpdateDltHz()
	//ana_ssq.BatchMoni(true, 4)
	//ana_ssq.SsqDataToExcel(true)
	//ana_ssq.SsqDataToExcel(false)
	//fmt.Println(gen.CrossComb([]string{"a", "b", "d"}, []string{"c"}, 1, 1))

	//lastDrawNum := GetString()
	//if lastDrawNum == "" {
	//	dlt := dbop.GetLastSsq()
	//	lastDrawNum = dlt.DrawNum
	//}
	//ana_ssq.SsqXHao(ana_ssq.NewXuHaoSt(), GetNum(), lastDrawNum)
}

func GetNum() int {
	reader := bufio.NewReader(os.Stdin) // 使用 bufio 读取输入
labelForInputNum:
	fmt.Print("请输入需要的注数1～9999999（按回车键进行确认）: ")
	input, _ := reader.ReadString('\n') // 读取整行输入
	input = strings.TrimSpace(input)

	// 转换输入为整数
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("\n输入错误，请确保输入的是有效整数！")
		goto labelForInputNum
	}

	if num < 1 || num > 9999999 {
		fmt.Println("\n请确保输入的整数范围在1～9999999之间！")
		goto labelForInputNum
	}

	return num
}

func GetString() string {
	reader := bufio.NewReader(os.Stdin) // 使用 bufio 读取输入
labelForInputNum:
	fmt.Print("请输入想要研究的最新一期期号,直接回车表示使用当前最新一期的期号: ")
	input, _ := reader.ReadString('\n') // 读取整行输入
	input = strings.TrimSpace(input)

	if len(input) != 0 && len(input) != 5 {
		fmt.Println("\n请确保输入的期号的长度为5,例如: 25001或25139等！")
		goto labelForInputNum
	}

	return input
}

func GetSpDrawNumString() string {
	reader := bufio.NewReader(os.Stdin) // 使用 bufio 读取输入
labelForInputNum:
	fmt.Print("请输入想要研究的那一期的期号,直接回车表示使用当前最新一期的期号: ")
	input, _ := reader.ReadString('\n') // 读取整行输入
	input = strings.TrimSpace(input)

	if len(input) != 0 && len(input) != 5 {
		fmt.Println("\n请确保输入的期号的长度为5,例如: 26001或26052等！")
		goto labelForInputNum
	}

	return input
}
