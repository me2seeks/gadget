package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

/*
- Name：熔断器的名字，用于区分不同的熔断器实例。
- MaxRequests：在半开状态时允许通过的最大请求数。熔断器从打开状态转换到半开状态时，会允许有限数量的请求通过以检测系统的健康状况。如果这些请求都成功了，熔断器会关闭，系统恢复正常状态。
- Interval：在打开状态时，熔断器尝试恢复的时间间隔。一旦过了这个时间，熔断器会转到半开状态，允许部分请求尝试执行。
- Timeout：熔断器开启状态的持续时间。在此期间所有尝试通过的请求都会立即被拒绝。一旦超时，熔断器会转换到半开状态。
- ReadyToTrip：一个函数，它定义了熔断器从关闭状态转换到开启状态的条件。gobreaker.Counts 包含失败、成功等请求的计数。
- OnStateChange：状态变化时的回调函数，当熔断器的状态发生变化（例如从关闭到开启）时会被调用。这可以用来记录日志或者进行其他操作。
- IsSuccessful：一个函数，用于判定一个操作是否成功。通常根据错误值来判断，如果没有错误，则认为操作成功。这个函数被用来更新成功和失败的计数，它们是决定是否熔断的关键因素。
*/
func main() {
	var logChan = make(chan string, 1) // 创建日志通道
	var changeChan = make(chan string, 1)

	settings := gobreaker.Settings{
		Name:        "Mock Breaker",
		ReadyToTrip: getReadyToTripFunc(logChan),
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			changeChan <- fmt.Sprintf("Circuit breaker '%s' state changed from %s to %sn", name, stateToString(from), stateToString(to))
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	for i := 0; i < 20; i++ {
		// 执行并通过断路器保护的操作
		_, err := cb.Execute(func() (interface{}, error) {
			// 随机返回错误以模拟失败请求
			if i%2 == 0 {
				return nil, errors.New("error")
			}
			return "success", nil
		})

		if err != nil {
			fmt.Println("Operation failed:", err.Error())
		} else {
			fmt.Println("Operation succeeded.")
		}

		printLog(logChan) // 检查并打印日志消息
		printChange(changeChan)

		time.Sleep(500 * time.Millisecond) // 等待一段时间再尝试下一个请求
	}
}

// getReadyToTripFunc 返回一个闭包，用于确定断路器何时触发
func getReadyToTripFunc(logChan chan string) func(gobreaker.Counts) bool {
	return func(counts gobreaker.Counts) bool {
		if counts.Requests < 3 {
			return false
		}
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		if failureRatio >= 0.6 {
			logChan <- "failure ratio too high"
			return true
		}
		return false
	}
}

// stateToString 将断路器状态转为字符串
func stateToString(state gobreaker.State) string {
	switch state {
	case gobreaker.StateClosed:
		return "Closed"
	case gobreaker.StateHalfOpen:
		return "Half-Open"
	case gobreaker.StateOpen:
		return "Open"
	}
	return "Unknown"
}

// printLog 检查日志通道是否有消息，并打印出来
func printLog(logChan chan string) {
	select {
	case log := <-logChan:
		fmt.Println("Received log:", log)
	default:
		// 无消息时不做任何事
	}
}

func printChange(changeChan chan string) {
	select {
	case log := <-changeChan:
		fmt.Println("Received log:", log)
	default:
		// 无消息时不做任何事
	}
}
