package dwn

import (
	"encoding/json"
	"net/http"
	"time"

	"log"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.FormValue("email")
	plainPassword := r.FormValue("password")
	var user User
	err := Db.One("Email", email, &user)
	if err != nil {
		log.Println("could not load user to build session:", err)
		http.Error(w, "incorrect email or password", http.StatusUnauthorized)
		return
	}
	if CheckPasswordHash(plainPassword, user.Password) {
		session := Session{
			Token:     uuid.NewV4(),
			User:      user,
			CreatedAt: time.Now(),
			HeartBeat: time.Now(),
		}
		err = Db.Save(&session)
		if err != nil {
			log.Println("could not save session:", err)
			http.Error(w, "could not save session", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(session)
	} else {
		log.Println("incorrect password:", plainPassword, "for user:", user)
		http.Error(w, "incorrect email or password", http.StatusUnauthorized)
		return
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
