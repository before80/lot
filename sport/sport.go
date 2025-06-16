package sport

import (
	"encoding/json"
	"fmt"
	"github.com/before80/lot/models"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type PrizeLevelListItem1 struct {
	AwardType         int    `json:"awardType"` //
	Group             string `json:"group"`     //
	LotteryCondition  string `json:"lotteryCondition"`
	PrizeLevel        string `json:"prizeLevel"`
	Sort              int    `json:"sort"`
	StakeAmount       string `json:"stakeAmount"`
	StakeAmountFormat string `json:"stakeAmountFormat"`
	StakeCount        string `json:"stakeCount"`
	TotalPrizeAmount  string `json:"totalPrizeamount"`
}

type PrizeLevelListItem2 struct {
	AwardType int `json:"awardType"`
	//Group     string `json:"group"`
	//LotteryCondition  string `json:"lotteryCondition"`
	PrizeLevel        string `json:"prizeLevel"`
	Sort              int    `json:"sort"`
	StakeAmount       string `json:"stakeAmount"`
	StakeAmountFormat string `json:"stakeAmountFormat"`
	StakeCount        string `json:"stakeCount"`
	TotalPrizeAmount  string `json:"totalPrizeamount"`
}

type LastPoolDraw struct {
	LotteryDrawNum       string                `json:"lotteryDrawNum"`
	LotteryDrawResult    string                `json:"lotteryDrawResult"`
	LotteryDrawTime      string                `json:"lotteryDrawTime"`
	LotteryGameNum       string                `json:"lotteryGameNum"`
	PoolBalanceAfterDraw string                `json:"poolBalanceAfterdraw"`
	PrizeLevelList       []PrizeLevelListItem1 `json:"prizeLevelList"`
}

// MyString 自定义类型处理空对象
type MyString string

// UnmarshalJSON 实现json.Unmarshaler接口
func (e *MyString) UnmarshalJSON(data []byte) error {
	// 如果数据是一个空对象`{}`，将其视为空字符串
	if string(data) == "{}" {
		*e = ""
		return nil
	}

	// 尝试解析为普通字符串
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("invalid string format: %w", err)
	}
	*e = MyString(s)
	return nil
}

type ListItem struct {
	LotteryGameName      string `json:"lotteryGameName"`      // 彩种名称
	LotteryGameNum       string `json:"lotteryGameNum"`       // 彩种编号
	LotteryDrawNum       string `json:"lotteryDrawNum"`       // 期号
	LotteryDrawResult    string `json:"lotteryDrawResult"`    // 开奖号码
	LotterySuspendedFlag int    `json:"lotterySuspendedFlag"` // 停售标志
	LotteryDrawStatus    int    `json:"lotteryDrawStatus"`    // 开奖状态
	LotterySaleEndTime   string `json:"lotterySaleEndtime"`   // 截止销售时间
	LotterySaleBeginTime string `json:"lotterySaleBeginTime"` // 开始销售时间
	//LotterySaleEndTimeUnix  string               `json:"lotterySaleEndTimeUnix"` // 截止销售时间
	LotteryDrawTime         string   `json:"lotteryDrawTime"`         // 开奖时间
	LotteryPaidBeginTime    string   `json:"lotteryPaidBeginTime"`    // 派奖开始时间
	LotteryPaidEndTime      string   `json:"lotteryPaidEndTime"`      // 派奖结束时间
	EstimateDrawTime        string   `json:"estimateDrawTime"`        // 预计开奖时间
	Verify                  int      `json:"verify"`                  // 是否验证
	LotteryPromotionFlag    int      `json:"lotteryPromotionFlag"`    // 是否是优惠彩种
	IsGetKjPdf              int      `json:"isGetKjpdf"`              // 是否获取开奖PDF
	IsGetXlPdf              int      `json:"isGetXlpdf"`              // 是否获取销售记录PDF
	PdfType                 int      `json:"pdfType"`                 // PDF类型
	LotteryUnSortDrawResult MyString `json:"lotteryUnsortDrawresult"` // 未排序开奖号码
	PoolBalanceAfterDraw    string   `json:"poolBalanceAfterdraw"`    // 奖池余额
	PoolBalanceAfterDrawRj  string   `json:"poolBalanceAfterdrawRj"`  // 奖池余额
	DrawFlowFund            string   `json:"drawFlowFund"`            // 本期资金
	DrawFlowFundRj          string   `json:"drawFlowFundRj"`          // 本期资金
	TotalSaleAmount         string   `json:"totalSaleAmount"`         // 本期销售金额
	TotalSaleAmountRj       string   `json:"totalSaleAmountRj"`       // 本期销售金额
	LotteryEquipmentCount   int      `json:"lotteryEquipmentCount"`   // 彩机数量
	LotteryGameProNum       int      `json:"lotteryGamePronum"`       // 彩种编号
	//MatchList               string               `json:"matchList"`
	PrizeLevelList []PrizeLevelListItem2 `json:"prizeLevelList"` // 奖级列表
	//TermList                string               `json:"termList"`
	//TermResultList          string               `json:"termResultList"`
	RuleType int `json:"ruleType"` // 规则类型
	//VToolsConfig            string               `json:"vtoolsConfig"`
	SurplusAmount          string `json:"surplusAmount"`          // 剩余金额
	SurplusAmountRj        string `json:"surplusAmountRj"`        // 剩余金额
	LotteryPromotionFlagRj int    `json:"lotteryPromotionFlagRj"` // 是否是优惠彩种
	DrawPdfUrl             string `json:"drawPdfUrl"`             // PDF地址
	LotteryNotice          int    `json:"lotteryNotice"`          // 是否有公告
	LotteryDrawStatusNo    string `json:"lotteryDrawStatusNo"`    // 开奖状态
	LotteryNoticeShowFlag  int    `json:"lotteryNoticeShowFlag"`  // 是否显示公告
}

