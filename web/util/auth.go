package util

import (
	"crypto/rand"

	"fmt"
	"io"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Store - secure cookie store
var Store = sessions.NewCookieStore(
	[]byte(securecookie.GenerateRandomKey(64)), //Signing key
	[]byte(securecookie.GenerateRandomKey(32)))

func init() {
	Store.Options.HttpOnly = true
	Store.MaxAge(3600 * 24) // max age is 24 hours of log tailing
}

// GenerateSecureKey - Key for CSRF Tokens
func GenerateSecureKey() string {
	// Inspired from gorilla/securecookie
	k := make([]byte, 32)
	io.ReadFull(rand.Reader, k)
	return fmt.Sprintf("%x", k) /////////
}
