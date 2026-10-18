package scan

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// walker 本结构体表示一次目录树扫描的遍历状态
type walker struct {
	root  string   // root 被扫描的根目录 用于识别根目录自身和生成相对路径
	files []string // files 扫描出的所有普通文件的相对路径列表
}

// visit 本方法用于处理 filepath.WalkDir 的回调逻辑,并将扫描出的普通文件路径收集到 walker.files 中
func (w *walker) visit(path string, d fs.DirEntry, err error) error {
	// WalkDir 自身读取出错 则直接终止 不产出结果
	if err != nil {
		return err
	}

	// 根目录不参与规则判定
	if path == w.root {
		return nil
	}

	// 跳过隐藏项 (目录整枝跳过 文件单个跳过)
	if strings.HasPrefix(d.Name(), ".") {
		if d.IsDir() {
			return filepath.SkipDir
		}

		return nil
	}

	// 跳过符号链接
	if d.Type()&fs.ModeSymlink != 0 {
		return nil
	}

	// 目录本身不产出
	if d.IsDir() {
		return nil
	}

	// 非普通文件(设备文件/socket)跳过
	if !d.Type().IsRegular() {
		return nil
	}

	rel, err := filepath.Rel(w.root, path)
	if err != nil {
		return err
	}
	w.files = append(w.files, filepath.ToSlash(rel))

	return nil
}
