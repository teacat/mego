# chunk

以區塊方式上傳檔案的範例。

## 範例

```bash
# 啟動 Mego 伺服器。
go run ./server/main.go
# 編譯並啟動客戶端，以區塊方式上傳資料夾內的 `example.jpg` 檔案並取得上傳後的暫存路徑。
go build ./client && ./client/client
```

```bash
2017/12/25 04:17:14 已成功將檔案上傳至遠端伺服器的 /tmp/3JbeKx 位置。
```