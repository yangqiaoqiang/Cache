# Cache 分布式缓存
Achieve Cache GO Like GroupCache、MemCache

## 1.LRU(Least Recently Used)

实现LRU淘汰算法两个核心数据结构

1.字典(map)，存储键(string)与值(list.Element链表节点)的关系。

2.双向链表实现的队列(list.List)

## 2.单机并发缓存

构建只读数据结构ByteView表示缓存值。用sync.Mutex封装LRU方法。

Group是最核心的数据结构，负责与用户的交互，并且控制缓存值存储和获取的流程。

Group.Get()实现缓存中存在直接获取，不存在则通过callback函数添加

```
                            是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值⑶
```

