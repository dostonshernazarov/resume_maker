package models

type RegisterReq struct {
	Fullname string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgetPassReq struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type RegisterRes struct {
	Content string `json:"content"`
}

type BacisResumeRes struct {
	Content      string `json:"content"`
	BasicRedisID string `json:"basic_redis_id"`
	MainRedisID  string `json:"main_redis_id"`
}

type Verify struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type TokenResp struct {
	ID      string `json:"user_id"`
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
	Role    string `json:"role"`
}

type UserResCreate struct {
	Id           string `json:"id"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	ProfileImg   string `json:"image_url"`
	PhoneNumber  string `json:"phone_num"`
	Role         string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ClientRedis struct {
	Fullname string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
