package model

import (
	"time"
)

type PointCountRecord struct {
	Id               int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ChainId          int64     `gorm:"column:chain_id"`
	LastExecDateTime time.Time `gorm:"column:last_exec_date_time"`
}
