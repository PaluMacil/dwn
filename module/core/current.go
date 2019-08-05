package core

import (
	"fmt"
	"github.com/PaluMacil/dwn/webserver/errs"
	"net/http"
)

type Current struct {
	User    UserInfo  `json:"user"`
	Session Session   `json:"session"`
	db      Providers `json:"-"`
}

var AnonymousCurrent = Current{}

func (c Current) String() string {
	return fmt.Sprintf("Current[Name:%s Email:%s]", c.User.DisplayName, c.User.PrimaryEmail)
}

func (c *Current) LogString() string {
	if c == nil {
		return "\tCurrent[none]\n"
	}
	return fmt.Sprintf("\t%s\n", c)
}

func GetCurrent(token string, db Providers) (Current, error) {
	if token == "" {
		return AnonymousCurrent, nil
	}
	session, err := db.Sessions.Get(token)
	if err != nil {
		return AnonymousCurrent, err
	}
	user, err := db.Users.Get(session.UserID)
	if err != nil {
		return AnonymousCurrent, err
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

func (c Current) Is(groupName string) (bool, error) {
	if c.Anonymous() {
		return false, nil
	}
	return c.db.UserGroups.Exists(c.User.ID, groupName)
}

func (c Current) Authenticated() bool {
	return !c.Anonymous()
}

func (c Current) Anonymous() bool {
	return c.User.ID == 0
}
