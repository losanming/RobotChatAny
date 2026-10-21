package com.robotchatany.android

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.lifecycleScope
import androidx.lifecycle.repeatOnLifecycle
import androidx.recyclerview.widget.LinearLayoutManager
import com.google.android.material.dialog.MaterialAlertDialogBuilder
import com.robotchatany.android.data.AppRepository
import com.robotchatany.android.databinding.FragmentTasksBinding
import com.robotchatany.android.model.TaskItem
import kotlinx.coroutines.launch

class TasksFragment : Fragment() {

    private var _binding: FragmentTasksBinding? = null
    private val binding get() = _binding!!
    private lateinit var adapter: TaskAdapter
    private val repository by lazy { AppRepository(requireContext().applicationContext) }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentTasksBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        adapter = TaskAdapter { task ->
            lifecycleScope.launch {
                repository.toggleTask(task.id)
            }
        }
        adapter.onLongPress = { task -> showTaskActions(task) }
        binding.taskList.layoutManager = LinearLayoutManager(requireContext())
        binding.taskList.adapter = adapter

        viewLifecycleOwner.lifecycleScope.launch {
            viewLifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
                repository.tasks.collect { tasks ->
                    adapter.submit(tasks)
                    binding.todayCount.text = "${tasks.count { !it.completed }} 项待办"
                    binding.doneCount.text = "${tasks.count { it.completed }} 项完成"
                    binding.summaryHint.text = if (tasks.any { it.completed }) "今天已经有推进，继续保持节奏。" else "先把最重要的一件事推进。"
                }
            }
        }
    }

    private fun showTaskActions(task: TaskItem) {
        if (!isAdded) return
        val items = arrayOf("编辑", "删除")
        MaterialAlertDialogBuilder(requireContext())
            .setTitle(task.title)
            .setMessage(task.detail)
            .setItems(items) { _, which ->
                when (which) {
                    0 -> showEditTaskDialog(task)
                    1 -> lifecycleScope.launch { repository.deleteTask(task.id) }
                }
            }
            .show()
    }

    private fun showEditTaskDialog(task: TaskItem) {
        val dialogView = layoutInflater.inflate(R.layout.dialog_edit_task, null)
        val titleInput = dialogView.findViewById<com.google.android.material.textfield.TextInputEditText>(R.id.editTaskTitle)
        val detailInput = dialogView.findViewById<com.google.android.material.textfield.TextInputEditText>(R.id.editTaskDetail)
        val timeInput = dialogView.findViewById<com.google.android.material.textfield.TextInputEditText>(R.id.editTaskTime)
        val priorityInput = dialogView.findViewById<com.google.android.material.textfield.TextInputEditText>(R.id.editTaskPriority)

        titleInput.setText(task.title)
        detailInput.setText(task.detail)
        timeInput.setText(task.schedule)
        priorityInput.setText(task.priority)

        MaterialAlertDialogBuilder(requireContext())
            .setTitle("编辑待办")
            .setView(dialogView)
            .setPositiveButton("保存") { _, _ ->
                val updated = task.copy(
                    title = titleInput.text?.toString().orEmpty().ifBlank { task.title },
                    detail = detailInput.text?.toString().orEmpty().ifBlank { task.detail },
                    schedule = timeInput.text?.toString().orEmpty().ifBlank { task.schedule },
                    priority = priorityInput.text?.toString().orEmpty().ifBlank { task.priority }
                )
                lifecycleScope.launch { repository.updateTask(updated) }
            }
            .setNegativeButton("取消", null)
            .show()
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
