package window

import (
	"github.com/yossoy/exciton/dialog"
	idialog "github.com/yossoy/exciton/internal/dialog"
)

func (w *Window) ShowMessageBoxAsync(message string, title string, messageBoxType dialog.MessageBoxType, cfg *dialog.MessageBoxConfig, handler func(int, error)) error {
	return idialog.ShowMessageBoxAsync(w, w.ID, message, title, messageBoxType, cfg, handler)
}

func (w *Window) ShowMessageBox(message string, title string, messageBoxType dialog.MessageBoxType, cfg *dialog.MessageBoxConfig) (int, error) {
	return idialog.ShowMessageBox(w.owner, w.ID, message, title, messageBoxType, cfg)
}

func (w *Window) ShowOpenDialogAsync(cfg *dialog.FileDialogConfig, handler func(*dialog.OpenFileResult, error)) error {
	return idialog.ShowOpenDialogAsync(w.owner, w.ID, cfg, handler)
}

func (w *Window) ShowOpenDialog(cfg *dialog.FileDialogConfig) (*dialog.OpenFileResult, error) {
	return idialog.ShowOpenDialog(w.owner, w.ID, cfg)
}

func (w *Window) ShowSaveDialogAsync(cfg *dialog.FileDialogConfig, handler func(string, error)) error {
	return idialog.ShowSaveDialogAsync(w.owner, w.ID, cfg, handler)
}

func (w *Window) ShowSaveDialog(cfg *dialog.FileDialogConfig) (string, error) {
	return idialog.ShowSaveDialog(w.owner, w.ID, cfg)
}
