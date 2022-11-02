package domain

type RestaurantData struct {
	RestaurantId int     `json:"restaurant_id"`
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	MenuItems    int     `json:"menu_items"`
	Menu         []Food  `json:"menu"`
	Rating       float64 `json:"rating"`
}

type RestaurantPayloadRating struct {
	OrderId              int     `json:"order_id"`
	Rating               int     `json:"rating"`
	EstimatedWaitingTime float64 `json:"estimated_waiting_time"`
	WaitingTime          float64 `json:"waiting_time"`
}

type RestaurantResponseRating struct {
	RestaurantId        int     `json:"restaurant_id"`
	RestaurantAvgRating float64 `json:"restaurant_avg_rating"`
	PreparedOrders      int     `json:"prepared_orders"`
}
