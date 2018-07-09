package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/satori/go.uuid"
)

type Session struct {
	Token     string
	Email     string
	Created   time.Time
	Heartbeat time.Time
}

func (s Session) Key() []byte {
	return append(s.Prefix(), []byte(s.Token)...)
}

func (s Session) Prefix() []byte {
	return []byte(sessionPrefix)
}

type SessionProvider struct {
	bgr *badger.DB
}

func (p *SessionProvider) Get(token string) (Session, error) {
	var session = Session{Token: token}
	if token == "" {
		return session, errors.New("SessionProvider.Get requires a token but got an empty string")
	}
	item, err := get(p.bgr, &session)
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
	return err != badger.ErrKeyNotFound, err
}

func (p *SessionProvider) Set(session Session) error {
	return set(p.bgr, &session)
}

func (p *SessionProvider) GenerateFor(email string) Session {
	t := uuid.Must(uuid.NewV4())
	session := Session{
		Token:     t.String(),
		Email:     email,
		Created:   time.Now(),
		Heartbeat: time.Now(),
	}
	return session
}

func (p *SessionProvider) All() ([]Session, error) {
	var items []DbItem
	err := all(p.bgr, Session{}.Prefix(), &items, true)
	sessions := make([]Session, len(items))
	for i, v := range items {
		sessions[i] = v.(Session)
	}

	return sessions, err
}

func (p SessionProvider) Delete(token string) error {
	return delete(p.bgr, Session{Token: token})
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
