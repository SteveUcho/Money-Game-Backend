package models

import (
	"github.com/google/uuid"
	ory "github.com/ory/kratos-client-go"
)

type User struct {
	ID      uuid.UUID
	Session *ory.Session
	Traits  IdentityTraits
}

type IdentityTraits struct {
	Email string `json:"email"`
	Name  struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
}
