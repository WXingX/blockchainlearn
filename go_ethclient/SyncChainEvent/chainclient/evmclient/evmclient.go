package evmclient

import (
	"SyncChainEvent/chainclient/types"
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EvmClient struct {
	client *ethclient.Client
}

func NewEvmClient(nodeUrl string) (*EvmClient, error) {
	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		return nil, errors.New("failed to connect to Ethereum client")
	}
	return &EvmClient{client: client}, nil
}

func (e *EvmClient) Client() *ethclient.Client {
	return e.client
}

func (e *EvmClient) FilterLogs(ctx context.Context, q types.FilterQuery) ([]interface{}, error) {
	var addresses []common.Address
	for _, addr := range q.Addresses {
		addresses = append(addresses, common.HexToAddress(addr))
	}

	var topicsHash [][]common.Hash
	for _, topics := range q.Topics {
		var topicHash []common.Hash
		for _, topic := range topics {
			topicHash = append(topicHash, common.HexToHash(topic))
		}
		topicsHash = append(topicsHash, topicHash)
	}
	// Topics 事件签名/索引参数，AND 关系；二维数组内再 OR。
	// 当 fromBlock == ToBlock 时，获取 FromBlock中的日志信息
	queryParam := ethereum.FilterQuery{
		FromBlock: q.FromBlock,
		ToBlock:   q.ToBlock,
		Addresses: addresses,
		Topics:    topicsHash,
	}

	logs, err := e.client.FilterLogs(ctx, queryParam)
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs: %w", err)
	}

	var logEvents []interface{}
	for _, log := range logs {
		logEvents = append(logEvents, log)
	}

	return logEvents, nil
}

func (e *EvmClient) BlockTimeByNumber(ctx context.Context, blockNum *big.Int) (uint64, error) {
	header, err := e.client.HeaderByNumber(ctx, blockNum)
	if err != nil {
		return 0, fmt.Errorf("failed on get block header: %w", err)
	}

	return header.Time, nil
}

func (e *EvmClient) CallContractByChain(ctx context.Context, param types.CallParam) (interface{}, error) {
	return e.CallContract(ctx, param.EVMParam, param.BlockNumber)
}

func (e *EvmClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return e.client.CallContract(ctx, msg, blockNumber)
}

func (e *EvmClient) BlockNumber() (uint64, error) {
	var err error
	blockNum, err := e.client.BlockNumber(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed on get evm block number: %w", err)
	}

	return blockNum, nil
}

func (e *EvmClient) BlockWithTxs(ctx context.Context, blockNumber uint64) (interface{}, error) {
	blockWithTxs, err := e.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, fmt.Errorf("failed on get evm block: %w", err)
	}
	return blockWithTxs, nil
}

func (e *EvmClient) BlockHashByNumber(ctx context.Context, blockNum *big.Int) (string, error) {
	header, err := e.client.HeaderByNumber(ctx, blockNum)
	if err != nil {
		return "", fmt.Errorf("failed on get block BlockHashByNumber: %w", err)
	}

	return header.Hash().Hex(), nil
}
