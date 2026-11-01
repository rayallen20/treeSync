package reporter

type Report struct {
	CatalogueReport *CatalogueReport // CatalogueReport 目录变更情况报告
	ArticleReport   *ArticleReport   // ArticleReport 文章变更情况报告
}
