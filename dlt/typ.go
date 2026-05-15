package dlt

import (
	"encoding/json"
	"fmt"
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
