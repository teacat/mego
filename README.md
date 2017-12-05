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

* [安裝方式](#安裝方式)
* [使用方式](#使用方式)
  * [初始化引擎](#初始化引擎)
  * [廣播與建立事件](#廣播與建立事件)
    * [多數廣播](#多數廣播)
    * [過濾廣播](#過濾廣播)
  * [映射資料與參數](#映射資料與參數)
    * [取得參數](#取得參數)
  * [處理請求與回應](#處理請求與回應)
    * [指定客戶端廣播事件](#指定客戶端廣播事件)
  * [中介軟體](#中介軟體)
    * [推遲執行與接續](#推遲執行與接續)
    * [終止請求](#終止請求)
    * [存放鍵值組](#存放鍵值組)
  * [檔案上傳處理](#檔案上傳處理)
  * [斷開客戶端](#斷開客戶端)
  * [複製並使用於 Goroutine](#複製並使用於-goroutine)
  * [錯誤處理](#錯誤處理)
  * [結束](#結束)
* [客戶端](#客戶端)
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
* [狀態碼](#狀態碼)
* [API 文件撰寫風格](#api-文件撰寫風格)


# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get github.com/TeaMeow/Mego
```

# 使用方式

Mego 的使用方式參考令人容易上手的 [Gin 網站框架](https://github.com/gin-gonic/gin)，令你撰寫 WebSocket 就像在撰寫普通的網站框架ㄧ樣容易。

###### Golang

```go
import "github.com/TeaMeow/Mego"

func main() {
	// 初始化一個基本的 Mego 引擎。
	e := mego.Default()
	// 定義一個 GetUser 方法供客戶端呼叫。
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

## 初始化引擎

初始化一個 Mego 引擎才能夠開始新增可供客戶端呼叫的方法，建立一個引擎有兩種方法。

```go
// 初始化一個帶有自動回復、紀錄中介軟體的基本引擎。
e := mego.Default()

// 初始化一個無預設中介軟體的引擎。
e := mego.New()
```

## 廣播與建立事件

由於 Mego 和傳統 HTTP 網站框架不同之處在於：Mego 透過 WebSocket 連線。這使你可以主動發送事件到客戶端，而不需要等待客戶端主動來發送請求。

```go
func main() {
	e := mego.Default()

	// 建立一個 `UpdateApp` 才能供客戶端監聽。
	e.Event("UpdateApp")

	// 廣播一個 `UpdateApp` 事件給所有監聽的客戶端，觸發其事件監聽函式。
	e.Emit("UpdateApp", nil)

	// 廣播一個 `UpdateApp` 事件並帶有指定的資料供客戶端讀取。
	e.Emit("UpdateApp", mego.H{
		"version": "1.0.0",
	})

	e.Run()
}
```

### 多數廣播

透過 `Emit` 會廣播指定事件給所有連線的客戶端，如果你希望廣播事件給指定的數個客戶端時，你可以透過 `EmitMultiple` 並傳入欲接收指定事件的客戶端階段達成。

```go
func main() {
	e := mego.Default()

	// 建立一個 `UpdateApp` 才能供客戶端監聽。
	e.Event("UpdateApp")

	// 對指定的 [0] 與 [1] 客戶端廣播 `UpdateApp` 事件，其他客戶端不會接收到此事件。
	e.EmitMultiple("UpdateApp", nil, []*mego.Session{
		e.Sessions[0],
		e.Sessions[1],
	})

	e.Run()
}
```

### 過濾廣播

同時，透過 `EmitFilter` 可以遍歷所有連線的客戶端，找出他們的相關資料並以此為依據決定是否要廣播指定事件給他們。

```go
func main() {
	e := mego.Default()

	// 建立一個 `UpdateApp` 才能供客戶端監聽。
	e.Event("UpdateApp")

	// 過濾所有客戶端，僅有 ID 為 `0k3LbEjL` 的才能夠接收到 UpdateApp 事件。
	e.EmitFilter("UpdateApp", nil, func(s *mego.Session) bool {
		// 此函式會過濾、遍歷每個客戶端，如果此函式回傳 `true`，此客戶端則會接收到此事件。
		return s.ID == "0k3LbEjL"
	})

	e.Run()
}
```

## 映射資料與參數

欲要接收客戶端傳來的資料，透過 `Bind` 可以將資料映射到本地的建構體。如果資料是重要且必須的，可以透過 `MustBind` 來映射資料，並在錯誤發生時自動呼叫 `panic` 終止此請求。

```go
// User 是一個使用者資料建構體。
type User struct {
	Username string
}

func main() {
	e := mego.Default()

	e.Register("CreateUser", func(c *mego.Context) {
		var u User

		// 將接收到的資料映射到本地的使用者資料建構體，這樣才能讀取其資料。
		if c.Bind(&u) == nil {
			fmt.Println(u.Username)
		}

		// 透過 `MustBind` 在映射錯誤時直接呼叫 `panic` 終止此請求繼續。
		c.MustBind(&u)
		fmt.Println(u.Username)
	})

	e.Run()
}
```

### 取得參數

當客戶端傳送來的不是資料，而是普通陣列時會被當作參數看待。參數可以讓你更快進行處理而不需映射。

透過 `Params.Get(index)` 取得指定索引的參數內容，亦有 `GetInt`、`GetString` 等方便的資料型態快捷函式。

```go
func main() {
	e := mego.Default()

	e.Register("Sum", func(c *mego.Context) {
		// 以正整數資料型態取得傳入的兩個參數。
		a := c.Params.GetInt(0)
		b := c.Params.GetInt(1)

		// 將兩個轉換成正整數的參數進行相加取得結果。
		fmt.Printf("A + B = %d", a+b)
	})

	e.Run()
}
```

## 處理請求與回應

透過 `Respond` 正常回應一個客戶端的請求。以 `Status` 可以僅傳送狀態碼並省去多餘的資料。而 `RespondWithError` 可以回傳一個錯誤發生的詳細資料用已告知客戶端發生了錯誤。

```go
func main() {
	e := mego.Default()

	e.Register("CreateUser", func(c *mego.Context) {
		// 針對此請求，回傳一個指定的狀態碼與資料。
		c.Respond(mego.StatusOK, mego.H{
			"message": "成功建立使用者！",
		})

		// 針對此請求，僅回傳一個狀態碼並省去多餘的資料內容。
		c.Status(mego.StatusOK)

		// 針對此請求，回傳一個錯誤相關資料與狀態碼，還有錯誤本身。
		c.RespondWithError(mego.StatusError, mego.H{
			"empty": "username",
		}, errors.New("使用者名稱不可為空。"))
	})

	e.Run()
}
```

### 指定客戶端廣播事件

在處理函式中使用上下文建構體的 `Emit` 可以僅對發送請求的客戶端進行指定的事件廣播。同時也可以透過上下文建構體內的 `EmitOthers` 來對此客戶端以外的所有其他人進行指定事件的廣播。

```go
func main() {
	e := mego.Default()

	// 建立一個 `UpdateApp` 才能供客戶端監聽。
	e.Event("UpdateApp")

	e.Register("CreateUser", func(c *mego.Context) {
		// 廣播 UpdateApp 事件給產生此請求的客戶端。
		c.Emit("UpdateApp", nil)

		// 廣播 UpdateApp 事件給產生此請求客戶端以外的「所有其他人」。
		c.EmitOthers("UpdateApp", nil)
	})

	e.Run()
}
```

## 中介軟體

透過中介軟體你可以很容易地集中管理一些函式，例如：請求驗證、連線紀錄、效能測量。簡單來說，中介軟體就是能夠在每個連線之前所執行的函式。

要謹記：在中介軟體中必須呼叫 `Next` 呼叫下一個中介軟體或函式，否則會無法進行下個動作。

```go
func main() {
	e := mego.Default()

	// logMw 是一個會紀錄請求編號的紀錄中介軟體。
	logMw := func(c *mego.Context) {
		fmt.Printf("%s 已連線。", c.ID)
		// ... 紀錄邏輯 ...
		c.Next()
	}

	// authMw 是一個驗證請求者是否有權限的身份檢查中介軟體。
	authMw := func(c *mego.Context) {
		if c.ID != "King" {
			// ... 驗證邏輯 ...
		}
		c.Next()
	}

	// 在 Mego 引擎中使用上述兩個全域中介軟體。
	e.Use(logMw, authMw)

	e.Run()
}
```

### 推遲執行與接續

`Next` 之後的程式會在其他中介軟體或處理函式執行完畢後反序執行，因此你可以透過這個特性測量一個請求到結束總共耗費了多少時間。

```go
func main() {
	e := mego.Default()

	mw := func(c *mego.Context) {
		// 從連線開始的時間。
		start := time.Now()
		// 呼叫下一個中介軟體或執行處理函式。
		c.Next()
		// 取得執行完畢後的時間。
		end := time.Now()

		// 將開始與結束的時間相減，取得此請求所耗費的時間。
		fmt.Printf("此請求共花費了 %v 時間", end.Sub(start))
	}

	// 在 Mego 引擎中使用上述全域中介軟體。
	e.Use(mw)

	e.Run()
}
```

### 終止請求

如果連線有異常不希望繼續呼叫下一個中介軟體或處理函式，則需要透過 `Abort` 終止此請求。透過 `AbortWithStatus` 可以在終止的時候回傳一個狀態碼給客戶端。以 `AbortWithRespond` 可以回傳正常的資料與狀態碼。

而 `AbortWithError` 則會在終止的時候回傳狀態碼與錯誤的詳細資料。

```go
func main() {
	e := mego.Default()

	mw := func(c *mego.Context) {
		if c.ID != "jn3Dl74eX" {
			// 直接終止此請求的繼續，並停止呼叫接下來的中介軟體與處理函式。
			c.Abort()

			// 終止此請求並回傳一個狀態碼。
			c.AbortWithStatus(mego.StatusNoPermission)

			// 終止此請求並且回傳狀態碼與資料。客戶端並不會知道是錯誤發生，仍會依一般回應處理。
			c.AbortWithRespond(mego.StatusOK, mego.H{
				"message": "嗨！雖然你不是 jn3Dl74eX 但我們還是很歡迎你！",
			})

			// 終止此請求並且回傳一個錯誤的資料與狀態碼和錯誤本體。
			c.AbortWithError(mego.StatusNoPermission, nil, errors.New("沒有相關權限可訪問。"))
		}
	}

	// 在 Mego 引擎中使用上述全域中介軟體。
	e.Use(mw)

	e.Run()
}
```

### 存放鍵值組

為了讓開發者更方便在不同中介軟體或處理函式中交換訊息，上下文建構體內建了可存放自訂值的鍵值組。透過 `Set` 存放指定的鍵值，而 `Get` 來取得指定的值。

如果有一個值是非常重要且必須的，透過 `MustGet` 來取得可以在找不到該值時自動呼叫 `panic` 終止此請求的繼續。

```go
func main() {
	e := mego.Default()

	mw := func(c *mego.Context) {
		// 在上下文建構體中存入名為 `Foo` 的 `Bar` 值。
		c.Set("Foo", "Bar")
		// 呼叫下一個中介軟體或處理函式。
		c.Next()
	}

	mw2 := func(c *mego.Context) {
		// 檢查是否有 `Foo` 這個鍵名，若沒有則終止此請求繼續。
		if v, ok := c.Get("Foo"); !ok {
			c.Abort()
		}
		fmt.Println(v) // 輸出：Bar
	}

	// 在 Mego 引擎中使用上述兩個全域中介軟體。
	e.Use(mw, mw2)

	e.Run()
}
```

和參數ㄧ樣，鍵值組也有相關的資料型態輔助函式，如：`GetString`、`GetInt`。

```go
func main() {
	e := mego.Default()

	mw := func(c *mego.Context) {
		fmt.Println(c.GetString("Foo")) // 輸出：Bar

		// c.GetInt()
		// c.GetBool()
		// c.GetStringMap()
		// ...
	}

	e.Run()
}
```

### 內建中介軟體

Mego 有內建數個方便的中介軟體能夠讓你限制同時連線數，或使用內建的簡易記錄系統⋯等。

如果你使用 `Default` 建立了一個 Mego 引擎，該引擎則包含自動回復與請求紀錄中介軟體。

#### 自動回復

單個方法發生 `panic` 時，`Recovery` 可以自動回復並在終端機顯示相關錯誤訊息供除錯。

此中介軟體能避免整個 Mego 引擎被強制關閉。

```go
func main() {
	e := mego.New()
	e.Use(mego.Recovery())
}
```

#### 請求紀錄

`Logger` 會在終端機即時顯示任何請求，方便開發者進行監控。

```go
func main() {
	e := mego.New()
	e.Use(mego.Logger())
}
```

#### 配額限制

如果你希望限制客戶端在單個方法的某個時間內僅能呼叫數次，透過 `RateLimit` 就能達成。

選項中的 `Period` 表示週期秒數，`Limit` 表示客戶端可在週期內呼叫幾次相同方法。

```go
func main() {
	e := mego.Default()

	e.Register("CreateUser", mego.RateLimit(mego.RateLimitOption{
		Period: 1 * time.Hour,
		Limit: 100
	}), func(c *mego.Context) {
		// ...
	})

	e.Run()
}
```

#### 同時請求數

為了避免某個方法特別繁忙，透過 `RequestLimit` 可以限制單個方法的最大同時連線數。

選項中的 `Limit` 表示單個方法最大僅能有幾個客戶端同時呼叫。

```go
func main() {
	e := mego.Default()

	e.Register("CreateUser", mego.RequestLimit(mego.RequestLimitConfig{
		Limit: 1000
	}), func(c *mego.Context) {
		// ...
	})

	e.Run()
}
```

## 檔案上傳處理

Mego 能夠自動幫你處理透過 WebSocket 所上傳的檔案。透過 `File` 建構體可以取得客戶端上傳的檔案資料。

```go
func main() {
	e := mego.Default()

	e.Recevie("Photo", func(c *mego.Context) {
		c.Respond(mego.StatusOK, mego.H{
			"filename":  c.File.Name,      // 檔案原始名稱。
			"size":      c.File.Size,      // 檔案大小（位元組）。
			"extension": c.File.Extension, // 檔案副檔名（無符號）。
			"path":      c.File.Path,      // 檔案本地路徑。
		})
	})

	e.Run()
}
```

## 斷開客戶端

欲要從伺服器斷開與指定客戶端的連線，請透過客戶端階段裡的 `Disconnect` 函式。

```go
func main() {
	e := mego.Default()

	// 斷開客戶端 `HkBE9lebt` 與伺服器的連線。
	e.Sessions["HkBE9lebt"].Disconnect()

	// 或者你也不需要指名道姓⋯⋯。
	e.Register("CreateUser", func(c *mego.Context) {
		// 直接斷開發送此請求的客戶端連線。
		c.Session.Disconnect()
	})

	e.Run()
}
```

## 複製並使用於 Goroutine

當你要將上下文建構體（Context）傳入 Goroutine 使用時，你必須透過 `Copy` 複製一份上下文建構體，這個建構體不會指向原本的上下文建構體，如此一來能夠避免資料競爭、衝突問題。

```go
func main() {
	e := mego.Default()

	e.Register("CreateUser", func(c *mego.Context) {
		// 複製一份目前的上下文建構體。
		copied := c.Copy()
		// 然後在 Goroutine 中使用這個複製體。
		go func() {
			time.Sleep(5 * time.Second)
			// ... 在這裡使用 copied 而非 c ...
		}()
	})

	e.Run()
}
```

## 錯誤處理

在中介軟體或函式中處理錯誤是非常簡單的一件事。透過 `Error` 存放或堆疊此請求所發生的所有錯誤，並能在最終一次讀取。

```go
func main() {
	e := mego.Default()

	// 這是一個錯誤中介軟體，會不斷地在上下文建構體中塞入一堆錯誤。
	errMw := func(c *mego.Context) {
		// 透過 `Error` 能將錯誤訊息堆放置上下文建構體中，
		// 並在最終一次獲得所有錯誤資訊。
		c.Error(errors.New("這是超重要的錯誤訊息啊！"))
		c.Error(errors.New("世界要毀滅啦！"))

		c.Next()
	}

	// 在 Mego 引擎中使用上述全域中介軟體。
	e.Use(errMw)

	e.Register("CreateUser", func(c *mego.Context) {
		// 檢查上下文建構體裡是否有任何錯誤。若有則回傳錯誤訊息給客戶端。
		if len(c.Errors) != 0 {
			c.AbortWithError(mego.StatusError, nil, c.Errors[0])
		}
	})

	e.Run()
}
```

## 結束

欲要結束整個 Mego 引擎，請使用 `Close` 函式。

```go
e.Close()
```

# 客戶端

Mego 已附帶了 JavaScript 客戶端，可供你在 Node.js 或瀏覽器中完好地與後端 Mego 進行溝通。

## 連線

透過 `MegoClient` 類別建立一個到伺服器的 Mego 連線。

```js
// 建立一個新的 Mego 客戶端，並經由 WebSocket 溝通。
ws = new MegoClient('ws://localhost/')
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

透過

```js
ws.file("photo", file).then(({result}) => {
	console.log('檔案上傳完畢。')
})
```

## 通知

透過 `notify` 來向伺服器廣播一個通知，且不要求回應與夾帶資料。

```js
// 伺服器會執行名為 `newUser` 的方法。
ws.notify('newUser')
```

# 狀態碼

在 Mego 中有內建這些狀態碼，但令 Mego 優於傳統 RESTful API 之處在於你可以自訂你自己的狀態碼。狀態碼是由數字奠定，從 `0` 到 `50` 都是成功，`51` 到 `100` 則是錯誤，如果你想要自訂自己的狀態碼，請從 `101` 開始。

| 狀態碼                   | 說明                                                                      |
|-------------------------|--------------------------------------------------------------------------|
| StatusOK                | 任何事情都很 Okay。                                                        |
| StatusProcessing        | 已成功傳送請求，但現在還不會完成，將會在背景中進行。                              |
| StatusNoChanges         | 這項請求沒有改變什麼事情，例如刪除了早就不存在的事物。                            |
| StatusError             | 內部錯誤發生，可能是非預期的錯誤。                                             |
| StatusFull              | 請求因為已滿而被拒絕，例如：好友清單已滿無法新增、聊天室人數已滿無法加入。            |
| StatusExists            | 已經有相同的事物存在而被拒絕。                                                |
| StatusInvalid           | 格式無效而被拒絕，通常是不符合表單驗證要求。                                     |
| StatusNotFound          | 找不到指定資源。                                                            |
| StatusNotAuthorized     | 請求被拒。沒有驗證的身份，需要登入以進行驗證。                                   |
| StatusNoPermission      | 已驗證身份，但沒有權限提出此請求因而被拒。                                      |
| StatusUnimplemented     | 此功能尚未實作完成，如果呼叫了一個不存在的函式即會回傳此狀態。                     |
| StatusTooManyRequests   | 使用者在短時間內發送太多的請求，請稍後重試。                                    |
| StatusResourceExhausted | 預定給使用者的配額已經被消耗殆盡，例如可用次數、分配容量。                        |
| StatusBusy              | 伺服器正處於忙碌狀態，請稍候重試。                                            |

# API 文件撰寫風格

由於不受傳統 RESTful 風格的拘束，你可以很輕易地就透過 Markdown 格式規劃出 Mego 的 API 文件。我們有現成的[使用者 API 範例](./examples/doc/User%20API.md)可供參考。

```markdown
# 會員系統

與會員、使用者有關的 API。

## 建立會員 [createUser]

建立一個新的使用者並在之後可透過註冊時所使用的帳號與密碼進行登入。

+ 請求

    + 欄位
        + username (string!) - 帳號。
        + password (string!) - 密碼。
        + email (string!) - 電子郵件地址。
        + gender (int!)

            使用者的性別採用編號以節省資料空間並在未來有更大的彈性進行變動。0 為男性、1 為女性。

    + 內容

            {
                "username": "YamiOdymel",
                "password": "yami123456",
                "email": "yamiodymel@yami.io",
                "gender": 1
            }
```