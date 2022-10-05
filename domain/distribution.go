package domain

type Distribution struct {
	OrderId        int             `json:"order_id"`
	TableId        int             `json:"table_id"`
	WaiterId       int             `json:"waiter_id"`
	Items          []int           `json:"items"`
	Priority       int             `json:"priority"`
	MaxWait        float64         `json:"max_wait"`
	PickUpTime     float64         `json:"pick_up_time"`
	CookingTime    float64         `json:"cooking_time"`
	CookingDetails []CookingDetail `json:"cooking_details"`
}
