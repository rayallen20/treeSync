package model

import (
	"errors"
	"fmt"
	"treeSync/db"
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

// FindByBookID 本方法用于根据给定的书籍ID查询该书籍下的所有目录
func (c *Catalogue) FindByBookID(bookID int64) ([]*Catalogue, error) {
	var catalogues []*Catalogue
	err := db.DB.Where("book_id = ?", bookID).Find(&catalogues).Error
	if err != nil {
		return nil, err
	}

	return catalogues, nil
}

func (c *Catalogue) CreateCatalogueTree(catalogueOrms []*Catalogue, articleOrms []*Article) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return errors.New("开启事务失败: " + tx.Error.Error())
	}

	for _, catalogueOrm := range catalogueOrms {
		err := tx.Create(catalogueOrm).Error
		if err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("创建目录 %d 失败: %s", catalogueOrm.ID, err.Error()))
		}
	}

	for _, articleOrm := range articleOrms {
		err := tx.Create(articleOrm).Error
		if err != nil {
			tx.Rollback()
			return errors.New(fmt.Sprintf("创建文章 %d 失败: %s", articleOrm.ID, err.Error()))
		}
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return errors.New("提交事务失败: " + err.Error())
	}

	return nil
}
