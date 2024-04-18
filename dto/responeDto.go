package dto

type ResponseUserProfile struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Bio        string `json:"bio"`
	ProfilePic string `json:"profile_pic"`
}
