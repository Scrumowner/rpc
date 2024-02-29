package ratelimit

import (
	"go.uber.org/ratelimit"
)

type RateLimiter struct {
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}
func (r *RateLimiter) Check(req chan string) {
	defer close(req)
	rl := ratelimit.New(5)
	rl.Take()
	for {
		select {
		case <-req:
			go func() {

			}()
		}

	}

}
