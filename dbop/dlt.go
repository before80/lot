package dbop

import (
	"fmt"
	"time"

	"github.com/before80/lot/db"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
	"gorm.io/gorm"
)

func InsertDlt(dlt models.Dlt) {
	db.DB.Create(&dlt)
}

func InsertTyp(typ models.Typ) {
	result := db.DB.Create(&typ)
	_ = result
	//if result.Error != nil {
	//	fmt.Println(result.Error)
	//}
}

func InsertAllDlt(allDlt models.AllDlt) {
	result := db.DB.Create(&allDlt)
	_ = result
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

func InsertDltMoni(moni models.DltMoni) {
	result := db.DB.Create(&moni)
	_ = result
	//if result.Error != nil {
	//	fmt.Println(result.Error)
	//}
}

func UpdateOrInsertDltMoni(moni models.DltMoni, groupNeedStr string) {
	var newMoni models.DltMoni
	var result *gorm.DB
	if groupNeedStr == "50000" {
		result = db.DB.Where("a = ?", moni.A).
			Where("typ = ?", moni.Typ).
			Where("method = ?", moni.Method).
			Find(&newMoni)
	} else {
		result = db.DB.Where("a = ?", moni.A).
			Where("b = ?", moni.B).
			Where("c = ?", moni.C).
			Where("d = ?", moni.D).
			Where("e = ?", moni.E).
			Where("typ = ?", moni.Typ).
			Where("method = ?", moni.Method).
			Find(&newMoni)
	}

	// 更新
	if result.RowsAffected != 0 {
		if groupNeedStr == "50000" {
			db.DB.Model(models.DltMoni{}).
				Where("a = ?", moni.A).
				Where("typ = ?", moni.Typ).
				Where("method = ?", moni.Method).
				Updates(models.DltMoni{A: moni.A, B: moni.B, C: moni.C, D: moni.D, E: moni.E, Cs: moni.Cs, UpdatedAt: time.Now()})
		} else {
			db.DB.Model(models.DltMoni{}).
				Where("a = ?", moni.A).
				Where("b = ?", moni.B).
				Where("c = ?", moni.C).
				Where("d = ?", moni.D).
				Where("e = ?", moni.E).
				Where("typ = ?", moni.Typ).
				Where("method = ?", moni.Method).
				Updates(models.DltMoni{Cs: moni.Cs, UpdatedAt: time.Now()})
		}
	} else {
		result1 := db.DB.Create(&moni)
		_ = result1
		//if result.Error != nil {
		//	fmt.Println(result.Error)
		//}
	}
}

func InsertAllDltBatch(allDlts []models.AllDlt, batchSize int) (insertedRow int, err error) {
	for i := 0; i < len(allDlts); i += batchSize {
		end := i + batchSize
		if end > len(allDlts) {
			end = len(allDlts)
		}
		batch := allDlts[i:end]

		if err = db.DB.Create(&batch).Error; err != nil {
			return insertedRow, err
		}
		insertedRow += len(batch)
		lg.InfoToFile(fmt.Sprintf("Inserted batch %d-%d\n", i, end-1))
	}
	return insertedRow, nil
}

func InsertDltBatch(dlts []models.Dlt, batchSize int) (insertedRow int, err error) {
	for i := 0; i < len(dlts); i += batchSize {
		end := i + batchSize
		if end > len(dlts) {
			end = len(dlts)
		}
		batch := dlts[i:end]

		if err = db.DB.Create(&batch).Error; err != nil {
			return insertedRow, err
		}
		insertedRow += len(batch)
		lg.InfoToFile(fmt.Sprintf("Inserted batch %d-%d\n", i, end-1))
	}
	return insertedRow, nil
}

func GetLastDlt() (lastDlt models.Dlt) {
	db.DB.Last(&lastDlt)
	return lastDlt
}

// ReadAllDlt 读取所有大乐透开奖数据
func ReadAllDlt(desc bool) (dlts []models.Dlt, err error) {
	if desc {
		if err = db.DB.Order("id desc").Find(&dlts).Error; err != nil {
			return dlts, err
		}
		return dlts, nil
	}

	if err = db.DB.Order("id asc").Find(&dlts).Error; err != nil {
		return dlts, err
	}
	return dlts, nil
}

func ReadDltGETDrawNum(drawNum string) (dlts []models.Dlt, err error) {
	if err = db.DB.Where("draw_num >= ?", drawNum).Order("id asc").Find(&dlts).Error; err != nil {
		return dlts, err
	}
	return dlts, nil
}

func ReadDlt(drawNum string) (models.Dlt, error) {
	var dlts []models.Dlt
	if err := db.DB.Where("draw_num = ?", drawNum).Find(&dlts).Error; err != nil {
		return models.Dlt{}, err
	}
	return dlts[0], nil
}

func ReadLastTyp() (lastTyp models.Typ) {
	db.DB.Last(&lastTyp)
	return lastTyp
	//var typs []models.Typ
	//if err := db.DB.Order("id desc").Find(&typs).Error; err != nil {
	//	return models.Typ{}, err
	//}
	//return typs[0], nil
}

func ReadAllTyp(desc bool) (typs []models.Typ, err error) {
	if desc {
		if err = db.DB.Order("id desc").Find(&typs).Error; err != nil {
			return typs, err
		}
		return typs, nil
	}

	if err = db.DB.Find(&typs).Error; err != nil {
		return typs, err
	}
	return typs, nil
}

func ReadAllDlts(desc bool) (dlts []models.AllDlt, err error) {
	if desc {
		if err = db.DB.Order("id desc").Find(&dlts).Error; err != nil {
			return dlts, err
		}
		return dlts, nil
	}

	if err = db.DB.Order("id desc").Find(&dlts).Error; err != nil {
		return dlts, err
	}
	return dlts, nil
}
