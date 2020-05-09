package core

import (
	"fmt"
	"github.com/PaluMacil/dwn/webserver/errs"
	"net/http"
)

// Current holds details about the current user. The zero value is an anonymous user.
type Current struct {
	User    UserInfo  `json:"user"`
	Session Session   `json:"session"`
	db      Providers `json:"-"`
}

var anonymousCurrent = Current{}

func (c Current) String() string {
	return fmt.Sprintf("Current[Name:%s Email:%s]", c.User.DisplayName, c.User.PrimaryEmail)
}

func (c *Current) LogString() string {
	if c == nil {
		return "\tCurrent[none]\n"
	}
	return fmt.Sprintf("\t%s\n", c)
}

// GetCurrent returns details about the current user based upon the token provided. If the token is empty or there is an
// error getting this information, the current values will be for the anonymous user.
func GetCurrent(token string, db Providers) (Current, error) {
	if token == "" {
		return anonymousCurrent, nil
	}
	session, err := db.Sessions.Get(token)
	if err != nil {
		return anonymousCurrent, fmt.Errorf("getting session: %w", err)
	}
	user, err := db.Users.Get(session.UserID)
	if err != nil {
		return anonymousCurrent, fmt.Errorf("getting user: %w", err)
	}
	return Current{
		User:    user.Info(),
		Session: session,
		db:      db,
	}, nil
}

// Can asks if a user can do something. It returns nil if a user is in a group with
// the specified permission. Admins always return nil because they can do anything.
// Otherwise can returns an appropriate StatusError.
//
// 	if err := cur.Can(core.PermissionViewAppSettings); err != nil {
//		return err
//	}
func (c Current) Can(permission string) errs.Error {
	if c.Anonymous() {
		return errs.StatusUnauthorized
	}
	groups, err := c.db.Groups.GroupsFor(c.User.ID)
	if err != nil {
		return errs.StatusError{http.StatusInternalServerError, err}
	}
	for _, g := range groups {
		if g.Name == BuiltInGroupAdmin || g.HasPermission(permission) {
			return nil
		}
	}
	return errs.StatusForbidden
}

// Is checks that that current user has a user group membership in the group named
func (c Current) Is(groupName string) (bool, error) {
	if c.Anonymous() {
		return false, nil
	}
	return c.db.UserGroups.Exists(c.User.ID, groupName)
}

// Authenticated returns whether the current user is not anonymous by checking that the user id is non-zero
func (c Current) Authenticated() bool {
	return !c.Anonymous()
}

// Anonymous returns whether the current user ID is 0, the anonymous user
func (c Current) Anonymous() bool {
	return c.User.ID == 0
}
