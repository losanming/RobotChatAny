package com.robotchatany.android

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.robotchatany.android.databinding.ItemTaskBinding
import com.robotchatany.android.model.TaskItem

class TaskAdapter(
    private val onToggle: (TaskItem) -> Unit
) : RecyclerView.Adapter<TaskAdapter.TaskViewHolder>() {

    var onLongPress: ((TaskItem) -> Unit)? = null

    private val items = mutableListOf<TaskItem>()

    fun submit(newItems: List<TaskItem>) {
        items.clear()
        items.addAll(newItems)
        notifyDataSetChanged()
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): TaskViewHolder {
        val binding = ItemTaskBinding.inflate(LayoutInflater.from(parent.context), parent, false)
        return TaskViewHolder(binding, onToggle)
    }

    override fun getItemCount(): Int = items.size

    override fun onBindViewHolder(holder: TaskViewHolder, position: Int) {
        holder.bind(items[position])
    }

    inner class TaskViewHolder(
        private val binding: ItemTaskBinding,
        private val onToggle: (TaskItem) -> Unit
    ) : RecyclerView.ViewHolder(binding.root) {

        fun bind(item: TaskItem) {
            binding.taskTitle.text = item.title
            binding.taskDetail.text = item.detail
            binding.taskTime.text = item.schedule
            binding.priorityTag.text = item.priority
            binding.taskCheck.isChecked = item.completed
            binding.priorityStripe.setBackgroundColor(
                when (item.priority) {
                    "高" -> 0xFFF4B860.toInt()
                    "中" -> 0xFFFF7A59.toInt()
                    else -> 0xFF6ED3A9.toInt()
                }
            )
            binding.taskCheck.setOnClickListener {
                onToggle(item)
            }
            binding.root.alpha = if (item.completed) 0.55f else 1f

            binding.root.setOnLongClickListener {
                this@TaskAdapter.onLongPress?.invoke(item)
                true
            }
        }
    }
}
