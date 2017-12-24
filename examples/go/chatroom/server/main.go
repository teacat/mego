package main

import mego "github.com/TeaMeow/Mego"

func main() {
	// 初始化 Mego 伺服器。
	e := mego.Default()

	// 建立一個 `SetName` 方法，允許客戶端更改自己的暱稱並廣播給其他客戶端得知。
	e.Register("SetName", func(c *mego.Context) {
		// 將使用者的新暱稱保存在 Mego 的鍵值組中。
		c.Session.Set("Nickname", c.Param(0).GetString())
		// 廣播一個新成員的事件給大家知道。
		e.Emit("NewMember", "Global", c.Param(0).GetString())
	})

	// 建立一個 `CreateMessage` 方法，接收訊息並傳遞給其他客戶端。
	e.Register("CreateMessage", func(c *mego.Context) {
		// 將接收到的訊息與其發送者一同廣播給所有客戶端。
		e.Emit("NewMessage", "Global", mego.H{
			"Sender":  c.Session.GetString("Nickname"),
			"Message": c.Param(0).GetString(),
		})
	})

	// 啟動伺服器。
	e.Run()
}
