package tracker

import (
	"github.com/AidarBabanov/wallets-tracker/internal/addrdb"
	"github.com/AidarBabanov/wallets-tracker/internal/explorer"
	"github.com/AidarBabanov/wallets-tracker/internal/notifier"
	"github.com/AidarBabanov/wallets-tracker/internal/uniswap"
	"github.com/asdine/storm/v3"
	"github.com/beego/beego/v2/core/logs"
	"gopkg.in/yaml.v2"
	"sort"
	"time"
)

const uniswapContractAddress = `0x7a250d5630b4cf539739df2c5dacb4c659f2488d`

type SwapTransaction struct {
	Transaction     string  `json:"transaction" yaml:"transaction"`
	Address         string  `json:"address" yaml:"address"`
	Gas             int64   `json:"gas" yaml:"gas"`
	GasUsed         int64   `json:"gas_used" yaml:"gas_used"`
	BlockNumber     int64   `json:"block_number" yaml:"block_number"`
	Timestamp       int64   `json:"timestamp" yaml:"timestamp"`
	FromTokenId     string  `json:"from_token_id" yaml:"from_token_id"`
	FromTokenSymbol string  `json:"from_token_symbol" yaml:"from_token_symbol"`
	FromTokenName   string  `json:"from_token_name" yaml:"from_token_name"`
	ToTokenId       string  `json:"to_token_id" yaml:"to_token_id"`
	ToTokenSymbol   string  `json:"to_token_symbol" yaml:"to_token_symbol"`
	ToTokenName     string  `json:"to_token_name" yaml:"to_token_name"`
	FromAmount      float64 `json:"from_amount" yaml:"from_amount"`
	ToAmount        float64 `json:"to_amount" yaml:"to_amount"`
}

type Tracker struct {
	AddressDatabase  addrdb.AddressDatabase
	Explorer         explorer.Explorer
	LastTransactions map[string]explorer.Transaction // key=address
	SwapViewer       *uniswap.SwapViewer
	Notifier         notifier.Notifier
	DB               *storm.DB
	AddressMutexMap  *MutexMap // saves from repeating
}

func New(
	addrDB addrdb.AddressDatabase,
	exp explorer.Explorer,
	swapViewer *uniswap.SwapViewer,
	notif notifier.Notifier,
) *Tracker {
	tracker := new(Tracker)
	tracker.LastTransactions = make(map[string]explorer.Transaction)
	tracker.AddressDatabase = addrDB
	tracker.Explorer = exp
	tracker.SwapViewer = swapViewer
	tracker.Notifier = notif
	tracker.AddressMutexMap = NewMutexMap()
	return tracker
}

func (t *Tracker) LoadLastTransactions() error {
	db, err := storm.Open("last_transactions.db")
	if err != nil {
		return err
	}
	var lastTransactions []explorer.Transaction
	err = db.All(&lastTransactions)
	if err != nil {
		return err
	}
	for _, trans := range lastTransactions {
		t.LastTransactions[trans.From] = trans
	}
	t.DB = db
	return nil
}

func (t *Tracker) RunTracker() {
	sem := make(chan struct{}, 10)
	for {
		address, err := t.AddressDatabase.Next()
		if err != nil {
			logs.Error(err)
		}
		sem <- struct{}{}
		go func(addr addrdb.Address) {
			var err error
			t.AddressMutexMap.Lock(addr.Address)
			defer func() {
				err = t.AddressMutexMap.Unlock(addr.Address)
				if err != nil {
					logs.Error(err)
				}
				<-sem
			}()
			last, existsLast := t.LastTransactions[addr.Address]
			if !existsLast {
				last.From = addr.Address
			}
			addrTrans := t.Explorer.NewTransactions(last)
			uniswapTrans := addrTrans[:0]
			isNewLast := false
			for _, tx := range addrTrans {
				if tx.To == uniswapContractAddress {
					uniswapTrans = append(uniswapTrans, tx)
				}
				if tx.Timestamp > last.Timestamp {
					last = tx
					isNewLast = true
				}
			}
			if isNewLast {
				t.LastTransactions[last.From] = last
				err = t.DB.Save(&last)
				if err != nil {
					logs.Error(err)
				}
			}

			// if previous last didn't exist, then probably we get very old transactions, so we don't want to show them
			// if previous last existed, then we want to make notification
			if existsLast && len(uniswapTrans) > 0 {
				sort.Sort(sort.Reverse(explorer.SortableTransactions(uniswapTrans)))
				for _, tx := range uniswapTrans {
					if tx.Timestamp > 3600 {
						break
					} // if transaction happened more than one hour ago, then skip it, it is too old
					swap, err := t.SwapViewer.ViewSwap(tx.Hash)
					if err != nil {
						logs.Error(err)
						continue
					}

					swapTransaction := SwapTransaction{
						Transaction:     tx.Hash,
						Address:         tx.From,
						Gas:             tx.Gas,
						GasUsed:         tx.GasUsed,
						BlockNumber:     tx.BlockNumber,
						Timestamp:       tx.Timestamp,
						FromTokenId:     swap.FromTokenId,
						FromTokenSymbol: swap.FromTokenSymbol,
						FromTokenName:   swap.FromTokenName,
						ToTokenId:       swap.ToTokenId,
						ToTokenSymbol:   swap.ToTokenSymbol,
						ToTokenName:     swap.ToTokenName,
						FromAmount:      swap.FromAmount,
						ToAmount:        swap.ToAmount,
					}

					go func(tx SwapTransaction) {
						yamlMsg, err := yaml.Marshal(tx)
						if err != nil {
							logs.Error(err)
							return
						}
						msg := string(yamlMsg)
						logs.Info(msg)
						t.Notifier.Notify(msg)
					}(swapTransaction)
				}
			}
		}(address)
		time.Sleep(time.Millisecond) // helps not to overload cpu
	}
}
