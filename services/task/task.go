package task

import (
	"sync"
)

type Pool struct {
	// 启用的协程个数
	quantity int
	// 推送任务的channel
	jobs chan *Job
	// 协程waitGroup
	wg sync.WaitGroup
	// 存放完成的任务
	finished []*Job
	// 更新finished的读写锁
	rwMutex sync.RWMutex
	// 标记是否正在运行
	isRun bool
}

// 创建一个新的任务池
func NewPool(quantity int, queueSize int) *Pool {
	pool := &Pool{
		quantity: quantity,
		jobs:     make(chan *Job, queueSize),
		finished: make([]*Job, 0),
	}
	pool.run()
	return pool
}

// 推送一个任务
func (receiver *Pool) Push(job *Job) {
	receiver.jobs <- job
}

// 完成推送
func (receiver *Pool) PushDone() {
	close(receiver.jobs)
}

// 等待
func (receiver *Pool) Wait() {
	receiver.wg.Wait()
}

// 异步运行
func (receiver *Pool) run() {
	if receiver.isRun {
		return
	}
	receiver.isRun = true
	for i := 0; i < receiver.quantity; i++ {
		receiver.wg.Add(1)
		go func(id int) {
			defer receiver.wg.Done()
			for {
				job, ok := <-receiver.jobs
				if !ok {
					break
				}
				job.Do()
				receiver.moveFinished(job)
			}
			return
		}(i)
	}
}

// 将任务移动至完成列表
func (receiver *Pool) moveFinished(job *Job) {
	receiver.rwMutex.Lock()
	defer receiver.rwMutex.Unlock()
	receiver.finished = append(receiver.finished, job)
}

// 获取已完成的任务
func (receiver *Pool) GetFinishedJobs() []*Job {
	return receiver.finished
}
