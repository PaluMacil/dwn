package dwn

import (
	"encoding/json"
	"net/http"
	"time"

	"log"

	"io/ioutil"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// TokenHandler takes in unencrypted credentials (email and password) via form values,
// hashes the password with bcrypt, compares to the stored hash, and returns Unauthorized
// or a token (uuid version 4) representing the session.
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("could not read login request json body:", err)
		http.Error(w, "could not read login request json body", http.StatusInternalServerError)
		return
	}
	var req TokenRequest
	err = json.Unmarshal(body, &req)

	var user User
	err = Db.One("Email", req.Email, &user)
	if err != nil {
		log.Println("could not load user to build session:", err, "for email", req.Email)
		http.Error(w, "incorrect email or password", http.StatusUnauthorized)
		return
	}
	if CheckPasswordHash(req.Password, user.Password) {
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
		log.Println("incorrect password:", req.Password, "for user:", user.Email)
		http.Error(w, "incorrect email or password", http.StatusBadRequest)
		return
	}
}

const (
	_ = iota
	// RoleAdmin denotes an administrator
	RoleAdmin = iota
	// RoleUser denotes an unpriviledged user
	RoleUser = iota
)

// User represents and application user
type User struct {
	ID        int    `storm:"id,increment"`
	Role      int    `storm:"index"`
	Email     string `storm:"unique"`
	Password  string `json:"-"`
	Name      string
	CreatedAt time.Time
}

// Session represents a user session
type Session struct {
	Token     uuid.UUID `storm:"id"`
	User      User      `storm:"index"`
	CreatedAt time.Time
	HeartBeat time.Time
}

// TokenRequest holds the incoming request for a token (session key)
type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HashPassword takes a plaintext password string and returns a hash from bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares the hash of a plain password with a stored
// hash, returning a bool match result
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// InRole checks to see that an http request comes from a user in a specified role.
func InRole(r *http.Request, role int) (bool, error) {
	token := r.Header.Get("X-DWN-TOKEN")
	var session Session
	if err := Db.One("Token", token, &session); err != nil {
		return false, err
	}
	return session.User.Role == role, nil
}

// LogoutHandler handles requests to log out. If the session exists, it is deleted.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-DWN-TOKEN")
	var session Session
	if err := Db.One("Token", token, &session); err != nil {
		http.Error(w, "Could not find session.", http.StatusBadRequest)
		return
	}
	if err := Db.DeleteStruct(&session); err != nil {
		log.Println("Could not delete session:", err)
		http.Error(w, "Could not delete session.", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
