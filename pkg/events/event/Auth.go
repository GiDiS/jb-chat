package event

import "jb_chat/pkg/models"

type AuthRegisterRequest struct {
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Title     string `json:"title"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatarUrl"`
}

type AuthSignInRequest struct {
	Nickname    string `json:"nickname,omitempty"`
	Password    string `json:"password,omitempty"`
	Service     string `json:"service,omitempty"`
	AccessToken string `json:"accessToken,omitempty"`
	SecretToken string `json:"secretToken,omitempty"`
	Ttl         int    `json:"ttl,omitempty"`
}

type AuthSignInResponse struct {
	AuthMeResponse
	Token string `json:"token,omitempty"`
}

type AuthSignOutResponse struct {
	ResultStatus
}

type AuthMeResponse struct {
	ResultStatus
	Me *models.UserInfo `json:"me,omitempty"`
}

func (r *AuthMeResponse) SetMe(user *models.UserInfo) {
	r.Ok = true
	r.Me = user
}
