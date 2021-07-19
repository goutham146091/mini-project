package view

//used for set password data's
type EmailPassword struct {
	Email           string `gorm:"unique" json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmpassword"`
	Role            string `json:"role" `
	InviteCode      int    `json:"invitecode"`
}
