package model

type AccountBalance struct {
	Id      int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Account string `gorm:"column:account"`
	Balance string `gorm:"column:balance"`
	Point   string `gorm:"column:point"`
}
