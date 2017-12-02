package mego

// Event 呈現了單一個事件與其監聽者。
type Event struct {
	// Name 是這個事件的名稱。
	Name string
	// Sessions 是監聽此事件的階段 ID 切片。
	Sessions []string
}

// Destroy 會摧毀一個事件避免其階段接收到相關事件。
func (e *Event) Destroy() {

}

// Disconnect 會斷開有註冊此事件指定階段，避免繼續接收到相關事件。
func (e *Event) Disconnect(id string) {

}
