package book

import (
	"errors"

	"gorm.io/gorm"
)

const (
	DisplayStatusVisible = "visible" // DisplayStatusVisible 书籍展示状态: 显示
	DisplayStatusHidden  = "hidden"  // DisplayStatusHidden 书籍展示状态: 隐藏
)

var ErrNotFound = errors.New("book不存在") // ErrNotFound 哨兵错误,表示指定的书籍不存在

// Book 书籍业务结构体
type Book struct {
	ID            int64  // ID 书籍ID
	Title         string // Title 书籍标题
	Sort          int    // Sort 排序字段
	DisplayStatus string // DisplayStatus 展示状态
}

// Get 本函数用于根据给定的书籍ID查询书籍信息,若不存在则返回 ErrNotFound
func Get(q *gorm.DB, id int64) (*Book, error) {
	return nil, nil
}

// Lock 本函数用于根据给定的书籍ID以 SELECT ... FOR UPDATE 的方式查询书籍信息,若不存在则返回 ErrNotFound
func Lock(tx *gorm.DB, id int64) (*Book, error) {
	return nil, nil
}
