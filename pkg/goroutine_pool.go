package pkg

import (
	"context"
	"runtime"
	"sync"
)

// Task 表示一个可以执行的任务
type Task func() error

// Result 表示任务执行的结果
type Result struct {
	Error error
}

// GoPool 协程池结构体
type GoPool struct {
	workerNum  int         // 工作协程数量
	taskChan   chan Task   // 任务通道
	resultChan chan Result // 结果通道
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewGoPool 创建一个新的协程池
func NewGoPool() *GoPool {
	workerNum := runtime.NumCPU() * 4
	if workerNum < 20 {
		workerNum = 20
	}
	ctx, cancel := context.WithCancel(context.Background())

	pool := &GoPool{
		workerNum:  workerNum,
		taskChan:   make(chan Task),
		resultChan: make(chan Result),
		ctx:        ctx,
		cancel:     cancel,
	}

	// 启动工作协程
	pool.startWorkers()

	return pool
}

// startWorkers 启动工作协程
func (p *GoPool) startWorkers() {
	for i := 0; i < p.workerNum; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case task, ok := <-p.taskChan:
					if !ok {
						return
					}
					// 执行任务并发送结果
					err := task()
					p.resultChan <- Result{Error: err}
				case <-p.ctx.Done():
					return
				}
			}
		}()
	}
}

// Submit 提交一个任务到协程池
func (p *GoPool) Submit(task Task) {
	p.taskChan <- task
}

// SubmitAndWait 提交一个任务并等待其完成，返回错误
func (p *GoPool) SubmitAndWait(task Task) error {
	p.Submit(task)
	result := <-p.resultChan
	return result.Error
}

// SubmitAll 批量提交任务并等待所有任务完成，返回所有错误
func (p *GoPool) SubmitAll(tasks []Task) []error {
	resultChan := make(chan Result, len(tasks))

	// 提交所有任务
	for _, task := range tasks {
		t := task // 创建副本避免闭包问题
		go func() {
			p.Submit(t)
			resultChan <- <-p.resultChan
		}()
	}

	// 收集结果
	errors := make([]error, 0)
	for i := 0; i < len(tasks); i++ {
		result := <-resultChan
		if result.Error != nil {
			errors = append(errors, result.Error)
		}
	}

	return errors
}

// Close 关闭协程池
func (p *GoPool) Close() {
	p.cancel()
	close(p.taskChan)
	p.wg.Wait()
	close(p.resultChan)
}
