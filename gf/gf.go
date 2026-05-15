package main

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/before80/lot/ana_dlt"
	"github.com/before80/lot/db"
	"github.com/before80/lot/dlt"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
)

func main() {
	// 找出dlts表中 data_src = 1的记录并更新其中的内容

	var dlts []models.Dlt
	db.DB.Where("data_src = ?", 1).Order("id asc").Find(&dlts)

	if len(dlts) > 0 {
		var drawNumSlice []string
		for _, dt := range dlts {
			drawNumSlice = append(drawNumSlice, dt.DrawNum)
		}
		lg.InfoToFileAndStdOut(fmt.Sprintf("drawNumSlice = %v\n", drawNumSlice))
		startDrawNum := dlts[0].DrawNum
		var endDrawNum string
		if len(dlts) == 1 {
			endDrawNum = strconv.Itoa(time.Now().Year())[2:] + "156"
		} else {
			endDrawNum = dlts[len(dlts)-1].DrawNum
		}
		ld := dlt.GetLotteryHistory(5, startDrawNum, endDrawNum)
		var parseErr error
		if len(ld) > 0 {
			// 按照 期号,从小到大
			sort.Slice(ld, func(i, j int) bool {
				return ld[i].LotteryDrawNum < ld[j].LotteryDrawNum
			})
			for _, v := range ld {
				//fmt.Println(v.LotteryDrawNum)
				if slices.Contains(drawNumSlice, v.LotteryDrawNum) {
					newDlt := models.Dlt{}
					newDlt.DrawNum = v.LotteryDrawNum
					newDlt.DrawTime = v.LotteryDrawTime
					newDlt.EquipmentCount = v.LotteryEquipmentCount
					newDlt.DrawPdfUrl = v.DrawPdfUrl
					newDlt.UnSortDrawResult = string(v.LotteryUnSortDrawResult)
					hmStrSlice := strings.Split(v.LotteryDrawResult, " ")
					if len(hmStrSlice) == 7 {
						newDlt.F1 = hmStrSlice[0]
						newDlt.F2 = hmStrSlice[1]
						newDlt.F3 = hmStrSlice[2]
						newDlt.F4 = hmStrSlice[3]
						newDlt.F5 = hmStrSlice[4]
						newDlt.B1 = hmStrSlice[5]
						newDlt.B2 = hmStrSlice[6]
						newDlt.Oe = ana_dlt.CalDltOe(hmStrSlice)
						hz := ana_dlt.CalDltHz(hmStrSlice)
						newDlt.Hz = hz
						newDlt.AeHz = ana_dlt.DltHzABCDE(hz)
						newDlt.Qzh = ana_dlt.CalDltQzh(hmStrSlice)
					}

					newDlt.PoolBalance, parseErr = strconv.ParseFloat(strings.ReplaceAll(v.PoolBalanceAfterDraw, ",", ""), 64)
					if parseErr != nil {
						newDlt.PoolBalance = 0
					}
					newDlt.TotalSaleAmount, parseErr = strconv.ParseFloat(strings.ReplaceAll(v.TotalSaleAmount, ",", ""), 64)
					if parseErr != nil {
						newDlt.TotalSaleAmount = 0
					}

					newDlt.StakeCount60, newDlt.StakeAmount60 = dlt.GetStake(v.PrizeLevelList, 60)
					newDlt.StakeCount80, newDlt.StakeAmount80 = dlt.GetStake(v.PrizeLevelList, 80)
					newDlt.StakeCount100, newDlt.StakeAmount100 = dlt.GetStake(v.PrizeLevelList, 100)
					newDlt.StakeCount101, newDlt.StakeAmount101 = dlt.GetStake(v.PrizeLevelList, 101)
					newDlt.StakeCount102, newDlt.StakeAmount102 = dlt.GetStake(v.PrizeLevelList, 102)
					newDlt.StakeCount201, newDlt.StakeAmount201 = dlt.GetStake(v.PrizeLevelList, 201)
					newDlt.StakeCount202, newDlt.StakeAmount202 = dlt.GetStake(v.PrizeLevelList, 202)
					newDlt.StakeCount301, newDlt.StakeAmount301 = dlt.GetStake(v.PrizeLevelList, 301)
					newDlt.StakeCount302, newDlt.StakeAmount302 = dlt.GetStake(v.PrizeLevelList, 302)
					newDlt.StakeCount401, newDlt.StakeAmount401 = dlt.GetStake(v.PrizeLevelList, 401)
					newDlt.StakeCount402, newDlt.StakeAmount402 = dlt.GetStake(v.PrizeLevelList, 402)
					newDlt.StakeCount501, newDlt.StakeAmount501 = dlt.GetStake(v.PrizeLevelList, 501)
					newDlt.StakeCount601, newDlt.StakeAmount601 = dlt.GetStake(v.PrizeLevelList, 601)
					newDlt.StakeCount701, newDlt.StakeAmount701 = dlt.GetStake(v.PrizeLevelList, 701)
					newDlt.StakeCount801, newDlt.StakeAmount801 = dlt.GetStake(v.PrizeLevelList, 801)
					newDlt.StakeCount901, newDlt.StakeAmount901 = dlt.GetStake(v.PrizeLevelList, 901)
					newDlt.StakeCount1001, newDlt.StakeAmount1001 = dlt.GetStake(v.PrizeLevelList, 1001)
					newDlt.StakeCount1101, newDlt.StakeAmount1101 = dlt.GetStake(v.PrizeLevelList, 1101)
					newDlt.DataSrc = 0

					for _, dt := range dlts {
						if dt.DrawNum == newDlt.DrawNum {
							db.DB.Model(dt).Updates(map[string]interface{}{
								"draw_time":          newDlt.DrawTime,
								"equipment_count":    newDlt.EquipmentCount,
								"draw_pdf_url":       newDlt.DrawPdfUrl,
								"unsort_draw_result": newDlt.UnSortDrawResult,
								"f1":                 newDlt.F1,
								"f2":                 newDlt.F2,
								"f3":                 newDlt.F3,
								"f4":                 newDlt.F4,
								"f5":                 newDlt.F5,
								"b1":                 newDlt.B1,
								"b2":                 newDlt.B2,
								"oe":                 newDlt.Oe,
								"hz":                 newDlt.Hz,
								"ae_hz":              newDlt.AeHz,
								"qzh":                newDlt.Qzh,
								"stake_count_101":    newDlt.StakeCount101,
								"stake_count_102":    newDlt.StakeCount101,
								"stake_count_201":    newDlt.StakeCount201,
								"stake_count_202":    newDlt.StakeCount202,
								"stake_count_301":    newDlt.StakeCount301,
								"stake_count_302":    newDlt.StakeCount302,
								"stake_count_401":    newDlt.StakeCount401,
								"stake_count_402":    newDlt.StakeCount402,
								"stake_count_501":    newDlt.StakeCount501,
								"stake_count_601":    newDlt.StakeCount601,
								"stake_count_701":    newDlt.StakeCount701,
								"stake_count_801":    newDlt.StakeCount801,
								"stake_count_901":    newDlt.StakeCount901,
								"stake_count_1001":   newDlt.StakeCount1001,
								"stake_count_1101":   newDlt.StakeCount1101,
								"stake_count_60":     newDlt.StakeCount60,
								"stake_count_80":     newDlt.StakeCount80,
								"stake_count_100":    newDlt.StakeCount100,
								"stake_amount_101":   newDlt.StakeAmount101,
								"stake_amount_102":   newDlt.StakeAmount102,
								"stake_amount_201":   newDlt.StakeAmount201,
								"stake_amount_202":   newDlt.StakeAmount202,
								"stake_amount_301":   newDlt.StakeAmount301,
								"stake_amount_302":   newDlt.StakeAmount302,
								"stake_amount_401":   newDlt.StakeAmount401,
								"stake_amount_402":   newDlt.StakeAmount401,
								"stake_amount_501":   newDlt.StakeAmount501,
								"stake_amount_601":   newDlt.StakeAmount601,
								"stake_amount_701":   newDlt.StakeAmount701,
								"stake_amount_801":   newDlt.StakeAmount801,
								"stake_amount_901":   newDlt.StakeAmount901,
								"stake_amount_1001":  newDlt.StakeAmount1001,
								"stake_amount_1101":  newDlt.StakeAmount1101,
								"stake_amount_60":    newDlt.StakeAmount60,
								"stake_amount_80":    newDlt.StakeAmount80,
								"stake_amount_100":   newDlt.StakeAmount100,

								"pool_balance":      newDlt.PoolBalance,
								"total_sale_amount": newDlt.TotalSaleAmount,
								"data_src":          0, // ✅ 一定会更新
							})
							lg.InfoToFileAndStdOut(fmt.Sprintf("已更新 %s ", newDlt.DrawNum))
						}
					}
				}
			}
		}
	}
}
