package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/darkcat013/pr-dining-hall/constants"
	"github.com/darkcat013/pr-dining-hall/domain"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var log *zap.Logger
var menu []domain.Food

func SetMenu(jsonPath string) {
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Fatal("Error opening " + jsonPath)
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &menu)

	if menu == nil {
		log.Fatal("Failed to decode menu from " + jsonPath)
	}
	log.Info("Menu decoded and set")
}

func GenerateOrders() {
	orderId := 1
	maxWaiters := 4
	currentWaiters := 0

	// var checkWaiterMutex sync.Mutex
	var addWaiterMutex sync.Mutex
	var subtractWaiterMutex sync.Mutex
	var orderIdMutex sync.Mutex

	for {
		if currentWaiters < maxWaiters {
			addWaiterMutex.Lock()
			currentWaiters++
			addWaiterMutex.Unlock()

			go func() {
				waiterDelay := rand.Intn(5) + 1
				log.Info(fmt.Sprintf("Waiter goes to pick-up order: delay %d seconds", waiterDelay))
				time.Sleep(time.Duration(waiterDelay) * constants.TIME_UNIT)

				itemsCount := rand.Intn(10) + 1
				var items []int
				var maxPrepTime int

				for i := 0; i < itemsCount; i++ {
					randomFood := menu[rand.Intn(len(menu))]
					if maxPrepTime < randomFood.PreparationTime {
						maxPrepTime = randomFood.PreparationTime
					}

					items = append(items, randomFood.Id)
				}
				orderIdMutex.Lock()
				log.Info(fmt.Sprintf("Started to create order with id %d", orderId))
				order := domain.Order{
					OrderId:    orderId,
					TableId:    0,
					WaiterId:   0,
					Items:      items,
					Priority:   0,
					MaxWait:    int(float64(maxPrepTime) * 1.3),
					PickUpTime: int(time.Now().Unix()),
				}
				orderId++
				orderIdMutex.Unlock()
				log.Info(fmt.Sprintf("Created order with id %d", orderId-1), zap.Any("order", order))

				body, err := json.Marshal(order)
				if err != nil {
					log.Fatal("Failed to convert order to JSON ", zap.String("error", err.Error()), zap.Any("order", order))
				}

				log.Info("Send order to kitchen", zap.Int("orderId", orderId))

				resp, err := http.Post(constants.KITCHEN_URL, "application/json", bytes.NewBuffer(body))

				if err != nil {
					log.Error("Failed to send order to kitchen", zap.String("error", err.Error()), zap.Int("orderId", orderId))
				} else {
					log.Info("Response from kitchen", zap.Int("statusCode", resp.StatusCode), zap.Int("orderId", orderId))
				}
				subtractWaiterMutex.Lock()
				currentWaiters--
				subtractWaiterMutex.Unlock()
			}()
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	log = zap.NewExample()
	defer log.Sync()

	SetMenu("config/menu.json")
	go GenerateOrders()

	router := mux.NewRouter()
	router.HandleFunc("/distribution", func(w http.ResponseWriter, r *http.Request) {
		var d domain.Distribution
		err := json.NewDecoder(r.Body).Decode(&d)

		if err != nil {
			log.Error("Failed to decode distribution",
				zap.String("error", err.Error()),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Info("Distribution decoded", zap.Any("distribution", d))

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("Requested",
			zap.String("method", r.Method),
			zap.String("endpoint", r.URL.String()),
		)
		router.ServeHTTP(w, r)
	})

	http.Handle("/", router)
	log.Info("Started web server at port :8081")
	http.ListenAndServe(":8081", handler)
}
