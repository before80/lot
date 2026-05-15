package ssq

type PrizeGrades struct {
	Type      int    `json:"type"`      // 几等奖
	TypeMoney string `json:"typemoney"` // 奖金
	TypeNum   string `json:"typenum"`   // 中奖注数
}
type ListItem struct {
	Code        string        `json:"code"`      // 期号
	Blue        string        `json:"blue"`      // 1个蓝色球号码
	Red         string        `json:"red"`       // 6个红色球号码
	Date        string        `json:"date"`      // 开奖日期
	Week        string        `json:"week"`      // 星期几
	Content     string        `json:"content"`   // 一等奖中奖情况
	PoolMoney   string        `json:"poolmoney"` // 奖池余额
	Sales       string        `json:"sales"`     // 销售额
	DetailsLink string        `json:"detailslink"`
	VideoLink   string        `json:"videolink"`
	PrizeGrades []PrizeGrades `json:"prizegrades"`
	AddMoney    string        `json:"addmoney"`
	AddMoney2   string        `json:"addmoney2"`
	Msg         string        `json:"msg"`
	Name        string        `json:"name"`
	M2Add       string        `json:"m2add"`
	Z2Add       string        `json:"z2add"`
}

type LotterySt struct {
	TFlag    int        `json:"TFlag"` // 数据来源
	Message  string     `json:"message"`
	PageNo   int        `json:"pageNo"`   // 错误码
	PageNum  int        `json:"pageNum"`  // 错误信息
	PageSize int        `json:"pageSize"` // 错误信息
	Result   []ListItem `json:"result"`   // 是否成功
	State    int        `json:"state"`
	Total    int        `json:"total"`
}
