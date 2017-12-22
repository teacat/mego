# Golang 客戶端

Mego 已附帶了 Golang 客戶端，可供你在後端與伺服端的 Mego 進行溝通或是單元測試。

# 索引

* [連線](#連線)
    * [重啟連線](#重啟連線)
    * [斷開連線](#斷開連線)
* [呼叫伺服端](#呼叫伺服端)
    * [透過 ES7 的 Async、Await](#透過-es7-的-async、await)
    * [錯誤處理](#錯誤處理)
* [事件監聽](#事件監聽)
    * [訂閱自訂事件](#訂閱自訂事件)
    * [取消訂閱](#取消訂閱)
* [檔案上傳](#檔案上傳)
* [通知](#通知)

## 連線

透過 `MegoClient` 類別建立一個到伺服器的 Mego 連線。在另一個參數可夾帶客戶端資料，這令你不需要每次發送資料都須額外參夾使用者資料。

```go
// 建立一個新的 Mego 客戶端，並經由 WebSocket 溝通。
ws := client.New("ws://localhost/")

// 第二個參數可以擺放資料供你在伺服器端更方便處理使用者身份。
ws := client.New("ws://localhost/", client.H{
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

透過 `Call` 可以呼叫伺服端並執行指定方法，且取得相關結果。倘若傳入的資料是一個陣列則會被當作參數對待，伺服端可以更加快速地使用該資料。

如果傳入的資料是物件，伺服端則能以更有條理的方式處理資料但稍微複雜了些。

```go
// 呼叫遠端的 `sum` 方法，並傳入兩個參數。
result, err := ws.Call("Sum", []int{5, 3})

// 呼叫遠端的 `createUser` 方法，並傳入一個物件。
result, err := ws.Call("CreateUser", client.H{
	"Username": "YamiOdymel",
})
```

## 事件監聽

基本的內建事件有 `open`（連線）, `close`（連線關閉）, `reopen`（重新連線）, `message`（收到資料）, `error`（錯誤、斷線）可供監聽。

```go
ws.On("open", func() {
    fmt.Println("成功連線！")
})
```

### 訂閱自訂事件

你也能以 `Subscribe` 訂閱並透過 `On` 來監聽伺服器的自訂事件。用在聊天室檢查新訊息是最方便不過的了。

```go
// 先告訴伺服器我們要訂閱 `NewMessage` 事件。
ws.Subscribe("NewMessage")

// 當接收到 `NewMessage` 事件時所執行的處理函式。
ws.On("NewMessage", func() {
    fmt.Println("收到 newMessage 事件！")
})
```

### 取消訂閱

透過 `Unsubscribe` 將自己從遠端伺服器上的指定事件監聽列表中移除。

```go
ws.Unsubscribe("NewMessage")
```

## 檔案上傳


```go
//
file, err := os.Open("myFile.mp4")
if err != nil {

}

result, err := ws.Upload("Photo", file)
if err == nil {
    fmt.Println("檔案上傳成功！")
}
```

## 通知

透過 `Notify` 來向伺服器廣播一個通知，且不要求回應與夾帶資料。

```go
// 伺服器會執行名為 `NewUser` 的方法。
ws.Notify("NewUser")
```