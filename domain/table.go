package domain

import (
	"math/rand"
	"sync/atomic"

	"github.com/darkcat013/pr-dining-hall/config"
	configGlobal "github.com/darkcat013/pr-dining-hall/config-global"
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

func InitializeTables() {
	if configGlobal.TABLES_ENABLED {
		for i := 0; i < config.TABLES; i++ {
			table := NewTable(i)
			Tables = append(Tables, table)
			utils.SleepBetween(config.TABLE_OCCUPATION_TIME_MIN, config.TABLE_OCCUPATION_TIME_MAX)
		}
	}
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
			KitchenOverloadMutex.Lock()
			if CurrentOrders == 0 {
				CurrentMaxOrders++
				CurrentOrders++
				t.State = MakingOrder
			} else if CurrentOrders < CurrentMaxOrders {
				t.State = MakingOrder
				CurrentOrders++
			}
			KitchenOverloadMutex.Unlock()
		} else if t.State == MakingOrder {
			utils.SleepBetween(config.TABLE_ORDERING_TIME_MIN, config.TABLE_ORDERING_TIME_MAX)
			t.newOrder()
			t.State = AwaitingOrder
		} else if t.State == AwaitingOrder {
			d := <-t.ReceiveOrderChan
			utils.Log.Info("Table received distribution", zap.Int("tableId", t.Id), zap.Int("orderId", d.OrderId))
			RatingChan <- d
			KitchenOverloadMutex.Lock()
			CurrentOrders--
			CurrentFoodAmount -= len(d.Items)
			KitchenOverloadMutex.Unlock()
			utils.SleepOneOf(config.TABLE_PICKING_ORDER_TIME, int(d.MaxWait))
			t.State = Free
		}
	}
}

func (t *Table) newOrder() {

	foodsCount := rand.Intn(config.MAX_FOODS) + 1
	var items []int
	var maxPrepTime float64

	KitchenOverloadMutex.Lock()

	for i := 0; i < foodsCount; i++ {
		randomFood := Menu[rand.Intn(len(Menu))]
		if randomFood.CookingApparatus != "" && rand.Intn(4) > 0 {
			i--
			continue
		}
		if maxPrepTime < randomFood.PreparationTime {
			maxPrepTime = randomFood.PreparationTime
		}

		items = append(items, randomFood.Id)
	}
	if len(items) >= config.MAX_FOODS/2 && CurrentMaxOrders > 2 {
		CurrentMaxOrders--
	}
	CurrentFoodAmount += len(items)
	KitchenOverloadMutex.Unlock()
	utils.Log.Info("Start creating order", zap.Int("tableId", t.Id))

	order := Order{
		OrderId: int(atomic.AddInt64(&OrderId, 1)),
		TableId: t.Id,
		Items:   items,
		MaxWait: maxPrepTime * config.MAX_PREP_TIME_COEFF,
	}

	NewOrderChan <- order
}
