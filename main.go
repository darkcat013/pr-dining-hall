package main

import (
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
	//check if kitchen is open
	rand.Seed(time.Now().UnixNano())
	domain.InitializeMenu("config/menu.json")

	for i := 0; i < config.TABLES; i++ {
		table := domain.NewTable(i)
		domain.Tables = append(domain.Tables, table)
		go table.Start()
	}

	for i := 0; i < config.WAITERS; i++ {
		waiter := domain.NewWaiter(i)
		domain.Waiters = append(domain.Waiters, waiter)
		go waiter.Start()
	}

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

	http.HandleFunc("/", unhandledRoutes)
	http.HandleFunc("/distribution", distribution)

	utils.Log.Info("Started web server at port :8081")

	if err := http.ListenAndServe(":8081", nil); err != nil {
		utils.Log.Fatal("Could not start web server", zap.String("error", err.Error()))
	}
}
