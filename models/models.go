package models

import "time"

type Dlt struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	DrawNum          string    `json:"draw_num"`
	DrawTime         string    `json:"draw_time"`
	EquipmentCount   int       `json:"equipment_count"`
	DrawPdfUrl       string    `json:"draw_pdf_url"`
	UnSortDrawResult string    `gorm:"column:unsort_draw_result" json:"unsort_draw_result"`
	F1               string    `json:"f1"`
	F2               string    `json:"f2"`
	F3               string    `json:"f3"`
	F4               string    `json:"f4"`
	F5               string    `json:"f5"`
	B1               string    `json:"b1"`
	B2               string    `json:"b2"`
	PoolBalance      float64   `gorm:"default:0" json:"pool_balance"`
	TotalSaleAmount  float64   `gorm:"default:0" json:"total_sale_amount"`
	StakeCount60     int       `gorm:"default:0;column:stake_count_60" json:"stake_count_60"`
	StakeCount80     int       `gorm:"default:0;column:stake_count_80" json:"stake_count_80"`
	StakeCount100    int       `gorm:"default:0;column:stake_count_100" json:"stake_count_100"`
	StakeCount101    int       `gorm:"default:0;column:stake_count_101" json:"stake_count_101"`
	StakeCount102    int       `gorm:"default:0;column:stake_count_102" json:"stake_count_102"`
	StakeCount201    int       `gorm:"default:0;column:stake_count_201" json:"stake_count_201"`
	StakeCount202    int       `gorm:"default:0;column:stake_count_202" json:"stake_count_202"`
	StakeCount301    int       `gorm:"default:0;column:stake_count_301" json:"stake_count_301"`
	StakeCount302    int       `gorm:"default:0;column:stake_count_302" json:"stake_count_302"`
	StakeCount401    int       `gorm:"default:0;column:stake_count_401" json:"stake_count_401"`
	StakeCount402    int       `gorm:"default:0;column:stake_count_402" json:"stake_count_402"`
	StakeCount501    int       `gorm:"default:0;column:stake_count_501" json:"stake_count_501"`
	StakeCount601    int       `gorm:"default:0;column:stake_count_601" json:"stake_count_601"`
	StakeCount701    int       `gorm:"default:0;column:stake_count_701" json:"stake_count_701"`
	StakeCount801    int       `gorm:"default:0;column:stake_count_801" json:"stake_count_801"`
	StakeCount901    int       `gorm:"default:0;column:stake_count_901" json:"stake_count_901"`
	StakeCount1001   int       `gorm:"default:0;column:stake_count_1001" json:"stake_count_1001"`
	StakeCount1101   int       `gorm:"default:0;column:stake_count_1101" json:"stake_count_1101"`
	StakeAmount60    int       `gorm:"default:0;column:stake_amount_60" json:"stake_amount_60"`
	StakeAmount80    int       `gorm:"default:0;column:stake_amount_80" json:"stake_amount_80"`
	StakeAmount100   int       `gorm:"default:0;column:stake_amount_100" json:"stake_amount_100"`
	StakeAmount101   int       `gorm:"default:0;column:stake_amount_101" json:"stake_amount_101"`
	StakeAmount102   int       `gorm:"default:0;column:stake_amount_102" json:"stake_amount_102"`
	StakeAmount201   int       `gorm:"default:0;column:stake_amount_201" json:"stake_amount_201"`
	StakeAmount202   int       `gorm:"default:0;column:stake_amount_202" json:"stake_amount_202"`
	StakeAmount301   int       `gorm:"default:0;column:stake_amount_301" json:"stake_amount_301"`
	StakeAmount302   int       `gorm:"default:0;column:stake_amount_302" json:"stake_amount_302"`
	StakeAmount401   int       `gorm:"default:0;column:stake_amount_401" json:"stake_amount_401"`
	StakeAmount402   int       `gorm:"default:0;column:stake_amount_402" json:"stake_amount_402"`
	StakeAmount501   int       `gorm:"default:0;column:stake_amount_501" json:"stake_amount_501"`
	StakeAmount601   int       `gorm:"default:0;column:stake_amount_601" json:"stake_amount_601"`
	StakeAmount701   int       `gorm:"default:0;column:stake_amount_701" json:"stake_amount_701"`
	StakeAmount801   int       `gorm:"default:0;column:stake_amount_801" json:"stake_amount_801"`
	StakeAmount901   int       `gorm:"default:0;column:stake_amount_901" json:"stake_amount_901"`
	StakeAmount1001  int       `gorm:"default:0;column:stake_amount_1001" json:"stake_amount_1001"`
	StakeAmount1101  int       `gorm:"default:0;column:stake_amount_1101" json:"stake_amount_1101"`
	CreatedAt        time.Time `json:"created_at"`
}
