package biz

import (
	"errors"
	"treeSync/model"
	"treeSync/reporter"
)

// CatalogueSyncer 本结构体用于表示一次目录树同步任务
type CatalogueSyncer struct {
	originTree  *Catalogue // originTree 旧目录树(即数据库中已存在的目录树)
	currentTree *Catalogue // currentTree 新目录树(即扫描文件系统后构建的目录树)
}

// Compare 本方法用于比较新旧目录树的差异,并返回目录树变更情况
func (s *CatalogueSyncer) Compare() (*model.CatalogueChanges, error) {
	if s == nil {
		return nil, errors.New("目录树同步对象为nil")
	}

	if s.currentTree == nil {
		return nil, errors.New("新目录树为nil")
	}

	if s.originTree != nil && s.originTree.BookID != s.currentTree.BookID {
		return nil, errors.New("新旧目录树的书籍ID不一致")
	}

	changes := &model.CatalogueChanges{}
	s.compareNodes(s.originTree, s.currentTree, 0, changes)

	return changes, nil
}

// compareNode 本方法用于比对2个节点的差异,并将差异记录到给定的目录树变更对象中
func (s *CatalogueSyncer) compareNodes(origin *Catalogue, current *Catalogue, parentID int64, changes *model.CatalogueChanges) {
	// 原节点不存在且新节点存在 则新增新节点及其全部后代节点
	if origin == nil && current != nil {
		current.flattenInto(parentID, &changes.AddCatalogues, &changes.AddArticles)
		return
	}

	// 新节点在原树中不存在 则删除原节点及其全部后代节点
	if current == nil {
		catalogues, articles := origin.flatten()

		for _, article := range articles {
			changes.DeleteArticleIDs = append(changes.DeleteArticleIDs, article.ID)
		}

		// Tips: flatten()得到的[]*model.Catalogue是父节点在前子节点在后的顺序,因此需要反向遍历,即先删除所有子节点再删除父节点
		for i := len(catalogues) - 1; i >= 0; i-- {
			changes.DeleteCatalogueIDs = append(changes.DeleteCatalogueIDs, catalogues[i].ID)
		}

		return
	}

	// 节点名称不同,则进入删除旧节点及其全部后代节点,新增新节点及其全部后代节点的逻辑
	if origin.Name != current.Name {
		// 删除旧节点及其全部后代节点
		s.compareNodes(origin, nil, parentID, changes)

		// 新增新节点及其全部后代节点
		s.compareNodes(nil, current, parentID, changes)
		return
	}

	// 节点名称相同,则复用旧节点的ID和Sort
	current.ID = origin.ID
	if origin.Sort != nil {
		current.Sort = origin.Sort
	}

	// 处理节点对应的文章信息
	// 两个节点均为叶子节点,则复用原节点的文章ID
	if origin.Article != nil && current.Article != nil {
		current.Article.ID = origin.Article.ID
	}

	// 原节点为叶子节点,新节点为非叶节点,则删除旧节点的文章
	if origin.Article != nil && current.Article == nil {
		changes.DeleteArticleIDs = append(changes.DeleteArticleIDs, origin.Article.ID)
	}

	// 原节点为非叶节点,新节点为叶子节点,则新增新节点的文章
	if origin.Article == nil && current.Article != nil {
		article := &model.Article{
			ID:          current.Article.ID,
			CatalogueID: current.ID,
			Path:        current.Article.Path,
		}
		changes.AddArticles = append(changes.AddArticles, article)
	}

	// 处理原节点的子节点: 匹配或删除
	for _, originChild := range origin.Children {
		currentChild := current.index[originChild.Name]
		s.compareNodes(originChild, currentChild, current.ID, changes)
	}

	// 处理新节点的子节点: 新增
	for _, currentChild := range current.Children {
		_, exists := origin.index[currentChild.Name]

		if !exists {
			s.compareNodes(nil, currentChild, current.ID, changes)
		}
	}
}

// Sync 本方法用于完成一次2颗目录树的比对与落盘,并返回变更结果报告
func (s *CatalogueSyncer) Sync() (*reporter.Report, error) {
	changes, err := s.Compare()
	if err != nil {
		return nil, errors.New("比较目录树失败: " + err.Error())
	}

	catalogueOrm := &model.Catalogue{}
	err = catalogueOrm.ApplyCatalogueChanges(changes)
	if err != nil {
		return nil, errors.New("同步目录树失败: " + err.Error())
	}

	report := &reporter.Report{
		CatalogueReport: &reporter.CatalogueReport{
			AddItemQuantity:    len(changes.AddCatalogues),
			DeleteItemQuantity: len(changes.DeleteCatalogueIDs),
		},
		ArticleReport: &reporter.ArticleReport{
			AddItemQuantity:    len(changes.AddArticles),
			DeleteItemQuantity: len(changes.DeleteArticleIDs),
		},
	}

	return report, nil
}
