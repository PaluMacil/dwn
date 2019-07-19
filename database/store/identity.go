package store

import (
	"encoding/binary"
	"strconv"
)

type Identity uint64

type Identities []Identity

func (i Identity) Bytes() []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(i))
	return buf[:]
}

func (i Identity) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func (i Identities) Distinct() Identities {
	set := map[Identity]bool{}
	for _, id := range i {
		set[id] = true
	}
	values := make([]Identity, 0, len(set))
	for id := range set {
		values = append(values, id)
	}
	return Identities(values)
}

func StringToIdentity(idString string) (Identity, error) {
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return 0, err
	}
	return Identity(idInt), nil
}
