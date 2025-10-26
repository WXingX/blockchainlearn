package eventsync

import (
	"SyncChainEvent/chainclient/evmclient"
	"SyncChainEvent/chainclient/types"
	"SyncChainEvent/logger/xzap"
	"SyncChainEvent/model"
	"SyncChainEvent/utils/calcutils"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"SyncChainEvent/utils/timeutils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethereumTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeromicro/go-zero/core/threading"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 事件同步

type SyncServer struct {
	ctx           context.Context
	name          string
	id            int64
	delayBlockNum uint64
	tokenAddr     string
	deployerAddr  string
	parseAbi      abi.ABI
	db            *gorm.DB
	evmClient     *evmclient.EvmClient
}

const (
	syncRecordTableName      = "tbl_sync_record"
	transferTablePrefix      = "tbl_transfer_details_"
	accountTablePrefix       = "tbl_account_"
	balanceChangeTablePrefix = "tbl_balance_change_"
	zeroAddress              = "0x0000000000000000000000000000000000000000"
	SleepInterval            = 10 // second
	SyncBlockPeriod          = 5  // 一次最多获取5个区块的数据，这个是由于 rpc服务的限制导致的
	ERC20TokenAbi            = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"allowance","type":"uint256"},{"internalType":"uint256","name":"needed","type":"uint256"}],"name":"ERC20InsufficientAllowance","type":"error"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"uint256","name":"balance","type":"uint256"},{"internalType":"uint256","name":"needed","type":"uint256"}],"name":"ERC20InsufficientBalance","type":"error"},{"inputs":[{"internalType":"address","name":"approver","type":"address"}],"name":"ERC20InvalidApprover","type":"error"},{"inputs":[{"internalType":"address","name":"receiver","type":"address"}],"name":"ERC20InvalidReceiver","type":"error"},{"inputs":[{"internalType":"address","name":"sender","type":"address"}],"name":"ERC20InvalidSender","type":"error"},{"inputs":[{"internalType":"address","name":"spender","type":"address"}],"name":"ERC20InvalidSpender","type":"error"},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"OwnableInvalidOwner","type":"error"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"OwnableUnauthorizedAccount","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"TransferLog","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"burn","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

func New(ctx context.Context, name string, id int64, tokenAddress string, deployerAddr string, rpcUrl string, delayBlockNum uint64, db *gorm.DB) (*SyncServer, error) {
	if name == "" || id == 0 || tokenAddress == "" || deployerAddr == "" || rpcUrl == "" {
		return nil, errors.New("invalid params")
	}

	if db == nil {
		return nil, errors.New("db is nil")
	}

	parsedAbi, _ := abi.JSON(strings.NewReader(ERC20TokenAbi))

	evmClient, err := evmclient.NewEvmClient(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("new evm client error: %w", err)
	}

	return &SyncServer{
		ctx:           ctx,
		name:          name,
		id:            id,
		delayBlockNum: delayBlockNum,
		tokenAddr:     tokenAddress,
		deployerAddr:  deployerAddr,
		parseAbi:      parsedAbi,
		db:            db,
		evmClient:     evmClient,
	}, nil
}

func (e *SyncServer) Start() error {
	//1. 判断同步记录表中是否有记录
	var syncRecord model.SyncRecord

	e.db.Where("chain_id = ?", e.id).First(&syncRecord)
	if syncRecord.LastSyncBlockHash == "" || syncRecord.LastSyncBlockNumber == 0 {
		return errors.New("sync record not exist")
	}

	//2. 判断是否存在转账记录表和账号表
	detailExists := e.db.Migrator().HasTable(e.GetTransferDetailTableName())
	if !detailExists {
		if !e.CreateTransferDetailTable() {
			return errors.New("create transfer detail table failed")
		}

	}

	accountExists := e.db.Migrator().HasTable(e.GetAccountTableName())
	if !accountExists {
		if !e.CreateAccountTable() {
			return errors.New("create account table failed")
		}
	}

	// 余额变化表
	balanceChangeExists := e.db.Migrator().HasTable(e.GetBalanceChangeTableName())
	if !balanceChangeExists {
		if !e.CreateBalanceChangeTable() {
			return errors.New("create balance change table failed")
		}
	}

	//3. 开始同步
	threading.GoSafe(func() {
		e.SyncEventLoop(syncRecord.LastSyncBlockNumber)
	})
	return nil
}

