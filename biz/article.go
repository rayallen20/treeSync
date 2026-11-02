package biz

import "treeSync/model"

type Article struct {
	ID   int64  // ID 主键雪花ID
	Path string // 相对存储根的路径(含文件扩展名)
}

// fillByOrm 本方法用于根据给定的文章ORM对象填充文章业务结构体
func (a *Article) fillByOrm(orm *model.Article) {
	a.ID = orm.ID
	a.Path = orm.Path
}
