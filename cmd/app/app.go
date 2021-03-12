package main

import (
	"github.com/AidarBabanov/wallets-tracker/internal/addrdb/csv"
	"github.com/AidarBabanov/wallets-tracker/internal/explorer/etherscan"
	"github.com/AidarBabanov/wallets-tracker/internal/logging"
	"github.com/AidarBabanov/wallets-tracker/internal/notifier/telegram"
	"github.com/AidarBabanov/wallets-tracker/internal/rest"
	"github.com/AidarBabanov/wallets-tracker/internal/tracker"
	"github.com/AidarBabanov/wallets-tracker/internal/uniswap"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	logging.Init()
	err := godotenv.Load() // load environment variables
	if err != nil {
		logs.Error(err)
	}
	addressDatabase := csv.New("addresses.csv")
	err = addressDatabase.ReadCSV()
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

	api := rest.API{
		Router: mux.NewRouter(),
		ApiKey: os.Getenv("API_KEY"),
		AddrDB: addressDatabase,
	}
	api.Handle()
	httpServer := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: api,
	}
	go func() {
		logs.Info("Maintaining the web application...")
		logs.Critical(httpServer.ListenAndServe())
	}()

	go func() {
		logs.Info("Started tracker session.")
		notif.Notify("Started tracker session.")
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		notif.Notify("Ended tracker session.")
		os.Exit(0)
	}()

	track.RunTracker()
}
