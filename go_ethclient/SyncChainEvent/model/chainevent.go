package model

type SyncChainEvent struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement;"`
	ChainID   int64  `gorm:"column:chain_id;not null"`
	EventHash string `gorm:"column:event_hash;not null"`
	EventName string `gorm:"column:event_name;not null"`
}

func (e *SyncChainEvent) TableName() string {
	return "tbl_chain_event"
}
