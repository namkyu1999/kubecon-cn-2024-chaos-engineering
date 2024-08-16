// dependencies
package main

import (
	"context"
	"strings"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MongoDBURI = "mongodb://my-release-mongodb-0.my-release-mongodb-headless.litmus:27017,my-release-mongodb-1.my-release-mongodb-headless.litmus:27017,my-release-mongodb-2.my-release-mongodb-headless.litmus:27017/admin"
const MongoDBDatabase = "delivery"
const MongoDBCollection = "delivery1"
const DBUser = "root"
const DBPassword = "1234"
const PubSubName = "pubsub1"
const TopicName = "orders1"

type Delivery struct {
	OrderID     string `bson:"orderID"`
	IsDelivered bool   `bson:"isDelivered"`
}

type OrderRequest struct {
	OrderID string `json:"order_id"`
}

func main() {
	// Mongo client
	mongoClient, err := MongoConnection()
	db := mongoClient.Database(MongoDBDatabase)
	log.Println("Connected to MongoDB")

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	err = CreateCollection(MongoDBCollection, db)
	if err != nil {
		log.Fatalf("failed to create collection: %v", err)
	}

	r := gin.Default()

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	r.POST("/order", func(c *gin.Context) {
		var req OrderRequest
		err := c.BindJSON(&req)
		if err != nil {
			panic(err)
		}

		newData := Delivery{
			OrderID:     req.OrderID,
			IsDelivered: false,
		}

		if _, err := db.Collection(MongoDBCollection).InsertOne(c.Request.Context(), newData); err != nil {
			log.Error("Error inserting data into MongoDB: " + err.Error())
			c.JSON(500, gin.H{
				"message": "Error inserting data into MongoDB",
			})
		}

		if err := client.PublishEvent(c.Request.Context(), PubSubName, TopicName, []byte(req.OrderID)); err != nil {
			c.JSON(500, gin.H{
				"message": "Error publishing event",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Order completed!",
			})
		}
	})

	if err = r.Run(); err != nil {
		panic(err)
	}
}

func MongoConnection() (*mongo.Client, error) {
	credential := options.Credential{
		Username: DBUser,
		Password: DBPassword,
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDBURI).SetAuth(credential))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CreateCollection(collectionName string, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.CreateCollection(ctx, collectionName)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Info(collectionName + "'s collection already exists, continuing with the existing mongo collection")
			return nil
		} else {
			return err
		}
	}

	log.Info(collectionName + "'s mongo collection created")
	return nil
}
