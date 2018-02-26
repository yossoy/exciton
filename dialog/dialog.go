package dialog

import (
	"github.com/yossoy/exciton/event"
	"github.com/yossoy/exciton/window"
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

type msgBoxOpt struct {
	Type     MessageBoxType `json:"type"`
	Title    string         `json:"title"`
	Message  string         `json:"message"`
	WindowID string         `json:"windowId"`
	MessageBoxConfig
}

func makeMsgBoxOpt(parent *window.Window, message string, title string, messageBoxType MessageBoxType, cfg *MessageBoxConfig) *msgBoxOpt {
	tmpl := &msgBoxOpt{Type: messageBoxType, Title: title, Message: message}
	if parent != nil {
		tmpl.WindowID = parent.ID
	}
	if cfg != nil {
		tmpl.MessageBoxConfig = *cfg
	}
	if tmpl.Buttons == nil {
		switch tmpl.Type {
		case MessageBoxTypeQuestion:
			tmpl.Buttons = []string{"YES", "NO"}
		default:
			tmpl.Buttons = []string{"OK"}
		}
	}
	return tmpl
}

func ShowMessageBoxAsync(window *window.Window, message string, title string, messageBoxType MessageBoxType, cfg *MessageBoxConfig, handler func(int, error)) error {
	opt := makeMsgBoxOpt(window, message, title, messageBoxType, cfg)
	err := event.EmitWithCallback("/dialog/-/showMessageBox", event.NewValue(opt), func(result event.Result) {
		if result.Error() != nil {
			handler(-1, result.Error())
			return
		}
		var button int
		err := result.Value().Decode(&button)
		if err != nil {
			handler(-1, err)
			return
		}
		handler(button, nil)
	})
	return err
}

func ShowMessageBox(window *window.Window, message string, title string, messageBoxType MessageBoxType, cfg *MessageBoxConfig) (int, error) {
	ch := make(chan interface{})
	err := ShowMessageBoxAsync(window, message, title, messageBoxType, cfg, func(result int, err error) {
		if err != nil {
			ch <- err
		}
		ch <- result
	})
	if err != nil {
		return -1, err
	}
	r := <-ch
	if err, ok := r.(error); ok {
		return -1, err
	}
	return r.(int), nil
}

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
type openDialogOpt struct {
	WindowID string `json:"windowId"`
	FileDialogConfig
}

func ShowOpenDialogAsync(parent *window.Window, cfg *FileDialogConfig, handler func([]string, error)) error {
	opt := &openDialogOpt{}
	if parent != nil {
		opt.WindowID = parent.ID
	}
	if cfg != nil {
		opt.FileDialogConfig = *cfg
	}
	if opt.Title == "" {
		opt.Title = "Open"
	}
	if opt.Properties == 0 {
		opt.Properties = OpenDialogForOpenFile
	}
	err := event.EmitWithCallback("/dialog/-/showOpenDialog", event.NewValue(opt), func(e event.Result) {
		if e.Error() != nil {
			handler(nil, e.Error())
			return
		}
		var files []string
		if err := e.Value().Decode(&files); err != nil {
			handler(nil, err)
		}
		handler(files, nil)
	})
	return err
}

func ShowOpenDialog(parent *window.Window, cfg *FileDialogConfig) ([]string, error) {
	ch := make(chan interface{})
	err := ShowOpenDialogAsync(parent, cfg, func(result []string, err error) {
		if err != nil {
			ch <- err
		}
		ch <- result
	})
	if err != nil {
		return nil, err
	}
	r := <-ch
	if err, ok := r.(error); ok {
		return nil, err
	}
	return r.([]string), nil
}

func ShowSaveDialogAsync(parent *window.Window, cfg *FileDialogConfig, handler func(string, error)) error {
	opt := &openDialogOpt{}
	if parent != nil {
		opt.WindowID = parent.ID
	}
	if cfg != nil {
		opt.FileDialogConfig = *cfg
	}
	if opt.Title == "" {
		opt.Title = "Save"
	}
	err := event.EmitWithCallback("/dialog/-/showSaveDialog", event.NewValue(opt), func(e event.Result) {
		if e.Error() != nil {
			handler("", e.Error())
			return
		}
		var file string
		if err := e.Value().Decode(&file); err != nil {
			handler("", err)
		}
		handler(file, nil)
	})
	return err
}

func ShowSaveDialog(parent *window.Window, cfg *FileDialogConfig) (string, error) {
	ch := make(chan interface{})
	err := ShowSaveDialogAsync(parent, cfg, func(result string, err error) {
		if err != nil {
			ch <- err
		}
		ch <- result
	})
	if err != nil {
		return "", err
	}
	r := <-ch
	if err, ok := r.(error); ok {
		return "", err
	}
	return r.(string), nil
}
