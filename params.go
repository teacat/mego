package mego

import "time"

// Param 代表著單一個參數。
type Param struct {
	// data 是這個參數在型態轉換之前的介面。
	data interface{}
	// len 是參數陣列的長度。
	len int
}

// Len 能夠取得參數的數量。
func (p *Param) Len() int {
	return p.len
}

// Get 能夠以介面取得指定的參數。
func (p *Param) Get() (v interface{}) {
	return p.data
}

// GetBool 能夠以布林值取得指定的參數。
func (p *Param) GetBool() (v bool) {
	if val, ok := p.data.(bool); ok {
		v = val
	}
	return
}

// GetFloat64 能夠以浮點數取得指定的參數。
func (p *Param) GetFloat64() (v float64) {
	if val, ok := p.data.(float64); ok {
		v = val
	}
	return
}

// GetInt 能夠以正整數取得指定的參數。
func (p *Param) GetInt() (v int) {
	if val, ok := p.data.(int); ok {
		v = val
	}
	return
}

// GetString 能夠以字串取得指定的參數。
func (p *Param) GetString() (v string) {
	if val, ok := p.data.(string); ok {
		v = val
	}
	return
}

// GetStringMap 能夠以 `map[string]interface{}` 取得指定的參數。
func (p *Param) GetStringMap() (v map[string]interface{}) {
	if val, ok := p.data.(map[string]interface{}); ok {
		v = val
	}
	return
}

// GetMapString 能夠以 `map[string]string` 取得指定的參數。
func (p *Param) GetMapString() (v map[string]string) {
	if val, ok := p.data.(map[string]string); ok {
		v = val
	}
	return
}

// GetStringSlice 能夠以字串切片取得指定的參數。
func (p *Param) GetStringSlice() (v []string) {
	if val, ok := p.data.([]string); ok {
		v = val
	}
	return
}

// GetTime 能夠以時間取得指定的參數。
func (p *Param) GetTime() (v time.Time) {
	if val, ok := p.data.(time.Time); ok {
		v = val
	}
	return
}

// GetDuration 能夠以時間長度取得指定的參數。
func (p *Param) GetDuration() (v time.Duration) {
	if val, ok := p.data.(time.Duration); ok {
		v = val
	}
	return
}
