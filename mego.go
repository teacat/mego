package mego

import (
	"net/http"

	"github.com/olahol/melody"
)

const (
	// 狀態碼範圍如下：
	// 0 ~ 50 正常、51 ~ 100 錯誤、101 ~ 999 自訂狀態碼。

	// StatusOK 表示正常。
	StatusOK = 0
	// StatusProcessing 表示請求已被接受並且正在處理中，並不會立即完成。
	StatusProcessing = 1
	// StatusNoChanges 表示這個請求沒有改變任何結果，例如：使用者刪除了一個早已被刪除的物件。
	StatusNoChanges = 2
	// StatusFileNext 表示此檔案區塊已處理完畢，需上傳下個區塊。
	StatusFileNext = 10
	// StatusFileAbort 表示終止整個檔案上傳進度。
	StatusFileAbort = 11

	// StatusError 表示有內部錯誤發生。
	StatusError = 51
	// StatusFull 表示此請求無法被接受，因為額度已滿。例如：使用者加入了一個已滿的聊天室、好友清單已滿。
	StatusFull = 52
	// StatusExists 表示請求的事物已經存在，例如：重複的使用者名稱、電子郵件地址。
	StatusExists = 53
	// StatusInvalid 表示此請求格式不正確。
	StatusInvalid = 54
	// StatusNotFound 表示找不到請求的資源。
	StatusNotFound = 55
	// StatusNotAuthorized 表示使用者需要登入才能進行此請求。
	StatusNotAuthorized = 56
	// StatusNoPermission 表示使用者已登入但沒有相關權限執行此請求。
	StatusNoPermission = 57
	// StatusUnimplemented 表示此功能尚未實作完成。
	StatusUnimplemented = 58
	// StatusTooManyRequests 表示使用者近期發送太多請求，需稍後再試。
	StatusTooManyRequests = 59
	// StatusResourceExhausted 表示使用者可用的額度已耗盡。
	StatusResourceExhausted = 60
	// StatusBusy 表示伺服器正繁忙無法進行執行請求。
	StatusBusy = 61
	// StatusFileRetry 表示檔案區塊發生錯誤，需要重新上傳相同區塊。
	StatusFileRetry = 70
	// StatusFileEmpty 表示上傳的檔案、區塊是空的。
	StatusFileEmpty = 71
	// StatusFileSize 表示檔案過大無法上傳。
	StatusFileSize = 72
)

// H 是常用的資料格式，簡單說就是 `map[string]interface{}` 的別名。
type H map[string]interface{}

// HandlerFunc 是方法處理函式的型態別名。
type HandlerFunc func(*Context)

// New 會建立一個新的 Mego 空白引擎。
func New() *Engine {
	return &Engine{}
}

// server 是基礎伺服器用來與基本 HTTP 連線進行互動。
type server struct {
	websocket *melody.Melody
}

// ServeHTTP 會將所有的 HTTP 請求轉嫁給 WebSocket。
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.websocket.HandleRequest(w, r)
}

// Engine 是 Mego 最主要的引擎結構體。
type Engine struct {
	// Sessions 儲存了正在連線的所有階段。
	Sessions map[string]*Session
	// Events 儲存了所有可用的事件與其監聽的客戶端資料。
	Events map[string]*Event
	// middlewares 是保存將會執行的全域中介軟體切片。
	middlewares []HandlerFunc
	// methods 是所有可用的方法切片。
	methods map[string][]HandlerFunc
	//
	noMethod []HandlerFunc
	//
	noEvent []HandlerFunc
}

// Run 會在指定的埠口執行 Mego 引擎。
func (e *Engine) Run(port string) {
	// 初始化一個 Melody 套件框架並當作 WebSocket 底層用途。
	m := melody.New()
	// 以 WebSocket 初始化一個底層伺服器。
	s := &server{
		websocket: m,
	}

	// 開始在指定埠口監聽 HTTP 請求並交由底層伺服器處理。
	http.ListenAndServe(port, s)
}

// Use 會使用傳入的中介軟體作為全域使用。
func (e *Engine) Use(middleware ...HandlerFunc) *Engine {
	e.middlewares = append(e.middlewares, middleware...)
	return e
}

// Len 會回傳目前有多少個連線數。
func (e *Engine) Len() int {
	return len(e.Sessions)
}

// Close 會結束此引擎的服務。
func (e *Engine) Close() error {
	return nil
}

// NoMethod 會在客戶端呼叫不存在方法時被執行。
func (e *Engine) NoMethod(handler ...HandlerFunc) *Engine {
	e.noMethod = handler
	return e
}

// NoEvent 會在客戶端監聽不存在事件時被執行。
func (e *Engine) NoEvent(handler ...HandlerFunc) *Engine {
	e.noEvent = handler
	return e
}

// Event 會建立一個新的事件，如此一來客戶端方能監聽。
func (e *Engine) Event(name string) {
	e.Events[name] = &Event{
		Name: name,
	}
}

// Register 會註冊一個指定的方法，並且允許客戶端呼叫此方法觸發指定韓式。
func (e *Engine) Register(method string, handler ...HandlerFunc) {
	e.methods[method] = handler
}

// Emit 會帶有指定資料並向所有人廣播指定事件。
func (e *Engine) Emit(event string, payload interface{}) error {
	return nil
}

// EmitMultiple 會將指定事件與資料向指定的客戶端切片進行廣播。
func (e *Engine) EmitMultiple(event string, payload interface{}, sessions []*Session) error {
	return nil
}

// EmitFilter 會以過濾函式來決定要將帶有指定資料的事件廣播給誰。
// 如果過濾函式回傳 `true` 則表示該客戶端會接收到該事件。
func (e *Engine) EmitFilter(event string, payload interface{}, filter func(*Session) bool) error {
	return nil
}

// Receive 會建立一個指定的方法，並且允許客戶端傳送檔案至此方法。
func (e *Engine) Receive(method string, handler ...HandlerFunc) {

}
