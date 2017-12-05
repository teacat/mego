package mego

import "time"

// RateLimitConfig 是用在流量限制的選項建構體。
type RateLimitConfig struct {
	// Period 是計時的週期。
	Period int
	// Limit 表示單個階段在指定週期內僅能發送相同請求的數量。
	Limit int
}

// RateLimit 會回傳一個流量限制的中介軟體，可以限制單個階段在限時時間內能夠發出多少個請求。
func RateLimit(conf RateLimitConfig) HandlerFunc {
	// session 存放單個階段的資料。
	type session struct {
		// count 是此階段在週期內的請求數。
		count int
		// start 是此階段重新計數開始的時間，用以在之後比對週期。
		start time.Time
	}
	var store map[string]*session

	return func(c *Context) {
		// reset 會重設此階段的流量限制資料。
		reset := func() {
			store[c.Session.ID] = &session{
				count: 0,
				start: time.Now(),
			}
		}

		// 檢查存儲裡有沒有這個階段的資料，沒有則初始化一個。
		if _, ok := store[c.Session.ID]; !ok {
			reset()
		}

		// 取得階段資料裡的起始計時時間，並比對現在取得中間所消耗的時間。
		// 如果中間消耗的時間秒數已經大於或等於設定，那麼就可以重設了。
		duration := int(time.Since(store[c.Session.ID].start).Seconds())
		if duration >= conf.Period {
			reset()
		}

		// 如果請求次數大於限制就終止這個請求。
		if store[c.Session.ID].count >= conf.Limit {
			c.AbortWithStatus(StatusTooManyRequests)
			return
		}

		// 不然就遞加請求計數器。
		store[c.Session.ID].count++
		c.Next()
	}
}
