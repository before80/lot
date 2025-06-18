package dlt

import (
	"fmt"
	"github.com/before80/lot/db"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/sport"
	"github.com/spf13/cobra"
	"time"
)

func UpdateDlt(cmd *cobra.Command) {
	var err error
	// 从命令行参数中获取ldn的值
	ldn, err := cmd.Flags().GetString("ldn")
	if ldn == "" {
		return
	}
	lg.InfoToFileAndStdOut(fmt.Sprintf("ldn=%s\n", ldn))
	startTime := time.Now()
	lastDlt := dbop.GetLastDlt()
	db.DB.Last(&lastDlt)
	lg.InfoToFileAndStdOut(fmt.Sprintf("当前数据库中最新的一条记录为 %v \n", lastDlt))
	dlts, err := sport.GetSomeDltFromWeb(lastDlt, ldn)
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("从网页获取开奖数据出现错误：%v\n", err), 3)
		return
	}
	insertedRow, err := dbop.InsertDltBatch(dlts, 100)
	if err != nil {
		lg.ErrorToFileAndStdOutWithSleepSecond(fmt.Sprintf("往数据表中插入数据出现错误：%v\n", err), 3)
		return
	} else {
		lg.InfoToFileAndStdOut(fmt.Sprintf("插入了 %d 条数据\n", insertedRow))
		lg.InfoToFileAndStdOut(fmt.Sprintf("程序运行时间：%.2f秒\n", time.Since(startTime).Seconds()))
	}
}
