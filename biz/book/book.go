package book

import "treeSync/biz/catalogue"

const (
	DisplayStatusVisible = "visible" // DisplayStatusVisible 书籍展示状态: 显示
	DisplayStatusHidden  = "hidden"  // DisplayStatusHidden 书籍展示状态: 隐藏
)

// Book 书籍业务结构体
type Book struct {
	ID            int64                // ID 书籍ID
	Title         string               // Title 书籍标题
	Sort          int                  // Sort 排序字段
	DisplayStatus string               // DisplayStatus 展示状态
	Root          *catalogue.Catalogue // Root 书籍对应目录树的根节点
}
