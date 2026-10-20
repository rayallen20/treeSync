package model

// Book 书籍表模型
type Book struct {
	ID            int64  `gorm:"column:id;primaryKey"`  // ID 主键雪花ID
	PlanetID      int64  `gorm:"column:planet_id"`      // PlanetID 书籍所属星球ID
	Title         string `gorm:"column:title"`          // Title 书籍标题
	Sort          int    `gorm:"column:sort"`           // 排序字段
	DisplayStatus string `gorm:"column:display_status"` // DisplayStatus 展示状态
}
