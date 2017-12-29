package mego

import (
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack"
)

const (
	abortIndex int8 = math.MaxInt8 / 2
)

// Context 是單次請求中的上下文建構體。
type Context struct {
	// Keys 存放透過 `Set` 儲存的鍵值組，僅限於單次請求。
	Keys map[string]interface{}
	// Errors 存放開發者自訂的錯誤，可用在中介軟體或處理函式中。
	Errors errorMsgs
	// Session 是產生此連線的客戶端階段建構體。
	Session *Session
	// isAborted 表示了這個請求是不是以已經被終止了。
	isAborted bool
	// ID 是本次請求的工作編號，用以讓客戶端呼叫相對應的處理函式。
	ID int
	// Request 是這個 WebSocket 的 HTTP 請求建構體。
	Request *http.Request
	// Method 是目前正執行的方法建構體。
	Method *Method

	// data 呈現了本次接收到的資料，主要是 MsgPack 內容格式。
	data []byte
	// index 是目前所執行的處理函式索引。
	index int8
	// handlers 存放此請求將會呼叫的所有處理函式。
	handlers []HandlerFunc
	// params 存放著已解序的參數陣列。
	params []interface{}
	// files 為檔案欄位切片，用以存放使用者上傳後且已解析的檔案。
	files map[string][]*File
	// engine 是主要引擎。
	engine *Engine
}

// Error 能夠將發生的錯誤保存到單次 Session 中。
func (c *Context) Error(err error) *Context {
	c.Errors = append(c.Errors, &Error{
		Err:  err,
		Type: ErrorTypePrivate,
	})
	return c
}

// ClientIP 會盡所能地取得到最真實的客戶端 IP 位置。
func (c *Context) ClientIP() string {
	clientIP := c.requestHeader("X-Forwarded-For")
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = strings.TrimSpace(c.requestHeader("X-Real-Ip"))
	if len(clientIP) > 0 {
		return clientIP
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// Param 能夠從參數陣列中透過指定索引取得特定的參數。
func (c *Context) Param(indexes ...int) *Param {
	// 預設索引為 `0` 省去使用者指定的困擾。
	var index int

	// 如果已解序的參數陣列長度為空，就表示先前可能還沒解序資料成參數陣列。
	// 因此我們需要透過 MsgPack 解序資料並存成參數陣列。
	if len(c.params) == 0 {
		// 如果錯誤發生回傳一個空的參數建構體。
		if err := msgpack.Unmarshal(c.data, &c.params); err != nil {
			return &Param{
				data: nil,
				len:  len(c.params),
			}
		}
	}

	// 如果使用者有指定索引則更改。
	if len(indexes) != 0 {
		index = indexes[0]
	}

	// 如果指定的索引位置大於參數陣列長度，則表示無此參數。
	if index+1 > len(c.params) {
		return &Param{
			data: nil,
			len:  len(c.params),
		}
	}

	// 不然就回傳指定的參數資料。
	return &Param{
		data: c.params[index],
		len:  len(c.params),
	}
}

// Copy 會複製一份 `Context` 供你在 Goroutine 中操作不會遇上資料競爭與衝突問題。
func (c *Context) Copy() *Context {
	ctx := *c
	return &ctx
}

// IsAborted 會回傳一個此請求是否已被終止的布林值。
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort 終止此次請求，避免繼續執行。
func (c *Context) Abort() {
	c.index = abortIndex
}

// AbortWithRespond 終止此次請求，避免繼續執行。並以資料回應特定的客戶端。
func (c *Context) AbortWithRespond(result interface{}) {
	c.Abort()
	c.Respond(result)
}

// AbortWithError 結束此次請求，避免繼續執行。並以指定的狀態碼、錯誤資料與訊息回應特定的客戶端並表示錯誤發生。
func (c *Context) AbortWithError(code int, data interface{}, err error) {
	c.Abort()
	c.RespondWithError(code, data, err)
}

// Respond 會以指定的狀態碼、資料回應特定的客戶端。
func (c *Context) Respond(result interface{}) {
	c.Session.write(Response{
		Result: result,
		ID:     c.ID,
	})
}

// RespondWithError 會以指定的狀態碼、錯誤資料與訊息回應特定的客戶端並表示錯誤發生。
func (c *Context) RespondWithError(code int, data interface{}, err error) {
	c.Session.write(Response{
		Error: ResponseError{
			Code:    code,
			Data:    data,
			Message: err.Error(),
		},
		ID: c.ID,
	})
}

// Bind 能夠將接收到的資料映射到本地建構體。
func (c *Context) Bind(dest interface{}) error {
	return msgpack.Unmarshal(c.data, dest)
}

// MustBind 和 `Bind` 相同，差異在於若映射失敗會呼叫 `panic` 並阻止此次請求繼續執行。
func (c *Context) MustBind(dest interface{}) error {
	if err := msgpack.Unmarshal(c.data, dest); err != nil {
		panic(err)
	}
	return nil
}

// Set 會在本次上下文建構體中存放指定的鍵值組內容，可供其他處理函式存取。
func (c *Context) Set(key string, value interface{}) {
	c.Keys[key] = value
}

// MustGet 和 `Get` 相同，但沒有該鍵值組時會呼叫 `panic`。
func (c *Context) MustGet(key string) interface{} {
	v, ok := c.Get(key)
	if !ok {
		panic(ErrKeyNotFound)
	}
	return v
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

// GetFile 會回傳一個指定檔案欄位中的檔案。
func (c *Context) GetFile(name string) (*File, error) {
	f, err := c.GetFiles(name)
	if err != nil {
		return nil, err
	}
	return f[0], nil
}

// MustGetFile 會回傳一個指定檔案欄位中的檔案，並且在沒有該檔案的時候呼叫 `panic` 終止此請求。
func (c *Context) MustGetFile(name string) *File {
	f, err := c.GetFile(name)
	if err != nil {
		panic(err)
	}
	return f
}

// GetFiles 會取得指定檔案欄位中的全部檔案，並回傳一個檔案切片供開發者遍歷。
func (c *Context) GetFiles(name string) ([]*File, error) {
	f, ok := c.files[name]
	if !ok {
		return nil, ErrFileNotFound
	}
	if len(f) == 0 {
		return nil, ErrFileNotFound
	}
	return f, nil
}

// MustGetFiles 會取得指定檔案欄位中的全部檔案，並回傳一個檔案切片供開發者遍歷。沒有該檔案欄位的時候會呼叫 `panic` 終止此請求。
func (c *Context) MustGetFiles(name string) []*File {
	f, err := c.GetFiles(name)
	if err != nil {
		panic(err)
	}
	return f
}

// Next 來呼叫下一個方法處理函式，常用於中介軟體中。
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Emit 會向此客戶端廣播一個事件。
func (c *Context) Emit(event string, result interface{}) {
	c.Session.write(Response{
		Event:  event,
		Result: result,
	})
}

// EmitOthers 會將帶有指定資料的事件向自己以外的所有人廣播。
func (c *Context) EmitOthers(event string, result interface{}) {
	c.writeOthers(Response{
		Event:  event,
		Result: result,
	})
}

// requestHeader 會回傳指定的 HTTP 標頭內容。
func (c *Context) requestHeader(key string) string {
	if values, _ := c.Request.Header[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}

// writeOthers 會將傳入的回應轉譯成 MessagePack 並且寫入除了自己以外的其他客戶端 WebSocket。
func (c *Context) writeOthers(resp Response) {
	if msg, err := msgpack.Marshal(resp); err == nil {
		c.engine.server.websocket.BroadcastOthers(msg, c.Session.websocket)
	}
}
