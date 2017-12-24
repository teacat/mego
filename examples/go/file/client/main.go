package main

import (
	"log"

	"github.com/TeaMeow/Mego/client"
)

func main() {
	// 建立新的 Mego 客戶端。
	client := client.NewClient("localhost:5000")

	var v string
	// 呼叫遠端的 `Upload` 方法。
	err := client.Call("Upload").
		// 將本範例底下的 `example.jpg` 上傳至遠端伺服器。
		SendFile("./example.jpg").
		// 將回傳的資料映射到本地的變數 `v`。
		Bind(&v).
		// 發送資料。
		End()
	if err != nil {
		panic(err)
	}

	log.Printf("已成功將檔案上傳至遠端伺服器的 %s 位置。", v)
}
