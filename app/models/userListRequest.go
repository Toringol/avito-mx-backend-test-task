package models

// UserListRequest is request for searching specific products
// by info in request
// swagger:model UserListRequest
type UserListRequest struct {
	SellerID int64  `json:"seller_id"`
	OfferID  int64  `json:"offer_id"`
	Name     string `json:"name"`
}
