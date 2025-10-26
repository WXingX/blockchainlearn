package model

import (
	"time"
)

type SyncRecord struct {
	ID                  int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ChainName           string    `gorm:"not null"`                      // 链名称
	ChainID             int64     `gorm:"index:idx_chain_id"`            // 链ID
	LastSyncBlockNumber uint64    `gorm:"type:bigint unsigned;not null"` // 上次同步区块号,这个区块已经同步，需要从 +1 区块继续同步
	LastSyncBlockHash   string    `gorm:"not null"`                      // 上次同步区块哈希
	LastSyncBlockTime   time.Time `gorm:"not null"`                      // 上次同步时间
}

func (SyncRecord) TableName() string {
	return "tbl_sync_record"
}
