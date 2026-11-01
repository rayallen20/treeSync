package model

import "treeSync/db"

// Book 书籍表模型
type Book struct {
	ID            int64  `gorm:"column:id;primaryKey"`  // ID 主键雪花ID
	PlanetID      int64  `gorm:"column:planet_id"`      // PlanetID 书籍所属星球ID
	Title         string `gorm:"column:title"`          // Title 书籍标题
	Sort          int    `gorm:"column:sort"`           // 排序字段
	DisplayStatus string `gorm:"column:display_status"` // DisplayStatus 展示状态
}

// GetByID 本方法用于根据给定的id获取书籍ORM
func (b *Book) GetByID(id int64) error {
	return db.DB.First(b, id).Error
}
