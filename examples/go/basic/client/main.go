package main

import (
	"log"

	"github.com/TeaMeow/Mego/client"
)

func main() {
	// 建立新的 Mego 客戶端。
	client := client.New("ws://localhost:5000")
	if err := client.Connect(); err != nil {
		panic(err)
	}

	var v int
	// 呼叫遠端的 `Sum` 方法。
	err := client.Call("Sum").
		// 並且傳入兩個數字參數要求加總。
		Send([]int{5, 3}).
		// 將回傳的資料映射到本地的變數 `v`。
		EndStruct(&v)
	if err != nil {
		panic(err)
	}

	log.Printf("從遠端 Mego 伺服器計算的 5 + 3 結果為：%d", v)
}
