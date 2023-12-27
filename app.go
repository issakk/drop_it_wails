package main

import (
	"context"
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
func (a *App) OpenFileDialog() []string {
	fmt.Println("openDialog")
	path, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})

	dir, _ := os.ReadDir(path)
	infos := lo.Map(dir, func(item os.DirEntry, index int) string {
		info, _ := item.Info()

		return filepath.Join(path, info.Name())
	})

	return infos

}
