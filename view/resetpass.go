package view

type Resetpass struct {
	Email           string `gorm:"unique" json:"email"`
	OldPassword     string `json:"oldpassword"`
	NewPassword     string `json:"newpassword"`
	ConfirmPassword string `json:"confirmpassword"`
}
