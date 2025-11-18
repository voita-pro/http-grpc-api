package domain

type TokenData struct {
	UID           string `json:"uid"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	DisplayName   string `json:"display_name"`
	Roles         string `json:"roles"`
}
type Country struct {
	Id      int32
	Title   string
	Code    string
	IsoCode *string
}
type Login struct {
	Email    string
	Password string
}
