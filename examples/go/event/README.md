# event

簡易的事件監聽範例。

## 範例

```bash
# 啟動 Mego 伺服器，且不斷廣播 `MyEvent` 事件至 `Channel1` 頻道。
go run ./server/main.go
# 啟動客戶端，監聽 `Channel1` 頻道中的 `MyEvent` 事件。
go run ./client/main.go
```

```bash
2017/12/25 04:17:14 接收到 MyEvent 事件。
2017/12/25 04:17:15 接收到 MyEvent 事件。
2017/12/25 04:17:16 接收到 MyEvent 事件。
```