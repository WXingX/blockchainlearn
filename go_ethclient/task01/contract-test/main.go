package main

import (
	"context"
	"contract-test/counter"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
	"time"
)

const (
	HexKey       = ""
	BlockUrl     = "https://quaint-side-daylight.ethereum-sepolia.quiknode.pro/b2b5094e028ada682797d2268e7116e13e7cb8c9"
	ContractAddr = "0xb94bb57735fc725207145012C203EC1f6F30Ba5f"
)

func main() {
	client, err := ethclient.Dial(BlockUrl)
	if err != nil {
		fmt.Println("Error connecting to ethereum client")
		return
	}
	//DeployContract(client)
	//LoadContract(client)
	//CallAdd(client)
	CallAddByAbi(client)
}

// DeployContract 部署合约
func DeployContract(client *ethclient.Client) {
	// 1. 获取账户地址
	privateKey, err := crypto.HexToECDSA(HexKey)
	if err != nil {
		fmt.Println("HexToECDSA err. " + err.Error())
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("error casting public key to ECDSA")
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 2. 获取nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("PendingNonceAt err. " + err.Error())
		return
	}

	// 3. 获取 gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("SuggestGasPrice err. " + err.Error())
		return
	}

	// 4. 获取chainId
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("ChainID err. " + err.Error())
		return
	}

	// 5. 生成 transaction opt
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		fmt.Println("NewKeyedTransactorWithChainID err. " + err.Error())
		return
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	address, tx, _, err := counter.DeployCounter(auth, client)
	if err != nil {
		fmt.Println("DeployCounter err. " + err.Error())
		return
	}
	fmt.Println("DeployCounter address " + address.Hex()) //  0xb94bb57735fc725207145012C203EC1f6F30Ba5f
	fmt.Println("DeployCounter tx " + tx.Hash().Hex())    //  0x2c349478b164710a6c5b6a51d1765c8dfac0a74a550984e85b9637b7c35bbb1a
}

// LoadContract 加载合约
func LoadContract(client *ethclient.Client) {
	counterContract, err := counter.NewCounter(common.HexToAddress(ContractAddr), client)
	if err != nil {
		fmt.Println("NewCounter err. " + err.Error())
		return
	}

	_ = counterContract
}

// CallAdd 通过生成的go代码调用合约中的add方法
func CallAdd(client *ethclient.Client) {
	//  1. 创建合约示例
	counterContract, err := counter.NewCounter(common.HexToAddress(ContractAddr), client)
	if err != nil {
		fmt.Println("NewCounter err. " + err.Error())
		return
	}

	//	2.获取私钥
	privateKey, err := crypto.HexToECDSA(HexKey)
	if err != nil {
		fmt.Println("HexToECDSA err. " + err.Error())
		return
	}

	//  3. 获取chainId
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("ChainID err. " + err.Error())
		return
	}

	//	4. 初始化交易opt
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		fmt.Println("NewKeyedTransactorWithChainID err. " + err.Error())
		return
	}
	fmt.Println(opt)
	//	 5. 调用合约方法
	tx, err := counterContract.Add(opt)
	//count, err := counterContract.GetCount(nil)
	if err != nil {
		fmt.Println("counterContract add err. " + err.Error())
		return
	}
	//fmt.Println(count)
	fmt.Println("tx hash:", tx.Hash().Hex())
}

// CallAddByAbi 通过abi文件调用合约中的add方法
func CallAddByAbi(client *ethclient.Client) {
	// 1. 获取账户地址
	privateKey, err := crypto.HexToECDSA(HexKey)
	if err != nil {
		fmt.Println("HexToECDSA err. " + err.Error())
		return
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("error casting public key to ECDSA")
		return
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 2. 获取nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println("PendingNonceAt err. " + err.Error())
		return
	}

	// 3. 获取 gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("SuggestGasPrice err. " + err.Error())
		return
	}

	// 4. 获取chainId
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		fmt.Println("ChainID err. " + err.Error())
		return
	}

	// 5. 准备交易data
	contractABI, err := abi.JSON(strings.NewReader("[{\"inputs\":[],\"name\":\"add\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cnt\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCount\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"))
	if err != nil {
		fmt.Println("contractABI err. " + err.Error())
		return
	}
	methodName := "add"
	input, err := contractABI.Pack(methodName)
	if err != nil {
		fmt.Println("contractABI.Pack err. " + err.Error())
		return
	}

	// 6. 获取合约地址
	contractAddr := common.HexToAddress(ContractAddr)

	// 7. 生成 transaction
	//tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), 300000, gasPrice, input)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &contractAddr,
		Value:    big.NewInt(0),
		Gas:      300000, // gaslimit
		GasPrice: gasPrice,
		Data:     input,
	})

	// 8. 对 tx 进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Println("SignTx err. " + err.Error())
		return
	}

	// 9. 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Println("SendTransaction err. " + err.Error())
		return
	}
	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())

	//	10. 查询交易结果
	_, err = waitForReceipt(client, signedTx.Hash())
	if err != nil {
		fmt.Println("waitForReceipt err. " + err.Error())
		return
	}

	callInput, err := contractABI.Pack("getCount")
	if err != nil {
		fmt.Println("callInput contractABI.pack err. " + err.Error())
		return
	}

	callMsg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: callInput,
	}

	// 解析返回值
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		fmt.Println("CallContract err. " + err.Error())
		return
	}
	var count *big.Int
	err = contractABI.UnpackIntoInterface(&count, "getCount", result)
	if err != nil {
		fmt.Println("UnpackIntoInterface err. " + err.Error())
		return
	}
	fmt.Println("value:", count.String())
}

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if !errors.Is(err, ethereum.NotFound) {
			return nil, err
		}
		// 等待一段时间后再次查询
		time.Sleep(1 * time.Second)
	}
}
