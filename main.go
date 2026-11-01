package main

import (
	"flag"
	"fmt"
	"os"
	"treeSync/biz"
	"treeSync/conf"
	"treeSync/db"
	"treeSync/snowflake"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	path := flag.String("path", "", "待处理的目录路径 (必填)")
	bookID := flag.Int64("bookID", 0, "书籍ID (必填, 正整数)")
	flag.Parse()

	book := &biz.Book{}
	report, err := book.Build(*bookID, *path)
	if err != nil {
		fmt.Printf("构建书籍目录失败: %v\n", err)
		return
	}

	fmt.Printf("构建书籍目录成功:\n")
	fmt.Printf("书籍ID: %d\n", book.ID)
	fmt.Printf("新增目录数: %d\n", report.CatalogueReport.AddItemQuantity)
	fmt.Printf("删除目录数: %d\n", report.CatalogueReport.DeleteItemQuantity)
	fmt.Printf("变更目录数: %d\n", report.CatalogueReport.UpdateItemQuantity)
	fmt.Printf("新增文章数: %d\n", report.ArticleReport.AddItemQuantity)
	fmt.Printf("删除文章数: %d\n", report.ArticleReport.DeleteItemQuantity)
	fmt.Printf("变更文章数: %d\n", report.ArticleReport.UpdateItemQuantity)
}

// run 本函数用于解析命令行参数,并执行对应的命令,返回错误信息
func run() error {
	err := conf.Load()
	if err != nil {
		return err
	}
	confInstance := conf.Get()

	err = db.Connect(confInstance.Database)
	if err != nil {
		return err
	}

	err = snowflake.Init(confInstance.Snowflake.NodeID)
	if err != nil {
		return err
	}
	return nil
}
