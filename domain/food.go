package domain

type Food struct {
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	PreparationTime  float64 `json:"preparation-time"`
	Complexity       int     `json:"complexity"`
	CookingApparatus string  `json:"cooking-apparatus"`
}
