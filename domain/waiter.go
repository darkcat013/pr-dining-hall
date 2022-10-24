package domain

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"sort"

	"github.com/darkcat013/pr-dining-hall/config"
	"github.com/darkcat013/pr-dining-hall/utils"
	"go.uber.org/zap"
)

type Waiter struct {
	Id                      int
	ReceiveDistributionChan chan Distribution
}

func NewWaiter(id int) *Waiter {
	waiter := &Waiter{
		Id:                      id,
		ReceiveDistributionChan: make(chan Distribution),
	}
	go waiter.Start()
	return waiter
}

func (w *Waiter) Start() {
	for {
		select {
		case order := <-NewOrderChan:
			utils.SleepBetween(config.WAITER_TAKING_ORDER_TIME_MIN, config.WAITER_TAKING_ORDER_TIME_MAX)

			order.Items = getSortedFoodsByPrepTime(order.Items)
			order.WaiterId = w.Id
			order.Priority = getOrderPriority(&order)
			order.PickUpTime = utils.GetCurrentTimeFloat()

			utils.Log.Info("Created order", zap.Any("order", order))

			go sendOrder(&order)

		case distribution := <-w.ReceiveDistributionChan:
			utils.Log.Info("Waiter received distribution", zap.Any("distribution", distribution))
			Tables[distribution.TableId].ReceiveOrderChan <- distribution
		}
	}
}

func getSortedFoodsByPrepTime(items []int) []int {
	foods := make([]Food, len(items))
	for i := 0; i < len(items); i++ {
		foods[i] = Menu[items[i]-1]
	}

	sort.Slice(foods, func(i, j int) bool {
		return foods[i].PreparationTime > foods[j].PreparationTime
	})

	sortedItems := make([]int, len(foods))
	for i := 0; i < len(foods); i++ {
		sortedItems[i] = foods[i].Id
	}

	return sortedItems
}

func getOrderPriority(order *Order) int {

	foodAmount := float64(len(order.Items))
	firstPriority := getPriorityBy(config.MAX_FOODS, foodAmount)

	if firstPriority >= 4 {
		return firstPriority
	}

	var sumPrepTime float64
	for i := 0; i < int(foodAmount); i++ {
		sumPrepTime += Menu[order.Items[i]-1].PreparationTime
	}
	avg := sumPrepTime / foodAmount
	maxPrepTime := order.MaxWait / config.MAX_PREP_TIME_COEFF
	secondPriority := getPriorityBy(maxPrepTime, avg)

	return int((firstPriority + secondPriority) / 2)
}

func getPriorityBy(maxValue, value float64) int {
	step := math.Ceil(maxValue / 5.0)
	if value <= step {
		return 5
	} else if value <= step*2 {
		return 4
	} else if value <= step*3 {
		return 3
	} else if value <= step*4 {
		return 2
	} else {
		return 1
	}
}

func sendOrder(order *Order) {
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
}
