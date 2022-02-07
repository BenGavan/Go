package token

import "time"

// for managing tokens
type Maker interface {
	// Creates a new token  for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
