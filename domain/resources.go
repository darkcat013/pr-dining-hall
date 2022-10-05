package domain

import "github.com/darkcat013/pr-dining-hall/config"

var Menu []*Food

var Tables = make([]*Table, 0, config.TABLES)
var OrderId int64

var Waiters = make([]*Waiter, 0, config.WAITERS)
