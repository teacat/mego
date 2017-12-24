package client

func NewClient(url string) *Client {
	return &Client{
		URL: url,
	}
}

type Client struct {
	URL string
}

func (c *Client) Call(method string) *Request {
	return &Request{
		Method: method,
	}
}

func (c *Client) Set(data interface{}) *Client {

}

func (c *Client) Connect() error {

}

func (c *Client) Reconnect() error {

}

func (c *Client) Close() error {

}

func (c *Client) Subscribe(event string, channel string) error {

}

func (c *Client) Unsubscribe(event string, channel string) error {

}

func (c *Client) On(event string, handler func(*Event)) *Client {

}

type Error struct {
}

func (e *Error) Error() string {

}

type Event struct {
}

func (e *Event) Bind(dest interface{}) error {
	return nil
}

type File struct {
}

type Request struct {
	Method string
	Params interface{}
	Files  []*File
	ID     int
}

func (r *Request) Send() {

}

func (r *Request) SendFile() {

}

func (r *Request) SendFileChunks() {

}

func (r *Request) End() error {
	return nil
}

func (r *Request) EndStruct(dest interface{}) error {
	return nil
}
