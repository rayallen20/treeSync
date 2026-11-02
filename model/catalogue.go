package model

import (
	"errors"
	"fmt"
	"treeSync/db"
)

// Catalogue 目录表模型
type Catalogue struct {
	ID            int64    `gorm:"column:id;primaryKey"`                 // ID 主键雪花ID
	BookID        int64    `gorm:"column:book_id"`                       // BookID 目录所属书籍ID
	ParentID      int64    `gorm:"column:parent_id"`                     // ParentID 目录的父节点ID(根节点此字段值为0)
	Sort          *int     `gorm:"column:sort"`                          // Sort 排序字段(未排序状态下本字段值为nil)
	Name          string   `gorm:"column:name"`                          // 章节名称(若节点为根节点,则该节点的名称即为书籍名称,即book.title)
	Intro         string   `gorm:"column:intro"`                         // Intro 章节介绍
	Type          string   `gorm:"column:type"`                          // Type 章节类型 folder: 文件夹; file: 文件
	DisplayStatus string   `gorm:"column:display_status"`                // DisplayStatus 展示状态
	Article       *Article `gorm:"foreignKey:CatalogueID;references:ID"` // Article 章节对应的文章信息(若节点为非叶节点,则该字段为nil)
}

// FindByBookID 本方法用于根据给定的书籍ID查询该书籍下的所有目录
func (c *Catalogue) FindByBookID(bookID int64) ([]*Catalogue, error) {
	var catalogues []*Catalogue
	err := db.DB.
		Where("book_id = ?", bookID).
		Preload("Article").
		Find(&catalogues).
		Error
	if err != nil {
		return nil, err
	}

	return catalogues, nil
}

// CreateCatalogueTree 本方法用于批量创建目录树(包括目录和文章)
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

// ApplyCatalogueChanges 本方法用于根据给定的目录变更情况应用到数据库中
func (c *Catalogue) ApplyCatalogueChanges(changes *CatalogueChanges) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return errors.New("开启事务失败: " + tx.Error.Error())
	}

	commited := false
	defer func() {
		if !commited {
			tx.Rollback()
		}
	}()

	// 先删除文章,再删除目录
	for _, articleID := range changes.DeleteArticleIDs {
		err := tx.Delete(&Article{}, articleID).Error
		if err != nil {
			return errors.New(fmt.Sprintf("删除文章 %d 失败: %s", articleID, err.Error()))
		}
	}

	for _, catalogueID := range changes.DeleteCatalogueIDs {
		err := tx.Delete(&Catalogue{}, catalogueID).Error
		if err != nil {
			return errors.New(fmt.Sprintf("删除目录 %d 失败: %s", catalogueID, err.Error()))
		}
	}

	// 写入新目录节点和文章节点(分成2步写入,写入目录节点时不操作文章节点)
	for _, catalogue := range changes.AddCatalogues {
		err := tx.Omit("Article").Create(catalogue).Error
		if err != nil {
			return errors.New(fmt.Sprintf("创建目录 %d 失败: %s", catalogue.ID, err.Error()))
		}
	}

	for _, article := range changes.AddArticles {
		err := tx.Create(article).Error
		if err != nil {
			return errors.New(fmt.Sprintf("创建文章 %d 失败: %s", article.ID, err.Error()))
		}
	}

	err := tx.Commit().Error
	if err != nil {
		return errors.New("提交事务失败: " + err.Error())
	}

	commited = true
	return nil
}
