package syncer

import (
	"errors"
	"treeSync/biz/catalogue"

	"gorm.io/gorm"
)

// ErrBookNotFound 哨兵错误,表示指定的书籍不存在
var ErrBookNotFound = errors.New("book不存在")

// Params 本结构体用于表示目录树同步器的参数
type Params struct {
	BookID  int64  // BookID 书籍id (即book表的ID字段值)
	Path    string // Path 给定的路径(相对路径或绝对路径)
	RootDir string // RootDir 存储根目录 即 conf.Storage.RootDir
	DryRun  bool   // DryRun 是否为模拟运行,若为true则不执行写入DB操作,仅打印变更明细
}

// Report 本结构体用于表示目录树同步器的执行结果
type Report struct {
	BookID          int64    // BookID 书籍id (即book表的ID字段值)
	DryRun          bool     // DryRun 是否为模拟运行,若为true则不执行写入DB操作,仅打印变更明细
	ScanArticles    int      // ScanArticles 新树的file节点数量
	ScanFolders     int      // ScanFolders 新树的folder节点数量
	AddedLines      []string // AddedLines 新增行明细 即 ChangeSet.AddedLines
	AddedArticles   int      // AddedArticles 新增文章数量 即 len(ChangeSet.InsertArticles)
	RemovedLines    []string // RemovedLines 删除行明细 即 ChangeSet.RemovedLines
	RemovedArticles int      // RemovedArticles 删除文章数量 即 len(ChangeSet.DeleteArticleIDs)
	RefreshedPaths  []string // RefreshedPaths 需要刷新路径的文章列表 即 ChangeSet.UpdateArticlePaths 中的 NewPath 字段
}

// Run 本函数用于执行一次完整的同步过程
func Run(params Params) (*Report, error) {
	return nil, nil
}

// resolveBookDir 本函数用于解析给定的路径,返回书籍根目录和书籍根路径相对存储根路径的前缀
func resolveBookDir(rawPath, rootDir string) (bookDir string, pathPrefix string, err error) {
	return "", "", nil
}

// plan 本函数用于根据给定的书籍ID和新树,生成变更集.若lock为true则在执行过程中对bookID对应的行加锁,以防止并发修改
// Tips: dryRun 同样使用本函数,差别在于 query 传入的 gorm.DB 是 db.DB 还是 db.Transaction,若为 db.Transaction 则会在事务中执行,
// 若为 db.DB 则不会执行事务
func plan(query *gorm.DB, bookID int64, newRoot *catalogue.Catalogue, lock bool) (*catalogue.ChangeSet, error) {
	return nil, nil
}

// apply 本函数用于根据给定的变更集执行数据库写入操作,若执行成功则返回nil,否则返回错误
func apply(tx *gorm.DB, changeSet *catalogue.ChangeSet) error {
	return nil
}

// countTree 本函数用于统计新树的规模,返回新树的文章数量和目录数量
func countTree(root *catalogue.Catalogue) (articles int, folders int) {
	return 0, 0
}

// fillReport 本函数用于根据给定的变更集填充报告,将变更集中的新增/删除行和更新行明细写入报告中
func fillReport(report *Report, changeSet *catalogue.ChangeSet) {

}
