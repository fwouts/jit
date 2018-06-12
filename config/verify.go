package config

func verify(config Config) bool {
	api, err := config.JiraClient()
	if err != nil {
		return false
	}
	_, _, err = api.Issue.Search("assignee = currentUser()", nil)
	if err != nil {
		return false
	}
	return true
}
