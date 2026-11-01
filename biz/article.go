package biz

type Article struct {
	ID   int64  // ID 主键雪花ID
	Path string // 相对存储根的路径(含文件扩展名)
}
