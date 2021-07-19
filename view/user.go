package view

type Users struct {
	UserID         int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Email          string `gorm:"unique" json:"email"`
	Password       string `json:"password"`
	Role           string `json:"role" `
	CreateDateTime string `json:"time"`
}

type Role struct {
	RoleID   int    `gorm:"primary_key;AUTO_INCREMENT" json:"roleId"`
	RoleName string `gorm:"not null;unique" json:"roleName"`
}

type ContactNumber struct {
	UserID         int    `gorm:"foriegn_key" json:"-"`
	NumberID       int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	NumberType     string `json:"numbertype"`
	ISD_Code       string `json:"isd_code"`
	ContactNum     string `json:"contactnumber" `
	CreateDateTime string `json:"time"`
}

type ContactAddress struct {
	UserID         int    `gorm:"foriegn_key" json:"-"`
	AddressID      int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	AddressLine1   string `json:"addressline1"`
	AddressLine2   string `json:"addressline2"`
	City           string `json:"city"`
	State          string `json:"state"`
	Postalcode     string `json:"postalcode" `
	CreateDateTime string `json:"time"`
}
