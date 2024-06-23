package entity

import "time"

type User struct {
	GUID        string
	Name        string
	Image       string
	Email       string
	PhoneNumber string
	Password    string
	Refresh     string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Users struct {
	Users []*User
	Total uint64
}

type IsUnique struct {
	Email string
}

type UpdateRefresh struct {
	UserID       string
	RefreshToken string
}

type UpdatePassword struct {
	UserID      string
	NewPassword string
}

type Response struct {
	Status bool
}
