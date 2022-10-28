package analysis

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"runtime"
	"time"
)

// forever 持续打印数字,直到ctx结束
func forever(ctx context.Context, i int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("forever stop now")
			return
		default:
		}
		time.Sleep(time.Second)
		fmt.Printf("goroutine %d is running\n", i)
		runtime.Gosched()
	}
}

// delayError 5秒后报错
func delayError() error {
	time.Sleep(time.Second * 5)
	fmt.Println("delayError return")
	return errors.New("should stop now")
}

func ReceiveFirstErr() {
	g, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < 2; i++ {
		i := i
		if i == 0 {
			g.Go(delayError)
			continue
		}
		g.Go(func() error {
			forever(ctx, i)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
