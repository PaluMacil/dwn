package shopping

import (
	"time"

	"github.com/PaluMacil/dwn/dwn"
)

const (
	ItemPrefix = "SHOPPING_ITEM:"
)

type Item struct {
	Name     string          `json:"name"`
	Quantity int             `json:"quantity"`
	Note     string          `json:"note"`
	AddedBy  dwn.DisplayName `json:"addedBy"`
	Added    time.Time       `json:"added"`
}

func (i Item) Key() []byte {
	return append(i.Prefix(), []byte(i.Name)...)
}

func (i Item) Prefix() []byte {
	return []byte(ItemPrefix)
}
