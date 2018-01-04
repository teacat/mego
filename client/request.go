package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/vmihailenco/msgpack"
)

// RequestOption 呈現了一個請求的設置。
type RequestOption struct {
	// ChunkSize 是此請求區塊分切基準位元組大小，會取代原先的客戶端設置。
	ChunkSize int
	// Timeout 是此請求逾期的秒數，會取代原先的客戶端設置。
	Timeout time.Duration
	// UploadTimeout 是每個區塊、所有檔案的上傳逾期秒數，`0` 表示無上限，會取代原先的客戶端設置。
	UploadTimeout time.Duration
}

type Response struct {
	// Event 是欲呼叫的客戶端事件名稱。
	Event string `codec:"v" msgpack:"v"`
	// Result 是正常回應時的資料酬載。
	Result []byte `codec:"r" msgpack:"r"`
	// Error 是錯誤回應時的資料酬載，與 Result 兩者擇其一，不會同時使用。
	Error Error `codec:"e" msgpack:"e"`
	// ID 是當時發送此請求的編號，用以讓客戶端比對是哪個請求所造成的回應。
	ID int `codec:"i" msgpack:"i"`
}

// Request 呈現了一個籲發送至遠端伺服器的請求。
type Request struct {
	// Method 是欲呼叫的方法名稱。
	Method string `codec:"m" msgpack:"m"`
	// Params 是資料或參數。
	Params []byte `codec:"p" msgpack:"p"`
	// Files 是此請求所包含的檔案欄位與其內容。
	Files map[string][]*File `codec:"f" msgpack:"f"`
	// ID 為本次請求編號，若無則為單次通知廣播不需回應。
	ID int `codec:"i" msgpack:"i"`
	// Option 是這個請求的選項設置。
	Option *RequestOption

	// response 是這個請求的回應。
	response chan *Response
	// client 是建立這個請求的客戶端。
	client *Client
	// fileNameID 是自動遞增的檔案欄位編號。
	fileNameID int
}

// Send 會將稍後要保存的資料轉換成 MessagePack 格式並存入請求中。
func (r *Request) Send(data interface{}) *Request {
	params, err := msgpack.Marshal(data)
	if err != nil {

	}
	r.Params = params
	return r
}

//
func (r *Request) storeFile(file interface{}, isChunk bool, fieldName ...string) {
	var f *File
	var total int

	// 遞增檔案編號。
	r.client.fileID++

	// 依照傳入的檔案型態做出不同的讀取方式。
	switch v := file.(type) {
	// *os.File 表示檔案寫入者。
	case *os.File:
		bin, err := ioutil.ReadAll(v)
		if err != nil {

		}
		// 如果本檔案要以區塊上傳，就計算區塊切片數。
		if isChunk {
			//total
		}
		f = &File{
			Binary: bin,
			ID:     r.client.fileID,
			Name:   filepath.Base(v.Name()),
		}

	// 位元組表示檔案二進制內容。
	case []byte:
		f = &File{
			Binary: v,
			ID:     r.client.fileID,
		}

	// 字串型態表示檔案路徑。
	case string:
		bin, err := ioutil.ReadFile(v)
		if err != nil {

		}
		f = &File{
			Binary: bin,
			ID:     r.client.fileID,
			Name:   filepath.Base(v),
		}
	}

	// 取得檔案欄位名稱，若無指定則自動編號取名。
	var n string
	if len(fieldName) > 0 {
		n = fieldName[0]
	} else {
		// 遞增檔案名稱編號。
		r.fileNameID++
		n = fmt.Sprintf("File%d", r.fileNameID)
	}

	// 取得指定檔案欄位並檢查是否存在，不存在則初始化。
	field, ok := r.Files[n]
	if !ok {
		r.Files[n] = []*File{}
	}
	// 將本檔案推入指定檔案欄位中。
	r.Files[n] = append(r.Files[n], f)
}

func (r *Request) readFile() {

}

// SendFile 會保存稍後將上傳的檔案。
func (r *Request) SendFile(file interface{}, fieldName ...string) *Request {
	r.storeFile(file, false, fieldName...)
	return r
}

// SendFiles 能夠保存多個檔案並將其歸納為同個檔案欄位。
func (r *Request) SendFiles(files []interface{}, fieldName ...string) *Request {
	for _, v := range files {
		r.SendFile(v, fieldName...)
	}
	return r
}

// SendFileChunks 會保存稍後將以區塊方式上傳的檔案。
// 注意：請求使用區塊檔案上傳時，不可使用 `SendFile` 夾帶其他檔案。
func (r *Request) SendFileChunks(file interface{}, fieldName ...string) *Request {
	r.storeFile(file, true, fieldName...)
	return r
}

// End 結束並發送這個請求且不求回應。
func (r *Request) End() error {
	// 向伺服端發送請求。
	err := r.client.writeMessage(r)
	if err != nil {
		return err
	}
	// 阻塞並等待此請求的回應。
	resp := <-r.response
	if resp.Error.Code != 0 {
		return resp.Error
	}
	// 如果回應有事件名稱則依照相對應方法處理。
	switch resp.Event {
	case "MegoChunkNext":
	case "MegoChunkAbort":
		return ErrAborted
	}

	return nil
}

// EndStruct 結束並發送這個請求，且將回應映射到本地建構體上。
func (r *Request) EndStruct(dest interface{}) error {
	// 向伺服端發送請求。
	err := r.client.writeMessage(r)
	if err != nil {
		return err
	}
	// 阻塞並等待此請求的回應。
	response := <-r.response
	if response.Error.Code != 0 {
		return response.Error
	}
	if err := msgpack.Unmarshal(response.Result, dest); err != nil {
		return err
	}
	return nil
}
