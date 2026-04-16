package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	ip         string
	canRequest int
	firstTime  time.Time
}

var ipRateLimiters = make(map[string]*rateLimiter)
var mutex = &sync.Mutex{}

func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		mutex.Lock()
		defer mutex.Unlock()
		rl, ok := ipRateLimiters[ip]
		if !ok {
			log.Println("ip :", ip)
			rl = &rateLimiter{ip, 10, time.Now()}
			ipRateLimiters[ip] = rl
		} else if time.Since(rl.firstTime).Minutes() > 3 {
			rl.firstTime = time.Now()
			rl.canRequest = 10
		} else if rl.canRequest <= 0 {
			http.Error(w, "max rate must be greater than zero", http.StatusTooManyRequests)
			return
		}
		rl.canRequest--

		next(w, r)
	}
}
