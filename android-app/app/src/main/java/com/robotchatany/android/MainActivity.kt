package com.robotchatany.android

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.fragment.app.Fragment
import androidx.fragment.app.commit
import androidx.lifecycle.lifecycleScope
import com.robotchatany.android.data.AppRepository
import com.robotchatany.android.databinding.ActivityMainBinding
import kotlinx.coroutines.launch

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private val repository by lazy { AppRepository(applicationContext) }
    private val chatFragment by lazy { ChatFragment() }
    private val tasksFragment by lazy { TasksFragment() }
    private val profileFragment by lazy { ProfileFragment() }
    private var currentTabId: Int = R.id.nav_chat

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        currentTabId = savedInstanceState?.getInt(KEY_TAB_ID, R.id.nav_chat) ?: R.id.nav_chat

        lifecycleScope.launch {
            repository.seedIfNeeded()
        }

        binding.bottomNav.setOnItemSelectedListener { item ->
            currentTabId = item.itemId
            when (item.itemId) {
                R.id.nav_chat -> openFragment(chatFragment)
                R.id.nav_tasks -> openFragment(tasksFragment)
                R.id.nav_profile -> openFragment(profileFragment)
            }
            true
        }

        if (savedInstanceState == null) {
            openFragment(chatFragment)
        }

        binding.bottomNav.menu.findItem(currentTabId)?.isChecked = true
    }

    override fun onSaveInstanceState(outState: Bundle) {
        outState.putInt(KEY_TAB_ID, currentTabId)
        super.onSaveInstanceState(outState)
    }

    private fun openFragment(fragment: Fragment) {
        val tag = fragment::class.java.name
        supportFragmentManager.commit {
            setReorderingAllowed(true)
            replace(R.id.contentContainer, fragment, tag)
        }
    }

    private companion object {
        const val KEY_TAB_ID = "current_tab_id"
    }
}
