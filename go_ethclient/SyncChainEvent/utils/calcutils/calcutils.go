package calcutils

import (
	"math/big"
)

// BigIntAdd 计算多个大数字符串的加法
func BigIntAdd(args ...string) string {
	if len(args) == 0 {
		return "0"
	}
	sum, _ := big.NewInt(0).SetString(args[0], 10)
	for _, arg := range args[1:] {
		n, _ := big.NewInt(0).SetString(arg, 10)
		sum.Add(sum, n)
	}
	return sum.String()
}

// BigIntSub 计算两个大数字符串的减法
func BigIntSub(a, b string) string {
	bigA, _ := big.NewInt(0).SetString(a, 10)
	bigB, _ := big.NewInt(0).SetString(b, 10)
	return bigA.Sub(bigA, bigB).String()
}

// BigIntMul 计算一个字符串和一个浮点数的乘法
func BigIntMul(a string, b float64) string {
	bigA, _ := big.NewInt(0).SetString(a, 10)
	bigB := big.NewFloat(b)
	bigC := big.NewFloat(0).Mul(big.NewFloat(0).SetInt(bigA), bigB)
	return bigC.Text('f', -1)
}
