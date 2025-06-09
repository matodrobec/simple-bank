package token

import "time"

type Maker interface {
	CrateToken(username string, duration time.Duration, role string) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
