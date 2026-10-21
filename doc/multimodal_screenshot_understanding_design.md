# 多模态截图理解与上下文融合设计文档

版本：v1.0  
状态：草案/可落地实施  
目录：`doc/multimodal_screenshot_understanding_design.md`

## 1. 目标与范围

- 全面解析“客户端截图 + 用户聊天历史 + 长期记忆 + 设备/日历上下文”，判断：
  - 正在做什么（Activity）
  - 计划做什么（Intent/NextStep）
  - 心情与状态（Mood/Stress）
  - 背景与上下文（Context：时间/应用/主题）
  - 风险与机会（Risk/Opportunity）
  - 建议与行动（Suggestions/Actions）
- 输出结构化结果（统一Schema），驱动主动对话与动作。
- 成本可控：本地质量筛选、批处理（Batch）降本、多帧聚合、最小化上传。
- 隐私合规：本地脱敏、最小化原则、端到端加密与可审计。

与现有模块对齐：
- 截图采集与压缩：`pkg/screenshot/screenshot.go`、`pkg/screenshot/service.go`
- 质量筛选：`pkg/screenshot/quality/analyzer.go`、`metrics.go`
- 文档参考：`doc/multi_screenshot_analysis_plan.md`

## 2. 总体架构

```
客户端（截图/质量筛选/脱敏）
   → Ingest MQ（图像引用或二进制 + 元数据）
      → 云端消费者（批处理/多模态推理/Qwen-VL）
         → Result MQ（结构化结果）
            → 客户端：TCP长连优先推送，定时拉取容灾
```

- 分区与顺序：MQ按 `user_id` 分区，保证单用户时序。
- 扩展：消费者组水平扩展；DLQ 处理异常消息。

## 3. 消息契约（JSON Schema, 建议）

时间字段均为时间戳（秒）。

### 3.1 入站（Ingest）
```json
{
  "version": "1.0",
  "msg_id": "uuid",
  "trace_id": "uuid",
  "user_id": "u123",
  "device_id": "d456",
  "ts": 1724228135,
  "image": {
    "bytes": "<binary-or-null>",
    "url": "https://obs/...&sig=...",
    "content_type": "image/jpeg",
    "width": 1280,
    "height": 720,
    "size": 122400,
    "jpeg_quality": 75,
    "hash": "ahash64",
    "quality_score": 73.2,
    "is_junk": false
  },
  "context": {
    "recent_chat_summary": "最近N轮聊天关键信息...",
    "long_term_profile_keys": ["Golang","xorm","习惯晚间开发"],
    "calendar": {"next_event_ts": 1724229000, "title": "站会"},
    "env": {"time_of_day": "afternoon", "weekday": 4}
  },
  "cfg": {
    "priority": "normal",
    "deadline_ts": 1724228335,
    "batch_hint": true
  }
}
```

### 3.2 出站（Result）
```json
{
  "version": "1.0",
  "msg_id": "uuid",
  "trace_id": "uuid",
  "ref_ingest_id": "uuid",
  "user_id": "u123",
  "ts": 1724228200,
  "observation": {
    "apps": [{"name":"VSCode","confidence":0.92}],
    "topics": ["Golang","xorm"],
    "text_snippets": ["函数签名","错误栈"],
    "mood": {"label":"focused","score":0.68},
    "environment": {"time_of_day":"afternoon","weekday":4}
  },
  "inference": {
    "doing": "在IDE阅读并修改Go服务代码",
    "intent": "准备接入并调参质量筛选模块",
    "next_2h_plan": ["联调质量筛选","编写埋点","验证性能"],
    "risks": ["模型调用成本上升","隐私泄露风险"],
    "opportunities": ["Batch降本50%","主动提醒减少中断"],
    "confidence": 0.81,
    "evidence": ["窗口标题包含.go","近期聊天提到筛选接入"]
  },
  "suggestions": [
    {"type":"reminder","text":"现在联调质量筛选并开启采样埋点？","urgency":"medium"},
    {"type":"break","text":"已连续工作50分钟，考虑3分钟走动。","urgency":"low"}
  ],
  "meta": {
    "model": "qwen-vl-plus@2025-08-16",
    "latency_ms": 620,
    "tokens_in": 1451,
    "tokens_out": 420,
    "cost_cny": 0.0012
  }
}
```

- 幂等：`msg_id` 唯一；结果用 `ref_ingest_id` 关联入站。
- 顺序：按 `user_id` 分区；客户端维持游标/去重。

## 4. 客户端侧策略

- 质量筛选（已实现，待接入服务回调）：`pkg/screenshot/quality/`。`IsJunk==true` 直接丢弃。
- 脱敏：
  - OCR 检出手机号/邮箱/姓名/账号 → 马赛克；
  - 窗口黑名单（支付/密码）不上传或仅上传摘要。
