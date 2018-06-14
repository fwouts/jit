package versioning

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const gitUser = "zenclabs"
const gitRepo = "jit"

var releasesURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", gitUser, gitRepo)

// Release is one particular release of `jit`.
type Release struct {
	URL     string `json:"url"`
	TagName string `json:"tag_name"`
}

// FetchReleases returns the list of available `jit` releases from GitHub.
func FetchReleases() (*[]Release, error) {
	response, err := http.Get(releasesURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	releases := make([]Release, 0)
	err = json.NewDecoder(response.Body).Decode(&releases)
	if err != nil {
		return nil, err
	}
	return &releases, nil
}
