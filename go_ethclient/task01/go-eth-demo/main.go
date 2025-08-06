package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func main() {
	url := "https://quaint-side-daylight.ethereum-sepolia.quiknode.pro/b2b5094e028ada682797d2268e7116e13e7cb8c9"
	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Printf("Dial err: %s", err.Error())
		return
	}

	//QueryBlockInfo(client, nil)
	NewTransaction(client)
}

// QueryBlockInfo 查询区块信息  task 1.2 查询区块
func QueryBlockInfo(client *ethclient.Client, number *big.Int) {
	header, err := client.HeaderByNumber(context.Background(), number)
	if err != nil {
		fmt.Printf("HeaderByNumber err: %s", err.Error())
		return
	}
	fmt.Printf("HeaderByNumber: %s\n", header.Number.String())
	fmt.Printf("header Number: %v\n", header.Number.Uint64())
	fmt.Printf("header time: %v\n", header.Time)
	fmt.Printf("header Difficulty: %v\n", header.Difficulty.Uint64())
	fmt.Printf("header hash: %v\n", header.Hash().Hex())
	//address := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	//
	//balance, err := client.BalanceAt(context.Background(), address, nil)
	//if err != nil {
	//	fmt.Printf("Balance err: %s", err.Error())
	//	return
	//}
	//
	//fmt.Printf("Address %s Balance: %s \n", address.Hex(), balance.String())
	//
	//nonce, err := client.PendingNonceAt(context.Background(), address)
	//if err != nil {
	//	fmt.Printf("Pending nonce err: %s", err.Error())
	//	return
	//}
	//fmt.Printf("Nonce: %d\n", nonce)

	block, err := client.BlockByNumber(context.Background(), number)
	if err != nil {
		fmt.Printf("HeaderByNumber err: %s", err.Error())
		return
	}
	fmt.Printf("block Number: %v\n", block.Number().Uint64())
	fmt.Printf("block time: %v\n", block.Time())
	fmt.Printf("block Difficulty: %v\n", block.Difficulty().Uint64())
	fmt.Printf("block hash: %v\n", block.Hash().Hex())
	fmt.Printf("block transctionCount: %v\n", block.Transactions().Len())

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		fmt.Printf("TransactionCount err: %s", err.Error())
		return
	}

	fmt.Printf("TransactionCount: %v\n", count)
}

// NewTransaction 新建一个 ETH 转账  task 1.3 发送交易
func NewTransaction(client *ethclient.Client) {
	//转账交易包括打算转账的以太币数量，燃气限额，燃气价格，一个自增数(nonce)，接收地址以及可选择性的添加的数据。
	//在发送到以太坊网络之前，必须使用发送方的私钥对该交易进行签名。
	//	1. 加载私钥
	//baea89ffdbb9dd77c8737a154c87cad181eb79241f142c2e86c9889629de9307
	privateKey, err := crypto.HexToECDSA("baea89ffdbb9dd77c8737a154c87cad181eb79241f142c2e86c9889629de9307")
	if err != nil {
		fmt.Printf("NewTransaction HexToECDSA err: %s \n", err.Error())
		return
	}

	//	2. 获取账户的随机数(nonce)
	//每笔交易都需要一个 nonce。 根据定义，nonce 是仅使用一次的数字。
	//如果是发送交易的新账户，则该随机数将为“0”。 来自账户的每个新事务都必须具有前一个 nonce 增加 1 的 nonce。
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return
	}

	from := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		fmt.Printf("PendingNonceAt err: %s\n", err.Error())
		return
	}

	//	3. 设置要交易的ETH 数量，ETH 最小单位为 wei, 1ETH= 10^18wei
	//  这里要注意，应该先判断账号是否有足够的 balance 进行交易
	value := big.NewInt(10000000000000000) // 0.01 ETH

	//	4. 设置gasLimit, 设置的最大允许消耗 gas（防止意外烧钱）
	// 上限为 21000 单位 gas
	gasLimit := uint64(21000)

	//	5.设置gas价格，gas价格以wei为单位，一般通过SuggestGasPrice 获取。
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Printf("SuggestGasPrice err: %s\n", err.Error())
		return
	}

	//	6. 获取接受方的地址
	to := common.HexToAddress("0xc58Ea738A46DA06D2dB97a40262e3275291a4883")

	//  7. 生成一个未签名的事务,在进行智能合约交互时使用最后一个字段data
	//tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, nil)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     nil,
	})

	//  8. 使用发送方私钥对事务进行签名
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Printf("NetworkID err: %s\n", err.Error())
		return
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Printf("SignTx err: %s\n", err.Error())
		return
	}

	//	9. 将已签名的事务广播到整个网络
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Printf("SendTransaction err: %s\n", err.Error())
		return
	}

	fmt.Printf("SendTransaction: %v\n", signedTx.Hash().Hex())
	//	 交易结果 0x656b67ef7d0836b29945196588f7f346396b281d0315a82813478dc2cb51c700
}
