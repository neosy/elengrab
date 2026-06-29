package dto

type SetUserRolesRequest struct {
	UserID  string   `json:"userId" validate:"required"`
	RoleIDs []string `json:"roles"`
}
