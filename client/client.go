package client

// Client 呈現了一個客戶端。
type Client struct {
}

// New 會初始化一個客戶端。
func New() *Client {
	return &Client{}
}

// Connect 會連線到遠端伺服器。
func (c *Client) Connect() {

}

// Reconnect 可以重新連線。
func (c *Client) Reconnect() {

}

// Close 會關閉連線。
func (c *Client) Close() {

}

// Call 會傳送指定資料並呼叫遠端的方法。當資料是一個陣列時，會被當作參數看待。
func (c *Client) Call() {

}

// On 可以監聽指定事件。
func (c *Client) On() {

}

// Subscribe 會向伺服器訂閱指定事件，如此一來才能夠監聽自訂事件。
func (c *Client) Subscribe() {

}

// Unsubscribe 會向伺服器取消訂與一個事件。
func (c *Client) Unsubscribe() {

}

// Upload 會上傳一個檔案到遠端伺服器。
func (c *Client) Upload() {

}

// Notify 來向伺服器廣播一個通知，且不要求回應與夾帶資料。
func (c *Client) Notify() {

}
