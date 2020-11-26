package tasks

import (
	"context"
	"log"
)

type Job struct {
	ID      uint `json:"id"`
	cancel  context.CancelFunc
	command string
}

func New(id uint, command string) *Job {
	return &Job{
		ID: id, command: command,
	}
}

func (j *Job) Run() {
	log.Println("任务开始")
	ctx, cancelFunc := context.WithCancel(context.Background())
	j.cancel = cancelFunc
	Command(ctx, j.command)
	j.cancel = nil
}

func (j *Job) Stop() bool {
	if j.cancel == nil {
		return false
	}
	j.cancel()
	j.cancel = nil
	return true
}

func (j *Job) Running() bool {
	return j.cancel != nil
}
