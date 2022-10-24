package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/darkcat013/pr-dining-hall/config"
	"github.com/darkcat013/pr-dining-hall/domain"
	"github.com/darkcat013/pr-dining-hall/utils"
	"go.uber.org/zap"
)

func StartServer() {
	unhandledRoutes := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		utils.Log.Warn("Path not found", zap.Int("statusCode", http.StatusNotFound))
		http.Error(w, "404 path not found.", http.StatusNotFound)
	}

	distribution := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		if r.Method != "POST" {
			utils.Log.Warn("Method not allowed", zap.Int("statusCode", http.StatusMethodNotAllowed))
			http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
			return
		}

		var d domain.Distribution
		err := json.NewDecoder(r.Body).Decode(&d)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			utils.Log.Fatal("Failed to decode distribution", zap.String("error", err.Error()))
			return
		}
		utils.Log.Info("Distribution decoded", zap.Any("distribution", d))

		domain.Waiters[d.WaiterId].ReceiveDistributionChan <- d

		w.WriteHeader(http.StatusOK)
	}

	v2order := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		if r.Method != "POST" {
			utils.Log.Warn("Method not allowed", zap.Int("statusCode", http.StatusMethodNotAllowed))
			http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
			return
		}

		var o domain.FoodServiceOrder
		err := json.NewDecoder(r.Body).Decode(&o)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			utils.Log.Fatal("Failed to decode order from food order service", zap.String("error", err.Error()))
			return
		}

		utils.Log.Info("Food service order decoded", zap.Any("data", o))

		var responseObj domain.FoodServiceOrderResponse

		responseObj = domain.FoodServiceOrderResponse{
			RestaurantId:         config.RESTAURANT_ID,
			OrderId:              0,
			EstimatedWaitingTime: 0,
			CreatedTime:          utils.GetCurrentTimeFloat(),
			RegisteredTime:       domain.RegisteredTime,
		}

		response, err := json.Marshal(responseObj)

		if err != nil {
			utils.Log.Fatal("Failed to convert response food service order to JSON ", zap.String("error", err.Error()), zap.Any("responseOrder", responseObj))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}

	v2orderGet := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		splitPath := strings.Split(r.URL.Path, "/")
		orderId := splitPath[len(splitPath)-1]

		switch r.Method {
		case "GET":
			utils.Log.Info("V2 order GET " + orderId)
			w.WriteHeader(http.StatusOK)
			return
		default:
			utils.Log.Warn("Method not allowed", zap.Int("statusCode", http.StatusMethodNotAllowed))
			http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
			return
		}
	}

	v2rating := func(w http.ResponseWriter, r *http.Request) {

		utils.Log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)

		switch r.Method {
		case "POST":
			utils.Log.Info("V2 rating POST")
			w.WriteHeader(http.StatusOK)
			return
		default:
			utils.Log.Warn("Method not allowed", zap.Int("statusCode", http.StatusMethodNotAllowed))
			http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
			return
		}
	}

	http.HandleFunc("/", unhandledRoutes)
	http.HandleFunc("/distribution", distribution)
	http.HandleFunc("/v2/order", v2order)
	http.HandleFunc("/v2/order/", v2orderGet) //v2/order/{id}
	http.HandleFunc("/v2/rating", v2rating)

	utils.Log.Info("Started web server at port :" + config.PORT)

	if err := http.ListenAndServe(":"+config.PORT, nil); err != nil {
		utils.Log.Fatal("Could not start web server", zap.String("error", err.Error()))
	}
}
