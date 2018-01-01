package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/vmihailenco/msgpack"
)

var (
	// ErrChunkWithFiles 表示不能上傳區塊檔案同時夾帶其他檔案實體。
	ErrChunkWithFiles = errors.New("mego: cannot send the file chunks with the file entities")
	// ErrTimeout 表示一個請求過久都沒有回應。
	ErrTimeout = errors.New("mego: request timeout")
	// ErrClosed 表示正在使用一個已經關閉的連線。
	ErrClosed = errors.New("mego: use of closed connection")
)

func NewClient(url string) *Client {
	return &Client{
		URL:  url,
		UUID: uuid.NewV4().String(),
	}
}

type ClientOption struct {
	// ChunkSize 是預設的區塊分切基準位元組大小。
	ChunkSize int
}

type Client struct {
	// URL 是遠端 Mego 伺服器的網址。
	URL string
	// UUID 是此客戶端的獨立編號。
	UUID string
	// Option 是這個客戶端的設置選項。
	Option *ClientOption

	// reponseLocks 是用來保存回應阻塞頻道的切片。
	reponseLocks map[int]*chan bool
	//
	fileID int
	//
	fileNameID int
	// taskID 是遞加的請求編號。
	taskID int
	// keys 是保存於遠端的鍵值組。
	keys map[string]interface{}
	// conn 是底層的 WebSocket 連線。
	conn *websocket.Conn
}

// Call 能夠建立一個呼叫遠端指定方法的空白請求。
func (c *Client) Call(method string) *Request {
	return &Request{
		Method: method,
	}
}

// Set 能夠將指定資料保存於遠端伺服器，避免每次請求都需傳遞相同資料供伺服端讀取。
// 由於資料保存在指定伺服器上，若使用負載平衡可能導致另一個伺服器找不到相關資料，因此使用負載平衡時請不要用上此方式。
func (c *Client) Set(data map[string]interface{}) *Client {
	for k, v := range data {
		c.keys[k] = v
	}
	return c
}

