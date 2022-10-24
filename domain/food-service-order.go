package domain

type FoodServiceOrder struct {
	Items       []int   `json:"items"`
	Priority    int     `json:"priority"`
	MaxWait     float64 `json:"max_wait"`
	CreatedTime float64 `json:"created_time"`
}

type FoodServiceOrderResponse struct {
	RestaurantId         int     `json:"restaurant_id"`
	OrderId              int     `json:"order_id"`
	EstimatedWaitingTime float64 `json:"estimated_waiting_time"`
	CreatedTime          float64 `json:"created_time"`
	RegisteredTime       float64 `json:"registered_time"`
}
