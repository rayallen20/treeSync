package model

const (
	// CatalogueTypeFolder 节点类型: 文件夹
	CatalogueTypeFolder = "folder"

	// CatalogueTypeFile 节点类型: 文件
	CatalogueTypeFile = "file"

	// CatalogueDisplayStatusVisible 章节展示状态: 显示
	CatalogueDisplayStatusVisible = "visible"

	// CatalogueDisplayStatusHidden 章节展示状态: 隐藏
	CatalogueDisplayStatusHidden = "hidden"
)

// Catalogue 目录表模型
type Catalogue struct {
	ID            int64  `gorm:"column:id;primaryKey"`  // ID 主键雪花ID
	BookID        int64  `gorm:"column:book_id"`        // BookID 目录所属书籍ID
	ParentID      int64  `gorm:"column:parent_id"`      // ParentID 目录的父节点ID(根节点此字段值为0)
	Sort          *int   `gorm:"column:sort"`           // Sort 排序字段(未排序状态下本字段值为nil)
	Name          string `gorm:"column:name"`           // 章节名称(若节点为根节点,则该节点的名称即为书籍名称,即book.title)
	Intro         string `gorm:"column:intro"`          // Intro 章节介绍
	Type          string `gorm:"column:type"`           // Type 章节类型 folder: 文件夹; file: 文件
	DisplayStatus string `gorm:"column:display_status"` // DisplayStatus 展示状态
}
