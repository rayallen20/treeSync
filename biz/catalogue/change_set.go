package catalogue

import (
	"treeSync/model"
)

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
