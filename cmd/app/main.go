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

	client := counter.NewCounter(os.Getenv("VK_TOKEN"))

	chatId, err := strconv.Atoi(os.Getenv("VK_CHAT_ID"))
	if err != nil {
		log.Fatalf("error convert chat_id from env to int: %s", err.Error())
	}

	stats, err := client.GetMessageStats(chatId, true)
	if err != nil {
		log.Fatalf("error occured while getting message stats: %s", err.Error())
	}

	prettyStats := stats.Format()
	fmt.Println(prettyStats)
}