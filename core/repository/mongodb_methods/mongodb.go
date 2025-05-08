package mongodb_methods

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"           // 用于构建 MongoDB 查询语法
	"go.mongodb.org/mongo-driver/bson/primitive" // 提供 MongoDB 原生类型，如 ObjectID
	"go.mongodb.org/mongo-driver/mongo"          // MongoDB 官方驱动
)

type YYMSteamInventory struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	SteamAID      string             `bson:"steam_aid"`
	InventoryDesc string             `bson:"inventory_desc"`
}

// MongoRepo 是一个通用的 MongoDB 仓库结构体，使用泛型 T 来支持任意文档类型
type MongoRepo[T any] struct {
	Collection *mongo.Collection // MongoDB 集合对象，指向具体集合
}

// Insert 插入一个文档到集合中，返回插入后的文档 ID
func (r *MongoRepo[T]) Insert(ctx context.Context, doc T) (primitive.ObjectID, error) {
	// 向 MongoDB 插入一条文档
	res, err := r.Collection.InsertOne(ctx, doc)
	if err != nil {
		// 插入失败则返回空的 ObjectID 和错误信息
		return primitive.NilObjectID, err
	}
	// 插入成功后返回插入文档的 ObjectID（需要断言类型）
	return res.InsertedID.(primitive.ObjectID), nil
}

// FindOne 根据给定的查询条件从集合中查找一条文档
func (r *MongoRepo[T]) FindOne(ctx context.Context, filter interface{}) (*T, error) {
	var result T // 声明一个泛型类型的变量，用于接收查找结果
	// 使用 filter 查询符合条件的一条文档
	err := r.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		// 如果查找或解码出错，则返回 nil 和错误信息
		return nil, err
	}
	// 查找成功，返回查询结果的指针
	return &result, nil
}

// UpdateOne 根据 filter 条件更新集合中的一条文档，使用 `$set` 更新字段
func (r *MongoRepo[T]) UpdateOne(ctx context.Context, filter interface{}, update interface{}) error {
	// 调用 UpdateOne 方法，使用 $set 方式更新字段
	_, err := r.Collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	// 返回错误信息（如果有）
	return err
}

// DeleteOne 根据 filter 条件从集合中删除一条文档
func (r *MongoRepo[T]) DeleteOne(ctx context.Context, filter interface{}) error {
	// 调用 DeleteOne 方法，删除匹配条件的文档
	_, err := r.Collection.DeleteOne(ctx, filter)
	// 返回错误信息（如果有）
	return err
}
