package router_test

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
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
	})

	log.Println(swapParams.MethodName)
	log.Println(swapParams.Args)
	log.Println(swapParams.Value)
}
