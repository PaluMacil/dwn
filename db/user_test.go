package db_test

import (
	"strings"
	"testing"

	"github.com/PaluMacil/dwn/db"
	"github.com/dgraph-io/badger"
)

func TestGetFakeUser(t *testing.T) {
	_, err := Db.Users.Get("123@example.com")
	if !strings.Contains(err.Error(), badger.ErrKeyNotFound.Error()) {
		t.Fatalf(`Expected '%s', but got '%s' error`, badger.ErrKeyNotFound, err)
	}
}

func TestGetUserNoEmail(t *testing.T) {
	_, err := Db.Users.Get("")
	const requireEmailErrFragment = "requires an email"
	isRequireEmailErr := strings.Contains(err.Error(), requireEmailErrFragment)
	if !isRequireEmailErr {
		t.Fatalf("Expected error containing text '%s' while getting fake user but got '%s' error", requireEmailErrFragment, err)
	}
}

func TestCreateUser(t *testing.T) {
	user := db.User{
		Email:       "exists@example.com",
		DisplayName: "Tester",
		GivenName:   "Testy",
		FamilyName:  "Testing",
	}
	err := Db.Users.Set(user)
	if err != nil {
		t.Fatalf("Error creating user: %s", err)
	}
}

func TestUserExists(t *testing.T) {
	user := db.User{
		Email:       "exists@example.com",
		DisplayName: "Tester",
		GivenName:   "Testy",
		FamilyName:  "Testing",
	}
	err := Db.Users.Set(user)
	if err != nil {
		t.Fatalf("Error creating user: %s", err)
	}
	Db.Users.Exists("exists@example.com")
}

func TestGetExistingUser(t *testing.T) {
	const existingEmail = "exists@example.com"
	const testerDisplayName = "Mr. Tester"
	user := db.User{
		Email:       existingEmail,
		DisplayName: testerDisplayName,
		GivenName:   "Testy",
		FamilyName:  "Testing",
	}
	err := Db.Users.Set(user)
	if err != nil {
		t.Fatalf("Error creating user: %s", err)
	}
	u, err := Db.Users.Get(existingEmail)
	if err != nil {
		t.Fatalf("Error getting user: %s", err)
	}
	if u.DisplayName != testerDisplayName {
		t.Fatalf(`DisplayName for retrieved user was "%s" but expected "%s"`,
			u.DisplayName, testerDisplayName)
	}
}

func TestUserCount(t *testing.T) {
	err := Db.Users.PurgeAll()
	if err != nil {
		t.Fatalf("Could not purge users: %s", err)
	}
	c1, err := Db.Users.Count()
	if err != nil {
		t.Fatalf("Could not count users: %s", err)
	}

	user1 := db.User{
		Email:       "exists1@example.com",
		DisplayName: "Tester",
		GivenName:   "Testy",
		FamilyName:  "Testing",
	}
	err = Db.Users.Set(user1)
	if err != nil {
		t.Fatalf("Error creating user: %s", err)
	}
	user2 := db.User{
		Email:       "exists2@example.com",
		DisplayName: "Tester",
		GivenName:   "Testy",
		FamilyName:  "Testing",
	}
	err = Db.Users.Set(user2)
	if err != nil {
		t.Fatalf("Error creating user: %s", err)
	}
	user3 := db.User{
		Email:       "exists3@example.com",
		DisplayName: "Tester",
		GivenName:   "Testy",
		FamilyName:  "Testing",
	}
	err = Db.Users.Set(user3)
	if err != nil {
		t.Fatalf("Error creating user: %s", err)
	}

	c2, err := Db.Users.Count()
	if err != nil {
		t.Fatalf("Could not count users: %s", err)
	}

	if c2-c1 != 3 {
		t.Fatalf("Added three users but user account did not increase by three")
	}
}
