# JavaScript 客戶端

Mego 已附帶了 JavaScript 客戶端，可供你在前端與伺服端的 Mego 進行溝通或是單元測試。

# 索引

* [連線](#連線)
    * [重啟連線](#重啟連線)
    * [斷開連線](#斷開連線)
* [呼叫伺服端](#呼叫伺服端)
    * [設置酬載](#設置酬載)
    * [送出資料](#送出資料)
* [檔案上傳](#檔案上傳)
    * [透過檔案欄位、File、Blob](#透過檔案欄位、File、Blob)
    * [透過 FileReader、ArrayBuffer](#透過-FileReader、ArrayBuffer)
    * [自訂欄位名稱](#自訂欄位名稱)
    * [區塊上傳](#區塊上傳)
* [事件監聽](#事件監聽)
    * [訂閱自訂事件](#訂閱自訂事件)
    * [取消訂閱](#取消訂閱)
* [通知](#通知)

## 連線

透過 `MegoClient` 類別建立一個到伺服器的 Mego 連線。在另一個參數可夾帶客戶端資料，這令你不需要每次發送資料都須額外參夾使用者資料。

```js
try {
	// 建立一個新的 Mego 客戶端，並經由 WebSocket 溝通。
	var ws = new MegoClient("ws://localhost/")
	// 透過 `set` 保存本客戶端的資料至遠端，
	// 如此一來就不需要每次都發送相同重複資料。
	ws.set({
		"token"  : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		"version": "1.0"
	})

	ws.connect()
} catch (err) {
	console.log(err)
}
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

透過 `call` 可以呼叫伺服端並執行指定方法，且取得相關結果。

### 設置酬載

透過 `send` 在呼叫遠端方法之前設置欲傳遞資料，完成後以 `end` 來真正地發送請求。倘若傳入的資料是一個陣列則會被當作參數對待，伺服端可以更加快速地使用該資料。

如果傳入的資料是物件，伺服端則能以更有條理的方式處理資料但稍微複雜了些。傳入的資料如果是字串則會被當作 JSON 看待。

```js
var result

// 呼叫遠端的 `sum` 方法，並傳入兩個參數。
result = await ws.call("sum")
    .send([5, 3])
    .end()

// 呼叫遠端的 `sum` 方法，並傳入一個物件。
result = await ws.call("sum")
    .send({
        "numberA": 5
    })
    .end()

// 呼叫遠端的 `Sum` 方法，並傳入一個 JSON 內容。
result = await ws.call("sum")
    .send('{"numberA": 5}')
    .end()
```

### 送出資料

設置資料時並不會直接送出請求，需透過 `end` 才能發送已設置好的資料請求。由於 `end` 會回傳一個 Promise 物件，因此你可以使用 ES6 的 Async/Await。

```js
try {
    // 呼叫遠端的 `sum` 方法，並傳入兩個參數。
    var result = await ws.call("sum")
        .send([5, 3])
        .end()
} catch (err) {
    console.log(err)
}
```

普通的 `.then()` 也能來好好地對待請求回應。

```js
ws.call("sum")
    .send([5, 3])
    .end()
    .then((res) => {
        // ... 在此處理成功請求 ...
    })
    .catch((err) => {
        // ... 在此處理錯誤請求 ...
    })
```

## 檔案上傳

透過 `sendFile` 上傳檔案至伺服端，此用法非常彈性。

### 透過檔案欄位、File、Blob

將一個 `File` 物件傳入 `sendFile` 可以將其資料與檔案本身上傳至伺服端。

```js
// 以 jQuery 從 `<input>` 中取得檔案資料。
var file = $("input").get(0).files[0]

// 呼叫遠端的 `upload` 方法。
var result = await ws.call("upload")
    // 並且傳遞下列物件資料。
    .send({
        "albumID": 30
    })
    // 且夾帶檔案。
    .sendFile(file)
    .end()
```

### 透過 FileReader、ArrayBuffer

當 `sendFile` 的第一個參數是 `FileReader`、`ArrayBuffer` 型態時會被當作二進制資料看待並上傳至伺服端。需注意的是 Mego 客戶端將無法取得檔案的名稱，因此在伺服端會接收到一個空白的檔案名稱。

```js
// 初始化一個 `FileReader` 稍後從檔案欄位中取得二進制資料緩衝區。
var reader = new FileReader()
// 開始讀取檔案資料。
reader.readAsArrayBuffer(file)

// 呼叫遠端的 `upload` 方法。
var result = await ws.call("upload")
	// 並且傳遞下列物件資料。
	.send({
		"albumID": 30
	})
	// 以 `FileReader` 並轉換成二進制上傳檔案。
	.sendFile(reader)
	.end()
```

### 自訂欄位名稱

上傳檔案時 Mego 客戶端會自動以 `file1`、`file2`、`file3` 等遞加方式替檔案欄位進行命名。如果你想要手動指定檔案欄位名稱，請在 `sendFile` 的第二個參數指派。

```js
// 呼叫遠端的 `upload` 方法。
var result = await ws.call("upload")
    // 這個檔案欄位會被命名為 `File1`。
    .sendFile(file)
    // 這個檔案欄位會被命名為 `TextFile`。
    .sendFile(buffer, "textFile")
    // 這個檔案欄位會被命名為 `File2`。
    .sendFile(anotherFile)
    .end()
```

### 區塊上傳

透過 `sendFileChunks` 將一個大型檔案以區塊的方式上傳至遠端伺服器。和 `sendFile` 一樣的是你可以傳入 `File`、`FileReader`、`ArrayBuffer` 的型態資料到第一個參數。

有個需要注意的地方：使用區塊上傳時僅能夾帶一個檔案，且不可與其他 `sendFile` 並用。

```js
// 呼叫遠端的 `upload` 方法。
var result = await ws.call("upload")
	// 並且傳遞下列物件資料。
	.send({
		"albumID": 30
	})
	// 這個檔案欄位會被命名為 `File1`。
	.sendFileChunks(file)
	.end()
```

## 事件監聽

基本的內建事件有 `open`（連線）, `close`（連線關閉）, `reopen`（重新連線）, `message`（收到資料）, `error`（錯誤、斷線）可供監聽。

```js
ws.on("open", () => {
	console.log("成功連線！")
})
```

### 訂閱自訂事件

你也能以 `subscribe` 訂閱並透過 `on` 來監聽伺服器的自訂事件。用在聊天室檢查新訊息是最方便不過的了。

```js
// 先告訴伺服器我們要訂閱 `NewMessage` 事件。
ws.subscribe("newMessage").catch((err) => {
	console.log(err)
})

// 當接收到 `newMessage` 事件時所執行的處理函式。
ws.on("newMessage", () => {
	console.log("收到 newMessage 事件！")
})
```

### 取消訂閱

透過 `unsubscribe` 將自己從遠端伺服器上的指定事件監聽列表中移除。

```js
ws.unsubscrible("newMessage").catch((err) => {
	console.log(err)
})
```

## 通知

透過 `notify` 來向伺服器廣播一個通知，且不要求回應與夾帶資料。

```js
// 伺服器會執行名為 `newUser` 的方法。
ws.notify("newUser")
```