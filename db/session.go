package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/satori/go.uuid"
)

type Session struct {
	Token       string    `json:"token"`
	Email       string    `json:"email"`
	CreatedDate time.Time `json:"createdDate"`
	Heartbeat   time.Time `json:"heartbeat"`
}

func (s Session) Key() []byte {
	return append(s.Prefix(), []byte(s.Token)...)
}

func (s Session) Prefix() []byte {
	return []byte(sessionPrefix)
}

type SessionProvider struct {
	Db *Db
}

func (p *SessionProvider) Get(token string) (Session, error) {
	var session = Session{Token: token}
	if token == "" {
		return session, errors.New("SessionProvider.Get requires a token but got an empty string")
	}
	item, err := p.Db.get(&session)
	if err != nil {
		return session, err
	}
	session, ok := item.(Session)
	if !ok {
		return session, fmt.Errorf("got data of type %T but wanted Session", session)
	}
	return session, err
}

func (p *SessionProvider) Exists(token string) (bool, error) {
	_, err := p.Get(token)
	if err == badger.ErrKeyNotFound {
		return false, nil
	}
	return true, err
}

func (p *SessionProvider) Set(session Session) error {
	return p.Db.set(&session)
}

func (p *SessionProvider) GenerateFor(email string) Session {
	t := uuid.Must(uuid.NewV4())
	session := Session{
		Token:       t.String(),
		Email:       email,
		CreatedDate: time.Now(),
		Heartbeat:   time.Now(),
	}
	return session
}

func (p *SessionProvider) All() ([]Session, error) {
	var items []DbItem
	err := p.Db.all(Session{}.Prefix(), &items, true)
	sessions := make([]Session, len(items))
	for i, v := range items {
		sessions[i] = v.(Session)
	}

	return sessions, err
}

func (p SessionProvider) Delete(token string) error {
	return p.Db.delete(Session{Token: token})
}

func (p SessionProvider) PurgeAll() error {
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
