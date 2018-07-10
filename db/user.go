package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
)

type User struct {
	GoogleID         string    `json:"google_id"`
	GoogleImportDate time.Time `json:"google_import_date"`
	Email            string    `json:"email"`
	PasswordHash     []byte    `json:"-"`
	VerifiedEmail    bool      `json:"verified_email"`
	Locked           bool      `json:"locked"`
	DisplayName      string    `json:"display_name"`
	GivenName        string    `json:"given_name"`
	FamilyName       string    `json:"family_name"`
	Link             string    `json:"link"`
	Picture          string    `json:"picture"`
	Gender           string    `json:"gender"`
	Locale           string    `json:"locale"`
	LastLogin        time.Time `json:"last_login"`
	ModifiedDate     time.Time `json:"modified_date"`
	CreatedDate      time.Time `json:"created_date"`
	Can              func() bool
}

func (u User) Key() []byte {
	return append(u.Prefix(), []byte(u.Email)...)
}

func (u User) Prefix() []byte {
	return []byte(userPrefix)
}

type UserProvider struct {
	bgr *badger.DB
	Db  *Db
}

func (p *UserProvider) Get(email string) (User, error) {
	var user = User{Email: email}
	if email == "" {
		return user, errors.New("UserProvider.Get requires an email but got an empty string")
	}
	item, err := get(p.bgr, &user)
	if err != nil {
		return user, err
	}
	user, ok := item.(User)
	if !ok {
		return user, fmt.Errorf("got data of type %T but wanted User", user)
	}
	return user, err
}

func (p *UserProvider) Exists(email string) (bool, error) {
	_, err := p.Get(email)
	return err != badger.ErrKeyNotFound, err
}

func (p *UserProvider) Set(user User) error {
	return set(p.bgr, &user)
}

func (p *UserProvider) Count() (int, error) {
	return count(p.bgr, User{}.Prefix())
}

func (p *UserProvider) All() ([]User, error) {
	var items []DbItem
	err := all(p.bgr, User{}.Prefix(), &items, true)
	users := make([]User, len(items))
	for i, v := range items {
		users[i] = v.(User)
	}

	return users, err
}

func (p UserProvider) Delete(email string) error {
	return delete(p.bgr, User{Email: email})
}

func (p UserProvider) PurgeAll() error {
	users, err := p.All()
	if err != nil {
		return err
	}
	for _, u := range users {
		p.Delete(u.Email)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
func (p UserProvider) Is(groupName string) bool {
	for _, g := range u.GroupNames {
		if g == groupName {
			return true
		}
	}
	return false
}
*/
