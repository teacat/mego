# JavaScript 客戶端

Mego 已附帶了 JavaScript 客戶端，可供你在 Node.js 或瀏覽器中完好地與後端 Mego 進行溝通。

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

```js
// 建立一個新的 Mego 客戶端，並經由 WebSocket 溝通。
ws = new MegoClient('ws://localhost/')

// 第二個參數可以擺放資料供你在伺服器端更方便處理使用者身份。
ws = new MegoClient('ws://localhost/', {
	token  : 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9',
	version: '1.0'
})
```

### 重啟連線

透過 `reconnect` 重啟一個新的連線。

```js
ws.reconnect()
```

### 斷開連線

透過 `close` 結束到伺服器的連線。

```js
ws.close()
```

## 呼叫伺服端

透過 `call` 可以呼叫伺服端並執行指定方法，且取得相關結果。倘若傳入的資料是一個陣列則會被當作參數對待，伺服端可以更加快速地使用該資料。

如果傳入的資料是物件，伺服端則能以更有條理的方式處理資料但稍微複雜了些。

```js
// 呼叫遠端的 `sum` 方法，並傳入兩個參數。
ws.call('sum', [5, 3]).then(({result}) => {
	// ... 執行完畢 ...
})

// 呼叫遠端的 `createUser` 方法，並傳入一個物件。
ws.call('createUser', {
	username: 'YamiOdymel'
}).then(({result}) => {
	// ... 執行完畢 ...
})
```

### 透過 ES7 的 Async、Await

Mego 使用 Promise 讓你能夠以 ES7 的 Async 異步函式良好地處理程式流程。

```js
// 像這樣。
var result = await ws.call('sum', [5, 3])

// 或這樣。
var result = await ws.call('createUser', {
	username: 'YamiOdymel'
})
```

### 錯誤處理

當遠端伺服器回傳的是一個錯誤時，Mego 會執行 `.catch` 而不是 `.then` 中的函式。這點與 JavaScript 標準的 Promise 用法一致。

```js
ws.call('sum', [5, 3]).then(({result}) => {
	// ... 正常回應 ...
}).catch((error) => {
	// ... 錯誤回應 ...
})
```

使用 ES7 的 Async 時，則需要透過 `try/catch` 區塊包覆。

```js
try {
	var result = await ws.call('sum', [5, 3])
	// ... 正常回應 ...
} catch(error) {
	// ... 錯誤回應 ...
}
```

## 事件監聽

基本的內建事件有 `open`（連線）, `close`（連線關閉）, `reopen`（重新連線）, `message`（收到資料）, `error`（錯誤、斷線）可供監聽。

```js
ws.on('open', () => {
	console.log('成功連線！')
})
```

### 訂閱自訂事件

你也能以 `subscribe` 訂閱並透過 `on` 來監聽伺服器的自訂事件。用在聊天室檢查新訊息是最方便不過的了。

```js
// 先告訴伺服器我們要訂閱 `newMessage` 事件。
ws.subscribe('newMessage')

// 當接收到 `newMessage` 事件時所執行的處理函式。
ws.on('newMessage', ({result}) => {
	console.log('收到 newMessage 事件！')
})
```

### 取消訂閱

透過 `unsubscribe` 將自己從遠端伺服器上的指定事件監聽列表中移除。

```js
ws.unsubscribe('newMessage')
```

## 檔案上傳



```js
ws.upload("photo", file).then(({result}) => {
	console.log('檔案上傳完畢。')
})
```

## 通知

透過 `notify` 來向伺服器廣播一個通知，且不要求回應與夾帶資料。

```js
// 伺服器會執行名為 `newUser` 的方法。
ws.notify('newUser')
```