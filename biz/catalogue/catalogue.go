package catalogue

import (
	"treeSync/biz/article"
)

const (
	TypeFolder           = "folder"  // TypeFolder 节点类型: 文件夹
	TypeFile             = "file"    // TypeFile 节点类型: 文件
	DisplayStatusVisible = "visible" // DisplayStatusVisible 章节展示状态: 显示
	DisplayStatusHidden  = "hidden"  // DisplayStatusHidden 章节展示状态: 隐藏
)

// Catalogue 目录业务结构体
type Catalogue struct {
	ID       int64                   // ID 目录ID
	Name     string                  // Name 章节名称(若节点为根节点,则该节点的名称即为书籍名称,即所属 book.Book 结构体的Title字段)
	Type     string                  // Type 章节类型
	RelPath  string                  // RelPath 章节相对路径(相对于书籍根目录的路径,仅新建树时用到该字段)
	Article  *article.Article        // Article 章节对应文章 (仅Type为file时,该字段才有值)
	Children []*Catalogue            // Children 子章节列表
	byKey    map[childKey]*Catalogue // byKey 子节点索引,和 Children 字段同步维护
}

// newCatalogue 本函数用于根据给定的节点名称和类型创建一个 Catalogue 实例
func newCatalogue(name, typ string) *Catalogue {
	return nil
}

// child 本方法用于根据给定的子章节名称和子章节类型在当前章节的子章节列表中查询子章节,若存在则返回该子章节,否则返回nil
func (c *Catalogue) child(name, typ string) *Catalogue {
	return nil
}

// addChild 本方法用于向当前章节下添加一个子章节
func (c *Catalogue) addChild(child *Catalogue) {

}

// Count 本方法用于统计当前章节及其所有子章节中文章和文件夹的数量,并返回统计结果
func (c *Catalogue) Count() (articles int, folders int) {
	return 0, 0
}

// Reconcile 本方法用于将当前目录树与旧目录树进行对比,并生成一个 ChangeSet 结构体,该结构体包含了所有需要进行的增删改操作,以便后续将变更写入数据库
// 其中: c 表示新目录树; old 表示旧目录树
func (c *Catalogue) Reconcile(old *Catalogue, bookID int64, bookTitle string, nextID func() int64) *ChangeSet {
	return nil
}
