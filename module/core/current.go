package core

import (
	"fmt"
	"net/http"
	"time"
)

type Current struct {
	User    UserInfo  `json:"user"`
	Session Session   `json:"session"`
	db      Providers `json:"-"`
}

func (c Current) String() string {
	return fmt.Sprintf("Current[Name:%s Email:%s]", c.User.DisplayName, c.User.Email)
}

func (c *Current) LogString() string {
	if c == nil {
		return "\tCurrent[none]\n"
	}
	return fmt.Sprintf("\t%s\n", c)
}

func GetCurrent(r *http.Request, db Providers) (*Current, error) {
	token := r.Header.Get("dwn-token")
	if token == "" {
		return nil, nil
	}
	session, err := db.Sessions.Get(token)
	if err != nil {
		return nil, err
	}
	session.Heartbeat = time.Now()
	session.IP = IP(r)
	if err = db.Sessions.Set(session); err != nil {
		return nil, err
	}
	user, err := db.Users.Get(session.Email)
	if err != nil {
		return nil, err
	}
	return &Current{
		User:    user.Info(),
		Session: session,
		db:      db,
	}, nil
}

// Can asks if a user can do something. It returns whether a user is in a group with
// the specified permission. Admins always return true because they can do anything.
func (c *Current) Can(permission string) (bool, error) {
	if c == nil {
		return false, nil
	}
	groups, err := c.db.Groups.GroupsFor(c.User.Email)
	if err != nil {
		return false, err
	}
	for _, g := range groups {
		if g.Name == BuiltInGroupAdmin || g.HasPermission(permission) {
			return true, nil
		}
	}
	return false, nil
}

func (c *Current) Is(groupName string) (bool, error) {
	if c == nil {
		return false, nil
	}
	return c.db.UserGroups.Exists(c.User.Email, groupName)
}

func (c *Current) Authenticated() bool {
	return c != nil
}
