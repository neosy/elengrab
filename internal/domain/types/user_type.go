package dtypes

type UserType uint8

const (
	UserTypeAnonymous UserType = iota
	UserTypeGuest
	UserTypeUser
	UserTypeAdmin
)
