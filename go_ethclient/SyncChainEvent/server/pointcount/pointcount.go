package pointcount

import (
	"SyncChainEvent/chainclient/evmclient"
	"SyncChainEvent/logger/xzap"
	"SyncChainEvent/model"
	"SyncChainEvent/utils/calcutils"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CountServer 计算积分
type CountServer struct {
	ctx       context.Context
	cancel    context.CancelFunc
	name      string
	id        int64
	rpcUrl    string // 用于获取最新的blocknum，判断事件同步是否正常运行
	db        *gorm.DB
	evmClient *evmclient.EvmClient
	cron      *cron.Cron
	isClosed  bool
}

const (
	accountTablePrefix          = "tbl_account_"
	pointCountRecordTablePrefix = "tbl_point_count_record_"
	balanceChangeTablePrefix    = "tbl_balance_change_"
	OneMinute                   = 60 // 1分钟
	CronJobString               = "0 0 * * * *"
	PointRate                   = 0.05           // 积分比例 每小时 5%
	HourlyPointRate             = PointRate / 60 // 每分钟积分比例 5% / 60
	// FormatString                = "2006-01-02 15:04:05"
)

func New(ctx context.Context, name string, id int64, rpcUrl string, db *gorm.DB) (*CountServer, error) {
	if name == "" || id == 0 || rpcUrl == "" {
		return nil, errors.New("invalid params")
	}

	if db == nil {
		return nil, errors.New("db is nil")
	}

	// 判断是否存在积分记录表
	pointCountTable := pointCountRecordTablePrefix + strings.Replace(strings.ToLower(name), " ", "_", -1)
	if !db.Migrator().HasTable(pointCountTable) {
		sqlStr := `CREATE TABLE IF NOT EXISTS ` + pointCountTable + ` (
			id 					bigint NOT NULL AUTO_INCREMENT PRIMARY KEY,
			chain_id 			int NOT NULL DEFAULT 0,
			last_exec_date_time datetime(0) NOT NULL
		) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic`
		if err := db.Exec(sqlStr).Error; err != nil {
			return nil, fmt.Errorf("create point count table failed: %w", err)
		}
	}

	evmClient, err := evmclient.NewEvmClient(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("new evm client error: %w", err)
	}
	ctx, cancel := context.WithCancel(ctx)

	return &CountServer{
		ctx:       ctx,
		name:      name,
		id:        id,
		rpcUrl:    rpcUrl,
		db:        db,
		evmClient: evmClient,
		cron:      cron.New(cron.WithSeconds()),
		cancel:    cancel,
	}, nil
}

func (p *CountServer) Start() {
	_, err := p.cron.AddFunc(CronJobString, func() {
		select {
		case <-p.ctx.Done():
			return
		default:
			p.CountLoop()
		}
	})
	if err != nil {
		xzap.WithContext(p.ctx).Error("add func job error", zap.Error(err))
		return
	}
	p.cron.Start()
}

func (p *CountServer) Stop() {
	if p.isClosed {
		return
	}
	p.cancel()
	p.cron.Stop()
	p.isClosed = true
	xzap.WithContext(p.ctx).Info("stop point count")
}

func (p *CountServer) CountLoop() {
	// 1.获取最新的区块，和区块时间，任务执行时间

	// 获取当期日期的日期
	now := time.Now()

	// 获取时间戳
	timeStamp := now.Unix() // 秒

	blockNum, err := p.evmClient.BlockNumber()
	if err != nil {
		xzap.WithContext(p.ctx).Error("get evm block number failed. err: ", zap.Error(err))
		return
	}

	blockTime, err := p.evmClient.BlockTimeByNumber(p.ctx, new(big.Int).SetUint64(blockNum))
	if err != nil {
		xzap.WithContext(p.ctx).Error("get evm block time failed. err: ", zap.Error(err))
		return
	}

	// 判断最新区块时间和当前时间是否超过1分钟，如果能获取到最新的区块，说明RPC 服务正常
	if blockTime > uint64(timeStamp) {
		xzap.WithContext(p.ctx).Error("get evm block time failed. err: ", zap.Error(err))
		return
	}
	if blockTime < uint64(timeStamp-OneMinute) {
		xzap.WithContext(p.ctx).Error("get evm block time failed. err: ", zap.Error(err))
		return
	}

	// 3.从tbl_point_count_${name} 表中获取任务上一次执行时间
	var lastCountRecord model.PointCountRecord
	result := p.db.Table(p.getPointCountTableName()).Where("chain_id = ?", p.id).Scan(&lastCountRecord)
	if result.Error != nil || lastCountRecord.LastExecDateTime.IsZero() {
		xzap.WithContext(p.ctx).Error("get last record failed. err: ", zap.Error(result.Error))
		return
	}

	// 4.从tbl_account 表中获取所有account,遍历所有account
	var accounts []model.AccountBalance
	result = p.db.Table(p.getAccountTableName()).Scan(&accounts)
	if result.Error != nil {
		xzap.WithContext(p.ctx).Error("get account balance failed. err: ", zap.Error(err))
		return
	}

	if len(accounts) == 0 {
		xzap.WithContext(p.ctx).Error("get account balance failed. err: ", zap.Error(err))
		return
	}

	// 5.从tbl_balance_change_ 表中根据任务执行的周期获取 account对应的记录
	for _, account := range accounts {
		var records []model.BalanceChange
		// 获取lastRecord.LastExecDateTime 和当前时间内的记录
		result = p.db.Table(p.getAccountBalanceChangeTableName()).
			Where("account = ?", account.Account).
			Where("time >= ?", lastCountRecord.LastExecDateTime).
			Where("time <= ?", now).
			Order("id DESC").Scan(&records)
		if result.Error != nil {
			xzap.WithContext(p.ctx).Error("get account balance change failed. err: ", zap.Error(result.Error))
			return
		}

		if len(records) == 0 {
			// 5.1 如果没有记录，则获取最新的一条记录A，积分计算为 A.balance * 0.05
			var newRecord model.BalanceChange
			result = p.db.Table(p.getAccountBalanceChangeTableName()).
				Where("account = ?", account.Account).
				Order("id DESC").Limit(1).Scan(&newRecord)
			if result.Error != nil {
				xzap.WithContext(p.ctx).Error("get account balance change failed. err: ", zap.Error(result.Error))
				continue
			}

			if account.Point == "" {
				account.Point = "0"
			}
			sumPoint := calcutils.BigIntMul(newRecord.Balance, 0.05)
			account.Point = calcutils.BigIntAdd(account.Point, sumPoint)

			result = p.db.Table(p.getAccountTableName()).
				Where("account = ?", account.Account).
				Update("point", account.Point)
			if result.RowsAffected == 0 {
				xzap.WithContext(p.ctx).Error("update account point failed. err: ", zap.Error(result.Error))
				continue
			}
		} else {
			// 5.2 获取记录切片中 最早的一条记录A，然后查询记录A的上一条记录B
			// 如果记录B 存在，则 积分为 B.balance * 0.05 *(A记录的时间 - 上一次任务执行的时间)的分钟数/60
			// 如果记录B不存在，则逆序处理 records, 计算 每一条记录的积分，累加起来。
			// 令记录records中A前面一条记录为 C， 则计算公式为 A.balance * 0.05 *(C记录的时间 - A记录的时间)的分钟数/60
			// 最后计算第一条记录 D 的积分公式为：D.balance * 0.05 *(本次任务执行时间 - D记录的时间)的分钟数/60

			// 获取records的最后一个元素
			lastChangeRecord := records[len(records)-1]
			// 从数据库中 获取lastRecord 同一个account的前一条记录
			var prevRecord model.BalanceChange
			result = p.db.Table(p.getAccountBalanceChangeTableName()).
				Where("account = ?", lastChangeRecord.Account).
				Where("time < ?", lastChangeRecord.Time).
				Order("id DESC").Limit(1).Scan(&prevRecord)
			if result.Error != nil {
				xzap.WithContext(p.ctx).Error("get account balance change failed. err: ", zap.Error(result.Error))
				continue
			}

			var prevPoint string = "0"
			if !prevRecord.Time.IsZero() {
				// 说明以前有数据，有数据的情况下要计算积分
				// 获取 balance, 用lastRecord的时间减去 上一次执行的时间，获得分钟数m1,计算 lastRecord.balance * 0.05 * m1 / 60
				//prevRecordBalance := prevRecord.Balance.Int64()

				lastRecordTimeToNow := lastChangeRecord.Time.Sub(lastCountRecord.LastExecDateTime)
				lastRecordTimeToNowMinutes := lastRecordTimeToNow.Minutes()
				prevPoint = calcutils.BigIntMul(prevRecord.Balance, lastRecordTimeToNowMinutes*HourlyPointRate)
			}

			// 循环遍历 records，用
			// records 反转
			var sumPoint string = "0"
			length := len(records)
			for i := length - 1; i >= 0; i-- {
				r1 := records[i]
				var recordToNowMinutes float64 = 0.0
				if i == 0 {
					// 说明是第一条记录
					recordToNowMinutes = now.Sub(r1.Time).Minutes()
				} else {
					r2 := records[i-1]
					recordToNowMinutes = r2.Time.Sub(r1.Time).Minutes()
					//fmt.Printf("recordToNowMinutes = %f\n", recordToNowMinutes)
				}
				crrPoint := calcutils.BigIntMul(r1.Balance, recordToNowMinutes*HourlyPointRate)
				sumPoint = calcutils.BigIntAdd(sumPoint, crrPoint)
			}

			if account.Point == "" {
				account.Point = "0"
			}

			account.Point = calcutils.BigIntAdd(account.Point, prevPoint, sumPoint)
			result = p.db.Table(p.getAccountTableName()).
				Where("account = ?", account.Account).
				Update("point", account.Point)
			if result.RowsAffected == 0 {
				xzap.WithContext(p.ctx).Error("update account point failed. err: ", zap.Error(result.Error))
			}
		}
	}

	// 执行完成后，记录结果
	result = p.db.Table(p.getPointCountTableName()).
		Where("chain_id = ?", p.id).
		Update("last_exec_date_time", now.Format("2006-01-02 15:00:00"))
	if result.Error != nil {
		xzap.WithContext(p.ctx).Error("update point count record failed. err: ", zap.Error(result.Error))
		return
	}
}

func (p *CountServer) getPointCountTableName() string {
	return pointCountRecordTablePrefix + strings.Replace(strings.ToLower(p.name), " ", "_", -1)
}

func (p *CountServer) getAccountTableName() string {
	return accountTablePrefix + strings.Replace(strings.ToLower(p.name), " ", "_", -1)
}

func (p *CountServer) getAccountBalanceChangeTableName() string {
	return balanceChangeTablePrefix + strings.Replace(strings.ToLower(p.name), " ", "_", -1)
}
