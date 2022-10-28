package analysis

import (
	"errors"
	"fmt"
	"sync"
)

// 多协程错误处理 - 接收多个错误

func ReceiveMoreErr() {
	var wg sync.WaitGroup

	errCh := make(chan error)
	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 模拟多协程部分失败，这里当i为偶数表示错误，奇数表示成功
			if i&1 == 0 {
				errCh <- errors.New(fmt.Sprintf("groutine %d error", i))
				return
			}
			fmt.Printf("groutine %d success exit\n", i)
		}()
	}

	// 另外开启一个协程，用于关闭channel
	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		fmt.Println(err)
	}

}
