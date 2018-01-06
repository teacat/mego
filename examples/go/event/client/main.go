package main

import (
	"log"

	c "github.com/TeaMeow/Mego/client"
)

func main() {
	// 建立新的 Mego 客戶端。
	client := c.New("ws://localhost:5000")
	if err := client.Connect(); err != nil {
		panic(err)
	}

	// 向伺服端訂閱 `Channel1` 的 `MyEvent` 事件。
	err := client.Subscribe("MyEvent", "Channel1")
	if err != nil {
		panic(err)
	}

	// 處理接收到的 `MyEvent` 事件。
	client.On("MyEvent", func(e *c.Event) {
		log.Println("接收到 MyEvent 事件。")
	})
}
