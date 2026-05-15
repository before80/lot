package arg

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetThreadNum 获取命令行中输入的线程数
func GetThreadNum() int {
	reader := bufio.NewReader(os.Stdin) // 使用 bufio 读取输入
labelForInputThread:
	fmt.Print("请输入需要启用的线程数1～60（按回车键进行确认）: ")
	input, _ := reader.ReadString('\n') // 读取整行输入
	input = strings.TrimSpace(input)

	// 转换输入为整数
	threads, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("\n输入错误，请确保输入的是有效整数！")
		goto labelForInputThread
	}

	if threads < 1 || threads > 60 {
		fmt.Println("\n请确保输入的整数范围在1～60之间！")
		goto labelForInputThread
	}

	return threads
}
