package counter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sort"
)

const (
	defaultBaseUrl = "https://api.vk.com/method"
	defaultLang    = "ru"
	defaultVersion = "5.131"

	tokenFieldName   = "access_token"
	versionFieldName = "v"
	langFieldName    = "lang"

	maxMessagesPerRequest = 200
)

type Counter struct {
	BaseURL string
	Lang    string
	Version string

	Token string

	HTTPClient *http.Client
}

type RequestParams map[string]interface{}

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
}

type GetChatResponse struct {
	Users []User `json:"users"`
}

type GetHistoryResponse struct {
	Count int `json:"count"`
	Items []struct {
		FromId int `json:"from_id"`
	} `json:"items"`
}

type MessageStats struct {
	TotalCount int
	List       map[string]int
}

func NewCounter(token string) *Counter {
	client := &Counter{
		BaseURL: defaultBaseUrl,
		Lang:    defaultLang,
		Version: defaultVersion,
		Token:   token,
	}

	client.HTTPClient = http.DefaultClient

	return client
}

func (c *Counter) GetMessageStats(chatId int, logging bool) (MessageStats, error) {
	var chat GetChatResponse
	err := c.CallMethod("messages.getChat", RequestParams{
		"chat_id": chatId,
		"fields":  "first_name,last_name",
	}, &chat)

	if err != nil {
		return MessageStats{}, err
	}

	var totalCount int
	counter := make(map[int]int)

	i := 0
	for {
		var history GetHistoryResponse
		err := c.CallMethod("messages.getHistory", RequestParams{
			"chat_id": chatId,
			"count":   maxMessagesPerRequest,
			"offset":  maxMessagesPerRequest * i,
		}, &history)
		if err != nil {
			return MessageStats{}, err
		}

		if len(history.Items) == 0 {
			break
		}

		for _, v := range history.Items {
			counter[v.FromId] += 1
		}

		totalCount = history.Count
		i += 1

		if logging {
			percent := float32(maxMessagesPerRequest*i) / float32(totalCount) * 100
			fmt.Printf("\r%d%%", int(percent))
		}
	}

	counterWithNames := make(map[string]int)
	for id, v := range counter {
		var user *User
		for _, u := range chat.Users {
			if u.Id == id {
				user = &u
				break
			}
		}

		if user == nil {
			continue
		}

		if user.Type == "group" {
			continue
		}

		fullName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		counterWithNames[fullName] = v
	}

	if logging {
		fmt.Println()
	}

	return MessageStats{
		TotalCount: totalCount,
		List:       counterWithNames,
	}, nil
}

func (c *Counter) CallMethod(method string, params RequestParams, response interface{}) error {
	values := url.Values{}

	for k, v := range params {
		values.Add(k, fmt.Sprint(v))
	}

	if !values.Has(tokenFieldName) {
		values.Add(tokenFieldName, c.Token)
	}

	if !values.Has(versionFieldName) {
		values.Add(versionFieldName, c.Version)
	}

	if !values.Has(langFieldName) {
		values.Add(langFieldName, c.Lang)
	}

	query := values.Encode()
	reqUrl := fmt.Sprintf("%s/%s?%s", c.BaseURL, method, query)

	resp, err := c.HTTPClient.Get(reqUrl)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var body struct {
		Response interface{}  `json:"response"`
		Error    *MethodError `json:"error"`
	}

	if response != nil {
		valueOfResponse := reflect.ValueOf(response)
		if valueOfResponse.Kind() != reflect.Ptr || valueOfResponse.IsNil() {
			return errors.New("invalid response")
		}

		body.Response = response
	}

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(rawBody, &body); err != nil {
		return err
	}

	if body.Error != nil {
		return body.Error
	}

	return err
}

type MethodError struct {
	Code          int64         `json:"error_code"`
	Message       string        `json:"error_msg"`
	RequestParams []interface{} `json:"request_params"`
}

func (err *MethodError) Error() string {
	return err.Message
}

func (s MessageStats) Format() string {
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range s.List {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	response := fmt.Sprintf("Total count - %d", s.TotalCount)
	for i, kv := range ss {
		index := i + 1
		response = fmt.Sprintf("%s\n%d) %s - %d", response, index, kv.Key, kv.Value)
	}

	return response
}
