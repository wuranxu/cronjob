package main

import (
	"cronjob/api"
	"cronjob/config"
	"cronjob/database"
	"cronjob/handler"
	"cronjob/tasks"
	"flag"
	"github.com/gin-gonic/gin"
)

var configPath = flag.String("conf", "./config.json", "配置文件")

func main() {
	flag.Parse()
	gin.ForceConsoleColor()
	engine := gin.New()
	api.RegisterRouter(engine)
	// 加载配置
	config.Use(*configPath)
	// 初始化db
	database.Use(config.Conf)
	// 加载定时任务
	tasks.InitTask()
	engine.Use(gin.Logger(), gin.Recovery(), handler.CORSMiddleware())
	engine.Run(":9999")
}
