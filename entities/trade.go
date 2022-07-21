package entities

import (
	"fmt"
	core "github.com/daoleno/uniswap-sdk-core/entities"
)

var (
	ErrInvalidSlippageTolerance = fmt.Errorf("invalid slippage tolerance")
	ErrInvalidCurrency          = fmt.Errorf("diff currency")
	ZeroFraction                = core.NewFraction(Zero, One)
)

// Trade Represents a trade executed against a list of pairs.
// Does not account for slippage, i.e. trades that front run this trade and move the price.
type Trade struct {
	/**
	 * The route of the trade, i.e. which pairs the trade goes through.
	 */
	Route *Route
	/**
	 * The type of the trade, either exact in or exact out.
	 */
	TradeType TradeType
	/**
	 * The input amount for the trade assuming no slippage.
	 */
	inputAmount *core.CurrencyAmount
	/**
	 * The output amount for the trade assuming no slippage.
	 */
	outputAmount *core.CurrencyAmount
	/**
	 * The price expressed in terms of output amount/input amount.
	 */
	ExecutionPrice *core.Price
	/**
	 * The mid price after the trade executes assuming no slippage.
	 */
	NextMidPrice *core.Price
	/**
	 * The percent difference between the mid price before the trade and the trade execution price.
	 */
	PriceImpact *core.Percent
}

func (t *Trade) InputAmount() *core.CurrencyAmount {
	return t.inputAmount
}

func (t *Trade) OutputAmount() *core.CurrencyAmount {
	return t.outputAmount
}

/**
 * Constructs an exact in trade with the given amount in and route
 * @param route route of the exact in trade
 * @param amountIn the amount being passed in
 */
func ExactIn(route *Route, amountIn *core.CurrencyAmount) (*Trade, error) {
	return NewTrade(route, amountIn, ExactInput)
}

/**
 * Constructs an exact out trade with the given amount out and route
 * @param route route of the exact out trade
 * @param amountOut the amount returned by the trade
 */
func ExactOut(route *Route, amountOut *core.CurrencyAmount) (*Trade, error) {
	return NewTrade(route, amountOut, ExactOutput)
}

// NewTrade creates a new trade
// nolint gocyclo
func NewTrade(route *Route, amount *core.CurrencyAmount, tradeType TradeType) (*Trade, error) {
	amounts := make([]*core.CurrencyAmount, len(route.Path))
	nextPairs := make([]*Pair, len(route.Pairs))
	var inputAmount, outputAmount *core.CurrencyAmount
	if tradeType == ExactInput {
		if !amount.Currency.Equal(route.Input) {
			return nil, ErrInvalidCurrency
		}
		amounts[0] = amount
		for i := 0; i < len(route.Path)-1; i++ {
			outputAmount, nextPair, err := route.Pairs[i].GetOutputAmount(amounts[i])
			if err != nil {
				return nil, err
			}
			amounts[i+1] = outputAmount
			nextPairs[i] = nextPair
		}
		inputAmount = core.FromFractionalAmount(route.Input, amount.Numerator, amount.Denominator)
		outputAmount = core.FromFractionalAmount(route.Output, amounts[len(amounts)-1].Numerator, amounts[len(amounts)-1].Denominator)
	} else {
		if !amount.Currency.Equal(route.Output) {
			return nil, ErrInvalidCurrency
		}
		amounts[len(amounts)-1] = amount
		for i := len(route.Path) - 1; i > 0; i-- {
			inputAmount, nextPair, err := route.Pairs[i-1].GetInputAmount(amounts[i])
			if err != nil {
				return nil, err
			}
			amounts[i-1] = inputAmount
			nextPairs[i-1] = nextPair
		}
		inputAmount = core.FromFractionalAmount(route.Input, amounts[0].Numerator, amounts[0].Denominator)
		outputAmount = core.FromFractionalAmount(route.Output, amount.Numerator, amount.Denominator)
	}

	nextRoute, err := NewRoute(nextPairs, route.Input, nil)
	if err != nil {
		return nil, err
	}
	nextMidPrice, err := nextRoute.MidPrice()
	if err != nil {
		return nil, err
	}
	price := core.NewPrice(inputAmount.Currency, outputAmount.Currency, inputAmount.Quotient(), outputAmount.Quotient())
	midPrice, err := route.MidPrice()
	if err != nil {
		return nil, err
	}
	return &Trade{
		Route:          route,
		TradeType:      tradeType,
		inputAmount:    inputAmount,
		outputAmount:   outputAmount,
		ExecutionPrice: price,
		NextMidPrice:   nextMidPrice,
		PriceImpact:    computePriceImpact(midPrice, inputAmount, outputAmount),
	}, nil
}

/**
 * Returns the percent difference between the mid price and the execution price, i.e. price impact.
 * @param midPrice mid price before the trade
 * @param inputAmount the input amount of the trade
 * @param outputAmount the output amount of the trade
 */
func computePriceImpact(midPrice *core.Price, inputAmount, outputAmount *core.CurrencyAmount) *core.Percent {
	exactQuote := midPrice.Fraction.Multiply(core.NewFraction(inputAmount.Quotient(), One))
	slippage := exactQuote.Subtract(core.NewFraction(outputAmount.Quotient(), One)).Divide(exactQuote)
	return &core.Percent{
		Fraction: slippage,
	}
}

/**
 * MinimumAmountOut - the minimum amount that must be received from this trade for the given slippage tolerance
 * @param slippageTolerance tolerance of unfavorable slippage from the execution price of this trade
 */
func (t *Trade) MinimumAmountOut(slippageTolerance *core.Percent) (*core.CurrencyAmount, error) {
	if slippageTolerance.LessThan(ZeroFraction) {
		return nil, ErrInvalidSlippageTolerance
	}

	if t.TradeType == ExactOutput {
		return t.outputAmount, nil
	}

	slippageAdjustedAmountOut := core.NewFraction(One, One).
		Add(slippageTolerance.Fraction).
		Invert().
		Multiply(core.NewFraction(t.outputAmount.Quotient(), One)).Quotient()
	return core.FromRawAmount(t.outputAmount.Currency, slippageAdjustedAmountOut), nil
}

/**
 * Get the maximum amount in that can be spent via this trade for the given slippage tolerance
 * @param slippageTolerance tolerance of unfavorable slippage from the execution price of this trade
 */
func (t *Trade) MaximumAmountIn(slippageTolerance *core.Percent) (*core.CurrencyAmount, error) {
	if slippageTolerance.LessThan(ZeroFraction) {
		return nil, ErrInvalidSlippageTolerance
	}

	if t.TradeType == ExactInput {
		return t.inputAmount, nil
	}

	slippageAdjustedAmountIn := core.NewFraction(One, One).
		Add(slippageTolerance.Fraction).
		Multiply(core.NewFraction(t.inputAmount.Quotient(), One)).Quotient()
	return core.FromRawAmount(t.inputAmount.Currency, slippageAdjustedAmountIn), nil
}
