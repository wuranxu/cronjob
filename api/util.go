package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CODE = iota
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
		jobRouter.POST("/insert", AddJob)
		jobRouter.GET("/stop/:id", StopJob)
	}
}
