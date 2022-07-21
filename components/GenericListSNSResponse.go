package components

type GenericListSNSResponse struct {
	ListItems                 []GenericListEntity                 `json:"ListItems"`
	ListKey                   string                              `json:"ListKey"`
}
