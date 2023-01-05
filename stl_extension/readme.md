- [stl\_extension 模块](#stl_extension-模块)
  - [LimitWaitGroup](#limitwaitgroup)
  - [OrderedMap](#orderedmap)

# stl_extension 模块

stl_extension 模块主要是对系统的 stl 的不足做补充。

## LimitWaitGroup
对原生 stl 的 WaitGroup 做补充，原生的 WaitGroup 没有办法限制协程的数量，比如如下处理，
```go
group := sync.WaitGroup{}
for i := 0; i < 10000; i++ {
  v := i
  group.Add(1)
  go func() {
    defer group.Done()
    fmt.Print(v)
  }()
}
group.Wait()
```
会在一瞬间创建大量的协程。而用 LimitWaitGroup，可以实现限制并发的协程数量。

```go
// 10 协程并发
group := stlextension.NewLimitWaitGroup(10)
for i := 0; i < 10000; i++ {
  v := i
  // 如果协程并发满了，会阻塞等待
  group.Add(1)
  go func() {
    defer group.Done()
    fmt.Print(v)
  }()
}
group.Wait()
```

## OrderedMap
golang stl 原生的 map 是基于 hash 的无序 map，OrderedMap 是对 hash map 的补充。支持按照 key 顺序遍历。

OrderedMap 使用 avl 树实现，是线程安全的。

和 c++ map 对比测试，1000 万的随机的增删查操作：
OrderedMap: 21806 ms, c++ map: 11592ms，效率比 c++ map 慢两倍。
next 遍历整个 map, OrderedMap: 676ms, c++ map: 171ms；
prev 遍历整个 map, OrderedMap: 663ms, c++ map: 198ms；
遍历的效率大概比 c++ map 的慢三倍。

c++ map 用红黑树实现，在随机数据上表现会比 avl 树要好。avl 树深度会比红黑树低一点，查询效率比较高。但是插入的效率没有红黑树高。

用法：
- 构建
```go
m := stlextension.OrderedMap{}
```

- 判断 map 是否为空
```go
m.Empty()
```

- 插入: 插入的 key 必须实现 stlextension.Key Less interface，也就是必须要实现比较方法，value 可以是任意 interface{}。如果插入的 key 和其他之前插入的旧的 key 比较的时候失败，插入也会失败。

```go
err := m.Insert(stlextension.IntKey(5), "BBB");
```

- 删除：插入的 key 必须实现 stlextension.Key Less interface，也就是必须要实现比较方法。
```go
err := m.Erase(stlextension.IntKey(5));
```

- 计数：计算 map 中的 key 的数量。
```go
num, err := m.Count(stlextension.IntKey(4))
```

- Find: 寻找 key 的 value。
```go
value, exixt, err := m.Find(stlextension.IntKey(5))
```

- Begin: Begin 获取 key 最小的元素的 key 和 value，配合 Next 进行迭代。
- Next: 按照 key 的顺序获取当前 key 的下一个 key 和 value。
```go
// test begin and next
for key, value := m.Begin(); key != nil; key, value, err = m.Next(key) {
  if err != nil {
    log.Fatalf("test next error: %v", err)
  }
  log.Printf("key: %v, value: %v", key, value)
}
```

- RBegin: RBegin 获取 key 最大的元素的 key 和 value，配合 Prev 进行方向迭代。
- Prev: 按照 key 的顺序获取当前 key 的前一个 key 和 value。
```go
// test rbegin and prev
for key, value := m.RBegin(); key != nil; key, value, err = m.Prev(key) {
  if err != nil {
    log.Fatalf("test prev error: %v", err)
  }
  log.Printf("key: %v, value: %v", key, value)
}
```

- Clear: 清空 map。
```go
m.Clear()
```

- Size: map 元素个数。
```go
s := m.Size()
```

- String: map 转换成 string 输出。
```go
log.Print(m.String())
```

example 参考：[TestOrderMap](https://github.com/memory-overflow/go-common-library/blob/main/stl_extension/stl_test.go#L13)。