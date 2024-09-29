package types

type UserDto struct {
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Pic         string   `json:"pic"`
	Permissions []string `json:"permissions"`
}
