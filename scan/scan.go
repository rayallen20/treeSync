package scan

import (
	"fmt"
	"path/filepath"
)

// Dir 本函数用于递归遍历给定的目录,返回所有普通文件的路径列表.该路径列表为相对给定目录 dir 的路径,分隔符统一为 /
// Tips:
// 1. 在递归过程中,会跳过隐藏项
// 2. 在递归过程中,会跳过符号链接
func Dir(dir string) ([]string, error) {
	w := &walker{
		root: dir,
	}

	err := filepath.WalkDir(dir, w.visit)
	if err != nil {
		return nil, fmt.Errorf("扫描目录失败: %w", err)
	}

	return w.files, nil
}
