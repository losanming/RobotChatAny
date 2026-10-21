package config

type Minio struct {
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	AccessKeyId     string `json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret"`
	BucketName      string `json:"bucket_name" yaml:"bucket_name"`
	UseSSL          bool   `json:"use_ssl" yaml:"use_ssl"`
	BasePath        string `json:"base_path" yaml:"base_path"`
	BucketUrl       string `json:"bucket_url" yaml:"bucket_url"`
}
