package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"new-robot-self/model"
	"new-robot-self/public"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CaptchaService 验证码服务
type CaptchaService struct{}

// 验证码缓存key前缀
const (
	SliderCaptchaPrefix = "slider_captcha:"
	SliderSuccessPrefix = "slider_success:"
)

// GenerateSliderCaptcha 生成滑块验证码
func (s *CaptchaService) GenerateSliderCaptcha() (*model.SliderCaptcha, error) {
	// 1. 随机选择背景图片
	originalImage, err := s.loadRandomBackgroundImage()
	if err != nil {
		return nil, err
	}

	// 2. 随机生成缺口位置(X坐标)
	xPosition := rand.Intn(150) + 100 // 100-250之间的随机数

	// 3. 生成滑块图片和带缺口的背景图
	sliderImage, processedOriginal, err := s.processImages(originalImage, xPosition)
	if err != nil {
		return nil, err
	}

	// 4. 生成验证token
	token := s.generateToken()

	// 5. 设置过期时间(当前时间戳+5分钟)
	expireTime := time.Now().Unix() + 300

	// 6. 将验证信息存入缓存(Redis)
	err = s.cacheVerificationInfo(token, xPosition, expireTime)
	if err != nil {
		return nil, err
	}

	// 7. 转换图片为base64
	originalBase64 := s.imageToBase64(processedOriginal)
	sliderBase64 := s.imageToBase64(sliderImage)

	return &model.SliderCaptcha{
		Token:         token,
		OriginalImage: originalBase64,
		SliderImage:   sliderBase64,
		XPosition:     xPosition,
		ExpireTime:    expireTime,
	}, nil
}

// VerifySliderCaptcha 验证滑块位置
func (s *CaptchaService) VerifySliderCaptcha(token string, userXPosition int) (*model.VerificationResult, error) {
	// 1. 从缓存中获取正确的位置信息
	correctPosition, expireTime, err := s.getVerificationInfo(token)
	if err != nil {
		return nil, err
	}

	// 2. 检查是否过期
	if time.Now().Unix() > expireTime {
		return &model.VerificationResult{
			Success: false,
			Message: "验证已过期",
		}, nil
	}

	// 3. 验证位置是否正确(允许一定误差)
	if math.Abs(float64(userXPosition-correctPosition)) <= 5 {
		// 验证成功，生成验证通过的标记
		successToken := s.generateSuccessToken()
		// 将成功标记存入缓存
		err = s.cacheSuccessToken(successToken)
		if err != nil {
			return nil, err
		}

		return &model.VerificationResult{
			Success: true,
			Message: "验证通过",
			Token:   successToken,
		}, nil
	}

	return &model.VerificationResult{
		Success: false,
		Message: "验证失败",
	}, nil
}

// CheckCaptchaToken 检查验证成功的token是否有效
func (s *CaptchaService) CheckCaptchaToken(token string) (bool, error) {
	key := SliderSuccessPrefix + token
	exists, err := public.Cache.Exists(public.Ctx, key)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// 加载随机背景图片
func (s *CaptchaService) loadRandomBackgroundImage() (image.Image, error) {
	// 在实际应用中，应该从配置的目录中随机选择一张图片
	// 这里为了简单，使用一个固定的图片路径
	imagePath := "public/captcha/background.png"

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		// 确保目录存在
		dir := filepath.Dir(imagePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}

		// 创建一个简单的背景图片
		img := image.NewRGBA(image.Rect(0, 0, 300, 200))
		blue := color.RGBA{0, 0, 255, 255}
		draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)

		// 保存图片
		f, err := os.Create(imagePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if err := png.Encode(f, img); err != nil {
			return nil, err
		}

		return img, nil
	}

	// 打开图片文件
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// 处理图片，生成滑块和带缺口的背景图
func (s *CaptchaService) processImages(originalImage image.Image, xPosition int) (image.Image, image.Image, error) {
	// 获取原图尺寸
	bounds := originalImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 创建滑块图片(一个简单的矩形)
	sliderWidth, sliderHeight := 50, 50
	yPosition := height/2 - sliderHeight/2

	// 创建新的RGBA图像
	slider := image.NewRGBA(image.Rect(0, 0, sliderWidth, sliderHeight))
	processed := image.NewRGBA(image.Rect(0, 0, width, height))

	// 复制原图到processed
	draw.Draw(processed, processed.Bounds(), originalImage, image.Point{}, draw.Src)

	// 创建滑块图像
	for y := 0; y < sliderHeight; y++ {
		for x := 0; x < sliderWidth; x++ {
			// 从原图对应位置复制像素到滑块
			slider.Set(x, y, originalImage.At(x+xPosition, y+yPosition))

			// 在原图对应位置绘制一个半透明的灰色块
			if x < sliderWidth-1 && y < sliderHeight-1 {
				processed.Set(x+xPosition, y+yPosition, color.RGBA{200, 200, 200, 150})
			}
		}
	}

	return slider, processed, nil
}

// 生成随机token
func (s *CaptchaService) generateToken() string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// 生成验证成功的token
func (s *CaptchaService) generateSuccessToken() string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return fmt.Sprintf("success_%x", bytes)
}

// 缓存验证信息
func (s *CaptchaService) cacheVerificationInfo(token string, position int, expireTime int64) error {
	key := SliderCaptchaPrefix + token
	value := fmt.Sprintf("%d:%d", position, expireTime)
	return public.Cache.Set(public.Ctx, key, value, 5*time.Minute)
}

// 获取验证信息
func (s *CaptchaService) getVerificationInfo(token string) (int, int64, error) {
	key := SliderCaptchaPrefix + token
	var value string
	err := public.Cache.Get(public.Ctx, key, value)
	if err != nil {
		return 0, 0, errors.New("验证信息不存在或已过期")
	}

	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return 0, 0, errors.New("验证信息格式错误")
	}

	var position int
	var expireTime int64
	fmt.Sscanf(parts[0], "%d", &position)
	fmt.Sscanf(parts[1], "%d", &expireTime)

	return position, expireTime, nil
}

// 缓存成功token
func (s *CaptchaService) cacheSuccessToken(token string) error {
	key := SliderSuccessPrefix + token
	return public.Cache.Set(public.Ctx, key, "1", 5*time.Minute)
}

// 将图片转换为base64
func (s *CaptchaService) imageToBase64(img image.Image) string {
	var buf strings.Builder

	// 将图片编码为PNG
	err := png.Encode(&buf, img)
	if err != nil {
		return ""
	}

	// 转换为base64
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(buf.String()))
}