type Value struct {
	LastPoolDraw LastPoolDraw `json:"lastPoolDraw"` // 最近一期开奖信息
	List         []ListItem   `json:"list"`         // 历史开奖列表
	PageNo       int          `json:"pageNo"`       // 页码
	PageSize     int          `json:"pageSize"`     // 每页条数
	Pages        int          `json:"pages"`        // 总页数
	Total        int          `json:"total"`        // 总条数
}

type LotterySt struct {
	DataFrom     string `json:"dataFrom"` // 数据来源
	EmptyFlag    bool   `json:"emptyFlag"`
	ErrorCode    string `json:"errorCode"`    // 错误码
	ErrorMessage string `json:"errorMessage"` // 错误信息
	Success      bool   `json:"success"`      // 是否成功
	Value        Value  `json:"value"`
}

func judgeCrossYearGt(startTerm, endTerm string, gtNum int) bool {
	startYear, _ := strconv.Atoi(startTerm[0:2])
	endYear, _ := strconv.Atoi(endTerm[0:2])
	if endYear-startYear > gtNum {
		return true
	}
	return false
}

func GenCrossTermYears(startTerm, endTerm string, gtNum int) (crossTermYears [][2]string) {
	curEndTerm := startTerm
	curStartTerm := startTerm
	var curEndTermNum int
	if judgeCrossYearGt(startTerm, endTerm, gtNum) {
	LabelForContinue0:
		curStartYear, _ := strconv.Atoi(curEndTerm[0:2])
		curEndTerm = strconv.Itoa(curStartYear+gtNum) + "001"
		curEndTermNum, _ = strconv.Atoi(curEndTerm)
		if curEndTerm <= endTerm {
			crossTermYears = append(crossTermYears, [][2]string{
				{curStartTerm, curEndTerm},
			}...)
			curEndTermNum += 1
			curStartTerm = strconv.Itoa(curEndTermNum)
			// 补0
			if len(curStartTerm) < len(startTerm) {
				curStartTerm = "0" + curStartTerm
			}
		} else {
			crossTermYears = append(crossTermYears, [][2]string{
				{curStartTerm, endTerm},
			}...)
		}

		if curEndTerm <= endTerm {
			if judgeCrossYearGt(curEndTerm, endTerm, gtNum) {
				goto LabelForContinue0
			} else {
				crossTermYears = append(crossTermYears, [][2]string{
					{curStartTerm, endTerm},
				}...)
			}
		}
	} else {
		crossTermYears = [][2]string{
			{startTerm, endTerm},
		}
	}
	return
}

