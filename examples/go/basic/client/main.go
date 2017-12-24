package main

import (
	"fmt"

	"github.com/TeaMeow/Mego/client"
)

func main() {
	// 建立新的 Mego 客戶端。
	client := client.NewClient("localhost:5000")

	var v int
	// 呼叫遠端的 `Sum` 方法。
	err := client.Call("Sum").
		// 並且傳入兩個數字參數要求加總。
		Send([]int{5, 3}).
		// 將回傳的資料映射到本地的變數 `v`。
		Bind(&v).
		// 發送資料。
		End()
	if err != nil {
		panic(err)
	}

	fmt.Printf("從遠端 Mego 伺服器計算的 5 + 3 結果為：%d", v)
}
