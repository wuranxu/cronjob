package tasks

import (
	"cronjob/database"
	log "cronjob/logger"
	"cronjob/models"
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
)

var (
	JobPool      sync.Map
	CronInstance *cron.Cron
)

func InitTask() {
	CronInstance = cron.New(cron.WithSeconds())
	AddJobs(CronInstance)
	CronInstance.Start()
}

func Status(jobId string) int {
	load, ok := JobPool.Load(jobId)
	if !ok {
		return 0
	}
	job := load.(*Job)
	if job.running {
		return 1
	}
	return 2
}

func AddJobs(cr *cron.Cron) {
	jobs := fetchJobs()
	for _, job := range jobs {
		task := &Job{
			ID:      job.ID,
			command: job.Command,
			pause:   job.Pause,
		}
		entry, err := cr.AddJob(job.CronExpr, task)
		if err != nil {
			log.Errorf("任务注册失败: %v\n", err)
		} else {
			log.Info("任务注册成功， ID：", entry)
			task.entry = entry
			JobPool.Store(fmt.Sprintf("%d", task.ID), task)
		}
	}
}

func fetchJobs() []*models.Job {
	jobs := make([]*models.Job, 0, 20)
	cursor := database.Conn.Find(&jobs)
	if cursor.Error != nil {
		panic("获取任务列表失败")
	}
	return jobs
}
