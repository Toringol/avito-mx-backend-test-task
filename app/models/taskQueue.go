package models

import "mime/multipart"

// Task is model for taskQueue
// When user loads files in running main goroutine, task adds
// to taskQueue futher we can process all tasks concurrently
// swagger:model Task
type Task struct {
	TaskID   int64
	SellerID int64
	Files    map[string][]*multipart.FileHeader
}
