package mego

import (
	"io"
	"net/http"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/vmihailenco/msgpack"

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
	// StatusFileTooLarge 表示檔案過大無法上傳。
	StatusFileTooLarge = 72
)

var (
	// DefaultPort 是 Mego 引擎的預設埠口。
	DefaultPort = ":5000"
	// DefaultWriter 是預設的紀錄、資料寫入者。
	DefaultWriter io.Writer = os.Stdout
	// DefaultErrorWriter 是預設的錯誤資料寫入者。
	DefaultErrorWriter io.Writer = os.Stderr
)

// H 是常用的資料格式，簡單說就是 `map[string]interface{}` 的別名。
type H map[string]interface{}

// HandlerFunc 是方法處理函式的型態別名。
type HandlerFunc func(*Context)

// New 會建立一個新的 Mego 空白引擎。
func New() *Engine {
	return &Engine{
		Sessions: make(map[string]*Session),
		Events:   make(map[string]*Event),
		Methods:  make(map[string]*Method),
		subscribeHandler: func(e, ch string, c *Context) bool {
			return true
		},
	}
}

// Default 會建立一個帶有 `Recovery` 和 `Logger` 中介軟體的 Mego 引擎。
func Default() *Engine {
	e := &Engine{
		Sessions: make(map[string]*Session),
		Events:   make(map[string]*Event),
		Methods:  make(map[string]*Method),
		subscribeHandler: func(e, ch string, c *Context) bool {
			return true
		},
	}
	return e.Use(Recovery(), Logger())
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
	// Methods 是所有可用的方法切片。
	Methods map[string]*Method
	// Option 是這個引擎的設置。
	Option *EngineOption

	// handlers 是保存將會執行的全域中介軟體切片。
	handlers []HandlerFunc
	// noMethod 是當呼叫不存在方式時所會呼叫的處理函式。
	noMethod []HandlerFunc
	// server 存放了最主要的底層伺服器與其 WebSocket 引擎。
	server *server
	// subscribeHandler 是處理所有事件訂閱的函式。
	subscribeHandler func(event string, channel string, c *Context) bool
}

// EngineOption 是引擎的選項設置。
type EngineOption struct {
	// MaxSize 是這個方法允許接收的最大位元組（Bytes）。
	MaxSize int
	// MaxChunkSize 是這個方法允許的區塊最大位元組（Bytes）。
	MaxChunkSize int
	// MaxFileSize 是這個方法允許的檔案最大位元組（Bytes），會在每次接收區塊時結算總計大小，
	// 如果超過此大小則停止接收檔案。
	MaxFileSize int
	// MaxSessions 是引擎能容忍的最大階段連線數量。
	MaxSessions int
	// CheckInterval 是每隔幾秒進行一次階段是否仍存在的連線檢查，
	// 此為輕量檢查而非發送回應至客戶端。
	CheckInterval int
}

// Method 呈現了一個方法。
type Method struct {
	// Name 是此方法的名稱。
	Name string
	// Handlers 是中介軟體與處理函式。
	Handlers []HandlerFunc
	// Option 是此方法的選項。
	Option *MethodOption
}

// MethodOption 是一個方法的選項。
type MethodOption struct {
	// MaxSize 是這個方法允許接收的最大位元組（Bytes）。此選項會覆蓋引擎設定。
	MaxSize int
	// MaxChunkSize 是這個方法允許的區塊最大位元組（Bytes）。此選項會覆蓋引擎設定。
	MaxChunkSize int
	// MaxFileSize 是這個方法允許的檔案最大位元組（Bytes），會在每次接收區塊時結算總計大小，
	// 如果超過此大小則停止接收檔案。此選項會覆蓋引擎設定。
	MaxFileSize int
}

// Run 會在指定的埠口執行 Mego 引擎。
func (e *Engine) Run(port ...string) {
	// 初始化一個 Melody 套件框架並當作 WebSocket 底層用途。
	m := melody.New()
	// 以 WebSocket 初始化一個底層伺服器並存入 Mego。
	e.server = &server{
		websocket: m,
	}

	// 設定預設埠口。
	p := DefaultPort
	if len(port) > 0 {
		p = port[0]
	}

	// 將接收到的所有訊息轉交給訊息處理函式。
	m.HandleMessage(e.messageHandler)
	// 將所有連線請求轉交給連線處理函式。
	m.HandleConnect(e.connectHandler)
	// 將所有斷線的請求轉交給斷線處理函式。
	m.HandleDisconnect(e.disconnectHandler)

	// 開始在指定埠口監聽 HTTP 請求並交由底層伺服器處理。
	http.ListenAndServe(p, e.server)
}

// disconnectHandler 會處理斷開連線的 WebSocket。
func (e *Engine) disconnectHandler(s *melody.Session) {
	// 如果客戶端離線了就自動移除他所監聽的事件和所有 Sessions
}

// connectHandler 處理連接起始的函式。
func (e *Engine) connectHandler(s *melody.Session) {
	// 替此階段建立一個獨立的 UUID。
	id := uuid.NewV4().String()
	// 在底層階段存放此階段的編號。
	s.Set("ID", id)
	// 將 Mego 階段放入引擎中保存。
	e.Sessions[id] = &Session{
		ID:        id,
		websocket: s,
	}
}

