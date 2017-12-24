package main

import (
	"log"

	"github.com/TeaMeow/Mego/client"
)

func main() {
	// 建立新的 Mego 客戶端。
	client := client.NewClient("localhost:5000")

	// 向伺服端訂閱 `Channel1` 的 `MyEvent` 事件。
	err := client.Subscribe("MyEvent", "Channel1")
	if err != nil {
		panic(err)
	}

	// 處理接收到的 `MyEvent` 事件。
	client.On("MyEvent", func(e *client.Event) {
		log.Println("接收到 MyEvent 事件。")
	})
}
