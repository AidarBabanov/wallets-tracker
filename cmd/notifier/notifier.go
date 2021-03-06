package main

import (
	"github.com/AidarBabanov/wallets-tracker/internal/logging"
	"github.com/AidarBabanov/wallets-tracker/internal/notifier/telegram"
	"github.com/beego/beego/v2/core/logs"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type subModel struct {
	SubMsg string `yaml:"sub_msg"`
}

type aModel struct {
	Msg       string   `yaml:"msg"`
	SubStruct subModel `yaml:"sub_struct"`
}

func main() {
	logging.Init()
	err := godotenv.Load()
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

	notif.Notify(aModel{
		Msg: "test msg",
		SubStruct: subModel{
			SubMsg: "test msg 2",
		},
	})
}
