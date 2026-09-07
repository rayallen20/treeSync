package model

const (
	// BookDisplayStatusVisible 书籍展示状态: 显示
	BookDisplayStatusVisible = "visible"

	// BookDisplayStatusHidden 书籍展示状态: 隐藏
	BookDisplayStatusHidden = "hidden"
)

// Book 书籍表模型
type Book struct {
	ID            int64  `gorm:"column:id;primaryKey"`  // ID 主键雪花ID
	PlanetID      int64  `gorm:"column:planet_id"`      // PlanetID 书籍所属星球ID
	Title         string `gorm:"column:title"`          // Title 书籍标题
	Sort          int    `gorm:"column:sort"`           // 排序字段
	DisplayStatus string `gorm:"column:display_status"` // DisplayStatus 展示状态
}
