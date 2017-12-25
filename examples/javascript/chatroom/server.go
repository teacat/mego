package main

import mego "github.com/TeaMeow/Mego"

func main() {
	// 初始化 Mego 伺服器。
	e := mego.Default()

	// 建立一個 `CreateMessage` 方法，接收訊息並傳遞給其他客戶端。
	e.Register("CreateMessage", func(c *mego.Context) {
		// 將接收到的訊息與其發送者一同廣播給所有客戶端。
		e.Emit("NewMessage", "Global", c.Param(0).GetString())
	})

	// 啟動伺服器。
	e.Run()
}
