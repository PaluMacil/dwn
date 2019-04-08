package core

import (
	"net/http"
	"strings"
	"time"
)

const SessionPrefix   = "SESSION:"

type Session struct {
	Token         string `json:"token"`
	Email         string `json:"email"`
	IP            string `json:"ip"`
	Proxy         bool
	VaultUnlocked bool
	CreatedDate   time.Time `json:"createdDate"`
	Heartbeat     time.Time `json:"heartbeat"`
}

func (s Session) Key() []byte {
	return append(s.Prefix(), []byte(s.Token)...)
}

func (s Session) Prefix() []byte {
	return []byte(SessionPrefix)
}

func IP(r *http.Request) string {
	ip := r.RemoteAddr
	forwardForIP := r.Header.Get("X-Forwarded-For")
	if forwardForIP != "" {
		// if the user if coming through one or more proxies, one or more IP addresses
		// could be set in this comma separated header. The first IP is the user's
		// original IP.
		ip = strings.Split(forwardForIP, ",")[0]
	}
	return ip
}
