package dto

type UserDTO struct {
	ID       uint64 `json:"id" form:"id"`
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password,omitempty" form:"password,omitempty"`
	Address  string `json:"address,omitempty" form:"address"`
	NoTelp   int64  `json:"no_telp,omitempty" form:"no_telp"`
}
