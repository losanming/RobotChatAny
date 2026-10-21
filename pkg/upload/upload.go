package upload

import (
	"log/slog"
	"mime/multipart"
	"new-robot-self/public"
)

// OSS 对象存储接口
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [ccfish86](https://github.com/ccfish86)
type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
	DownloadFileByByte(key string) ([]byte, error)
	ListBuckets() ([]string, error)
}

// NewOss OSS的实例化方法
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [ccfish86](https://github.com/ccfish86)
func NewOss(t string) OSS {
	switch t {
	case "local":
		return &Local{}
	case "huawei":
		return HuaWeiObs
	case "minio":
		minioClient, err := GetMinio(public.Config.MinIo.Endpoint, public.Config.MinIo.AccessKeyId, public.Config.MinIo.AccessKeySecret, public.Config.MinIo.BucketName, public.Config.MinIo.UseSSL)
		if err != nil {
			slog.Warn("你配置了使用minio，但是初始化失败，请检查minio可用性或安全配置: " + err.Error())
			panic("minio初始化失败") // 建议这样做，用户自己配置了minio，如果报错了还要把服务开起来，使用起来也很危险
		}
		return minioClient
	default:
		return &Local{}
	}
}
