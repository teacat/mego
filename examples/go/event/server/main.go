package main

import (
	"time"

	mego "github.com/TeaMeow/Mego"
)

func main() {
	// 初始化 Mego 伺服器。
	e := mego.Default()

	// 透過 Goroutine 執行另一個函式。
	go func() {
		// 不斷地執行。
		for {
			// 每一秒。
			<-time.After(time.Second * 1)
			// 對 `Channel1` 頻道廣播 `MyEvent` 事件且不帶資料。
			e.Emit("MyEvent", "Channel1", nil)
		}
	}()

	// 啟動伺服器。
	e.Run()
}
