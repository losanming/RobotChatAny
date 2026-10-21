package config

// RAGConfig RAG相关配置
type RAGConfig struct {
	Embedding EmbeddingConfig `yaml:"embedding"`
	Qdrant   QdrantConfig   `yaml:"qdrant"`
}

// EmbeddingConfig Embedding服务配置
type EmbeddingConfig struct {
	Provider            string  `yaml:"provider"`
	BaseURL             string  `yaml:"baseURL"`
	APIKey              string  `yaml:"apiKey"`
	Model               string  `yaml:"model"`
	EndpointPath        string  `yaml:"endpointPath"`
	TimeoutSeconds      int     `yaml:"timeoutSeconds"`
	Dimensions          int     `yaml:"dimensions"`
	QueryInstruction    string  `yaml:"queryInstruction"`
	DocumentInstruction string  `yaml:"documentInstruction"`
}

// QdrantConfig Qdrant向量数据库配置
type QdrantConfig struct {
	BaseURL         string  `yaml:"baseURL"`
	APIKey          string  `yaml:"apiKey"`
	Collection      string  `yaml:"collection"`
	TimeoutSeconds  int     `yaml:"timeoutSeconds"`
	TopK            int     `yaml:"topK"`
	ScoreThreshold  float64 `yaml:"scoreThreshold"`
}
