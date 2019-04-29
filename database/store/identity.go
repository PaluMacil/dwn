package store

import (
	"encoding/binary"
)

type Identity uint64

func (i Identity) Bytes() []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(i))
	return buf[:]
}
