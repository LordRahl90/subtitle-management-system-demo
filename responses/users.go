package responses

import (
	"time"
)

// User response DTO for users
type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	UserType  string     `json:"user_type"`
	Token     string     `json:"token"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
