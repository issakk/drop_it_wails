package main

import (
	"context"
	"dropit/model"
	"fmt"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path/filepath"
)

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
