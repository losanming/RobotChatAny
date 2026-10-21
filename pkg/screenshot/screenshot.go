package screenshot

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/nfnt/resize"
)

// ScreenshotConfig 截图配置
type ScreenshotConfig struct {
	// 目标宽度，如果为0则保持原始宽度
	TargetWidth uint
	// 目标高度，如果为0则保持原始高度
	TargetHeight uint
	// JPEG压缩质量 (1-100)
	JpegQuality int
	// 显示器索引，-1表示所有显示器
	DisplayIndex int
}

// DefaultConfig 返回默认配置
func DefaultConfig() *ScreenshotConfig {
	return &ScreenshotConfig{
		TargetWidth:  1280,
		TargetHeight: 720,
		JpegQuality:  75,
		DisplayIndex: 0, // 主显示器
	}
}

// CaptureResult 截图结果
type CaptureResult struct {
	// 压缩后的图像数据
	ImageData []byte
	// 原始图像尺寸
	OriginalWidth int
	// 原始图像高度
	OriginalHeight int
	// 压缩后的图像尺寸
	CompressedWidth int
	// 压缩后的图像高度
	CompressedHeight int
	// 压缩后的数据大小（字节）
	CompressedSize int
	// 截图时间戳
	Timestamp int64
}

// Capture 捕获屏幕截图并压缩
func Capture(config *ScreenshotConfig) (*CaptureResult, error) {
	if config == nil {
		config = DefaultConfig()
	}

	var img image.Image
	var err error
	var bounds image.Rectangle

	// 获取显示器数量
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		return nil, fmt.Errorf("没有检测到活动的显示器")
	}

	// 检查显示器索引是否有效
	if config.DisplayIndex >= n {
		return nil, fmt.Errorf("无效的显示器索引: %d, 最大值: %d", config.DisplayIndex, n-1)
	}

	// 捕获指定显示器的截图
	if config.DisplayIndex >= 0 {
		bounds = screenshot.GetDisplayBounds(config.DisplayIndex)
		img, err = screenshot.CaptureRect(bounds)
		if err != nil {
			return nil, fmt.Errorf("截图失败: %w", err)
		}
	} else {
		// 捕获所有显示器的截图（这里只是示例，实际上需要更复杂的处理来合并多个显示器的图像）
		// 这里简化为只捕获第一个显示器
		bounds = screenshot.GetDisplayBounds(0)
		img, err = screenshot.CaptureRect(bounds)
		if err != nil {
			return nil, fmt.Errorf("截图失败: %w", err)
		}
	}

	// 记录原始尺寸
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// 调整图像大小
	var resized image.Image
	if config.TargetWidth > 0 && config.TargetHeight > 0 {
		resized = resize.Resize(config.TargetWidth, config.TargetHeight, img, resize.Lanczos3)
	} else if config.TargetWidth > 0 {
		// 保持宽高比，只指定宽度
		ratio := float64(originalHeight) / float64(originalWidth)
		targetHeight := uint(float64(config.TargetWidth) * ratio)
		resized = resize.Resize(config.TargetWidth, targetHeight, img, resize.Lanczos3)
	} else if config.TargetHeight > 0 {
		// 保持宽高比，只指定高度
		ratio := float64(originalWidth) / float64(originalHeight)
		targetWidth := uint(float64(config.TargetHeight) * ratio)
		resized = resize.Resize(targetWidth, config.TargetHeight, img, resize.Lanczos3)
	} else {
		// 不调整大小
		resized = img
	}

	// 压缩为JPEG
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resized, &jpeg.Options{
		Quality: config.JpegQuality,
	})
	if err != nil {
		return nil, fmt.Errorf("JPEG编码失败: %w", err)
	}

	// 获取压缩后的尺寸
	compressedSize := buf.Len()
	compressedBounds := resized.Bounds()
	compressedWidth := compressedBounds.Dx()
	compressedHeight := compressedBounds.Dy()

	return &CaptureResult{
		ImageData:       buf.Bytes(),
		OriginalWidth:   originalWidth,
		OriginalHeight:  originalHeight,
		CompressedWidth: compressedWidth,
		CompressedHeight:compressedHeight,
		CompressedSize:  compressedSize,
		Timestamp:       time.Now().Unix(),
	}, nil
}
