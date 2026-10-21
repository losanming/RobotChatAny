package quality

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"
)

// Result 表示一次质量评估的结果
// Score 0-100; IsJunk 表示是否判定为垃圾图片（不建议上传云端）
// Reasons 记录判定依据，便于日志排查
// Hash 用于与前一张图片对比（aHash）
// Sharpness/Contrast 等为主要指标，便于埋点

type Result struct {
	Score      int
	IsJunk     bool
	Reasons    []string
	Hash       uint64
	Sharpness  float64
	Contrast   float64
	Uniformity float64
	Mean       float64
	Flatness   float64
	EdgeDensity float64
	ChangeHamming int
	BlockVarRatio float64
	Colorfulness  float64
}

// Analyzer 质量评估器，带状态（上一张图片的哈希值）
// 阈值均可根据实际数据再调整

type Analyzer struct {
	// 阈值
	MinSharpness float64 // 拉普拉斯方差阈值，过低判定为模糊
	MinContrast  float64 // 亮度标准差阈值，过低判定为低对比度
	MaxUniform   float64 // 图像均匀度阈值，过高意味着近似全黑/全白/大面积单色
	MinChange    int     // 与上一帧 aHash 的最小海明距离，过低认为重复

	// 新增阈值
	MaxFlatness       float64 // 平坦度（接近均值像素占比）过高则近单色
	MinEdgeDensity    float64 // 边缘密度过低且对比度低，视作无效
	NearBlackMean     float64 // 近黑屏均值阈值
	NearWhiteMean     float64 // 近白屏均值阈值
	LowContrastStrong float64 // 强低对比度阈值
	// 桌面场景专用
	UniformMax        float64 // 近单色判据之一（熵反向）
	BlockVarGrid      int     // Block 网格
	BlockVarThr       float64 // Block 方差阈值（0-255 灰度方差）
	BlockVarMaxRatio  float64 // 小方差块占比阈值（纯色/灰屏）
	MinContrastStrong float64 // 无结构判据搭配的更低对比度阈值
	MinChangeStrong   int     // 重复帧强判据

	// 评分权重（总分100）
	WeightSharpness int
	WeightContrast  int
	WeightChange    int
	WeightContent   int // 内容丰富度（由均匀度反向度量）

	// 状态
	lastHash    uint64
	initialized bool
}

// NewAnalyzer 返回带默认阈值与权重的 Analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		MinSharpness:    25.0,
		MinContrast:     18.0,
		MaxUniform:      0.92,
		MinChange:       6,
		MaxFlatness:     0.88,
		MinEdgeDensity:  0.005,
		NearBlackMean:   5.0,
		NearWhiteMean:   250.0,
		LowContrastStrong: 10.0,
		UniformMax:      0.95,
		BlockVarGrid:    16,
		BlockVarThr:     20.0,
		BlockVarMaxRatio: 0.85,
		MinContrastStrong: 16.0,
		MinChangeStrong:  4,
		WeightSharpness: 35,
		WeightContrast:  25,
		WeightChange:    25,
		WeightContent:   15,
	}
}

// AnalyzeBytes 从 JPEG/PNG/GIF 等字节流评估质量
func (a *Analyzer) AnalyzeBytes(b []byte) (Result, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return Result{}, err
	}
	return a.AnalyzeImage(img)
}

// AnalyzeJPEGReader 针对 JPEG 流更高效
func (a *Analyzer) AnalyzeJPEGReader(r io.Reader) (Result, error) {
	img, err := jpeg.Decode(r)
	if err != nil {
		return Result{}, err
	}
	return a.AnalyzeImage(img)
}

