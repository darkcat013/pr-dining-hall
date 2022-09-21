package domain

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/darkcat013/pr-dining-hall/config"
	"github.com/darkcat013/pr-dining-hall/utils"
	"go.uber.org/zap"
)

var Waiters = make([]*Waiter, 0, config.WAITERS)

type Waiter struct {
	Id                      int
	ReceiveDistributionChan chan Distribution
}

func NewWaiter(id int) *Waiter {
	return &Waiter{
		Id:                      id,
		ReceiveDistributionChan: make(chan Distribution),
	}
}

func (w *Waiter) Start() {
	for {
		select {
		case order := <-NewOrderChan:
			order.WaiterId = w.Id
			order.PickUpTime = int(time.Now().Unix())
			order.Priority = rand.Intn(5) + 1

			utils.Log.Info("Created order", zap.Any("order", order))

			body, err := json.Marshal(order)
			if err != nil {
				utils.Log.Fatal("Failed to convert order to JSON ", zap.String("error", err.Error()), zap.Any("order", order))
			}

			utils.Log.Info("Send order to kitchen", zap.Any("order", order))

			resp, err := http.Post(config.KITCHEN_URL, "application/json", bytes.NewBuffer(body))

			if err != nil {
				utils.Log.Fatal("Failed to send order to kitchen", zap.String("error", err.Error()), zap.Int("orderId", order.OrderId))
			} else {
				utils.Log.Info("Response from kitchen", zap.Int("statusCode", resp.StatusCode), zap.Int("orderId", order.OrderId))
			}
		case distribution := <-w.ReceiveDistributionChan:
			utils.Log.Info("Waiter received distribution", zap.Any("distribution", distribution))
			Tables[distribution.TableId].ReceiveOrderChan <- distribution
		}
	}
}
