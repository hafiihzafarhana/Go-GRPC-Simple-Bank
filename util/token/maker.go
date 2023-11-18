package token

import (
	"time"
)

// Implementasi interface untuk mengatur token
type Maker interface {
	// Method untuk sign token
	CreateToken(username string, duration time.Duration) (string, error)
	// Verify token
	VerifyToken(token string) (*Payload, error)
}
