package dialog

import (
	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/event"
	// ievent "github.com/yossoy/exciton/internal/event"
)

type msgBoxOpt struct {
	Type     dialog.MessageBoxType `json:"type"`
	Title    string                `json:"title"`
	Message  string                `json:"message"`
	WindowID string                `json:"windowId"`
	dialog.MessageBoxConfig
}

func makeMsgBoxOpt(parent string, message string, title string, messageBoxType dialog.MessageBoxType, cfg *dialog.MessageBoxConfig) *msgBoxOpt {
	tmpl := &msgBoxOpt{Type: messageBoxType, Title: title, Message: message}
	if parent != "" {
		tmpl.WindowID = parent
	}
	if cfg != nil {
		tmpl.MessageBoxConfig = *cfg
	}
	if tmpl.Buttons == nil {
		switch tmpl.Type {
		case dialog.MessageBoxTypeQuestion:
			tmpl.Buttons = []string{"YES", "NO"}
		default:
			tmpl.Buttons = []string{"OK"}
		}
	}
	return tmpl
}

// TODO: owner は実質App. package flat化でなんとかしよう……
type Owner interface {
	event.EventTarget
}

func ShowMessageBoxAsync(owner Owner, windowID string, message string, title string, messageBoxType dialog.MessageBoxType, cfg *dialog.MessageBoxConfig, handler func(int, error)) error {
	opt := makeMsgBoxOpt(windowID, message, title, messageBoxType, cfg)
	err := event.EmitWithCallback(owner, "showMessageBox", event.NewValue(opt), func(result event.Result) {
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

func ShowMessageBox(owner Owner, windowID string, message string, title string, messageBoxType dialog.MessageBoxType, cfg *dialog.MessageBoxConfig) (int, error) {
	ch := make(chan interface{})
	err := ShowMessageBoxAsync(owner, windowID, message, title, messageBoxType, cfg, func(result int, err error) {
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

type openDialogOpt struct {
	WindowID string `json:"windowId"`
	dialog.FileDialogConfig
}

func ShowOpenDialogAsync(owner Owner, windowID string, cfg *dialog.FileDialogConfig, handler func(*dialog.OpenFileResult, error)) error {
	opt := &openDialogOpt{}
	if windowID != "" {
		opt.WindowID = windowID
	}
	if cfg != nil {
		opt.FileDialogConfig = *cfg
	}
	if opt.Title == "" {
		opt.Title = "Open"
	}
	if opt.Properties == 0 {
		opt.Properties = dialog.OpenDialogForOpenFile
	}
	err := event.EmitWithCallback(owner, "showOpenDialog", event.NewValue(opt), func(e event.Result) {
		if e.Error() != nil {
			handler(nil, e.Error())
			return
		}

		var result *dialog.OpenFileResult
		if err := e.Value().Decode(&result); err != nil {
			handler(nil, err)
		}
		handler(result, nil)
	})
	return err
}

func ShowOpenDialog(owner Owner, windowID string, cfg *dialog.FileDialogConfig) (*dialog.OpenFileResult, error) {
	ch := make(chan interface{})
	err := ShowOpenDialogAsync(owner, windowID, cfg, func(result *dialog.OpenFileResult, err error) {
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
	return r.(*dialog.OpenFileResult), nil
}

func ShowSaveDialogAsync(owner Owner, windowID string, cfg *dialog.FileDialogConfig, handler func(string, error)) error {
	opt := &openDialogOpt{}
	if windowID != "" {
		opt.WindowID = windowID
	}
	if cfg != nil {
		opt.FileDialogConfig = *cfg
	}
	if opt.Title == "" {
		opt.Title = "Save"
	}
	err := event.EmitWithCallback(owner, "showSaveDialog", event.NewValue(opt), func(e event.Result) {
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

func ShowSaveDialog(owner Owner, windowID string, cfg *dialog.FileDialogConfig) (string, error) {
	ch := make(chan interface{})
	err := ShowSaveDialogAsync(owner, windowID, cfg, func(result string, err error) {
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
