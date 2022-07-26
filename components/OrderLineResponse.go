package components

type OrderLineResponse struct {
	MatchingOrderLines        []OrderLineEntity                   `json:"MatchingOrderLines"`
}
