package entities

import (
	"fmt"
	core "github.com/daoleno/uniswap-sdk-core/entities"
)

var (
	ErrInvalidPairs         = fmt.Errorf("invalid pairs")
	ErrInvalidPairsChainIDs = fmt.Errorf("invalid pairs chainIDs")
	ErrInvalidInput         = fmt.Errorf("invalid token input")
	ErrInvalidOutput        = fmt.Errorf("invalid token output")
	ErrInvalidPath          = fmt.Errorf("invalid pairs for path")
)

type Route struct {
	Pairs    []*Pair
	Path     []*core.Token
	Input    core.Currency
	Output   core.Currency
	midPrice *core.Price
}

func NewRoute(pairs []*Pair, input, output core.Currency) (*Route, error) {
	if len(pairs) == 0 {
		return nil, ErrInvalidPairs
	}

	for i := range pairs {
		if pairs[i].ChainID() != pairs[0].ChainID() {
			return nil, ErrInvalidPairsChainIDs
		}
	}

	if !pairs[0].InvolvesToken(input.Wrapped()) {
		return nil, ErrInvalidInput
	}
	if !(output == nil || pairs[len(pairs)-1].InvolvesToken(output.Wrapped())) {
		return nil, ErrInvalidOutput
	}

	path := make([]*core.Token, len(pairs)+1)
	path[0] = input.Wrapped()
	for i := range pairs {
		currentInput := path[i]
		if !(currentInput.Equal(pairs[i].Token0()) || currentInput.Equal(pairs[i].Token1())) {
			return nil, ErrInvalidPath
		}
		currentOutput := pairs[i].Token0()
		if currentInput.Equal(pairs[i].Token0()) {
			currentOutput = pairs[i].Token1()
		}
		path[i+1] = currentOutput
	}

	if output == nil {
		output = path[len(pairs)]
	}

	route := &Route{
		Pairs:  pairs,
		Path:   path,
		Input:  input,
		Output: output,
	}
	return route, nil
}

func (r *Route) MidPrice() (*core.Price, error) {
	if r.midPrice != nil {
		return r.midPrice, nil
	}
	length := len(r.Pairs)
	// NOTE: check route Pairs len?
	prices := make([]*core.Price, length)
	for i := range r.Pairs {
		if r.Path[i].Equal(r.Pairs[i].Token0()) {
			prices[i] = core.NewPrice(r.Pairs[i].Reserve0().Currency, r.Pairs[i].Reserve1().Currency,
				r.Pairs[i].Reserve0().Quotient(), r.Pairs[i].Reserve1().Quotient())
		} else {
			prices[i] = core.NewPrice(r.Pairs[i].Reserve1().Currency, r.Pairs[i].Reserve0().Currency,
				r.Pairs[i].Reserve1().Quotient(), r.Pairs[i].Reserve0().Quotient())
		}
	}
	price := prices[0]
	var err error
	for i := 1; i < length; i++ {
		price, err = price.Multiply(prices[i])
		if err != nil {
			return nil, err
		}
	}
	r.midPrice = price
	return r.midPrice, nil
}

func (r *Route) ChainID() uint {
	return r.Pairs[0].ChainID()
}
