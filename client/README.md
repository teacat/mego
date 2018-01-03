# Golang 客戶端

Mego 已附帶了 Golang 客戶端，可供你在後端與伺服端的 Mego 進行溝通或是單元測試。

# 索引

* [連線](#連線)
    * [重啟連線](#重啟連線)
    * [斷開連線](#斷開連線)
* [呼叫伺服端](#呼叫伺服端)
    * [設置酬載](#設置酬載)
    * [送出資料](#送出資料)
* [檔案上傳](#檔案上傳)
    * [透過檔案路徑](#透過檔案路徑)
    * [透過 *os.File](#透過-*os.File)
    * [透過二進制](#透過二進制)
    * [自訂欄位名稱](#自訂欄位名稱)
    * [區塊上傳](#區塊上傳)
* [錯誤處理](#錯誤處理)
* [事件監聽](#事件監聽)
    * [訂閱自訂事件](#訂閱自訂事件)
    * [取消訂閱](#取消訂閱)

## 連線

透過 `MegoClient` 類別建立一個到伺服器的 Mego 連線。在另一個參數可夾帶客戶端資料，這令你不需要每次發送資料都須額外參夾使用者資料。

```go
// 建立一個新的 Mego 客戶端，並經由 WebSocket 溝通。
ws := client.NewClient("ws://localhost/").
	// 透過 `Set` 保存本客戶端的資料至遠端，
	// 如此一來就不需要每次都發送相同重複資料。
	Set(client.H{
		"Token":   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		"Version": "1.0",
	})

if err := ws.Connect(); err != nil  {
    panic(err)
}
```

### 重啟連線

透過 `Reconnect` 重啟一個新的連線。

```go
ws.Reconnect()
```

### 斷開連線

透過 `Close` 結束到伺服器的連線。

```go
ws.Close()
```

## 呼叫伺服端

透過 `Call` 可以呼叫伺服端並執行指定方法，且取得相關結果。

### 設置酬載

透過 `Send` 在呼叫遠端方法之前設置欲傳遞資料，完成後以 `End` 來真正地發送請求。倘若傳入的資料是一個陣列則會被當作參數對待，伺服端可以更加快速地使用該資料。

如果傳入的資料是物件，伺服端則能以更有條理的方式處理資料但稍微複雜了些。傳入的資料如果是字串則會被當作 JSON 看待。

```go
var resp Response

// 呼叫遠端的 `Sum` 方法，並傳入兩個參數。
err := ws.Call("Sum").
	Send([]int{5, 3}).
	End()

// 呼叫遠端的 `Sum` 方法，並傳入一個物件。
err = ws.Call("Sum").
	Send(client.H{
		"NumberA": 5,
	}).
	End()

// 呼叫遠端的 `Sum` 方法，並傳入一個 JSON 內容。
err = ws.Call("Sum").
	Send(`{"NumberA": 5}`).
	End()
```

### 送出資料

設置資料時並不會直接送出請求，需透過 `End` 才能發送已設置好的資料請求。

```go
var resp Response

// 呼叫遠端的 `Sum` 方法，並傳入兩個參數。
err := ws.Call("Sum").
	Send([]int{5, 3}).
	End()
```

以 `EndStruct` 則可以將設置好的資料送出並將接收到的回應映射到本地端的建構體。

```go
var resp Response

// 呼叫遠端的 `Sum` 方法，並傳入兩個參數。
// 且將接收到的資料映射到 `Response` 建構體。
err := ws.Call("Sum").
	Send([]int{5, 3}).
	EndStruct(&resp)
```

## 檔案上傳

透過 `SendFile` 上傳檔案至伺服端，此用法非常彈性。

### 透過檔案路徑

當 `SendFile` 的第一個參數是 `string` 型態時會被當作檔案路徑看待。Mego 客戶端會自動讀取該路徑的檔案二進制內容與其檔案名稱並上傳至伺服端。

```go
// 呼叫遠端的 `Upload` 方法。
err := ws.Call("Upload").
	// 並且傳遞下列物件資料。
	Send(client.H{
		"AlbumID": 30,
	}).
	// 且夾帶名為 `photo.png` 的檔案。
	SendFile("./photo.png").
	End()
```

### 透過 *os.File

Mego 客戶端能在 `SendFile` 的第一個參數是 `*os.File` 時自動讀取其內容，並且取得該檔案名稱一同上傳至遠端。

```go
// 開啟一個檔案並取得其 *os.File 資料。
file, _ := os.Open("video.mp4")

// 呼叫遠端的 `Upload` 方法。
err := ws.Call("Upload").
	// 並且傳遞下列物件資料。
	Send(client.H{
		"AlbumID": 30,
	}).
	// 以 *os.File 方式夾帶名為 `video.mp4` 的檔案。
	SendFile(file).
	End()
```

### 透過二進制

你亦能直接傳入二進制內容至 `SendFile` 的第一個參數，需注意的是 Mego 客戶端將無法取得檔案的名稱，因此在伺服端會接收到一個空白的檔案名稱。

```go
// 開啟一個檔案並取得其二進制資料。
path, _ := filepath.Abs("./textfile.txt")
bytes, _ := ioutil.ReadFile(path)

// 呼叫遠端的 `Upload` 方法。
err := ws.Call("Upload").
	// 並且傳遞下列物件資料。
	Send(client.H{
		"Author": "YamiOdymel",
	}).
	// 以二進制方式夾帶並上傳 `textfile.txt` 的內容。
	SendFile(bytes).
	End()
```

### 自訂欄位名稱

上傳檔案時 Mego 客戶端會自動以 `File1`、`File2`、`File3` 等遞加方式替檔案欄位進行命名。如果你想要手動指定檔案欄位名稱，請在 `SendFile` 的第二個參數指派。

```go
// 呼叫遠端的 `Upload` 方法。
err := ws.Call("Upload").
	// 這個檔案欄位會被命名為 `File1`。
	SendFile(file).
	// 這個檔案欄位會被命名為 `TextFile`。
	SendFile(bytes, "TextFile").
	// 這個檔案欄位會被命名為 `File2`。
	SendFile(anotherFile).
	End()
```

### 區塊上傳

透過 `SendFileChunks` 將一個大型檔案以區塊的方式上傳至遠端伺服器。和 `SendFile` 一樣的是你可以傳入 `[]byte`、`string`、`*os.File` 的型態資料到第一個參數。

有個需要注意的地方：使用區塊上傳時僅能夾帶一個檔案，且不可與其他 `SendFile` 並用。

```go
// 呼叫遠端的 `Upload` 方法。
err := ws.Call("Upload").
	// 並且傳遞下列物件資料。
	Send(client.H{
		"AlbumID": 30,
	}).
	// 這個檔案欄位會被命名為 `File1`。
	SendFileChunks("./largeFile.zip").
	End()
```

## 錯誤處理

```go
err := ws.Call("Sum").
	Send([]int{5, 3}).
	End()
if err != nil {

}
```

## 事件監聽

基本的內建事件有 `Open`（連線）, `Close`（連線關閉）, `Reopen`（重新連線）, `Message`（收到資料）, `Error`（錯誤、斷線）可供監聽。

```go
ws.On("Open", func() {
	fmt.Println("成功連線！")
})
```

### 移除監聽器

透過 `Off` 可以移除先前新增的監聽函式，但這仍會接收到來自伺服端的事件。欲要完全終止請使用 `Unsubscribe`。

```go
ws.Off("Open")
```

### 訂閱自訂事件

透過 `Subscribe` 告訴遠端伺服器我們要訂閱自訂事件與指定頻道，並透過 `On` 處理接收到的事件。一個事件可以有很多個頻道，類似一個「新訊息」事件來自很多個聊天室。

```go
// 先告訴伺服器我們要訂閱 `NewMessage` 事件，且表明這是頻道 `Chatroom1`。
// 指定頻道可以避免接收到其他人的事件。
err := ws.Subscribe("NewMessage", "Chatroom1")
if err != nil {
	panic(err)
}

// 當接收到 `NewMessage` 事件時所執行的處理函式。
ws.On("NewMessage", func(e *client.Event) {
	fmt.Println("收到 newMessage 事件！")
})
```

### 取消訂閱

透過 `Unsubscribe` 將自己從遠端伺服器上的指定事件監聽列表中移除，這將會停止接收到之後的事件。

```go
err := ws.Unsubscribe("NewMessage", "Chatroom1")
```