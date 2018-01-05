package client

import (
	"errors"

	"github.com/vmihailenco/msgpack"
)

var (
	// ErrChunkWithFiles 表示不能上傳區塊檔案同時夾帶其他檔案實體。
	ErrChunkWithFiles = errors.New("mego: cannot send the file chunks with the file entities")
	// ErrTimeout 表示一個請求過久都沒有回應。
	ErrTimeout = errors.New("mego: request timeout")
	// ErrClosed 表示正在使用一個已經關閉的連線。
	ErrClosed = errors.New("mego: use of closed connection")
	// ErrSent 表示正在使用一個早已發送出去的請求。
	ErrSent = errors.New("mego: use of sent request")
	// ErrSubscriptionRefused 表示欲訂閱的事件請求被拒。
	ErrSubscriptionRefused = errors.New("mego: the event subscription was refused")
	// ErrAborted 表示請求已被終止。
	ErrAborted = errors.New("mego: the request has been aborted")
)

const (
	// StatusError 表示有內部錯誤發生。
	StatusError = -1000
	// StatusFull 表示此請求無法被接受，因為額度已滿。例如：使用者加入了一個已滿的聊天室、好友清單已滿。
	StatusFull = -1001
	// StatusExists 表示請求的事物已經存在，例如：重複的使用者名稱、電子郵件地址。
	StatusExists = -1002
	// StatusInvalid 表示此請求格式不正確。
	StatusInvalid = -1003
	// StatusNotFound 表示找不到請求的資源。
	StatusNotFound = -1004
	// StatusNotAuthorized 表示使用者需要登入才能進行此請求。
	StatusNotAuthorized = -1005
	// StatusNoPermission 表示使用者已登入但沒有相關權限執行此請求。
	StatusNoPermission = -1006
	// StatusUnimplemented 表示此功能尚未實作完成。
	StatusUnimplemented = -1007
	// StatusTooManyRequests 表示使用者近期發送太多請求，需稍後再試。
	StatusTooManyRequests = -1008
	// StatusResourceExhausted 表示使用者可用的額度已耗盡。
	StatusResourceExhausted = -1009
	// StatusBusy 表示伺服器正繁忙無法進行執行請求。
	StatusBusy = -1010
	// StatusFileRetry 表示檔案區塊發生錯誤，需要重新上傳相同區塊。
	StatusFileRetry = -1011
	// StatusFileEmpty 表示上傳的檔案、區塊是空的。
	StatusFileEmpty = -1012
	// StatusFileTooLarge 表示檔案過大無法上傳。
	StatusFileTooLarge = -1013
	// StatusTimeout 表示這個請求費時過久導致逾期而被終止。
	StatusTimeout = -1014
)

// Error 呈現了一個遠端所傳回的錯誤。
type Error struct {
	// Code 是錯誤代號。
	Code int `codec:"c" msgpack:"c"`
	// Message 是人類可讀的簡略錯誤訊息。
	Message string `codec:"m" msgpack:"m"`
	// Data 是錯誤的詳細資料。基於 Message Pack 格式需要透過 `Bind` 映射到本地建構體。
	Data []byte `codec:"d" msgpack:"d"`
}

// Error 可以取得錯誤中的文字敘述訊息。
func (e Error) Error() string {
	return e.Message
}

// Bind 能夠將錯誤的詳細資料映射到本地建構體。
func (e Error) Bind(dest interface{}) error {
	return msgpack.Unmarshal(e.Data, dest)
}
