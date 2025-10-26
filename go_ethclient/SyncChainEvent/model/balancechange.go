package model

import (
	"time"
)

type BalanceChange struct {
	Id      int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Account string    `gorm:"column:account"`
	Time    time.Time `gorm:"column:time"`
	Balance string    `gorm:"column:balance"`
}
