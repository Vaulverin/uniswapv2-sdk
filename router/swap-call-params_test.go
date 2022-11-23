package router_test

import (
	core "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vaulverin/uniswapv2-sdk/entities"
	"github.com/vaulverin/uniswapv2-sdk/router"
	"math/big"
	"reflect"
	"testing"
	"time"
)

var (
	recipient = common.HexToAddress("0x0000000000000000000000000000000000000004")
	ether     = core.EtherOnChain(1)
	token0    = core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000001"), 18, "t0", "t0")
	token1    = core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000002"), 18, "t1", "t1")

	amount0   = core.FromRawAmount(token0, big.NewInt(1000))
	amount1   = core.FromRawAmount(token1, big.NewInt(1000))
	ethAmount = core.FromRawAmount(core.WETH9[1], big.NewInt(1000))

	pair_0_1, _    = entities.NewPair(amount0, amount1, nil)
	pair_weth_0, _ = entities.NewPair(ethAmount, amount0, nil)

	slippage = core.NewPercent(big.NewInt(1), big.NewInt(100))
	deadline = big.NewInt(time.Now().Add(30 * time.Second).Unix())

	testNumber = 0
)

func check(t *testing.T, expect interface{}, output interface{}) {
	if !reflect.DeepEqual(output, expect) {
		t.Errorf("test #%d: FAILED expect[%+v], but got[%+v]", testNumber, expect, output)
	} else {
		t.Logf("test #%d: PASSED", testNumber)
	}
	testNumber++
}

func TestExactInEtherToToken1(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_weth_0, pair_0_1}, ether, token1)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(ether, big.NewInt(100))
	trade, err := entities.ExactIn(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	swapParams, err := router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
	})
	check(t, "swapExactETHForTokens", swapParams.MethodName)
	check(t, hexutil.MustDecodeBig("0x51"), swapParams.Args[0])
	check(t, []common.Address{core.WETH9[1].Address, token0.Address, token1.Address}, swapParams.Args[1])
	check(t, recipient, swapParams.Args[2])
	check(t, deadline, swapParams.Args[3])
	check(t, hexutil.MustDecodeBig("0x64"), swapParams.Value)
}

func TestExactInToken1ToEther(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_0_1, pair_weth_0}, token1, ether)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(token1, big.NewInt(100))
	trade, err := entities.ExactIn(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	swapParams, err := router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
	})
	check(t, "swapExactTokensForETH", swapParams.MethodName)
	check(t, hexutil.MustDecodeBig("0x64"), swapParams.Args[0])
	check(t, hexutil.MustDecodeBig("0x51"), swapParams.Args[1])
	check(t, []common.Address{token1.Address, token0.Address, core.WETH9[1].Address}, swapParams.Args[2])
	check(t, recipient, swapParams.Args[3])
	check(t, deadline, swapParams.Args[4])
	check(t, big.NewInt(0), swapParams.Value)
}

func TestExactInToken0ToToken1(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_0_1}, token0, token1)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(token0, big.NewInt(100))
	trade, err := entities.ExactIn(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	swapParams, err := router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
	})
	check(t, "swapExactTokensForTokens", swapParams.MethodName)
	check(t, hexutil.MustDecodeBig("0x64"), swapParams.Args[0])
	check(t, hexutil.MustDecodeBig("0x59"), swapParams.Args[1])
	check(t, []common.Address{token0.Address, token1.Address}, swapParams.Args[2])
	check(t, recipient, swapParams.Args[3])
	check(t, deadline, swapParams.Args[4])
	check(t, big.NewInt(0), swapParams.Value)
}

func TestExactOutEtherToToken1(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_weth_0, pair_0_1}, ether, token1)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(token1, big.NewInt(100))
	trade, err := entities.ExactOut(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	swapParams, err := router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
	})
	check(t, "swapETHForExactTokens", swapParams.MethodName)
	check(t, hexutil.MustDecodeBig("0x64"), swapParams.Args[0])
	check(t, []common.Address{core.WETH9[1].Address, token0.Address, token1.Address}, swapParams.Args[1])
	check(t, recipient, swapParams.Args[2])
	check(t, deadline, swapParams.Args[3])
	check(t, hexutil.MustDecodeBig("0x80"), swapParams.Value)
}

func TestExactOutToken1ToEther(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_0_1, pair_weth_0}, token1, ether)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(ether, big.NewInt(100))
	trade, err := entities.ExactOut(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	swapParams, err := router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
	})
	check(t, "swapTokensForExactETH", swapParams.MethodName)
	check(t, hexutil.MustDecodeBig("0x64"), swapParams.Args[0])
	check(t, hexutil.MustDecodeBig("0x80"), swapParams.Args[1])
	check(t, []common.Address{token1.Address, token0.Address, core.WETH9[1].Address}, swapParams.Args[2])
	check(t, recipient, swapParams.Args[3])
	check(t, deadline, swapParams.Args[4])
	check(t, big.NewInt(0), swapParams.Value)
}

