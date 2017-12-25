package mego

// Event 呈現了單一個事件。
type Event struct {
	// Name 是這個事件的名稱。
	Name string
	// Channels 是這個事件的所有頻道。
	Channels map[string]*Channel
	// engine 是這個事件所依存的主要引擎。
	engine *Engine
}

// Destroy 會摧毀一個事件和所有頻道避免其階段接收到相關事件。
func (e *Event) Destroy() {
	delete(e.engine.Events, e.Name)
}

// Channel 呈現了事件中的一個頻道與其監聽者。
type Channel struct {
	// Name 是這個頻道的名稱。
	Name string
	// Sessions 是監聽此頻道的階段 ID 切片。
	Sessions []*Session
	// Event 是這個頻道的父事件。
	Event *Event
}

// Destroy 會摧毀一個頻道避免其階段接收到相關事件。
func (c *Channel) Destroy() {
	delete(c.Event.Channels, c.Name)
}

// Disconnect 會斷開有註冊此頻道指定階段，避免繼續接收到相關事件。
func (c *Channel) Disconnect(id string) {
	for i, v := range c.Sessions {
		if v.ID == id {
			c.Sessions = append(c.Sessions[:i], c.Sessions[i+1:]...)
		}
	}
}
