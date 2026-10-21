package com.robotchatany.android.db

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "profile")
data class ProfileEntity(
    @PrimaryKey val id: Int = 1,
    val displayName: String,
    val role: String,
    val bio: String,
    val prompt: String,
    val voiceStyle: String,
    val cloudSync: Boolean,
    val proactiveReminder: Boolean
)

@Entity(tableName = "tasks")
data class TaskEntity(
    @PrimaryKey(autoGenerate = true) val id: Long = 0,
    val title: String,
    val detail: String,
    val schedule: String,
    val priority: String,
    val completed: Boolean
)

@Entity(tableName = "messages")
data class MessageEntity(
    @PrimaryKey(autoGenerate = true) val id: Long = 0,
    val sender: String,
    val content: String,
    val timeLabel: String,
    val type: String
)
