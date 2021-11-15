package main

import (
	"flag"
	"fmt"
	"github.com/pavel-trbv/vk-message-counter/pkg"
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
		lang = pkg.DefaultLang
	}

	apiClient := pkg.NewHTTPAPIClient(
		*token,
		pkg.DefaultBaseUrl,
		lang,
		pkg.DefaultVersion,
	)

	service := pkg.NewService(apiClient, true)
	stats, err := service.GetMessageStats(chatIdInt)
	if err != nil {
		log.Fatalf("error occured while getting message stats: %s", err.Error())
	}

	var formatter pkg.Formatter = pkg.NewDefaultFormatter()
	fmt.Println(formatter.FormatText(stats))
}
