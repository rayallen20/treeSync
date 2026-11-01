package biz

import (
	"errors"
	"strings"
	"treeSync/conf"
	"treeSync/model"
	"treeSync/reporter"

	"gorm.io/gorm"
)

// Book 书籍业务结构体
type Book struct {
	ID            int64      // ID 主键雪花ID
	Title         string     // Title 书籍标题
	Sort          int        // 排序字段
	DisplayStatus string     // DisplayStatus 展示状态
	Catalogue     *Catalogue // Catalogue 书籍对应目录
}

// Build 本方法用于根据给定的书籍ID和目录路径构建书籍目录
func (b *Book) Build(id int64, path string) (*reporter.Report, error) {
	confInstance := conf.Get()
	isInLocation := strings.HasPrefix(path, confInstance.Storage.RootDir)
	if !isInLocation {
		return nil, errors.New("给定的路径不在存储根目录下")
	}

	orm := &model.Book{}
	err := orm.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("不存在给定ID的书籍")
		}

		return nil, errors.New("获取书籍信息失败: " + err.Error())
	}

	b.fillByORM(orm)
	report, err := b.Catalogue.Build(b.ID, b.Title, path)
	if err != nil {
		return nil, errors.New("构建书籍目录失败: " + err.Error())
	}

	return report, nil
}

// fillByORM 本方法用于根据给定的书籍ORM对象填充书籍业务结构体
func (b *Book) fillByORM(orm *model.Book) {
	b.ID = orm.ID
	b.Title = orm.Title
	b.Sort = orm.Sort
	b.DisplayStatus = orm.DisplayStatus
	b.Catalogue = &Catalogue{}
}
