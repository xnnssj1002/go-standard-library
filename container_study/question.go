package container_study

import (
	"container/list"
	"fmt"
)

// Question 这里是一个问题
// 希望输出0，但是输出的是-1
func Question() {
	l := list.List{}
	e := l.PushBack(10)

	l = list.List{}
	l.Remove(e)
	fmt.Println("list len: ", l.Len())
}

/*
图片看考 node.md 文件 链表重赋值
当我们执行完 l.PushBack(10)  这段代码直接，就如左图所示了，节点 e 指向了链表 l，链表 l 的值的长度是 1
但是当我们执行完 l = list.List{}  我们重置了链表 l 的值，并没有修改他的地址，所以 e 还是指向的链表 l，当执行 Remove  方法的时候，e 的链表和当前执行的 l 的地址判断是可以对应的上的，但是实际上链表的值已经发生了变化，链表的长度已经不为 1 的，并且新的链表的根节点也没有指向 e
所以最后得到的长度才是 -1
最后细心的同学应该已经发现，这里其实还可能存在内存泄漏的问题，元素 e -> root(0x2) 其实已经没有用到了
*/
