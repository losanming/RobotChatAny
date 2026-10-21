package com.robotchatany.android

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.lifecycleScope
import androidx.lifecycle.repeatOnLifecycle
import com.robotchatany.android.data.AppRepository
import com.robotchatany.android.databinding.FragmentProfileBinding
import com.robotchatany.android.model.ProfileState
import kotlinx.coroutines.launch

class ProfileFragment : Fragment() {

    private var _binding: FragmentProfileBinding? = null
    private val binding get() = _binding!!
    private val repository by lazy { AppRepository(requireContext().applicationContext) }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentProfileBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                repository.profile.collect { renderProfile(it) }
            }
        }

        binding.saveButton.setOnClickListener {
            val updated = ProfileState(
                displayName = binding.nameInput.text?.toString().orEmpty().ifBlank { "威猛先生" },
                role = binding.roleInput.text?.toString().orEmpty().ifBlank { "独立开发者 / 主动式 AI 伙伴" },
                bio = binding.bioInput.text?.toString().orEmpty().ifBlank { "喜欢把个人系统做成一个能主动帮忙的助手。" },
                prompt = binding.promptInput.text?.toString().orEmpty().ifBlank { "你是我的私人 AI 助手，优先主动提醒、整理待办、总结关键信息，语气简洁且温和。" },
                voiceStyle = binding.voiceInput.text?.toString().orEmpty().ifBlank { "自然、简短、有一点点情绪表达" },
                cloudSync = binding.cloudSwitch.isChecked,
                proactiveReminder = binding.reminderSwitch.isChecked
            )
            viewLifecycleOwner.lifecycleScope.launch {
                repository.updateProfile(updated)
                Toast.makeText(requireContext(), "已保存本地配置", Toast.LENGTH_SHORT).show()
            }
        }
    }

    private fun renderProfile(profile: ProfileState) {
        binding.nameInput.setText(profile.displayName)
        binding.roleInput.setText(profile.role)
        binding.bioInput.setText(profile.bio)
        binding.promptInput.setText(profile.prompt)
        binding.voiceInput.setText(profile.voiceStyle)
        binding.cloudSwitch.isChecked = profile.cloudSync
        binding.reminderSwitch.isChecked = profile.proactiveReminder
        binding.previewText.text = "${profile.displayName} · ${profile.role}\n${profile.prompt}"
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
