package entity

type User struct {
	ID       uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Email    string `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password string `gorm:"->;<-;not null" json:"password"`
	Address  string `gorm:"type:varchar(255)" json:"address"`
	NoTelp   int64  `gorm:"type:int(10)" json:"no_telp"`
	Token    string `gorm:"-" json:"token,omitempty"`
	Book     []Book `json:"books"`
}
