package Cache

//只读数据结构 ByteView 用来表示缓存值
type ByteView struct{
	b []byte
}
//实现Value接口
func (v ByteView)Len()int{
	return len(v.b)
}

// ByteSlice 方法返回一个拷贝,防止缓存值被修改
func (v ByteView)ByteSlice()[]byte{
	return cloneBytes(v.b)
}
func cloneBytes(b []byte)[]byte{
	c:=make([]byte,len(b))
	copy(c,b)
	return c
}
func (v ByteView)String()string{
	return string(v.b)
}
