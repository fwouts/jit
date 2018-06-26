package config

// Config describes the format of the YAML config stored in ~/.jit/config.yaml.
type Config struct {
	Jira struct {
		Host  string
		User  string
		Token string
	}
}