func (e *SyncServer) GetTransferDetailTableName() string {
	return transferTablePrefix + strings.Replace(strings.ToLower(e.name), " ", "_", -1)
}

func (e *SyncServer) GetAccountTableName() string {
	return accountTablePrefix + strings.Replace(strings.ToLower(e.name), " ", "_", -1)
}

func (e *SyncServer) GetBalanceChangeTableName() string {
	return balanceChangeTablePrefix + strings.Replace(strings.ToLower(e.name), " ", "_", -1)
}

func (e *SyncServer) CreateTransferDetailTable() bool {
	tableName := e.GetTransferDetailTableName()
	sqlStr := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		id				bigint NOT NULL AUTO_INCREMENT,
		chain_id int 	NOT NULL DEFAULT 0,
		from_account	varchar(255) DEFAULT NULL,
		to_account		varchar(255) DEFAULT NULL,
		amount 			varchar(255) NOT NULL,
		block_number 	bigint unsigned  NOT NULL,
		block_hash 		varchar(255) NOT NULL
		block_time 		datetime(0) NOT NULL,
		tx_hash 		varchar(255) NOT NULL
	) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic`
	result := e.db.Exec(sqlStr)
	if result.Error != nil {
		xzap.WithContext(e.ctx).Error("Failed to create table:"+tableName, zap.Error(result.Error))
		return false
	}
	return true
}

func (e *SyncServer) CreateAccountTable() bool {
	tableName := e.GetAccountTableName()
	sqlStr := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		id				bigint NOT NULL AUTO_INCREMENT,
		account			varchar(255) DEFAULT NULL,
		balance 		varchar(255) DEFAULT NULL,
		point 			varchar(255) DEFAULT NULL
	) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic`
	result := e.db.Exec(sqlStr)
	if result.Error != nil {
		xzap.WithContext(e.ctx).Error("Failed to create table:"+tableName, zap.Error(result.Error))
		return false
	}
	return true
}

