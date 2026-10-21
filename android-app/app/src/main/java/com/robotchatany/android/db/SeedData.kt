package com.robotchatany.android.db

object SeedData {
    fun defaultProfile() = ProfileEntity(
        id = 1,
        displayName = "威猛先生",
        role = "独立开发者 / 主动式 AI 伙伴",
        bio = "喜欢把个人系统做成一个能主动帮忙的助手。",
        prompt = "你是我的私人 AI 助手，优先主动提醒、整理待办、总结关键信息，语气简洁且温和。",
        voiceStyle = "自然、简短、有一点点情绪表达",
        cloudSync = false,
        proactiveReminder = true
    )

    fun defaultTasks() = listOf(
        TaskEntity(title = "完成晨间复盘", detail = "回顾昨天的 3 件事，写下今天最重要的一件。", schedule = "09:00", priority = "高", completed = false),
        TaskEntity(title = "午间散步 20 分钟", detail = "给大脑一个短暂的恢复窗口。", schedule = "12:20", priority = "中", completed = true),
        TaskEntity(title = "整理 brain 待办", detail = "把今天新增想法放入统一清单。", schedule = "15:00", priority = "高", completed = false),
        TaskEntity(title = "晚上 10 点前关闭屏幕", detail = "进入低打扰模式，准备睡眠。", schedule = "22:00", priority = "低", completed = false)
    )

    fun defaultMessages() = listOf(
        MessageEntity(sender = "AI", content = "早上好，我已经准备好帮你安排今天的节奏。", timeLabel = "08:30", type = "TEXT"),
        MessageEntity(sender = "USER", content = "帮我看看今天有什么计划？", timeLabel = "08:31", type = "TEXT"),
        MessageEntity(sender = "AI", content = "我会把待办、提醒和你的偏好都放在一起看。", timeLabel = "08:31", type = "TEXT")
    )
}
