package upload

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/pkg/errors"
	"mime/multipart"
	"new-robot-self/pkg"
	"new-robot-self/public"
	"time"
)

var HuaWeiObs = new(Obs)

type Obs struct{}

func NewHuaWeiObsClient() (client *obs.ObsClient, err error) {
	return obs.New(public.Config.HuaweiCloudOBS.AccessKey, public.Config.HuaweiCloudOBS.SecretKey, public.Config.HuaweiCloudOBS.Endpoint)
}

func (o *Obs) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// var open multipart.File
	open, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer open.Close()
	file_path := public.Config.HuaweiCloudOBS.Path + "/" + "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + pkg.GetRandStr(8) + file.Filename
	input := &obs.PutObjectInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket: public.Config.HuaweiCloudOBS.BucketName,
				Key:    file_path,
			},
			HttpHeader: obs.HttpHeader{
				ContentType: file.Header.Get("content-type"),
			},
		},
		Body: open,
	}

	var client *obs.ObsClient
	client, err = NewHuaWeiObsClient()
	if err != nil {
		return "", "", errors.Wrap(err, "获取华为对象存储对象失败!")
	}

	_, err = client.PutObject(input)
	if err != nil {
		return "", "", errors.Wrap(err, "文件上传失败!")
	}
	//生成带有预签名的url
	var sign_in = obs.CreateSignedUrlInput{
		Bucket:  public.Config.HuaweiCloudOBS.BucketName,
		Key:     file_path,
		Method:  "GET",
		Expires: 630710000,
	}
	url, err := client.CreateSignedUrl(&sign_in)
	if err != nil {
		return "", "", errors.Wrap(err, "生成带有预签名的url失败!")
	}
	return url.SignedUrl, file_path, err
}

func (o *Obs) DeleteFile(key string) error {
	client, err := NewHuaWeiObsClient()
	if err != nil {
		return errors.Wrap(err, "获取华为对象存储对象失败!")
	}
	input := &obs.DeleteObjectInput{
		Bucket: public.Config.HuaweiCloudOBS.BucketName,
		Key:    key,
	}
	var output *obs.DeleteObjectOutput
	output, err = client.DeleteObject(input)
	if err != nil {
		return errors.Wrapf(err, "删除对象(%s)失败!, output: %v", key, output)
	}
	return nil
}

func (o *Obs) DownloadFile(key string) (string, error) {
	client, err := NewHuaWeiObsClient()
	if err != nil {
		return "", errors.Wrap(err, "获取华为对象存储对象失败!")
	}
	var sign_in = obs.CreateSignedUrlInput{
		Bucket:  public.Config.HuaweiCloudOBS.BucketName,
		Key:     key,
		Method:  "GET",
		Expires: 31536000,
	}
	url, err := client.CreateSignedUrl(&sign_in)
	if err != nil {
		return "", errors.Wrap(err, "生成带有预签名的url失败!")
	}
	return url.SignedUrl, nil
}
func (o *Obs) DownloadFileByByte(key string) ([]byte, error) {
	client, err := NewHuaWeiObsClient()
	if err != nil {
		return nil, errors.Wrap(err, "获取华为对象存储对象失败!")
	}
	input := &obs.GetObjectMetadataInput{
		Bucket: public.Config.HuaweiCloudOBS.BucketName,
		Key:    key,
	}
	output, err := client.GetObjectMetadata(input)
	if err != nil {
		return nil, errors.Wrap(err, "生成带有预签名的url失败!")
	}
	fmt.Printf("Metadata:%v\n", output.Metadata)

	return nil, nil
}

func (o *Obs) ListBuckets() ([]string, error) {
	client, err := NewHuaWeiObsClient()
	if err != nil {
		return nil, errors.Wrap(err, "获取华为对象存储对象失败!")
	}
	input := &obs.ListBucketsInput{}
	// 指定桶列表中是否存在Location字段，此处以true为例，默认是false
	input.QueryLocation = true
	// 指定桶类型，此处以为obs.OBJECT为例，即获取所有桶列表，默认不设置则获取所有桶和并行文件系统列表
	input.BucketType = obs.OBJECT
	output, err := client.ListBuckets(input)
	if err != nil {
		return nil, errors.Wrap(err, "生成带有预签名的url失败!")
	}
	var res []string
	for _, v := range output.Buckets {
		res = append(res, v.Name)
	}
	return res, nil
}
