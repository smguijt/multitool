package components

type OrderResponse struct {
	MatchingOrders            []OrderEntity                       `json:"MatchingOrders"`
}
