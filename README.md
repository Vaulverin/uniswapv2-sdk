# Uniswap V2 SDK

[![API Reference](https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667)](https://pkg.go.dev/github.com/vaulverin/uniswapv2-sdk)
[![Test](https://github.com/vaulverin/uniswapv2-sdk/actions/workflows/test.yml/badge.svg)](https://github.com/vaulverin/uniswapv2-sdk/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/vaulverin/uniswapv2-sdk)](https://goreportcard.com/report/github.com/vaulverin/uniswapv2-sdk)

🛠 A Go SDK for building applications on top of Uniswap V2

## Installation

```sh
go get github.com/vaulverin/uniswapv2-sdk
```

## Usage
```go
package main

import (
	core "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vaulverin/uniswapv2-sdk/entities"
	"github.com/vaulverin/uniswapv2-sdk/router"
	"log"
	"math/big"
)

func main() {
	ether = core.EtherOnChain(1)
	token0 = core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000001"), 18, "t0", "t0")
	token1 = core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000002"), 18, "t1", "t1")

	pair_0_1, _ = entities.NewPair(
		core.FromRawAmount(token0, big.NewInt(1000)),
		core.FromRawAmount(token1, big.NewInt(1000)),
		nil,
	)
	pair_weth_0, _ = entities.NewPair(
		core.FromRawAmount(core.WETH9[1], big.NewInt(1000)),
		core.FromRawAmount(token0, big.NewInt(1000)),
		nil,
	)

	route, _ := entities.NewRoute([]*entities.Pair{pair_weth_0, pair_0_1}, ether, token1)
	trade, _ := entities.ExactIn(route, core.FromRawAmount(ether, big.NewInt(100)))
	swapParams, _ := router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: core.NewPercent(big.NewInt(1), big.NewInt(100)),
		Recipient:       common.HexToAddress("0x0000000000000000000000000000000000000004"),
		Deadline:        big.NewInt(time.Now().Add(30 * time.Second).Unix()),
	})

	log.Println(swapParams.MethodName)
	log.Println(swapParams.Args)
	log.Println(swapParams.Value)
}

```
Look for more examples in tests.