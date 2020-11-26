package api

import (
	"cronjob/database"
	"cronjob/models"
	"cronjob/tasks"
	"github.com/gin-gonic/gin"
)

const (
	AddJobParamsError = iota + 10000
	AddExistsJob
	JobNotExists
	JobNotRunning
)

const (
	AddJobParamsErrorInfo = "添加任务参数错误"
	AddExistsJobInfo      = "任务已存在"
	JobNotExistsInfo      = "任务不存在"
	JobNotRunningInfo     = "任务未开始"
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
	Success(c, job)
}

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
