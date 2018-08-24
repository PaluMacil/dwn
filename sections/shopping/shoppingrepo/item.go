package shoppingrepo

import (
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/sections/shopping"
)

type ItemRepo struct {
	store database.Storer
	db    *database.Database
}

func NewItemRepo(store database.Storer, db *database.Database) *ItemRepo {
	return &ItemRepo{store, db}
}

func (p ItemRepo) All() ([]shopping.Item, error) {
	var items []database.Item
	pfx := shopping.Item{}.Prefix()
	err := p.store.All(pfx, &items, true)
	shoppingItems := make([]shopping.Item, len(items))
	for i, v := range items {
		shoppingItems[i] = v.(shopping.Item)
	}

	return shoppingItems, err
}

func (p ItemRepo) Delete(name string) error {
	return p.store.Delete(shopping.Item{Name: name})
}
