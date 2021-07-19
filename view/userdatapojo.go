package view

//for handling user update
type UserDataPojo struct {
	Client         Client         `json:"client"`
	ContactNumber  ContactNumber  `json:"contactnumbers"`
	ContactAddress ContactAddress `json:"contactaddresses"`
}

type BranchDataPojo struct {
	Branch         Branch         `json:"branch"`
	ContactNumber  ContactNumber  `json:"contactnumbers"`
	ContactAddress ContactAddress `json:"contactaddresses"`
}