func GetLotteryHistory(gameNo, pageSize, pageNo, intervalSecond int, startTerm, endTerm string) (lotteryHistory []ListItem) {
	crossTermYears := GenCrossTermYears(startTerm, endTerm, 5)
	for _, crossTermYear := range crossTermYears {
		// 组合返回的结果
		curPageNo := pageNo
		lotterySt := getLotteryHistory(gameNo, pageSize, pageNo, crossTermYear[0], crossTermYear[1])
		lotteryHistory = append(lotteryHistory, lotterySt.Value.List...)
		if curPageNo <= lotterySt.Value.Pages {
		LabelForContinue1:
			curPageNo += 1
			for curPageNo <= lotterySt.Value.Pages {
				time.Sleep(time.Duration(intervalSecond) * time.Second)
				lotterySt = getLotteryHistory(gameNo, pageSize, curPageNo, crossTermYear[0], crossTermYear[1])
				lotteryHistory = append(lotteryHistory, lotterySt.Value.List...)
				goto LabelForContinue1
			}
		}
	}

	// 排序
	slices.SortFunc(lotteryHistory, func(a, b ListItem) int {
		aTerm, _ := strconv.Atoi(a.LotteryDrawNum)
		bTerm, _ := strconv.Atoi(b.LotteryDrawNum)
		if aTerm > bTerm {
			return 1
		} else if aTerm < bTerm {
			return -1
		}
		return 0
	})

	return
}

func getLotteryHistory(gameNo, pageSize, pageNo int, startTerm, endTerm string) (lotteryData LotterySt) {
	client := &http.Client{
		Timeout: 20 * time.Second,
		//Transport: &http.Transport{
		//	MaxIdleConns: 100,
		//},
	}

	url := fmt.Sprintf("https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=%d&provinceId=0&pageSize=%d&isVerify=1&pageNo=%d&startTerm=%s&endTerm=%s", gameNo, pageSize, pageNo, startTerm, endTerm)
	// 创建HTTP GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Sprintf("创建请求失败: %v\n", err))
	}

	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	////req.Header.Set("Access-Control-Request-Method", "GET")
	//req.Header.Set("Access-Control-Request-Headers", "x-dev-fp,x-full-ref,x-sensors-id")
	//////req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	//req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	//req.Header.Set("Priority", "u=1, i")
	//req.Header.Set("Origin", "https://static.sporttery.cn")
	//req.Header.Set("Referer", "https://static.sporttery.cn/")
	//req.Header.Set("Sec-Ch-Ua-Platform", "Windows")
	//req.Header.Set("Sec-Fetch-Dest", "")
	//req.Header.Set("Sec-Fetch-Mode", "cors")
	//req.Header.Set("Sec-Fetch-Site", "same-site")
	//req.Header.Set("x-dev-fp", "8e26f46ee3eceeac0b38f819dcce5fc7")
	//req.Header.Set("x-full-ref", "//www.lottery.gov.cn/kj/kjlb.html?dlt")
	//req.Header.Set("x-sensors-id", "196bd9a41b7842-064d65a09f91dfc-26011f51-3686400-196bd9a41b9136b")

	//req.Header.Set(":authority", "webapi.sporttery.cn")
	//req.Header.Set(":method", "OPTIONS")
	//req.Header.Set(":path", fmt.Sprintf("/gateway/lottery/getHistoryPageListV1.qry?gameNo=85&provinceId=0&pageSize=100&isVerify=1&pageNo=%d", pageNo))
	//req.Header.Set(":scheme", "https")
	//req.Header.Set(":authority", "webapi.sporttery.cn")
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("打不开%s，遇到错误：%v\n", url, err))
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		// 处理非200响应（如401、500等）
		panic(fmt.Sprintf("打开%s遇到响应状态码：%d\n", url, resp.StatusCode))
	}

	body, _ := io.ReadAll(resp.Body)
	//fmt.Printf("%v\n", string(body))

	err = json.Unmarshal(body, &lotteryData)
	if err != nil {
		panic(fmt.Sprintf("解析%s返回结果，遇到错误：%v\n", url, err))
	}
	fmt.Printf("已完成第%d页（%d） %s~%s\n", pageNo, lotteryData.Value.Pages, startTerm, endTerm)
	return
}

func GetStake(pls []PrizeLevelListItem2, sort int) (stakeCount int, stakeAmount int) {
	var parseErr error
	for _, pl := range pls {
		if pl.Sort == sort {
			stakeCount, parseErr = strconv.Atoi(strings.ReplaceAll(pl.StakeCount, ",", ""))
			if parseErr != nil || stakeCount < 0 {
				stakeCount = 0
			}
			stakeAmount, parseErr = strconv.Atoi(pl.StakeAmountFormat)
			if parseErr != nil || stakeAmount < 0 {
				stakeAmount = 0
			}
			return
		}
	}
	return
}

