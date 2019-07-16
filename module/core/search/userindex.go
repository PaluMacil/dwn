package search

import (
	"fmt"
	"path"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/store"
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
	id := u.ID.String()
	err := ui.idx.Index(id, u)
	if err != nil {
		return err
	}
	return nil
}

func (ui UserIndex) Deindex(u core.User) error {
	id := u.ID.String()
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
		id, err := store.StringToIdentity(res.ID)
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

func (ui UserIndex) WithEmail(email string) ([]core.User, error) {
	searchQuery := bleve.NewMatchQuery(email)
	search := bleve.NewSearchRequest(searchQuery)
	result, err := ui.idx.Search(search)
	if err != nil {
		return nil, err
	}
	users := make([]core.User, len(result.Hits))
	for i, res := range result.Hits {
		id, err := store.StringToIdentity(res.ID)
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

func (ui UserIndex) FromEmail(email string) (core.User, error) {
	users, err := ui.WithEmail(email)
	if err != nil {
		return core.User{}, err
	}
	for _, user := range users {
		for _, e := range user.Emails {
			// Check if matched AND verified.
			if e.Email == email && e.Verified {
				return user, nil
			}
		}
	}

	return core.User{}, ui.db.KeyNotFoundErr()
}

func (ui UserIndex) EmailExists(email string) (bool, error) {
	searchQuery := bleve.NewMatchQuery(email)
	search := bleve.NewSearchRequest(searchQuery)
	result, err := ui.idx.Search(search)
	if err != nil {
		return false, err
	}
	return len(result.Hits) > 0, nil
}
