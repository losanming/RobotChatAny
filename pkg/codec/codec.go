package codec

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/golang/snappy"
	"github.com/klauspost/compress/zstd"
)

// Algorithm 压缩算法类型
// 可选: none/gzip/zstd/snappy
// 默认推荐 zstd：压缩比与速度兼顾，纯Go实现

type Algorithm string

const (
	AlgoNone   Algorithm = "none"
	AlgoGzip   Algorithm = "gzip"
	AlgoZstd   Algorithm = "zstd"
	AlgoSnappy Algorithm = "snappy"
)

// CompressToBase64 将原始数据压缩后再进行Base64编码
// 适合将图片等二进制通过文本通道传输，显著降低体积
func CompressToBase64(data []byte, algo Algorithm) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	var compressed []byte
	var err error
	switch algo {
	case AlgoNone:
		compressed = data
	case AlgoGzip:
		compressed, err = gzipCompress(data)
	case AlgoZstd:
		compressed, err = zstdCompress(data)
	case AlgoSnappy:
		compressed = snappy.Encode(nil, data)
	default:
		return "", fmt.Errorf("未知算法: %s", algo)
	}
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(compressed), nil
}

// DecompressFromBase64 对Base64数据先解码再解压，得到原始数据
func DecompressFromBase64(b64 string, algo Algorithm) ([]byte, error) {
	if b64 == "" {
		return nil, nil
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %w", err)
	}
	switch algo {
	case AlgoNone:
		return raw, nil
	case AlgoGzip:
		return gzipDecompress(raw)
	case AlgoZstd:
		return zstdDecompress(raw)
	case AlgoSnappy:
		return snappy.Decode(nil, raw)
	default:
		return nil, fmt.Errorf("未知算法: %s", algo)
	}
}

// ChunkBase64 将较长的Base64字符串按 maxLen 分片，便于分包传输
func ChunkBase64(s string, maxLen int) []string {
	if maxLen <= 0 || len(s) <= maxLen {
		return []string{s}
	}
	chunks := make([]string, 0, (len(s)+maxLen-1)/maxLen)
	for i := 0; i < len(s); i += maxLen {
		end := i + maxLen
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

// JoinChunks 合并Base64分片
func JoinChunks(chunks []string) string {
	if len(chunks) == 0 {
		return ""
	}
	// 预估容量，减少拷贝
	var total int
	for _, c := range chunks {
		total += len(c)
	}
	b := make([]byte, 0, total)
	for _, c := range chunks {
		b = append(b, c...)
	}
	return string(b)
}

// gzipCompress 使用gzip压缩
func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	if _, err = zw.Write(data); err != nil {
		_ = zw.Close()
		return nil, err
	}
	if err = zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// gzipDecompress 使用gzip解压
func gzipDecompress(data []byte) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	var out bytes.Buffer
	if _, err = io.Copy(&out, zr); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// zstdCompress 使用klauspost纯Go zstd
func zstdCompress(data []byte) ([]byte, error) {
	// 速度优先等级: zstd.SpeedDefault 或 SpeedFastest
	enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedDefault))
	if err != nil {
		return nil, err
	}
	defer enc.Close()
	return enc.EncodeAll(data, make([]byte, 0, len(data)/3)), nil
}

// zstdDecompress 解压zstd
func zstdDecompress(data []byte) ([]byte, error) {
	dec, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}
	defer dec.Close()
	out, err := dec.DecodeAll(data, nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompressToBase64Adaptive 自适应压缩：评估多种算法，选择体积最小者
// 返回选中的算法与Base64字符串
func CompressToBase64Adaptive(data []byte) (Algorithm, string, error) {
	if len(data) == 0 {
		return AlgoNone, "", nil
	}
	candidates := []Algorithm{AlgoZstd, AlgoGzip, AlgoSnappy, AlgoNone}
	var (
		bestAlgo Algorithm = AlgoNone
		bestEnc  string
		bestSize int = 1<<31 - 1
	)
	for _, a := range candidates {
		enc, err := CompressToBase64(data, a)
		if err != nil {
			continue
		}
		if l := len(enc); l < bestSize {
			bestSize = l
			bestEnc = enc
			bestAlgo = a
		}
	}
	if bestEnc == "" { // 兜底
		s := base64.StdEncoding.EncodeToString(data)
		return AlgoNone, s, nil
	}
	return bestAlgo, bestEnc, nil
}

// EncodeWithAlgoHeader 便捷打包: "algo:base64" 形式，便于透传与解码
func EncodeWithAlgoHeader(b64 string, algo Algorithm) string {
	return string(algo) + ":" + b64
}

// DecodeWithAlgoHeader 解析 "algo:base64" 格式
func DecodeWithAlgoHeader(s string) (Algorithm, string) {
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			return Algorithm(s[:i]), s[i+1:]
		}
	}
	return AlgoNone, s
}

// CompressSmartToBase64 仅在“有意义”节省时才压缩
// - minSize: 小于此字节数直接不压缩（如 < 1KB）
// - minGainPct: 最小节省百分比(0~100)，如设为 5 表示至少节省5%才采用压缩
// 返回：选用算法与Base64编码
func CompressSmartToBase64(data []byte, minSize int, minGainPct float64) (Algorithm, string, error) {
	// 边界处理
	if len(data) == 0 {
		return AlgoNone, "", nil
	}
	if len(data) < minSize {
		return AlgoNone, base64.StdEncoding.EncodeToString(data), nil
	}
	// 原始base64长度
	rawB64 := base64.StdEncoding.EncodeToString(data)
	rawLen := len(rawB64)

	// 选择体积最小者
	algo, enc, err := CompressToBase64Adaptive(data)
	if err != nil {
		return AlgoNone, "", err
	}
	if algo == AlgoNone {
		return AlgoNone, rawB64, nil
	}

	// 计算节省百分比
	saved := float64(rawLen-len(enc)) / float64(rawLen) * 100.0
	if saved < minGainPct {
		// 节省不达标，回退为不压缩
		return AlgoNone, rawB64, nil
	}
	return algo, enc, nil
}
