package models

type ProductInfo struct {
	sellerID  int64   `json:"seller_id"`
	offerID   int64   `json:"offer_id"`
	name      string  `json:"name"`
	price     float64 `json:"price"`
	quantity  int64   `json:"quantity"`
	available bool    `json:"available"`
}
