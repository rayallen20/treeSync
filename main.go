package main

import (
	"fmt"
	"os"
	"treeSync/biz/syncer"
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

// printReport 本函数用于打印结果
func printReport(report *syncer.Report) {

}
