package main

import (
	"github.com/AidarBabanov/wallets-tracker/internal/addrdb/csv"
	"github.com/AidarBabanov/wallets-tracker/internal/explorer/etherscan"
	"github.com/AidarBabanov/wallets-tracker/internal/logging"
	"github.com/AidarBabanov/wallets-tracker/internal/notifier/telegram"
	"github.com/AidarBabanov/wallets-tracker/internal/tracker"
	"github.com/AidarBabanov/wallets-tracker/internal/uniswap"
	"github.com/beego/beego/v2/core/logs"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func main() {
	logging.Init()
	err := godotenv.Load() // load environment variables
	if err != nil {
		logs.Error(err)
	}
	addressDatabase := csv.New()
	err = addressDatabase.ReadCSV("addresses.csv")
	if err != nil {
		logs.Critical(err)
		os.Exit(-1)
	}

	exp := etherscan.New(os.Getenv("ETHERSCAN_API_KEY"))
	err = exp.Run()
	if err != nil {
		logs.Critical(err)
		os.Exit(-1)
	}
	defer exp.Close()

	swapViewer := uniswap.New()
	err = swapViewer.Run()
	if err != nil {
		logs.Critical(err)
		os.Exit(-1)
	}
	botToken := os.Getenv("TLG_NOTIFIER_TOKEN")
	chatID, err := strconv.Atoi(os.Getenv("TLG_NOTOFIER_CHAT_ID"))
	if err != nil {
		logs.Critical(err)
		os.Exit(-1)
	}
	notif, err := telegram.New(botToken, int64(chatID))
	if err != nil {
		logs.Critical(err)
		os.Exit(-1)
	}
	track := tracker.New(addressDatabase, exp, swapViewer, notif)
	err = track.LoadLastTransactions()
	if err != nil {
		logs.Error(err)
	}
	track.RunTracker()
}
