package com.robotchatany.android

import android.os.Bundle
import android.content.ClipData
import android.content.ClipboardManager
import android.text.Editable
import android.text.TextWatcher
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.core.view.isVisible
import com.google.android.material.dialog.MaterialAlertDialogBuilder
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.lifecycle.repeatOnLifecycle
import androidx.lifecycle.Lifecycle
import com.robotchatany.android.data.AppRepository
import com.robotchatany.android.databinding.FragmentChatBinding
import com.robotchatany.android.model.ChatMessage
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import kotlinx.coroutines.Job

class ChatFragment : Fragment() {

    private var _binding: FragmentChatBinding? = null
    private val binding get() = _binding!!
    private lateinit var adapter: ChatAdapter
    private val repository by lazy { AppRepository(requireContext().applicationContext) }
    private var isRecording = false
    private var voiceJob: Job? = null

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentChatBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        adapter = ChatAdapter()
        adapter.onLongPress = { message -> showMessageActions(message) }
        binding.chatList.layoutManager = LinearLayoutManager(requireContext()).apply {
            stackFromEnd = true
        }
        binding.chatList.adapter = adapter

        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                launch {
                    repository.messages.collect { messages ->
                        adapter.submit(messages)
                        binding.messageCount.text = "${messages.size} 条对话"
                        binding.chatList.post {
                            if (adapter.itemCount > 0) {
                                binding.chatList.smoothScrollToPosition(adapter.itemCount - 1)
                            }
                        }
                    }
                }
            }
        }

        binding.sendButton.setOnClickListener {
            val text = binding.messageInput.text?.toString()?.trim().orEmpty()
            if (text.isBlank()) return@setOnClickListener
            sendText(text)
        }

        binding.emojiButton.setOnClickListener {
            binding.messageInput.append("😊")
        }

        binding.imageButton.setOnClickListener {
            viewLifecycleOwner.lifecycleScope.launch {
                repository.addUserImage("图片消息 · 预览占位")
                delay(100)
                repository.addAiReply("图片")
            }
        }

        binding.voiceButton.setOnClickListener {
            if (isRecording) return@setOnClickListener
            isRecording = true
            binding.voiceHint.isVisible = true
            binding.voiceHint.text = "正在聆听你的语音..."
            voiceJob?.cancel()
            voiceJob = viewLifecycleOwner.lifecycleScope.launch {
                delay(1200)
                if (_binding == null) return@launch
                isRecording = false
                binding.voiceHint.isVisible = false
                val transcript = listOf(
                    "帮我把今晚的计划整理一下。",
                    "提醒我明天上午优先处理待办。",
                    "把这个想法记成一个任务。"
                ).random()
                repository.addUserVoice(transcript)
                delay(100)
                repository.addAiReply("语音")
            }
        }

        binding.messageInput.addTextChangedListener(object : TextWatcher {
            override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) = Unit
            override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) = Unit
            override fun afterTextChanged(s: Editable?) {
                binding.sendButton.isEnabled = !s.isNullOrBlank()
            }
        })

        binding.sendButton.isEnabled = false
    }

    private fun sendText(text: String) {
        viewLifecycleOwner.lifecycleScope.launch {
            repository.addUserText(text)
            delay(100)
            repository.addAiReply(text)
        }
        binding.messageInput.setText("")
    }

    private fun showMessageActions(message: ChatMessage) {
        if (!isAdded) return
        val options = if (message.sender.name == "USER") arrayOf("复制", "删除") else arrayOf("复制", "删除")
        MaterialAlertDialogBuilder(requireContext())
            .setTitle(if (message.sender.name == "USER") "你的消息" else "AI 消息")
            .setMessage(message.content)
            .setItems(options) { _, which ->
                when (which) {
                    0 -> copyToClipboard(message.content)
                    1 -> lifecycleScope.launch { repository.deleteMessage(message.id) }
                }
            }
            .show()
    }

    private fun copyToClipboard(text: String) {
        val clipboard = requireContext().getSystemService(ClipboardManager::class.java)
        clipboard?.setPrimaryClip(ClipData.newPlainText("chat", text))
    }

    override fun onDestroyView() {
        voiceJob?.cancel()
        super.onDestroyView()
        _binding = null
    }
}
