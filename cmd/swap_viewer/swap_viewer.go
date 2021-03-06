package main

import (
	"github.com/AidarBabanov/wallets-tracker/internal/logging"
	"github.com/AidarBabanov/wallets-tracker/internal/uniswap"
	"github.com/beego/beego/v2/core/logs"
	"os"
)

func main() {
	logging.Init()
	viewer := uniswap.New()
	err := viewer.Run()
	if err != nil {
		logs.Critical(err)
		os.Exit(-1)
	}
	txHash := os.Args[1] // submit as argument
	swap, err := viewer.ViewSwap(txHash)
	if err != nil {
		logs.Error(err)
	} else {
		logs.Info("%+v", swap)
	}
}
