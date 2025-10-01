# RobotChatAny：下一代主动智能聊天机器人

> 让AI主动陪伴，赋予每一次对话温度与个性。

## 项目简介
RobotChatAny 是一款创新的智能聊天机器人，专注于主动对话和个性化陪伴体验。我们不仅让AI能主动发起自然交流，还能无缝对接你的手表、手机等智能设备，采集健康、运动、睡眠等数据，生成专属报告，真正实现“懂你”的AI伙伴。

## 为什么选择 RobotChatAny？
- **主动对话体验**：突破传统被动应答，AI可根据情境、时间、事件主动开启话题，打造更自然的互动氛围。
- **多维度权重系统**：通过情感、兴趣、行为等多维度动态调整AI个性，让每个用户都能拥有独一无二的AI朋友。
- **极致可扩展性**：插件化架构，支持天气、音乐、热搜、新闻等多种扩展，开发者可轻松自定义新功能。
- **健康数据融合**：原生支持对接智能硬件，自动采集和分析健康数据，助力用户科学生活。
- **高性能技术栈**：Golang+xorm后端，TCP长连接，PostgreSQL、Redis、OBS等企业级组件，保障高并发和稳定性。

## 面向开发者的邀请
RobotChatAny 仍在快速迭代中，欢迎每一位热爱AI、对人机交互充满想象的你加入我们：
- 参与核心功能开发，体验高性能Golang后端实践
- 设计和实现有趣的AI插件，扩展无限应用场景
- 共建健康数据、主动对话、情感计算等创新模块
- 贡献文档、测试、社区支持，共同打造开源生态

## 项目愿景
- 让AI成为真正有温度的数字伙伴，陪伴每一个需要关怀的人
- 推动主动智能对话技术落地，服务个人、企业与社会
- 打造开放、协作、创新的开发者社区

## 产品初步方案
<img width="4832" height="5796" alt="初版总体方案" src="https://github.com/user-attachments/assets/059536d3-86f5-4f08-8fd7-ff50ca0aff38" />

## 技术亮点
- 多模态能力：视觉（桌面/摄像头、质量筛选、OCR摘要）、听觉（VAD/唤醒词/ASR）、TTS（多音色/多情绪）
- 位置与环境：GPS（WGS84→GCJ02）、高德逆地理/天气/实时路况、围栏/停留/通勤学习
- 事件总线：NSQ 实时消费，指针+元数据传输图片，OBS/S3/Redis 临时 Key 存储（TTL 10–30 分钟）
- 图像质量筛选：清晰度/对比度/均匀度/aHash 去重与阈值权重（`pkg/screenshot/quality/*`）
- 上下文黑板：AI 大脑统一读取“当前状态 + 最近时间线”，多模态置信度融合后决策
- 技术栈：Golang 1.21+ / xorm ORM / PostgreSQL / Redis / OBS / TLS / 自定义协议 + TCP 长连接
- 时间表示：全局统一使用秒级时间戳（int64, s），获取用 `time.Now().Unix()`；与外部 API 的毫秒/字符串时间在 provider 层转换为“秒”后再进入域模型与数据库
- Awareness 缓存与重试 + GPS 条件触发：内存 TTL + 重试封装（逆地理/周边/天气），GPS 上报仅在时间/位移/TTL 触发时刷新昂贵数据，其他时刻复用上次结果并即时更新位置。

## 今日更新（2025-10-01）
- **[Eino 工具集]** 新增并增强 `einoUtils/brainTools/`：
  - `todo_tool`：
    - `action` 改为可选；支持 `title`/`titles`/`text` 自然语言解析与批量新增；
    - 自动从中文句子中拆出多个待办（如“今天上午10点去逛动物园,下午3点看电影”），逐条写入；
    - 内置简单中文时间解析（今天/明天 + 上午/下午 + X点），为每条待办推断 `due_ts`（秒）；
    - 统一返回结构并在批量时返回 `count`；所有时间戳统一为“秒”。
  - `events_tool`：完成 `AddEvent`/`ListEvents` 的调用通路。
  - `task_tool`：完成 `SetCurrentTask`/`ClearCurrentTask`/`GetCurrentTask` 的调用通路。
