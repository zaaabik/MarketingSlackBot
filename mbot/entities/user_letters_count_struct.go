package entities

type UserLettersCount struct {
	HostId       string `json:"host_id"`
	Provider     string `json:"provider"`
	LettersCount string `json:"lettersCount"`
	UserId       string `json:"user_id"`
}

type UserSendGrid struct {
	HostId   string `json:"host_id"`
	Provider string `json:"provider"`
	Email    string `json:"email"`
	UserId   string `json:"user_id"`
}
