package config

type NsqConfig struct {
	Nsqd        string   `yaml:"nsqd"`
	Nsqlookupd  []string `yaml:"nsqlookupd"`
	Topic       string   `yaml:"topic"`
	Channel     string   `yaml:"channel"`
	Dlq         string   `yaml:"dlq"`
	Workers     int      `yaml:"workers"`
	MaxInFlight int      `yaml:"max_in_flight"`
	MaxAttempts uint16   `yaml:"max_attempts"`
	MsgTimeout  int      `yaml:"msg_timeout"`
	Off         bool     `yaml:"off"` // 是否开启截图的功能
}
