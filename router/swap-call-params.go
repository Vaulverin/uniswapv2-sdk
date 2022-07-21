package router

import (
	"errors"
	core "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vaulverin/uniswapv2-sdk/entities"
	"math/big"
	"strings"
	"time"
)

const V2Router02ABI = "[ { \"inputs\": [ { \"internalType\": \"address\", \"name\": \"_factory\", \"type\": \"address\" }, { \"internalType\": \"address\", \"name\": \"_WETH\", \"type\": \"address\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"constructor\" }, { \"inputs\": [], \"name\": \"WETH\", \"outputs\": [ { \"internalType\": \"address\", \"name\": \"\", \"type\": \"address\" } ], \"stateMutability\": \"view\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"address\", \"name\": \"tokenA\", \"type\": \"address\" }, { \"internalType\": \"address\", \"name\": \"tokenB\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"amountADesired\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountBDesired\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountAMin\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountBMin\", \"type\": \"uint256\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"addLiquidity\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountA\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountB\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"liquidity\", \"type\": \"uint256\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"address\", \"name\": \"token\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"amountTokenDesired\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountTokenMin\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountETHMin\", \"type\": \"uint256\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"addLiquidityETH\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountToken\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountETH\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"liquidity\", \"type\": \"uint256\" } ], \"stateMutability\": \"payable\", \"type\": \"function\" }, { \"inputs\": [], \"name\": \"factory\", \"outputs\": [ { \"internalType\": \"address\", \"name\": \"\", \"type\": \"address\" } ], \"stateMutability\": \"view\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountOut\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"reserveIn\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"reserveOut\", \"type\": \"uint256\" } ], \"name\": \"getAmountIn\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountIn\", \"type\": \"uint256\" } ], \"stateMutability\": \"pure\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountIn\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"reserveIn\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"reserveOut\", \"type\": \"uint256\" } ], \"name\": \"getAmountOut\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountOut\", \"type\": \"uint256\" } ], \"stateMutability\": \"pure\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountOut\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" } ], \"name\": \"getAmountsIn\", \"outputs\": [ { \"internalType\": \"uint256[]\", \"name\": \"amounts\", \"type\": \"uint256[]\" } ], \"stateMutability\": \"view\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountIn\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" } ], \"name\": \"getAmountsOut\", \"outputs\": [ { \"internalType\": \"uint256[]\", \"name\": \"amounts\", \"type\": \"uint256[]\" } ], \"stateMutability\": \"view\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountA\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"reserveA\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"reserveB\", \"type\": \"uint256\" } ], \"name\": \"quote\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountB\", \"type\": \"uint256\" } ], \"stateMutability\": \"pure\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"address\", \"name\": \"tokenA\", \"type\": \"address\" }, { \"internalType\": \"address\", \"name\": \"tokenB\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"liquidity\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountAMin\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountBMin\", \"type\": \"uint256\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"removeLiquidity\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountA\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountB\", \"type\": \"uint256\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"address\", \"name\": \"token\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"liquidity\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountTokenMin\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountETHMin\", \"type\": \"uint256\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"removeLiquidityETH\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountToken\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountETH\", \"type\": \"uint256\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"address\", \"name\": \"token\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"liquidity\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountTokenMin\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountETHMin\", \"type\": \"uint256\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"removeLiquidityETHSupportingFeeOnTransferTokens\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountETH\", \"type\": \"uint256\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"address\", \"name\": \"token\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"liquidity\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountTokenMin\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountETHMin\", \"type\": \"uint256\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" }, { \"internalType\": \"bool\", \"name\": \"approveMax\", \"type\": \"bool\" }, { \"internalType\": \"uint8\", \"name\": \"v\", \"type\": \"uint8\" }, { \"internalType\": \"bytes32\", \"name\": \"r\", \"type\": \"bytes32\" }, { \"internalType\": \"bytes32\", \"name\": \"s\", \"type\": \"bytes32\" } ], \"name\": \"removeLiquidityETHWithPermit\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountToken\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountETH\", \"type\": \"uint256\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"address\", \"name\": \"token\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"liquidity\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountTokenMin\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountETHMin\", \"type\": \"uint256\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" }, { \"internalType\": \"bool\", \"name\": \"approveMax\", \"type\": \"bool\" }, { \"internalType\": \"uint8\", \"name\": \"v\", \"type\": \"uint8\" }, { \"internalType\": \"bytes32\", \"name\": \"r\", \"type\": \"bytes32\" }, { \"internalType\": \"bytes32\", \"name\": \"s\", \"type\": \"bytes32\" } ], \"name\": \"removeLiquidityETHWithPermitSupportingFeeOnTransferTokens\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountETH\", \"type\": \"uint256\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"address\", \"name\": \"tokenA\", \"type\": \"address\" }, { \"internalType\": \"address\", \"name\": \"tokenB\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"liquidity\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountAMin\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountBMin\", \"type\": \"uint256\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" }, { \"internalType\": \"bool\", \"name\": \"approveMax\", \"type\": \"bool\" }, { \"internalType\": \"uint8\", \"name\": \"v\", \"type\": \"uint8\" }, { \"internalType\": \"bytes32\", \"name\": \"r\", \"type\": \"bytes32\" }, { \"internalType\": \"bytes32\", \"name\": \"s\", \"type\": \"bytes32\" } ], \"name\": \"removeLiquidityWithPermit\", \"outputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountA\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountB\", \"type\": \"uint256\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountOut\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"swapETHForExactTokens\", \"outputs\": [ { \"internalType\": \"uint256[]\", \"name\": \"amounts\", \"type\": \"uint256[]\" } ], \"stateMutability\": \"payable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountOutMin\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"swapExactETHForTokens\", \"outputs\": [ { \"internalType\": \"uint256[]\", \"name\": \"amounts\", \"type\": \"uint256[]\" } ], \"stateMutability\": \"payable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountOutMin\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"swapExactETHForTokensSupportingFeeOnTransferTokens\", \"outputs\": [], \"stateMutability\": \"payable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountIn\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountOutMin\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"swapExactTokensForETH\", \"outputs\": [ { \"internalType\": \"uint256[]\", \"name\": \"amounts\", \"type\": \"uint256[]\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountIn\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountOutMin\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"swapExactTokensForETHSupportingFeeOnTransferTokens\", \"outputs\": [], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountIn\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountOutMin\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"swapExactTokensForTokens\", \"outputs\": [ { \"internalType\": \"uint256[]\", \"name\": \"amounts\", \"type\": \"uint256[]\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountIn\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountOutMin\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"swapExactTokensForTokensSupportingFeeOnTransferTokens\", \"outputs\": [], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountOut\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountInMax\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"swapTokensForExactETH\", \"outputs\": [ { \"internalType\": \"uint256[]\", \"name\": \"amounts\", \"type\": \"uint256[]\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"inputs\": [ { \"internalType\": \"uint256\", \"name\": \"amountOut\", \"type\": \"uint256\" }, { \"internalType\": \"uint256\", \"name\": \"amountInMax\", \"type\": \"uint256\" }, { \"internalType\": \"address[]\", \"name\": \"path\", \"type\": \"address[]\" }, { \"internalType\": \"address\", \"name\": \"to\", \"type\": \"address\" }, { \"internalType\": \"uint256\", \"name\": \"deadline\", \"type\": \"uint256\" } ], \"name\": \"swapTokensForExactTokens\", \"outputs\": [ { \"internalType\": \"uint256[]\", \"name\": \"amounts\", \"type\": \"uint256[]\" } ], \"stateMutability\": \"nonpayable\", \"type\": \"function\" }, { \"stateMutability\": \"payable\", \"type\": \"receive\" } ]"

var (
	ErrEtherInOut  = errors.New("the router does not support both ether in and out")
	ErrExactOutFot = errors.New("EXACT_OUT_FOT")
)

// TradeOptions for producing the arguments to send call to the router.
type TradeOptions struct {
	AllowedSlippage *core.Percent  // How much the execution price is allowed to move unfavorably from the trade execution price.
	Recipient       common.Address // The account that should receive the output.
	Deadline        *big.Int       // When the transaction expires, in epoch seconds.
	FeeOnTransfer   bool           // Whether any of the tokens in the path are fee on transfer tokens, which should be handled with special methods
}

// SwapParameters to use in the call to the Uniswap V2 Router to execute a trade.
type SwapParameters struct {
	MethodName string        // The method to call on the Uniswap V2 Router.
	Args       []interface{} // The arguments to pass to the method.
	Value      *big.Int      // The amount of wei to send.
}

// toHex converts a big int to a hex string
func toHex(i *big.Int) string {
	if i == nil {
		return "0x00"
	}

	hex := i.String()
	if len(hex)%2 != 0 {
		hex = "0" + hex
	}
	return "0x" + hex
}

// SwapCallParameters produces the on-chain method name to call and the hex encoded parameters to pass as arguments for a given trade.
func SwapCallParameters(trade *entities.Trade, options TradeOptions) (*SwapParameters, error) {
	etherIn := trade.InputAmount().Currency.IsNative()
	etherOut := trade.OutputAmount().Currency.IsNative()
	if etherIn && etherOut {
		return nil, ErrEtherInOut
	}
	to := options.Recipient.Hex()
	slippage := options.AllowedSlippage
	if slippage == nil {
		slippage = core.NewPercent(big.NewInt(0), big.NewInt(1))
	}
	maxAmountIn, err := trade.MaximumAmountIn(slippage)
	if err != nil {
		return nil, err
	}
	amountIn := maxAmountIn.Quotient()
	minAmountOut, err := trade.MinimumAmountOut(slippage)
	if err != nil {
		return nil, err
	}
	amountOut := minAmountOut.Quotient()
	var path []string
	for _, token := range trade.Route.Path {
		path = append(path, token.Address.Hex())
	}
	deadline := options.Deadline
	if options.Deadline == nil {
		deadline = big.NewInt(time.Now().Add(5 * time.Minute).Unix())
	}

	var (
		methodName string
		args       []interface{}
		value      *big.Int
	)
	switch trade.TradeType {
	case entities.ExactInput:
		if etherIn {
			methodName = "swapExactETHForTokens"
			if options.FeeOnTransfer {
				methodName = "swapExactETHForTokensSupportingFeeOnTransferTokens"
			}
			// (uint amountOut, address[] calldata path, address to, uint deadline)
			args = []interface{}{amountOut, path, to, deadline}
			value = amountIn
			break
		} else if etherOut {
			methodName = "swapExactTokensForETH"
			if options.FeeOnTransfer {
				methodName = "swapExactTokensForETHSupportingFeeOnTransferTokens"
			}
			// (uint amountOut, uint amountInMax, address[] calldata path, address to, uint deadline)
			args = []interface{}{amountIn, amountOut, path, to, deadline}
			value = big.NewInt(0)
			break
		}
		methodName = "swapExactTokensForTokens"
		if options.FeeOnTransfer {
			methodName = "swapExactTokensForTokensSupportingFeeOnTransferTokens"
		}
		// (uint amountIn, uint amountOutMin, address[] calldata path, address to, uint deadline)
		args = []interface{}{amountIn, amountOut, path, to, deadline}
		value = big.NewInt(0)
	case entities.ExactOutput:
		if options.FeeOnTransfer {
			return nil, ErrExactOutFot
		}
		if etherIn {
			methodName = "swapETHForExactTokens"
			// (uint amountOut, address[] calldata path, address to, uint deadline)
			args = []interface{}{amountOut, path, to, deadline}
			value = amountIn
			break
		} else if etherOut {
			methodName = "swapTokensForExactETH"
			// (uint amountOut, uint amountInMax, address[] calldata path, address to, uint deadline)
			args = []interface{}{amountOut, amountIn, path, to, deadline}
			value = big.NewInt(0)
			break
		}
		methodName = "swapTokensForExactTokens"
		// (uint amountIn, uint amountOutMin, address[] calldata path, address to, uint deadline)
		args = []interface{}{amountOut, amountIn, path, to, deadline}
		value = big.NewInt(0)
	}
	return &SwapParameters{
		MethodName: methodName,
		Args:       args,
		Value:      value,
	}, nil
}

// SwapCallParametersPacked packs swap parameters.
// 	Returns value and data to use in transaction body.
func SwapCallParametersPacked(trade *entities.Trade, options TradeOptions) (*big.Int, []byte, error) {
	params, err := SwapCallParameters(trade, options)
	if err != nil {
		return nil, nil, err
	}
	routerABI, err := abi.JSON(strings.NewReader(V2Router02ABI))
	if err != nil {
		return nil, nil, err
	}
	data, err := routerABI.Pack(params.MethodName, params.Args)
	if err != nil {
		return nil, nil, err
	}
	return params.Value, data, nil
}
