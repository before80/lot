package dbop

import (
	"fmt"
	"github.com/before80/lot/db"
	"github.com/before80/lot/lg"
	"github.com/before80/lot/models"
)

func InsertDlt(dlt models.Dlt) {
	db.DB.Create(&dlt)
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
