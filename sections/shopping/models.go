package shopping

import (
	"time"
)

const (
	ItemPrefix = "SHOPPING_ITEM:"
)

type Item struct {
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
	Note     string    `json:"note"`
	AddedBy  string    `json:"addedBy"`
	Added    time.Time `json:"added"`
}

func (i Item) Key() []byte {
	return append(i.Prefix(), []byte(i.Name)...)
}

func (i Item) Prefix() []byte {
	return []byte(ItemPrefix)
}
