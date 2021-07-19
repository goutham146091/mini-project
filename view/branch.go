package view

//used for branch data's
type Branch struct {
	UserID         int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	ClientID       int    `gorm:"foriegn_key" json:"-"`
	Email          string `gorm:"unique" json:"email"`
	BranchName     string `json:"branchname"`
	BranchLocation string `json:"branchlocation"`
	Role           string `json:"role" `
	CreateDateTime string `json:"time"`
	InviteCode     int    `json:"invitecode"`
	Is_Registered  bool   `json:"isregistered"`
	Is_Deleted     bool   `json:"isdeleted"`
	Is_Active      bool   `json:"isactive"`
}
