package mego

// Request 呈現了一個客戶端所傳送過來的請求內容。
type Request struct {
	// Method 是欲呼叫的方法名稱。
	Method string
	// Params 是資料或參數。
	Params []byte
	// ID 為本次請求編號，若無則為單次通知廣播不需回應。
	ID int
}

// Response 呈現了 Mego 將會回應給客戶端的內容。
type Response struct {
	// Event 是欲呼叫的客戶端事件名稱。
	Event string
	// Result 是正常回應時的資料酬載。
	Result interface{}
	// Error 是錯誤回應時的資料酬載，與 Result 兩者擇其一，不會同時使用。
	Error ResponseError
	// ID 是當時發送此請求的編號，用以讓客戶端比對是哪個請求所造成的回應。
	ID int
}

// ResponseError 是回應錯誤資料建構體。
type ResponseError struct {
	// Code 是錯誤代號。
	Code int
	// Message 是人類可讀的簡略錯誤訊息。
	Message string
	// Data 是本次錯誤的詳細資料。
	Data interface{}
}

// Chunk 是一個檔案區塊。
type Chunk struct {
	// Name 是這個區塊應該被推入的目的檔案名稱。
	Name string
	// Part 是檔案的分塊編號。
	Part int
	// Bin 是此區塊的二進制資料。
	Bin []byte
}
