package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	CODE = iota
)

const (
	DEFAULT_PAGE      = 1
	DEFAULT_PAGE_SIZE = 8
)

var (
	PageNotValid     = errors.New("page参数非法")
	PageSizeNotValid = errors.New("size参数非法")
)

func msg(desc string, err error) string {
	return fmt.Sprintf("%s: %v", desc, err)
}

// success response
func Success(c *gin.Context, data interface{}, msg ...interface{}) {
	if len(msg) == 0 || msg[0] == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": CODE, "msg": "操作成功", "data": data,
		})
		return
	}
	switch msg[0].(type) {
	case string:
		c.JSON(http.StatusOK, gin.H{
			"code": CODE, "msg": msg, "data": data,
		})
	case error:
		c.JSON(http.StatusOK, gin.H{
			"code": CODE, "msg": msg[0].(error).Error(), "data": data,
		})
	default:
		c.JSON(http.StatusOK, gin.H{
			"code": CODE, "msg": fmt.Sprintf("%v", msg[0]), "data": data,
		})
	}
}

func Failed(c *gin.Context, code int, data interface{}, msg interface{}) {
	switch msg.(type) {
	case string:
		c.JSON(http.StatusOK, gin.H{
			"code": code, "msg": msg, "data": data,
		})
	case error:
		c.JSON(http.StatusOK, gin.H{
			"code": code, "msg": msg.(error).Error(), "data": data,
		})
	default:
		c.JSON(http.StatusOK, gin.H{
			"code": code, "msg": fmt.Sprintf("%v", msg), "data": data,
		})
	}
}

func RegisterRouter(engine *gin.Engine) {
	// job api
	jobRouter := engine.Group("/job")
	{
		jobRouter.POST("/", AddJob)
		jobRouter.PUT("/:id", EditJob)
		jobRouter.GET("/start/:id", StartJob)
		jobRouter.GET("/stop/:id", StopJob)
		jobRouter.GET("/list", ListJob)
		jobRouter.GET("/log/:id", Websocket)
	}
}

func PageUtil(c *gin.Context) (int, int, error) {
	page := c.Query("page")
	size := c.Query("size")
	current, err := strconv.ParseInt(page, 10, 32)
	if err == nil && current < 0 {
		return DEFAULT_PAGE, DEFAULT_PAGE_SIZE, PageNotValid
	}
	pageSize, err := strconv.ParseInt(size, 10, 32)
	if err == nil && pageSize > 500 || pageSize < 0 {
		return DEFAULT_PAGE, DEFAULT_PAGE_SIZE, PageSizeNotValid
	}
	if current == 0 {
		current = DEFAULT_PAGE
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	return int(current), int(pageSize), nil
}
