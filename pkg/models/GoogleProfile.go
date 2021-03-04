package models

type GoogleProfile struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func (gp GoogleProfile) ToUser() User {
	return User{
		Title:     gp.Name,
		Email:     gp.Email,
		AvatarUrl: gp.Picture,
		Nickname:  gp.Email,
	}
}
