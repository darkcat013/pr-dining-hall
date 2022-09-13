package domain

type Distribution struct {
	OrderId        int             `json:"order_id"`
	TableId        int             `json:"table_id"`
	WaiterId       int             `json:"waiter_id"`
	Items          []int           `json:"items"`
	Priority       int             `json:"priority"`
	MaxWait        int             `json:"max_wait"`
	PickUpTime     int             `json:"pick_up_time"`
	CookingTime    int             `json:"cooking_time"`
	CookingDetails []CookingDetail `json:"cooking_details"`
}
