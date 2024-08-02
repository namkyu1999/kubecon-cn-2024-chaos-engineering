package main

import (
	"context"
	"strconv"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gin-gonic/gin"
)

const (
	pubsubComponentName = "pubsub"
	pubsubTopic         = "orders"
)

func main() {
	r := gin.Default()

	// Create a new client for Dapr using the SDK
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	r.GET("/pub", func(c *gin.Context) {
		contents := ""
		for i := 1; i <= 10; i++ {
			order := `{"orderId":` + strconv.Itoa(i) + `}`

			err := client.PublishEvent(context.Background(), pubsubComponentName, pubsubTopic, []byte(order))
			if err != nil {
				panic(err)
			}

			contents += order + "\n"

			time.Sleep(time.Second)
		}

		c.JSON(200, gin.H{
			"message": contents,
		})
	})

	r.Run()
}
