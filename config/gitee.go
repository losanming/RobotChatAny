package config

type GiteeConfig struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectUri  string `yaml:"redirect_uri"`
	Scope        string `yaml:"scope"`
}
