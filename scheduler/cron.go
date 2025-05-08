package scheduler

import (
	"github.com/robfig/cron/v3"
	"go-task-service/cmd/global"
	"log"
)

func InitScheduler() {
	// 初始化定时任务调度器
	global.Cron = cron.New(cron.WithSeconds())

	// 注册定时任务
	registerTasks()

	// 启动定时任务调度器
	global.Cron.Start()
	log.Println("定时任务调度器已启动")
}
