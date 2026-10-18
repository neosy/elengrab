package dto

import pservices "github.com/neosy/elengrab/internal/ports/services"

type (
	AuthUserRequest          = pservices.AuthUserRequest
	AuthUserResponse         = pservices.AuthUserResponse
	AuthRegisterUserRequest  = pservices.AuthRegisterUserRequest
	AuthRegisterAdminRequest = pservices.AuthRegisterAdminRequest
)
