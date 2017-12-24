# file

簡易的檔案上傳範例。

## 範例

```bash
# 啟動 Mego 伺服器。
go run ./server/main.go
# 編譯並啟動客戶端，上傳資料夾內的 `example.jpg` 檔案至伺服器並取得上傳後的暫存路徑。
go build ./client && ./client/client
```

```bash
2017/12/25 04:17:14 已成功將檔案上傳至遠端伺服器的 /tmp/3JbeKx 位置。
```