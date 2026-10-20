package catalogue

// childKey 本结构体用于表示子节点的唯一标识
// 相同父节点下,只能存在唯一一个name+typ组合的子节点
type childKey struct {
	name string // name 子节点名称
	typ  string // typ 子节点类型
}
