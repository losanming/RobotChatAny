package com.robotchatany.android

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.core.view.isVisible
import androidx.recyclerview.widget.RecyclerView
import com.robotchatany.android.databinding.ItemChatMessageBinding
import com.robotchatany.android.model.ChatMessage
import com.robotchatany.android.model.MessageType
import com.robotchatany.android.model.Sender

class ChatAdapter : RecyclerView.Adapter<ChatAdapter.ChatViewHolder>() {

    var onLongPress: ((ChatMessage) -> Unit)? = null

    private val items = mutableListOf<ChatMessage>()

    fun submit(newItems: List<ChatMessage>) {
        items.clear()
        items.addAll(newItems)
        notifyDataSetChanged()
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ChatViewHolder {
        val binding = ItemChatMessageBinding.inflate(LayoutInflater.from(parent.context), parent, false)
        return ChatViewHolder(binding)
    }

    override fun getItemCount(): Int = items.size

    override fun onBindViewHolder(holder: ChatViewHolder, position: Int) {
        holder.bind(items[position])
    }

    inner class ChatViewHolder(private val binding: ItemChatMessageBinding) : RecyclerView.ViewHolder(binding.root) {

        fun bind(item: ChatMessage) {
            val isUser = item.sender == Sender.USER
            binding.senderLabel.isVisible = !isUser
            binding.avatar.isVisible = !isUser
            binding.root.gravity = if (isUser) android.view.Gravity.END else android.view.Gravity.START
            binding.senderLabel.text = if (isUser) "你" else "AI"
            binding.avatarLabel.text = if (isUser) "你" else "AI"
            binding.messageText.text = item.content
            binding.timeText.text = item.timeLabel
            binding.imagePreview.isVisible = item.type == MessageType.IMAGE

            if (item.type == MessageType.IMAGE) {
                binding.messageText.text = "${item.content}\n点击后续可接入真实图片选择"
                binding.imagePreview.setImageResource(android.R.drawable.ic_menu_gallery)
            }

            if (item.type == MessageType.VOICE) {
                binding.messageText.text = "🎤 ${item.content}"
            }

            val bubbleColor = if (isUser) 0xFFF4B860.toInt() else 0xFF1B2331.toInt()
            binding.bubbleCard.setCardBackgroundColor(bubbleColor)
            binding.messageText.setTextColor(if (isUser) 0xFF1A1307.toInt() else 0xFFF5F1E8.toInt())
            binding.timeText.setTextColor(if (isUser) 0x661A1307 else 0x66F5F1E8)

            itemView.setOnLongClickListener {
                this@ChatAdapter.onLongPress?.invoke(item)
                true
            }
        }
    }
}
