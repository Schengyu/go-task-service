package scheduler

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go-task-service/cmd/global"
	"go-task-service/core/curd_methods"
	"go.uber.org/zap"
	"log"
	"strconv"
)

func UpdateInventoryTask() {
	log.Println("[定时任务] 开始执行 UpdateInventoryTask")

	//第一步查询六小时没有更新的库存信息 然后把他们的 steam_aid 发送到消息队列
	inventoryList, err := curd_methods.QueryOutdatedInventory()
	if err != nil {
		log.Println("查询过期库存失败:", err)
		zap.L().Error("查询过期库存失败", zap.Error(err))
		return
	}
	fmt.Println("查询到的过期库存信息:", inventoryList)
	//查询到之后将他们的 steam_aid 发送到消息队列
	for _, v := range inventoryList {
		//调用全局生产者向队列当中发送一条信息
		_, err = global.RocketMQProducer.SendSync(context.Background(), &primitive.Message{
			Topic: "inventory_desc",
			Body:  []byte(strconv.Itoa(int(v.SteamAID))),
		})
		if err != nil {
			fmt.Println("发送消息失败", err)
			zap.S().Error("发送消息失败", err)
			return
		}
		fmt.Println("消息发送成功")
	}
	log.Println("[定时任务] UpdateInventoryTask 执行完成")
}
