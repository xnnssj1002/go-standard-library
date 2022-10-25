package container_study

import (
	"container/list"
	"sync"
)

// List 链表
type List struct {
	*list.List
	mu sync.Mutex
}

// New 新建链表
func New() *List {
	return &List{List: list.New()}
}

// PushBack 向链表尾部插入值
func (l *List) PushBack(v interface{}) {
	// 加锁
	l.mu.Lock()
	defer l.mu.Unlock()
	l.List.PushBack(v)
}
