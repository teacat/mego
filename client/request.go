package client

import (
	"fmt"
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
	// isChunking 表示這個請求是否為區塊上傳。
	isChunking bool
	// originalParams 是個可供儲藏原先 `Params` 的地方，這是因為區塊上傳只有在最後一塊才會跟 `Params` 一起上傳。
	// 因此在那之前我們需要先把 `Params` 藏起來。
	originalParams []byte
}

// Send 會將稍後要保存的資料轉換成 MessagePack 格式並存入請求中。
func (r *Request) Send(data interface{}) *Request {
	params, err := msgpack.Marshal(data)
	if err != nil {

	}
	r.Params = params
	return r
}

// storeFile 會將檔案存放至請求建構體中並在送出時進行總結。
func (r *Request) storeFile(file interface{}, isChunk bool, fieldName ...string) {
	// 遞增檔案編號。
	r.client.fileID++

	// 初始化一個檔案。
	f := &File{
		ID:        r.client.fileID,
		source:    file,
		chunkSize: r.Option.ChunkSize,
	}
	// 裝載檔案內容，並且讀取下個區塊片段（如果是區塊上傳的話）。
	f.load(isChunk)
	f.next()

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
	_, ok := r.Files[n]
	if !ok {
		r.Files[n] = []*File{}
	}
	// 將本檔案推入指定檔案欄位中。
	r.Files[n] = append(r.Files[n], f)
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
	r.isChunking = true
	r.storeFile(file, true, fieldName...)
	return r
}

// timeout 會在指定的逾期時間後發送逾期錯誤給自己。
func (r *Request) timeout() {
	go func() {
		// 等待逾期時間。
		<-time.After(r.Option.Timeout)
		// 傳遞一個逾期回應給自己。
		r.response <- &Response{
			ID: r.ID,
			Error: Error{
				Code:    StatusTimeout,
				Message: ErrTimeout.Error(),
			},
		}
	}()
}

// End 結束並發送這個請求且不求回應。
func (r *Request) End() error {
	return r.EndStruct(nil)
}

// EndStruct 結束並發送這個請求，且將回應映射到本地建構體上。
func (r *Request) EndStruct(dest interface{}) error {
	var chunk *File
	// 如果本請求是區塊上傳，那麼先暫時把資料藏在緩衝區。
	// 只發送區塊檔案，直至最後一個區塊才將資料從緩衝區取出並一同上傳。
	if r.isChunking {
		for _, v := range r.Files {
			chunk = v[0]
		}
		// 如果這是第一個區塊，而且不是最後一塊。
		if chunk.Parts[0] == 0 && chunk.Parts[0] > chunk.Parts[1] {
			r.originalParams = r.Params
			r.Params = []byte{}
		}
		// 如果這是最後一塊。
		if chunk.Parts[0] == chunk.Parts[1] {
			r.Params = r.originalParams
			r.originalParams = []byte{}
		}
	}

	// 向伺服端發送請求。
	err := r.client.writeMessage(r)
	if err != nil {
		return err
	}
	// 啟動逾時檢查。
	//r.timeout()

	// 阻塞並等待此請求的回應。
	resp := <-r.response

	// 如果回應有事件名稱則依照相對應方法處理。
	switch resp.Event {
	case "MegoChunkNext":
		// 在區塊中載入下一段內容。
		chunk.next()
		// 呼叫自己重新發送相同的內容。
		if err := r.End(); err != nil {
			return err
		}
		return nil
	case "MegoChunkAbort":
		return ErrAborted
	}

	//
	if resp.Error.Code != 0 {
		return resp.Error
	}

	// 將最終回應映射到本地建構體上。
	if err := msgpack.Unmarshal(resp.Result, dest); err != nil {
		return err
	}
	return nil
}
