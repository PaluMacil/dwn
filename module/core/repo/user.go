package repo

import (
	"errors"
	"fmt"

	"github.com/PaluMacil/dwn/module/core"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/store"
)

type UserRepo struct {
	store database.Storer
	db    *database.Database
	core.UserSearcher
}

func NewUserRepo(store database.Storer, db *database.Database, search core.UserSearcher) *UserRepo {
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
			u, err := p.Get(ug.UserID)
			if err != nil {
				return nil, err
			}
			users = append(users, u)
		}
	}
	return users, nil
}

func (p UserRepo) Get(userID store.Identity) (core.User, error) {
	var user = core.User{ID: userID}
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

func (p UserRepo) Exists(userID store.Identity) (bool, error) {
	_, err := p.Get(userID)
	if p.db.IsKeyNotFoundErr(err) {
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

func (p UserRepo) Delete(userID store.Identity) error {
	u := core.User{UserID: userID}
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
		p.Delete(u.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
