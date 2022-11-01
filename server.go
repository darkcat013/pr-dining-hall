package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

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

		var clientOrder domain.FoodServiceOrder
		err := json.NewDecoder(r.Body).Decode(&clientOrder)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			utils.Log.Fatal("Failed to decode order from food order service", zap.String("error", err.Error()))
			return
		}

		utils.Log.Info("Food service order decoded", zap.Any("data", clientOrder))

		oId := atomic.AddInt64(&domain.OrderId, 1)
		order := domain.Order{
			OrderId:  int(oId),
			TableId:  -1,
			Items:    clientOrder.Items,
			Priority: clientOrder.Priority,
			MaxWait:  clientOrder.MaxWait,
		}

		utils.Log.Info("Food service order converted to order", zap.Any("clientOrder", clientOrder), zap.Any("order", order))

		domain.ClientOrderChan <- order

		utils.Log.Info("Food service order sent to waiters", zap.Any("data", order))

		responseObj := domain.FoodServiceOrderResponse{
			RestaurantId:         config.RESTAURANT_ID,
			OrderId:              int(oId),
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
			oId, err := strconv.Atoi(orderId)

			if err != nil {
				utils.Log.Fatal("Failed to convert orderId to int ", zap.String("error", err.Error()), zap.String("orderId", orderId))
			}

			utils.Log.Info("V2 order GET", zap.Int("orderId", oId))

			order := domain.ReadyClientOrders[oId]
			var responseObj domain.DistributionClient
			if order.CookingDetails == nil {
				responseObj = domain.DistributionClient{
					OrderId:              order.OrderId,
					IsReady:              false,
					EstimatedWaitingTime: 10,
					Priority:             order.Priority,
					MaxWait:              order.MaxWait,
					CreatedTime:          order.PickUpTime,
					RegisteredTime:       domain.RegisteredTime,
				}
			} else {
				responseObj = domain.DistributionClient{
					OrderId:              order.OrderId,
					IsReady:              true,
					EstimatedWaitingTime: 0,
					Priority:             order.Priority,
					MaxWait:              order.MaxWait,
					CreatedTime:          order.PickUpTime,
					RegisteredTime:       domain.RegisteredTime,
					PreparedTime:         utils.GetCurrentTimeFloat(),
					CookingTime:          order.CookingTime,
					CookingDetails:       order.CookingDetails,
				}
			}

			response, err := json.Marshal(responseObj)

			if err != nil {
				utils.Log.Fatal("Failed to convert distribution to JSON for client", zap.String("error", err.Error()), zap.Any("distribution", responseObj))
			}

			utils.Log.Info("Converted distribution to JSON for client", zap.String("error", err.Error()), zap.Any("distribution", responseObj))

			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
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
