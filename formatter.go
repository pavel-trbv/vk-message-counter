package counter

import (
	"fmt"
	"sort"
)

type Formatter interface {
	Format(stats MessageStats) string
}

type DefaultFormatter struct {
}

func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{}
}

func (f *DefaultFormatter) Format(stats MessageStats) string {
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
