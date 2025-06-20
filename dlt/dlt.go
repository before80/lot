package dlt

import (
	"fmt"
	"github.com/before80/lot/db"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/sport"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UpdateDlt(cmd *cobra.Command) {
	var err error
	// 从命令行参数中获取ldn的值
	ldn, err := cmd.Flags().GetString("ldn")
	lg.InfoToFileAndStdOut(fmt.Sprintf("ldn=%s\n", ldn))
	if ldn == "" {
		return
	}

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
		lg.InfoToFileAndStdOut(fmt.Sprintf("处理pdf中...\n"))
		dealPdfNum := 0
		for _, dlt := range dlts {
			if dlt.DrawNum >= "19081" {
				err = DownloadPDF(dlt.DrawPdfUrl)
				if err != nil {
					lg.ErrorToFile(fmt.Sprintf("下载 %s 所在PDF出现错误：%v\n", dlt.DrawPdfUrl, err))
				}
				dealPdfNum++
			}
		}

		lg.InfoToFileAndStdOut(fmt.Sprintf("成功下载了%d个pdf文件\n", dealPdfNum))
		lg.InfoToFileAndStdOut(fmt.Sprintf("程序运行时间：%.2f秒\n", time.Since(startTime).Seconds()))
	}
}

func DownloadPDF(url string) error {
	// 创建HTTP客户端
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头，模拟浏览器访问
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 从URL中提取文件名
	filename := filepath.Base(url)

	// 如果提取的文件名不包含.pdf后缀，则添加后缀
	if !hasPDFExtension(filename) {
		filename += ".pdf"
	}

	_ = os.MkdirAll("download", os.ModePerm)
	// 创建文件
	out, err := os.Create("download/" + filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	// 下载文件内容并写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	//fmt.Printf("文件已成功下载并保存为: %s\n", filename)
	return nil
}

// 检查文件名是否有.pdf后缀
func hasPDFExtension(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".pdf" || ext == ".PDF"
}
