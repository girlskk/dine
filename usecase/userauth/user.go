package userauth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthToken struct {
	ID uuid.UUID
	jwt.RegisteredClaims
}
