package models

import "time"

type Ip struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Ip        string `gorm:"column:ip;size:60;not null"`
	CreatedAt time.Time
}
