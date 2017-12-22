// 初始化一個 Mego 連線。
ws = new MegoClient('ws://localhost')

// 監聽 WebSocket 開始運作時的事件。
ws.on('open', async () => {
    // 訂閱 `calculating` 事件。
    ws.subscribe('calculating')

    // 當 `calculating` 事件發生時就執行指定函式。
    ws.on('calculating', () => {
        console.log('正在進行計算。')
    })

    // 呼叫遠端 `sum` 進行 1 跟 2 的加總計算。
    response = await ws.call('sum', [1, 2])

    // 當完成時顯示加總後的答案。
    console.log('計算完畢，答案是：' + response.result)

    // 取消訂閱 `calculating` 事件。
    ws.unsubscribe('calculating')

    // 關閉連線。
    ws.close()
})