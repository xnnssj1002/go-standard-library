## 数据结构
我们先看一下Group的数据结构：
``` go
type Group struct {
    cancel func() // 这个存的是context的cancel方法
    wg sync.WaitGroup // 封装sync.WaitGroup
    errOnce sync.Once // 保证只接受一次错误
    err     error // 保存第一个返回的错误
}
```
## 方法集
``` go
func WithContext(ctx context.Context) (*Group, context.Context)
func (g *Group) Go(f func() error)
func (g *Group) Wait() error
```
### 1、WithContext方法
``` go
func WithContext(ctx context.Context) (*Group, context.Context) {
    ctx, cancel := context.WithCancel(ctx)
    return &Group{cancel: cancel}, ctx
}
```
> 这个方法只有两步：
> * 使用context的WithCancel()方法创建一个可取消的Context
> * 创建cancel()方法赋值给Group对象
### 2、Go方法
``` go
func (g *Group) Go(f func() error) {
    g.wg.Add(1)
    
    go func() {
        defer g.wg.Done()
        if err := f(); err != nil {
            g.errOnce.Do(func() {
                g.err = err
                if g.cancel != nil {
                    g.cancel()
                }
            })
        }
    }()
}
``` 
> Go方法中运行步骤如下：
> * 执行Add()方法增加一个计数器
> * 开启一个协程，运行我们传入的函数f，使用waitGroup的Done()方法控制是否结束
> * 如果有一个函数f运行出错了，我们把它保存起来，如果有cancel()方法，则执行cancel()取消其他goroutine

这里大家应该会好奇为什么使用errOnce，也就是sync.Once，这里的目的就是保证获取到第一个出错的信息，避免被后面的Goroutine的错误覆盖。

### 3、wait方法
``` go
func (g *Group) Wait() error {
    g.wg.Wait()
    if g.cancel != nil {
        g.cancel()
    }
    return g.err
}
```
> 总结一下wait方法的执行逻辑：
> * 调用waitGroup的Wait()等待一组Goroutine的运行结束
> * 这里为了保证代码的健壮性，如果前面赋值了cancel，要执行cancel()方法
> * 返回错误信息，如果有goroutine出现了错误才会有值

## 小结
> errGroup包，总共就1个结构体和3个方法，理解起来还是比较简单的，针对上面的知识点我们做一个小结：
> * 我们可以使用withContext方法创建一个可取消的Group，也可以直接使用一个零值的Group或new一个Group，不过直接使用零值的Group和new出来的Group出现错误之后就不能取消其他Goroutine了。
> * 如果多个Goroutine出现错误，我们只会获取到第一个出错的Goroutine的错误信息，晚于第一个出错的Goroutine的错误信息将不会被感知到。
> * errGroup中没有做panic处理，我们在Go方法中传入func() error方法时要保证程序的健壮性