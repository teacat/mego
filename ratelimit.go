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
		// interval 是剩餘的週期秒數。
		interval int
	}
	var store map[string]*session

	// 透過 Goroutine 開啟新執行緒啟動一個計時器。
	go func() {
		// 不斷地重複。
		for {
			// 每一秒執行一次。
			<-time.After(1 * time.Second)

			// 遍歷整個階段流量限制存儲。
			for _, s := range store {
				// 如果時間到了，就重設此客戶端的計數和剩餘週期時間。
				if s.interval <= 0 {
					s.count = 0
					s.interval = conf.Period
				}

				// 繼續倒數時間。
				s.interval--
			}
		}
	}()

	return func(c *Context) {
		// 檢查存儲裡有沒有這個階段的資料，沒有則初始化一個。
		if _, ok := store[c.Session.ID]; !ok {
			store[c.Session.ID] = &session{
				count:    0,
				interval: conf.Period,
			}
		}

		// 如果請求次數大於限制就終止這個請求。
		if store[c.Session.ID].count > conf.Limit {
			c.AbortWithStatus(StatusTooManyRequests)
			return
		}

		// 不然就遞加請求計數器。
		store[c.Session.ID].count++
		c.Next()
	}
}
