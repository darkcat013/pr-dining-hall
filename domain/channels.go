package domain

import "github.com/darkcat013/pr-dining-hall/config"

var NewOrderChan = make(chan Order, config.TABLES)
var RatingChan = make(chan Distribution)
var ClientOrderChan = make(chan Order, config.WAITERS)
var ClientRatingChan = make(chan int)
