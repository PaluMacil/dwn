package dwn

import (
	"net/http"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.FormValue("email")
	plainPassword := r.FormValue("password")
	var user User
	err := Db.One("email", email, &user)
	if err != nil {
		http.Error(w, "incorrect email or password", http.StatusUnauthorized)
	}
	if !CheckPasswordHash(plainPassword, user.Password) {
		//TODO: create session, return token
	} else {
		http.Error(w, "incorrect email or password", http.StatusUnauthorized)
	}
}

const (
	_         = iota
	RoleAdmin = iota
	RoleUser  = iota
)

type User struct {
	ID        int    `storm:"id,increment"`
	Role      int    `storm:"index"`
	Email     string `storm:"unique"`
	Password  string
	Name      string
	CreatedAt time.Time
}

type Session struct {
	Token     uuid.UUID `storm:"id"`
	User      User      `storm:"index"`
	CreatedAt time.Time
	HeartBeat time.Time
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
