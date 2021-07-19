package view

//used for invoice data's
type Invoice struct {
	InvoiceID    int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	ProductID    int    `gorm:"foriegn_key" json:"-"`
	ProductName  string `json:"productname"`
	Quantity     int    `json:"quantity"`
	PricePerUnit int    `json:"price_per_unit" `
	WholePrice   int    `json:"whole_price" `
	PurchaseDate string `json:"time"`
	Discount     int    `json:"discount"`
	TotalPrice   int    `json:"total_price"`
}
