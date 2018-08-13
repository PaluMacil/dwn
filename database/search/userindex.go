package search

import (
	"log"

	"github.com/PaluMacil/dwn/dwn"
	"github.com/PaluMacil/dwn/database"
	"github.com/blevesearch/bleve"
)

type UserIndex struct {
	idx bleve.Index
	db  database.Database
}

func (idx UserIndex) ReIndex() error {
	users, err := idx.db.Users.All()
	if err != nil {
		return err
	}
	//TODO: do a reindex
	log.Println(users)
	return nil
}

func NewUserIndex(db database.Database, dataDir string) (*UserIndex, error) {
	return &UserIndex{}, nil
}

func (idx UserIndex) CompletionSuggestions(query string) ([]dwn.User, error) {
	return []dwn.User{}, nil
}
