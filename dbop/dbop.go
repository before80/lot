package dbop

import (
	"github.com/before80/lot/db"
	"github.com/before80/lot/models"
)

func GetLastDlt() (lastDlt models.Dlt) {
	db.DB.Last(&lastDlt)
	return lastDlt
}
