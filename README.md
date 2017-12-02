&nbsp;

<p align="center">
  <img src="https://user-images.githubusercontent.com/7308718/33408556-0f04a56c-d5b2-11e7-8a14-b556cf283b43.png"/>
</p>
<p align="center">
  <i>Meet Go.</i>
</p>

&nbsp;

# // WIP - 尚未完成 //

# Mego

用以取代傳統 [RESTful](https://zh.wikipedia.org/zh-tw/REST) 溝通結構，並且以[類 RPC](https://zh.wikipedia.org/zh-tw/%E9%81%A0%E7%A8%8B%E9%81%8E%E7%A8%8B%E8%AA%BF%E7%94%A8) 的方式以 [MessagePack](http://msgpack.org/) 將資料進行壓縮並透過 [WebSocket](https://developer.mozilla.org/zh-TW/docs/WebSockets/WebSockets_reference/WebSocket) 傳遞使資料更加輕量。同時亦支援分塊檔案上傳。

# 這是什麼？

煤果是一個由 [Golang](https://golang.org/) 所撰寫的前後端雙向溝通方式，這很適合用於聊天室、即時連線、單頁應用程式。目前僅支援前端 JavaScript 與後端 Golang 的溝通。

* 可取代傳統 RESTful API。
* 以 MessagePack 編碼資料與達到更輕量的傳遞。
* 可雙向溝通與廣播以解決傳統僅單向溝通的障礙。
* 檔案分塊上傳方式。
* 支援中介軟體（Middleware）。
* 可使用 ES7 Async/Await 處理前端的呼叫。
* 友善的錯誤處理環境。

# 為什麼？

在部份時候傳統 RESTful API 確實能派上用場，但久而久之就會為了如何命名、遵循 REST 風格而產生困擾，例如網址必須是名詞而非動詞，登入必須是 `GET /token` 而非 `POST /login`，但使用 `GET` 傳遞機密資料極為不安全，因此只好更改為 `POST /session` 等，陷入如此地窘境。

且多個時候 RESTful API 都是單向，而非雙向。這造成了在即時連線、雙向互動上的溝通困擾，不得不使用 Comet、Long Polling 達到相關的要求，令 API 更加分散難以管理。

Mego 試圖以 WebSocket 解決雙向溝通和多數重複無用的 HTTP 標頭導致浪費頻寬的問題。Mego 同時亦支援透過 WebSocket 上傳檔案的功能。

# 檔案上傳與廣播

因為 Mego 是透過 WebSocket 進行溝通，這令你可以主動式廣播訊息至客戶端（讓你更友善地達成通知系統）。

Mego 亦支援檔案上傳。但在 Mego 中，上傳檔案和送出資料是分開的，這意味著當你想要上傳帶有圖片的表單時，你需要先上傳圖片，接著取得已上傳圖片的以檔案編號的方式夾帶到另一個表單方可傳遞相關資訊。這在上傳大型檔案如影片時非常有用。

```
[Client]      [Server]
    |             |      客戶端：建立不重複金鑰，將檔案切分成塊。
    |------------>|      客戶端：傳送區塊 1／2 與金鑰。
    |<------------|      伺服器：完成 #1，請傳送下一塊。
    |------------>|      客戶端：傳送區塊 2／2 與金鑰。
    |             |      伺服器：組合所有區塊。
    |<------------|      伺服器：上傳程序完成，呼叫完成函式進行檔案處理（如：縮圖、轉檔）並回傳檔案資料。
    |             |      客戶端：取得檔案編號，存至新表單資料。
   ~~~~~~新請求~~~~~~~
    |------------>|      客戶端：將帶有檔案編號的表單資料傳遞至伺服器。
    v             v
```

# 效能如何？

Mego 是基於 [`net/http`](https://golang.org/pkg/net/http/) 和 [`olahol/melody`](https://github.com/olahol/melody) 的 WebSocket 框架作為基礎，並由 [`vmihailenco/msgpack`](https://github.com/vmihailenco/msgpack) 作為傳遞訊息的基本格式。

這裡有份簡略化的效能測試報表。

```
測試規格：
1.7 GHz Intel Core i7 (4650U)
8 GB 1600 MHz DDR3
```

# 索引

# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get github.com/TeaMeow/Mego
```

# 開始使用

Mego 的使用方式參考令人容易上手的 [Gin 網站框架](https://github.com/gin-gonic/gin)，令你撰寫 WebSocket 就像在撰寫普通的網站框架ㄧ樣容易。

###### Golang

```go
import "github.com/TeaMeow/Mego"

func main() {
	// 初始化一個基本的 Mego 引擎。
	e := mego.Default()
	// 定義一個 User 方法供客戶端呼叫。
	e.Register("GetUser", func(c *mego.Context) {
		// 回應 `{"username": "YamiOdymel"}` 資料。
		c.Respond(mego.StatusOK, mego.H{
			"username": "YamiOdymel",
		})
	})
	// 執行引擎在埠口 5000 上。
	e.Run()
}
```

###### JavaScript

```javascript
ws = new Mego('ws://localhost/')

ws.on('open', () => {
    ws.call('getUser').then(({result}) => {
        console.log(result.username) // 輸出：YamiOdymel
    })
})
```
