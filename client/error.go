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
