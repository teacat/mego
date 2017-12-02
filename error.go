package mego

// Error 呈現了一個路由中發生的錯誤。
type Error struct {
	Err      error
	Metadata interface{}
}

// Error 會回傳一個錯誤的代表訊息。
func (e Error) Error() string {
	return e.Err.Error()
}
