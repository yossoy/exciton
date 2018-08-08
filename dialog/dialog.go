package dialog

import (
	"io"
)

type MessageBoxConfig struct {
	Buttons   []string `json:"buttons"`
	DefaultID int      `json:"defaultId"`
	CancelID  int      `json:"cancelId"`
	Detail    string   `json:"detail"`
	NoLink    bool     `json:"noLink"`
	IconPath  string   `json:"iconPath"`
}

type MessageBoxType int

const (
	MessageBoxTypeNone MessageBoxType = iota
	MessageBoxTypeInfo
	MessageBoxTypeWarning
	MessageBoxTypeError
	MessageBoxTypeQuestion
)

type OpenDialogProperty int

const (
	OpenDialogForOpenFile OpenDialogProperty = 1 << iota
	OpenDialogForOpenDirectory
	OpenDialogWithMultiSelections
	OpenDialogWithCreateDirectory
	OpenDialogWithShowHiddenFiles
)

type FileDialogFilter struct {
	Name       string   `json:"name"`
	Extensions []string `json:"extensions"`
}

type FileDialogConfig struct {
	Title       string             `json:"title"`
	DefaultPath string             `json:"defaultPath"`
	ButtonLabel string             `json:"buttonLabel"`
	Filters     []FileDialogFilter `json:"filters,omitempty"`
	Properties  OpenDialogProperty `json:"properties"`
}

type OpenFileItem interface {
	Name() string
	Open() (io.ReadCloser, error)
	Size() int64
	LocalFilePath() (string, error)
	IsTemporary() bool
}

type TemporaryFileItem interface {
	Cleanup()
}

type OpenFileResult struct {
	Items []OpenFileItem `json:"items"`
}

func (ofr *OpenFileResult) Cleanup() {
	for _, itm := range ofr.Items {
		if tfi, ok := itm.(TemporaryFileItem); ok {
			tfi.Cleanup()
		}
	}
}
