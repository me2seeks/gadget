package ratelimiting_test

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestT(t *testing.T) {
	r := rate.NewLimiter(5, 10) //r:每秒生成令牌数  b:令牌容量
	for i := 0; i < 100; i++ {
		if r.Allow() {
			fmt.Println("allow", i)
		} else {
			fmt.Println("not allow")
			time.Sleep(time.Second)
		}
	}
}
