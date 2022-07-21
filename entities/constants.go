package entities

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	FactoryAddress = common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")
	InitCodeHash   = common.FromHex("0x96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f")

	MinimumLiquidity = big.NewInt(1000)

	Zero  = big.NewInt(0)
	One   = big.NewInt(1)
	Five  = big.NewInt(5)
	B997  = big.NewInt(997)
	B1000 = big.NewInt(1000)
)

type TradeType int

const (
	ExactInput TradeType = iota
	ExactOutput
)
