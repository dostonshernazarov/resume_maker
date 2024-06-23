package entity

import "time"

type AppVersion struct {
	ID             int64
	AndroidVersion string
	IOSVersion     string
	IsForceUpdate    bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RefreshToken struct {
	GUID string
	RefreshToken string
	ExpiryDate time.Time
	CreatedAt time.Time

}