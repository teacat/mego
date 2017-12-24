# basic

一個最簡單的請求範例。

## 範例

```bash
# 啟動 Mego 伺服器。
go run ./server/main.go
# 啟動客戶端，傳送資料至伺服器並呼叫 `Sum` 方法且取得回應。
go run ./client/main.go
```

```bash
2017/12/25 04:17:14 從遠端 Mego 伺服器計算的 5 + 3 結果為：8
```