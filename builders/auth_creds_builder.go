package builders

import (
	"github.com/CameronHonis/chess-arbitrator/models"
	"time"
)

type AuthCredsBuilder struct {
	authCreds models.AuthCreds
}

func NewAuthCredsBuilder() *AuthCredsBuilder {
	return &AuthCredsBuilder{}
}

func (b *AuthCredsBuilder) WithClientKey(clientKey models.Key) *AuthCredsBuilder {
	b.authCreds.ClientKey = clientKey
	b.authCreds.CreatedAt = time.Now()
	return b
}

func (b *AuthCredsBuilder) WithPrivateKey(privateKey models.Key) *AuthCredsBuilder {
	b.authCreds.PrivateKey = privateKey
	b.authCreds.PriKeyCreatedAt = time.Now()
	return b
}

func (b *AuthCredsBuilder) WithRole(role models.RoleName) *AuthCredsBuilder {
	b.authCreds.Role = role
	return b
}

func (b *AuthCredsBuilder) FromAuthCreds(authCreds models.AuthCreds) *AuthCredsBuilder {
	b.authCreds = authCreds
	return b
}

func (b *AuthCredsBuilder) Build() *models.AuthCreds {
	return &b.authCreds
}
