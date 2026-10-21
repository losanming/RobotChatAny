package brain

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// AutoScheduler 自动调度器
type AutoScheduler struct {
	interval       int // 小时
	isRunning      bool
	stopChan       chan struct{}
	mutex          sync.RWMutex
	scheduledTasks []*ScheduledTask
}

// ScheduledTask 调度任务
type ScheduledTask struct {
	ID          string
	Name        string
	Interval    time.Duration
	LastRun     time.Time
	NextRun     time.Time
	Handler     func(ctx context.Context) error
	IsEnabled   bool
}

// NewAutoScheduler 创建自动调度器
func NewAutoScheduler(intervalHours int) *AutoScheduler {
	return &AutoScheduler{
		interval:       intervalHours,
		stopChan:       make(chan struct{}),
		scheduledTasks: make([]*ScheduledTask, 0),
	}
}

// Start 启动调度器
func (as *AutoScheduler) Start() {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	
	if as.isRunning {
		return
	}
	
	as.isRunning = true
	go as.run()
	
	log.Printf("自动调度器已启动，间隔: %d小时", as.interval)
}

// Stop 停止调度器
func (as *AutoScheduler) Stop() {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	
	if !as.isRunning {
		return
	}
	
	as.isRunning = false
	close(as.stopChan)
	
	log.Println("自动调度器已停止")
}

// run 运行调度器
func (as *AutoScheduler) run() {
	ticker := time.NewTicker(time.Duration(as.interval) * time.Hour)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			as.executeScheduledTasks()
		case <-as.stopChan:
			return
		}
	}
}

// executeScheduledTasks 执行调度任务
func (as *AutoScheduler) executeScheduledTasks() {
	as.mutex.RLock()
	defer as.mutex.RUnlock()
	
	now := time.Now()
	
	for _, task := range as.scheduledTasks {
		if !task.IsEnabled {
			continue
		}
		
		if now.After(task.NextRun) || now.Equal(task.NextRun) {
			log.Printf("执行调度任务: %s", task.Name)
			
			ctx := context.Background()
			if err := task.Handler(ctx); err != nil {
				log.Printf("任务执行失败 [%s]: %v", task.Name, err)
			} else {
				log.Printf("任务执行成功 [%s]", task.Name)
			}
			
			// 更新下次运行时间
			task.LastRun = now
			task.NextRun = now.Add(task.Interval)
		}
	}
}

// AddTask 添加调度任务
func (as *AutoScheduler) AddTask(name string, interval time.Duration, handler func(ctx context.Context) error) {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	
	task := &ScheduledTask{
		ID:        generateTaskID(name),
		Name:      name,
		Interval:  interval,
		LastRun:   time.Time{},
		NextRun:   time.Now().Add(interval),
		Handler:   handler,
		IsEnabled: true,
	}
	
	as.scheduledTasks = append(as.scheduledTasks, task)
	log.Printf("添加调度任务: %s, 间隔: %v", name, interval)
}

// RemoveTask 移除调度任务
func (as *AutoScheduler) RemoveTask(taskID string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	
	for i, task := range as.scheduledTasks {
		if task.ID == taskID {
			as.scheduledTasks = append(as.scheduledTasks[:i], as.scheduledTasks[i+1:]...)
			log.Printf("移除调度任务: %s", task.Name)
			return
		}
	}
}

// EnableTask 启用任务
func (as *AutoScheduler) EnableTask(taskID string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	
	for _, task := range as.scheduledTasks {
		if task.ID == taskID {
			task.IsEnabled = true
			log.Printf("启用调度任务: %s", task.Name)
			return
		}
	}
}

// DisableTask 禁用任务
func (as *AutoScheduler) DisableTask(taskID string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()
	
	for _, task := range as.scheduledTasks {
		if task.ID == taskID {
			task.IsEnabled = false
			log.Printf("禁用调度任务: %s", task.Name)
			return
		}
	}
}

// GetTasks 获取所有任务
func (as *AutoScheduler) GetTasks() []*ScheduledTask {
	as.mutex.RLock()
	defer as.mutex.RUnlock()
	
	tasksCopy := make([]*ScheduledTask, len(as.scheduledTasks))
	copy(tasksCopy, as.scheduledTasks)
	
	return tasksCopy
}

// IsRunning 检查是否运行中
func (as *AutoScheduler) IsRunning() bool {
	as.mutex.RLock()
	defer as.mutex.RUnlock()
	
	return as.isRunning
}

// generateTaskID 生成任务ID
func generateTaskID(name string) string {
	return fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
}
