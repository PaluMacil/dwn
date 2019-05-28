package search

import (
	"fmt"
	"path"
	"strconv"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
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

func (ui UserIndex) Index(u core.User) error {
	id := strconv.Itoa(u.ID)
	err := ui.idx.Index(id, u)
	if err != nil {
		return err
	}
	return nil
}

func (ui UserIndex) Deindex(u core.User) error {
	id := strconv.Itoa(u.ID)
	err := ui.idx.Delete(id)
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

func (ui UserIndex) CompletionSuggestions(query string) ([]core.User, error) {
	queryWithoutPrefix := bleve.NewPrefixQuery(query)
	queryWithPrefix := bleve.NewPrefixQuery("@" + query)
	searchQuery := bleve.NewDisjunctionQuery(queryWithoutPrefix, queryWithPrefix)
	search := bleve.NewSearchRequest(searchQuery)
	result, err := ui.idx.Search(search)
	if err != nil {
		return []core.User{}, err
	}
	users := make([]core.User, len(result.Hits))
	for i, res := range result.Hits {
		id, err := strconv.Atoi(res.ID)
		if err != nil {
			return users, err
		}
		u, err := ui.db.Users.Get(id)
		if err != nil {
			return users, err
		}
		users[i] = u
	}
	return users, nil
}

func (ui UserIndex) FromEmail(email string) (User, error) {
	searchQuery := bleve.NewMatchQuery(email)
}
