# Cache 分布式缓存
Achieve Cache GO Like GroupCache、MemCache

## 1.LRU(Least Recently Used)

实现LRU淘汰算法两个核心数据结构

1.字典(map)，存储键(string)与值(list.Element链表节点)的关系。

2.双向链表实现的队列(list.List)

