package tasks

import (
	"context"
	log "cronjob/logger"
	"fmt"
	"github.com/robfig/cron/v3"
)

type Job struct {
	ID      uint `json:"id"`
	cancel  context.CancelFunc
	command string
	running bool
	pause   bool
	entry   cron.EntryID
	Spec    string
}

func New(id uint, spec, command string, pause bool) *Job {
	return &Job{
		ID: id, command: command, Spec: spec, pause: pause,
	}
}

func Stop(jobId string) {
	load, ok := JobPool.Load(jobId)
	if ok {
		j := load.(*Job)
		j.Stop()
	}
}

func Update(job *Job) error {
	load, ok := JobPool.Load(fmt.Sprintf("%d", job.ID))
	if ok {
		j := load.(*Job)
		if j.running {
			j.Stop()
		}
		CronInstance.Remove(j.entry)
	}
	entry, err := CronInstance.AddJob(job.Spec, job)
	if err != nil {
		return err
	}
	job.entry = entry
	JobPool.Store(fmt.Sprintf("%d", job.ID), job)
	return nil
}

func (j *Job) Run() {
	if j.running {
		log.Info("任务仍在进行，不继续操作")
		return
	}
	if j.pause {
		log.Info("任务已被暂停")
		return
	}
	j.running = true
	ctx, cancelFunc := context.WithCancel(context.Background())
	j.cancel = cancelFunc
	Command(ctx, j.ID, j.command, false)
	j.cancel = nil
	j.running = false
}

func (j *Job) RunForWebSocket() {
	if j.pause {
		log.Info("任务已被暂停")
		return
	}
	if j.running {
		log.Info("任务仍在进行，不继续操作")
		return
	}
	j.running = true
	ctx, cancelFunc := context.WithCancel(context.Background())
	j.cancel = cancelFunc
	Command(ctx, j.ID, j.command, true)
	j.cancel = nil
	j.running = false
}

func (j *Job) Stop() bool {
	if !j.running {
		j.cancel = nil
		return false
	}
	if j.cancel != nil {
		j.cancel()
	}
	j.cancel = nil
	j.running = false
	return true
}

func (j *Job) Running() bool {
	return j.running
}
