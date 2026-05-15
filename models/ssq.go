package models

import "time"

type Ssq struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	DrawNum         string    `json:"draw_num"`
	DrawTime        string    `json:"draw_time"`
	Week            string    `json:"week"`
	VideoUrl        string    `json:"video_url"`
	DetailsUrl      string    `json:"details_url"`
	Content         string    `json:"content"`
	F1              string    `json:"f1"`
	F2              string    `json:"f2"`
	F3              string    `json:"f3"`
	F4              string    `json:"f4"`
	F5              string    `json:"f5"`
	F6              string    `json:"f6"`
	B1              string    `json:"b1"`
	Oe              string    `json:"oe"`
	Hz              int       `json:"hz"`
	AeHz            string    `json:"ae_hz"`
	Qzh             string    `json:"qzh"`
	PoolBalance     float64   `gorm:"default:0" json:"pool_balance"`
	TotalSaleAmount float64   `gorm:"default:0" json:"total_sale_amount"`
	StakeCount1     int       `gorm:"default:0;column:stake_count_1" json:"stake_count_1"`
	StakeCount2     int       `gorm:"default:0;column:stake_count_2" json:"stake_count_2"`
	StakeCount3     int       `gorm:"default:0;column:stake_count_3" json:"stake_count_3"`
	StakeCount4     int       `gorm:"default:0;column:stake_count_4" json:"stake_count_4"`
	StakeCount5     int       `gorm:"default:0;column:stake_count_5" json:"stake_count_5"`
	StakeCount6     int       `gorm:"default:0;column:stake_count_6" json:"stake_count_6"`
	StakeAmount1    int       `gorm:"default:0;column:stake_amount_1" json:"stake_amount_1"`
	StakeAmount2    int       `gorm:"default:0;column:stake_amount_2" json:"stake_amount_2"`
	StakeAmount3    int       `gorm:"default:0;column:stake_amount_3" json:"stake_amount_3"`
	StakeAmount4    int       `gorm:"default:0;column:stake_amount_4" json:"stake_amount_4"`
	StakeAmount5    int       `gorm:"default:0;column:stake_amount_5" json:"stake_amount_5"`
	StakeAmount6    int       `gorm:"default:0;column:stake_amount_6" json:"stake_amount_6"`
	CreatedAt       time.Time `json:"created_at"`
}

type AllSsq struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Hm        string    `json:"hm"`
	Oe        string    `json:"oe"`
	Hz        int       `json:"hz"`
	AeHz      string    `json:"ae_hz"`
	Qzh       string    `json:"qzh"`
	CreatedAt time.Time `json:"created_at"`
}

type SsqMoni struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	A         string    `json:"a"`
	B         string    `json:"b"`
	C         string    `json:"c"`
	D         string    `json:"d"`
	E         string    `json:"e"`
	F         string    `json:"f"`
	Cs        int       `json:"cs"`
	Comb      int       `json:"comb"`
	Typ       string    `json:"typ"`
	Method    string    `json:"method"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
