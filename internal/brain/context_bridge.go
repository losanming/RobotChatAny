package brain

import (
	"context"
	"sort"
	"time"

	ctxmgr "new-robot-self/internal/context"
)

// DialogSource 是 brain 从外部上下文管理模块同步对话所需的最小接口
// 适配 internal/context.ContextManager（其已实现 GetContext(userID)）
// 之所以定义接口，是为了方便单元测试与未来替换实现。
type DialogSource interface {
	GetContext(userID string) (*ctxmgr.ChatContext, error)
}

// SyncDialogFromSource 拉取用户的对话上下文（短期记忆为主），注入到 WorldState.Dialog
// - maxRecent: 仅同步最近 N 条短期消息，避免冗长；传 <=0 则默认为 10 条
// - role 映射：优先使用 ctx.Message.Type 作为 role；为空则默认为 "user"
func (b *Brain) SyncDialogFromSource(ctx context.Context, userID string, src DialogSource, maxRecent int) error {
	if src == nil {
		return nil
	}
	cc, err := src.GetContext(userID)
	if err != nil || cc == nil {
		return err
	}
	if b.state.Dialog == nil {
		b.state.Dialog = &DialogContext{}
	}
	// 复制并截取最近 N 条
	recent := make([]ctxmgr.Message, len(cc.ShortTermMemory))
	copy(recent, cc.ShortTermMemory)
	sort.Slice(recent, func(i, j int) bool { return recent[i].Time < recent[j].Time })
	if maxRecent <= 0 {
		maxRecent = 10
	}
	if len(recent) > maxRecent {
		recent = recent[len(recent)-maxRecent:]
	}
	// 注入到 Dialog.History
	for _, m := range recent {
		role := m.Type
		if role == "" {
			role = "user"
		}
		b.state.Dialog.History = append(b.state.Dialog.History, DialogTurn{
			TS:      m.Time, // 统一时间戳秒
			Role:    role,
			Content: m.Content,
		})
	}
	b.state.UpdatedTS = time.Now().Unix()
	return nil
}
