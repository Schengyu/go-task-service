package scheduler

import (
	"github.com/robfig/cron/v3"
	"go-task-service/cmd/global"
	"log"
)

func InitScheduler() {
	// 初始化定时任务调度器
	global.Cron = cron.New(cron.WithSeconds())

	// 每6小时执行一次
	_, err := global.Cron.AddFunc("0 0 */6 * * *", UpdateInventoryTask)
	if err != nil {
		log.Fatalf("注册任务失败: %v", err)
	}

	// 启动定时任务调度器
	global.Cron.Start()
	log.Println("定时任务调度器已启动")
}
