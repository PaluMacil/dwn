package repo

import (
	"errors"
	"fmt"
	"time"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/dwn"
	"github.com/satori/go.uuid"
)

type SessionRepo struct {
	store database.Storer
	db    *database.Database
}

func NewSessionRepo(store database.Storer, db *database.Database) *SessionRepo {
	return &SessionRepo{store, db}
}

func (p SessionRepo) Get(token string) (dwn.Session, error) {
	var session = dwn.Session{Token: token}
	if token == "" {
		return session, errors.New("SessionRepo.Get requires a token but got an empty string")
	}
	item, err := p.store.Get(&session)
	if err != nil {
		return session, err
	}
	session, ok := item.(dwn.Session)
	if !ok {
		return session, fmt.Errorf("got data of type %T but wanted dwn.Session", session)
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

func (p SessionRepo) Set(session dwn.Session) error {
	return p.store.Set(&session)
}

func (p SessionRepo) GenerateFor(email, ip string) dwn.Session {
	t := uuid.Must(uuid.NewV4())
	session := dwn.Session{
		Token:       t.String(),
		Email:       email,
		IP:          ip,
		CreatedDate: time.Now(),
		Heartbeat:   time.Now(),
	}
	return session
}

func (p SessionRepo) All() ([]dwn.Session, error) {
	var items []database.Item
	pfx := dwn.Session{}.Prefix()
	err := p.store.All(pfx, &items, true)
	sessions := make([]dwn.Session, len(items))
	for i, v := range items {
		sessions[i] = v.(dwn.Session)
	}

	return sessions, err
}

func (p SessionRepo) Delete(token string) error {
	return p.store.Delete(dwn.Session{Token: token})
}

func (p SessionRepo) PurgeAll() error {
	sessions, err := p.All()
	if err != nil {
		return err
	}
	for _, s := range sessions {
		p.Delete(s.Token)
		if err != nil {
			return err
		}
	}
	return nil
}
