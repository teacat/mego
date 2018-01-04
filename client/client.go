package client

import (
	"time"

	"github.com/TeaMeow/Mego"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/vmihailenco/msgpack"
)

// H 是 `map[string]interface{}` 簡寫。
type H map[string]interface{}

var (
	// ChunkSize 是預設的區塊分切基準位元組大小。
	ChunkSize = mego.MB * 1
	// Timeout 是預設的逾期秒數。
	Timeout = time.Second * 15
	// UploadTimeout 是每個區塊、所有檔案的上傳逾期秒數，`0` 表示無上限。
	UploadTimeout = time.Second * 30
)

// New 能夠回傳一個新的客戶端。
func New(url string) *Client {
	return &Client{
		URL:  url,
		UUID: uuid.NewV4().String(),
		Option: &ClientOption{
			ChunkSize:     ChunkSize,
			Timeout:       Timeout,
			UploadTimeout: UploadTimeout,
		},
	}
}

// ClientOption 呈現了一個客戶端的設置。
type ClientOption struct {
	// ChunkSize 是預設的區塊分切基準位元組大小。
	ChunkSize int
	// Timeout 是預設的逾期秒數。
	Timeout time.Duration
	// UploadTimeout 是每個區塊、所有檔案的上傳逾期秒數，`0` 表示無上限。
	UploadTimeout time.Duration
}

// Client 是一個客戶端。
type Client struct {
	// URL 是遠端 Mego 伺服器的網址。
	URL string
	// UUID 是此客戶端的獨立編號。
	UUID string
	// Option 是這個客戶端的設置選項。
	Option *ClientOption

	// requests 是用來保存請求的儲藏區。
	requests map[int]Request
	// fileID 是自動遞增的檔案編號。
	fileID int
	// taskID 是遞加的請求編號。
	taskID int
	// listeners 是已註冊的事件監聽器列表，會在接收事件時呼叫相對應的函式。
	listeners map[string]func(*Event)
	// keys 是保存於遠端的鍵值組。
	keys map[string]interface{}
	// conn 是底層的 WebSocket 連線。
	conn *websocket.Conn
}

// Call 能夠建立一個呼叫遠端指定方法的空白請求。
func (c *Client) Call(method string) *Request {
	c.taskID++
	return &Request{
		Method: method,
		ID:     c.taskID,
		Option: &RequestOption{
			ChunkSize:     c.Option.ChunkSize,
			Timeout:       c.Option.Timeout,
			UploadTimeout: c.Option.UploadTimeout,
		},
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
	// 持續接收訊息。
	_, msg, err := c.conn.ReadMessage()
	if err != nil {

	}
	// 將接收到的訊息從 MessagePack 格式映射回本地的回應建構體。
	var resp *Response
	if err := msgpack.Unmarshal(msg, &resp); err != nil {

	}
	// 如果回應沒有編號，又有事件名稱則表示自訂事件。
	if resp.ID == 0 && resp.Event != "" {
		//
		return
	}
	// 如果回應有編號，取得並確定相對應的請求存在。
	req, ok := c.requests[resp.ID]
	if !ok {
		return
	}
	// 將回應傳入給請求中，解除阻塞狀況。
	req.response <- resp
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
	params, err := msgpack.Marshal([]string{event, channel})
	if err != nil {
		return err
	}
	return c.writeMessage(Request{
		Method: "MegoSubscribe",
		Params: params,
	})
}

// Unsubscribe 會取消訂閱指定的遠端事件，避免接收到無謂的事件。
func (c *Client) Unsubscribe(event string, channel string) error {
	params, err := msgpack.Marshal([]string{event, channel})
	if err != nil {
		return err
	}
	return c.writeMessage(Request{
		Method: "MegoUnsubscribe",
		Params: params,
	})
}

// On 能夠監聽系統或透過 `Subscribe` 所訂閱的事件，並做出相對應的動作。
func (c *Client) On(event string, handler func(*Event)) *Client {
	c.listeners[event] = handler
	return c
}

// Off 能移除當初以 `On` 所新增的事件監聽器，但這仍會繼續接收事件資料。
// 若要停止接收事件資料請使用 `Unsubscribe` 函式。
func (c *Client) Off(event string) *Client {
	delete(c.listeners, event)
	return c
}
