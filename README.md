&nbsp;

<p align="center">
  <img src="https://user-images.githubusercontent.com/7308718/33408556-0f04a56c-d5b2-11e7-8a14-b556cf283b43.png"/>
</p>
<p align="center">
  <i>Meet Go.</i>
</p>

&nbsp;

# // WIP - 尚未完成 //
# // WIP - 尚未完成 //
# // WIP - 尚未完成 //
# // WIP - 尚未完成 //
# // WIP - 尚未完成 //
# // WIP - 尚未完成 //

# Mego

類似 [Gin](https://github.com/gin-gonic/gin) 用法，用以取代傳統 [RESTful](https://zh.wikipedia.org/zh-tw/REST) 溝通結構，並且以[類 RPC](https://zh.wikipedia.org/zh-tw/%E9%81%A0%E7%A8%8B%E9%81%8E%E7%A8%8B%E8%AA%BF%E7%94%A8) 的方式以 [MessagePack](http://msgpack.org/) 將資料進行壓縮並透過 [WebSocket](https://developer.mozilla.org/zh-TW/docs/WebSockets/WebSockets_reference/WebSocket) 傳遞使資料更加輕量。同時亦支援分塊檔案上傳。

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

因為 Mego 是透過 WebSocket 進行溝通，這令你可以主動式廣播訊息至客戶端（讓你更友善地達成通知系統）。Mego 亦支援檔案上傳，同時也能夠透過最簡單的方式妥善地以分塊方式上傳大型檔案。

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
	* [存取階段資料](#存取階段資料)
  * [處理請求與回應](#處理請求與回應)
    * [指定客戶端廣播事件](#指定客戶端廣播事件)
  * [中介軟體](#中介軟體)
    * [推遲執行與接續](#推遲執行與接續)
    * [終止請求](#終止請求)
    * [存放鍵值組](#存放鍵值組)
	* [內建中介軟體](#內建中介軟體)
	  * [自動回復](#自動回復)
	  * [請求紀錄](#請求紀錄)
	  * [配額限制](#配額限制)
	  * [同時請求數](#同時請求數)
  * [檔案上傳處理](#檔案上傳處理)
	* [單一檔案](#單一檔案)
	* [多個檔案](#多個檔案)
	* [區塊上傳](#區塊上傳)
  * [斷開客戶端](#斷開客戶端)
  * [複製並使用於 Goroutine](#複製並使用於-goroutine)
  * [錯誤處理](#錯誤處理)
  * [結束](#結束)
* [客戶端](#客戶端)
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
		c.Respond(mego.H{
			"username": "YamiOdymel",
		})
	})
	// 執行引擎在埠口 5000 上。
	e.Run()
}
```

###### JavaScript

```javascript
ws = new MegoClient('ws://localhost/')

ws.on('open', () => {
    ws.call('getUser').end().then(({result}) => {
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

一個事件可以有很多個頻道，類似一個「新訊息」事件可以發送到很多個「聊天室」。透過不同的頻道你可以避免不相關的事件廣播到其他人手中。

```go
func main() {
	e := mego.Default()

	// 廣播一個 `UpdateApp` 事件給所有監聽 `Channel1` 頻道的客戶端，觸發其事件監聽函式。
	e.Emit("UpdateApp", "Channel1", nil)

	// 廣播一個 `UpdateApp` 事件給 `Channel1` 頻道並帶有指定的資料供客戶端讀取。
	e.Emit("UpdateApp", "Channel1", mego.H{
		"version": "1.0.0",
	})

	e.Run()
}
```

### 多數廣播

透過 `Emit` 會廣播指定事件給指定頻道的所有連線的客戶端，如果你希望廣播事件給指定頻道中的某個客戶端時，你可以透過 `EmitMultiple` 並傳入欲接收指定事件的客戶端階段達成。

```go
func main() {
	e := mego.Default()

	// 對頻道 `Channel1` 內的 [0] 與 [1] 客戶端廣播 `UpdateApp` 事件，
	// 位於此頻道的其他客戶端不會接收到此事件。
	e.EmitMultiple("UpdateApp", "Channel1", nil, []*mego.Session{
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

	// 過濾所有客戶端，僅有 ID 為 `0k3LbEjL` 且在 `Channel1` 內才能夠接收到 UpdateApp 事件。
	e.EmitFilter("UpdateApp", "Channel1", nil, func(s *mego.Session) bool {
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

透過 `Param(index).Get()` 取得指定索引的參數內容，亦有 `GetInt`、`GetString` 等方便的資料型態快捷函式。

```go
func main() {
	e := mego.Default()

	e.Register("Sum", func(c *mego.Context) {
		// 以正整數資料型態取得傳入的兩個參數。
		a := c.Param(0).GetInt()
		b := c.Param(1).GetInt()

		// 將兩個轉換成正整數的參數進行相加取得結果。
		fmt.Printf("A + B = %d", a+b)
	})

	e.Run()
}
```

### 存取階段資料

你可以透過將資料存放至階段中，在下次相同客戶端呼叫時進行存取。

同時，你也不需要在每次傳送請求時夾帶一長串的 Token 或使用者身份驗證字串。

因為當 Mego 在初始化連線時就會將客戶端傳入的資料保存於階段中供你存取。

```go
func main() {
	e := mego.Default()

	e.Register("Sum", func(c *mego.Context) {
		// 取得客戶端初始化連線時所傳入的 `Token` 資料。
		c.Session.GetString("Token")

		// 在階段中保存名為 `Foo` 且值為 `Bar` 的內容，下次相同客戶端呼叫時可以進行存取。
		c.Session.Set("Foo", "Bar")
	})

	e.Run()
}
```

## 處理請求與回應

透過 `Respond` 正常回應一個客戶端的請求。而 `RespondWithError` 可以回傳一個錯誤發生的詳細資料用已告知客戶端發生了錯誤。

```go
func main() {
	e := mego.Default()

	e.Register("CreateUser", func(c *mego.Context) {
		// 針對此請求，回傳一個指定的狀態碼與資料。
		c.Respond(mego.H{
			"message": "成功建立使用者！",
		})

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

	e.Register("CreateUser", func(c *mego.Context) {
		// 廣播 `UpdateApp` 事件到此請求客戶端所監聽的 `Channel1` 中。
		c.Emit("UpdateApp", "Channel1", nil)

		// 廣播 `UpdateApp` 事件到此請求客戶端以外「其他所有人」所監聽的 `Channel1` 中。
		c.EmitOthers("UpdateApp", "Channel1", nil)
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
		if c.Session.ID != "jn3Dl74eX" {
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

如果連線有異常不希望繼續呼叫下一個中介軟體或處理函式，則需要透過 `Abort` 終止此請求。以 `AbortWithRespond` 可以回傳正常的資料與狀態碼。

而 `AbortWithError` 則會在終止的時候回傳狀態碼與錯誤的詳細資料。

```go
func main() {
	e := mego.Default()

	mw := func(c *mego.Context) {
		if c.Session.ID != "jn3Dl74eX" {
			// 直接終止此請求的繼續，並停止呼叫接下來的中介軟體與處理函式。
			c.Abort()

			// 終止此請求並且回傳資料。客戶端並不會知道是錯誤發生，仍會依一般回應處理。
			c.AbortWithRespond(mego.H{
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

Mego 能夠自動幫你處理透過 WebSocket 所上傳的檔案。

### 單一檔案

客戶端可以上傳多個且獨立的檔案。透過 `GetFile` 取得單一個檔案，需要傳遞一個檔案欄位的名稱至此函式。當沒有指定名稱時以第一個檔案為主，如果方法永遠只接收一個檔案，那麼就可以省略名稱。

如果檔案是必要的，那麼透過 `MustGetFile` 可以在沒有接收到檔案的情況下自動呼叫 `panic` 終止此請求。

```go
func main() {
	e := mego.Default()

	e.Register("UploadPhoto", func(c *mego.Context) {
		// 透過 `GetFile` 取得客戶端上傳的 `Photo` 檔案。
		if file, err := c.GetFile("Photo"); err == nil {
			c.Respond(mego.H{
				"filename":  file.Name,      // 檔案原始名稱。
				"size":      file.Size,      // 檔案大小（位元組）。
				"extension": file.Extension, // 檔案副檔名（無符號）。
				"path":      file.Path,      // 檔案本地路徑。
			})
		}
	})

	e.Run()
}
```

### 多個檔案

一個檔案欄位可以擺放多個檔案，這會使其成為檔案陣列。透過 `GetFiles` 以切片的方式取得指定檔案欄位的多個檔案。同樣地，當這個函式沒有傳入指定檔案欄位名稱時，以第一個檔案欄位為主。

這個函式會回傳一個切片方便進行遍歷。如果有個檔案欄位是必須的，透過 `MustGetFiles` 來在沒有接收到檔案的情況下自動呼叫 `panic` 中止請求。

```go
func main() {
	e := mego.Default()

	e.Register("UploadPhoto", func(c *mego.Context) {
		// 透過 `GetFiles` 取得客戶端上傳的多個檔案。
		if files, err := c.GetFiles("Photos"); err == nil {
			for _, file := range files {
				fmt.Println(file.Name) // 檔案原始名稱。
			}
		}
	})

	e.Run()
}
```

### 區塊上傳

如果客戶端中的檔案過大，區塊上傳便是最好的方法。你不需要在 Mego 特別設置，Mego 即會自動組合區塊。使用這個區塊上傳時方法只會接收到一個完整檔案與起初傳遞的資料。透過 `GetFile` 直接取得完整的檔案。

如果你希望能夠手動處理區塊，例如搭配 Amazon S3 的 Multi-part 將接收到的每個區塊都各自上傳至雲端時，請更改方法中的 `ChunkHandler`。

一個自訂的區塊處理函式內部結構看起來應該要像這樣。

```go
// 建立一個自訂的區塊處理函式。
myChunkHandler := func(c *mego.Context, raw *mego.RawFile, dest *mego.File) mego.ChunkStatus {
	// ... 區塊處理邏輯 ...
	// fmt.Println(raw.Binary)

	// 如果這個區塊是最後一塊的話。
	if raw.Last {
		// 將完成的資料保存到檔案建構體。
		dest.Name = "MyFile"
		dest.Size = 128691

		// 表示區塊處理完畢。
		return mego.ChunkDone
	}

	// 如果尚未完成則要求客戶端傳送下一塊。
	return mego.ChunkNext
}
```

照理來說 Mego 就能協助你處理區塊組合了，不過如果你像上述ㄧ樣建立了自訂的區塊組合函式，透過 `HandleChunk` 能夠覆蓋全域的預設區塊處理函式。透過修改方法中的 `ChunkHandler` 能夠獨立修改各個方法的區塊處理函式。

```go
func main() {
	e := mego.Default()

	// 更改全域的區塊處理函式。
	e.HandleChunk(myChunkHandler)

	// 註冊一個接受檔案的方法供客戶端上傳檔案至此。
	e.Register("UploadVideo", func(c *mego.Context) {
		var u User

		// 取得已從區塊組成的完整檔案。
		file := c.MustGetFile()
		fmt.Println(file.Name)

		// 你仍可以將檔案夾帶的資料映射到本地建構體。
		c.MustBind(&u)
		fmt.Println(u.Username)
	})

	// 你也能獨立更改 `UploadVideo` 方法的區塊處理函式。
	e.Method["UploadVideo"].ChunkHandler = myChunkHandler

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
			c.AbortWithError(mego.StatusError, nil, c.Errors[0].Err)
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

Mego 有附帶數個官方客戶端。

* [JavaScript](./client/js/README.md)
* [Golang](./client/README.md)

# 狀態碼

在 Mego 中有內建這些狀態碼，但令 Mego 優於傳統 RESTful API 之處在於你可以自訂你自己的狀態碼。狀態碼是由數字奠定，如果你想要自訂自己的狀態碼，請從 `101` 開始。

| 狀態碼                   | 說明                                                                      |
|-------------------------|--------------------------------------------------------------------------|
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

由於不受傳統 RESTful 風格的拘束，你可以很輕易地就透過 Markdown 格式規劃出 Mego 的 API 文件。我們有現成的[使用者 API 範例](./docs/Example%20API%20Docs.md)可供參考。

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