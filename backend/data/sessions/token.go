package sessions

import "time"

type Token struct {
	UserID     int
	Expiry     time.Time
	Attributes map[string]string
}
