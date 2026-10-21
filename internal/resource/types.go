package resource

// ResourceType 资源类型
type ResourceType string

const (
	Game     ResourceType = "game"     // 游戏相关
	Anime    ResourceType = "anime"    // 动漫相关
	Social   ResourceType = "social"   // 社会热点
	Tech     ResourceType = "tech"     // 科技资讯
	Life     ResourceType = "life"     // 生活相关
	Work     ResourceType = "work"     // 职场相关
	Movie    ResourceType = "movie"    // 电影相关
	Music    ResourceType = "music"    // 音乐相关
	Shopping ResourceType = "shopping" // 购物相关
)

// Resource 资源结构
type Resource struct {
	Types     []ResourceType `json:"types"`     // 资源类型列表
	Title     string         `json:"title"`     // 标题
	Content   string         `json:"content"`   // 内容
	URL       string         `json:"url"`       // 链接
	HotValue  int64          `json:"hot_value"` // 热度值
	Timestamp int64          `json:"timestamp"` // 获取时间
}

// ResourceSource 资源来源配置
type ResourceSource struct {
	Type            ResourceType      `json:"type"`             // 资源类型
	Name            string            `json:"name"`             // 来源名称
	URL             string            `json:"url"`              // 来源URL
	Selector        string            `json:"selector"`         // 主选择器
	TitleSelector   string            `json:"title_selector"`   // 标题选择器
	URLSelector     string            `json:"url_selector"`     // 链接选择器
	ContentSelector string            `json:"content_selector"` // 内容选择器
	HotSelector     string            `json:"hot_selector"`     // 热度选择器
	Headers         map[string]string `json:"headers"`          // 自定义请求头
}

func GetAllType() []ResourceType {
	return []ResourceType{
		Game,
		Anime,
		Social,
		Tech,
		Life,
		Work,
		Movie,
		Music,
		Shopping,
	}
}
