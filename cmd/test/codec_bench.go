package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"new-robot-self/pkg/codec"
)

func humanize(n int) string {
	suffixes := []string{"B", "KB", "MB", "GB"}
	s := 0
	val := float64(n)
	for val >= 1024 && s < len(suffixes)-1 {
		val /= 1024
		s++
	}
	return fmt.Sprintf("%.2f %s", val, suffixes[s])
}

// 计算百分比
// de为现在 base为原先
func pct(delta, base int) string {
	if base <= 0 {
		return "-"
	}
	return fmt.Sprintf("%.1f%%", (1-(float64(delta)/float64(base)))*100)
}

func main() {
	fileByte, err := ioutil.ReadFile("asset/everasm-system.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法读取文件: %v\n", err)
		os.Exit(1)
	}
	// 原始Base64
	start := time.Now()
	rawB64 := base64.StdEncoding.EncodeToString(fileByte)
	fmt.Printf("原始Base64: 长度=%s (%d)  耗时=%s\n", humanize(len(rawB64)), len(rawB64), time.Since(start))

	// zstd
	start = time.Now()
	toBase64, s, err := codec.CompressSmartToBase64(fileByte, 1024, 5.0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "zstd压缩失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("压缩方式: %s\n", toBase64)
	fmt.Printf("压缩: 长度=%s (%d)  耗时=%s  压缩比=%s\n", humanize(len(s)), len(s), time.Since(start), pct(len(s), len(rawB64)))
	fromBase64, err := codec.DecompressFromBase64(s, toBase64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "解压失败: %v\n", err)
		os.Exit(1)
	}
	if !bytes.Equal(fromBase64, fileByte) {
		fmt.Fprintf(os.Stderr, "解压结果与原始数据不一致\n")
		os.Exit(1)
	}
	// 解压后落地成文件
	err = ioutil.WriteFile("asset/everasm-system_zstd.sql", fromBase64, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法写入文件: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("over!")
}