- **[Brain 接口与时间戳统一]** `internal/brain/service.go`
  - 新增：`AddEvent(e Event)`、`ListEvents() []Event`。
  - 新增：`SetCurrentTask(t Task)`、`ClearCurrentTask()`、`GetCurrentTask() *Task`。
  - 修正：`SetHealth(h HealthState)`、`PatchHealth(...)` 全部统一使用 `time.Now().Unix()`（秒）。
- **[Demo 覆盖]** 更新 `cmd/test/eino_todo_demo.go` 示例问句，验证一句多待办与时间抽取的落库效果。

> 说明：新增目录 `einoUtils/brainTools/` 存放与 Eino 框架对接的工具集（todo/events/task/profile/health）。

## 今日更新（2025-09-30）
- **[时间戳统一]** 全项目统一“秒级时间戳”：一律使用 `time.Now().Unix()`。若外部 API 返回毫秒或字符串时间，均在 provider 层转换为“秒”。
- **[GPSService 初始化]** 在 `internal/brain/gps_service.go` 新增 `NewGPSService(parent context.Context, brain *Brain, aware *awareness.Service, opts ...GPSOption)`，支持可选项：队列大小、汇聚超时与策略。
- **[AMap 兼容策略]** 仅在 `pkg/location/provider/amap/` 层兼容高德“字段类型不稳定”（如 `address` 可能为 string 或 []string），通过自定义反序列化类型统一为稳定输出；对外域模型保持不变。
- **[通勤习惯讨论]** 形成方案：记录单次通勤（起点/终点/出发/到达/用时/路过地点），按周聚合生成“通勤习惯档案”（出发时间分布、用时P50/P90、主路线）。


## 快速开始
1. 克隆本仓库，安装依赖
2. 配置数据库、Redis、OBS等服务
3. 启动服务端，体验AI主动聊天与插件扩展

