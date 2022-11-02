package domain

import (
	"sync"

	"github.com/darkcat013/pr-dining-hall/config"
)

var Menu []Food

var Tables = make([]*Table, 0, config.TABLES)
var OrderId int64

var Waiters = make([]*Waiter, 0, config.WAITERS)

var RegisteredTime float64

var ReadyClientOrders map[int]*Distribution = make(map[int]*Distribution)

var KitchenOverloadMutex sync.Mutex
var CurrentOrders = 0
var CurrentMaxOrders = config.TABLES
var KitchenOverloaded = false
