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

// buildFromPath 本方法用于根据给定的目录路径集合构建书籍目录
func (c *Catalogue) buildFromPath(bookID int64, bookTitle, rootPath string) error {
	c.ID = snowflake.Next()
	c.Name = bookTitle
	c.Type = enum.TypeFolder
	c.BookID = bookID
	c.DisplayStatus = enum.DisplayHidden

	subPaths, err := scan.Dir(rootPath)
	if err != nil {
		return err
	}

	for _, subPath := range subPaths {
		subPath = strings.Trim(subPath, "/")
		if subPath == "" {
			continue
		}

		parts := strings.Split(subPath, "/")
		err = c.insertCatalogue(parts, subPath, rootPath)
		if err != nil {
			return err
		}
	}

	return nil
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
		// 叶子节点不需要简介(因为直接显示文章了)
		child.Intro = ""

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

// saveTree 本方法用于将当前节点及其所有子节点展开为一个平面列表,并将该列表及其对应的(叶子节点的)文章列表落盘到数据库中
func (c *Catalogue) saveTree() (*reporter.Report, error) {
	catalogueOrms, articleOrms := c.flatten()
	orm := &model.Catalogue{}
	err := orm.CreateCatalogueTree(catalogueOrms, articleOrms)
	if err != nil {
		return nil, errors.New("落盘书籍目录失败: " + err.Error())
	}

	report := &reporter.Report{
		CatalogueReport: &reporter.CatalogueReport{
			AddItemQuantity:    len(catalogueOrms),
			DeleteItemQuantity: 0,
		},
		ArticleReport: &reporter.ArticleReport{
			AddItemQuantity:    len(articleOrms),
			DeleteItemQuantity: 0,
		},
	}

	return report, nil
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

// buildFromRows 本方法用于根据给定的目录ORM集合构建书籍目录
func (c *Catalogue) buildFromRows(rows []*model.Catalogue) error {
	if len(rows) == 0 {
		return errors.New("目录ORM集合为空")
	}

	// 遍历ORM集合 将ORM转换为业务对象(也将叶子节点的文章ORM转换为业务对象)
	catalogues := make(map[int64]*Catalogue, len(rows))
	
	for _, row := range rows {
		if row == nil {
			return errors.New("目录ORM集合中存在nil")
		}

		_, exists := catalogues[row.ID]
		if exists {
			return errors.New(fmt.Sprintf("目录ORM集合中存在重复的ID: %d", row.ID))
		}

		catalogue := &Catalogue{}
		catalogue.fillByORM(row)

		catalogues[row.ID] = catalogue
	}

	// 遍历ORM集合 将业务对象组织为树形结构
	var root *Catalogue

	for _, row := range rows {
		catalogue := catalogues[row.ID]

		// 处理根节点
		if row.ParentID == 0 {
			if root != nil {
				return errors.New("目录ORM集合中存在多个根节点")
			}

			root = catalogue
			continue
		}

		// 处理非根节点
		parent, exists := catalogues[row.ParentID]
		if !exists {
			msg := fmt.Sprintf("目录 %d 的父节点 %d 不存在", row.ID, row.ParentID)
			return errors.New(msg)
		}

		_, exists = parent.index[catalogue.Name]
		if exists {
			msg := fmt.Sprintf("父目录 %d 下存在同名节点 %d", parent.ID, catalogue.ID)
			return errors.New(msg)
		}

		parent.Children = append(parent.Children, catalogue)
		parent.index[catalogue.Name] = catalogue
	}

	if root == nil {
		return errors.New("目录ORM集合中不存在根节点")
	}

	// 检查是否存在从根节点出发不可达的节点
	isReachable := c.isReachable(root, len(rows))
	if !isReachable {
		return errors.New("目录ORM集合中存在不可达的节点")
	}

	// 检查通过后 将根节点赋值给当前接收者
	*c = *root
	return nil
}

// fillByORM 本方法用于根据给定的目录ORM对象填充目录业务结构体
func (c *Catalogue) fillByORM(orm *model.Catalogue) {
	c.ID = orm.ID
	c.Sort = orm.Sort
	c.Name = orm.Name
	c.Intro = orm.Intro
	c.Type = orm.Type
	c.DisplayStatus = orm.DisplayStatus
	c.BookID = orm.BookID
	c.index = make(map[string]*Catalogue)

	if orm.Article != nil {
		c.Article = &Article{}
		c.Article.fillByOrm(orm.Article)
	}
}

// isReachable 本方法用于判断给定的树中的节点是否都可以从根节点出发到达,若存在不可达的节点则返回false,否则返回true
func (c *Catalogue) isReachable(root *Catalogue, catalogueQuantity int) bool {
	reachable := make([]*Catalogue, 0, catalogueQuantity)
	reachable = append(reachable, root)
	for i := 0; i < len(reachable); i++ {
		reachable = append(reachable, reachable[i].Children...)
	}

	return len(reachable) == catalogueQuantity
}
