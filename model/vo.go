package model

type TreeNode struct {
	Value    string     `json:"value"`
	Label    string     `json:"label"`
	Children []TreeNode `json:"children,omitempty"`
}

type FileInfo struct {
	Size int64  `json:"size"`
	Name string `json:"name"`
	Date string `json:"date" `
}
