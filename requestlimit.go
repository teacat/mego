package mego

// RequestLimitConfig 是總請求限制的選項建構體。
type RequestLimitConfig struct {
	// Limit 是指定方法最大的同時連線數。
	Limit int
}

// RequestLimit 會回傳一個能夠限制單個方法最大同時請求數的中介軟體。
func RequestLimit(conf RequestLimitConfig) HandlerFunc {
	// 建立一個計數器。
	count := 0

	return func(c *Context) {
		// 如果大於限制就終止這個請求。
		if count > conf.Limit {
			c.AbortWithStatus(StatusBusy)
			return
		}

		// 不然就遞加計數器。
		count++
		c.Next()

		// 執行完畢後計數器遞減。
		count--
	}
}
