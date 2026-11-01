package biz

import (
	"errors"
	"fmt"
	"strings"
	"treeSync/enum"
	"treeSync/lib"
	"treeSync/model"
	"treeSync/reporter"
	"treeSync/scan"
	"treeSync/snowflake"
)

// Catalogue 章节业务结构体
type Catalogue struct {
	ID            int64                 // ID 主键雪花ID
	Sort          *int                  // Sort 排序字段(未排序状态下本字段值为nil)
	Name          string                // 章节名称(若节点为根节点,则该节点的名称即为书籍名称,即book.title)
	Intro         string                // Intro 章节介绍
	Type          string                // Type 章节类型 folder: 非叶节点; file: 叶子节点
	DisplayStatus string                // DisplayStatus 展示状态
	Children      []*Catalogue          // Children 子章节列表
	BookID        int64                 // BookID 章节所属书籍ID
	index         map[string]*Catalogue // index 子节点索引
	Article       *Article              // Article 章节对应的文章信息(若节点为非叶节点,则该字段为nil)
}

// Build 本方法用于根据给定的
// - 书籍ID
// - 书籍标题
// - 目录路径集合
// 构建书籍目录
func (c *Catalogue) Build(bookID int64, bookTitle, rootPath string) (*reporter.Report, error) {
	orm := &model.Catalogue{}
	orms, err := orm.FindByBookID(bookID)
	if err != nil {
		return nil, errors.New("获取书籍目录失败: " + err.Error())
	}

	if len(orms) == 0 {
		report, err := c.buildFromPath(bookID, bookTitle, rootPath)
		if err != nil {
			return nil, errors.New("构建书籍目录失败: " + err.Error())
		}

		return report, nil
	}

	return nil, nil
}

// buildFromPath 本方法用于根据给定的目录路径集合构建书籍目录
func (c *Catalogue) buildFromPath(bookID int64, bookTitle, rootPath string) (*reporter.Report, error) {
	c.ID = snowflake.Next()
	c.Name = bookTitle
	c.Type = enum.TypeFolder
	c.BookID = bookID
	c.DisplayStatus = enum.DisplayHidden

	subPaths, err := scan.Dir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, subPath := range subPaths {
		subPath = strings.Trim(subPath, "/")
		if subPath == "" {
			continue
		}

		parts := strings.Split(subPath, "/")
		err = c.insertCatalogue(parts, subPath, rootPath)
		if err != nil {
			return nil, err
		}
	}

	catalogueOrms, articleOrms := c.flatten()
	orm := &model.Catalogue{}
	err = orm.CreateCatalogueTree(catalogueOrms, articleOrms)
	if err != nil {
		return nil, errors.New("落盘书籍目录失败: " + err.Error())
	}

	report := &reporter.Report{
		CatalogueReport: &reporter.CatalogueReport{
			AddItemQuantity:    len(catalogueOrms),
			DeleteItemQuantity: 0,
			UpdateItemQuantity: 0,
		},
		ArticleReport: &reporter.ArticleReport{
			AddItemQuantity:    len(articleOrms),
			DeleteItemQuantity: 0,
			UpdateItemQuantity: 0,
		},
	}

	return report, nil
}

// insertCatalogue 本方法用于在当前节点下根据给定路径递归插入子节点
func (c *Catalogue) insertCatalogue(parts []string, articlePath, rootPath string) error {
	// 递归终止条件
	if len(parts) == 0 {
		return nil
	}

	name := parts[0]
	if name == "" {
		return errors.New("路径中包含空的章节名称")
	}

	if c.index == nil {
		c.index = make(map[string]*Catalogue)
	}

	// 判断节点类型
	catalogueType := enum.TypeFolder
	if len(parts) == 1 && lib.IsMarkdown(name) {
		catalogueType = enum.TypeFile
	}

	// 查找当前节点下是否存在同名子节点
	// 存在则在同名子节点下递归插入子节点的子节点
	// 不存在则在当前节点下创建新的子节点
	child, exists := c.index[name]
	if exists {
		return child.insertCatalogue(parts[1:], articlePath, rootPath)
	}

	// 创建新的子节点
	child = &Catalogue{
		ID:            snowflake.Next(),
		Name:          name,
		Intro:         name + "的章节介绍",
		Type:          catalogueType,
		DisplayStatus: enum.DisplayHidden,
		BookID:        c.BookID,
	}
	if catalogueType == enum.TypeFile {
		child.Article = &Article{
			ID:   snowflake.Next(),
			Path: rootPath + "/" + articlePath,
		}
	}

	c.Children = append(c.Children, child)
	c.index[name] = child

	// 递归插入子节点
	return child.insertCatalogue(parts[1:], articlePath, rootPath)
}

// flatten 本方法用于将当前节点及其所有子节点展开为一个平面列表,并返回该列表及其对应的(叶子节点的)文章列表
func (c *Catalogue) flatten() ([]*model.Catalogue, []*model.Article) {
	var catalogueOrms []*model.Catalogue
	var articleOrms []*model.Article

	c.flattenInto(0, &catalogueOrms, &articleOrms)
	return catalogueOrms, articleOrms
}

// flattenInto 本方法用于将当前节点及其所有子节点展开为一个平面列表,并将该列表及其对应的(叶子节点的)文章列表追加到给定的ORM切片中
func (c *Catalogue) flattenInto(parentID int64, catalogueOrms *[]*model.Catalogue, articleOrms *[]*model.Article) {
	if c == nil {
		return
	}

	catalogueOrm := &model.Catalogue{
		ID:            c.ID,
		BookID:        c.BookID,
		ParentID:      parentID,
		Name:          c.Name,
		Intro:         c.Intro,
		Type:          c.Type,
		DisplayStatus: c.DisplayStatus,
	}
	*catalogueOrms = append(*catalogueOrms, catalogueOrm)

	if c.Article != nil {
		articleOrm := &model.Article{
			ID:          c.Article.ID,
			CatalogueID: c.ID,
			Path:        c.Article.Path,
		}
		*articleOrms = append(*articleOrms, articleOrm)
	}

	for _, child := range c.Children {
		child.flattenInto(c.ID, catalogueOrms, articleOrms)
	}
}

// PrintTree 本方法用于打印目录树,包含根节点
func (c *Catalogue) PrintTree() {
	if c == nil {
		return
	}

	fmt.Println(c.Name)
	c.printChildren("")
}

// printChildren 本方法用于打印当前节点的子节点,并递归打印子节点的子节点
func (c *Catalogue) printChildren(prefix string) {
	for i, child := range c.Children {
		branch := "├── "
		nextPrefix := prefix + "│   "

		if i == len(c.Children)-1 {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		fmt.Printf("%s%s%s (%s) %#v\n", prefix, branch, child.Name, child.Type, child.Article)
		child.printChildren(nextPrefix)
	}
}
