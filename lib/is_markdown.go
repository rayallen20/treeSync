package lib

import (
	"path"
	"strings"
)

// IsMarkdown 本方法用于判断给定的文件名是否为 Markdown 文件
// Tips: .md/.MD/.Md/.mD 均会被判定为 Markdown 文件
func IsMarkdown(filename string) bool {
	extName := path.Ext(filename)
	return strings.EqualFold(extName, ".md")
}
