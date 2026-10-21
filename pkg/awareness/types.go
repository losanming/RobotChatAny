package awareness

import (
	loc "new-robot-self/pkg/location"
)

// Input 输入：一次环境感知请求
// - RawLoc: 端侧上报的原始 WGS-84 坐标与元信息（时间戳必须为毫秒）
// - ImageJPEG(optional): 同步提供的一帧桌面/相机 JPEG 图像（可为空）
// - NeedAround: 是否需要拉取更丰富的周边 POI（除逆地理内置 POI 外）
// - AroundKeywords/Types/Radius: 周边搜索的可选参数
// 说明：本模块不做截图采集，仅做整合。

type Input struct {
	RawLoc        loc.RawLocation
	ImageJPEG     []byte
	NeedAround    bool
	AroundKeywords string
	AroundTypes    string
	AroundRadius   int // meters
}

// ImageInsight 图片识别的输出（占位：后续接入本地图像理解或云端视觉）

type ImageInsight struct {
	Labels   []string          // Top-n 语义标签
	Objects  []string          // 粗粒度物体类型
	Summary  string            // 一句话总结
	Extra    map[string]string // 额外元信息（可选）
}

// Context 汇总后的环境上下文（用于喂给大模型）

type Context struct {
	// 时间与设备
	TS       int64  // ms
	DeviceID string
	Source   string

	// 地理语义
	LocTrans        loc.TransformedLocation
	Regeo           *loc.ReverseGeocodeResult
	Around          *loc.PlaceAroundResult
	Weather         *loc.WeatherResult

	// 视觉语义（可空）
	Vision *ImageInsight
}

// Result 汇总返回

type Result struct {
	Ctx    Context
	Prompt string
}
