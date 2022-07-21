package entities_test

import (
	core "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/vaulverin/uniswapv2-sdk/entities"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// nolint funlen
func TestTrade(t *testing.T) {
	ether := core.EtherOnChain(1)
	token0 := core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000001"), 18, "t0", "")
	token1 := core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000002"), 18, "t1", "")
	token2 := core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000003"), 18, "t2", "")
	token3 := core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000004"), 18, "t3", "")

	pair_0_1, _ := entities.NewPair(
		core.FromRawAmount(token0, big.NewInt(1000)),
		core.FromRawAmount(token1, big.NewInt(1000)),
		nil,
	)
	pair_0_2, _ := entities.NewPair(
		core.FromRawAmount(token0, big.NewInt(1000)),
		core.FromRawAmount(token2, big.NewInt(1100)),
		nil,
	)
	pair_0_3, _ := entities.NewPair(
		core.FromRawAmount(token0, big.NewInt(1000)),
		core.FromRawAmount(token3, big.NewInt(900)),
		nil,
	)
	pair_1_2, _ := entities.NewPair(
		core.FromRawAmount(token1, big.NewInt(1200)),
		core.FromRawAmount(token2, big.NewInt(1000)),
		nil,
	)
	pair_1_3, _ := entities.NewPair(
		core.FromRawAmount(token1, big.NewInt(1200)),
		core.FromRawAmount(token3, big.NewInt(1300)),
		nil,
	)

	pair_weth_0, _ := entities.NewPair(
		core.FromRawAmount(core.WETH9[1], big.NewInt(1000)),
		core.FromRawAmount(token0, big.NewInt(1000)),
		nil,
	)
	empty_pair_0_1, _ := entities.NewPair(
		core.FromRawAmount(token0, big.NewInt(0)),
		core.FromRawAmount(token1, big.NewInt(0)),
		nil,
	)

	{
		route, _ := entities.NewRoute([]*entities.Pair{pair_weth_0}, ether, token0)
		trade, _ := entities.NewTrade(
			route,
			core.FromRawAmount(core.EtherOnChain(1), big.NewInt(100)),
			entities.ExactInput,
		)

		// can be constructed with ETHER as input
		{
			expect := core.EtherOnChain(1)
			output := trade.InputAmount().Currency
			if !expect.Equal(output) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		{
			expect := token0
			output := trade.OutputAmount().Currency
			if !expect.Equal(output) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		// can be constructed with ETHER as input for exact output
		route, _ = entities.NewRoute([]*entities.Pair{pair_weth_0}, ether, token0)
		trade, _ = entities.NewTrade(
			route,
			core.FromRawAmount(token0, big.NewInt(100)),
			entities.ExactOutput,
		)
		{
			expect := core.EtherOnChain(1)
			output := trade.InputAmount().Currency
			if !expect.Equal(output) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		{
			expect := token0
			output := trade.OutputAmount().Currency
			if !expect.Equal(output) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		route, _ = entities.NewRoute([]*entities.Pair{pair_weth_0}, token0, ether)
		// can be constructed with ETHER as output
		trade, _ = entities.NewTrade(
			route,
			core.FromRawAmount(core.EtherOnChain(1), big.NewInt(100)),
			entities.ExactOutput,
		)
		{
			expect := token0
			output := trade.InputAmount().Currency
			if !expect.Equal(output) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		{
			expect := core.EtherOnChain(1)
			output := trade.OutputAmount().Currency
			if !expect.Equal(output) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		// can be constructed with ETHER as output for exact input
		trade, _ = entities.NewTrade(
			route,
			core.FromRawAmount(token0, big.NewInt(100)),
			entities.ExactInput,
		)
		{
			expect := token0
			output := trade.InputAmount().Currency
			if !expect.Equal(output) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		{
			expect := core.EtherOnChain(1)
			output := trade.OutputAmount().Currency
			if !expect.Equal(output) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
	}
	tokenAmount_0_100 := core.FromRawAmount(token0, big.NewInt(100))
	//entities.BestTradeExactIn
	{
		var pairs []*entities.Pair
		_, output := entities.BestTradeExactIn(pairs, tokenAmount_0_100, token2,
			entities.NewDefaultBestTradeOptions(), nil, tokenAmount_0_100, nil)
		//throws with empty pairs
		{
			expect := entities.ErrInvalidPairs
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		pairs = []*entities.Pair{pair_0_2}
		_, output = entities.BestTradeExactIn(pairs, tokenAmount_0_100, token2, &entities.BestTradeOptions{},
			nil, tokenAmount_0_100, nil)
		// throws with max hops of 0
		{
			expect := entities.ErrInvalidOption
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		pairs = []*entities.Pair{pair_0_1, pair_0_2, pair_1_2}
		result, err := entities.BestTradeExactIn(pairs, tokenAmount_0_100, token2,
			entities.NewDefaultBestTradeOptions(), nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		// provides best route
		{
			{
				var tests = []struct {
					expect int
					output int
				}{
					{2, len(result)},
					{1, len(result[0].Route.Pairs)},
					{2, len(result[1].Route.Pairs)},
				}
				for i, test := range tests {
					if test.expect != test.output {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
			{
				var tests = []struct {
					expect []*core.Token
					output []*core.Token
				}{
					{[]*core.Token{token0, token2}, result[0].Route.Path},
					{[]*core.Token{token0, token1, token2}, result[1].Route.Path},
				}
				for i, test := range tests {
					if len(test.expect) != len(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, len(test.expect), (test.output))
					}
					for j := range test.expect {
						if !test.expect[j].Equal(test.output[j]) {
							t.Errorf("test #%d#%d: expect[%+v], but got[%+v]", i, j, test.expect[j], test.output[j])
						}
					}
				}
			}

			{
				tokenAmount_2_99 := core.FromRawAmount(token2, big.NewInt(99))
				tokenAmount_2_69 := core.FromRawAmount(token2, big.NewInt(69))
				var tests = []struct {
					output *core.CurrencyAmount
					expect *core.CurrencyAmount
				}{
					{result[0].InputAmount(), tokenAmount_0_100},
					{result[0].OutputAmount(), tokenAmount_2_99},
					{result[1].InputAmount(), tokenAmount_0_100},
					{result[1].OutputAmount(), tokenAmount_2_69},
				}
				for i, test := range tests {
					if !test.expect.EqualTo(test.output.Fraction) {
						t.Error(test.output.Numerator, test.output.Denominator, test.output.Quotient())
						t.Error(test.expect.Numerator, test.expect.Denominator, test.expect.Quotient())
						t.Errorf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
		}

		// doesnt throw for zero liquidity pairs
		// throws with max hops of 0
		{
			pairs := []*entities.Pair{empty_pair_0_1}
			results, err := entities.BestTradeExactIn(pairs, tokenAmount_0_100, token1,
				entities.NewDefaultBestTradeOptions(), nil, tokenAmount_0_100, nil)
			if err != nil {
				t.Fatalf("err should be nil, got[%+v]", err)
			}
			expect := 0
			output := len(results)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		tokenAmount := core.FromRawAmount(token0, big.NewInt(10))
		result, _ = entities.BestTradeExactIn(pairs, tokenAmount, token2,
			&entities.BestTradeOptions{MaxNumResults: 3, MaxHops: 1}, nil, tokenAmount, nil)
		// respects maxHops
		{
			{
				var tests = []struct {
					expect int
					output int
				}{
					{1, len(result)},
					{1, len(result[0].Route.Pairs)},
				}
				for i, test := range tests {
					if test.expect != test.output {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
			{
				var tests = []struct {
					expect []*core.Token
					output []*core.Token
				}{
					{[]*core.Token{token0, token2}, result[0].Route.Path},
				}
				for i, test := range tests {
					if len(test.expect) != len(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, len(test.expect), (test.output))
					}
					for j := range test.expect {
						if !test.expect[j].Equal(test.output[j]) {
							t.Errorf("test #%d#%d: expect[%+v], but got[%+v]", i, j, test.expect[j], test.output[j])
						}
					}
				}
			}
		}

		tokenAmount = core.FromRawAmount(token0, big.NewInt(1))
		result, _ = entities.BestTradeExactIn(pairs, tokenAmount, token2,
			nil, nil, nil, nil)
		// insufficient input for one pair
		{
			{
				var tests = []struct {
					expect int
					output int
				}{
					{1, len(result)},
					{1, len(result[0].Route.Pairs)},
				}
				for i, test := range tests {
					if test.expect != test.output {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
			{
				var tests = []struct {
					expect []*core.Token
					output []*core.Token
				}{
					{[]*core.Token{token0, token2}, result[0].Route.Path},
				}
				for i, test := range tests {
					if len(test.expect) != len(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, len(test.expect), (test.output))
					}
					for j := range test.expect {
						if !test.expect[j].Equal(test.output[j]) {
							t.Errorf("test #%d#%d: expect[%+v], but got[%+v]", i, j, test.expect[j], test.output[j])
						}
					}
				}
			}
			{
				expect := core.FromRawAmount(token2, big.NewInt(1))
				output := result[0].OutputAmount()
				if !expect.EqualTo(output.Fraction) {
					t.Errorf("expect[%+v], but got[%+v]", expect, output)
				}
			}
		}

		tokenAmount = core.FromRawAmount(token0, big.NewInt(10))
		result, _ = entities.BestTradeExactIn(pairs, tokenAmount, token2,
			&entities.BestTradeOptions{MaxNumResults: 1, MaxHops: 3}, nil, nil, nil)
		// respects n
		{
			expect := 1
			output := len(result)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		pairs = []*entities.Pair{pair_0_1, pair_0_3, pair_1_3}
		result, _ = entities.BestTradeExactIn(pairs, tokenAmount, token2,
			&entities.BestTradeOptions{MaxNumResults: 1, MaxHops: 3}, nil, nil, nil)
		// no path
		{
			expect := 0
			output := len(result)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		tokenAmountETHER := core.FromRawAmount(core.EtherOnChain(1), big.NewInt(100))
		pairs = []*entities.Pair{pair_weth_0, pair_0_1, pair_0_3, pair_1_3}
		result, _ = entities.BestTradeExactIn(pairs, tokenAmountETHER, token3,
			nil, nil, nil, nil)
		// works for ETHER currency input
		{
			{
				expect := 2
				output := len(result)
				if expect != output {
					t.Fatalf("expect[%+v], but got[%+v]", expect, output)
				}
			}
			{
				var tests = []struct {
					expect core.Currency
					output core.Currency
				}{
					{ether, result[0].InputAmount().Currency},
					{token3, result[0].OutputAmount().Currency},
					{ether, result[1].InputAmount().Currency},
					{token3, result[1].OutputAmount().Currency},
				}
				for i, test := range tests {
					if !test.expect.Equal(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
			{
				var tests = []struct {
					expect []*core.Token
					output []*core.Token
				}{
					{[]*core.Token{ether.Wrapped(), token0, token1, token3}, result[0].Route.Path},
					{[]*core.Token{ether.Wrapped(), token0, token3}, result[1].Route.Path},
				}
				for i, test := range tests {
					if len(test.expect) != len(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, len(test.expect), (test.output))
					}
					for j := range test.expect {
						if !test.expect[j].Equal(test.output[j]) {
							t.Errorf("test #%d#%d: expect[%+v], but got[%+v]", i, j, test.expect[j], test.output[j])
						}
					}
				}
			}
		}

		tokenAmount = core.FromRawAmount(token3, big.NewInt(100))
		result, _ = entities.BestTradeExactIn(pairs, tokenAmount, ether,
			nil, nil, nil, nil)
		// works for ETHER currency output
		{
			{
				expect := 2
				output := len(result)
				if expect != output {
					t.Fatalf("expect[%+v], but got[%+v]", expect, output)
				}
			}
			{
				var tests = []struct {
					expect core.Currency
					output core.Currency
				}{
					{token3, result[0].InputAmount().Currency},
					{ether, result[0].OutputAmount().Currency},
					{token3, result[1].InputAmount().Currency},
					{ether, result[1].OutputAmount().Currency},
				}
				for i, test := range tests {
					if !test.expect.Equal(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
			{
				var tests = []struct {
					expect []*core.Token
					output []*core.Token
				}{
					{[]*core.Token{token3, token0, ether.Wrapped()}, result[0].Route.Path},
					{[]*core.Token{token3, token1, token0, ether.Wrapped()}, result[1].Route.Path},
				}
				for i, test := range tests {
					if len(test.expect) != len(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, len(test.expect), (test.output))
					}
					for j := range test.expect {
						if !test.expect[j].Equal(test.output[j]) {
							t.Errorf("test #%d#%d: expect[%+v], but got[%+v]", i, j, test.expect[j], test.output[j])
						}
					}
				}
			}
		}
	}

	// maximumAmountIn
	{
		// tradeType = EXACT_INPUT
		route, _ := entities.NewRoute([]*entities.Pair{pair_0_1, pair_1_2}, token0, nil)
		exactIn, _ := entities.ExactIn(route, tokenAmount_0_100)

		// throws if less than 0
		{
			percent := core.NewPercent(big.NewInt(-1), big.NewInt(100))
			_, output := exactIn.MaximumAmountIn(percent)
			expect := entities.ErrInvalidSlippageTolerance
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		// returns exact if 0
		{
			percent := core.NewPercent(big.NewInt(0), big.NewInt(100))
			output, _ := exactIn.MaximumAmountIn(percent)
			expect := exactIn.InputAmount()
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		// returns exact if nonzero
		{
			percent := core.NewPercent(big.NewInt(0), big.NewInt(100))
			output, _ := exactIn.MaximumAmountIn(percent)
			expect := core.FromRawAmount(token0, big.NewInt(100))
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}

			percent = core.NewPercent(big.NewInt(5), big.NewInt(100))
			output, _ = exactIn.MaximumAmountIn(percent)
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}

			percent = core.NewPercent(big.NewInt(200), big.NewInt(100))
			output, _ = exactIn.MaximumAmountIn(percent)
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		// tradeType = EXACT_OUTPUT
		tokenAmount := core.FromRawAmount(token2, big.NewInt(100))
		exactOut, _ := entities.ExactOut(route, tokenAmount)

		// throws if less than 0
		{
			percent := core.NewPercent(big.NewInt(-1), big.NewInt(100))
			_, output := exactOut.MaximumAmountIn(percent)
			expect := entities.ErrInvalidSlippageTolerance
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		// returns exact if 0
		{
			percent := core.NewPercent(big.NewInt(0), big.NewInt(100))
			output, _ := exactOut.MaximumAmountIn(percent)
			expect := exactOut.InputAmount()
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		// returns slippage amount if nonzero
		{
			percent := core.NewPercent(big.NewInt(0), big.NewInt(100))
			output, _ := exactOut.MaximumAmountIn(percent)
			expect := core.FromRawAmount(token0, big.NewInt(156))
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}

			percent = core.NewPercent(big.NewInt(5), big.NewInt(100))
			output, _ = exactOut.MaximumAmountIn(percent)
			expect = core.FromRawAmount(token0, big.NewInt(163))
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}

			percent = core.NewPercent(big.NewInt(200), big.NewInt(100))
			output, _ = exactOut.MaximumAmountIn(percent)
			expect = core.FromRawAmount(token0, big.NewInt(468))
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
	}

	// #minimumAmountOut
	{
		// tradeType = EXACT_INPUT
		route, _ := entities.NewRoute([]*entities.Pair{pair_0_1, pair_1_2}, token0, nil)
		exactIn, _ := entities.ExactIn(route, tokenAmount_0_100)

		// throws if less than 0
		{
			percent := core.NewPercent(big.NewInt(-1), big.NewInt(100))
			_, output := exactIn.MinimumAmountOut(percent)
			expect := entities.ErrInvalidSlippageTolerance
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		// returns exact if 0
		{
			percent := core.NewPercent(big.NewInt(0), big.NewInt(100))
			output, _ := exactIn.MinimumAmountOut(percent)
			expect := exactIn.OutputAmount()
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		// returns exact if nonzero
		{
			percent := core.NewPercent(big.NewInt(0), big.NewInt(100))
			output, _ := exactIn.MinimumAmountOut(percent)
			expect := core.FromRawAmount(token2, big.NewInt(69))
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}

			percent = core.NewPercent(big.NewInt(5), big.NewInt(100))
			output, _ = exactIn.MinimumAmountOut(percent)
			expect = core.FromRawAmount(token2, big.NewInt(65))
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}

			percent = core.NewPercent(big.NewInt(200), big.NewInt(100))
			output, _ = exactIn.MinimumAmountOut(percent)
			expect = core.FromRawAmount(token2, big.NewInt(23))
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		// tradeType = EXACT_OUTPUT
		tokenAmount := core.FromRawAmount(token2, big.NewInt(100))
		exactOut, _ := entities.ExactOut(route, tokenAmount)

		// throws if less than 0
		{
			percent := core.NewPercent(big.NewInt(-1), big.NewInt(100))
			_, output := exactOut.MinimumAmountOut(percent)
			expect := entities.ErrInvalidSlippageTolerance
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		// returns exact if 0
		{
			percent := core.NewPercent(big.NewInt(0), big.NewInt(100))
			output, _ := exactOut.MinimumAmountOut(percent)
			expect := exactOut.OutputAmount()
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
		// returns slippage amount if nonzero
		{
			percent := core.NewPercent(big.NewInt(0), big.NewInt(100))
			output, _ := exactOut.MinimumAmountOut(percent)
			expect := core.FromRawAmount(token2, big.NewInt(100))
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}

			percent = core.NewPercent(big.NewInt(5), big.NewInt(100))
			output, _ = exactOut.MinimumAmountOut(percent)
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}

			percent = core.NewPercent(big.NewInt(200), big.NewInt(100))
			output, _ = exactOut.MinimumAmountOut(percent)
			if !expect.EqualTo(output.Fraction) {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
	}

	// #bestTradeExactOut
	{
		pairs := []*entities.Pair{}
		tokenAmount_1_100 := core.FromRawAmount(token1, big.NewInt(100))
		tokenAmount_2_100 := core.FromRawAmount(token2, big.NewInt(100))
		_, output := entities.BestTradeExactOut(pairs, token2, tokenAmount_2_100,
			nil, nil, nil, nil)
		//throws with empty pairs
		{
			expect := entities.ErrInvalidPairs
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		pairs = []*entities.Pair{pair_0_2}
		_, output = entities.BestTradeExactOut(pairs, token0, tokenAmount_2_100,
			&entities.BestTradeOptions{MaxNumResults: 3}, nil, nil, nil)
		// throws with max hops of 0
		{
			expect := entities.ErrInvalidOption
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		pairs = []*entities.Pair{pair_0_1, pair_0_2, pair_1_2}
		result, _ := entities.BestTradeExactOut(pairs, token0, tokenAmount_2_100,
			nil, nil, nil, nil)
		// provides best route
		{
			{
				var tests = []struct {
					expect int
					output int
				}{
					{2, len(result)},
					{1, len(result[0].Route.Pairs)},
					{2, len(result[1].Route.Pairs)},
				}
				for i, test := range tests {
					if test.expect != test.output {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
			{
				var tests = []struct {
					expect []*core.Token
					output []*core.Token
				}{
					{[]*core.Token{token0, token2}, result[0].Route.Path},
					{[]*core.Token{token0, token1, token2}, result[1].Route.Path},
				}
				for i, test := range tests {
					if len(test.expect) != len(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, len(test.expect), (test.output))
					}
					for j := range test.expect {
						if !test.expect[j].Equal(test.output[j]) {
							t.Errorf("test #%d#%d: expect[%+v], but got[%+v]", i, j, test.expect[j], test.output[j])
						}
					}
				}
			}

			{
				tokenAmount_0_101 := core.FromRawAmount(token0, big.NewInt(101))
				tokenAmount_0_156 := core.FromRawAmount(token0, big.NewInt(156))
				var tests = []struct {
					expect *core.CurrencyAmount
					output *core.CurrencyAmount
				}{
					{result[0].InputAmount(), tokenAmount_0_101},
					{result[0].OutputAmount(), tokenAmount_2_100},
					{result[1].InputAmount(), tokenAmount_0_156},
					{result[1].OutputAmount(), tokenAmount_2_100},
				}
				for i, test := range tests {
					if !test.expect.EqualTo(test.output.Fraction) {
						t.Errorf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
		}

		// doesnt throw for zero liquidity pairs
		{
			pairs := []*entities.Pair{empty_pair_0_1}
			results, err := entities.BestTradeExactOut(pairs, token1, tokenAmount_1_100,
				nil, nil, nil, nil)
			if err != nil {
				t.Fatalf("err should be nil, got[%+v]", err)
			}
			expect := 0
			output := len(results)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		tokenAmount := core.FromRawAmount(token2, big.NewInt(10))
		result, _ = entities.BestTradeExactOut(pairs, token0, tokenAmount,
			&entities.BestTradeOptions{MaxNumResults: 3, MaxHops: 1}, nil, nil, nil)
		// respects maxHops
		{
			{
				var tests = []struct {
					expect int
					output int
				}{
					{1, len(result)},
					{1, len(result[0].Route.Pairs)},
				}
				for i, test := range tests {
					if test.expect != test.output {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
			{
				var tests = []struct {
					expect []*core.Token
					output []*core.Token
				}{
					{[]*core.Token{token0, token2}, result[0].Route.Path},
				}
				for i, test := range tests {
					if len(test.expect) != len(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, len(test.expect), (test.output))
					}
					for j := range test.expect {
						if !test.expect[j].Equal(test.output[j]) {
							t.Errorf("test #%d#%d: expect[%+v], but got[%+v]", i, j, test.expect[j], test.output[j])
						}
					}
				}
			}
		}

		tokenAmount = core.FromRawAmount(token2, big.NewInt(1200))
		result, _ = entities.BestTradeExactOut(pairs, token0, tokenAmount,
			nil, nil, nil, nil)
		// insufficient liquidity
		{
			expect := 0
			output := len(result)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		tokenAmount = core.FromRawAmount(token2, big.NewInt(1050))
		result, _ = entities.BestTradeExactOut(pairs, token0, tokenAmount,
			nil, nil, nil, nil)
		// insufficient liquidity in one pair but not the other
		{
			expect := 1
			output := len(result)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		tokenAmount = core.FromRawAmount(token2, big.NewInt(10))
		result, _ = entities.BestTradeExactOut(pairs, token0, tokenAmount,
			&entities.BestTradeOptions{MaxNumResults: 1, MaxHops: 3}, nil, nil, nil)
		// respects n
		{
			expect := 1
			output := len(result)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		pairs = []*entities.Pair{pair_0_1, pair_0_3, pair_1_3}
		result, _ = entities.BestTradeExactOut(pairs, token0, tokenAmount,
			nil, nil, nil, nil)
		// no path
		{
			expect := 0
			output := len(result)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		pairs = []*entities.Pair{pair_weth_0, pair_0_1, pair_0_3, pair_1_3}
		tokenAmount = core.FromRawAmount(token3, big.NewInt(100))
		result, _ = entities.BestTradeExactOut(pairs, ether, tokenAmount,
			nil, nil, nil, nil)
		// works for ETHER currency input
		{
			{
				expect := 2
				output := len(result)
				if expect != output {
					t.Fatalf("expect[%+v], but got[%+v]", expect, output)
				}
			}
			{
				var tests = []struct {
					expect core.Currency
					output core.Currency
				}{
					{ether, result[0].InputAmount().Currency},
					{token3, result[0].OutputAmount().Currency},
					{ether, result[1].InputAmount().Currency},
					{token3, result[1].OutputAmount().Currency},
				}
				for i, test := range tests {
					if !test.expect.Equal(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
			{
				var tests = []struct {
					expect []*core.Token
					output []*core.Token
				}{
					{[]*core.Token{ether.Wrapped(), token0, token1, token3}, result[0].Route.Path},
					{[]*core.Token{ether.Wrapped(), token0, token3}, result[1].Route.Path},
				}
				for i, test := range tests {
					if len(test.expect) != len(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, len(test.expect), (test.output))
					}
					for j := range test.expect {
						if !test.expect[j].Equal(test.output[j]) {
							t.Errorf("test #%d#%d: expect[%+v], but got[%+v]", i, j, test.expect[j], test.output[j])
						}
					}
				}
			}
		}

		tokenAmount = core.FromRawAmount(ether, big.NewInt(100))
		result, _ = entities.BestTradeExactOut(pairs, token3, tokenAmount,
			nil, nil, nil, nil)
		// works for ETHER currency output
		{
			{
				expect := 2
				output := len(result)
				if expect != output {
					t.Fatalf("expect[%+v], but got[%+v]", expect, output)
				}
			}
			{
				var tests = []struct {
					expect core.Currency
					output core.Currency
				}{
					{token3, result[0].InputAmount().Currency},
					{ether, result[0].OutputAmount().Currency},
					{token3, result[1].InputAmount().Currency},
					{ether, result[1].OutputAmount().Currency},
				}
				for i, test := range tests {
					if !test.expect.Equal(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, test.expect, test.output)
					}
				}
			}
			{
				var tests = []struct {
					expect []*core.Token
					output []*core.Token
				}{
					{[]*core.Token{token3, token0, ether.Wrapped()}, result[0].Route.Path},
					{[]*core.Token{token3, token1, token0, ether.Wrapped()}, result[1].Route.Path},
				}
				for i, test := range tests {
					if len(test.expect) != len(test.output) {
						t.Fatalf("test #%d: expect[%+v], but got[%+v]", i, len(test.expect), (test.output))
					}
					for j := range test.expect {
						if !test.expect[j].Equal(test.output[j]) {
							t.Errorf("test #%d#%d: expect[%+v], but got[%+v]", i, j, test.expect[j], test.output[j])
						}
					}
				}
			}
		}
	}
}