// messageHandler 處理所有接收到的訊息，並轉接給相對應的方法處理函式。
func (e *Engine) messageHandler(s *melody.Session, msg []byte) {
	var req Request

	// 將接收到的訊息映射到本地端的請求建構體。
	err := msgpack.Unmarshal(msg, &req)
	if err != nil {
		return
	}

	// 取得這個 WebSocket 階段對應的 Mego 階段。
	id, ok := s.Get("ID")
	if !ok {
		return
	}
	// 透過獨有編號在引擎中找出相對應的階段資料。
	sess, ok := e.Sessions[id.(string)]
	if !ok {
		return
	}

	// 依接收到的方法處理指定的事情。
	switch methodName := strings.ToUpper(req.Method); methodName {
	// 呼叫 Mego 初始化方法。
	case "MEGOINITIALIZE":
		var keys map[string]interface{}

		// 將接收到的資料映射到本地的 map 型態，並保存到階段資料中的鍵值組。
		if err := msgpack.Unmarshal(req.Params, &keys); err != nil {
			return
		}
		sess.Keys = keys

	// 呼叫 Mego 訂閱方法。
	case "MEGOSUBSCRIBE":
		// 取得事件訂閱資料，此為陣列。索引 0 為事件名稱、索引 1 為頻道名稱。
		if len(req.Event) != 2 {
			return
		}
		evt, ch := req.Event[0], req.Event[1]

		// 建立一個上下文建構體。
		ctx := &Context{
			Session: sess,
			ID:      req.ID,
			Request: s.Request,
			data:    req.Params,
		}

		// 呼叫訂閱處理函式，如果回傳的是 `true` 才繼續。
		if ok := e.subscribeHandler(evt, ch, ctx); !ok {
			return
		}

		// 如果欲訂閱的事件不存在，就建立一個。
		if _, ok = e.Events[evt]; !ok {
			e.Events[evt] = &Event{
				Name:     evt,
				Channels: make(map[string]*Channel),
				engine:   e,
			}
		}

		// 如果欲訂閱的頻道不存在，就建立一個。
		if _, ok = e.Events[evt].Channels[ch]; !ok {
			e.Events[evt].Channels[ch] = &Channel{
				Name:  ch,
				Event: e.Events[evt],
			}
		}

		// 如果欲訂閱的階段不在其頻道內，就將該階段存至該頻道作為訂閱者。
		var has bool
		for _, v := range e.Events[evt].Channels[ch].Sessions {
			if v == sess {
				has = true
			}
		}
		if !has {
			e.Events[evt].Channels[ch].Sessions = append(e.Events[evt].Channels[ch].Sessions, sess)
		}

	// 呼叫伺服端現有的方法。
	default:
		// 檢查此方法是否存在於伺服器中。
		method, ok := e.Methods[methodName]
		if !ok {
			return
		}

		// 建立一個上下文建構體。
		ctx := &Context{
			Session:  sess,
			Method:   method,
			ID:       req.ID,
			Request:  s.Request,
			data:     req.Params,
			handlers: e.handlers,
		}
		// 將該方法的處理函式推入上下文建構體中供依序執行。
		ctx.handlers = append(ctx.handlers, method.Handlers...)

		// 如果處理函式數量大於零的話就可以開始執行了。
		if len(ctx.handlers) > 0 {
			ctx.handlers[0](ctx)
		}
	}
}

// HandleRequest 會更改
// func (e *Engine) HandleRequest() *Engine {
// 	return e
// }

// func (e *Engine) HandleConnect() *Engine {
// 	return e
// }

// HandleSubscribe 會更改預設的事件訂閱檢查函式，開發者可傳入一個回呼函式並接收客戶端欲訂閱的事件與頻道和相關資料。
// 回傳一個 `false` 即表示客戶端的資格不符，將不納入訂閱清單中。該客戶端將無法接收到指定的事件。
func (e *Engine) HandleSubscribe(handler func(event string, channel string, c *Context) bool) *Engine {
	e.subscribeHandler = handler
	return e
}

// Use 會使用傳入的中介軟體作為全域使用。
func (e *Engine) Use(handlers ...HandlerFunc) *Engine {
	e.handlers = append(e.handlers, handlers...)
	return e
}

// Len 會回傳目前有多少個連線數。
func (e *Engine) Len() int {
	return len(e.Sessions)
}

// Close 會結束此引擎的服務。
func (e *Engine) Close() error {
	e.server.websocket.Close()
	return nil
}

// NoMethod 會在客戶端呼叫不存在方法時被執行。
func (e *Engine) NoMethod(handler ...HandlerFunc) *Engine {
	e.noMethod = handler
	return e
}

// Event 會建立一個新的事件，如此一來客戶端方能監聽。
func (e *Engine) Event(name string) {
	e.Events[name] = &Event{
		Name: name,
	}
}

// Register 會註冊一個指定的方法，並且允許客戶端呼叫此方法觸發指定韓式。
func (e *Engine) Register(method string, handler ...HandlerFunc) *Method {
	m := &Method{
		Name:     method,
		Handlers: handler,
	}
	e.Methods[method] = m
	return m
}

// Emit 會帶有指定資料並廣播指定事件與頻道，當頻道為空字串時則廣播到所有頻道。
func (e *Engine) Emit(event string, channel string, result interface{}) error {
	evt, ok := e.Events[event]
	if !ok {
		return ErrEventNotFound
	}
	ch, ok := evt.Channels[channel]
	if !ok {
		return ErrChannelNotFound
	}
	var firstErr error
	for _, v := range ch.Sessions {
		v.write(Response{
			Event:  event,
			Result: result,
		})
		//err := v.websocket.WriteBinary()
		//if firstErr == nil {
		//	firstErr = err
		//}
	}

	return firstErr
}

// EmitMultiple 會將指定事件與資料向指定的客戶端切片進行廣播。
func (e *Engine) EmitMultiple(event string, channel string, result interface{}, sessions []*Session) error {
	return nil
}

// EmitFilter 會以過濾函式來決定要將帶有指定資料的事件廣播給誰。
// 如果過濾函式回傳 `true` 則表示該客戶端會接收到該事件。
func (e *Engine) EmitFilter(event string, channel string, payload interface{}, filter func(*Session) bool) error {
	return nil
}
