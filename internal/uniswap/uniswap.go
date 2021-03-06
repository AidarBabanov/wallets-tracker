package uniswap

import (
	"fmt"
	"github.com/AidarBabanov/wallets-tracker/internal/httpclient"
	"time"
)

const (
	baseURL = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"
	// submit transaction hash instead of %s
	getTransaction = "{\"query\":\"{\\n  transaction(id:\\\"%s\\\"){\\n    id\\n    swaps{\\n      pair{\\n        token0{\\n          id\\n          symbol\\n          name\\n        }\\n        token1{\\n          id\\n          symbol\\n          name\\n        }\\n      }\\n      \\n      amount0In\\n      amount1In\\n      \\n      amount0Out\\n      amount1Out\\n    }\\n  }\\n}\\n\",\"variables\":{}}"
)

type Swap struct {
	FromTokenId     string  `json:"from_token_id"`
	FromTokenSymbol string  `json:"from_token_symbol"`
	FromTokenName   string  `json:"from_token_name"`
	ToTokenId       string  `json:"to_token_id"`
	ToTokenSymbol   string  `json:"to_token_symbol"`
	ToTokenName     string  `json:"to_token_name"`
	FromAmount      float64 `json:"from_amount"`
	ToAmount        float64 `json:"to_amount"`
}

type SwapViewer struct {
	client *httpclient.GraphqlClient
}

func New() *SwapViewer {
	viewer := new(SwapViewer)
	viewer.client = httpclient.NewGraphQLClient(baseURL, 10*time.Second, 275*time.Millisecond)
	return viewer
}

func (v *SwapViewer) Run() error {
	return v.client.StartDelayer()
}

func (v *SwapViewer) Close() {
	v.client.Close()
}

func (v *SwapViewer) ViewSwap(txHash string) (Swap, error) {
	var response Response
	err := v.client.DoGraphqlRequest(fmt.Sprintf(getTransaction, txHash), &response)
	if err != nil {
		return Swap{}, err
	}

	if len(response.Data.Transaction.Swaps) == 0 {
		return Swap{}, fmt.Errorf("no swaps in the transaction")
	}
	from := response.Data.Transaction.Swaps[0]
	to := response.Data.Transaction.Swaps[len(response.Data.Transaction.Swaps)-1]
	fromToken, fromAmount, err := getFromData(from)
	if err != nil {
		return Swap{}, err
	}
	toToken, toAmount, err := getToData(to)
	if err != nil {
		return Swap{}, err
	}
	swap := Swap{
		FromTokenId:     fromToken.Id,
		FromTokenSymbol: fromToken.Symbol,
		FromTokenName:   fromToken.Name,
		ToTokenId:       toToken.Id,
		ToTokenSymbol:   toToken.Symbol,
		ToTokenName:     toToken.Name,
		FromAmount:      fromAmount,
		ToAmount:        toAmount,
	}
	return swap, nil
}

func getFromData(from SwapResponse) (TokenResponse, float64, error) {
	fromAmount0, err := from.Amount0In.Float64()
	if err != nil {
		return TokenResponse{}, 0, err
	}
	fromAmount1, err := from.Amount1In.Float64()
	if err != nil {
		return TokenResponse{}, 0, err
	}
	fromAmount := 0.0
	fromToken := TokenResponse{}
	if fromAmount0 > 0 && fromAmount1 == 0 {
		fromAmount = fromAmount0
		fromToken = from.Pair.Token0
	} else if fromAmount0 == 0 && fromAmount1 > 0 {
		fromAmount = fromAmount1
		fromToken = from.Pair.Token1
	} else {
		return TokenResponse{}, 0, fmt.Errorf("wrong amounts")
	}
	return fromToken, fromAmount, nil
}

func getToData(to SwapResponse) (TokenResponse, float64, error) {
	toAmount0, err := to.Amount0Out.Float64()
	if err != nil {
		return TokenResponse{}, 0, err
	}
	toAmount1, err := to.Amount1Out.Float64()
	if err != nil {
		return TokenResponse{}, 0, err
	}
	toAmount := 0.0
	toToken := TokenResponse{}
	if toAmount0 > 0 && toAmount1 == 0 {
		toAmount = toAmount0
		toToken = to.Pair.Token0
	} else if toAmount0 == 0 && toAmount1 > 0 {
		toAmount = toAmount1
		toToken = to.Pair.Token1
	} else {
		return TokenResponse{}, 0, fmt.Errorf("wrong amounts")
	}
	return toToken, toAmount, nil
}
