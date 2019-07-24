package core

import (
	"log"
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
	LoginResultEmailNotVerified
	LoginResult2FA
	LoginResultChangePassword
	LoginResultLockedOrDisabled
	LoginResultError
)

const (
	// TODO: make max login attempts before lock configurable
	// TODO: save failures to database and require captcha for suspect / repeat offender IPs
	maxLoginAttempts    = 5
	failedLoginDuration = 20 * time.Hour
)

// Do executes a login request. It returns a user, session (possibly both empty) and
// a LoginResult. The error will be nil if (and only if) the result is LoginResultError
func (req LoginRequest) Do(db Providers, ip string) (UserInfo, Session, LoginResult, error) {
	exists, err := db.Users.EmailExists(req.Email)
	if err != nil {
		// error checking if user exists
		return UserInfo{}, Session{}, LoginResultError, err
	}
	usersWithEmail, err := db.Users.WithEmail(req.Email)
	if err != nil && exists {
		// error getting user, but user exists
		return UserInfo{}, Session{}, LoginResultError, err
	}

	// if user doesn't exist (same response as the password is incorrect)
	if !exists {
		// TODO: count non-existent user attempts towards suspicion score of an IP
		return UserInfo{}, Session{}, LoginResultBadCredentials, nil
	}

	// Check whether any of the users with this email have verified it.
	var user User
	noVerifiedUser := true
	for _, u := range usersWithEmail {
		for _, email := range u.Emails {
			// Check if matched AND verified.
			if email.Email == req.Email && email.Verified {
				user = u
				noVerifiedUser = false
			}
		}
	}
	if noVerifiedUser {
		return UserInfo{}, Session{}, LoginResultEmailNotVerified, nil
	}

	// if user cannot log in, respond with this information before checking credentials
	// (otherwise bruteforce attempts on a locked account could be possible)
	if !user.CanLogin() {
		return UserInfo{}, Session{}, LoginResultLockedOrDisabled, nil
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
		err := db.Users.Set(user)
		if err != nil {
			log.Println("saving user with failed login attempt:", err.Error())
		}
		return UserInfo{}, Session{}, LoginResultBadCredentials, nil
	}

	// TODO: check for 2FA required

	// TODO: check for password change required

	// success, no further steps required
	session := db.Sessions.GenerateFor(user.ID, ip)
	err = db.Sessions.Set(session)
	if err != nil {
		return UserInfo{}, Session{}, LoginResultError, err
	}
	return user.Info(), session, LoginResultSuccess, nil
}
