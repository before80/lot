package dbop

import (
	"fmt"
	"time"

	"github.com/before80/lot/db"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
	"gorm.io/gorm"
)

func InsertSsqBatch(ssqs []models.Ssq, batchSize int) (insertedRow int, err error) {
	for i := 0; i < len(ssqs); i += batchSize {
		end := i + batchSize
		if end > len(ssqs) {
			end = len(ssqs)
		}
		batch := ssqs[i:end]

		if err = db.DB.Create(&batch).Error; err != nil {
			return insertedRow, err
		}
		insertedRow += len(batch)
		lg.InfoToFile(fmt.Sprintf("Inserted batch %d-%d\n", i, end-1))
	}
	return insertedRow, nil
}

func InsertAllSsqBatch(allSsqs []models.AllSsq, batchSize int) (insertedRow int, err error) {
	for i := 0; i < len(allSsqs); i += batchSize {
		end := i + batchSize
		if end > len(allSsqs) {
			end = len(allSsqs)
		}
		batch := allSsqs[i:end]

		if err = db.DB.Create(&batch).Error; err != nil {
			return insertedRow, err
		}
		insertedRow += len(batch)
		lg.InfoToFile(fmt.Sprintf("Inserted batch %d-%d\n", i, end-1))
	}
	return insertedRow, nil
}

func UpdateOrInsertSsqMoni(moni models.SsqMoni, groupNeedStr string) {
	var newMoni models.SsqMoni
	var result *gorm.DB
	if groupNeedStr == "500000" {
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
			Where("f = ?", moni.F).
			Where("typ = ?", moni.Typ).
			Where("method = ?", moni.Method).
			Find(&newMoni)
	}

	// 更新
	if result.RowsAffected != 0 {
		if groupNeedStr == "500000" {
			db.DB.Model(models.SsqMoni{}).
				Where("a = ?", moni.A).
				Where("typ = ?", moni.Typ).
				Where("method = ?", moni.Method).
				Updates(models.SsqMoni{A: moni.A, B: moni.B, C: moni.C, D: moni.D, E: moni.E, Cs: moni.Cs, UpdatedAt: time.Now()})
		} else {
			db.DB.Model(models.SsqMoni{}).
				Where("a = ?", moni.A).
				Where("b = ?", moni.B).
				Where("c = ?", moni.C).
				Where("d = ?", moni.D).
				Where("e = ?", moni.E).
				Where("f = ?", moni.F).
				Where("typ = ?", moni.Typ).
				Where("method = ?", moni.Method).
				Updates(models.SsqMoni{Cs: moni.Cs, UpdatedAt: time.Now()})
		}
	} else {
		result1 := db.DB.Create(&moni)
		_ = result1
		//if result.Error != nil {
		//	fmt.Println(result.Error)
		//}
	}
}

func GetLastSsq() (lastSsq models.Ssq) {
	db.DB.Last(&lastSsq)
	return lastSsq
}

func ReadAllSsq(desc bool) (ssqs []models.Ssq, err error) {
	if desc {
		if err = db.DB.Order("id desc").Find(&ssqs).Error; err != nil {
			return ssqs, err
		}
		return ssqs, nil
	} else {
		if err = db.DB.Order("id asc").Find(&ssqs).Error; err != nil {
			return ssqs, err
		}
		return ssqs, nil
	}

}

func ReadSsqGETDrawNum(drawNum string) (ssqs []models.Ssq, err error) {
	if err = db.DB.Where("draw_num >= ?", drawNum).Order("id asc").Find(&ssqs).Error; err != nil {
		return ssqs, err
	}
	return ssqs, nil
}

func ReadSsq(drawNum string) (models.Ssq, error) {
	var ssqs []models.Ssq
	if err := db.DB.Where("draw_num = ?", drawNum).Find(&ssqs).Error; err != nil {
		return models.Ssq{}, err
	}
	return ssqs[0], nil
}

func ReadAllSsqs(desc bool) (ssqs []models.AllSsq, err error) {
	if desc {
		if err = db.DB.Order("id desc").Find(&ssqs).Error; err != nil {
			return ssqs, err
		}
		return ssqs, nil
	}

	if err = db.DB.Find(&ssqs).Error; err != nil {
		return ssqs, err
	}
	return ssqs, nil
}
