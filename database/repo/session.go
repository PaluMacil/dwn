package repo

import (
	"errors"
	"fmt"
	"github.com/PaluMacil/dwn/module/core"
	"time"

	"github.com/PaluMacil/dwn/database"
	"github.com/gofrs/uuid"
)

type SessionRepo struct {
	store database.Storer
	db    *database.Database
}

func NewSessionRepo(store database.Storer, db *database.Database) *SessionRepo {
	return &SessionRepo{store, db}
}

func (p SessionRepo) Get(token string) (core.Session, error) {
	var session = core.Session{Token: token}
	if token == "" {
		return session, errors.New("SessionRepo.Get requires a token but got an empty string")
	}
	item, err := p.store.Get(&session)
	if err != nil {
		return session, err
	}
	session, ok := item.(core.Session)
	if !ok {
		return session, fmt.Errorf("got data of type %T but wanted core.Session", session)
	}
	return session, err
}

func (p SessionRepo) Exists(token string) (bool, error) {
	_, err := p.Get(token)
	if p.db.Util.IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

func (p SessionRepo) Set(session core.Session) error {
	return p.store.Set(&session)
}

func (p SessionRepo) GenerateFor(email, ip string) core.Session {
	t := uuid.Must(uuid.NewV4())
	session := core.Session{
		Token:       t.String(),
		Email:       email,
		IP:          ip,
		CreatedDate: time.Now(),
		Heartbeat:   time.Now(),
	}
	return session
}

func (p SessionRepo) All() ([]core.Session, error) {
	var items []database.Item
	pfx := core.Session{}.Prefix()
	err := p.store.All(pfx, &items, true)
	sessions := make([]core.Session, len(items))
	for i, v := range items {
		sessions[i] = v.(core.Session)
	}

	return sessions, err
}

func (p SessionRepo) Delete(token string) error {
	return p.store.Delete(core.Session{Token: token})
}

func (p SessionRepo) PurgeAll() error {
	sessions, err := p.All()
	if err != nil {
		return err
	}
	for _, s := range sessions {
		err = p.Delete(s.Token)
		if err != nil {
			return err
		}
	}
	return nil
}
