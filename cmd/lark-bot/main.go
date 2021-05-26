package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sysu-ecnc/lark-bot/internal/app"
	"github.com/sysu-ecnc/lark-bot/internal/pkg/lark"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
		return
	}

	err = lark.InitConfig(
		os.Getenv("LARK_APP_ID"),
		os.Getenv("LARK_APP_SECRET"),
		os.Getenv("LARK_VERIFICATION_TOKEN"),
		os.Getenv("LARK_ENCRYPT_KEY"),
	)
	if err != nil {
		log.Fatal("Error init Lark config: ", err.Error())
		return
	}

	app.Start()
}
