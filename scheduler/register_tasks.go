package scheduler

import (
	"go-task-service/cmd/global"
	"log"
)

func registerTasks() {
	tasks := []struct {
		Spec string
		Job  func()
	}{
		{"0 0 */6 * * *", UpdateInventoryTask},
	}

	for _, task := range tasks {
		_, err := global.Cron.AddFunc(task.Spec, task.Job)
		if err != nil {
			log.Fatalf("注册任务失败: %v", err)
		}
	}
}
