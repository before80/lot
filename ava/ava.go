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
	ana_dlt.DltDataToExcel(false)

	// 预先执行模拟,再整理大乐透相关数据到Excel表中
	//ana_dlt.DltDataToExcel(true)

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
