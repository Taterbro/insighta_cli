package api

import "time"

var BackendUrl = "https://backendstagezero.pxxl.click"

type Account struct {
	ID          string    `json:"id" db:"id"`
	GitHubID    int       `json:"github_id" db:"github_id"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	AvatarURL   string    `json:"avatar_url" db:"avatar_url"`
	Role        string    `json:"role" db:"role"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	LastLoginAt time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Credentials struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	UserDetails  Account `json:"user_details"`
}