## 目录结构
 ``` go
 ├── Readme.md
├── api // 外部api接口定义
│   ├── auth
│   │   ├── captcha.go
│   │   └── login.go
│   ├── contract // 协议定义
│   │   ├── auth
│   │   │   ├── request.go
│   │   │   └── response.go
│   │   └── user
│   └── user
├── asset // 静态文件
│   ├── pic
│   │   └── city.jpeg
│   ├── resources.json
│   └── triggers.json
├── cmd // 程序主入口
│   ├── main.go
│   ├── test // 测试文件
│   │   ├── commendLine.go
│   │   ├── documentLoader.go
│   │   ├── duckduckgo.go
│   │   ├── prompt.go
│   │   ├── resource.go
│   │   ├── test_eine.go
│   │   └── 需求.md
│   └── testFile
│       ├── food.db
│       ├── test_boheFood.go
│       └── utils
│           └── goroutine_pool.go
├── config //配置
│   ├── db.go
│   ├── gitee.go
│   ├── init.go
│   ├── log.go
│   ├── oss_local.go
│   ├── oss_minio.go
│   ├── oss_obs.go
│   ├── redis.go
│   ├── system.go
│   ├── tcp.go
│   └── web.go
├── config.yaml // 配置文件
├── doc // 相关文档
│   ├── ai_chat_design.md
│   ├── ai_weight_system_design.md
│   ├── context_management_design.md
│   ├── hardware_integration_design.md
│   ├── tcp_design.md
│   ├── technical_design.md
│   └── user_personality_design.md
├── eineUtils // ai框架eine实现的mcp工具
│   ├── duckSeach
│   │   └── duck.go
│   └── httpSearch
│       ├── httpCommen.go
│       ├── httpGet.go
│       └── searchEngine.go
├── global // 全局默认参数
│   └── common.go
├── go.mod
├── go.sum
├── initialize // 初始化
│   ├── ai.go
│   ├── config.go
│   ├── connection.go
│   ├── context.go
│   ├── ctx.go
│   ├── fetcher.go
│   ├── messageQueue.go
│   ├── postgres.go
│   ├── prompt.go
│   ├── redis.go
│   ├── tcpServer.go
│   ├── trigger.go
│   └── web.go
├── internal // 业务逻辑
│   ├── ai // ai服务
│   │   └── service.go
│   ├── connection // 连接和心跳管理器
│   │   ├── heartbeat.go
│   │   └── manager.go
│   ├── context // 用户上下文管理器
│   │   ├── manager.go
│   │   └── types.go
│   ├── prompt // prompt生成管理器
│   │   ├── context_prompt.go
│   │   └── generator.go
│   ├── protocol // 协议以及加解密管理器
│   │   ├── codec.go
│   │   ├── codec_test.go
│   │   ├── message.go
│   │   ├── protocol.go
│   │   ├── secure_codec.go
│   │   └── security.go
│   ├── push // 推送服务管理器
│   │   └── service.go
│   ├── queue // 消息队列管理器
│   │   ├── dispatcher.go
│   │   ├── message.go
│   │   └── retry.go
│   ├── resource // 资源管理器
│   │   ├── classifier.go
│   │   ├── fetcher.go
│   │   └── types.go
│   ├── tcp // TCP服务管理器
│   │   └── server.go
│   └── trigger // 触发器管理器
│       ├── config.go
│       ├── manager.go
│       ├── types.go
│       └── utils.go
├── load // 加载管理器
│   └── load.go
├── logs // 日志管理器
│   └── app.log
├── middleware // 中间件管理器
│   ├── cors.go
│   └── jwt.go
├── model // 模型
│   ├── login.go
│   ├── slider_captcha.go
│   ├── third_party.go
│   └── user.go
├── pkg // 工具包
│   ├── banner // 指纹
│   │   └── banner.go
│   ├── cache // 缓存
│   │   └── cache.go
│   ├── goroutine_pool.go // 协程池
│   ├── hash.go // 加密
│   ├── httpClient // http客户端
│   │   └── httpClient.go
│   ├── logger // 日志
│   │   └── logger.go
│   ├── postgres // 数据库
│   │   └── postgres.go
│   ├── response // 响应
│   │   ├── api.go 
│   │   ├── init.go 
│   │   └── response.go 
│   ├── str.go // 字符串处理
│   ├── screenshot // 图片质量分析
│   │   ├── quality // 图片质量分析
│   │   │   └── analyzer.go // 质量评估器
│   │   │   └── metrics.go // 拉普拉斯卷积等质量分析算法
│   │   ├── screenshot.go // 截图压缩函数
│   │   ├── service.go // 截图压缩管理器
│   ├── upload // 文件上传
│   │   ├── local.go // 本地上传
│   │   ├── minio_oss.go //minio上传
│   │   ├── obs.go // 华为obs上传
│   │   └── upload.go // 上传接口
│   ├── utils // 部分工具
│   │   └── aes.go
│   ├── uuid.go
│   └── web
│       └── web.go
├── public // 公共会话存储
│   ├── ai.go
│   ├── captcha
│   │   └── background.png
│   ├── config.go
│   ├── connection.go
│   ├── context.go
│   ├── ctx.go
│   ├── fetcher.go
│   ├── messageQueue.go
│   ├── postgres.go
│   ├── prompt.go
│   ├── redis.go
│   ├── tcpServer.go
│   ├── trigger.go
│   └── web.go
├── router // 路由
│   ├── auth
│   │   └── authRouter.go
│   ├── router.go
│   ├── system
│   │   └── systemRouter.go
│   └── user
│       └── userRouter.go
├── service // 服务
│   ├── captcha_service.go
│   ├── login_service.go
│   ├── third_party_service.go
│   └── user_service.go
└── webHandler // web控制器
    └── webHander.go
```
## 项目计划以及完成度
* __总体目标__：打造具备“主动对话 + 情境感知 + 个性化权重 + 设备与健康数据融合”的AI助手，支持长连接实时交互与可扩展插件生态。

### 路线图（Roadmap）
* __阶段1·基础架构__（已完成 100%）
  - TCP长连接、认证、心跳与连接管理器，安全编解码与协议栈（`internal/tcp`、`internal/protocol`、`initialize/tcpServer.go`）。
  - 配置、日志、缓存、数据库与资源加载（`initialize/`、`pkg/logger`、`pkg/postgres`、`pkg/cache`）。
* __阶段2·AI对话与权重系统__（进行中 ≈70%）
  - 上下文管理与Token优化设计已完成文档与实现框架（`doc/context_management_design.md`，`internal/context`、`internal/prompt`）。
  - 权重系统与触发机制设计完成（`doc/ai_weight_system_design.md`、`doc/ai_chat_design.md`），接入与调参中。
