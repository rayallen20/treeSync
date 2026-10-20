package catalogue

// pair 本结构体表示BFS匹配队列元素: 一对已经匹配的folder节点(新树和旧树)及其相对书籍根的路径(书籍根的 display 字段为空字符串)
type pair struct {
	old     *Catalogue // old 旧树的folder节点
	new     *Catalogue // new 新树的folder节点
	display string     // display 相对书籍根的路径
}
