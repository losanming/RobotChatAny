# 用户个性化老哥设计方案

## 一、设计目标

通过用户问卷和交互数据，为每个用户定制独特的"老哥"性格特征，实现个性化的AI对话体验。

## 二、用户问卷设计

### 1. 基础性格定制
```
1. 你希望老哥是什么性格？
   □ 社恐老实人 - "我超，有点紧张..."
   □ 话痨段子手 - "笑死，这事我可熟..."
   □ 佛系养生派 - "摆烂了都，佛了佛了..."
   □ 热血动漫迷 - "燃起来了！这波必须..."
   □ 资深吐槽王 - "属实难绷，这也太离谱..."
   □ 文艺青年派 - "细品了一下，感觉..."

2. 老哥的兴趣爱好倾向？（多选）
   □ 二次元 □ 游戏 □ 美食 □ 运动 
   □ 数码 □ 电影 □ 音乐 □ 旅游

3. 老哥的工作设定？
   □ 程序员 □ 设计师 □ 运营 
   □ 产品 □ 市场 □ 其他
```

### 2. 情感交互倾向
```
1. 你平时最需要倾诉的话题是？
   □ 工作压力 □ 人际关系 
   □ 生活琐事 □ 个人理想

2. 你希望老哥如何回应？
   □ 认真倾听和共情
   □ 幽默化解和调节
   □ 给出建议和方案
   □ 纯发泄不需要回应

3. 你喜欢什么样的交流方式？
   □ 简短有趣 □ 深度交流
   □ 段子调侃 □ 感性共鸣
```

### 3. 互动习惯设置
```
1. 希望老哥多久找你聊天？
   □ 经常互动 □ 适度交流 
   □ 重要时刻 □ 等我主动

2. 最喜欢的网络用语风格？
   □ 二次元风 □ 网络热梗
   □ 文艺风 □ 简单直接

3. 特别关注的话题领域？
   （可自由填写）
```

## 三、数据存储结构

```go
type UserPersonality struct {
    // 基础信息
    UserId        string    `json:"user_id"`
    CreateTime    time.Time `json:"create_time"`
    UpdateTime    time.Time `json:"update_time"`

    // 性格设定
    PersonalityType string   `json:"personality_type"` // 性格类型
    Interests       []string `json:"interests"`        // 兴趣爱好
    Occupation      string   `json:"occupation"`       // 工作设定

    // 情感倾向
    EmotionalTopics []string `json:"emotional_topics"` // 倾诉话题
    ResponseStyle   string   `json:"response_style"`   // 回应方式
    ChatStyle       string   `json:"chat_style"`       // 交流方式

    // 互动设置
    ChatFrequency   string   `json:"chat_frequency"`   // 聊天频率
    LanguageStyle   string   `json:"language_style"`   // 用语风格
    FocusTopics     []string `json:"focus_topics"`     // 关注话题

    // 动态数据
    InteractionHistory []Interaction `json:"interaction_history"` // 互动历史
    EmotionalState    float64       `json:"emotional_state"`     // 当前情绪值
    TopicPreferences  map[string]float64 `json:"topic_preferences"` // 话题偏好权重
}
```

## 四、Prompt生成逻辑

### 1. 基础模板
```
你是一个{personality_type}性格的年轻人，热爱{interests}。
作为一个{occupation}，你经常会{daily_activities}。
你说话的风格是{language_style}，经常使用的表达方式包括{expressions}。
```

### 2. 动态组装规则
1. 基础人设（固定）
2. 个性化性格（来自问卷）
3. 用户兴趣（来自问卷+交互学习）
4. 历史互动记忆（动态累积）
5. 当前话题（实时）
6. 情境因素（实时）

### 3. 示例Prompt
```
你是一个社恐但内心热血的年轻人，特别喜欢二次元和游戏。
作为一个程序员，你经常熬夜写代码，对工作压力深有体会。
你说话喜欢用二次元风格，但在倾听别人烦恼时会认真共情。
现在是周一早上，你看到了{hot_topics}，想到自己也是...
```

## 五、交互记忆系统

### 1. 记忆类型
- 重要对话内容
- 用户偏好话题
- 特殊互动经历
- 共同的梗和笑话

### 2. 记忆更新机制
- 定期保存重要对话
- 更新话题偏好权重
- 记录情感波动节点
- 生成个性化梗库

## 六、实现步骤

1. **用户注册阶段**
   - 完成性格问卷
   - 初始化个性化数据
   - 生成基础Prompt模板

2. **日常交互阶段**
   - 记录互动数据
   - 更新用户画像
   - 动态调整Prompt

3. **定期优化阶段**
   - 分析互动效果
   - 优化个性化参数
   - 更新交互策略

## 七、注意事项

1. 保持性格特征的一致性
2. 避免过度个性化导致的不自然
3. 注意隐私数据的保护
4. 定期清理无用的记忆数据
5. 保持角色设定的基本底线
