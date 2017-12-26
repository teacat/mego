package mego

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"runtime"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery 會回傳一個回復中介軟體，
// 這能在請求發生 `panic` 錯誤時自動回復並顯示相關錯誤訊息供除錯，而不會導致整個伺服器強制關閉。
func Recovery() HandlerFunc {
	return RecoveryWithWriter(DefaultErrorWriter)
}

// RecoveryWithWriter 可以傳入自訂的 `io.Writer` 決定回復紀錄的輸出位置（如：`Stdout` 或檔案）。
func RecoveryWithWriter(out io.Writer) HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					stack := stack(3)
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					logger.Printf("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), err, stack, reset)
				}
				c.RespondWithError(StatusError, nil, ErrPanicRecovered)
			}
		}()
		c.Next()
	}
}

// stack 會回傳一個夭壽讚的已排序函式呼叫堆疊，且從指定的地方開始略過部分堆疊。
func stack(skip int) []byte {
	// 欲回傳的資料。
	buf := new(bytes.Buffer)

	// 當我們進行遍歷的時候我們會開啟檔案。這些變數則用來擺放已載入的檔案。
	var lines [][]byte
	var lastFile string

	// 從指定的地方略過。
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 輸出指定的資料。如果沒有該來源則不會輸出。
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source 回傳一個沒有多餘空白的指定行數文字。
func source(lines [][]byte, n int) []byte {
	// 在堆疊中的索引從 1 開始，但我們的陣列是 0。
	n--

	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function 會盡可能地回傳呼叫者的名稱。
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())

	// 名稱已經包含了到套件的路徑，不過檔案名稱已經有概括此名稱了所以這實際上是多餘的。
	// 加上還有個中間點很礙眼。
	//	runtime/debug.*T·ptrmethod
	// 而我們想要的是：
	//	*T.ptrmethod
	// 而且套件路徑還有可能包含點符號（例如：code.google.com/...）
	// 所以首先要移除路徑前輟。
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
