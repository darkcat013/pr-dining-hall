package utils

import (
	"math/rand"
	"time"

	"github.com/darkcat013/pr-dining-hall/config"
)

func SleepBetween(min, max int) {
	time.Sleep(time.Duration(GetRandomInt(min, max)) * config.TIME_UNIT)
}

func SleepOneOf(params ...int) {
	time.Sleep(time.Duration(params[rand.Intn(len(params))]) * config.TIME_UNIT)
}

func GetCurrentTimeFloat() float64 {
	if config.TIME_UNIT >= time.Millisecond && config.TIME_UNIT < time.Second {
		return float64(time.Now().UnixMilli())
	} else {
		return float64(time.Now().Unix())
	}
}
