package tasks

import (
	"cronjob/database"
	"cronjob/models"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
)

var JobPool sync.Map

func InitTask() {
	c := cron.New(cron.WithSeconds())
	AddJobs(c)
	c.Start()
}

func AddJobs(cr *cron.Cron) {
	jobs := fetchJobs()
	for _, job := range jobs {
		task := &Job{
			ID:      job.ID,
			command: job.Command,
		}
		JobPool.Store(fmt.Sprintf("%d", task.ID), task)
		addJob, err := cr.AddJob(job.CronExpr, task)
		if err != nil {
			log.Printf("任务注册失败: %v", err)
		} else {
			log.Println("任务注册成功， ID：", addJob)
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
