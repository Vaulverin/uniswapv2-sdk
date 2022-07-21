package entities_test

import (
	core "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vaulverin/uniswapv2-sdk/entities"
	"math/big"
	"strings"
	"testing"
)

var (
	USDC = core.NewToken(1, common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), 18, "USDC", "USD Coin")
	DAI  = core.NewToken(1, common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"), 18, "DAI", "DAI Stablecoin")
	B100 = big.NewInt(100)
)

func TestGetAddress(t *testing.T) {
	var tests = []struct {
		Input  [2]*core.Token
		Output string
	}{
		{
			[2]*core.Token{USDC, DAI},
			"0xb50b5182D6a47EC53a469395AF44e371d7C76ed4",
		},
		{
			[2]*core.Token{DAI, USDC},
			"0xb50b5182D6a47EC53a469395AF44e371d7C76ed4",
		},
	}
	factory := common.HexToAddress("0x1111111111111111111111111111111111111111")
	for i, test := range tests {
		output, err := entities.GetAddress(test.Input[0], test.Input[1], factory, entities.InitCodeHash)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.EqualFold(output.Hex(), test.Output) {
			t.Errorf("test #%d: failed to match when it should (%s != %s)", i, output, test.Output)
		}
	}
}

func TestPair(t *testing.T) {
	tokenAmountUSDC := core.FromRawAmount(USDC, B100)
	tokenAmountDAI := core.FromRawAmount(DAI, B100)
	tokenAmountUSDC101 := core.FromRawAmount(USDC, big.NewInt(101))
	tokenAmountDAI101 := core.FromRawAmount(DAI, big.NewInt(101))
	usdcDaiPairAddress := common.HexToAddress("0xAE461cA67B15dc8dc81CE7615e0320dA1A9aB8D5")

	// cannot be used for tokens on different chains
	{
		tokenAmountB := core.FromRawAmount(core.WETH9[4], B100)
		_, output := entities.NewCurrencyAmounts(tokenAmountUSDC, tokenAmountB)
		expect := core.ErrDifferentChain
		if expect != output {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
	}

	// returns the correct address
	{
		output, err := entities.GetAddress(DAI, USDC, entities.FactoryAddress, entities.InitCodeHash)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.EqualFold(output.Hex(), usdcDaiPairAddress.Hex()) {
			t.Errorf("expect[%+v], but got[%+v]", usdcDaiPairAddress, output)
		}
	}

	{
		pairA, _ := entities.NewPair(tokenAmountUSDC, tokenAmountDAI, nil)
		pairB, _ := entities.NewPair(tokenAmountDAI, tokenAmountUSDC, nil)
		expect := DAI
		// always is the token that sorts before
		output := pairA.Token0()
		if !expect.Equal(output) {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
		output = pairB.Token0()
		if !expect.Equal(output) {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}

		expect = USDC
		// always is the token that sorts after
		output = pairA.Token1()
		if !expect.Equal(output) {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
		output = pairB.Token1()
		if !expect.Equal(output) {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
	}

	{
		pairA, _ := entities.NewPair(tokenAmountUSDC, tokenAmountDAI101, nil)
		pairB, _ := entities.NewPair(tokenAmountDAI101, tokenAmountUSDC, nil)
		expect := tokenAmountDAI101
		// always comes from the token that sorts before
		output := pairA.Reserve0()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output = pairB.Reserve0()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		expect = tokenAmountUSDC
		// always comes from the token that sorts after
		output = pairA.Reserve1()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output = pairB.Reserve1()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
	}

	{
		pairA, _ := entities.NewPair(tokenAmountUSDC101, tokenAmountDAI, nil)
		pairB, _ := entities.NewPair(tokenAmountDAI, tokenAmountUSDC101, nil)
		b100 := big.NewInt(100)
		b101 := big.NewInt(101)
		expect := core.NewPrice(DAI, USDC, b100, b101)
		// returns price of token0 in terms of token1
		output := pairA.Token0Price()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output = pairB.Token0Price()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		expect = core.NewPrice(USDC, DAI, b101, b100)
		// returns price of token1 in terms of token0
		output = pairA.Token1Price()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output = pairB.Token1Price()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
	}

	{
		p, _ := entities.NewPair(tokenAmountUSDC101, tokenAmountDAI, nil)
		// returns price of token in terms of other token
		expect := p.Token0Price()
		output, _ := p.PriceOf(tokenAmountDAI.Currency.Wrapped())
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		expect = p.Token1Price()
		output, _ = p.PriceOf(tokenAmountUSDC101.Currency.Wrapped())
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		{
			// throws if invalid token
			expect := entities.ErrDiffToken
			_, output := p.PriceOf(core.WETH9[1])
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
	}

	{
		pairA, _ := entities.NewPair(tokenAmountUSDC, tokenAmountDAI101, nil)
		pairB, _ := entities.NewPair(tokenAmountDAI101, tokenAmountUSDC, nil)
		expect := tokenAmountUSDC
		// returns reserves of the given token
		output, _ := pairA.ReserveOf(USDC)
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output, _ = pairB.ReserveOf(USDC)
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		expect = tokenAmountUSDC
		// always comes from the token that sorts after
		output = pairA.Reserve1()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output = pairB.Reserve1()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		{
			// throws if not in the pair
			expect := entities.ErrDiffToken
			_, output := pairB.ReserveOf(core.WETH9[1])
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
	}

	{
		pairA, _ := entities.NewPair(tokenAmountUSDC, tokenAmountDAI, nil)
		pairB, _ := entities.NewPair(tokenAmountDAI, tokenAmountUSDC, nil)
		expect := uint(1)
		// returns the token0 chainId
		output := pairA.ChainID()
		if expect != output {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
		output = pairB.ChainID()
		if expect != output {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}

		{
			expect := true
			// involvesToken
			output := pairA.InvolvesToken(USDC)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
			output = pairA.InvolvesToken(DAI)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
			expect = false
			output = pairA.InvolvesToken(core.WETH9[1])
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		{
			tokenA := core.NewToken(4, common.HexToAddress("0x0000000000000000000000000000000000000001"), 18, "", "")
			tokenB := core.NewToken(4, common.HexToAddress("0x0000000000000000000000000000000000000002"), 18, "", "")
			tokenAmountA := core.FromRawAmount(tokenA, big.NewInt(0))
			tokenAmountB := core.FromRawAmount(tokenB, big.NewInt(0))
			p, _ := entities.NewPair(tokenAmountA, tokenAmountB, nil)
			{
				tokenAmount := core.FromRawAmount(p.LiquidityToken, big.NewInt(0))
				tokenAmountA := core.FromRawAmount(tokenA, big.NewInt(1000))
				tokenAmountB := core.FromRawAmount(tokenB, big.NewInt(1000))
				// getLiquidityMinted:0
				expect := entities.ErrInsufficientInputAmount
				_, output := p.GetLiquidityMinted(tokenAmount, tokenAmountA, tokenAmountB)
				if expect != output {
					t.Errorf("expect[%+v], but got[%+v]", expect, output)
				}

				tokenAmountA = core.FromRawAmount(tokenA, big.NewInt(1000000))
				tokenAmountB = core.FromRawAmount(tokenB, big.NewInt(1))
				_, output = p.GetLiquidityMinted(tokenAmount, tokenAmountA, tokenAmountB)
				if expect != output {
					t.Errorf("expect[%+v], but got[%+v]", expect, output)
				}

				tokenAmountA = core.FromRawAmount(tokenA, big.NewInt(1001))
				tokenAmountB = core.FromRawAmount(tokenB, big.NewInt(1001))
				{
					expect := "1"
					liquidity, _ := p.GetLiquidityMinted(tokenAmount, tokenAmountA, tokenAmountB)
					output := liquidity.Quotient().String()
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
			}

			// getLiquidityMinted:!0
			tokenAmountA = core.FromRawAmount(tokenA, big.NewInt(10000))
			tokenAmountB = core.FromRawAmount(tokenB, big.NewInt(10000))
			p, _ = entities.NewPair(tokenAmountA, tokenAmountB, nil)
			{
				tokenAmount := core.FromRawAmount(p.LiquidityToken, big.NewInt(10000))
				tokenAmountA = core.FromRawAmount(tokenA, big.NewInt(2000))
				tokenAmountB = core.FromRawAmount(tokenB, big.NewInt(2000))
				expect := "2000"
				liquidity, _ := p.GetLiquidityMinted(tokenAmount, tokenAmountA, tokenAmountB)
				output := liquidity.Quotient().String()
				if expect != output {
					t.Errorf("expect[%+v], but got[%+v]", expect, output)
				}
			}

			// getLiquidityValue:!feeOn
			tokenAmountA = core.FromRawAmount(tokenA, big.NewInt(1000))
			tokenAmountB = core.FromRawAmount(tokenB, big.NewInt(1000))
			p, _ = entities.NewPair(tokenAmountA, tokenAmountB, nil)
			tokenAmount := core.FromRawAmount(p.LiquidityToken, big.NewInt(1000))
			tokenAmount500 := core.FromRawAmount(p.LiquidityToken, big.NewInt(500))
			{
				liquidityValue, _ := p.GetLiquidityValue(tokenA, tokenAmount, tokenAmount, false, nil)
				{
					expect := true
					output := liquidityValue.Currency.Wrapped().Equal(tokenA)
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
				{
					expect := "1000"
					output := liquidityValue.Quotient().String()
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}

				liquidityValue, _ = p.GetLiquidityValue(tokenA, tokenAmount, tokenAmount500, false, nil)
				// 500
				{
					expect := true
					output := liquidityValue.Currency.Wrapped().Equal(tokenA)
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
				{
					expect := "500"
					output := liquidityValue.Quotient().String()
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}

				liquidityValue, _ = p.GetLiquidityValue(tokenB, tokenAmount, tokenAmount, false, nil)
				// tokenB
				{
					expect := true
					output := liquidityValue.Currency.Wrapped().Equal(tokenB)
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
				{
					expect := "1000"
					output := liquidityValue.Quotient().String()
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
			}

			// getLiquidityValue:feeOn
			{
				liquidityValue, _ := p.GetLiquidityValue(tokenA, tokenAmount500, tokenAmount500, true, big.NewInt(250000))
				{
					expect := true
					output := liquidityValue.Currency.Wrapped().Equal(tokenA)
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
				{
					expect := "917" // ceiling(1000 - (500 * (1 / 6)))
					output := liquidityValue.Quotient().String()
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
			}
		}
	}
}
