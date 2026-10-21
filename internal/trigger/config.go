package trigger

import (
	"encoding/json"
	"io/ioutil"
)

// Config 触发器配置
type Config struct {
	TimePoints   []TimePointTrigger   `json:"time_points"`
	RealTime     []RealTimeTrigger    `json:"real_time"`
	UserBehavior []UserBehaviorTrigger `json:"user_behavior"`
	RandomTopics []RandomTopicTrigger  `json:"random_topics"`
}

// LoadFromFile 从文件加载配置
func LoadFromFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveToFile 保存配置到文件
func (c *Config) SaveToFile(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
