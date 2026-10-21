package com.robotchatany.android.model

enum class Sender {
    USER,
    AI
}

enum class MessageType {
    TEXT,
    IMAGE,
    VOICE
}

data class ChatMessage(
    val id: Long,
    val sender: Sender,
    val content: String,
    val timeLabel: String,
    val type: MessageType = MessageType.TEXT
)

data class TaskItem(
    val id: Long,
    val title: String,
    val detail: String,
    val schedule: String,
    val priority: String,
    var completed: Boolean = false
)

data class ProfileState(
    val displayName: String,
    val role: String,
    val bio: String,
    val prompt: String,
    val voiceStyle: String,
    val cloudSync: Boolean,
    val proactiveReminder: Boolean
)
