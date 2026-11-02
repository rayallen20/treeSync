package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"treeSync/biz"
	"treeSync/conf"
	"treeSync/db"
	"treeSync/snowflake"
)

func main() {
	err := run()
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "%v\n", err)
		if err != nil {
			log.Print("写入标准错误失败: ", err)
		}

		os.Exit(1)
	}

	path := flag.String("path", "", "待处理的目录路径 (必填)")
	bookID := flag.Int64("bookID", 0, "书籍ID (必填, 正整数)")
	flag.Parse()

	book := &biz.Book{}
	report, err := book.Build(*bookID, *path)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "构建书籍目录失败: %v\n", err)
		if err != nil {
			log.Print("写入标准错误失败: ", err)
		}

		os.Exit(1)
	}

	fmt.Printf("构建书籍目录成功:\n")
	fmt.Printf("书籍ID: %d\n", book.ID)
	fmt.Printf("新增目录数: %d\n", report.CatalogueReport.AddItemQuantity)
	fmt.Printf("删除目录数: %d\n", report.CatalogueReport.DeleteItemQuantity)
	fmt.Printf("新增文章数: %d\n", report.ArticleReport.AddItemQuantity)
	fmt.Printf("删除文章数: %d\n", report.ArticleReport.DeleteItemQuantity)
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