- 触发：
  - 定时（默认1min）+ 显著变化触发（aHash/Hamming）；
  - 低功耗/锁屏时暂停或降频。
- 传输：
  - 小流量阶段可直传 `[]byte`；中长期建议 OBS/S3 + 短期签名URL；
  - 全链路加密；配置化超时与重试。

## 5. 云端消费与推理

- 批处理器：
  - 按 `user_id` & 时间窗口聚合 3-5 张截图，调用 Qwen-VL Batch（半价）；
  - 帧间去重（aHash/Hamming），优先显著变化帧；
  - 合并上下文：最近聊天摘要 + 长期画像 Top-K 检索 + 设备/日历上下文。
- 推理范式：
  - ReAct/Chain-of-Thought 受控模板；
  - 强制 JSON 输出（与“出站Schema”一致）；
  - 自一致多采样 n=2~3，合并取中位/投票；
  - 角色化提示词：编码/文档/会议等场景模板（可配置）。

### 5.1 提示词模板（示意）
```
系统：你是多模态全能助手。输入包含：截图序列、最近聊天摘要、用户长期画像、设备/日历上下文。
要求：先做 Observation，再做 Inference，最后给出 Suggestions。严格按JSON Schema输出，所有时间为时间戳。

用户内容：
- 截图摘要（最多5张，按时间升序）：{screens}
- 最近聊天摘要：{recent_chat}
- 长期画像TopK：{profile}
- 设备/日历：{env}
```

## 6. 成本与SLO

- 参考：`doc/multi_screenshot_analysis_plan.md` 第十章。
- 降本手段：
  - 本地质量筛选（预期丢弃40-60%）；
  - Batch多图调用（半价）；
  - 去重与缓存（重复场景复用结论）；
  - 分层推理：轻量前处理 + 云端全面推理。
- SLO（建议）：
  - 端到端 P95 ≤ 2-3s；
  - 误触发率 < 3%；
  - 建议点击通过率 > 20%。

## 7. 隐私与安全

- 最小化：优先上传摘要与裁剪；原图短期存储并设置生命周期策略。
- 加密：端到端加密；密钥轮换与权限最小化；访问审计。
- 用户控制：可暂停、可清除、透明披露；敏感窗口黑名单。

## 8. 可观测性与治理

- 追踪：全链路 `trace_id`（客户端生成）。
- 指标：队列Lag、失败率、DLQ计数、端到端延迟、单图成本、IsJunk比例、Batch命中率。
- 日志：结构化日志 + 采样；`Reasons`（质量筛选）仅在调试或采样打印。

## 9. 评测方案

- 数据集：合成/匿名真实会话 + 截图序列（IDE、文档、会议、浏览器、视频）。
- 指标：
  - 任务识别准确率、意图预测准确率、情绪分类准确率；
  - 建议有效性（用户采纳率/满意度）；
  - 成本/时延对比（Batch on/off）。
- A/B：
  - 单帧 vs 多帧；
  - 无聊天 vs 有聊天摘要；
  - 去重策略开/关。

## 10. 落地里程碑（MVP→P2）

- P0（1周）：
  - 单帧 + 最近聊天摘要 → 结构化JSON结果；
  - 直传 `[]byte`；结果通过TCP推送；
  - 基础观测（延迟/失败率）。
- P1（2-3周）：
  - 多帧聚合 + Batch；
  - 质量筛选接入服务回调；
  - 结果置信度与证据链；
  - 简单隐私脱敏。
- P2（3-4周）：
  - 长期画像融合/Top-K检索；
  - 动作触发策略（冷却/优先级/场景白名单）；
  - 完整隐私与SLO治理、DLQ补偿流程。

## 11. 与项目其他组件的接口

- 长连接通信：`internal/connection/manager.go`、`public/connectionManager/`（推送结果）。
- 上下文与提示词：`internal/prompt/`、`internal/context/`；
- AI服务编排：`internal/ai/service.go`（可在此新增多模态入口）。
- 配置：`config/`（阈值/权重/批量参数/隐私策略）。

## 12. 未来扩展

- 场景专用Agent与Graph：编码、会议助理、学习助理等子Agent协同；
- 多模型协作：小模型做关键词提取/应用识别，大模型做全面推理；
- 本地轻量模型探索：端侧OCR/分类器以进一步降本与低延迟；
- 长期习惯建模：周/月模式分析与主动建议。

---

附录A：字段说明（节选）
- `ts`：消息创建时间戳（秒）。
- `deadline_ts`：任务截止时间戳（秒）。
- `hash`：aHash 64-bit（字符串表示）。
- `quality_score`：0-100。
- `is_junk`：布尔，垃圾帧不建议上云。

附录B：安全基线
- MQ 不存明文：客户端可选本地加密；云侧消费解密。
- OBS/S3 使用短期签名与生命周期策略；访问记录全审计。
