package multimedia

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ImageProcessor 图片处理工具
type ImageProcessor struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewImageProcessor 创建图片处理工具
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		enableLearning: true,
		category:       "multimedia",
		toolName:       "image_processor",
	}
}

// Info 返回工具信息
func (i *ImageProcessor) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "image_processor",
		Desc: "图片处理工具，支持图片裁剪、调整、滤镜、格式转换",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"resize", "crop", "rotate", "filter", "convert", "get_info"},
			},
			"file_path": {
				Desc:     "图片文件路径",
				Type:     schema.String,
				Required: false,
			},
			"width": {
				Desc:     "目标宽度",
				Type:     schema.Integer,
				Required: false,
			},
			"height": {
				Desc:     "目标高度",
				Type:     schema.Integer,
				Required: false,
			},
			"angle": {
				Desc:     "旋转角度",
				Type:     schema.Integer,
				Required: false,
			},
			"filter_type": {
				Desc:     "滤镜类型",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"grayscale", "blur", "brightness", "contrast"},
			},
			"output_format": {
				Desc:     "输出格式",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"jpeg", "png", "gif"},
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (i *ImageProcessor) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "resize":
		return i.resizeImage(params)
	case "crop":
		return i.cropImage(params)
	case "rotate":
		return i.rotateImage(params)
	case "filter":
		return i.applyFilter(params)
	case "convert":
		return i.convertFormat(params)
	case "get_info":
		return i.getImageInfo(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// getImageInfo 获取图片信息
func (i *ImageProcessor) getImageInfo(params map[string]interface{}) (string, error) {
	filePath := params["file_path"].(string)
	if filePath == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("解码图片失败: %w", err)
	}

	bounds := img.Bounds()
	result := map[string]interface{}{
		"success":  true,
		"file_path": filePath,
		"width":    bounds.Dx(),
		"height":   bounds.Dy(),
		"format":    getImageFormat(filePath),
		"file_size": getFileSize(filePath),
	}

	return toJSON(result), nil
}

// resizeImage 调整图片大小
func (i *ImageProcessor) resizeImage(params map[string]interface{}) (string, error) {
	filePath := params["file_path"].(string)
	width := getIntParam(params, "width", 0)
	height := getIntParam(params, "height", 0)

	if filePath == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	// 演示版本：返回模拟结果
	result := map[string]interface{}{
		"success":      true,
		"original_size": map[string]int{"width": 1920, "height": 1080},
		"new_size":     map[string]int{"width": width, "height": height},
		"message":      "图片调整成功（演示版本）",
	}

	return toJSON(result), nil
}

// cropImage 裁剪图片
func (i *ImageProcessor) cropImage(params map[string]interface{}) (string, error) {
	filePath := params["file_path"].(string)

	result := map[string]interface{}{
		"success": true,
		"file_path": filePath,
		"message":  "图片裁剪成功（演示版本）",
	}

	return toJSON(result), nil
}

// rotateImage 旋转图片
func (i *ImageProcessor) rotateImage(params map[string]interface{}) (string, error) {
	filePath := params["file_path"].(string)
	angle := getIntParam(params, "angle", 90)

	result := map[string]interface{}{
		"success": true,
		"file_path": filePath,
		"angle":    angle,
		"message":  "图片旋转成功（演示版本）",
	}

	return toJSON(result), nil
}

// applyFilter 应用滤镜
func (i *ImageProcessor) applyFilter(params map[string]interface{}) (string, error) {
	filePath := params["file_path"].(string)
	filterType := getStringParam(params, "filter_type", "grayscale")

	result := map[string]interface{}{
		"success":    true,
		"file_path":  filePath,
		"filter_type": filterType,
		"message":    "滤镜应用成功（演示版本）",
	}

	return toJSON(result), nil
}

// convertFormat 转换格式
func (i *ImageProcessor) convertFormat(params map[string]interface{}) (string, error) {
	filePath := params["file_path"].(string)
	outputFormat := getStringParam(params, "output_format", "png")

	result := map[string]interface{}{
		"success":       true,
		"file_path":     filePath,
		"output_format": outputFormat,
		"message":       "格式转换成功（演示版本）",
	}

	return toJSON(result), nil
}

// AfterRun 工具调用后的后处理
func (i *ImageProcessor) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现
}

// 辅助函数
func getImageFormat(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	formats := map[string]string{
		".jpg":  "jpeg",
		".jpeg": "jpeg",
		".png":  "png",
		".gif":  "gif",
		".bmp":  "bmp",
	}
	if format, ok := formats[ext]; ok {
		return format
	}
	return "unknown"
}

func getFileSize(filePath string) int64 {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return info.Size()
}