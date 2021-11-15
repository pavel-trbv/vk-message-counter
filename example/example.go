package main

import (
	"fmt"
	"github.com/joho/godotenv"
	counter "github.com/pavel-trbv/vk-message-counter"
	"log"
	"os"
	"strconv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error load .env")
	}

	token := os.Getenv("VK_TOKEN")
	chatId, err := strconv.Atoi(os.Getenv("VK_CHAT_ID"))
	if err != nil {
		log.Fatalf("error convert chat_id from env to int: %s", err.Error())
	}

	counterService := counter.Default(token)
	stats, err := counterService.GetMessageStats(chatId)
	if err != nil {
		log.Fatalf("error occured while getting message stats: %s", err.Error())
	}

	formatter := counter.NewDefaultFormatter()
	output := formatter.FormatText(stats)
	fmt.Println(output)
}
