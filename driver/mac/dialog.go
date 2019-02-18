package mac

/*
#include "driver.h"
#include "dialog.h"
*/
import "C"
import (
	"io"
	"os"
	"path/filepath"

	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/event"
)

type openFileItem struct {
	filePath string
}

func (ofi *openFileItem) Name() string {
	_, fn := filepath.Split(ofi.filePath)
	return fn
}

func (ofi *openFileItem) Open() (io.ReadCloser, error) {
	f, err := os.Open(ofi.filePath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (ofi *openFileItem) Size() int64 {
	s, err := os.Stat(ofi.filePath)
	if err != nil {
		return 0
	}
	return s.Size()
}

func (ofi *openFileItem) LocalFilePath() (string, error) {
	return ofi.filePath, nil
}

func (ofi *openFileItem) IsTemporary() bool {
	return false
}

func initializeDialog(serializer driver.DriverEventSerializer) {
	app.AppClass.AddHandlerWithResult("showMessageBox", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, callback)
	})
	app.AppClass.AddHandlerWithResult("showOpenDialog", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, func(ee event.Result) {
			if ee.Error() != nil {
				callback(event.NewErrorResult(ee.Error()))
				return
			}
			var filePaths []string
			if err := ee.Value().Decode(&filePaths); err != nil {
				callback(event.NewErrorResult(err))
			}
			ofr := &dialog.OpenFileResult{
				Items: make([]dialog.OpenFileItem, len(filePaths)),
			}
			for i, fp := range filePaths {
				ofr.Items[i] = &openFileItem{filePath: fp}
			}
			callback(event.NewValueResult(event.NewValue(ofr)))
		})
	})
	app.AppClass.AddHandlerWithResult("showSaveDialog", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, callback)
	})
	
	C.Dialog_Init()
}
