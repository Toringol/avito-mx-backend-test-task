package models

import "mime/multipart"

type Task struct {
	TaskID   int64
	SellerID string
	Files    map[string][]*multipart.FileHeader
}
