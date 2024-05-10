package rest

type Token struct {
	Token string `json:"token"`
}

func (t Token) TokenToResponse(token string) Token {
	return Token{Token: token}
}
