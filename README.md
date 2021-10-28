# vk-message-counter
Simple golang script for getting VK message statistics
## Example
```go
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
```
## Result
Something like this will be output to the console:
```
Total count - 398
1) Some Person 1 - 85
2) Some Person 2 - 48
3) Some Person 3 - 30
...
```
## Thanks
https://github.com/go-vk-api/vk
