// 建立新的 Mego 客戶端。
ws = new MegoClient("localhost:5000")

// 呼叫遠端的 `Sum` 方法。
ws.call("sum")
    // 並且傳入兩個數字參數要求加總。
    .send([5, 3])
    // 發送資料。
    .end()
    // 處理接收到的回應。
    .then((result) => {
        alert(`從遠端 Mego 伺服器計算的 5 + 3 結果為：${result}`)
    })