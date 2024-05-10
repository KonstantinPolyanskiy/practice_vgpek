package rest

import "practice_vgpek/internal/model/domain"

type Token struct {
	Token string      `json:"token"`
	Role  domain.Role `json:"role"`
}

func (t Token) TokenToResponse(token string, role domain.Role) Token {
	return Token{
		Token: token,
		Role:  role,
	}
}
