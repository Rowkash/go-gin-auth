package auth

type BaseAuth struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type LoginDto struct {
	BaseAuth
}

type RegisterDto struct {
	BaseAuth
	Name            string `json:"name" binding:"required,min=2"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}
