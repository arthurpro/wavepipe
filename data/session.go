package data

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"code.google.com/p/go.crypto/pbkdf2"
)

// Session represents an API session for a specific user on wavepipe
type Session struct {
	ID        int    `json:"id"`
	UserID    int    `db:"user_id" json:"userId"`
	Client    string `json:"client"`
	PublicKey string `db:"public_key" json:"publicKey"`
	SecretKey string `db:"secret_key" json:"secretKey"`
}

// NewSession generates and saves a new session for the specified user, with the specified
// client name. This function also randomly generates public and private keys.
func NewSession(userID int, password string, client string) (*Session, error) {
	// Generate session
	session := new(Session)
	session.UserID = userID
	session.Client = client

	// Generate salts for use with PBKDF2
	saltBuf := make([]byte, 16)
	if _, err := rand.Read(saltBuf); err != nil {
		return nil, err
	}
	salt1 := saltBuf

	saltBuf = make([]byte, 16)
	if _, err := rand.Read(saltBuf); err != nil {
		return nil, err
	}
	salt2 := saltBuf

	// Use PBKDF2 to generate a public key and secret key based off the user's password
	session.PublicKey = fmt.Sprintf("%x", pbkdf2.Key([]byte(password), salt1, 4096, 16, sha1.New))
	session.SecretKey = fmt.Sprintf("%x", pbkdf2.Key([]byte(password), salt2, 4096, 16, sha1.New))

	// Save session
	if err := session.Save(); err != nil {
		return nil, err
	}

	return session, nil
}

// Delete removes an existing Session from the database
func (u *Session) Delete() error {
	return DB.DeleteSession(u)
}

// Load pulls an existing Session from the database
func (u *Session) Load() error {
	return DB.LoadSession(u)
}

// Save creates a new Session in the database
func (u *Session) Save() error {
	return DB.SaveSession(u)
}

// ToJSON generates a JSON representation of a Session
func (u Session) ToJSON() ([]byte, error) {
	// Marshal into JSON
	out, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	// Return JSON
	return out, nil
}

// FromJSON generates a Session from its JSON representation
func (u *Session) FromJSON(in []byte) error {
	return json.Unmarshal(in, &u)
}