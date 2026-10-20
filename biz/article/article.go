package article

const MaxPathChars = 255 // MaxPathChars 本常量用于定义文章路径的最大字符数限制,超过该限制的路径将被视为无效路径

// Article 文章业务结构体
type Article struct {
	ID   int64  // ID 文章ID
	Path string // Path 文章在文件系统中的路径(相对配置文件根路径的路径,含扩展名)
}

// IsMarkdown 本函数用于判断给定的文件名是否为 Markdown 文件,若是则返回true,否则返回false
// Tips: .md/.MD/.Md 均会被判定为 Markdown 文件
func IsMarkdown(filename string) bool {
	return false
}

// NameOf 本函数用于获取给定文件名的名称(不含扩展名)
func NameOf(filename string) string {
	return ""
}
