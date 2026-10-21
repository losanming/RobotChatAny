package resource

import "strings"

// 关键词到类型的映射


var categoryKeywords = map[ResourceType][]string{
	Game: {
		"游戏", "手游", "端游", "网游", "主机", "steam", "xbox", "ps5", "switch",
		"lol", "英雄联盟", "王者荣耀", "原神", "崩坏", "战双", "mmorpg", "fps",
	},
	Anime: {
		"动漫", "番剧", "漫画", "动画", "二次元", "cos", "cosplay",
		"漫展", "acg", "轻小说", "日漫", "国漫",
	},
	Social: {
		"社会", "新闻", "时事", "政治", "经济", "法律", "事件",
		"热点", "政策", "民生", "调查", "曝光",
	},
	Tech: {
		"科技", "技术", "编程", "开发", "软件", "硬件", "数码",
		"ai", "人工智能", "区块链", "云计算", "大数据", "物联网",
		"编程语言", "框架", "工具", "开源", "github",
	},
	Life: {
		"生活", "美食", "旅游", "健康", "运动", "穿搭", "护肤",
		"美妆", "家居", "装修", "养生", "美容", "减肥", "健身",
	},
	Work: {
		"工作", "职场", "求职", "面试", "简历", "薪资", "办公",
		"职业", "创业", "管理", "团队", "晋升", "加班",
	},
	Movie: {
		"电影", "影视", "院线", "netflix", "剧集", "综艺",
		"导演", "演员", "票房", "影评", "电视剧",
	},
	Music: {
		"音乐", "歌曲", "专辑", "演唱会", "歌手", "乐队",
		"说唱", "嘻哈", "rap", "流行", "摇滚", "古典",
	},
	Shopping: {
		"购物", "优惠", "折扣", "促销", "淘宝", "京东",
		"拼多多", "双十一", "618", "海淘", "好物",
	},
}

// classifyByKeywords 根据关键词对资源进行分类
func classifyByKeywords(title, content string) []ResourceType {
	// 将标题和内容转换为小写，便于匹配
	text := strings.ToLower(title + " " + content)
	
	// 记录匹配到的类型
	var types []ResourceType
	var maxScore float64
	typeScores := make(map[ResourceType]float64)

	// 对每个类型进行评分
	for typ, keywords := range categoryKeywords {
		var score float64
		for _, keyword := range keywords {
			// 计算关键词在文本中出现的次数
			count := strings.Count(text, strings.ToLower(keyword))
			if count > 0 {
				// 标题中的关键词权重更高
				titleCount := strings.Count(strings.ToLower(title), strings.ToLower(keyword))
				score += float64(titleCount) * 2.0
				score += float64(count-titleCount) * 1.0
			}
		}
		if score > 0 {
			typeScores[typ] = score
			if score > maxScore {
				maxScore = score
			}
		}
	}

	// 选择得分最高的类型，以及得分达到最高分70%的其他类型
	threshold := maxScore * 0.7
	for typ, score := range typeScores {
		if score >= threshold {
			types = append(types, typ)
		}
	}

	// 如果没有匹配到任何类型，返回原始类型
	if len(types) == 0 {
		types = append(types, Social) // 默认为社会类型
	}

	return types
}
