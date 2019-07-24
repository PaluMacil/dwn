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

// Reindex takes all users in the database and reindexes them.
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

// Index will index a single user
func (ui UserIndex) Index(u core.User) error {
	id := u.ID.String()
	err := ui.idx.Index(id, u)
	if err != nil {
		return err
	}
	return nil
}

// Deindex deletes the index for a single user
func (ui UserIndex) Deindex(u core.User) error {
	id := u.ID.String()
	err := ui.idx.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// NewUserIndex returns a new *UserIndex based upon the database and data directory given.
// Indexes will be stored at: {dataDir}/indexes/user
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

// CompletionSuggestions returns user suggestions starting with the query string or the query
// string with an ampersand in front.
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

// WithEmail returns all users with this email regardless of whether they have verified
// the email.
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

// FromEmail returns the one verified user with this email or else a KeyNotFoundErr
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

// EmailExists returns whether a user with an email exists, regardless of whether the
// existing user has verified that email.
func (ui UserIndex) EmailExists(email string) (bool, error) {
	searchQuery := bleve.NewMatchQuery(email)
	search := bleve.NewSearchRequest(searchQuery)
	result, err := ui.idx.Search(search)
	if err != nil {
		return false, err
	}
	return len(result.Hits) > 0, nil
}

// VerifiedEmailExists returns whether a user with an email exists if this
// existing user has also verified that email.
func (ui UserIndex) VerifiedEmailExists(email string) (bool, error) {
	users, err := ui.WithEmail(email)
	if err != nil {
		return false, err
	}
	for _, user := range users {
		for _, e := range user.Emails {
			// Check if matched AND verified.
			if e.Email == email && e.Verified {
				return true, nil
			}
		}
	}
	return false, nil
}
