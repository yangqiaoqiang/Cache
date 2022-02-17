package lru

import "container/list"

// Cache 结构体包含字典和双向链表,并发访问不安全。
type Cache struct{
	maxBytes int64 //允许使用的最大内存
	nbytes int64 //当前使用的内存
	ll *list.List //双向链表
	cache map[string]*list.Element //字典*list.Element 当前key对应的结点指针
	OnEvicted func(key string, value Value)//某条记录删除的CallBack func

}
//返回 值所占用的内存大小
type Value interface {
	Len() int
}
//双向链表的数据类型
type entry struct {
	key string
	value Value
}

//实例化 Cache
func New(maxBytes int64,onEvicted func(string,Value))*Cache{
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}
//查找
func (c *Cache)Get(key string)(value Value,ok bool){
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv :=ele.Value.(*entry)
		return kv.value,true
	}
	return
}
//删除
func (c *Cache)RemoveOldest(){
	ele :=c.ll.Back()
	if ele!=nil{
		c.ll.Remove(ele)
		kv:=ele.Value.(*entry)
		delete(c.cache,kv.key)
		c.nbytes -=int64(len(kv.key))+int64(kv.value.Len())
		if c.OnEvicted!=nil{
			c.OnEvicted(kv.key,kv.value)
		}
	}
}
//新增/修改
func (c *Cache)Add(key string ,value Value){
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv :=ele.Value.(*entry)
		c.nbytes +=int64(value.Len())-int64(kv.value.Len())
		kv.value=value
	}else{
		ele:=c.ll.PushFront(&entry{key,value})
		c.cache[key]=ele
		c.nbytes+=int64(len(key))+int64(value.Len())

	}
	for c.maxBytes!=0&&c.maxBytes<c.nbytes{
		c.RemoveOldest()
	}
}

// Len 调用 list.List.len 记录条数
func (c *Cache) Len() int {
	return c.ll.Len()
}