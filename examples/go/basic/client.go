package main

import (
	"fmt"

	"github.com/TeaMeow/Mego/client"
)

func main() {
	// 初始化一個客戶端，並且嘗試連線。
	c := client.New()
	if err := c.Connect(); err != nil {
		panic(err)
	}

	// 呼叫遠端 `Sum` 函式並且傳入兩個數字參數。
	result, err := c.Call("Sum", []int{3, 4})
	if err != nil {
		panic(err)
	}

	fmt.Printf("計算結果是：%d", result)
}