func GetSomeDltFromWeb(lastDlt models.Dlt, endTerm string) (dlts []models.Dlt, err error) {
	if lastDlt.DrawNum != endTerm {
		startDrawNum := "07001"
		if lastDlt.DrawNum != "" {
			startDrawNum = lastDlt.DrawNum
		}
		ld := GetLotteryHistory(85, 100, 1, 10, startDrawNum, endTerm)
		fmt.Printf("%v\n", len(ld))
		//fmt.Printf("%v\n", ld)
		//fmt.Printf("%v\n", GenCrossTermYears("07001", "25053", 5))
		var parseErr error
		if len(ld) > 0 {
			for _, v := range ld {
				if lastDlt.DrawNum != "" && v.LotteryDrawNum == lastDlt.DrawNum {
					continue
				}
				newDlt := models.Dlt{}
				newDlt.DrawNum = v.LotteryDrawNum
				newDlt.DrawTime = v.LotteryDrawTime
				newDlt.EquipmentCount = v.LotteryEquipmentCount
				newDlt.DrawPdfUrl = v.DrawPdfUrl
				newDlt.UnSortDrawResult = string(v.LotteryUnSortDrawResult)
				lotteryDrawResults := strings.Split(v.LotteryDrawResult, " ")
				if len(lotteryDrawResults) == 7 {
					newDlt.F1 = lotteryDrawResults[0]
					newDlt.F2 = lotteryDrawResults[1]
					newDlt.F3 = lotteryDrawResults[2]
					newDlt.F4 = lotteryDrawResults[3]
					newDlt.F5 = lotteryDrawResults[4]
					newDlt.B1 = lotteryDrawResults[5]
					newDlt.B2 = lotteryDrawResults[6]
				}

				newDlt.PoolBalance, parseErr = strconv.ParseFloat(strings.ReplaceAll(v.PoolBalanceAfterDraw, ",", ""), 64)
				if parseErr != nil {
					newDlt.PoolBalance = 0
				}
				newDlt.TotalSaleAmount, parseErr = strconv.ParseFloat(strings.ReplaceAll(v.TotalSaleAmount, ",", ""), 64)
				if parseErr != nil {
					newDlt.TotalSaleAmount = 0
				}

				newDlt.StakeCount60, newDlt.StakeAmount60 = GetStake(v.PrizeLevelList, 60)
				newDlt.StakeCount80, newDlt.StakeAmount80 = GetStake(v.PrizeLevelList, 80)
				newDlt.StakeCount100, newDlt.StakeAmount100 = GetStake(v.PrizeLevelList, 100)
				newDlt.StakeCount101, newDlt.StakeAmount101 = GetStake(v.PrizeLevelList, 101)
				newDlt.StakeCount102, newDlt.StakeAmount102 = GetStake(v.PrizeLevelList, 102)
				newDlt.StakeCount201, newDlt.StakeAmount201 = GetStake(v.PrizeLevelList, 201)
				newDlt.StakeCount202, newDlt.StakeAmount202 = GetStake(v.PrizeLevelList, 202)
				newDlt.StakeCount301, newDlt.StakeAmount301 = GetStake(v.PrizeLevelList, 301)
				newDlt.StakeCount302, newDlt.StakeAmount302 = GetStake(v.PrizeLevelList, 302)
				newDlt.StakeCount401, newDlt.StakeAmount401 = GetStake(v.PrizeLevelList, 401)
				newDlt.StakeCount402, newDlt.StakeAmount402 = GetStake(v.PrizeLevelList, 402)
				newDlt.StakeCount501, newDlt.StakeAmount501 = GetStake(v.PrizeLevelList, 501)
				newDlt.StakeCount601, newDlt.StakeAmount601 = GetStake(v.PrizeLevelList, 601)
				newDlt.StakeCount701, newDlt.StakeAmount701 = GetStake(v.PrizeLevelList, 701)
				newDlt.StakeCount801, newDlt.StakeAmount801 = GetStake(v.PrizeLevelList, 801)
				newDlt.StakeCount901, newDlt.StakeAmount901 = GetStake(v.PrizeLevelList, 901)
				newDlt.StakeCount1001, newDlt.StakeAmount1001 = GetStake(v.PrizeLevelList, 1001)
				newDlt.StakeCount1101, newDlt.StakeAmount1101 = GetStake(v.PrizeLevelList, 1101)

				// 追加到切片
				dlts = append(dlts, newDlt)
			}
		}
	}
	return
}
