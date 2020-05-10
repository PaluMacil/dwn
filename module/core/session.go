package core

import (
	"net/http"
	"strings"
	"time"

	"github.com/PaluMacil/dwn/database/store"
)

const SessionPrefix = "SESSION:"

type Session struct {
	Token         string         `json:"token"`
	UserID        store.Identity `json:"userID"`
	IP            string         `json:"ip"`
	Proxy         bool           `json:"proxy"`
	VaultUnlocked bool           `json:"vaultUnlocked"`
	CreatedDate   time.Time      `json:"createdDate"`
	Heartbeat     time.Time      `json:"heartbeat"`
}

func (s Session) Key() []byte {
	return append(s.Prefix(), []byte(s.Token)...)
}

func (s Session) Prefix() []byte {
	return []byte(SessionPrefix)
}

func (s Session) Info() SessionInfo {
	return SessionInfo{
		IP:          s.IP,
		Proxy:       s.Proxy,
		CreatedDate: s.CreatedDate,
		Heartbeat:   s.Heartbeat,
	}
}

type SessionInfo struct {
	IP          string    `json:"ip"`
	Proxy       bool      `json:"proxy"`
	CreatedDate time.Time `json:"createdDate"`
	Heartbeat   time.Time `json:"heartbeat"`
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult int

const (
	LoginResultSuccess          LoginResult = iota // 0
	LoginResultBadCredentials                      // 1
	LoginResultEmailNotVerified                    // 2
	LoginResult2FA                                 // 3
	LoginResultChangePassword                      // 4
	LoginResultLockedOrDisabled                    // 5
	LoginResultError                               // 6
)

type LoginResponse struct {
	LoginResult       LoginResult `json:"loginResult"`
	IntermediateToken *string     `json:"intermediateToken,omitempty"`
	Me                *Me         `json:"me,omitempty"`
}
