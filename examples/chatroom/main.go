package main

import mego "github.com/TeaMeow/Mego"

func main() {
	// 建立一個 Mego 引擎。
	e := mego.New()

	// Chat - 聊天
	//
	// 參數：A string
	// 說明：將客戶端傳入的 A 參數以 `Message` 事件發送給所有客戶端。
	e.Register("Chat", func(c *mego.Context) {
		v, _ := c.Params.GetString(0)
		e.Emit("Message", v)
	})

	// 在 :80 埠口上執行 Mego 引擎。
	e.Run()
}
