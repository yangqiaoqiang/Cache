package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}
//测试get方法
func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}
//测试超过设定时是否触发超出结点删除
func TestRemoveOldest(t *testing.T) {
	k1,k2,k3:="key1","key2","key3"
	v1,v2,v3:="value1","value2","value3"
	cap:=len(k1+k2+v1+v2)
	//fmt.Println(int64(cap))
	lru:=New(int64(cap),nil)
	lru.Add(k1,String(v1))
	lru.Add(k2,String(v2))
	lru.Add(k3,String(v3))
	//fmt.Println(lru.Len())
	if _,ok:=lru.Get("key1");ok||lru.Len()!=2{
		t.Fatalf("Removeoldest key1 failed")
	}
}
//测试回调函数能否使用
//回调函数在 RemoveOldest()中，超过maxBytes时回调
//插入key1时，nbytes=10
//插入k2时，nbytes=14，删除key1，nbytes=4，append("key1")
//插入k3时，nbytes=8
//插入k4时，nbytes=12，删除k2，nbytes=8,append("k2")
func TestOnEvicted(t *testing.T){
	keys :=make([]string,0)
	callback:=func(key string,value Value){
		keys =append(keys,key)
	}
	lru :=New(int64(10),callback)
	lru.Add("key1",String("123456"))
	lru.Add("k2",String("k2"))
	lru.Add("k3",String("k3"))
	lru.Add("k4",String("k4"))
	expect :=[]string{"key1","k2"}
	if !reflect.DeepEqual(expect,keys){
		t.Fatalf("Call OnEvicted failed,expect keys equals to %s",expect)
	}

}
