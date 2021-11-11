package counter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

const (
	DefaultBaseUrl = "https://api.vk.com/method"
	DefaultLang    = "ru"
	DefaultVersion = "5.131"

	tokenFieldName   = "access_token"
	versionFieldName = "v"
	langFieldName    = "lang"

	maxMessagesPerRequest = 200
)

type requestParams map[string]interface{}

type chatUser struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
}

type ChatResponse struct {
	Users []chatUser `json:"users"`
}

type HistoryResponse struct {
	Count int `json:"count"`
	Items []struct {
		FromId int `json:"from_id"`
	} `json:"items"`
}

type APIClient interface {
	CallMethod(method string, params requestParams, response interface{}) error
	GetChat(chatId int) (ChatResponse, error)
	GetHistory(chatId, count, offset int) (HistoryResponse, error)
}

type HTTPAPIClient struct {
	BaseURL string
	Lang    string
	Version string
	Token   string

	HTTPClient *http.Client
}

func NewHTTPAPIClient(token string, baseURL string, lang string, version string) *HTTPAPIClient {
	return &HTTPAPIClient{
		BaseURL:    baseURL,
		Lang:       lang,
		Version:    version,
		Token:      token,
		HTTPClient: http.DefaultClient,
	}
}

func (c *HTTPAPIClient) GetChat(chatId int) (ChatResponse, error) {
	var chat ChatResponse
	err := c.CallMethod("messages.getChat", requestParams{
		"chat_id": chatId,
		"fields":  "first_name,last_name",
	}, &chat)

	return chat, err
}

func (c *HTTPAPIClient) GetHistory(chatId, count, offset int) (HistoryResponse, error) {
	var history HistoryResponse
	err := c.CallMethod("messages.getHistory", requestParams{
		"chat_id": chatId,
		"count":   count,
		"offset":  offset,
	}, &history)

	return history, err
}

func (c *HTTPAPIClient) CallMethod(method string, params requestParams, response interface{}) error {
	httpParams := url.Values{}

	for k, v := range params {
		httpParams.Add(k, fmt.Sprint(v))
	}

	if !httpParams.Has(tokenFieldName) {
		httpParams.Add(tokenFieldName, c.Token)
	}

	if !httpParams.Has(versionFieldName) {
		httpParams.Add(versionFieldName, c.Version)
	}

	if !httpParams.Has(langFieldName) {
		httpParams.Add(langFieldName, c.Lang)
	}

	query := httpParams.Encode()
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
