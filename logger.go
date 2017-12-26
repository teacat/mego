package mego

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// Logger 會回傳一個紀錄中介軟體，能夠簡略地紀錄每個請求的 IP 和日期與耗費時間。
func Logger() HandlerFunc {
	return LoggerWithWriter(DefaultWriter)
}

// LoggerWithWriter 可以傳入自訂的 `io.Writer` 決定紀錄的輸出位置（如：`Stdout` 或檔案）。
func LoggerWithWriter(out io.Writer) HandlerFunc {
	return func(c *Context) {
		// 從接收到請求的時候就開始記錄時間。
		start := time.Now()

		// 繼續請求的執行。
		c.Next()

		// 取得完成的時間。
		end := time.Now()
		// 取得本次請求的總費時。
		latency := time.Since(start)
		// 取得本次呼叫的方法名稱。
		method := c.Method.Name
		// 取得產生此請求的客戶端 IP 位置。
		ip := c.ClientIP()
		// 取得此請求的回應錯誤狀態碼。
		code := 0
		// 依照狀態碼回傳顏色。
		codeColor := green
		if code != 0 {
			codeColor = red
		}

		// 取得本次請求的所有錯誤。
		var buffer bytes.Buffer
		for k, v := range c.Errors {
			fmt.Fprintf(&buffer, "Error #%02d: %s\n", (k + 1), v.Err)
			// 如果這個錯誤有中繼資料也一併輸出。
			if v.Meta != nil {
				fmt.Fprintf(&buffer, "     Meta: %v\n", v.Meta)
			}
		}
		comment := buffer.String()

		// 將訊息輸出到指定的寫入者。
		fmt.Fprintf(out, "[MEGO] %v |%s %3d %s| %13v | %s | %s\n%s",
			end.Format("2006/01/02 - 15:04:05"),
			codeColor, code, reset,
			latency,
			ip,
			method,
			comment,
		)
	}
}
