package main

import (
	"flag"
	"fmt"
	counter "github.com/pavel-trbv/vk-message-counter"
	"log"
	"strconv"
)

func main() {
	token := flag.String("token", "", "")
	chatId := flag.String("chat_id", "", "")
	flag.Parse()

	if *token == "" || *chatId == "" {
		log.Fatal("empty arguments")
	}

	chatIdInt, err := strconv.Atoi(*chatId)
	if err != nil {
		log.Fatalf("error convert chat_id from env to int: %s", err.Error())
	}

	client := counter.NewCounter(*token)
	stats, err := client.GetMessageStats(chatIdInt, true)
	if err != nil {
		log.Fatalf("error occured while getting message stats: %s", err.Error())
	}

	prettyStats := stats.Format()
	fmt.Println(prettyStats)
}
