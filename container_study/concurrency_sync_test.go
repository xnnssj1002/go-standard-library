package container_study

import (
	"sync"
	"testing"
)

func TestList_PushBack(t *testing.T) {
	l := New()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			l.PushBack(i)
		}
		wg.Done()
	}()
	l.PushBack(11)
	wg.Wait()
}
