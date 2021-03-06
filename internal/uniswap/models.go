package uniswap

import jsoniter "github.com/json-iterator/go"

type TokenResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type PairResponse struct {
	Token0 TokenResponse `json:"token0"`
	Token1 TokenResponse `json:"token1"`
}

type SwapResponse struct {
	Amount0In  jsoniter.Number `json:"amount0In"`
	Amount1In  jsoniter.Number `json:"amount1In"`
	Amount0Out jsoniter.Number `json:"amount0Out"`
	Amount1Out jsoniter.Number `json:"amount1Out"`
	Pair       PairResponse    `json:"pair"`
}

type TransactionResponse struct {
	Id    string         `json:"id"`
	Swaps []SwapResponse `json:"swaps"`
}

type Response struct {
	Data struct {
		Transaction TransactionResponse `json:"transaction"`
	} `json:"data"`
}
