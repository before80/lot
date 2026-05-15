package ana_dlt

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/before80/lot/dbop"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
)

// GetDltKj (从500.com网页中)获取大乐透开奖数据, 并批量插入到数据表中
func GetDltKj() (executeRes bool) {
	//fmt.Printf("run 1")
	lastDlt := dbop.GetLastDlt()
	startDn := lastDlt.DrawNum
	//startDn := "25140"
	ldn := strconv.Itoa(time.Now().Year())[2:] + "156"

	url := fmt.Sprintf("https://datachart.500.com/dlt/history/newinc/history.php?start=%s&end=%s", startDn, ldn)

	// 创建自定义客户端，设置超时时间
	client := &http.Client{
		Timeout: 30 * time.Second, // 设置30秒超时
	}

	// 创建请求对象
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("创建请求失败:%v", err))
		return
	}

	// 设置请求头，模拟浏览器请求
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("sec-ch-ua", fmt.Sprintf(`"Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"`))
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("Referer", "https://datachart.500.com/dlt/history/history.shtml")

	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("请求网页%s失败:%v", url, err))
		return
	}
	defer res.Body.Close()

	defer res.Body.Close()
	if res.StatusCode != 200 {
		lg.ErrorToFile(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
		return
	}
	//fmt.Printf("run 3")
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		lg.ErrorToFile(fmt.Sprintf("载入网页%s内容 遇到错误:%v", url, err))
		return
	}

	var willInsertDlts []models.Dlt
	//fmt.Printf("run 4")

	doc.Find("#tdata tr").Each(func(i int, s *goquery.Selection) {
		newDlt := models.Dlt{}
		s.Find("td").Each(func(j int, is *goquery.Selection) {
			text := strings.TrimSpace(is.Text())
			if j == 0 {
				newDlt.DrawNum = text
			}
			if j == 1 {
				newDlt.F1 = text
			}
			if j == 2 {
				newDlt.F2 = text
			}
			if j == 3 {
				newDlt.F3 = text
			}
			if j == 4 {
				newDlt.F4 = text
			}
			if j == 5 {
				newDlt.F5 = text
			}
			if j == 6 {
				newDlt.B1 = text
			}
			if j == 7 {
				newDlt.B2 = text
			}
			if j == 8 {
				newDlt.PoolBalance, _ = strconv.ParseFloat(strings.ReplaceAll(text, ",", ""), 64)
			}
			if j == 9 {
				newDlt.StakeCount101, _ = strconv.Atoi(text)
			}
			if j == 10 {
				stakeAmount101, _ := strconv.ParseInt(strings.ReplaceAll(text, ",", ""), 10, 64)
				newDlt.StakeAmount101 = int(stakeAmount101)
			}
			if j == 11 {
				newDlt.StakeCount301, _ = strconv.Atoi(text)
			}
			if j == 12 {
				stakeAmount301, _ := strconv.ParseInt(strings.ReplaceAll(text, ",", ""), 10, 64)
				newDlt.StakeAmount301 = int(stakeAmount301)
			}
			if j == 13 {
				newDlt.TotalSaleAmount, _ = strconv.ParseFloat(strings.ReplaceAll(text, ",", ""), 64)
			}
			if j == 14 {
				newDlt.DrawTime = text
			}
		})

		if newDlt.DrawNum > startDn {
			newDlt.DataSrc = 1 // 设置数据来源 1=500.com
			newDlt.Hz = CalDltHz([]string{newDlt.F1, newDlt.F2, newDlt.F3, newDlt.F4, newDlt.F5, newDlt.B1, newDlt.B2})
			newDlt.Oe = CalDltOe([]string{newDlt.F1, newDlt.F2, newDlt.F3, newDlt.F4, newDlt.F5, newDlt.B1, newDlt.B2})
			newDlt.Qzh = CalDltQzh([]string{newDlt.F1, newDlt.F2, newDlt.F3, newDlt.F4, newDlt.F5})
			newDlt.AeHz = DltHzABCDE(newDlt.Hz)
			willInsertDlts = append(willInsertDlts, newDlt)
		}
	})

	//lg.InfoToFileAndStdOut(fmt.Sprintf("len(willInsertDlts)=%d willInsertDlts = %#v , \n", len(willInsertDlts), willInsertDlts))
	//time.Sleep(9999 * time.Hour)

	if len(willInsertDlts) > 0 {
		insertedRow, err1 := dbop.InsertDltBatch(willInsertDlts, 100)
		if err1 != nil {
			lg.ErrorToFile(fmt.Sprintf("往数据表中插入数据出现错误：%v\n", err))
			return
		} else {
			executeRes = true
			lg.InfoToFile(fmt.Sprintf("插入时间: %v ,其插入了 %d 条数据\n", time.Now(), insertedRow))
		}
	} else {
		executeRes = true
	}
	return
}
