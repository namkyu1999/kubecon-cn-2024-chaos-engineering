package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var sub = &common.Subscription{
	PubsubName: "pubsub",
	Topic:      "orders1",
	Route:      "/orders1",
}

const AppPort = "6005"
const MongoDBURI = "mongodb://mongo-mongodb.common.svc.cluster.local:27017/delivery"
const MongoDBDatabase = "delivery"
const MongoDBCollection = "delivery1"

func main() {
	// Mongo client
	client, err := MongoConnection()
	db := client.Database(MongoDBDatabase)
	log.Println("Connected to MongoDB")

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	err = CreateCollection(MongoDBCollection, db)
	if err != nil {
		log.Fatalf("failed to create collection: %v", err)
	}

	// Create the new server on appPort and add a topic listener
	s := daprd.NewService(":" + AppPort)

	err = s.AddServiceInvocationHandler("/result", func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		log.Info("Received request for result")

		res, err := db.Collection(MongoDBCollection).CountDocuments(ctx, bson.M{"isDelivered": false})
		if err != nil {
			return nil, err
		}

		result := strconv.FormatInt(res, 10)

		out = &common.Content{
			Data:        []byte(result),
			ContentType: in.ContentType,
			DataTypeURL: in.DataTypeURL,
		}
		return
	})
	if err != nil {
		log.Fatalf("error adding invocation handler: %v", err)
	}

	err = s.AddTopicEventHandler(sub, func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		orderID := fmt.Sprint(e.Data)
		log.Println("Received orderID: " + orderID)
		filter := bson.D{{"orderID", orderID}}
		updates := bson.D{{"$set", bson.D{{"isDelivered", true}}}}

		if _, err := db.Collection(MongoDBCollection).UpdateOne(ctx, filter, updates); err != nil {
			return false, err
		}

		return false, nil
	})
	if err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	// Start the server
	err = s.Start()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error listenning: %v", err)
	}
}

func MongoConnection() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDBURI))
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
