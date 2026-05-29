package dto

type SetUserRolesRequest struct {
	UserID  string   `json:"userId" validate:"required,uuid"`
	RoleIDs []string `json:"roles"`
}
