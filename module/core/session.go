package core

import (
	"net/http"
	"strings"
	"time"
)

const SessionPrefix = "SESSION:"

type Session struct {
	Token         string    `json:"token"`
	Email         string    `json:"email"`
	IP            string    `json:"ip"`
	Proxy         bool      `json:"proxy"`
	VaultUnlocked bool      `json:"vaultUnlocked"`
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult int

const (
	LoginResultSuccess LoginResult = iota
	LoginResultBadCredentials
	LoginResult2FA
	LoginResultChangePassword
	LoginResultLockedOrDisabled
	LoginResultError
)

const (
	// TODO: make max login attempts before lock configurable
	maxLoginAttempts    = 5
	failedLoginDuration = 20 * time.Hour
)

func (req LoginRequest) Do(db Providers, ip string) (Session, LoginResult, error) {
	exists, err := db.Users.Exists(req.Email)
	if err != nil {
		// error checking if user exists
		return Session{}, LoginResultError, err
	}
	user, err := db.Users.Get(req.Email)
	if err != nil && exists {
		// error getting user, but user exists
		return Session{}, LoginResultError, err
	}

	// if user cannot log in, respond with this information before checking credentials
	// (otherwise bruteforce attempts on a locked account could be possible)
	if !user.CanLogin() {
		return Session{}, LoginResultLockedOrDisabled, nil
	}

	// if user doesn't exist or the password is incorrect
	if !exists {
		return Session{}, LoginResultBadCredentials, nil
	}
	if !user.PasswordHash.Check(req.Password) {
		// if it's been longer than required since the last failure, reset failures to 1
		if user.LastFailedLogin.Add(failedLoginDuration).Before(time.Now()) {
			user.LoginAttempts = 1
		} else {
			user.LoginAttempts++
		}
		// if user has had to many failed attempts, lock account
		if user.LoginAttempts > maxLoginAttempts {
			user.Locked = true
		}
		user.LastFailedLogin = time.Now()
		return Session{}, LoginResultBadCredentials, nil
	}

	// check for 2FA required

	// check for password change required

	// success, no further steps required
	session := db.Sessions.GenerateFor(user.Email, ip)
	return
}

//TODO: save user when credentials are bad, and save session on success
