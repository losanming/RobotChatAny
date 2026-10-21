package quality

import (
	"image"
	"image/color"
	"math"
	"math/bits"
)

// clamp 将 x 裁剪到 [min, max]
func clamp(x, min, max float64) float64 {
    if x < min {
        return min
    }
    if x > max {
        return max
    }
    return x
}

// laplacianVariance 计算拉普拉斯卷积响应的方差，衡量清晰度（越大越清晰）
func laplacianVariance(g *image.Gray) float64 {
    kernel := [3][3]int{
        {0, -1, 0},
        {-1, 4, -1},
        {0, -1, 0},
    }
    b := g.Bounds()
    if b.Dx() < 3 || b.Dy() < 3 {
        return 0
    }
    var sum float64
    var sumSq float64
    n := 0
    for y := b.Min.Y + 1; y < b.Max.Y-1; y++ {
        for x := b.Min.X + 1; x < b.Max.X-1; x++ {
            var conv int
            for ky := -1; ky <= 1; ky++ {
                for kx := -1; kx <= 1; kx++ {
                    p := int(g.GrayAt(x+kx, y+ky).Y)
                    conv += kernel[ky+1][kx+1] * p
                }
            }
            v := float64(conv)
            sum += v
            sumSq += v * v
            n++
        }
    }
    if n == 0 {
        return 0
    }
    mean := sum / float64(n)
    variance := sumSq/float64(n) - mean*mean
    if variance < 0 {
        return 0
    }
    return variance
}

// luminanceStdDev 计算亮度标准差（对比度，越大对比度越高）
func luminanceStdDev(g *image.Gray) float64 {
    b := g.Bounds()
    var sum float64
    var sumSq float64
    n := 0
    for y := b.Min.Y; y < b.Max.Y; y++ {
        for x := b.Min.X; x < b.Max.X; x++ {
            v := float64(g.GrayAt(x, y).Y)
            sum += v
            sumSq += v * v
            n++
        }
    }
    if n == 0 {
        return 0
    }
    mean := sum / float64(n)
    variance := sumSq/float64(n) - mean*mean
    if variance < 0 {
        variance = 0
    }
    return math.Sqrt(variance)
}

// uniformity 使用熵的反向指标来度量均匀度：1 - H/8，取值[0,1]，越大越均匀
func uniformity(g *image.Gray) float64 {
    var hist [256]int
    b := g.Bounds()
    n := 0
    for y := b.Min.Y; y < b.Max.Y; y++ {
        for x := b.Min.X; x < b.Max.X; x++ {
            hist[g.GrayAt(x, y).Y]++
            n++
        }
    }
    if n == 0 {
        return 1
    }
    var entropy float64
    invN := 1.0 / float64(n)
    for i := 0; i < 256; i++ {
        if hist[i] == 0 {
            continue
        }
        p := float64(hist[i]) * invN
        entropy -= p * math.Log2(p)
    }
    if entropy < 0 {
        entropy = 0
    }
    if entropy > 8 {
        entropy = 8
    }
    return 1.0 - entropy/8.0
}

// averageHash64 计算 8x8 aHash，返回 64-bit 哈希
func averageHash64(g *image.Gray) uint64 {
    // 采样为 8x8（简单均匀下采样）
    b := g.Bounds()
    w := b.Dx()
    h := b.Dy()
    if w == 0 || h == 0 {
        return 0
    }
    var block [64]float64
    idx := 0
    for gy := 0; gy < 8; gy++ {
        y0 := b.Min.Y + gy*h/8
        y1 := b.Min.Y + (gy+1)*h/8
        if gy == 7 {
            y1 = b.Max.Y
        }
        for gx := 0; gx < 8; gx++ {
            x0 := b.Min.X + gx*w/8
            x1 := b.Min.X + (gx+1)*w/8
            if gx == 7 {
                x1 = b.Max.X
            }
            var sum float64
            n := 0
            for y := y0; y < y1; y++ {
                for x := x0; x < x1; x++ {
                    sum += float64(g.GrayAt(x, y).Y)
                    n++
                }
            }
            if n == 0 {
                block[idx] = 0
            } else {
                block[idx] = sum / float64(n)
            }
            idx++
        }
    }
    // 计算均值
    var total float64
    for i := 0; i < 64; i++ {
        total += block[i]
    }
    avg := total / 64.0
    // 生成位图
    var hash uint64
    for i := 0; i < 64; i++ {
        if block[i] >= avg {
            hash |= 1 << uint(i)
        }
    }
    return hash
}

// hammingDistance64 计算 64-bit 哈希的海明距离
func hammingDistance64(a, b uint64) int {
    return bits.OnesCount64(a ^ b)
}

// toGray 将任意图像转换为灰度图
func toGray(src image.Image) *image.Gray {
	bounds := src.Bounds()
	g := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, gch, b, _ := src.At(x, y).RGBA()
			// 使用加权灰度（ITU-R BT.601）
			l := 0.299*float64(r) + 0.587*float64(gch) + 0.114*float64(b)
			g.SetGray(x, y, color.Gray{Y: uint8((l / 257.0) + 0.5)})
		}
	}
	return g
}

