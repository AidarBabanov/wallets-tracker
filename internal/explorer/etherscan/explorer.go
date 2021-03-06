package etherscan

import (
	"github.com/AidarBabanov/wallets-tracker/internal/explorer"
	"github.com/AidarBabanov/wallets-tracker/internal/httpclient"
	"github.com/beego/beego/v2/core/logs"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	getAccountTransactionsURL = "https://api.etherscan.io/api"
)

type transactionResponse struct {
	BlockNumber       jsoniter.Number `json:"blockNumber"`
	TimeStamp         jsoniter.Number `json:"timeStamp"`
	Hash              string          `json:"hash"`
	Nonce             string          `json:"nonce"`
	BlockHash         string          `json:"blockHash"`
	TransactionIndex  string          `json:"transactionIndex"`
	From              string          `json:"from"`
	To                string          `json:"to"`
	Value             string          `json:"value"`
	Gas               jsoniter.Number `json:"gas"`
	GasPrice          jsoniter.Number `json:"gasPrice"`
	IsError           string          `json:"isError"`
	TxreceiptStatus   string          `json:"txreceipt_status"`
	Input             string          `json:"input"`
	ContractAddress   string          `json:"contractAddress"`
	CumulativeGasUsed jsoniter.Number `json:"cumulativeGasUsed"`
	GasUsed           jsoniter.Number `json:"gasUsed"`
	Confirmations     string          `json:"confirmations"`
}

type transactionsResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Result  []transactionResponse `json:"result"`
}

type Explorer struct {
	httpClient *httpclient.Client
	apiKey     string
}

func New(apiKey string) *Explorer {
	exp := new(Explorer)
	exp.apiKey = apiKey
	httpClient := httpclient.New(10*time.Second, 200*time.Millisecond)
	exp.httpClient = httpClient
	return exp
}

func (e *Explorer) Run() error {
	return e.httpClient.StartDelayer()
}

func (e *Explorer) Close() {
	e.httpClient.Close()
}

func (e *Explorer) NewTransactions(last explorer.Transaction) []explorer.Transaction {
	var response transactionsResponse
	var err error
	params := map[string]string{
		"module":     "account",
		"action":     "txlist",
		"startblock": strconv.Itoa(int(last.BlockNumber) + 1),
		"sort":       "desc",
		"apiKey":     e.apiKey,
		"address":    last.From,
	}
	err = e.httpClient.DoRequest(http.MethodGet, getAccountTransactionsURL, params, nil, nil, &response)
	for count := 0; err != nil && count < 10; count++ {
		err = e.httpClient.DoRequest(http.MethodGet, getAccountTransactionsURL, params, nil, nil, &response)
	}
	var newTransactions []explorer.Transaction
	for _, responseTx := range response.Result {
		timestamp, err := responseTx.TimeStamp.Int64()
		if err != nil {
			logs.Error(err)
		}
		blockNumber, err := responseTx.BlockNumber.Int64()
		if err != nil {
			logs.Error(err)
		}
		gas, err := responseTx.Gas.Int64()
		if err != nil {
			logs.Error(err)
		}
		gasUsed, err := responseTx.GasUsed.Int64()
		if err != nil {
			logs.Error(err)
		}
		gasPrice, err := responseTx.GasPrice.Int64()
		if err != nil {
			logs.Error(err)
		}
		newTransactions = append(newTransactions, explorer.Transaction{
			From:        last.From,
			Hash:        responseTx.Hash,
			Timestamp:   timestamp,
			BlockNumber: blockNumber,
			To:          responseTx.To,
			Gas:         gas,
			GasUsed:     gasUsed,
			GasPrice:    gasPrice,
		})
	}

	sort.Sort(explorer.SortableTransactions(newTransactions))
	return newTransactions
}
