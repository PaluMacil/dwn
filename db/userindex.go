package db

import (
	"github.com/blevesearch/bleve"
)

type UserIndex struct {
	idx bleve.Index
	db *db.Db
}

func (idx *UserIndex) ReIndex() error {
	users, err := idx.db.Users.All()
	if err != nil {
		return err
	}

}

func NewUserIndex(dataDir string) (UserIndex, error) {
	return UserIndex{
		
	}
}

func (idx *UserIndex) CompletionSuggestions(query string) ([]db.User, error) {

}