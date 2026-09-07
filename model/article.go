package model

// Article 文章表模型
type Article struct {
	ID          int64  `gorm:"column:id;primaryKey"` // ID 主键雪花ID
	CatalogueID int64  `gorm:"column:catalogue_id"`  // CatalogueID 文章对应的叶子节点ID
	Path        string `gorm:"column:path"`          // 相对存储根的路径(含文件扩展名)
}
