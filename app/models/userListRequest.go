package models

type UserListRequest struct {
	sellerID int64  `json:"seller_id"`
	offerID  int64  `json:"offer_id"`
	name     string `json:"name"`
}
