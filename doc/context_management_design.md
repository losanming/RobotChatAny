# 上下文管理与Token优化设计

## 一、上下文分级存储

### 1. 存储结构
```go
type ChatContext struct {
    // 短期记忆（最近的对话）
    ShortTermMemory []Message `json:"short_term"`
    
    // 中期记忆（重要对话摘要）
    MediumTermMemory []Summary `json:"medium_term"`
    
    // 长期记忆（关键信息）
    LongTermMemory []KeyInfo `json:"long_term"`
}

type Message struct {
    Content   string    `json:"content"`
    Time      int64    `json:"time"`      // 时间戳
    Priority  int      `json:"priority"`  // 优先级
    Type      string   `json:"type"`      // 消息类型
    Keywords  []string `json:"keywords"`  // 关键词
}
```

### 2. 压缩策略
- 定期总结对话要点
- 只保留关键信息
- 删除重复内容
- 合并相似主题

## 二、Token优化策略

### 1. 消息压缩
```go
type CompressedMessage struct {
    KeyPoints []string  // 关键信息点
    Sentiment float64   // 情感值
    Time      int64     // 时间戳
}
```

### 2. Token配额控制
```go
const (
    MaxTokensForRecent   = 1000 // 最近对话
    MaxTokensForSummary  = 500  // 摘要信息
    MaxTokensForUserInfo = 300  // 用户信息
    MaxTokensForPrompt   = 200  // 基础prompt
)
```

### 3. 场景区分
```go
// 不同场景的token限制
var SceneTokenLimits = map[string]int{
    "日常闲聊": 800,   // 较少token
    "情感倾诉": 1500,  // 需要更多情感上下文
    "深度话题": 2000,  // 需要更多知识背景
    "默认": 1000,
}
```

## 三、存储优化

### 1. 分级缓存
- Redis缓存热数据
- MongoDB存储历史数据
- 定期归档冷数据

### 2. 数据结构
```go
type ContextCache struct {
    // Redis缓存
    HotContext struct {
        RecentMessages []Message
        UserState     UserState
        ActiveTopics  []string
    }
    
    // MongoDB存储
    HistoricalContext struct {
        Summaries    []Summary
        KeyMemories  []Memory
        Statistics   UserStats
    }
}
```

## 四、智能提取机制

### 1. 上下文提取
```go
func extractContext(userId string, currentTopic string) string {
    // 1. 获取最近对话
    recentChat := getRecentChat(userId, 5)
    
    // 2. 获取相关的重要记忆
    relevantMemories := getRelevantMemories(userId, currentTopic)
    
    // 3. 获取用户关键信息
    userInfo := getUserKeyInfo(userId)
    
    // 4. 组装上下文
    return buildContext(recentChat, relevantMemories, userInfo)
}
```

### 2. 动态裁剪规则
- 优先保留最近的对话
- 保留高情感值的内容
- 保留关键决策点
- 删除重复信息

## 五、成本控制

### 1. 存储限制
- 每个用户的最大存储容量
- 历史记录保存时间
- 冷数据清理策略

### 2. Token使用优化
- 智能摘要生成
- 关键信息提取
- 冗余信息删除
- 上下文相关性过滤

## 六、实现步骤

1. **初始化阶段**
   - 设置存储限制
   - 初始化缓存系统
   - 配置Token限制

2. **运行时处理**
   - 实时消息压缩
   - 动态Token控制
   - 上下文智能提取

3. **维护阶段**
   - 定期数据清理
   - 性能监控
   - 优化策略调整

## 七、注意事项

1. 保证信息提取的准确性
2. 避免过度压缩导致上下文丢失
3. 定期监控Token使用情况
4. 注意存储成本控制
5. 保持系统响应时间
6. 考虑数据备份机制

## 八、监控指标

1. **性能指标**
   - Token使用量
   - 存储空间使用
   - 响应时间
   - 压缩比率

2. **质量指标**
   - 对话连贯性
   - 上下文准确性
   - 用户满意度
   - 系统稳定性
