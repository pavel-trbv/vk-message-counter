package counter

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Formatter interface {
	FormatText(stats MessageStats) string
	FormatJson(stats MessageStats) string
}

type DefaultFormatter struct {
}

func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{}
}

func (f *DefaultFormatter) FormatText(stats MessageStats) string {
	type kv struct {
		Key   string
		Value int
	}

	var statsList []kv
	for k, v := range stats.List {
		statsList = append(statsList, kv{k, v})
	}

	sort.Slice(statsList, func(i, j int) bool {
		return statsList[i].Value > statsList[j].Value
	})

	response := fmt.Sprintf("Total count - %d", stats.TotalCount)
	for i, kv := range statsList {
		index := i + 1
		response = fmt.Sprintf("%s\n%d) %s - %d", response, index, kv.Key, kv.Value)
	}

	return response
}

func (f *DefaultFormatter) FormatJson(stats MessageStats) string {
	type item struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	type output struct {
		TotalCount int    `json:"totalCount"`
		List       []item `json:"list"`
	}

	type kv struct {
		Key   string
		Value int
	}

	var statsList []kv
	for k, v := range stats.List {
		statsList = append(statsList, kv{k, v})
	}

	sort.Slice(statsList, func(i, j int) bool {
		return statsList[i].Value > statsList[j].Value
	})

	var list []item
	for _, v := range statsList {
		list = append(list, item{
			Name:  v.Key,
			Count: v.Value,
		})
	}

	response := output{
		TotalCount: stats.TotalCount,
		List:       list,
	}

	body, err := json.Marshal(response)
	if err != nil {
		return ""
	}

	return string(body)
}
