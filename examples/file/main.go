package main

import mego "github.com/TeaMeow/Mego"

func main() {
	// 建立一個 Mego 引擎。
	e := mego.New()

	// File - 檔案
	//
	// 輸出接收到的檔案名稱、大小、資料與路徑。
	e.Receive("File", func(c *mego.Context) {
		c.Respond(mego.StatusOK, mego.H{
			"name":      c.File.Name,
			"size":      c.File.Size,
			"extension": c.File.Extension,
			"path":      c.File.Path,
		})
	})

	// 在 :80 埠口上執行 Mego 引擎。
	e.Run()
}
