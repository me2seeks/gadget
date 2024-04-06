package read_aggregation_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"
)

func TestT(t *testing.T) {

	var g singleflight.Group
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			v, err, shared := g.Do("objectkey", func() (interface{}, error) {
				fmt.Printf("协程ID:%v 开始进行读操作\n", idx)
				time.Sleep(2 * time.Second)
				fmt.Printf("协程ID:%v 完成进行读操作\n", idx)
				return "objectvalue", nil
			})
			if err != nil {
				t.Logf("err:%v", err)
			}
			fmt.Printf("协程ID:%v 请求结果: %v, 是否共享结果: %v\n", idx, v, shared)
		}(i)
		time.Sleep(time.Second)
	}
	wg.Wait()
}
