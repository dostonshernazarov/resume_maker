package models

type UserReq struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DateOfBirth string `json:"date_of_birth"`
	Card        string `json:"card"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
}

type UserRes struct {
	Id           string `json:"id"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	DateOfBirth  string `json:"date_of_birth"`
	ProfileImg   string `json:"profile_img"`
	Card         string `json:"card"`
	Gender       string `json:"gender"`
	PhoneNumber  string `json:"phone_number"`
	Role         string `json:"role"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	DeletedAt    string `json:"deleted_at"`
}

type ListUsersRes struct{}

type Users struct {
	Users []*UserRes `json:"users"`
}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type FieldValues struct {
	Column string `json:"column"`
	Value  string `json:"value"`
}
