package view

//used for client data's
type Client struct {
	UserID         int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Email          string `gorm:"unique" json:"email"`
	ClientName     string `json:"clientname"`
	Role           string `json:"role" `
	CreateDateTime string `json:"time"`
	Shoptype       string `json:"shoptype"`
	NoOfBranches   int    `json:"noofbranches"`
	InviteCode     int    `json:"invitecode"`
	Is_Registered  bool   `json:"isregistered"`
	Is_Deleted     bool   `json:"isdeleted"`
	Is_Active      bool   `json:"isactive"`
}
