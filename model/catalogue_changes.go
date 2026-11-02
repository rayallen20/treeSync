package model

// CatalogueChanges 本结构体用于记录目录树变更情况
type CatalogueChanges struct {
	AddCatalogues      []*Catalogue // AddCatalogues 新增的目录列表
	AddArticles        []*Article   // AddArticles 新增的文章列表
	DeleteCatalogueIDs []int64      // DeleteCatalogueIDs 删除的目录ID列表
	DeleteArticleIDs   []int64      // DeleteArticleIDs 删除的文章ID列表
}