// blockVarianceRatio 将图像划分为 grid x grid 小块，统计亮度方差小于 varThr 的“小方差块”占比
// 该指标用于检测大面积单色/灰屏：小方差块比例高意味着画面平坦无细节
func blockVarianceRatio(g *image.Gray, grid int, varThr float64) float64 {
	if grid <= 1 {
		return 1
	}
	b := g.Bounds()
	w := b.Dx()
	h := b.Dy()
	if w < grid || h < grid {
		return 1
	}
	blockW := w / grid
	blockH := h / grid
	total := 0
	lowVar := 0
	for gy := 0; gy < grid; gy++ {
		for gx := 0; gx < grid; gx++ {
			x0 := b.Min.X + gx*blockW
			y0 := b.Min.Y + gy*blockH
			x1 := x0 + blockW
			y1 := y0 + blockH
			if gx == grid-1 {
				x1 = b.Max.X
			}
			if gy == grid-1 {
				y1 = b.Max.Y
			}
			// 块内统计
			var sum float64
			var sumSq float64
			n := 0
			for y := y0; y < y1; y++ {
				for x := x0; x < x1; x++ {
					v := float64(g.GrayAt(x, y).Y)
					sum += v
					sumSq += v * v
					n++
				}
			}
			if n == 0 {
				continue
			}
			mean := sum / float64(n)
			variance := sumSq/float64(n) - mean*mean
			if variance < 0 {
				variance = 0
			}
			if variance <= varThr {
				lowVar++
			}
			total++
		}
	}
	if total == 0 {
		return 1
	}
	return float64(lowVar) / float64(total)
}

// colorfulnessRGB 计算颜色丰富度（Hasler-Suesstrunk 指标简化版）
// 返回值越大表示色彩越丰富，纯灰/单色趋近于 0-10 的小值
func colorfulnessRGB(img image.Image) float64 {
    b := img.Bounds()
    if b.Empty() {
        return 0
    }
    var sumRG, sumYB float64
    var sumRG2, sumYB2 float64
    n := 0
    for y := b.Min.Y; y < b.Max.Y; y++ {
        for x := b.Min.X; x < b.Max.X; x++ {
            r16, g16, b16, _ := img.At(x, y).RGBA()
            // 转 0-255
            r := float64(r16) / 257.0
            gch := float64(g16) / 257.0
            bl := float64(b16) / 257.0
            rg := r - gch
            yb := 0.5*(r+gch) - bl
            sumRG += rg
            sumYB += yb
            sumRG2 += rg * rg
            sumYB2 += yb * yb
            n++
        }
    }
    if n == 0 {
        return 0
    }
    meanRG := sumRG / float64(n)
    meanYB := sumYB / float64(n)
    varRG := sumRG2/float64(n) - meanRG*meanRG
    varYB := sumYB2/float64(n) - meanYB*meanYB
    if varRG < 0 {
        varRG = 0
    }
    if varYB < 0 {
        varYB = 0
    }
    stdRG := math.Sqrt(varRG)
    stdYB := math.Sqrt(varYB)
    // Hasler-Suesstrunk: C = sqrt(std_rg^2 + std_yb^2) + 0.3*sqrt(mean_rg^2 + mean_yb^2)
    color := math.Sqrt(stdRG*stdRG+stdYB*stdYB) + 0.3*math.Sqrt(meanRG*meanRG+meanYB*meanYB)
    return color
}

// normalize 将 value 在 [min, max] 映射到 [0,1]，超界会裁剪
func normalize(value, min, max float64) float64 {
    if max <= min {
        return 0
    }
    return clamp((value-min)/(max-min), 0, 1)
}

// luminanceMean 计算亮度均值（0-255）
func luminanceMean(g *image.Gray) float64 {
	bounds := g.Bounds()
	var sum float64
	count := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			sum += float64(g.GrayAt(x, y).Y)
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// flatnessRatio 统计接近均值的像素占比。delta 为容忍范围（0-255）
func flatnessRatio(g *image.Gray, delta uint8) float64 {
	mean := luminanceMean(g)
	lo := mean - float64(delta)
	hi := mean + float64(delta)
	if lo < 0 {
		lo = 0
	}
	if hi > 255 {
		hi = 255
	}
	bounds := g.Bounds()
	total := 0
	flat := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			v := float64(g.GrayAt(x, y).Y)
			if v >= lo && v <= hi {
				flat++
			}
			total++
		}
	}
	if total == 0 {
		return 1
	}
	return float64(flat) / float64(total)
}

// edgeDensity 统计强边缘像素占比，基于拉普拉斯响应绝对值大于阈值
func edgeDensity(g *image.Gray, threshold int) float64 {
	kernel := [3][3]int{
		{0, -1, 0},
		{-1, 4, -1},
		{0, -1, 0},
	}
	bounds := g.Bounds()
	if bounds.Dx() < 3 || bounds.Dy() < 3 {
		return 0
	}
	strong := 0
	total := 0
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			var conv int
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					p := int(g.GrayAt(x+kx, y+ky).Y)
					conv += kernel[ky+1][kx+1] * p
				}
			}
			if conv < 0 {
				conv = -conv
			}
			if conv > threshold {
				strong++
			}
			total++
		}
	}
	if total == 0 {
		return 0
	}
	return float64(strong) / float64(total)
}
