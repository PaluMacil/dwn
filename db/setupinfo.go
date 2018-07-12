package db

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
)

type SetupInfo struct {
	InitializedDate time.Time
}

func (s SetupInfo) Key() []byte {
	return s.Prefix()
}

func (s SetupInfo) Prefix() []byte {
	return []byte(setupInfoPrefix)
}

type SetupInfoProvider struct {
	bgr *badger.DB
	Db  *Db
}

func (p *SetupInfoProvider) Get() (SetupInfo, error) {
	var setupInfo = SetupInfo{}
	item, err := get(p.bgr, &setupInfo)
	if err != nil {
		return setupInfo, err
	}
	setupInfo, ok := item.(SetupInfo)
	if !ok {
		return setupInfo, fmt.Errorf("got data of type %T but wanted SetupInfo", setupInfo)
	}
	return setupInfo, err
}

func (p *SetupInfoProvider) Completed() (bool, error) {
	_, err := p.Get()
	if err == badger.ErrKeyNotFound {
		return false, nil
	}
	return true, err
}

func (p *SetupInfoProvider) Set(setupInfo SetupInfo) error {
	return set(p.bgr, &setupInfo)
}

func (p SetupInfoProvider) Delete() error {
	return delete(p.bgr, SetupInfo{})
}
