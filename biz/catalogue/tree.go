package catalogue

import (
	"treeSync/model"
)

// FromPaths 本函数用于根据给定的相对路径列表和路径前缀创建一个 Catalogue 树结构,并返回该树结构的根节点
// Tips: 本函数用于创建新树
func FromPaths(relPaths []string, pathPrefix string) (*Catalogue, error) {
	return nil, nil
}

// FromRows 本函数用于根据给定的目录和文章数据创建一个 Catalogue 树结构,并返回该树结构的根节点
// Tips: 本函数用于从数据库中读取数据并创建树(创建旧树)
func FromRows(catalogues []model.Catalogue, articles []model.Article) (*Catalogue, error) {
	return nil, nil
}
