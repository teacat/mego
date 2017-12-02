package mego

import (
	"fmt"
	"time"
)

// Context 是單次請求中的上下文建構體。
type Context struct {
	// Keys 存放透過 `Set` 儲存的鍵值組，僅限於單次請求。
	Keys map[string]interface{}
	// Errors 存放開發者自訂的錯誤，可用在中介軟體或處理函式中。
	Errors []error
	// Params 呈現了本次接收到的參數（如果客戶端傳送的內容是陣列而不是物件）。
	Params Params
	// Data 呈現了本次接收到的資料（如果客戶端傳送的內容是物件而不是一個參數陣列）。
	Data interface{}
	// Session 是產生此連線的客戶端階段建構體。
	Session *Session

	// isAborted 表示了這個請求是不是以已經被終止了。
	isAborted bool
}

// Error 能夠將發生的錯誤保存到單次 Session 中。
func (c *Context) Error(err error) *Context {
	c.Errors = append(c.Errors, err)
	return c
}

// Copy 會複製一份 `Context` 供你在 Goroutine 中操作不會遇上資料競爭與衝突問題。
func (c *Context) Copy() *Context {
	ctx := *c
	return &ctx
}

// IsAborted 會回傳一個此請求是否已被終止的布林值。
func (c *Context) IsAborted() bool {
	return c.isAborted
}

// Abort 終止此次請求，避免繼續執行。
func (c *Context) Abort() error {
	c.isAborted = true
	return nil
}

// AbortWithStatus 終止此次請求，避免繼續執行。並以指定的狀態碼回應特定客戶端。
func (c *Context) AbortWithStatus(code int) error {
	c.isAborted = true
	return nil
}

// AbortWithRespond 終止此次請求，避免繼續執行。並以指定的狀態碼、資料回應特定的客戶端。
func (c *Context) AbortWithRespond(code int, result interface{}) error {
	c.isAborted = true
	return nil
}

// AbortWithError 結束此次請求，避免繼續執行。並以指定的狀態碼、錯誤資料與訊息回應特定的客戶端並表示錯誤發生。
func (c *Context) AbortWithError(code int, err interface{}, msg string) error {
	c.isAborted = true
	return nil
}

// Respond 會以指定的狀態碼、資料回應特定的客戶端。
func (c *Context) Respond(code int, result interface{}) error {
	return nil
}

// Status 會回傳指定的狀態碼給客戶端。
func (c *Context) Status(code int) error {
	return nil
}

// RespondWithError 會以指定的狀態碼、錯誤資料與訊息回應特定的客戶端並表示錯誤發生。
func (c *Context) RespondWithError(code int, err interface{}, msg string) error {
	return nil
}

// Bind 能夠將接收到的資料映射到本地建構體。
func (c *Context) Bind(dest interface{}) error {
	return nil
}

// MustBind 和 `Bind` 相同，差異在於若映射失敗會呼叫 `panic` 並阻止此次請求繼續執行。
func (c *Context) MustBind(dest interface{}) error {
	panic(err)
}

// Set 會在本次的 Session 中存放指定的鍵值組內容，可供下次相同客戶端呼叫時存取。
func (c *Context) Set(key string, value interface{}) {
	c.Keys[key] = value
}

// MustGet 和 `Get` 相同，但沒有該鍵值組時會呼叫 `panic`。
func (c *Context) MustGet(key string) interface{} {
	if v, ok := c.Get(key); ok {
		return v
	}
	panic(fmt.Sprintf("Key %s does not exist", key))
}

// Get 會取得先前以 `Set` 存放在本次 Session 中的指定鍵值組。
func (c *Context) Get(key string) (v interface{}, ok bool) {
	v, ok = c.Keys[key]
	return
}

// GetBool 能夠以布林值取得指定的參數。
func (c *Context) GetBool(key string) (v bool) {
	if val, ok := c.Get(key); ok && val != nil {
		v, _ = val.(bool)
	}
	return
}

// GetFloat64 能夠以浮點數取得指定的參數。
func (c *Context) GetFloat64(key string) (v float64) {
	if val, ok := c.Get(key); ok && val != nil {
		v, _ = val.(float64)
	}
	return
}

// GetInt 能夠以正整數取得指定的參數。
func (c *Context) GetInt(key string) (v int) {
	if val, ok := c.Get(key); ok && val != nil {
		v, _ = val.(int)
	}
	return
}

// GetString 能夠以字串取得指定的參數。
func (c *Context) GetString(key string) (v string) {
	if val, ok := c.Get(key); ok && val != nil {
		v, _ = val.(string)
	}
	return
}

// GetStringMap 能夠以 `map[string]interface{}` 取得指定的參數。
func (c *Context) GetStringMap(key string) (v map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		v, _ = val.(map[string]interface{})
	}
	return
}

// GetMapString 能夠以 `map[string]string` 取得指定的參數。
func (c *Context) GetMapString(key string) (v map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		v, _ = val.(map[string]string)
	}
	return
}

// GetStringSlice 能夠以字串切片取得指定的參數。
func (c *Context) GetStringSlice(key string) (v []string) {
	if val, ok := c.Get(key); ok && val != nil {
		v, _ = val.([]string)
	}
	return
}

// GetTime 能夠以時間取得指定的參數。
func (c *Context) GetTime(key string) (v time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		v, _ = val.(time.Time)
	}
	return
}

// GetDuration 能夠以時間長度取得指定的參數。
func (c *Context) GetDuration(key string) (v time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		v, _ = val.(time.Duration)
	}
	return
}

// Next 來呼叫下一個方法處理函式，常用於中介軟體中。
func (c *Context) Next() {

}

// Emit 會向此客戶端廣播一個事件。
func (c *Context) Emit(event string, payload interface{}) error {
	return nil
}

// EmitOthers 會將帶有指定資料的事件向自己以外的所有人廣播。
func (c *Context) EmitOthers(event string, payload interface{}) error {
	return nil
}
