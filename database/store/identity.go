package store

import (
	"encoding/binary"
	"strconv"
)

type Identity uint64

func (i Identity) Bytes() []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(i))
	return buf[:]
}

func (i Identity) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

func StringToIdentity(idString string) (Identity, error) {
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return 0, err
	}
	return Identity(idInt), nil
}