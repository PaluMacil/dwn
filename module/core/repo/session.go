package repo

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/PaluMacil/dwn/module/core"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/database/store"
	"github.com/gofrs/uuid"
)

type SessionRepo struct {
	store database.Storer
	db    *database.Database
}

func NewSessionRepo(store database.Storer, db *database.Database) *SessionRepo {
	return &SessionRepo{store, db}
}

func (p SessionRepo) Get(token string) (core.Session, error) {
	var session = core.Session{Token: token}
	if token == "" {
		return session, errors.New("SessionRepo.Get requires a token but got an empty string")
	}
	item, err := p.store.Get(&session)
	if err != nil {
		return session, err
	}
	session, ok := item.(core.Session)
	if !ok {
		return session, fmt.Errorf("got data of type %T but wanted core.Session", session)
	}
	return session, err
}

// Exists checks whether a specific core.Session exists
func (p SessionRepo) Exists(token string) (bool, error) {
	_, err := p.Get(token)
	if p.db.IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

// Set saves a core.Session to the database
func (p SessionRepo) Set(session core.Session) error {
	return p.store.Set(&session)
}

// GenerateFor makes a new core.Session for the specified userID but does
// not check if the id exists or persist anything to the database.
func (p SessionRepo) GenerateFor(userID store.Identity, ip string) core.Session {
	t := uuid.Must(uuid.NewV4())
	session := core.Session{
		Token:       t.String(),
		UserID:      userID,
		IP:          ip,
		CreatedDate: time.Now(),
		Heartbeat:   time.Now(),
	}
	return session
}

func (p SessionRepo) All() ([]core.Session, error) {
	var items []database.Item
	pfx := core.Session{}.Prefix()
	err := p.store.All(pfx, &items, true)
	sessions := make([]core.Session, len(items))
	for i, v := range items {
		sessions[i] = v.(core.Session)
	}

	return sessions, err
}

func (p SessionRepo) Delete(token string) error {
	return p.store.Delete(core.Session{Token: token})
}

func (p SessionRepo) PurgeAll() error {
	sessions, err := p.All()
	if err != nil {
		return err
	}
	for _, s := range sessions {
		err = p.Delete(s.Token)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p SessionRepo) UpdateHeartbeat(session *core.Session, ip string) error {
	session.Heartbeat = time.Now()
	session.IP = ip
	return p.Set(*session)
}

const (
	// TODO: make max login attempts before lock configurable
	// TODO: save failures to database and require captcha for suspect / repeat offender IPs
	maxLoginAttempts    = 5
	failedLoginDuration = 20 * time.Hour
)

// DoLogin executes a login request. It returns a user, session (possibly both empty) and
// a LoginResult. The error will be nil if (and only if) the result is LoginResultError.
func (p SessionRepo) DoLogin(req core.LoginRequest, ip string) (core.UserInfo, core.Session, core.LoginResult, error) {
	exists, err := p.db.Users.EmailExists(req.Email)
	if err != nil {
		// error checking if user exists
		return core.UserInfo{}, core.Session{}, core.LoginResultError, err
	}
	usersWithEmail, err := p.db.Users.WithEmail(req.Email)
	if err != nil && exists {
		// error getting user, but user exists
		return core.UserInfo{}, core.Session{}, core.LoginResultError, err
	}

	// if user doesn't exist, send same response as the password is incorrect
	if !exists {
		// TODO: count non-existent user attempts towards suspicion score of an IP
		return core.UserInfo{}, core.Session{}, core.LoginResultBadCredentials, nil
	}

	// Check whether any of the users with this email have verified it.
	var user core.User
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
		return core.UserInfo{}, core.Session{}, core.LoginResultEmailNotVerified, nil
	}

	// if user cannot log in, respond with this information before checking credentials
	// (otherwise bruteforce attempts on a locked account could be possible)
	if !user.CanLogin() {
		return core.UserInfo{}, core.Session{}, core.LoginResultLockedOrDisabled, nil
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
		err := p.db.Users.Set(user)
		if err != nil {
			log.Println("saving user with failed login attempt:", err.Error())
		}
		return core.UserInfo{}, core.Session{}, core.LoginResultBadCredentials, nil
	}

	// TODO: check for 2FA required

	// TODO: check for password change required

	// success, no further steps required
	session := p.db.Sessions.GenerateFor(user.ID, ip)
	err = p.db.Sessions.Set(session)
	if err != nil {
		return core.UserInfo{}, core.Session{}, core.LoginResultError, err
	}
	return user.Info(), session, core.LoginResultSuccess, nil
}
