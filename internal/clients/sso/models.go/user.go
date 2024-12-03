package ssomodels

type RegisterUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password_complexity"`
	Name     string `json:"name,omitempty"`
}

type RegisterResp struct {
	Success bool
}

type LogIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LogInResp struct {
	AccessToken  string
	RefreshToken string
}

type LogInViaTg struct {
	Email string `json:"email" validate:"required,email"`
}

type LogInViaTgResp struct {
	Success bool
	Info    string
}

type CheckOTPAndLogIn struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type CheckOTPAndLogInResp struct {
	AccessToken  string
	RefreshToken string
}

type UpdateUserInfo struct {
	ID      string `json:"id" validate:"required"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	TgLink  string `json:"tg_link,omitempty"`
	IsAdmin bool   `json:"is_admin,omitempty"`
}

type UpdateUserInfoResp struct {
	Success bool
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResp struct {
	AccessToken string
}

type IsAdmin struct {
	UserID string `json:"user_id" validate:"required"`
}

type IsAdminResp struct {
	IsAdmin bool
}

type AuthCheck struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type AuthCheckResp struct {
	IsValid bool
	UserID  string
}
