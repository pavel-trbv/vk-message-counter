package pkg

import (
	"fmt"
)

type Service struct {
	Logging   bool
	APIClient APIClient
}

type MessageStats struct {
	TotalCount int
	List       map[string]int
}

func NewService(APIClient APIClient, Logging bool) *Service {
	client := &Service{
		Logging:   Logging,
		APIClient: APIClient,
	}

	return client
}

func (s *Service) GetMessageStats(chatId int) (MessageStats, error) {
	chat, err := s.APIClient.GetChat(chatId)
	if err != nil {
		return MessageStats{}, err
	}

	var totalCount int
	counter := make(map[int]int)

	i := 0
	for {
		history, err := s.APIClient.GetHistory(chatId, maxMessagesPerRequest, maxMessagesPerRequest*i)
		if err != nil {
			return MessageStats{}, err
		}

		if len(history.Items) == 0 {
			break
		}

		totalCount = history.Count
		for _, v := range history.Items {
			counter[v.FromId] += 1
		}

		if s.Logging {
			percent := float32(maxMessagesPerRequest*i) / float32(totalCount) * 100
			fmt.Printf("\r%d%%", int(percent))
		}

		i += 1
	}

	counterWithNames := make(map[string]int)
	for id, v := range counter {
		var user *chatUser
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

	if s.Logging {
		fmt.Println()
	}

	return MessageStats{
		TotalCount: totalCount,
		List:       counterWithNames,
	}, nil
}
