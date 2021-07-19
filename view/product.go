package view

//used for product data's
type Product struct {
	ProductID       int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	BranchID        int    `gorm:"foriegn_key" json:"-"`
	ProductName     string `json:"productname"`
	Quantity        int    `json:"quantity"`
	PricePerUnit    int    `json:"price_per_unit" `
	WholePrice      int    `json:"whole_price" `
	CreateDateTime  string `json:"time"`
	AvailableStatus string `json:"available_status"`
}