// Connect 會開始連線到遠端伺服器。
func (c *Client) Connect() error {
	// 開啟一個 WebSocket 連線。
	conn, _, err := websocket.DefaultDialer.Dial(c.URL, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	// 將連線保存至客戶端建構體內。
	c.conn = conn

	// 不斷呼叫訊息處理函式。
	go func() {
		for {
			c.messageHandler()
		}
	}()

	// 初始化這個連線，並傳遞鍵值組資料。
	c.initializeConn()

	return nil
}

// initializeConn 會初始化一個連線，並傳送初始鍵值組至遠端伺服器保存。
func (c *Client) initializeConn() error {
	// 將鍵值組以 Message Pack 封裝，方能於伺服端解開。
	keys, err := msgpack.Marshal(c.keys)
	if err != nil {
		return err
	}
	// 傳送鍵值組至伺服端，啟動連線後的第一個訊息會被作為初始訊息。
	return c.writeMessage(Request{
		Params: keys,
	})
}

// writeMessage 會以 Message Pack 包裝訊息並傳遞至遠端伺服器。
func (c *Client) writeMessage(data interface{}) error {
	msg, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

//
func (c *Client) messageHandler() {
	_, m, err := c.conn.ReadMessage()
	if err != nil {

	}
}

// Reconnect 會重新連線，能在斷線或結束連線時使用。
func (c *Client) Reconnect() error {

}

// Close 會結束並關閉連線。
func (c *Client) Close() error {
	return c.conn.Close()
}

// Subscribe 可以訂閱指定的遠端事件，並在之後能透過 `On` 接收。
func (c *Client) Subscribe(event string, channel string) error {
	return c.writeMessage(Request{
		Method: "MegoSubscribe",
		Event:  []string{event, channel},
	})
}

// Unsubscribe 會取消訂閱指定的遠端事件，避免接收到無謂的事件。
func (c *Client) Unsubscribe(event string, channel string) error {

}

// On 能夠監聽系統或透過 `Subscribe` 所訂閱的事件，並做出相對應的動作。
func (c *Client) On(event string, handler func(*Event)) *Client {

}

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
func (e *Error) Error() string {
	return e.Message
}

// Bind 能夠將錯誤的詳細資料映射到本地建構體。
func (e *Error) Bind(dest interface{}) error {
	return msgpack.Unmarshal(e.Data, dest)
}

// Event 呈現了一個接收到的事件。
type Event struct {
	// Data 是事件所夾帶的資料。基於 Message Pack 格式需要透過 `Bind` 映射到本地建構體。
	Data []byte
}

// Bind 能將事件所夾帶的資料映射到本地的建構體。
func (e *Event) Bind(dest interface{}) error {
	return nil
}

// File 呈現了一個欲上傳的檔案資料。
type File struct {
	// Binary 是檔案的二進制。
	Binary []byte `codec:"b" msgpack:"b"`
	// ID 是檔案編號，用於區塊組合。
	ID int `codec:"i" msgpack:"i"`
	// Parts 呈現區塊的分塊進度。索引 0 表示總共區塊數，索引 1 則是本區塊編號。
	// 如果這個切片是空的表示此為實體檔案而非區塊。
	Parts []int `codec:"p" msgpack:"p"`
	// Name 是檔案的原始名稱。
	Name string `codec:"n" msgpack:"n"`

	// source 是這個檔案的源頭，也許是 `string`、`[]byte`、`*os.File`
	source interface{}
}

// Request 呈現了一個籲發送至遠端伺服器的請求。
type Request struct {
	// Method 是欲呼叫的方法名稱。
	Method string `codec:"m" msgpack:"m"`
	// Params 是資料或參數。
	Params []byte `codec:"p" msgpack:"p"`
	// Files 是此請求所包含的檔案欄位與其內容。
	Files map[string][]*File `codec:"f" msgpack:"f"`
	// Event 是欲註冊的事件名稱（索引 0）與頻道（索引 1）。
	Event []string `codec:"e" msgpack:"e"`
	// ID 為本次請求編號，若無則為單次通知廣播不需回應。
	ID int `codec:"i" msgpack:"i"`

	// client 是建立這個請求的客戶端。
	client *Client
}

// Send 會保存稍後將發送的資料。
func (r *Request) Send(data interface{}) *Request {
	params, err := msgpack.Marshal(data)
	if err != nil {

	}
	r.Params = params
	return r
}

func (r *Request) storeFile(file interface{}, isChunk bool, fieldName ...string) {
	var f *File
	var total int

	// 遞增檔案編號。
	r.client.fileID++

	// 依照傳入的檔案型態做出不同的讀取方式。
	switch v := file.(type) {
	// *os.File 表示檔案寫入者。
	case *os.File:
		bin, err := ioutil.ReadAll(v)
		if err != nil {

		}
		// 如果本檔案要以區塊上傳，就計算區塊切片數。
		if isChunk {
			total
		}
		f = &File{
			Binary: bin,
			ID:     r.client.fileID,
			Name:   filepath.Base(v.Name()),
		}

	// 位元組表示檔案二進制內容。
	case []byte:
		f = &File{
			Binary: v,
			ID:     r.client.fileID,
		}

	// 字串型態表示檔案路徑。
	case string:
		bin, err := ioutil.ReadFile(v)
		if err != nil {

		}
		f = &File{
			Binary: bin,
			ID:     r.client.fileID,
			Name:   filepath.Base(v),
		}
	}

	// 取得檔案欄位名稱，若無指定則自動編號取名。
	var n string
	if len(fieldName) > 0 {
		n = fieldName[0]
	} else {
		// 遞增檔案名稱編號。
		r.client.fileNameID++
		n = fmt.Sprintf("File%d", r.client.fileNameID)
	}

	// 取得指定檔案欄位並檢查是否存在，不存在則初始化。
	field, ok := r.Files[n]
	if !ok {
		r.Files[n] = []*File{}
	}
	// 將本檔案推入指定檔案欄位中。
	r.Files[n] = append(r.Files[n], f)
}

func (r *Request) readFile() {

}

// SendFile 會保存稍後將上傳的檔案。
func (r *Request) SendFile(file interface{}, fieldName ...string) *Request {
	r.storeFile(file, false, fieldName...)
	return r
}

// SendFiles 能夠保存多個檔案並將其歸納為同個檔案欄位。
func (r *Request) SendFiles(files []interface{}, fieldName ...string) *Request {
	for _, v := range files {
		r.SendFile(v, fieldName...)
	}
	return r
}

// SendFileChunks 會保存稍後將以區塊方式上傳的檔案。
// 注意：請求使用區塊檔案上傳時，不可使用 `SendFile` 夾帶其他檔案。
func (r *Request) SendFileChunks(file interface{}, fieldName ...string) *Request {
	r.storeFile(file, true, fieldName...)
	return r
}

// End 結束並發送這個請求且不求回應。
func (r *Request) End() error {
	return nil
}

// EndStruct 結束並發送這個請求，且將回應映射到本地建構體上。
func (r *Request) EndStruct(dest interface{}) error {
	return nil
}
