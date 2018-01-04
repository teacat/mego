package client

// Event 呈現了一個接收到的事件。
type Event struct {
	// Data 是事件所夾帶的資料。基於 Message Pack 格式需要透過 `Bind` 映射到本地建構體。
	Data []byte
}

// Bind 能將事件所夾帶的資料映射到本地的建構體。
func (e *Event) Bind(dest interface{}) error {
	return nil
}
