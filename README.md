# vk-message-counter
Simple golang script for getting VK message statistics
## Run as web application
```
npm install
npm run dev
```
Open http://localhost:3000 in your browser
## Usage from CLI
Use from CLI
```
go run cmd/cli/main.go -token=<string> -chat_id=<number> [-lang=<string>]
```
Something like this will be output:
```
Total count - 398
1) Some Person 1 - 85
2) Some Person 2 - 48
3) Some Person 3 - 30
...
```
## Usage from application
Example use from application is [here](https://github.com/pavel-trbv/vk-message-counter/blob/master/example/example.go)
## Thanks
https://github.com/go-vk-api/vk
