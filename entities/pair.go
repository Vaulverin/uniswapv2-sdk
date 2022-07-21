package entities

import (
	"errors"
	"github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

var (
	// ErrInvalidLiquidity invalid liquidity
	ErrInvalidLiquidity = errors.New("invalid liquidity")
	// ErrInvalidKLast invalid kLast
	ErrInvalidKLast            = errors.New("invalid kLast")
	ErrDiffToken               = errors.New("diff token")
	ErrInsufficientReserves    = errors.New("doesn't have insufficient reserves")
	ErrInsufficientInputAmount = errors.New("the input amount insufficient reserves")
)

// Tokens warps Token array
type Tokens [2]*entities.Token
type CurrencyAmounts [2]*entities.CurrencyAmount

func GetAddress(tokenA, tokenB *entities.Token, factory common.Address, initHash []byte) (common.Address, error) {
	ok, err := tokenA.SortsBefore(tokenB)
	if err != nil {
		return [20]byte{}, err
	}
	tokens := Tokens{tokenA, tokenB}
	if !ok {
		tokens[0], tokens[1] = tokens[1], tokens[0]
	}
	var salt [32]byte
	copy(salt[:], crypto.Keccak256(append(tokens[0].Address.Bytes(), tokens[1].Address.Bytes()...)))
	return crypto.CreateAddress2(factory, salt, initHash), nil
}

// NewCurrencyAmounts creates a CurrencyAmounts
func NewCurrencyAmounts(amountA, amountB *entities.CurrencyAmount) (CurrencyAmounts, error) {
	ok, err := amountA.Currency.Wrapped().SortsBefore(amountB.Currency.Wrapped())
	if err != nil {
		return CurrencyAmounts{}, err
	}
	if ok {
		return CurrencyAmounts{amountA, amountB}, nil
	}
	return CurrencyAmounts{amountB, amountA}, nil
}

// Pair warps uniswap pair
type Pair struct {
	LiquidityToken *entities.Token
	TokenAmounts   CurrencyAmounts // sorted tokens
	Address        common.Address  // Pair address
	Options        *PairOptions
}

// PairOptions for generating pair address
type PairOptions struct {
	Factory      common.Address  // Uniswap factory address
	InitCodeHash []byte          // Chain init code
	Address      *common.Address // Pair address if already known. Leave empty if not
}

// NewPair creates Pair
func NewPair(amountA, amountB *entities.CurrencyAmount, options *PairOptions) (*Pair, error) {
	amounts, err := NewCurrencyAmounts(amountA, amountB)
	if err != nil {
		return nil, err
	}
	opts := options
	if opts == nil {
		opts = &PairOptions{
			Factory:      FactoryAddress,
			InitCodeHash: InitCodeHash,
		}
	}
	var pairAddress common.Address
	if opts.Address != nil {
		pairAddress = *opts.Address
	} else {
		pairAddress, err = GetAddress(amountA.Currency.Wrapped(), amountB.Currency.Wrapped(), opts.Factory, opts.InitCodeHash)
		if err != nil {
			return nil, err
		}
	}
	liquidityToken := entities.NewToken(amounts[0].Currency.Wrapped().ChainId(), pairAddress,
		18, "UNI-V2", "Uniswap V2")
	return &Pair{
		TokenAmounts:   amounts,
		LiquidityToken: liquidityToken,
		Address:        pairAddress,
		Options:        opts,
	}, nil
}

// GetAddress returns a contract's address for a pair
func (p *Pair) GetAddress() common.Address {
	return p.Address
}

// InvolvesToken Returns true if the token is either token0 or token1
// @param token to check
func (p *Pair) InvolvesToken(token *entities.Token) bool {
	return token.Equal(p.TokenAmounts[0].Currency.Wrapped()) || token.Equal(p.TokenAmounts[1].Currency.Wrapped())
}

// Token0Price Returns the current mid price of the pair in terms of token0, i.e. the ratio of reserve1 to reserve0
func (p *Pair) Token0Price() *entities.Price {
	return entities.NewPrice(p.Token0(), p.Token1(), p.TokenAmounts[0].Numerator, p.TokenAmounts[1].Numerator)
}

// Token1Price Returns the current mid price of the pair in terms of token1, i.e. the ratio of reserve0 to reserve1
func (p *Pair) Token1Price() *entities.Price {
	return entities.NewPrice(p.Token1(), p.Token0(), p.TokenAmounts[1].Numerator, p.TokenAmounts[0].Numerator)
}

// PriceOf Returns the price of the given token in terms of the other token in the pair.
// @param token token to return price of
func (p *Pair) PriceOf(token *entities.Token) (*entities.Price, error) {
	if !p.InvolvesToken(token) {
		return nil, ErrDiffToken
	}

	if token.Equal(p.Token0()) {
		return p.Token0Price(), nil
	}
	return p.Token1Price(), nil
}

// ChainID Returns the chain ID of the tokens in the pair.
func (p *Pair) ChainID() uint {
	return p.Token0().ChainId()
}

// Token0 returns the first token in the pair
func (p *Pair) Token0() *entities.Token {
	return p.TokenAmounts[0].Currency.Wrapped()
}

// Token1 returns the last token in the pair
func (p *Pair) Token1() *entities.Token {
	return p.TokenAmounts[1].Currency.Wrapped()
}

// Reserve0 returns the first CurrencyAmount in the pair
func (p *Pair) Reserve0() *entities.CurrencyAmount {
	return p.TokenAmounts[0]
}

// Reserve1 returns the last CurrencyAmount in the pair
func (p *Pair) Reserve1() *entities.CurrencyAmount {
	return p.TokenAmounts[1]
}

// ReserveOf returns the CurrencyAmount that equals to the token
func (p *Pair) ReserveOf(token *entities.Token) (*entities.CurrencyAmount, error) {
	if !p.InvolvesToken(token) {
		return nil, ErrDiffToken
	}

	if token.Equal(p.Token0()) {
		return p.Reserve0(), nil
	}
	return p.Reserve1(), nil
}

// GetOutputAmount returns OutputAmount and a Pair for the InputAmout
func (p *Pair) GetOutputAmount(inputAmount *entities.CurrencyAmount) (*entities.CurrencyAmount, *Pair, error) {
	if !p.InvolvesToken(inputAmount.Currency.Wrapped()) {
		return nil, nil, ErrDiffToken
	}

	if p.Reserve0().Quotient().Cmp(Zero) == 0 ||
		p.Reserve1().Quotient().Cmp(Zero) == 0 {
		return nil, nil, ErrInsufficientReserves
	}

	inputReserve, err := p.ReserveOf(inputAmount.Currency.Wrapped())
	if err != nil {
		return nil, nil, err
	}
	token := p.Token0()
	if inputAmount.Currency.Wrapped().Equal(p.Token0()) {
		token = p.Token1()
	}
	outputReserve, err := p.ReserveOf(token)
	if err != nil {
		return nil, nil, err
	}

	inputAmountWithFee := big.NewInt(0).Mul(inputAmount.Quotient(), B997)
	numerator := big.NewInt(0).Mul(inputAmountWithFee, outputReserve.Quotient())
	denominator := big.NewInt(0).Add(big.NewInt(0).Mul(inputReserve.Quotient(), B1000), inputAmountWithFee)
	outputAmount := entities.FromRawAmount(token, big.NewInt(0).Div(numerator, denominator))
	if outputAmount.Quotient().Cmp(Zero) == 0 {
		return nil, nil, ErrInsufficientInputAmount
	}

	tokenAmountA := inputAmount.Add(inputReserve)
	tokenAmountB := outputReserve.Subtract(outputAmount)
	pair, err := NewPair(tokenAmountA, tokenAmountB, p.Options)
	if err != nil {
		return nil, nil, err
	}

	return outputAmount, pair, nil
}

// GetInputAmount returns InputAmout and a Pair for the OutputAmount
func (p *Pair) GetInputAmount(outputAmount *entities.CurrencyAmount) (*entities.CurrencyAmount, *Pair, error) {
	if !p.InvolvesToken(outputAmount.Currency.Wrapped()) {
		return nil, nil, ErrDiffToken
	}

	outputReserve, err := p.ReserveOf(outputAmount.Currency.Wrapped())
	if err != nil {
		return nil, nil, err
	}
	if p.Reserve0().Quotient().Cmp(Zero) == 0 ||
		p.Reserve1().Quotient().Cmp(Zero) == 0 ||
		outputAmount.Quotient().Cmp(outputReserve.Quotient()) >= 0 {
		return nil, nil, ErrInsufficientReserves
	}

	token := p.Token0()
	if outputAmount.Currency.Wrapped().Equal(p.Token0()) {
		token = p.Token1()
	}
	inputReserve, err := p.ReserveOf(token)
	if err != nil {
		return nil, nil, err
	}

	numerator := big.NewInt(0).Mul(inputReserve.Quotient(), outputAmount.Quotient())
	numerator.Mul(numerator, B1000)
	denominator := big.NewInt(0).Sub(outputReserve.Quotient(), outputAmount.Quotient())
	denominator.Mul(denominator, B997)
	amount := big.NewInt(0).Div(numerator, denominator)
	amount.Add(amount, One)
	inputAmount := entities.FromRawAmount(token, amount)
	if err != nil {
		return nil, nil, err
	}

	tokenAmountA := inputAmount.Add(inputReserve)
	tokenAmountB := outputReserve.Subtract(outputAmount)
	pair, err := NewPair(tokenAmountA, tokenAmountB, p.Options)
	if err != nil {
		return nil, nil, err
	}

	return inputAmount, pair, nil
}

// GetLiquidityMinted returns liquidity minted CurrencyAmount
func (p *Pair) GetLiquidityMinted(totalSupply, tokenAmountA, tokenAmountB *entities.CurrencyAmount) (*entities.CurrencyAmount, error) {
	if !p.LiquidityToken.Equal(totalSupply.Currency.Wrapped()) {
		return nil, ErrDiffToken
	}

	tokenAmounts, err := NewCurrencyAmounts(tokenAmountA, tokenAmountB)
	if err != nil {
		return nil, err
	}
	if !(tokenAmounts[0].Currency.Wrapped().Equal(p.Token0()) && tokenAmounts[1].Currency.Wrapped().Equal(p.Token1())) {
		return nil, ErrDiffToken
	}

	var liquidity *big.Int
	if totalSupply.Quotient().Cmp(Zero) == 0 {
		liquidity = big.NewInt(0).Mul(tokenAmounts[0].Quotient(), tokenAmounts[1].Quotient())
		liquidity.Sqrt(liquidity)
		liquidity.Sub(liquidity, MinimumLiquidity)
	} else {
		amount0 := big.NewInt(0).Mul(tokenAmounts[0].Quotient(), totalSupply.Quotient())
		amount0.Div(amount0, p.Reserve0().Quotient())
		amount1 := big.NewInt(0).Mul(tokenAmounts[1].Quotient(), totalSupply.Quotient())
		amount1.Div(amount1, p.Reserve1().Quotient())
		liquidity = amount0
		if liquidity.Cmp(amount1) > 0 {
			liquidity = amount1
		}
	}

	if liquidity.Cmp(Zero) <= 0 {
		return nil, ErrInsufficientInputAmount
	}

	return entities.FromRawAmount(p.LiquidityToken, liquidity), nil
}

// GetLiquidityValue returns liquidity value CurrencyAmount
func (p *Pair) GetLiquidityValue(token *entities.Token, totalSupply, liquidity *entities.CurrencyAmount, feeOn bool, kLast *big.Int) (*entities.CurrencyAmount, error) {
	if !p.InvolvesToken(token) || !p.LiquidityToken.Equal(totalSupply.Currency.Wrapped()) || !p.LiquidityToken.Equal(liquidity.Currency.Wrapped()) {
		return nil, ErrDiffToken
	}
	if liquidity.Quotient().Cmp(totalSupply.Quotient()) > 0 {
		return nil, ErrInvalidLiquidity
	}

	totalSupplyAdjusted, err := p.adjustTotalSupply(totalSupply, feeOn, kLast)
	if err != nil {
		return nil, err
	}

	tokenAmount, err := p.ReserveOf(token)
	if err != nil {
		return nil, err
	}

	amount := big.NewInt(0).Mul(liquidity.Quotient(), tokenAmount.Quotient())
	amount.Div(amount, totalSupplyAdjusted.Quotient())
	return entities.FromRawAmount(token, amount), nil
}

func (p *Pair) adjustTotalSupply(totalSupply *entities.CurrencyAmount, feeOn bool, kLast *big.Int) (*entities.CurrencyAmount, error) {
	if !feeOn {
		return totalSupply, nil
	}

	if kLast == nil {
		return nil, ErrInvalidKLast
	}
	if kLast.Cmp(Zero) == 0 {
		return totalSupply, nil
	}

	rootK := big.NewInt(0).Mul(p.Reserve0().Quotient(), p.Reserve1().Quotient())
	rootK.Sqrt(rootK)
	rootKLast := big.NewInt(0).Sqrt(kLast)
	if rootK.Cmp(rootKLast) <= 0 {
		return totalSupply, nil
	}

	numerator := big.NewInt(0).Sub(rootK, rootKLast)
	numerator.Mul(numerator, totalSupply.Quotient())
	denominator := big.NewInt(0).Mul(rootK, Five)
	denominator.Add(denominator, rootKLast)
	tokenAmount := entities.FromFractionalAmount(p.LiquidityToken, numerator, denominator)
	return totalSupply.Add(tokenAmount), nil
}
