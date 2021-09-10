package task

import (
	"fmt"
	"testing"
	"time"
)

func TestPool_Run(t *testing.T) {
	pool := NewPool(3, 100)
	for i := 0; i < 10; i++ {
		copyI := i * 10
		job := NewJob(func(job *Job) {
			// 3个协程，10个任务，每个3秒
			// 至少需要12秒
			time.Sleep(time.Second * 3)
			fmt.Printf("%x %x %d\r\n", &copyI, &i, i)
			//设置任务结果
			job.SetResult(copyI)
		})
		pool.Push(job)
	}
	pool.PushDone()
	fmt.Println("hahaha")
	time.Sleep(time.Second * 5)
	for _, job := range pool.GetFinishedJobs() {
		fmt.Printf("计算结果: %d\r\n", job.GetResult())
	}
	pool.Wait()
	fmt.Println("done")
}
