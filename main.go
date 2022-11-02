package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/darkcat013/pr-dining-hall/config"
	"github.com/darkcat013/pr-dining-hall/domain"
	"github.com/darkcat013/pr-dining-hall/utils"
	"go.uber.org/zap"
)

func main() {

	utils.InitializeLogger()
	rand.Seed(time.Now().UnixNano())
	domain.InitializeMenu(config.MENU_PATH)

	go domain.StartRatingLogging()
	domain.InitializeWaiters()
	go domain.InitializeTables()
	// go RegisterRestaurant()

	StartServer()
}

func RegisterRestaurant() {

	utils.Log.Info("Start registering restaurant", zap.Int("restaurantId", config.RESTAURANT_ID))
	restaurantData := domain.RestaurantData{
		RestaurantId: config.RESTAURANT_ID,
		Name:         config.RESTAURANT_NAME,
		Address:      config.ADDRESS,
		MenuItems:    len(domain.Menu),
		Menu:         domain.Menu,
		Rating:       domain.AvgRating,
	}

	body, err := json.Marshal(restaurantData)
	if err != nil {
		utils.Log.Fatal("Failed to convert restaurant data to JSON ", zap.String("error", err.Error()), zap.Any("data", restaurantData))
	}

	utils.Log.Info("Send restaurant data to food ordring service", zap.Any("restaurantData", restaurantData))

	resp, err := http.Post(config.FOOD_ORDERING_SERVICE_URL+"/register", "application/json", bytes.NewBuffer(body))

	if err != nil {
		utils.Log.Fatal("Failed to send restaurant data to food ordering service", zap.String("error", err.Error()), zap.Any("data", restaurantData))
	} else {
		utils.Log.Info("Response from food order service", zap.Int("statusCode", resp.StatusCode), zap.Int("restaurantId", config.RESTAURANT_ID))
	}
}
