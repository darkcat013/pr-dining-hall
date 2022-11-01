package domain

import (
	"math/rand"
	"sync/atomic"

	"github.com/darkcat013/pr-dining-hall/config"
	"github.com/darkcat013/pr-dining-hall/utils"
	"go.uber.org/zap"
)

type TableState int

const (
	Free          TableState = 0
	MakingOrder   TableState = 1
	AwaitingOrder TableState = 2
)

type Table struct {
	Id               int
	ReceiveOrderChan chan Distribution
	State            TableState
}

func NewTable(id int) *Table {
	table := &Table{
		Id:               id,
		ReceiveOrderChan: make(chan Distribution),
		State:            Free,
	}
	go table.Start()
	return table
}

func (t *Table) Start() {
	for {
		if t.State == Free {
			utils.SleepBetween(config.TABLE_OCCUPATION_TIME_MIN, config.TABLE_OCCUPATION_TIME_MAX)
			t.State = MakingOrder
		} else if t.State == MakingOrder {
			utils.SleepBetween(config.TABLE_ORDERING_TIME_MIN, config.TABLE_ORDERING_TIME_MAX)
			t.newOrder()
			t.State = AwaitingOrder
		} else if t.State == AwaitingOrder {
			d := <-t.ReceiveOrderChan
			utils.Log.Info("Table received distribution", zap.Int("tableId", t.Id), zap.Int("orderId", d.OrderId))
			RatingChan <- d
			utils.SleepOneOf(config.TABLE_PICKING_ORDER_TIME, int(d.MaxWait))
			t.State = Free
		}
	}
}

func (t *Table) newOrder() {

	foodsCount := rand.Intn(config.MAX_FOODS) + 1
	var items []int
	var maxPrepTime float64
	probability := float64(100)

	for float64(rand.Intn(100)) <= probability && len(items) < foodsCount {

		randomFood := Menu[rand.Intn(len(Menu))]
		if maxPrepTime < randomFood.PreparationTime {
			maxPrepTime = randomFood.PreparationTime
		}

		items = append(items, randomFood.Id)
		probability /= 1.2
	}
	utils.Log.Info("Start creating order", zap.Int("tableId", t.Id))

	order := Order{
		OrderId: int(atomic.AddInt64(&OrderId, 1)),
		TableId: t.Id,
		Items:   items,
		MaxWait: maxPrepTime * config.MAX_PREP_TIME_COEFF,
	}

	NewOrderChan <- order
}
