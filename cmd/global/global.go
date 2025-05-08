package global

import (
	"encoding/base64"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/bwmarrin/snowflake"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/robfig/cron/v3"
	"go-task-service/cmd/appconf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

var (
	Nacos            *appconf.Nacos
	AppConfig        *appconf.AppConfig
	Router           *gin.Engine
	DB               *gorm.DB
	RedisDB          *redis.Client
	AppConfigMaster  *appconf.AppConfigMaster
	AESKey, _        = base64.StdEncoding.DecodeString("hRgcXGXelyYzRPvMVwHfJfJ0pj+2mhJoH0QYcOGlrcY=")
	AESIv            = []byte("0000000000000000") // 16 字节 IV
	SnowflakeNode    *snowflake.Node
	SnowflakeOnce    sync.Once
	ESClient         *elasticsearch.Client
	Once             sync.Once
	RocketMQProducer rocketmq.Producer
	RocketMQConsumer rocketmq.PushConsumer
	ZapLog           *zap.Logger
	MongoDB          *mongo.Database
	Cron             *cron.Cron
)
