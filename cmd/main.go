package main

import (
	"github.com/gin-gonic/gin"
	_ "go-task-service/cmd/initialize"
	"go-task-service/scheduler"
	"log"
)

func main() {
	// 初始化定时任务调度器
	scheduler.InitScheduler()

	r := gin.Default()

	log.Println("定时任务服务已启动，监听端口 8082")
	r.Run(":8082")
}
