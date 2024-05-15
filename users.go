package main

import (
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-webauthn/webauthn/webauthn"
)

// User entity.
type User struct {
	ID          uint32 // Immutable
	Email       string // Mutable
	DisplayName string // Nickname

	webAuthnCredentials []webauthn.Credential
}

// for use as cookie. example only.
func (u *User) login() string {
	return fmt.Sprint(u.ID)
}

var _ webauthn.User = (*User)(nil)

func (u *User) WebAuthnID() []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, u.ID)
	return buf
}
func (u *User) WebAuthnName() string {
	return u.Email
}
func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.webAuthnCredentials
}
func (u *User) WebAuthnIcon() string {
	return ""
}

// TODO: concurrency
type Store struct {
	// simulate database rows
	users map[uint32]*User
}

func NewStore() *Store {
	return &Store{
		users: map[uint32]*User{},
	}
}

func (s *Store) AuthRequest(r *http.Request) *User {
	login, _ := r.Cookie(`login`)
	if login == nil {
		return nil
	}
	for _, u := range s.users {
		if u.login() == login.Value {
			return u
		}
	}
	return nil
}

func (s *Store) AddWebAuthnCredentialFor(u *User, credential *webauthn.Credential) {
	u.webAuthnCredentials = append(u.webAuthnCredentials, *credential)
}

func (s *Store) MakeCookie(u *User, w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  `login`,
		Value: u.login(),
		Path:  "/",
	})
}

func (s *Store) RemoveCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  `login`,
		Value: "",
		Path:  "/",
	})
}

func (s *Store) GetUserByID(id uint32) *User {
	return s.users[id]
}

func (s *Store) AddNewUser(email string) (*User, error) {
	for _, u := range s.users {
		if strings.EqualFold(u.Email, email) {
			return nil, errors.New(`email address taken by someone else`)
		}
	}
	u := &User{
		ID:          uint32(len(s.users) + 1), // TODO
		Email:       email,
		DisplayName: "",
	}
	s.users[u.ID] = u
	return u, nil
}
