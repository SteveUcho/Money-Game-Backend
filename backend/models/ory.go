package models

import ory "github.com/ory/kratos-client-go"

type User struct {
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
