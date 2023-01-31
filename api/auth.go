package api

type SignUPRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
