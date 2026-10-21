package push

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"new-robot-self/internal/ai"
	"new-robot-self/internal/connection"
	"new-robot-self/internal/prompt"
	"new-robot-self/internal/protocol"
	"new-robot-self/internal/queue"
	"new-robot-self/internal/resource"
	"new-robot-self/internal/trigger"
	"sync"
	"time"
)

// Service 推送服务
type Service struct {
	// 连接管理器
	connMgr *connection.ConnectionManager
	// 触发器管理器
	triggerMgr *trigger.Manager
	// 资源获取器
	fetcher *resource.Fetcher
	// AI服务
	aiService *ai.Service
	// Prompt生成器
	promptGen *prompt.ContextPromptGenerator
	// 用户最后收到推送的时间
	lastPushTime sync.Map // key: userId, value: int64(timestamp)
	// 用户推送间隔(分钟)
	pushIntervals sync.Map // key: userId, value: int
}

// NewService 创建推送服务
func NewService(
	connMgr *connection.ConnectionManager,
	triggerMgr *trigger.Manager,
	fetcher *resource.Fetcher,
	aiService *ai.Service,
	promptGen *prompt.ContextPromptGenerator,
) *Service {
	return &Service{
		connMgr:    connMgr,
		triggerMgr: triggerMgr,
		fetcher:    fetcher,
		aiService:  aiService,
		promptGen:  promptGen,
	}
}

// Start 启动推送服务
func (s *Service) Start() {
	// 启动定时推送检查
	go s.checkTimePointTriggers()
	// 启动随机推送检查
	go s.checkRandomPush()
}

// checkTimePointTriggers 检查定时触发器
func (s *Service) checkTimePointTriggers() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// 获取当前时间点的触发器
		triggers := s.triggerMgr.GetRandomTriggers(1, []trigger.TriggerType{trigger.TimePoint})
		if len(triggers) == 0 {
			continue
		}

		// 向所有在线用户推送
		s.pushToAllOnlineUsers(triggers[0])
	}
}

// checkRandomPush 检查随机推送
func (s *Service) checkRandomPush() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 获取所有在线用户
		onlineUsers := s.connMgr.GetUserIDs()
		if len(onlineUsers) == 0 {
			continue
		}

		// 为每个用户检查是否需要推送
		for _, userId := range onlineUsers {
			if s.shouldPushToUser(userId) {
				// 获取随机触发器
				triggers := s.triggerMgr.GetRandomTriggers(1, []trigger.TriggerType{
					trigger.RealTime,
					trigger.RandomTopic,
				})
				if len(triggers) > 0 {
					// 推送给用户
					s.pushToUser(userId, triggers[0])
					// 更新用户最后推送时间
					s.updateUserPushTime(userId)
				}
			}
		}
	}
}

// shouldPushToUser 判断是否应该向用户推送
func (s *Service) shouldPushToUser(userId int64) bool {
	now := time.Now().Unix()

	// 获取用户最后推送时间
	lastTime, ok := s.lastPushTime.Load(userId)
	if !ok {
		return true
	}

	// 获取用户推送间隔
	interval, ok := s.pushIntervals.Load(userId)
	if !ok {
		// 默认间隔30-120分钟
		interval = 30 + rand.Intn(90)
		s.pushIntervals.Store(userId, interval)
	}

	// 检查是否达到推送间隔
	return now-lastTime.(int64) >= int64(interval.(int)*60)
}

// updateUserPushTime 更新用户最后推送时间
func (s *Service) updateUserPushTime(userId int64) {
	now := time.Now().Unix()
	s.lastPushTime.Store(userId, now)

	// 重新随机设置下次推送间隔
	interval := 30 + rand.Intn(90)
	s.pushIntervals.Store(userId, interval)
}

// pushToUser 向指定用户推送消息
func (s *Service) pushToUser(userId int64, tr trigger.TriggerResult) {
	//index1
	// 获取触发器类型和上下文
	var tmp_ctx = context.Background()
	var (
		resType     resource.ResourceType
		triggerDesc string
		timeNow     = time.Now()
		emotionType string // 情感状态
		topicType   string // 话题类型
	)

	// 根据触发器类型设置上下文
	switch tr.Type {
	case trigger.TimePoint:
		// 时间点触发
		rt := tr.Content.(trigger.TimePointTrigger)
		resType = resource.ResourceType(rt.Topic)
		triggerDesc = rt.Topic

		// 根据时间点设置情感和话题
		hour := timeNow.Hour()
		switch {
		case hour == 7:
			emotionType = "吐槽" // 上班时间
			topicType = "工作生活"
		case hour == 12:
			emotionType = "快乐" // 午饭时间
			topicType = "工作生活"
		case hour == 15:
			emotionType = "共情" // 下午困
			topicType = "社会现象"
		case hour == 18:
			emotionType = "愤怒" // 下班时间
			topicType = "社会现象"
		}

	case trigger.RandomTopic:
		// 随机话题触发
		rt := tr.Content.(trigger.RandomTopicTrigger)
		resType = resource.ResourceType(rt.Type)
		triggerDesc = rt.Content

		// 根据资源类型设置情感和话题
		switch rt.Type {
		case "game":
			emotionType = "快乐"
			topicType = "游戏"
		case "anime":
			emotionType = "快乐"
			topicType = "动漫"
		case "social":
			emotionType = "吐槽"
			topicType = "社会现象"
		case "work":
			emotionType = "愤怒"
			topicType = "工作生活"
		}
	}

	var tr_msg = fmt.Sprintf("聊天内容是:%s,情感状态为:%s,话题类型为:%s", triggerDesc, emotionType, topicType)

	// 从缓存中获取一个随机资源
	resources, err := s.fetcher.GetRandomResources([]resource.ResourceType{resType}, 1)
	if err != nil || len(resources) == 0 {
		return
	}
	// index2
	// 获取用户ID字符串
	userIDStr := fmt.Sprintf("%d", userId)

	// 生成个性化prompt
	prompt, err := s.promptGen.GeneratePrompt(userIDStr, tr_msg)
	if err != nil {
		log.Printf("生成prompt失败: %v", err)
		return
	}

	// 调用AI服务获取流式回复
	resp, err := s.aiService.Chat(tmp_ctx, ai.Message{
		Role:    "user",
		Content: prompt,
	})
	if err != nil {
		return
	}

	// 记录AI回答
	if err := s.promptGen.RecordAIResponse(userIDStr, resp); err != nil {
		log.Printf("记录AI回答失败: %v", err)
	}

	// 如果有有效的回复内容，发送给用户
	if resp != "" {
		dis, b := s.connMgr.GetDispatcher()
		if !b {
			log.Println("dispatcher not found")
			return
		}
		dis.AddUser(userId)
		defer dis.RemoveUser(userId)
		err = s.connMgr.SendMessage(tmp_ctx, &queue.Message{
			From:      10086001, // 系统消息
			To:        userId,
			Type:      protocol.MsgTypeText,
			Content:   []byte(resp),
			Timestamp: time.Now().Unix(),
			Status:    0,
		})
		if err != nil {
			log.Printf("send message to user %d failed: %v", userId, err)
		}
	}
}

// pushToAllOnlineUsers 向所有在线用户推送消息
func (s *Service) pushToAllOnlineUsers(tr trigger.TriggerResult) {
	onlineUsers := s.connMgr.GetUserIDs()
	for _, userId := range onlineUsers {
		s.pushToUser(userId, tr)
	}
}
