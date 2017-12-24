package main

import mego "github.com/TeaMeow/Mego"

func main() {
	// 初始化 Mego 伺服器。
	e := mego.Default()

	// 註冊一個會接收檔案的基本 `Upload` 方法。
	e.Register("Upload", func(c *mego.Context) {
		// 回傳接收到的檔案其上傳後的暫存路徑。
		c.Respond(c.MustGetFile().Path)
	})

	// 啟動伺服器。
	e.Run()
}
