package shoppingrepo

import (
	"errors"
	"fmt"

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

func (p ItemRepo) Get(name string) (shopping.Item, error) {
	var shoppingItem = shopping.Item{Name: name}
	if name == "" {
		return shoppingItem, errors.New("ItemRepo.Get requires a name but got an empty string")
	}
	item, err := p.store.Get(&shoppingItem)
	if err != nil {
		return shoppingItem, err
	}
	shoppingItem, ok := item.(shopping.Item)
	if !ok {
		return shoppingItem, fmt.Errorf("got data of type %T but wanted shopping.Item", shoppingItem)
	}
	return shoppingItem, err
}

func (p ItemRepo) Set(item shopping.Item) error {
	return p.store.Set(&item)
}