func (e *SyncServer) CreateBalanceChangeTable() bool {
	tableName := e.GetBalanceChangeTableName()
	sqlStr := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		id				bigint NOT NULL AUTO_INCREMENT,
		account			varchar(255) DEFAULT NULL,
		time			datetime(0) DEFAULT NULL,
		balance 		varchar(255) DEFAULT NULL,
		PRIMARY KEY (id) USING BTREE
	) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic`
	result := e.db.Exec(sqlStr)
	if result.Error != nil {
		xzap.WithContext(e.ctx).Error("Failed to create table: "+tableName, zap.Error(result.Error))
		return false
	}
	return true
}

func (e *SyncServer) SyncEventLoop(lastBlockNumber uint64) {
	xzap.WithContext(e.ctx).Info("SyncEventLoop start lastBlockNumber = " + strconv.FormatUint(lastBlockNumber, 10))

	// 从数据库中获取同步事件的 hash值
	var chainEvent model.SyncChainEvent
	e.db.Where("chain_id = ?", e.id).First(&chainEvent)
	if chainEvent.EventHash == "" {
		xzap.WithContext(e.ctx).Error("SyncEventLoop end event is empty.")
		return
	}

	for {
		select {
		case <-e.ctx.Done():
			xzap.WithContext(e.ctx).Info("SyncEventLoop stopped due to context cancellation")
			return
		default:
		}
		// 获取当前最新的区块
		currentBlockNum, err := e.evmClient.BlockNumber()
		if err != nil {
			xzap.WithContext(e.ctx).Error("failed on get current block number", zap.Error(err))
			time.Sleep(SleepInterval * time.Second)
			continue
		}

		if lastBlockNumber >= currentBlockNum-e.delayBlockNum {
			xzap.WithContext(e.ctx).Info("lastBlockNumber more than currentBlockNum - delayBlockNum...")
			time.Sleep(SleepInterval * time.Second)
			continue
		}

		// 获取当前应该要同步到哪个区块
		startBlock := lastBlockNumber + 1
		endBlock := startBlock + SyncBlockPeriod
		if endBlock > currentBlockNum-e.delayBlockNum {
			endBlock = currentBlockNum - e.delayBlockNum
		}

		query := types.FilterQuery{
			FromBlock: new(big.Int).SetUint64(startBlock),
			ToBlock:   new(big.Int).SetUint64(endBlock),
			Addresses: []string{e.tokenAddr},
		}

		// 获取事件
		logs, err := e.evmClient.FilterLogs(e.ctx, query)
		if err != nil {
			xzap.WithContext(e.ctx).Error("failed to get log, error:", zap.Error(err))
			time.Sleep(SleepInterval * time.Second)
			continue
		}

		for _, log := range logs {
			ethLog := log.(ethereumTypes.Log)
			if ethLog.Topics[0].String() == chainEvent.EventHash {
				e.handleTransferLogEvent(ethLog)
			}
		}

		// 区块同步完成，记录到数据库中
		lastBlockNumber = endBlock
		// 获取lastBlockNumber 的参数
		blockTime, err1 := e.evmClient.BlockTimeByNumber(e.ctx, new(big.Int).SetUint64(endBlock))
		if err1 != nil {
			xzap.WithContext(e.ctx).Error("failed on get current block time, endBlock = "+strconv.FormatUint(endBlock, 10), zap.Error(err1))
			continue
		}
		blockHash, err2 := e.evmClient.BlockHashByNumber(e.ctx, new(big.Int).SetUint64(endBlock))
		if err2 != nil {
			xzap.WithContext(e.ctx).Error("failed on get current block hash, endBlock = "+strconv.FormatUint(endBlock, 10), zap.Error(err1))
			continue
		}
		if err := e.db.Table(syncRecordTableName).
			Where("chain_id = ? ", e.id).
			Update("last_sync_block_hash", blockHash).
			Update("last_sync_block_number", lastBlockNumber).
			Update("last_sync_block_time", timeutils.UnixToTime(int64(blockTime))).Error; err != nil {
			xzap.WithContext(e.ctx).Error("failed on update event sync block info",
				zap.Error(err))
			return
		}

		xzap.WithContext(e.ctx).Info("sync event ...",
			zap.Uint64("start_block", startBlock),
			zap.Uint64("end_block", endBlock))
	}
}

func (e *SyncServer) handleTransferLogEvent(log ethereumTypes.Log) bool {
	var transferEvent struct {
		//From   common.Address
		//To     common.Address
		Amount *big.Int
	}
	err := e.parseAbi.UnpackIntoInterface(&transferEvent, "TransferLog", log.Data)
	if err != nil {
		xzap.WithContext(e.ctx).Error("failed to unpack TransferLog", zap.Error(err))
		return false
	}

	from := common.BytesToAddress(log.Topics[1].Bytes())
	to := common.BytesToAddress(log.Topics[2].Bytes())
	blockTime := timeutils.Format(timeutils.UnixToTime(int64(log.BlockTimestamp)), timeutils.FormatString)

	//插入数据库
	data := map[string]interface{}{
		"chain_id":     e.id,
		"from_account": from.Hex(),
		"to_account":   to.Hex(),
		"amount":       transferEvent.Amount.String(),
		"block_number": log.BlockNumber,
		"block_hash":   log.BlockHash.Hex(),
		"block_time":   blockTime,
		"tx_hash":      log.TxHash.Hex(),
	}
	result := e.db.Table(e.GetTransferDetailTableName()).Create(&data)
	if result.Error != nil {
		xzap.WithContext(e.ctx).Error("failed to create transferDetail", zap.Error(result.Error))
		return false
	}

	//	更新账号表
	if from.Hex() == zeroAddress && to.Hex() != zeroAddress { // mint
		var toAcc model.AccountBalance
		e.db.Table(e.GetAccountTableName()).Where("account = ?", to.Hex()).Find(&toAcc)
		// 有就更新
		if toAcc.Account == "" {
			// 新增
			toAcc.Balance = transferEvent.Amount.String()
			toAcc.Account = to.Hex()
			res := e.UpdateDataTable(toAcc, "create", log.BlockTimestamp)
			if !res {
				xzap.WithContext(e.ctx).Error("UpdateDataTable create mint failed.")
				return false
			}
		} else {
			//	更新
			toAcc.Balance = calcutils.BigIntAdd(toAcc.Balance, transferEvent.Amount.String())

			res := e.UpdateDataTable(toAcc, "update", log.BlockTimestamp)
			if !res {
				xzap.WithContext(e.ctx).Error("UpdateDataTable update mint failed.")
				return false
			}
		}
	}

	if from.Hex() != zeroAddress && to.Hex() == zeroAddress { // burn
		var fromAcc model.AccountBalance
		e.db.Table(e.GetAccountTableName()).Where("account = ?", from.Hex()).Find(&fromAcc)
		if fromAcc.Account == "" { // burn时，表中要有 fromAccount
			xzap.WithContext(e.ctx).Error("burn event from account not in account table.")
			return false
		} else {
			fromAcc.Balance = calcutils.BigIntSub(fromAcc.Balance, transferEvent.Amount.String())
			res := e.UpdateDataTable(fromAcc, "update", log.BlockTimestamp)
			if !res {
				xzap.WithContext(e.ctx).Error("UpdateDataTable update burn failed.")
				return false
			}
		}
	}

	if from.Hex() != zeroAddress && to.Hex() != zeroAddress { // transfer
		if from.Hex() != e.deployerAddr { // 如果不是合约部署者的地址，就检测
			var fromAcc model.AccountBalance
			e.db.Table(e.GetAccountTableName()).Where("account = ?", from.Hex()).Find(&fromAcc)
			if fromAcc.Account == "" {
				xzap.WithContext(e.ctx).Error("transfer event from account not in account table.")
				return false
			} else {
				fromAcc.Balance = calcutils.BigIntSub(fromAcc.Balance, transferEvent.Amount.String())
				res := e.UpdateDataTable(fromAcc, "update", log.BlockTimestamp)
				if !res {
					xzap.WithContext(e.ctx).Error("UpdateDataTable update transfer failed.")
					return false
				}
			}
		}

		var toAcc model.AccountBalance
		e.db.Table(e.GetAccountTableName()).Where("account = ?", to.Hex()).Find(&toAcc)
		// 有就更新
		if toAcc.Account == "" {
			// 新增
			toAcc.Balance = transferEvent.Amount.String()
			toAcc.Account = to.Hex()
			res := e.UpdateDataTable(toAcc, "create", log.BlockTimestamp)
			if !res {
				xzap.WithContext(e.ctx).Error("UpdateDataTable create transfer failed.")
				return false
			}
		} else {
			//	更新
			toAcc.Balance = calcutils.BigIntAdd(toAcc.Balance, transferEvent.Amount.String())

			res := e.UpdateDataTable(toAcc, "update", log.BlockTimestamp)
			if !res {
				xzap.WithContext(e.ctx).Error("UpdateDataTable update transfer failed.")
				return false
			}
		}
	}

	return true
}

func (e *SyncServer) UpdateDataTable(acc model.AccountBalance, op string, blockTimestamp uint64) bool {
	var result *gorm.DB
	if op == "create" {
		result = e.db.Table(e.GetAccountTableName()).Create(&acc)
		if result.Error != nil {
			xzap.WithContext(e.ctx).Error("UpdateDataTable failed to create account balance, acc:"+acc.Account, zap.Error(result.Error))
			return false
		}
	} else if op == "update" {
		result = e.db.Table(e.GetAccountTableName()).Save(&acc)
		if result.Error != nil {
			xzap.WithContext(e.ctx).Error("UpdateDataTable failed to update account balance, acc:"+acc.Account, zap.Error(result.Error))
			return false
		}
	}

	var change model.BalanceChange
	change.Account = acc.Account
	change.Balance = acc.Balance
	change.Time = timeutils.UnixToTime(int64(blockTimestamp))
	result = e.db.Table(e.GetBalanceChangeTableName()).Create(&change)
	if result.Error != nil {
		xzap.WithContext(e.ctx).Error("UpdateDataTable failed to update account change table, account:"+change.Account, zap.Error(result.Error))
		return false
	}
	return true
}