// AnalyzeImage 对解码后的 image.Image 进行评估
func (a *Analyzer) AnalyzeImage(img image.Image) (Result, error) {
	if img == nil {
		return Result{}, errors.New("nil image")
	}
	g := toGray(img)

	// 指标计算
	sharp := laplacianVariance(g)
	contrast := luminanceStdDev(g)
	uniform := uniformity(g)
	hash := averageHash64(g)
	mean := luminanceMean(g)
	flat := flatnessRatio(g, 5)
	edges := edgeDensity(g, 25)
	blkRatio := blockVarianceRatio(g, a.BlockVarGrid, a.BlockVarThr)
	colorf := colorfulnessRGB(img)

	reasons := []string{}
	junk := false

	if sharp < a.MinSharpness {
		reasons = append(reasons, "模糊: sharpness below threshold")
	}
	if contrast < a.MinContrast {
		reasons = append(reasons, "低对比度: contrast below threshold")
	}
	if uniform > a.MaxUniform {
		reasons = append(reasons, "近似单色/全黑全白: uniformity too high")
	}

	// 与上一帧对比（若已初始化）
	changeScore := 1.0
	hdist := 0
	if a.initialized {
		hdist = hammingDistance64(a.lastHash, hash)
		if hdist < a.MinChange {
			reasons = append(reasons, "重复内容: similar to previous frame")
		}
		// 变化程度的归一化分数（0-1），以 32 位差异为满分
		if hdist > 32 {
			hdist = 32
		}
		changeScore = float64(hdist) / 32.0
	}

	// 内容丰富度（均匀度越高越差），将其反向映射
	contentScore := 1.0 - clamp(uniform, 0.0, 1.0)

	// 近黑/近白/近灰与边缘稀少检测
	if mean <= a.NearBlackMean {
		reasons = append(reasons, "近黑屏: luminance mean very low")
	}
	if mean >= a.NearWhiteMean {
		reasons = append(reasons, "近白屏: luminance mean very high")
	}
	if flat >= a.MaxFlatness && contrast <= a.MinContrast*0.8 {
		reasons = append(reasons, "近单色/灰屏: flatness high & contrast low")
	}
	if edges <= a.MinEdgeDensity && contrast <= a.MinContrast {
		reasons = append(reasons, "边缘稀少: edge density very low")
	}
	if uniform >= a.UniformMax {
		reasons = append(reasons, "近单色: uniformity very high")
	}
	if blkRatio >= a.BlockVarMaxRatio {
		reasons = append(reasons, "大片纯色: low-variance blocks ratio high")
	}
	if colorf <= 7.0 && (edges <= 0.006 || contrast <= 18.0) {
		reasons = append(reasons, "低色彩: colorfulness low")
	}

	// 汇总评分
	score := int(
		float64(a.WeightSharpness)*normalize(sharp, a.MinSharpness, a.MinSharpness*3.0)+
			float64(a.WeightContrast)*normalize(contrast, a.MinContrast, a.MinContrast*3.0)+
			float64(a.WeightChange)*changeScore+
			float64(a.WeightContent)*contentScore,
	) / 1 // 保持整型

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	// 判定垃圾图片：任一硬阈值严重不达标，或综合得分过低
	if sharp < a.MinSharpness*0.6 || contrast < a.MinContrast*0.6 || uniform > 0.98 {
		junk = true
	}
	// 严重近黑/近白
	if mean <= a.NearBlackMean || mean >= a.NearWhiteMean {
		junk = true
	}
	// 近灰/单色大片 & 低对比度
	if flat >= a.MaxFlatness && contrast <= a.MinContrast*0.8 {
		junk = true
	}
	// 边缘极少且对比度低
	if edges <= a.MinEdgeDensity && contrast <= a.MinContrast*0.9 {
		junk = true
	}
	// 近单色（熵反向）
	if uniform >= a.UniformMax {
		junk = true
	}
	// 大面积纯色（Block 小方差高占比）
	if blkRatio >= a.BlockVarMaxRatio {
		junk = true
	}
	// 低色彩度（配合低结构或低对比）
	if colorf <= 7.0 && (edges <= 0.006 || contrast <= 18.0) {
		junk = true
	}
	if a.initialized && hdist < a.MinChange/2 {
		junk = true
	}
	if a.initialized && hdist < a.MinChangeStrong {
		junk = true
	}
	if score < 55 {
		junk = true
	}

	// 更新状态
	a.lastHash = hash
	a.initialized = true

	return Result{
		Score:      score,
		IsJunk:     junk,
		Reasons:    reasons,
		Hash:       hash,
		Sharpness:  sharp,
		Contrast:   contrast,
		Uniformity: uniform,
		Mean:       mean,
		Flatness:   flat,
		EdgeDensity: edges,
		ChangeHamming: hdist,
		BlockVarRatio: blkRatio,
		Colorfulness:  colorf,
	}, nil
}

// ResetState 清空内部状态（上一帧哈希）
func (a *Analyzer) ResetState() {
	a.initialized = false
	a.lastHash = 0
}
