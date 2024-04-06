package writeaggregation

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestT(t *testing.T) {
	file, err := os.OpenFile("hello.txt", os.O_RDWR, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 初始化 Sync 聚合服务
	syncJob := NewSyncJob(func(interface{}) error {
		fmt.Printf("do sync...\n")
		time.Sleep(time.Second)
		return file.Sync()
	})

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 执行写操作 write ...
			fmt.Printf("write...\n")
			// 触发 sync 操作
			syncJob.Do(file)
		}()
	}
	wg.Wait()
}
