package api

import (
	"cronjob/database"
	log "cronjob/logger"
	"cronjob/models"
	"cronjob/tasks"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

const (
	AddJobParamsError  = iota + 10000
	EditJobParamsError = iota + 10000
	AddExistsJob
	JobNotExists
	JobUpdateError
	JobNotRunning
	JobIsRunning
	JobInitError
	PageError
	ListJobError
)

const (
	AddJobParamsErrorInfo  = "添加任务参数错误"
	EditJobParamsErrorInfo = "编辑任务参数错误"
	AddExistsJobInfo       = "任务已存在"
	JobNotExistsInfo       = "任务不存在"
	JobNotRunningInfo      = "任务未开始"
	JobInitErrorInfo       = "定时任务创建失败"
	JobIsRunningInfo       = "定时任务已经开始执行了"
	ListJobErrorInfo       = "定时任务查询失败"
	JobUpdateErrorInfo     = "定时任务更新失败"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type JobDto struct {
	*models.Job
	Status int `json:"status"`
}

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
	task := tasks.New(job.ID, job.CronExpr, job.Command, job.Pause)
	if err = tasks.Update(task); err != nil {
		Failed(c, JobInitError, job, msg(JobInitErrorInfo, err))
		return
	}
	Success(c, job)
}

// 编辑job
func EditJob(c *gin.Context) {
	jobId := c.Param("id")
	id, err := strconv.ParseInt(jobId, 10, 64)
	if err != nil {
		Failed(c, JobNotExists, nil, JobNotExistsInfo)
		return
	}
	job := make(map[string]interface{})
	if err := c.ShouldBindJSON(&job); err != nil {
		Failed(c, EditJobParamsError, nil, msg(EditJobParamsErrorInfo, err))
		return
	}
	oldJob := new(models.Job)
	if err := database.Conn.First(oldJob, `id = ?`, id).Error; err != nil || oldJob.ID != uint(id) {
		Failed(c, JobNotExists, nil, JobNotExistsInfo)
		return
	}
	if _, err := database.Conn.Updates(oldJob, job); err != nil {
		Failed(c, JobUpdateError, nil, msg(JobUpdateErrorInfo, err))
		return
	}
	// 添加到任务池
	task := tasks.New(oldJob.ID, oldJob.CronExpr, oldJob.Command, oldJob.Pause)
	if err := tasks.Update(task); err != nil {
		Failed(c, JobInitError, job, msg(JobInitErrorInfo, err))
		return
	}
	log.Info(tasks.CronInstance.Entries())
	Success(c, oldJob)

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
	go data.RunForWebSocket()
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

// 查看所有job
func ListJob(c *gin.Context) {
	page, size, err := PageUtil(c)
	if err != nil {
		Failed(c, PageError, nil, err)
		return
	}
	jobs := make([]*models.Job, 0, size)
	total, err := database.Conn.FindPagination(page, size, &jobs)
	if err != nil {
		Failed(c, ListJobError, nil, msg(ListJobErrorInfo, err))
		return
	}
	res := make([]*JobDto, 0, total)
	for _, job := range jobs {
		res = append(res, &JobDto{
			Job:    job,
			Status: tasks.Status(fmt.Sprintf("%d", job.ID)),
		})
	}
	Success(c, res)
}

// websocket
func Websocket(c *gin.Context) {
	id := c.Param("id")
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("websocket连接失败, error: %v\n", err)
		return
	}
	defer ws.Close()
MESSAGE:
	for {
		ch := tasks.StreamData.Read(id)
		if ch == nil {
			break MESSAGE
		}
		for {
			select {
			case info := <-ch:
				if info == nil {
					ws.WriteMessage(1, []byte("finished"))
					tasks.StreamData.Close(id)
					return
				}
				ws.WriteMessage(1, []byte(*info))
			}
		}

	}

}
