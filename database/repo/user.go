package repo

import (
	"errors"
	"fmt"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/dwn"
)

type UserRepo struct {
	store database.Storer
	db    *database.Database
	dwn.UserSearcher
}

func NewUserRepo(store database.Storer, db *database.Database, search dwn.UserSearcher) *UserRepo {
	return &UserRepo{store, db, search}
}

func (p UserRepo) UsersFor(groupName string) ([]dwn.User, error) {
	userGroups, err := p.db.UserGroups.All()
	if err != nil {
		return nil, err
	}
	var users []dwn.User
	for _, ug := range userGroups {
		if ug.GroupName == groupName {
			u, err := p.Get(ug.Email)
			if err != nil {
				return nil, err
			}
			users = append(users, u)
		}
	}
	return users, nil
}

func (p UserRepo) Get(email string) (dwn.User, error) {
	var user = dwn.User{Email: email}
	if email == "" {
		return user, errors.New("UserRepo.Get requires an email but got an empty string")
	}
	item, err := p.store.Get(&user)
	if err != nil {
		return user, err
	}
	user, ok := item.(dwn.User)
	if !ok {
		return user, fmt.Errorf("got data of type %T but wanted User", user)
	}
	return user, err
}

//TODO: exists can be pushed into a db method
func (p UserRepo) Exists(email string) (bool, error) {
	_, err := p.Get(email)
	if p.db.Util.IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

func (p UserRepo) Set(user dwn.User) error {
	err := p.Index(user)
	if err != nil {
		return err
	}
	return p.store.Set(&user)
}

func (p UserRepo) Count() (int, error) {
	return p.store.Count(dwn.User{}.Prefix())
}

func (p UserRepo) All() ([]dwn.User, error) {
	var items []database.Item
	err := p.store.All(dwn.User{}.Prefix(), &items, true)
	users := make([]dwn.User, len(items))
	for i, v := range items {
		users[i] = v.(dwn.User)
	}

	return users, err
}

func (p UserRepo) Delete(email string) error {
	u := dwn.User{Email: email}
	err := p.Deindex(u)
	if err != nil {
		return err
	}
	return p.store.Delete(u)
}

func (p UserRepo) PurgeAll() error {
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
