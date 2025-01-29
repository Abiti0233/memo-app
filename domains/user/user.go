package user

import (
	"time"
)

type User struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	EmailLowerCase string     `json:"-"`
	EmailVerified  bool `json:"emailVerified.omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
