package database

import (
    "context"
    "sync"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    clientInstance     *mongo.Client
    clientInstanceErr  error
    mongoOnce          sync.Once
    connectionTimeout  = 10 * time.Second           // 连接超时时间
)

// GetMongoClient 返回 MongoDB 的 client 实例
func GetMongoClient(mongoURI string ) (*mongo.Client, error) {
    // 使用 sync.Once 确保只初始化一次
    mongoOnce.Do(func() {
        clientOptions := options.Client().ApplyURI(mongoURI)
        ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
        defer cancel()

        clientInstance, clientInstanceErr = mongo.Connect(ctx, clientOptions)
        if clientInstanceErr != nil {
            // log.Fatal("Failed to connect to MongoDB:", clientInstanceErr)
        }

        // 检查连接是否成功
        if err := clientInstance.Ping(ctx, nil); err != nil {
            clientInstanceErr = err
            // log.Fatal("Could not ping MongoDB:", err)
        }
    })

    return clientInstance, clientInstanceErr
}
