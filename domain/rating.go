package domain

import (
	"sync/atomic"

	"github.com/darkcat013/pr-dining-hall/config"
	"github.com/darkcat013/pr-dining-hall/utils"
	"go.uber.org/zap"
)

var CompletedOrders int64
var AvgRating float64

func StartRatingLogging() {
	var summedRatings int64

	for {
		d := <-RatingChan
		atomic.AddInt64(&CompletedOrders, 1)
		waitTime := (utils.GetCurrentTimeFloat() - d.PickUpTime) / config.TIME_UNIT_COEFF

		currRating := int64(0)

		if waitTime < d.MaxWait {
			currRating = 5
		} else if waitTime < d.MaxWait*1.1 {
			currRating = 4
		} else if waitTime < d.MaxWait*1.2 {
			currRating = 3
		} else if waitTime < d.MaxWait*1.3 {
			currRating = 2
		} else if waitTime < d.MaxWait*1.4 {
			currRating = 1
		}

		atomic.AddInt64(&summedRatings, currRating)
		avg := float64(atomic.LoadInt64(&summedRatings)) / float64(atomic.LoadInt64(&CompletedOrders))
		KitchenOverloadMutex.Lock()
		if avg < AvgRating {
			CurrentMaxOrders -= 2
		} else {
			CurrentMaxOrders += 1
		}
		AvgRating = avg
		utils.LogRep.Info("AVG RATING", zap.Float64("rating", AvgRating), zap.Float64("waitTime", waitTime), zap.Float64("maxWait", d.MaxWait), zap.Float64("cookingTime", d.CookingTime), zap.Int("orderId", d.OrderId), zap.Int("currentOrders", CurrentOrders), zap.Int("currentMaxOrders", CurrentMaxOrders))

		KitchenOverloadMutex.Unlock()

	}
}
