class MegoClient
    #
    @StatusOK        : 0
    @StatusProcessing: 1
    @StatusNoChanges : 2
    @StatusFileNext  : 10
    @StatusFileAbort : 11

    @StatusError            : 51
    @StatusFull             : 52
    @StatusExists           : 53
    @StatusInvalid          : 54
    @StatusNotFound         : 55
    @StatusNotAuthorized    : 56
    @StatusNoPermission     : 57
    @StatusUnimplemented    : 58
    @StatusTooManyRequests  : 59
    @StatusResourceExhausted: 60
    @StatusBusy             : 61
    @StatusFileRetry        : 70
    @StatusFileEmpty        : 71
    @StatusFileSize         : 72

    # 索引實際上會從 1 開始，因為對 Go 來說 0 是零值。
    id: 0

    #
    events: {}

    #
    tasks: []

    # 目前 Mego 客戶端的連線狀態。
    status: 'disconnected'

    #
    constructor: (@url) =>

    # 手動重新連線。
    reconnect: =>

    # 關閉並結束此連線。
    close: =>

    # 呼叫遠端方法並帶有指定參數或物件。
    call: (method, params) =>

    # 監聽特定事件。
    on: (event, handler) =>
        switch event
            when 'open', 'reopen', 'close', 'message', 'error'
            else

    # 像遠端伺服器表明欲訂閱指定事件。
    subscribe: (event) =>

    # 取消訂閱指定事件。
    unsubscribe: (event) =>

    # 上傳指定檔案。
    file: (method, file) =>

    # 對遠端伺服器廣播，簡單說就是無資料無回應的方法呼叫。
    notify: (method) =>