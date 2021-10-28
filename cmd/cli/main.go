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
	inputLang := flag.String("lang", "", "")
	flag.Parse()

	if *token == "" || *chatId == "" {
		log.Fatal("empty arguments")
	}

	chatIdInt, err := strconv.Atoi(*chatId)
	if err != nil {
		log.Fatalf("error convert chat_id from env to int: %s", err.Error())
	}

	lang := *inputLang
	if lang == "" {
		lang = counter.DefaultLang
	}

	apiClient := counter.NewHTTPAPIClient(
		*token,
		counter.DefaultBaseUrl,
		lang,
		counter.DefaultVersion,
	)

	service := counter.NewService(apiClient, true)
	stats, err := service.GetMessageStats(chatIdInt)
	if err != nil {
		log.Fatalf("error occured while getting message stats: %s", err.Error())
	}

	var formatter counter.Formatter = counter.NewDefaultFormatter()
	output := formatter.Format(stats)
	fmt.Println(output)
}
