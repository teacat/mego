# basic

一個最簡單的請求範例。

## 範例

```bash
# 啟動生產者，這會向 NSQ 傳送一個 Protobuf 資料到 `msg` 主題。
go run ./producer/main.go
# 啟動消費者，這會向 NSQ 註冊一個基於 `msg` 主題的 `user` 頻道，
# 所以就能接收 `msg` 主題的訊息，接著透過 Protobuf 解碼。
go run ./consumer/main.go
```

```bash
2017/03/14 04:17:09 INF    2 [msg/user] querying nsqlookupd http://127.0.0.1:4161/lookup?topic=msg
2017/03/14 04:17:09 INF    2 [msg/user] (127.0.0.1:4150) connecting to nsqd
2017/03/14 04:17:14 接收到了一個 Protobuf 資料，而 Content 是：你好，世界！
```

## Protobuf

```proto
// Msg 呈現了一個訊息資料。
message Msg {
    string content = 1;
}
```