package transform

import "math"

// InChina 判断坐标点是否处于中国大陆范围（粗略BBox判断）。
// 说明：该判断仅用于是否需要进行 WGS84→GCJ02 的转换；并非精确国界判定。
// 参数：lat（纬度）、lon（经度）
// 返回：true 表示在中国大陆范围内，false 表示不在（境外或边界外）
func InChina(lat, lon float64) bool {
	// Rough bbox: lat 0.8293~55.8271, lon 72.004~137.8347
	if lon < 72.004 || lon > 137.8347 || lat < 0.8293 || lat > 55.8271 {
		return false
	}
	return true
}

// WGS84ToGCJ02 将 WGS-84 坐标转换为 GCJ-02（火星坐标）。
// 行为：若坐标不在中国大陆范围内（InChina=false），则直接返回原始坐标（不做转换）。
// 注意：算法为通用开源实现，适用于常见场景；严肃测绘请遵循相关法律法规。
func WGS84ToGCJ02(lat, lon float64) (float64, float64) {
	if !InChina(lat, lon) {
		return lat, lon
	}
	return transformLatLon(lat, lon)
}

// 以下实现基于社区常见开源公式改写，适配 Go 语言环境。
// 只做经纬度偏移计算，不涉及更高阶纠偏（如地形、时间等）。
const (
	a = 6378245.0
	ees = 0.00669342162296594323
)

func transformLatLon(wgLat, wgLon float64) (mgLat, mgLon float64) {
	dLat := transformLat(wgLon-105.0, wgLat-35.0)
	dLon := transformLon(wgLon-105.0, wgLat-35.0)
	radLat := wgLat / 180.0 * math.Pi
	magic := math.Sin(radLat)
	magic = 1 - ees*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ees)) / (magic * sqrtMagic) * math.Pi)
	dLon = (dLon * 180.0) / (a / sqrtMagic * math.Cos(radLat) * math.Pi)
	mgLat = wgLat + dLat
	mgLon = wgLon + dLon
	return
}

func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*math.Pi) + 40.0*math.Sin(y/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*math.Pi) + 320*math.Sin(y*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

func transformLon(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*math.Pi) + 40.0*math.Sin(x/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*math.Pi) + 300.0*math.Sin(x/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}
