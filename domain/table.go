package domain

import (
	"math/rand"
	"sync/atomic"

	"github.com/darkcat013/pr-dining-hall/config"
	"github.com/darkcat013/pr-dining-hall/utils"
	"go.uber.org/zap"
)

var orderId int64
var Tables = make([]*Table, 0, config.TABLES)

type Table struct {
	Id               int
	ReceiveOrderChan chan Distribution
	IsFree           bool
}

func NewTable(id int) *Table {
	return &Table{
		Id:               id,
		ReceiveOrderChan: make(chan Distribution),
		IsFree:           true,
	}
}

func (t *Table) Start() {
	for {
		if t.IsFree {
			t.newOrder()
		} else {
			d := <-t.ReceiveOrderChan
			utils.Log.Info("Received distribution", zap.Int("tableId", t.Id), zap.Int("orderId", d.OrderId))
			//calculate rating
			t.IsFree = true
		}
	}
}

func (t *Table) newOrder() {

	t.IsFree = false

	foodsCount := rand.Intn(10) + 1
	var items []int
	var maxPrepTime int

	for i := 0; i < foodsCount; i++ {
		randomFood := Menu[rand.Intn(len(Menu))]
		if maxPrepTime < randomFood.PreparationTime {
			maxPrepTime = randomFood.PreparationTime
		}

		items = append(items, randomFood.Id)
	}
	utils.Log.Info("Start creating order", zap.Int("tableId", t.Id))

	order := Order{
		OrderId: int(atomic.AddInt64(&orderId, 1)),
		TableId: t.Id,
		Items:   items,
		MaxWait: float64(maxPrepTime) * 1.3,
	}

	NewOrderChan <- order
}
