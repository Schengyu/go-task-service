package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/bwmarrin/snowflake"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go-task-service/cmd/global"
	"go-task-service/core/tools"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	//初始化AppConfig配置信息
	InitAppConfig()
	//初始化zap日志
	InitZapLogger()
	//初始化mysql连接
	InitMysql()
	//初始化redis连接
	InitRedis()
	//初始化mongodb连接
	InitMongoDB()
	//初始化elasticsearch连接
	InitElasticsearchClient()
	//初始化rocketmq连接
	InitRocketmqProducer()
	//初始化rocketmq消费者
	InitRocketmqConsumer()
}

// 初始化AppConfig配置信息
func InitAppConfig() {
	// 设置config file path
	viper.SetConfigFile("cmd/appconfig.yaml")
	//通过viper读取nacos的配置信息
	err := viper.ReadInConfig()
	if err != nil {
		panic(err.Error() + " 读取appConfig.yaml配置文件失败")
	}
	//将nacos的配置信息反序列化给global.Nacos(全局变量)
	err = viper.Unmarshal(&global.Nacos)
	if err != nil {
		panic(err.Error() + " 解析appConfig.yaml配置文件中的Nacos信息失败")
	}
	//做个标记 证明nacos已经读取成功
	fmt.Println("Nacos 配置信息读取成功! 开始运行应用。")
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.Nacos.Address, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            global.Nacos.User,
		Password:            global.Nacos.Pass,
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.Nacos.Host,
			Port:   uint64(global.Nacos.Port),
		},
	}
	// 创建动态配置客户端
	Nacosclient, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	//通过nacos客户端去获取config
	config, err := Nacosclient.GetConfig(vo.ConfigParam{
		DataId: global.Nacos.DataId,
		Group:  global.Nacos.Group,
	})
	if err != nil {
		panic(err.Error() + " 获取nacos config 失败")
	}
	var configs map[string]interface{}
	err = json.Unmarshal([]byte(config), &configs)
	if err != nil {
		panic(err.Error() + " 反序列化nacos config 到 AppConfig 失败")
	}
	for k, v := range configs {
		//先判断v是否是字符串
		if str, ok := v.(string); ok {
			//再判断他是否是加密数据
			if str[:3] == "{e}" {
				//对加密项进行解密
				prefix := tools.RemoveEncryptionPrefix(str)
				decryptedValue := tools.DecryptAES(prefix, "base64")
				configs[k] = decryptedValue
			} else {
				configs[k] = v
			}
		} else {
			configs[k] = v
		}
	}
	// 将解密后的配置重新序列化并更新到 global.AppConfigMaster
	configBytes, err := json.Marshal(configs)
	if err != nil {
		panic(err.Error() + " 重新序列化解密后的配置失败")
	}
	err = json.Unmarshal(configBytes, &global.AppConfigMaster)
	fmt.Println("AppConfig 读取成功!", global.AppConfigMaster.DBUSER)
}

// InitZapLogger 初始化 zap 日志系统
func InitZapLogger() {
	// 日志配置
	cfg := &tools.LoggerConfig{
		Mode:       "dev",
		LogDir:     "./core/logs\\",
		MaxSize:    100,
		MaxBackups: 100000,
		MaxAge:     90,
		Compress:   true,
	}
	// 初始化日志组件
	writeSyncer := tools.GetLogWriter(cfg)
	encoder := tools.GetEncoder(cfg.Mode)

	// 设置日志级别
	level := zapcore.InfoLevel
	if cfg.Mode == "dev" {
		level = zapcore.DebugLevel
	}

	// 创建核心 Core
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writeSyncer),
		level,
	)

	// 初始化全局 Logger
	global.ZapLog = zap.New(
		core,
		zap.AddCaller(),      // 打印调用信息
		zap.AddCallerSkip(1), // 跳过一层调用栈，定位更准确
	)

	// 替换全局 zap 实例
	zap.ReplaceGlobals(global.ZapLog)

	fmt.Println("✅ Zap 彩色日志系统初始化完成!")
}

//初始化连接mysql

