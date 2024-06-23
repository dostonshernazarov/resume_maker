package models

type UserReq struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

type UserUpdateReq struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type UserRes struct {
	Id           string `json:"id"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Role         string `json:"role"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ListUsersResponse struct{}

type Users struct {
	Users []*UserRes `json:"users"`
	Count uint64     `json:"count"`
}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type FieldValues struct {
	Column string `json:"column"`
	Value  string `json:"value"`
}
