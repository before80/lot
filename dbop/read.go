package dbop

import (
	"github.com/before80/lot/db"
	"github.com/before80/lot/models"
)

func ReadAllDlt(desc bool) (dlts []models.Dlt, err error) {
	if desc {
		if err = db.DB.Order("id desc").Find(&dlts).Error; err != nil {
			return dlts, err
		}
		return dlts, nil
	}

	if err = db.DB.Find(&dlts).Error; err != nil {
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
