package search

import (
	"fmt"
	"path"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/dwn"
	"github.com/blevesearch/bleve"
)

type UserIndex struct {
	idx bleve.Index
	db  *database.Database
}

func (ui UserIndex) Reindex() error {
	users, err := ui.db.Users.All()
	if err != nil {
		return err
	}
	for _, u := range users {
		err = ui.Index(u)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ui UserIndex) Index(u dwn.User) error {
	err := ui.idx.Index(u.Email, u)
	if err != nil {
		return err
	}
	return nil
}

func (ui UserIndex) Deindex(u dwn.User) error {
	err := ui.idx.Delete(u.Email)
	if err != nil {
		return err
	}
	return nil
}

func NewUserIndex(db *database.Database, dataDir string) (*UserIndex, error) {
	indexPath := path.Join(dataDir, "indexes", "user")
	index, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, fmt.Errorf(`creating user index: %s`, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf(`opening user index: %s`, err)
	}
	return &UserIndex{index, db}, nil
}

func (ui UserIndex) CompletionSuggestions(query string) ([]dwn.User, error) {
	queryWithoutPrefix := bleve.NewPrefixQuery(query)
	queryWithPrefix := bleve.NewPrefixQuery("@" + query)
	searchQuery := bleve.NewDisjunctionQuery(queryWithoutPrefix, queryWithPrefix)
	search := bleve.NewSearchRequest(searchQuery)
	result, err := ui.idx.Search(search)
	if err != nil {
		return []dwn.User{}, err
	}
	users := make([]dwn.User, len(result.Hits))
	for i, res := range result.Hits {
		email := res.ID
		u, err := ui.db.Users.Get(email)
		if err != nil {
			return users, err
		}
		users[i] = u
	}
	return users, nil
}