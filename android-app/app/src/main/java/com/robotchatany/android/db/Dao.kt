package com.robotchatany.android.db

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import androidx.room.Update
import kotlinx.coroutines.flow.Flow

@Dao
interface ProfileDao {
    @Query("SELECT * FROM profile WHERE id = 1 LIMIT 1")
    fun observeProfile(): Flow<ProfileEntity?>

    @Query("SELECT COUNT(*) FROM profile")
    suspend fun count(): Int

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsert(profile: ProfileEntity)
}

@Dao
interface TaskDao {
    @Query("SELECT * FROM tasks ORDER BY id DESC")
    fun observeTasks(): Flow<List<TaskEntity>>

    @Query("SELECT COUNT(*) FROM tasks")
    suspend fun count(): Int

    @Query("SELECT * FROM tasks WHERE id = :id LIMIT 1")
    suspend fun getById(id: Long): TaskEntity?

    @Query("DELETE FROM tasks WHERE id = :id")
    suspend fun deleteById(id: Long)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(tasks: List<TaskEntity>)

    @Update
    suspend fun update(task: TaskEntity)

    @Query("DELETE FROM tasks")
    suspend fun clear()
}

@Dao
interface MessageDao {
    @Query("SELECT * FROM messages ORDER BY id ASC")
    fun observeMessages(): Flow<List<MessageEntity>>

    @Query("SELECT COUNT(*) FROM messages")
    suspend fun count(): Int

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(message: MessageEntity): Long

    @Query("DELETE FROM messages WHERE id = :id")
    suspend fun deleteById(id: Long)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(messages: List<MessageEntity>)

    @Query("DELETE FROM messages")
    suspend fun clear()
}
