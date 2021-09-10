package task

type Job struct {
	callback func(job *Job)
	result   interface{}
	err      error
}

func (receiver *Job) Do() {
	receiver.callback(receiver)
}

// 设置任务执行的具体方法逻辑
func (receiver *Job) SetCallback(callback func(job *Job)) *Job {
	receiver.callback = callback
	return receiver
}

// 设置错误信息
func (receiver *Job) SetError(err error) *Job {
	receiver.err = err
	return receiver
}

// 获取错误
func (receiver *Job) GetError() error {
	return receiver.err
}

// 设置结果
func (receiver *Job) SetResult(result interface{}) *Job {
	receiver.result = result
	return receiver
}

func (receiver *Job) GetResult() interface{} {
	return receiver.result
}

// 创建一个新的任务
func NewJob(callback func(job *Job)) *Job {
	job := &Job{
		callback: callback,
		result:   nil,
		err:      nil,
	}
	return job
}
