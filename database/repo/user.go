package repo

import (
	"errors"
	"fmt"
	"github.com/PaluMacil/dwn/module/core"

	"github.com/PaluMacil/dwn/database"
)

type UserRepo struct {
	store database.Storer
	db    *database.Database
	database.UserSearcher
}

func NewUserRepo(store database.Storer, db *database.Database, search database.UserSearcher) *UserRepo {
	return &UserRepo{store, db, search}
}

func (p UserRepo) UsersFor(groupName string) ([]core.User, error) {
	userGroups, err := p.db.UserGroups.All()
	if err != nil {
		return nil, err
	}
	var users []core.User
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

func (p UserRepo) Get(email string) (core.User, error) {
	var user = core.User{Email: email}
	if email == "" {
		return user, errors.New("UserRepo.Get requires an email but got an empty string")
	}
	item, err := p.store.Get(&user)
	if err != nil {
		return user, err
	}
	user, ok := item.(core.User)
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

func (p UserRepo) Set(user core.User) error {
	err := p.Index(user)
	if err != nil {
		return err
	}
	return p.store.Set(&user)
}

func (p UserRepo) Count() (int, error) {
	return p.store.Count(core.User{}.Prefix())
}

func (p UserRepo) All() ([]core.User, error) {
	var items []database.Item
	err := p.store.All(core.User{}.Prefix(), &items, true)
	users := make([]core.User, len(items))
	for i, v := range items {
		users[i] = v.(core.User)
	}

	return users, err
}

func (p UserRepo) Delete(email string) error {
	u := core.User{Email: email}
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
