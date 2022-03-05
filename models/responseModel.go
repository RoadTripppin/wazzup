package models

type ErrResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
	Token   string `json:"token"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
	Token   string `json:"token"`
}
