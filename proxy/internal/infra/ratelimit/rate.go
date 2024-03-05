package ratelimit

import (
	"go.uber.org/ratelimit"
	"sync"
	"time"
)

type RateLimiter struct {
	limiter ratelimit.Limiter
	users   map[string]time.Time
	mutex   sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiter: ratelimit.New(5),
		users:   make(map[string]time.Time),
	}
}

func (rl *RateLimiter) Allow(userID string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	if _, ok := rl.users[userID]; !ok {
		rl.users[userID] = time.Now()
		return true
	}
	take := rl.limiter.Take()
	between := take.Sub(rl.users[userID])
	if between < time.Millisecond*200 {
		rl.users[userID] = time.Now()
		return false
	}
	rl.users[userID] = time.Now()
	return true
}

func RateWorker(us chan string, res chan bool) {
	limiter := NewRateLimiter()
	for {
		select {
		case userid := <-us:
			resolve := limiter.Allow(userid)
			res <- resolve
		default:
			continue
		}

	}
}

//
//	for i := 0; i <= 1; i++ {
//		for j := 0; j <= 50; j++ {
//			if i == 0 {
//				userChan <- userID1
//				fmt.Println(<-resolveChan)
//			} else if i == 1 {
//				userChan <- userID
//				fmt.Println(<-resolveChan)
//			}
//		}
//
// }
