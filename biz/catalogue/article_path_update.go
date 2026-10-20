package catalogue

// ArticlePathUpdate 本结构体用于表示一条待更新的文章路径
type ArticlePathUpdate struct {
	ArticleID int64  // ArticleID 目标文章id (即article表的ID字段值)
	NewPath   string // NewPath 新的文章路径 (即article表的path字段值)
}
