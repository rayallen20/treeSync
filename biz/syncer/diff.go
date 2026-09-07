package syncer

import "treeSync/model"

// ArticlePathUpdate 本结构体用于表示一条待更新的文章路径
type ArticlePathUpdate struct {
	ArticleID int64  // ArticleID 目标文章id (即article表的ID字段值)
	NewPath   string // NewPath 新的文章路径 (即article表的path字段值)
}

// ChangeSet 本结构体用于表示1次目录树同步(写入DB)操作的全部变更及其明细
type ChangeSet struct {
	InsertCatalogues   []model.Catalogue   // InsertCatalogues 待插入的目录行
	InsertArticles     []model.Article     // InsertArticles 待插入的文章行
	DeleteCatalogueIDs []int64             // DeleteCatalogueIDs 待删除的目录行ID列表(含整枝后代)
	DeleteArticleIDs   []int64             // DeleteArticleIDs 待删除的文章行ID列表
	UpdateArticlePaths []ArticlePathUpdate // UpdateArticlePaths 匹配成功但path变更的文章列表
	AddedLines         []string            // AddedLines 新增行明细
	RemovedLines       []string            // RemovedLines 删除行明细
}

// pair 本结构体表示BFS匹配队列元素 一对已经匹配的folder节点(新树和旧树)及其相对书籍根的路径(书籍根的 display 字段为空字符串)
type pair struct {
	old     *Node  // old 旧树的folder节点
	new     *Node  // new 新树的folder节点
	display string // display 相对书籍根的路径
}

// Reconcile 本函数用于对比新旧目录树,生成变更集(包含新增目录/文章,删除目录/文章,以及文章路径变更)
// oldRoot 可为 nil (表示首次构建)
func Reconcile(oldRoot, newRoot *Node, bookID int64, bookTitle string, nextID func() int64) *ChangeSet {
	return nil
}

// insertSubtree 本函数用于将给定的子树变更为待插入的行
func insertSubtree(node *Node, parentID int64, bookID int64, displayPrefix string, nextID func() int64, changeSet *ChangeSet) {

}

// deleteSubtree 本函数用于把旧树中给定的一整支子树收集进变更集的删除集(目录id与文章id分别收集),
// 并按 displayPrefix 拼出 RemovedLines(file 节点的文件名从 article.path 末段取回)
func deleteSubtree(node *Node, displayPrefix string, changeSet *ChangeSet) {}