* __阶段3·情境感知（桌面识别）__（进行中 ≈40%）
  - 截图服务与压缩（`pkg/screenshot/screenshot.go`、`pkg/screenshot/service.go`）。
  - 本地图片质量筛选：清晰度/对比度/均匀度/aHash去重等（`pkg/screenshot/quality/analyzer.go`、`metrics.go`，参见`doc/multi_screenshot_analysis_plan.md`第十一章）。
  - 待将质量筛选接入服务回调并丢弃垃圾帧，接阿里云视觉分析。
* __阶段4·设备/健康数据融合__（规划中）
  - 通过硬件与MQTT/语音链路集成，形成语音-文本-语音闭环（`doc/hardware_integration_design.md`）。
* __阶段5·产品化与生态__（规划中）
  - 插件体系完善、监控与埋点、A/B调参、权限与合规。

### 已完成（Done）
* __通信与安全__：TCP长连接、认证、心跳、加密编解码（`internal/connection`、`internal/protocol`）。
* __初始化与基础设施__：配置、Redis/PostgreSQL、资源与触发器（`initialize/`、`internal/trigger`、`internal/resource`）。
* __上下文与Prompt框架__：上下文类型与生成器（`internal/context`、`internal/prompt`）。
* __AI设计文档__：`doc/ai_chat_design.md`、`doc/ai_weight_system_design.md`、`doc/context_management_design.md`。
* __桌面截图基础能力__：定时截图、压缩、回调扩展（`pkg/screenshot/`）。
* __本地图片质量筛选模块__：评分/判废/去重，内置阈值与权重（`pkg/screenshot/quality/*`）。

### 进行中（In Progress）
* __质量筛选接入链路__：在 `pkg/screenshot/service.go` 回调处接 `quality.Analyzer`，`IsJunk` 直接丢弃。
* __云端视觉识别__：接入阿里云视觉/多模态识别，结合本地筛选降低成本（参考 `doc/multi_screenshot_analysis_plan.md` 成本估算）。
* __权重系统落地与埋点__：将权重参数配置化（阈值/频率/情绪与话题权重），埋点数据用于闭环调参。
* __数据源整合__：热搜/天气/资讯等外部API接入与权重融合（见 `doc/ai_chat_design.md`）。

### 下一步计划（Next）
1. 在 `pkg/screenshot/service.go` 接入 `quality.Analyzer` 并新增埋点：Score/IsJunk/Reasons 分布。
2. 接入阿里云视觉分析，按场景做批量/合并请求与缓存，控制成本（详见 `doc/multi_screenshot_analysis_plan.md` 第十章）。
3. 将权重阈值与触发策略配置化（环境变量/配置文件），支持运行时热更新与灰度。
4. 补齐主动触发策略上线闭环：触发→生成→推送→反馈→自适应（`internal/trigger` 联动 `internal/ai`）。
5. 增加监控告警与性能压测：长连接稳定性、截图管线吞吐、视觉调用SLA（Grafana/Prometheus 规划）。

### 里程碑与完成度
* __M1 基础通信稳定上线__（100%）：长连接、认证、心跳与协议、安全编解码。
* __M2 上下文与Prompt框架__（100%）：上下文结构、生成器与优化策略文档。
* __M3 截图与质量筛选__（80%）：截图服务与质量评估器完成，服务侧接入与埋点进行中。
* __M4 情境感知上线（首版）__（40%）：接入云视觉、落地主动触发若干场景。
* __M5 权重系统实装与调参平台__（30%）：配置化与数据闭环在线调参。
* __M6 设备/健康数据融合（语音链路）__（0% 规划）：按 `doc/hardware_integration_design.md` 打通链路。

说明：上述百分比为工程实现进度的粗略估算，会随代码合并与联调推进而调整。

## 加入我们
无论你是后端高手、AI算法爱好者、产品设计师还是热衷体验的普通用户，都欢迎加入RobotChatAny！一起探索主动智能的无限可能。

> 让AI不再冷冰冰，让陪伴触手可及。


> 我的联系方式qq:1321291309
---
如有建议、想法或合作意向，请提交 Issue 或 PR，或直接联系项目维护者。



