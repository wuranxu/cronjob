package api

import (
	"cronjob/database"
	"cronjob/models"
	"cronjob/tasks"
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	AddJobParamsError = iota + 10000
	AddExistsJob
	JobNotExists
	JobNotRunning
	JobIsRunning
	JobInitError
)

const (
	AddJobParamsErrorInfo = "添加任务参数错误"
	AddExistsJobInfo      = "任务已存在"
	JobNotExistsInfo      = "任务不存在"
	JobNotRunningInfo     = "任务未开始"
	JobInitErrorInfo      = "定时任务创建失败"
	JobIsRunningInfo      = "定时任务已经开始执行了"
)

// 新增job
func AddJob(c *gin.Context) {
	job := new(models.Job)
	if err := c.ShouldBindJSON(job); err != nil {
		Failed(c, AddJobParamsError, nil, msg(AddJobParamsErrorInfo, err))
		return
	}
	err := database.Conn.Insert(job)
	if err != nil {
		Failed(c, AddExistsJob, nil, msg(AddExistsJobInfo, err))
		return
	}
	// 添加到任务池
	task := tasks.New(job.ID, job.Command)
	tasks.JobPool.Store(fmt.Sprintf("%d", job.ID), task)
	if _, err := tasks.CronInstance.AddJob(job.CronExpr, task); err != nil {
		Failed(c, JobInitError, job, msg(JobInitErrorInfo, err))
		return
	}
	Success(c, job)
}

// 开始job
func StartJob(c *gin.Context) {
	jobId := c.Param("id")
	load, ok := tasks.JobPool.Load(jobId)
	if !ok {
		Failed(c, JobNotExists, nil, JobNotExistsInfo)
		return
	}
	data := load.(*tasks.Job)
	if data.Running() {
		Failed(c, JobIsRunning, nil, JobIsRunningInfo)
		return
	}
	go data.Run()
	Success(c, nil)
}

// 停止job
func StopJob(c *gin.Context) {
	param := c.Param("id")
	load, ok := tasks.JobPool.Load(param)
	if !ok {
		Failed(c, JobNotExists, nil, JobNotExistsInfo)
		return
	}
	data := load.(*tasks.Job)
	if !data.Stop() {
		Failed(c, JobNotRunning, nil, JobNotRunningInfo)
		return
	}
	Success(c, nil)
}
