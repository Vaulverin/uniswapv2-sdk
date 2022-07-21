package entities_test

import (
	core "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vaulverin/uniswapv2-sdk/entities"
	"math/big"
	"testing"
)

func TestRoute(t *testing.T) {
	ether := core.EtherOnChain(1)
	token0 := core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000001"), 18, "t0", "t0")
	token1 := core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000002"), 18, "t1", "t1")
	weth := core.WETH9[1]
	b100 := big.NewInt(100)

	pair01, err := entities.NewPair(
		core.FromRawAmount(token0, b100),
		core.FromRawAmount(token1, big.NewInt(200)),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	pair0Weth, err := entities.NewPair(
		core.FromRawAmount(token0, b100),
		core.FromRawAmount(weth, b100),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	pair1Weth, err := entities.NewPair(
		core.FromRawAmount(token1, big.NewInt(175)),
		core.FromRawAmount(weth, b100),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	// constructs a path from the tokens
	{
		route, err := entities.NewRoute([]*entities.Pair{pair01}, token0, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(route.Pairs) != 1 || route.Pairs[0] != pair01 {
			t.Error("wrong pairs for route")
		}
		if len(route.Path) != 2 || route.Path[0] != token0 || route.Path[1] != token1 {
			t.Error("wrong path for route")
		}
		if route.Input != token0 {
			t.Error("wrong input for route")
		}
		if route.Output != token1 {
			t.Error("wrong output for route")
		}
		if route.ChainID() != 1 {
			t.Error("wrong chain id for route")
		}
	}

	// can have a token as both input and output
	{
		pairs := []*entities.Pair{pair0Weth, pair01, pair1Weth}
		route, err := entities.NewRoute(pairs, weth, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(route.Pairs) != len(pairs) {
			t.Fatal("wrong pairs for route")
		}
		for i, pair := range route.Pairs {
			if pair != pairs[i] {
				t.Error("wrong pairs for route")
			}
		}
		if route.Input != weth {
			t.Error("wrong input for route")
		}
		if route.Output != weth {
			t.Error("wrong output for route")
		}
	}

	{
		// supports ether input
		pairs := []*entities.Pair{pair0Weth}
		route, err := entities.NewRoute(pairs, ether, token0)
		if err != nil {
			t.Fatal(err)
		}
		if len(route.Pairs) != len(pairs) {
			t.Fatal("wrong pairs for route")
		}
		for i, pair := range route.Pairs {
			if pair != pairs[i] {
				t.Error("wrong pairs for route")
			}
		}
		if !route.Input.Equal(ether) {
			t.Error("wrong input for route")
		}
		if !route.Output.Equal(token0) {
			t.Error("wrong output for route")
		}
	}

	{
		// supports ether output
		pairs := []*entities.Pair{pair0Weth}
		route, err := entities.NewRoute(pairs, token0, ether)
		if err != nil {
			t.Fatal(err)
		}
		if len(route.Pairs) != len(pairs) {
			t.Fatal("wrong pairs for route")
		}
		for i, pair := range route.Pairs {
			if pair != pairs[i] {
				t.Error("wrong pairs for route")
			}
		}
		if !route.Input.Equal(token0) {
			t.Error("wrong input for route")
		}
		if !route.Output.Equal(ether) {
			t.Error("wrong output for route")
		}
	}
}
