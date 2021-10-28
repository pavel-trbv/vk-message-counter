package counter

func Default(token string) *Service {
	apiClient := NewHTTPAPIClient(token, DefaultBaseUrl, DefaultLang, DefaultVersion)
	service := NewService(apiClient, false)

	return service
}
