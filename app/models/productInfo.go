package models

// ProductInfo - DB model description of product
// swagger:model ProductInfo
type ProductInfo struct {
	SellerID  int64   `json:"seller_id"`
	OfferID   int64   `json:"offer_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int64   `json:"quantity"`
	Available bool    `json:"available"`
}
