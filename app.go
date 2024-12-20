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
	"图片":    {"png", "jpg", "jpeg", "gif", "bmp", "tiff", "svg", "ico"},
	"电子书":   {"pdf", "mobi", "azw3", "epub", "djvu", "cbz", "cbr"},
	"文档资料":  {"txt", "md", "doc", "docx", "xlsx", "csv", "ppt", "pptx", "rtf", "odt", "ods", "odp"},
	"压缩包":   {"zip", "rar", "7z", "tar", "gz", "bz2", "xz"},
	"音频":    {"mp3", "wmv", "m4a", "flac", "wav", "ogg", "aac", "m4b"},
	"视频":    {"mp4", "mkv", "mov", "flv", "wmv", "rmvb"},
	"exe程序": {"exe", "sh", "bat", "py", "jar", "app", "dmg", "msi", "apk", "ipa"},
	"编程相关":  {"out", "log", "json", "yml", "yaml"},
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

func fileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
}
func copyFiles(m map[string][]fs.FileInfo, path string) int {
	count := 0
	backupPath := filepath.Join(path, "备份")
	if createDirIfNotExist(backupPath) {
		fmt.Println("创建备份目录成功")
		return count
	}

	for k, v := range m {
		tempPath := filepath.Join(path, k)
		createDirIfNotExist(tempPath)
		for _, i := range v {
			oldPath := filepath.Join(path, i.Name())
			file, _ := ioutil.ReadFile(oldPath)

			err := ioutil.WriteFile(filepath.Join(tempPath, i.Name()), file, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return -1
			}
			err = os.Rename(oldPath, filepath.Join(backupPath, i.Name()))
			count++
			if err != nil {
				fmt.Println(err)
				return -1
			}

		}
	}
	return count
}

func createDirIfNotExist(backupPath string) bool {
	if !fileIsExisted(backupPath) {
		err := os.Mkdir(backupPath, fs.ModeDir)
		if err != nil {
			fmt.Println(err)
			return true
		}
	}
	return false
}

func (a *App) Drop(path string) int {
	return copyFiles(readFiles(path), path)

}
func readFiles(path string) map[string][]fs.FileInfo {
	m := make(map[string][]fs.FileInfo)
	dir, _ := ioutil.ReadDir(path)
	for _, i := range dir {
		name := i.Name()
		if i.IsDir() {
			continue
		}
		ext := filepath.Ext(name)
		// 去掉扩展名前面的点号
		if strings.Contains(ext, ".") {
			ext = ext[1:]
		}
		for index, v := range fileMap {
			if lo.Contains(v, ext) {
				//fmt.Println(name, "是个", index)
				infos, ok := m[index]
				if !ok {
					infos = make([]fs.FileInfo, 0, 0)

				}
				infos = append(infos, i)
				m[index] = infos
			}
		}
		//fmt.Println(name+" 是文件夹 ", i.IsDir())

	}
	return m
}
