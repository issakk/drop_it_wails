package main

import (
	"context"
	"dropit/model"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var fileMap = map[string][]string{
	"图片":    {"png", "jpg", "jpeg", "gif", "bmp", "tiff", "svg", "ico", "webp", "psd", "raw", "heic"},
	"电子书":   {"pdf", "mobi", "azw3", "epub", "djvu", "cbz", "cbr", "txt", "chm"},
	"文档资料":  {"txt", "md", "doc", "docx", "xlsx", "csv", "ppt", "pptx", "rtf", "odt", "ods", "odp", "pages", "numbers", "key", "one", "xls", "xmind", "mindnode"},
	"压缩包":   {"zip", "rar", "7z", "tar", "gz", "bz2", "xz", "iso", "pkg"},
	"音频":    {"mp3", "wmv", "m4a", "flac", "wav", "ogg", "aac", "m4b", "ape", "aiff", "wma", "opus"},
	"视频":    {"mp4", "mkv", "mov", "flv", "wmv", "rmvb", "avi", "m4v", "webm", "3gp", "ts", "f4v"},
	"exe程序": {"exe", "sh", "bat", "py", "jar", "app", "dmg", "msi", "apk", "ipa", "deb", "rpm", "run", "bin"},
	"编程相关":  {"out", "log", "json", "yml", "yaml", "xml", "sql", "db", "sqlite", "ini", "conf", "config", "properties", "env", "toml"},
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
} /**/

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (a *App) OpenFileDialog() []model.TreeNode {
	fmt.Println("openDialog")
	path, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	//dir, _ := os.ReadDir(path)
	//dirList := lo.Filter(dir, func(item os.DirEntry, index int) bool {
	//	info, _ := item.Info()
	//	return info.IsDir()
	//})
	//return lo.Map(dirList, func(item os.DirEntry, index int) model.TreeNode {
	//	info, _ := item.Info()
	//	filePath := filepath.Join(path, info.Name())
	//	return model.TreeNode{Value: filePath, Label: info.Name()}
	//})

	node := NewNode(path)

	return []model.TreeNode{node}
}

func NewNode(path string) model.TreeNode {
	if path == "" {
		return model.TreeNode{}
	}
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	node := model.TreeNode{
		Value: path,
		Label: filepath.Base(path),
	}

	if info.IsDir() {
		childrenInfo, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}

		for _, childFileInfo := range childrenInfo {
			childPath := filepath.Join(path, childFileInfo.Name())
			if childFileInfo.IsDir() {
				node.Children = append(node.Children, NewNode(childPath))
			}
		}
	}

	return node
}

func (a *App) ListFileInfo(path string) []model.FileInfo {
	dir, _ := os.ReadDir(path)
	fileList := lo.Filter(dir, func(item os.DirEntry, index int) bool {
		info, _ := item.Info()
		return !info.IsDir()
	})
	return lo.Map(fileList, func(item os.DirEntry, index int) model.FileInfo {
		info, _ := item.Info()
		return model.FileInfo{Name: info.Name(), Size: info.Size(), Date: info.ModTime().Format("2006-01-02 15:04:05")}
	})
}

func (a *App) Drop(path string) (int, error) {
	files, err := readFiles(path)
	if err != nil {
		return 0, fmt.Errorf("读取文件失败: %w", err)
	}
	return copyFiles(files, path)
}

func copyFiles(m map[string][]fs.FileInfo, path string) (int, error) {
	count := 0
	backupPath := filepath.Join(path, "备份")

	if err := ensureDir(backupPath); err != nil {
		return 0, fmt.Errorf("创建备份目录失败: %w", err)
	}

	for category, files := range m {
		categoryPath := filepath.Join(path, category)
		if err := ensureDir(categoryPath); err != nil {
			return 0, fmt.Errorf("创建分类目录失败: %w", err)
		}

		for _, fileInfo := range files {
			if err := moveFile(path, categoryPath, backupPath, fileInfo); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}

func moveFile(basePath, categoryPath, backupPath string, fileInfo fs.FileInfo) error {
	fileName := fileInfo.Name()
	oldPath := filepath.Join(basePath, fileName)

	// 先复制到分类目录
	if err := copyFileToDir(oldPath, categoryPath); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	// 再移动到备份目录
	if err := os.Rename(oldPath, filepath.Join(backupPath, fileName)); err != nil {
		return fmt.Errorf("移动到备份目录失败: %w", err)
	}

	return nil
}

func copyFileToDir(src, dstDir string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	dst := filepath.Join(dstDir, filepath.Base(src))
	return ioutil.WriteFile(dst, data, os.ModePerm)
}

func ensureDir(path string) error {
	if exists, err := pathExists(path); err != nil {
		return err
	} else if !exists {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readFiles(path string) (map[string][]fs.FileInfo, error) {
	m := make(map[string][]fs.FileInfo)

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	for _, fileInfo := range dir {
		if fileInfo.IsDir() {
			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileInfo.Name()), "."))

		for category, extensions := range fileMap {
			if lo.Contains(extensions, ext) {
				m[category] = append(m[category], fileInfo)
				break
			}
		}
	}
	return m, nil
}
