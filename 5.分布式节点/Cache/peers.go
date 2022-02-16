package Cache

//抽象两个接口 PeerPicker 根据传入的key选择节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 用于从对应的group查找缓存值。对应HTTP客户端
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
