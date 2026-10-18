package scan_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"treeSync/scan"
)

// fixture 本结构体用于表示一棵待构造的目录树,用于测试
type fixture struct {
	files    []string          // files 要创建的普通文件(父目录自动创建)
	dirs     []string          // dirs 要创建的空目录
	symlinks map[string]string // symlinks 要创造的符号链接
}

// makeTree 本函数用于在临时目录中按 fixture 中的描述创建真实文件树结构,返回临时根目录的绝对路径
// Tips: 临时目录会在测试结束后自动删除
func makeTree(t *testing.T, fx fixture) string {
	t.Helper()
	base := t.TempDir()

	for _, file := range fx.files {
		p := filepath.Join(base, filepath.FromSlash(file))

		err := os.MkdirAll(filepath.Dir(p), 0o755)
		if err != nil {
			t.Fatal(err)
		}

		err = os.WriteFile(p, []byte("test"), 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, dir := range fx.dirs {
		err := os.MkdirAll(filepath.Join(base, filepath.FromSlash(dir)), 0o755)
		if err != nil {
			t.Fatal(err)
		}
	}

	for link, target := range fx.symlinks {
		l := filepath.Join(base, filepath.FromSlash(link))
		err := os.Symlink(filepath.Join(base, filepath.FromSlash(target)), l)
		if err != nil {
			t.Fatal(err)
		}
	}

	return base
}

// TestDir 本函数用于测试 Dir 函数的正确性
func TestDir(t *testing.T) {
	cases := []struct {
		name string
		fx   fixture
		want []string
	}{
		{
			name: "嵌套文件_相对路径_斜杠分隔",
			fx: fixture{
				files: []string{
					"book/b.md",
					"book/a/y.md",
					"book/a/x/deep.md",
					"book/a/z.md",
				},
				dirs:     nil,
				symlinks: nil,
			},
			want: []string{
				"a/x/deep.md",
				"a/y.md",
				"a/z.md",
				"b.md",
			},
		},
		{
			name: "顺序_逐层字典序_非整条路径字典序",
			fx: fixture{
				files: []string{
					"book/a.md",
					"book/a-b/c.md",
					"book/a/z.md",
				},
				dirs:     nil,
				symlinks: nil,
			},
			want: []string{
				"a/z.md",
				"a-b/c.md",
				"a.md",
			},
		},
		{
			name: "隐藏文件_单个跳过_不影响同目录其他条目",
			fx: fixture{
				files: []string{
					"book/a/.hidden.md",
					"book/a/b/c.md",
					"book/a/z.md",
				},
				dirs:     nil,
				symlinks: nil,
			},
			want: []string{
				"a/b/c.md",
				"a/z.md",
			},
		},
		{
			name: "隐藏目录_整支跳过",
			fx: fixture{
				files: []string{
					"book/.git/x/in.md",
					"book/keep.md",
				},
				dirs:     nil,
				symlinks: nil,
			},
			want: []string{
				"keep.md",
			},
		},
		{
			name: "符号链接_不产出_不下钻",
			fx: fixture{
				files: []string{
					"book/real.md",
					"outside/o.md",
				},
				dirs: nil,
				symlinks: map[string]string{
					"book/linkfile.md": "outside/o.md", // 指向文件
					"book/linkdir":     "outside",      // 指向 book 之外的目录
					"book/loop":        "book",         // 指向自身, 若跟随会成环
				},
			},
			want: []string{
				"real.md",
			},
		},
		{
			name: "空目录_返回空_不报错",
			fx: fixture{
				files: nil,
				dirs: []string{
					"book",
				},
				symlinks: nil,
			},
			want: nil,
		},
		{
			name: "非md文件_照样返回",
			fx: fixture{
				files: []string{
					"book/img/p.jpg",
				},
				dirs:     nil,
				symlinks: nil,
			},
			want: []string{
				"img/p.jpg",
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			base := makeTree(t, testCase.fx)

			got, err := scan.Dir(filepath.Join(base, "book"))
			if err != nil {
				t.Fatalf("Dir()返回了意外的错误: %v", err)
			}

			if !slices.Equal(got, testCase.want) {
				t.Errorf("Dir()结果不符\n got: %q\nwant: %q", got, testCase.want)
			}
		})
	}
}

// TestDirHiddenRoot 本函数用于测试 Dir 函数在隐藏根目录下的行为
func TestDirHiddenRoot(t *testing.T) {
	f := fixture{
		files: []string{
			".book/a.md",
		},
		dirs:     nil,
		symlinks: nil,
	}
	base := makeTree(t, f)

	got, err := scan.Dir(filepath.Join(base, ".book"))
	if err != nil {
		t.Fatalf("Dir()返回了意外的错误: %v", err)
	}

	want := []string{
		"a.md",
	}
	if !slices.Equal(got, want) {
		t.Errorf("书籍目录自身以 . 开头不应被当作隐藏项\n got: %q\nwant: %q", got, want)
	}
}

// TestDirNotExist 本函数用于测试 Dir 函数在目录不存在时的行为
func TestDirNotExist(t *testing.T) {
	got, err := scan.Dir(filepath.Join(t.TempDir(), "nope"))

	if err == nil {
		t.Fatal("目录不存在时应返回错误, 实际为 nil")
	}

	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("错误应包装底层的不存在错误(%%w), 实际: %v", err)
	}

	if got != nil {
		t.Errorf("出错时不应产出结果, 实际: %q", got)
	}
}

// TestDirUnreadableSubDir 本函数用于测试 Dir 函数在子目录不可读时的行为
func TestDirUnreadableSubDir(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root 用户不受权限限制, 本用例无意义")
	}

	f := fixture{
		files: []string{
			"book/a.md",
			"book/locked/inside.md",
		},
		dirs:     nil,
		symlinks: nil,
	}
	base := makeTree(t, f)
	locked := filepath.Join(base, "book", "locked")

	err := os.Chmod(locked, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(locked, 0o755)
	}) // 先于 TempDir 的清理执行, 否则删不掉 locked/

	got, err := scan.Dir(filepath.Join(base, "book"))
	if err == nil {
		t.Fatal("子目录不可读时应返回错误, 实际为 nil")
	}
	if !errors.Is(err, fs.ErrPermission) {
		t.Errorf("错误应包装底层的权限错误(%%w), 实际: %v", err)
	}
	if got != nil {
		t.Errorf("整体失败时不应产出部分结果(a.md 已在 locked/ 之前被扫到), 实际: %q", got)
	}
}
