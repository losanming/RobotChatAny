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

## 技术亮点
- Golang 1.21+ / xorm ORM / PostgreSQL / Redis / OBS
- TCP长连接保障实时通信
- 插件系统支持热插拔与快速扩展
- 健康数据采集与分析
- 高安全性（TLS、加密存储、权限控制等）

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

## 加入我们
无论你是后端高手、AI算法爱好者、产品设计师还是热衷体验的普通用户，都欢迎加入RobotChatAny！一起探索主动智能的无限可能。

> 让AI不再冷冰冰，让陪伴触手可及。


> 我的联系方式qq:1321291309
---
![IMG_1867](https://github.com/user-attachments/assets/4685c2ea-53c8-425b-a6d6-d478b27cf7a8)
如有建议、想法或合作意向，请提交 Issue 或 PR，或直接联系项目维护者。
