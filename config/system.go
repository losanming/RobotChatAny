package config

type System struct {
	Db      string `yaml:"db-type"`
	Oss     string `yaml:"oss-type"`
	Version string `yaml:"version"`
}
