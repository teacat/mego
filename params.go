package mego

import "time"

// Params 保存著接收到的參數。
type Params struct {
	// data 是接收到的參數切片。
	Data []interface{}
}

// Len 能夠取得參數的數量。
func (p *Params) Len() int {
	return len(p.Data)
}

// Get 能夠以介面取得指定的參數。
func (p *Params) Get(index int) (v interface{}, has bool) {
	if !p.Has(index) {
		v, has = nil, false
		return
	}
	v, has = p.Data[index], true
	return
}

// GetBool 能夠以布林值取得指定的參數。
func (p *Params) GetBool(index int) (v bool) {
	if val, ok := p.Get(index); ok && val != nil {
		v, _ = val.(bool)
	}
	return
}

// GetFloat64 能夠以浮點數取得指定的參數。
func (p *Params) GetFloat64(index int) (v float64) {
	if val, ok := p.Get(index); ok && val != nil {
		v, _ = val.(float64)
	}
	return
}

// GetInt 能夠以正整數取得指定的參數。
func (p *Params) GetInt(index int) (v int) {
	if val, ok := p.Get(index); ok && val != nil {
		v, _ = val.(int)
	}
	return
}

// GetString 能夠以字串取得指定的參數。
func (p *Params) GetString(index int) (v string) {
	if val, ok := p.Get(index); ok && val != nil {
		v, _ = val.(string)
	}
	return
}

// GetStringMap 能夠以 `map[string]interface{}` 取得指定的參數。
func (p *Params) GetStringMap(index int) (v map[string]interface{}) {
	if val, ok := p.Get(index); ok && val != nil {
		v, _ = val.(map[string]interface{})
	}
	return
}

// GetMapString 能夠以 `map[string]string` 取得指定的參數。
func (p *Params) GetMapString(index int) (v map[string]string) {
	if val, ok := p.Get(index); ok && val != nil {
		v, _ = val.(map[string]string)
	}
	return
}

// GetStringSlice 能夠以字串切片取得指定的參數。
func (p *Params) GetStringSlice(index int) (v []string) {
	if val, ok := p.Get(index); ok && val != nil {
		v, _ = val.([]string)
	}
	return
}

// GetTime 能夠以時間取得指定的參數。
func (p *Params) GetTime(index int) (v time.Time) {
	if val, ok := p.Get(index); ok && val != nil {
		v, _ = val.(time.Time)
	}
	return
}

// GetDuration 能夠以時間長度取得指定的參數。
func (p *Params) GetDuration(index int) (v time.Duration) {
	if val, ok := p.Get(index); ok && val != nil {
		v, _ = val.(time.Duration)
	}
	return
}

// Has 會返回一個指定的參數是否存在的布林值。
func (p *Params) Has(index int) bool {
	return len(p.Data) >= index
}
