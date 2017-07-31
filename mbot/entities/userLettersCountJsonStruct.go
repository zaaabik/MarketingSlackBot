package entities

type UserLettersCount struct {
	HostId       string `json:"host_id"`
	Provider     string `json:"provider"`
	LettersCount string `json:"lettersCount"`
}

type UserSendGrid struct {
	HostId   string `json:"host_id"`
	Provider string `json:"provider"`
	Email    string `json:"email"`
}
