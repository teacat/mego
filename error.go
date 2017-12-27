package mego

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type errorMsgs []*Error

// ErrorType 是錯誤種類。
type ErrorType uint64

var (
	// ErrEventNotFound 表示欲發送的事件沒有被初始化或任何客戶端監聽而無法找到因此發送失敗。
	ErrEventNotFound = errors.New("the event doesn't exist")
	// ErrChannelNotFound 表示欲發送的事件存在，但目標頻道沒有被初始化或任何客戶端監聽而無法找到因此發送失敗。
	ErrChannelNotFound = errors.New("the channel doesn't exist")
	// ErrFileNotFound 表示欲取得的檔案並不存在，可能是客戶端上傳不完整。
	ErrFileNotFound = errors.New("the file was not found")
	// ErrKeyNotFound 表示欲從鍵值組中取得的鍵名並不存在。
	ErrKeyNotFound = errors.New("the key was not found")
	// ErrPanicRecovered 表示 Panic 發生了但已回復正常。
	ErrPanicRecovered = errors.New("panic recovered")
)

const (
	// ErrorTypeBind 是 `Bind` 錯誤。
	ErrorTypeBind ErrorType = 1 << 63
	// ErrorTypePrivate 是建議僅供開發內部顯示的錯誤。
	ErrorTypePrivate ErrorType = 1 << 0
	// ErrorTypePublic 是可公開的錯誤。
	ErrorTypePublic ErrorType = 1 << 1
	// ErrorTypeAny 是任意、不限公開或較私密的錯誤。
	ErrorTypeAny ErrorType = 1<<64 - 1
)

// Error 呈現了一個路由中發生的錯誤。
type Error struct {
	// Err 是實際的錯誤資料。
	Err error
	// Meta 是這個錯誤的中繼資料，可傳遞至其他函式讀取運用。
	Meta interface{}
	// Type 是這個錯誤的種類，預設為 `ErrorTypePrivate`。
	Type ErrorType
}

// Error 會回傳一個錯誤的代表訊息。
func (e *Error) Error() string {
	return e.Err.Error()
}

// IsType 會表示此錯誤是否為指定的種類。
func (e *Error) IsType(flags ErrorType) bool {
	return (e.Type & flags) > 0
}

// SetMeta 可以將自訂的資料保存至錯誤建構體中傳遞到其他地方進行運用。
func (e *Error) SetMeta(data interface{}) *Error {
	e.Meta = data
	return e
}

func (e *Error) JSON() interface{} {
	json := H{}
	if e.Meta != nil {
		value := reflect.ValueOf(e.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return e.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				json[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			json["meta"] = e.Meta
		}
	}
	if _, ok := json["error"]; !ok {
		json["error"] = e.Error()
	}
	return json
}

// MarshalJSON implements the json.Marshaller interface.
func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.JSON())
}

// ByType 回傳一個以指定錯誤種類過濾後的現存錯誤清單。
func (a errorMsgs) ByType(t ErrorType) errorMsgs {
	if len(a) == 0 {
		return nil
	}
	if t == ErrorTypeAny {
		return a
	}
	var result errorMsgs
	for _, msg := range a {
		if msg.IsType(t) {
			result = append(result, msg)
		}
	}
	return result
}

// Last 回傳最後發生的一個錯誤，如果沒有錯誤發生則是 `nil`。
// 簡單說這就是 `errors[len(errors)-1]` 的縮寫用法。
func (a errorMsgs) Last() *Error {
	if length := len(a); length > 0 {
		return a[length-1]
	}
	return nil
}

// Errors 會回傳一個錯誤訊息切片。
// 舉例來說：
// 		c.Error(errors.New("一"))
// 		c.Error(errors.New("二"))
// 		c.Error(errors.New("三"))
// 		c.Errors.Errors() // == []string{"一", "二", "三"}
func (a errorMsgs) Errors() []string {
	if len(a) == 0 {
		return nil
	}
	errorStrings := make([]string, len(a))
	for i, err := range a {
		errorStrings[i] = err.Error()
	}
	return errorStrings
}

//
func (a errorMsgs) JSON() interface{} {
	switch len(a) {
	case 0:
		return nil
	case 1:
		return a.Last().JSON()
	default:
		json := make([]interface{}, len(a))
		for i, err := range a {
			json[i] = err.JSON()
		}
		return json
	}
}

func (a errorMsgs) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.JSON())
}

// String 會從錯誤切片輸出成一個可供人類閱讀的錯誤列表字串，其中亦包含所有的中繼資料。
func (a errorMsgs) String() string {
	if len(a) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for i, msg := range a {
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", i+1, msg.Err)
		if msg.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", msg.Meta)
		}
	}
	return buffer.String()
}
