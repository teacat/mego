package mego

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	mirror "github.com/TeaMeow/Mirror"
	uuid "github.com/satori/go.uuid"
	"github.com/vmihailenco/msgpack"

	"github.com/olahol/melody"
)

// SubscribeHandler 是訂閱事件處理函式。
type SubscribeHandler func(evt string, ch string, ctx *Context) bool

// ChunkHandler 是區塊處理函式。
type ChunkHandler func(ctx *Context, raw *RawFile, dest *File) ChunkStatus

// H 是常用的資料格式，簡單說就是 `map[string]interface{}` 的別名。
type H map[string]interface{}

// HandlerFunc 是方法處理函式的型態別名。
type HandlerFunc func(*Context)

// ChunkStatus 是區塊處理結果狀態碼，用以讓 Mego 得知下一步該如何執行。
type ChunkStatus int

const (
	// ChunkNext 表示本次處理成功，請求下一個檔案區塊。
	ChunkNext ChunkStatus = iota
	// ChunkDone 表示所有區塊皆處理完畢，結束檔案處理。
	ChunkDone
	// ChunkAbort 表示不打算處理本檔案了，結束此檔案的處理手續並停止上傳。
	ChunkAbort
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

var (
	// DefaultPort 是 Mego 引擎的預設埠口。
	DefaultPort = ":5000"
	// DefaultWriter 是預設的紀錄、資料寫入者。
	DefaultWriter io.Writer = os.Stdout
	// DefaultErrorWriter 是預設的錯誤資料寫入者。
	DefaultErrorWriter io.Writer = os.Stderr
)

// New 會建立一個新的 Mego 空白引擎。
func New() *Engine {
	return &Engine{
		Sessions:     make(map[string]*Session),
		Events:       make(map[string]*Event),
		Methods:      make(map[string]*Method),
		chunkHandler: chunkHandler,
	}
}

// Default 會建立一個帶有 `Recovery` 和 `Logger` 中介軟體的 Mego 引擎。
func Default() *Engine {
	e := New()
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
	subscribeHandler SubscribeHandler
	// chunkHandler 是預設的方法區塊處理回呼函式，能被各個方法獨立覆蓋。
	chunkHandler ChunkHandler
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
	// ChunkHandler 是本方法的區塊處理回呼函式。
	ChunkHandler ChunkHandler
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
	// 將所有斷線的請求轉交給斷線處理函式。
	m.HandleDisconnect(e.disconnectHandler)
	fmt.Println("Running...")
	// 開始在指定埠口監聽 HTTP 請求並交由底層伺服器處理。
	http.ListenAndServe(p, e.server)
}

// disconnectHandler 會處理斷開連線的 WebSocket。
func (e *Engine) disconnectHandler(s *melody.Session) {
	// 如果客戶端離線了就自動移除他所監聽的事件和所有 Sessions
}

// subscribe 會替傳入的 Session 訂閱指定的事件與頻道。
func (e *Engine) subscribe(sess *Session, evtName string, chName string) {
	// 如果欲訂閱的事件不存在，就建立一個。
	evt, ok := e.Events[evtName]
	if !ok {
		e.Events[evtName] = &Event{
			Name:     evtName,
			Channels: make(map[string]*Channel),
			engine:   e,
		}
	}
	// 如果欲訂閱的頻道不存在，就建立一個。
	ch, ok := evt.Channels[chName]
	if !ok {
		evt.Channels[chName] = &Channel{
			Name:  chName,
			Event: evt,
		}
	}
	// 如果欲訂閱的階段不在其頻道內，就將該階段存至該頻道作為訂閱者。
	var has bool
	for _, v := range ch.Sessions {
		if v == sess {
			has = true
		}
	}
	if !has {
		ch.Sessions = append(ch.Sessions, sess)
	}
}

// unsubscribe 會替傳入的 Session 取消訂閱指定的事件與頻道。
func (e *Engine) unsubscribe(sess *Session, evtName string, chName string) {
	if ch, ok := e.Events[evtName].Channels[chName]; ok {
		ch.Kick(sess.ID)
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
	// 如果沒有的話則當此請求為初次設置。
	id, ok := s.Get("MegoID")
	if !ok {
		var keys map[string]interface{}
		// 將接收到的資料映射到本地的 map 型態，並保存到階段資料中的鍵值組。
		mirror.Cast(req.Params, &keys)

		// 如果鍵值組中沒有 Mego 客戶端建立的獨立編號就離開。
		i, ok := keys["MegoID"]
		if !ok {
			return
		}
		// 如果獨立編號不是字串也離開。
		id, ok := i.(string)
		if !ok {
			return
		}
		// 如果 UUID 不正確也離開。
		_, err := uuid.FromString(id)
		if err != nil {
			return
		}

		// 在底層階段存放此階段的編號。
		s.Set("MegoID", id)

		// 將 Mego 階段放入引擎中保存。
		e.Sessions[id] = &Session{
			ID:        id,
			engine:    e,
			websocket: s,
		}
	}

	// 重新取得一次此客戶端的獨立 UUID 編號。
	id, _ = s.Get("MegoID")
	if id == nil {
		return
	}

	// 透過獨有編號在引擎中找出相對應的階段資料。
	sess, ok := e.Sessions[id.(string)]
	if !ok {
		return
	}

	// 依接收到的方法處理指定的事情。
	switch methodName := strings.ToUpper(req.Method); methodName {
	// 呼叫 Mego 取消訂閱方法。
	case "MEGOUNSUBSCRIBE":
		// 建立一個上下文建構體。
		ctx := &Context{
			Session: sess,
			ID:      req.ID,
			Request: s.Request,
			data:    req.Params,
		}
		// 取得事件訂閱資料，此為陣列。索引 0 為事件名稱、索引 1 為頻道名稱。
		evt := ctx.Param(0).GetString()
		ch := ctx.Param(1).GetString()

		// 執行此客戶端的取消訂閱方法。
		e.unsubscribe(sess, evt, ch)

	// 呼叫 Mego 訂閱方法。
	case "MEGOSUBSCRIBE":
		// 建立一個上下文建構體。
		ctx := &Context{
			Session: sess,
			ID:      req.ID,
			Request: s.Request,
			data:    req.Params,
		}
		// 取得事件訂閱資料，此為陣列。索引 0 為事件名稱、索引 1 為頻道名稱。
		evt := ctx.Param(0).GetString()
		ch := ctx.Param(1).GetString()

		// 呼叫訂閱處理函式，如果回傳的是 `true` 才繼續。
		if e.subscribeHandler != nil {
			if ok := e.subscribeHandler(evt, ch, ctx); !ok {
				return
			}
		}
		// 執行此客戶端的訂閱方法。
		e.subscribe(sess, evt, ch)

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
			files:    make(map[string][]*File),
			handlers: e.handlers,
		}
		// 將該方法的處理函式推入上下文建構體中供依序執行。
		ctx.handlers = append(ctx.handlers, method.Handlers...)

		// 解析上傳的檔案。
		if done := e.fileHandler(ctx, req.Files); !done {
			// 如果是區塊檔案且尚未處理完畢，就先不要繼續執行。
			// 告訴客戶端上傳下一個區塊。
			return
		}

		// 如果處理函式數量大於零的話就可以開始執行了。
		if len(ctx.handlers) > 0 {
			ctx.handlers[0](ctx)
		}
	}
}

// chunkHandler 是預設的區塊處理函式，這會接收區塊並組成一個檔案。
func chunkHandler(c *Context, raw *RawFile, dest *File) ChunkStatus {
	// 在系統中建立並開啟一個新的暫存檔案。
	t, err := os.OpenFile(fmt.Sprintf("%s/MEGO_CHUNK_%s_%d", os.TempDir(), c.Session.ID, raw.ID), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {

	}
	// 將使用者上傳的位元組內容寫入暫存檔案中，並取得該位元組長度做為檔案大小。
	size, err := t.Write(raw.Binary)
	if err != nil {

	}
	// 關閉檔案。
	if err := t.Close(); err != nil {

	}

	// 取得總區塊數。
	total := raw.Parts[0]
	// 取得本區塊編號。
	current := raw.Parts[1]

	// 如果這是最後一個區塊就回傳處理完成狀態碼。
	if current == total {
		// 從檔案名稱中取得名稱與副檔名。
		ext := filepath.Ext(raw.Name)
		name := strings.TrimSuffix(raw.Name, ext)

		// 如果副檔名長度過短就是出問題了。
		if len(ext) < 2 {

		}

		// 將正確的檔案資料配置到檔案建構體中。
		dest.Name = name
		dest.Extension = ext[1:]
		dest.Path = t.Name()
		dest.Size = size

		//
		return ChunkDone
	}
	return ChunkNext
}

// fileHandler 會處理並解析接收到的檔案。如果接收到了區塊內容，則會呼叫額外的區塊處理函式。
// 在這種情況本函式會回傳 `false` 來終止請求的繼續，直到所有區塊都處理完為止。
func (e *Engine) fileHandler(c *Context, fields map[string][]*RawFile) bool {
	// 遍歷每個檔案欄位。
	for field, files := range fields {
		// 如果這個檔案欄位不存在於本地的上下文建構體中，
		// 那麼就初始化一個。
		if _, ok := c.files[field]; !ok {
			c.files[field] = []*File{}
		}

		// 遍歷這個檔案欄位中的所有檔案。
		for _, f := range files {
			// 如果這個檔案內容不是最後結果，即表示這是個區塊內容。
			if len(f.Parts) > 0 {
				// 初始化一個目標檔案，在區塊組合完畢後就使用這個檔案建構體。
				var dest *File
				var status ChunkStatus

				// 呼叫區塊處理函式。
				switch {
				// 如果此方法有自訂的區塊處理函式則優先呼叫。
				case c.Method.ChunkHandler != nil:
					status = c.Method.ChunkHandler(c, f, dest)
				// 沒有則就呼叫全域區塊處理函式。
				case e.chunkHandler != nil:
					status = e.chunkHandler(c, f, dest)
				}

				// 依照區塊處理的狀態碼決定下一步。
				switch status {
				// ChunkNext 表示本次處理成功，請求下一個檔案區塊。
				case ChunkNext:
					c.Session.write(Response{
						Event: "MegoChunkNext",
						ID:    c.ID,
					})
					// 結束本次請求，避免區塊還沒處理完畢就繼續呼叫了接下來的方法函式。
					return false
				// ChunkAbort 表示不打算處理本檔案了，結束此檔案的處理手續並停止上傳。
				case ChunkAbort:
					c.Session.write(Response{
						Event: "MegoChunkAbort",
						ID:    c.ID,
					})
					// 結束本次請求，避免區塊還沒處理完畢就繼續呼叫了接下來的方法函式。
					return false
				// ChunkDone 表示所有區塊皆處理完畢，結束檔案處理。
				case ChunkDone:
					// 將這個檔案整理後推入至上下文建構體中的檔案欄位。
					c.files[field] = append(c.files[field], dest)
					// 繼續本次請求，並呼叫接下來的方法函式。
					return true
				}
			}

			// 在系統中建立並開啟一個新的暫存檔案。
			t, err := ioutil.TempFile("", "")
			if err != nil {

			}

			// 將使用者上傳的位元組內容寫入暫存檔案中，並取得該位元組長度做為檔案大小。
			size, err := t.Write(f.Binary)
			if err != nil {

			}

			// 從檔案名稱中取得名稱與副檔名。
			ext := filepath.Ext(f.Name)
			name := strings.TrimSuffix(f.Name, ext)
			fmt.Printf("%+v", f)
			// 如果副檔名長度過短就是出問題了。
			if len(ext) < 2 {
				panic("Fuck")
			}

			// 將這個檔案整理後推入至上下文建構體中的檔案欄位。
			c.files[field] = append(c.files[field], &File{
				Name:      name,
				Extension: ext[1:],
				Path:      t.Name(),
				Size:      size,
			})
		}
	}

	return true
}

// HandleSubscribe 會更改預設的事件訂閱檢查函式，開發者可傳入一個回呼函式並接收客戶端欲訂閱的事件與頻道和相關資料。
// 回傳一個 `false` 即表示客戶端的資格不符，將不納入訂閱清單中。該客戶端將無法接收到指定的事件。
func (e *Engine) HandleSubscribe(handler SubscribeHandler) *Engine {
	e.subscribeHandler = handler
	return e
}

// HandleChunk 會更改預設的區塊處理函式，開發者可以傳入一個回呼函式並接收區塊內容。
// 回傳 `ChunkStatus` 來告訴 Mego 區塊的處理狀態如何。
func (e *Engine) HandleChunk(handler ChunkHandler) *Engine {
	e.chunkHandler = handler
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
	e.Methods[strings.ToUpper(method)] = m
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
