// dependencies
package main

import (
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const StateStoreName = "redis-outbox"

type OrderRequest struct {
	OrderID string `json:"order_id"`
}

func main() {
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

		messageID := uuid.NewString()

		ops := make([]*dapr.StateOperation, 0)

		operation := &dapr.StateOperation{
			Type: dapr.StateOperationTypeUpsert,
			Item: &dapr.SetStateItem{
				Key:   messageID,
				Value: []byte(req.OrderID),
			},
		}
		ops = append(ops, operation)
		meta := map[string]string{}
		err = client.ExecuteStateTransaction(c.Request.Context(), StateStoreName, meta, ops)

		c.JSON(200, gin.H{
			"message": "Order completed!",
		})
	})

	if err = r.Run(); err != nil {
		panic(err)
	}
}