func TestExactOutToken0ToToken1(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_0_1}, token0, token1)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(token1, big.NewInt(100))
	trade, err := entities.ExactOut(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	swapParams, err := router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
	})
	check(t, "swapTokensForExactTokens", swapParams.MethodName)
	check(t, hexutil.MustDecodeBig("0x64"), swapParams.Args[0])
	check(t, hexutil.MustDecodeBig("0x71"), swapParams.Args[1])
	check(t, []common.Address{token0.Address, token1.Address}, swapParams.Args[2])
	check(t, recipient, swapParams.Args[3])
	check(t, deadline, swapParams.Args[4])
	check(t, big.NewInt(0), swapParams.Value)
}

func TestFeeOnTransferExactInEtherToToken1(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_weth_0, pair_0_1}, ether, token1)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(ether, big.NewInt(100))
	trade, err := entities.ExactIn(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	swapParams, err := router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
		FeeOnTransfer:   true,
	})
	check(t, "swapExactETHForTokensSupportingFeeOnTransferTokens", swapParams.MethodName)
	check(t, hexutil.MustDecodeBig("0x51"), swapParams.Args[0])
	check(t, []common.Address{core.WETH9[1].Address, token0.Address, token1.Address}, swapParams.Args[1])
	check(t, recipient, swapParams.Args[2])
	check(t, deadline, swapParams.Args[3])
	check(t, hexutil.MustDecodeBig("0x64"), swapParams.Value)
}

func TestFeeOnTransferExactInToken1ToEther(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_0_1, pair_weth_0}, token1, ether)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(token1, big.NewInt(100))
	trade, err := entities.ExactIn(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	swapParams, err := router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
		FeeOnTransfer:   true,
	})
	check(t, "swapExactTokensForETHSupportingFeeOnTransferTokens", swapParams.MethodName)
	check(t, hexutil.MustDecodeBig("0x64"), swapParams.Args[0])
	check(t, hexutil.MustDecodeBig("0x51"), swapParams.Args[1])
	check(t, []common.Address{token1.Address, token0.Address, core.WETH9[1].Address}, swapParams.Args[2])
	check(t, recipient, swapParams.Args[3])
	check(t, deadline, swapParams.Args[4])
	check(t, big.NewInt(0), swapParams.Value)
}

func TestFeeOnTransferExactInToken0ToToken1(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_0_1}, token0, token1)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(token0, big.NewInt(100))
	trade, err := entities.ExactIn(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	swapParams, err := router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
		FeeOnTransfer:   true,
	})
	check(t, "swapExactTokensForTokensSupportingFeeOnTransferTokens", swapParams.MethodName)
	check(t, hexutil.MustDecodeBig("0x64"), swapParams.Args[0])
	check(t, hexutil.MustDecodeBig("0x59"), swapParams.Args[1])
	check(t, []common.Address{token0.Address, token1.Address}, swapParams.Args[2])
	check(t, recipient, swapParams.Args[3])
	check(t, deadline, swapParams.Args[4])
	check(t, big.NewInt(0), swapParams.Value)
}

func TestFeeOnTransferExactOutEtherToToken1(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_weth_0, pair_0_1}, ether, token1)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(token1, big.NewInt(100))
	trade, err := entities.ExactOut(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	_, err = router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
		FeeOnTransfer:   true,
	})
	check(t, router.ErrExactOutFot, err)
}

func TestFeeOnTransferExactOutToken1ToEther(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_0_1, pair_weth_0}, token1, ether)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(ether, big.NewInt(100))
	trade, err := entities.ExactOut(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	_, err = router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
		FeeOnTransfer:   true,
	})
	check(t, router.ErrExactOutFot, err)
}

func TestFeeOnTransferExactOutToken0ToToken1(t *testing.T) {
	testNumber = 0
	route, err := entities.NewRoute([]*entities.Pair{pair_0_1}, token0, token1)
	if err != nil {
		t.Fatal(err)
	}
	amount := core.FromRawAmount(token1, big.NewInt(100))
	trade, err := entities.ExactOut(route, amount)
	if err != nil {
		t.Fatal(err)
	}
	_, err = router.SwapCallParameters(trade, router.TradeOptions{
		AllowedSlippage: slippage,
		Recipient:       recipient,
		Deadline:        deadline,
		FeeOnTransfer:   true,
	})
	check(t, router.ErrExactOutFot, err)
}
