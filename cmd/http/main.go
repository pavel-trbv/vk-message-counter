package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	counter "github.com/pavel-trbv/vk-message-counter"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error load .env")
	}

	http.HandleFunc("/stats", handler)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		log.Fatal("http port is empty")
	}

	fmt.Printf("Server started at port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

type response struct {
	Message string `json:"message"`
}

func responseMessage(w http.ResponseWriter, statusCode int, message string) {
	res := response{Message: message}
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}

func handler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	chatId := r.URL.Query().Get("chat_id")
	inputLang := r.URL.Query().Get("lang")

	if token == "" || chatId == "" {
		responseMessage(w, http.StatusBadRequest, "empty token or chat id")
		return
	}

	chatIdInt, err := strconv.Atoi(chatId)
	if err != nil {
		responseMessage(w, http.StatusBadRequest, "fail convert chat id to int")
		return
	}

	lang := inputLang
	if lang == "" {
		lang = counter.DefaultLang
	}

	apiClient := counter.NewHTTPAPIClient(
		token,
		counter.DefaultBaseUrl,
		lang,
		counter.DefaultVersion,
	)

	service := counter.NewService(apiClient, false)
	stats, err := service.GetMessageStats(chatIdInt)
	if err != nil {
		fmt.Println(err.Error())
		responseMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	formatter := counter.NewDefaultFormatter()
	statsJson := formatter.FormatJson(stats)
	if err != nil {
		responseMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(statsJson))
}
