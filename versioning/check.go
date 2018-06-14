package versioning

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"strings"
)

// CheckNewRelease shows a log message to the user if a new version is available.
func CheckNewRelease() {
	releases, err := FetchReleases()
	if err != nil {
		fmt.Printf("Could not check latest `jit` releases from GitHub. Are you offline?\n")
		return
	}
	if len(*releases) == 0 {
		// No releases available on GitHub.
		fmt.Printf("It looks like no releases of `jit` are available on GitHub yet.\n")
		return
	}
	lastRelease := (*releases)[0]
	current, err := version.NewVersion(strings.TrimPrefix(CurrentVersion, "v"))
	last, err := version.NewVersion(strings.TrimPrefix(lastRelease.TagName, "v"))
	if current.LessThan(last) {
		message := fmt.Sprintf("You are using `jit` %s. A newer version (%s) is available. Download it from %s.\n", CurrentVersion, lastRelease.TagName, lastRelease.URL)
		coloredMessage := color.MagentaString(message)
		fmt.Printf(coloredMessage)
	}
}
