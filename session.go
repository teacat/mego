package mego

import "time"

// Session 是接收請求時的關聯內容，其包含了指向到特定客戶端的函式。
type Session struct {
	// Jar 包含了發送此請求的客戶端初始連線資料，此資料由客戶端連線時自訂。可用以取得用戶身份和相關資料。
	Jar map[string]interface{}
	// identifier 是此客戶端初始化時由伺服端所建立的不重複隨機名稱，供辨識匿名身份用。
	ID string
}

// Disconnect 會結束掉這個階段的連線。
func (s *Session) Disconnect() error {

}

// Copy 會複製一份 `Session` 供你在 Goroutine 中操作不會遇上資料競爭與衝突問題。
func (s *Session) Copy() *Session {

}

// Get 會取得客戶端當初建立連線時所傳入的資料特定鍵值組。
func (s *Session) Get(name string) (v interface{}, ok bool) {
	v, ok = s.Jar[name]
}

// GetBool 能夠以布林值取得指定的參數。
func (s *Session) GetBool(key string) (v bool) {
	if val, ok := s.Get(key); ok && val != nil {
		v, _ = val.(bool)
	}
	return
}

// GetFloat64 能夠以浮點數取得指定的參數。
func (s *Session) GetFloat64(key string) (v float64) {
	if val, ok := s.Get(key); ok && val != nil {
		v, _ = val.(float64)
	}
	return
}

// GetInt 能夠以正整數取得指定的參數。
func (s *Session) GetInt(key string) (v int) {
	if val, ok := s.Get(key); ok && val != nil {
		v, _ = val.(int)
	}
	return
}

// GetString 能夠以字串取得指定的參數。
func (s *Session) GetString(key string) (v string) {
	if val, ok := s.Get(key); ok && val != nil {
		v, _ = val.(string)
	}
	return
}

// GetStringMap 能夠以 `map[string]interface{}` 取得指定的參數。
func (s *Session) GetStringMap(key string) (v map[string]interface{}) {
	if val, ok := s.Get(key); ok && val != nil {
		v, _ = val.(map[string]interface{})
	}
	return
}

// GetMapString 能夠以 `map[string]string` 取得指定的參數。
func (s *Session) GetMapString(key string) (v map[string]string) {
	if val, ok := s.Get(key); ok && val != nil {
		v, _ = val.(map[string]string)
	}
	return
}

// GetStringSlice 能夠以字串切片取得指定的參數。
func (s *Session) GetStringSlice(key string) (v []string) {
	if val, ok := s.Get(key); ok && val != nil {
		v, _ = val.([]string)
	}
	return
}

// GetTime 能夠以時間取得指定的參數。
func (s *Session) GetTime(key string) (v time.Time) {
	if val, ok := s.Get(key); ok && val != nil {
		v, _ = val.(time.Time)
	}
	return
}

// GetDuration 能夠以時間長度取得指定的參數。
func (s *Session) GetDuration(key string) (v time.Duration) {
	if val, ok := s.Get(key); ok && val != nil {
		v, _ = val.(time.Duration)
	}
	return
}
