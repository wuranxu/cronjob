package main

import (
	"cronjob/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.ForceConsoleColor()
	engine := gin.New()

	engine.Use(gin.Logger(), gin.Recovery(), handler.CORSMiddleware())
	engine.Run(":9999")
}
