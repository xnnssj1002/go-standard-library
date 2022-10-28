package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Result string
type Search func(ctx context.Context, query string) (Result, error)

func fakeSearch(kind string) Search {
	return func(_ context.Context, query string) (Result, error) {
		return Result(fmt.Sprintf("%s result for %q", kind, query)), nil
	}
}

func main() {
	Google := func(ctx context.Context, query string) ([]Result, error) {
		g, ctx := errgroup.WithContext(ctx)

		searches := []Search{Web, Image, Video}
		results := make([]Result, len(searches))
		for i, search := range searches {
			i, search := i, search // https://golang.org/doc/faq#closures_and_goroutines
			g.Go(func() error {
				result, err := search(ctx, query)
				if err == nil {
					results[i] = result
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return nil, err
		}
		return results, nil
	}

	results, err := Google(context.Background(), "golang")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}

}

/*
注意：fakeSearch 方法中没有做panic处理，我们在Go方法中传入func() error方法时要保证程序的健壮性，进行recover()接收
*/

/*
上面这个例子来自官方文档，代码量有点多，但是核心主要是在Google这个闭包中，
首先我们使用 errgroup.WithContext 创建一个errGroup对象和ctx对象，
然后我们直接调用errGroup对象的Go方法就可以启动一个协程了，Go方法中已经封装了waitGroup的控制操作，不需要我们手动添加了，
最后我们调用Wait方法，其实就是调用了waitGroup方法。
这个包不仅减少了我们的代码量，而且还增加了错误处理，对于一些业务可以更好的进行并发处理。
*/
