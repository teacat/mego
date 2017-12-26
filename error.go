package mego

import "errors"

var (
	// ErrEventNotFound 表示欲發送的事件沒有被初始化或任何客戶端監聽而無法找到因此發送失敗。
	ErrEventNotFound = errors.New("the event doesn't exist")
	// ErrChannelNotFound 表示欲發送的事件存在，但目標頻道沒有被初始化或任何客戶端監聽而無法找到因此發送失敗。
	ErrChannelNotFound = errors.New("the channel doesn't exist")
	// ErrFileNotFound 表示欲取得的檔案並不存在，可能是客戶端上傳不完整。
	ErrFileNotFound = errors.New("the file was not found")
	// ErrKeyNotFound 表示欲從鍵值組中取得的鍵名並不存在。
	ErrKeyNotFound = errors.New("the key was not found")
)

// Error 呈現了一個路由中發生的錯誤。
type Error struct {
	Err      error
	Metadata interface{}
}

// Error 會回傳一個錯誤的代表訊息。
func (e Error) Error() string {
	return e.Err.Error()
}
