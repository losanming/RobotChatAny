package http_client

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"time"
)

// HTTPClient 封装resty.Client
var HTTPClient = resty.New()

func init() {
	HTTPClient.SetTimeout(10 * time.Second)
}

// HTTPResponse 统一响应结构体
// Data字段为接口返回的原始内容
// Err为请求或解析过程中的错误
// StatusCode为HTTP状态码
// Raw为原始响应体
type HTTPResponse struct {
	Data       interface{}
	Err        error
	StatusCode int
	Raw        []byte
}

// Get 发送GET请求，并自动解析JSON响应
func Get(ctx context.Context, url string, headers map[string]string, respObj interface{}) *HTTPResponse {
	resp, err := HTTPClient.R().
		SetContext(ctx).
		SetHeaders(headers).
		Get(url)
	if err != nil {
		return &HTTPResponse{Err: err}
	}
	if respObj != nil {
		err = json.Unmarshal(resp.Body(), respObj)
	}
	return &HTTPResponse{
		Data:       respObj,
		Err:        err,
		StatusCode: resp.StatusCode(),
		Raw:        resp.Body(),
	}
}

// PostJSON 发送POST JSON请求，并自动解析JSON响应
func PostJSON(ctx context.Context, url string, headers map[string]string, body interface{}, respObj interface{}) *HTTPResponse {
	resp, err := HTTPClient.R().
		SetContext(ctx).
		SetHeaders(headers).
		SetBody(body).
		Post(url)
	if err != nil {
		return &HTTPResponse{Err: err}
	}
	if respObj != nil {
		err = json.Unmarshal(resp.Body(), respObj)
	}
	return &HTTPResponse{
		Data:       respObj,
		Err:        err,
		StatusCode: resp.StatusCode(),
		Raw:        resp.Body(),
	}
}

// PostForm 发送POST表单请求（application/x-www-form-urlencoded）
func PostForm(ctx context.Context, url string, headers map[string]string, form map[string]string, respObj interface{}) *HTTPResponse {
	resp, err := HTTPClient.R().
		SetContext(ctx).
		SetHeaders(headers).
		SetFormData(form).
		Post(url)
	if err != nil {
		return &HTTPResponse{Err: err}
	}
	if respObj != nil {
		err = json.Unmarshal(resp.Body(), respObj)
	}
	return &HTTPResponse{
		Data:       respObj,
		Err:        err,
		StatusCode: resp.StatusCode(),
		Raw:        resp.Body(),
	}
}
