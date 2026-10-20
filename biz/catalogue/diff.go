package catalogue

// insertSubtree 本函数用于将以node为根的子树变更为待插入的行
func insertSubtree(node *Catalogue, parentID int64, bookID int64, displayPrefix string, nextID func() int64, changeSet *ChangeSet) {

}

// deleteSubtree 本函数用于在旧树中将以node为根的子树收集进变更集的删除集(目录id与文章id分别收集),
// 并按 displayPrefix 拼出 RemovedLines(file 节点的文件名从 article.path 末段取回)
func deleteSubtree(node *Catalogue, displayPrefix string, changeSet *ChangeSet) {}
