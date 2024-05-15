package dto

type LoginRes struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}