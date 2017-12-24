package main

import mego "github.com/TeaMeow/Mego"

func main() {
	// 初始化 Mego 伺服器。
	e := mego.Default()

	// 註冊一個會加總兩個參數的 `Sum` 方法。
	e.Register("Sum", func(c *mego.Context) {
		c.Respond(c.Param(0).GetInt() + c.Param(1).GetInt())
	})

	// 啟動伺服器。
	e.Run()
}