func InitMysql() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", global.AppConfigMaster.DBUSER, global.AppConfigMaster.DBPASSWORD, global.AppConfigMaster.DBHOST, global.AppConfigMaster.DBPORT, global.AppConfigMaster.DBNAME)
	fmt.Println(dsn)
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("MySQL数据库连接失败", err)
		zap.L().Error("MySQL数据库连接异常失败" + err.Error())
		panic("MySQL数据库连接失败" + err.Error())
	}
	fmt.Println("MySQL数据库连接成功")
	zap.L().Info("MySQL数据库连接成功")
}

//初始化redis连接

func InitRedis() {
	//配置Redis连接信息
	global.RedisDB = redis.NewClient(&redis.Options{
		Addr:     global.AppConfigMaster.REDISHOST + ":" + strconv.Itoa(global.AppConfigMaster.REDISPORT),
		Password: global.AppConfigMaster.REDISPASSWORD, // no password set
		DB:       0,                                    // use default DB
	})
	//连接Redis数据库
	pong, err := global.RedisDB.Ping().Result()
	if err != nil {
		fmt.Println("Redis数据库连接失败", err)
		zap.L().Error("Redis数据库连接异常失败" + err.Error())
		panic("Redis数据库连接失败" + err.Error())
	}
	fmt.Println(pong, err)
	fmt.Println("Redis数据库连接成功")
	zap.L().Info("MySQL数据库连接成功")
}

// InitSnowflake 初始化 Snowflake 节点
func InitSnowflake(nodeID int64) {
	global.SnowflakeOnce.Do(func() {
		var err error
		global.SnowflakeNode, err = snowflake.NewNode(nodeID)
		if err != nil {
			fmt.Println("初始化 Snowflake 节点失败", err)
			zap.L().Error("初始化 Snowflake 节点失败" + err.Error())
			return
		}
	})
}

// InitElasticsearchClient 初始化 Elasticsearch 客户端
func InitElasticsearchClient() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			global.AppConfigMaster.ESNODE, // 你可以放多个节点
		},
		Username: "elastic",       // 如果开启了认证
		Password: "yuanyoumaoaaa", // 替换为实际密码
	}
	var err error
	global.ESClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println("创建Elasticsearch客户端失败", err)
		zap.L().Error("创建Elasticsearch客户端失败" + err.Error())
		return
	}
	// 测试连接
	res, err := global.ESClient.Info()
	if err != nil {
		fmt.Println("获取响应失败", err)
		zap.L().Info("获取响应失败" + err.Error())
		return
	}
	defer res.Body.Close()
	fmt.Println("Elasticsearch客户端连接成功")
	zap.L().Info("Elasticsearch客户端连接成功")
}

// 初始化rocketmq生产者
func InitRocketmqProducer() {
	var err error
	global.RocketMQProducer, err = rocketmq.NewProducer(producer.WithNameServer([]string{"118.89.80.34:9876"}), producer.WithGroupName("SCY"))
	if err != nil {
		fmt.Println("初始化失败", err)
		zap.L().Error("初始化rocketmq生产者失败" + err.Error())
		return
	}
	fmt.Println("初始化rocketmq生产者成功")
	zap.L().Info("初始化rocketmq生产者成功")
	err = global.RocketMQProducer.Start()
	if err != nil {
		fmt.Println("生产者启动失败", err)
		zap.S().Error("生产者启动失败", err)
		return
	}
	fmt.Println("rocketmq生产者启动成功")
	zap.L().Info("rocketmq生产者启动成功")

}

// 初始化rocketmq消费者
func InitRocketmqConsumer() {
	var err error
	global.RocketMQConsumer, err = rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"118.89.80.34:9876"}), consumer.WithGroupName("SCY"))
	if err != nil {
		fmt.Println("初始化失败")
		zap.L().Error("初始化rocketmq消费者失败" + err.Error())
		return
	}
	fmt.Println("初始化rocketmq消费者成功")
	zap.L().Info("初始化rocketmq消费者成功")
}

// InitMongoDB 初始化 MongoDB 连接，并赋值给全局变量 global.MongoDB
func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI("mongodb://mongo:mongo_scy@118.89.80.34:27017")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("连接 MongoDB 失败: %v", err)
	}

	// 检查连接
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB Ping 失败: %v", err)
	}

	fmt.Println("MongoDB 连接成功")
	global.MongoDB = client.Database("test")
}
