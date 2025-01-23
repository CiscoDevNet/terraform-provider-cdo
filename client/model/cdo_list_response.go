package model

type CdoListResponse[T any] struct {
	Count int `json:"count"`
	Items []T `json:"items"`
}
