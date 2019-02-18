package web

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/yossoy/exciton/log"

	"github.com/yossoy/exciton/app"
	"github.com/yossoy/exciton/dialog"
	"github.com/yossoy/exciton/driver"
	"github.com/yossoy/exciton/event"
)

type fakeReadCloser struct {
	tempData *bytes.Buffer
}

func (frc *fakeReadCloser) Read(p []byte) (n int, err error) {
	return frc.tempData.Read(p)
}

func (frc *fakeReadCloser) Close() error {
	return nil
}

type openFileItem struct {
	origName     string
	tempFilePath string
	tempData     []byte
	size         int64
}

func (ofi *openFileItem) Name() string {
	return ofi.origName
}

func (ofi *openFileItem) Open() (io.ReadCloser, error) {
	if ofi.tempData != nil {
		return &fakeReadCloser{tempData: bytes.NewBuffer(ofi.tempData)}, nil
	}
	return os.Open(ofi.tempFilePath)
}

func (ofi *openFileItem) Size() int64 {
	return ofi.size
}

func (ofi *openFileItem) LocalFilePath() (string, error) {
	return "", fmt.Errorf("File %q is temporary file", ofi.origName)
}

func (ofi *openFileItem) IsTemporary() bool {
	return true
}

func (ofi *openFileItem) Cleanup() {
	if ofi.tempFilePath != "" {
		log.PrintDebug("ofi::Cleanup: %q", ofi.tempFilePath)
		os.Remove(ofi.tempFilePath)
		ofi.tempFilePath = ""
	}
}

func initializeDialog(serializer driver.DriverEventSerializer) error {
	app.AppClass.AddHandlerWithResult("showMessageBox", func(e *event.Event, callback event.ResponceCallback) {
		// TODO: 非NativeなWebの場合、シリアライザを通す意味はない?
		serializer.RelayEventWithResult(e, callback)
	})
	app.AppClass.AddHandlerWithResult("showOpenDialog", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, func(ee event.Result) {
			if ee.Error() != nil {
				callback(event.NewErrorResult(ee.Error()))
				return
			}
			var result *dialog.OpenFileResult
			if err := ee.Value().Decode(&result); err != nil {
				callback(event.NewErrorResult(err))
			}
			callback(event.NewValueResult(event.NewValue(result)))
		})
	})
	app.AppClass.AddHandlerWithResult("showSaveDialog", func(e *event.Event, callback event.ResponceCallback) {
		serializer.RelayEventWithResult(e, callback)
	})
	driver.AddInitProc(driver.InitProcTimingPostStartup, addOpenDialogUploadForm)

	return nil
}

func addOpenDialogUploadForm(si *driver.StartupInfo) error {
	si.Router.HandleFunc("/webFileOpenDialog", func(w http.ResponseWriter, r *http.Request) {
		log.PrintDebug("webFileOpenDialog called")
		if r.Method != "POST" {
			http.Error(w, "Allowed POST method only", http.StatusMethodNotAllowed)
			return
		}
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.PrintError("ParseMultipartForm failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rns, ok := r.MultipartForm.Value["openDialogResponceNo"]
		if !ok || len(rns) != 1 {
			log.PrintError("Invalid Form value")
			http.Error(w, "Invalid Form value", http.StatusInternalServerError)
			return
		}
		responceNo, err := strconv.ParseInt(rns[0], 10, 64)
		if err != nil {
			log.PrintError("%q parseInt failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.PrintDebug("ReponceNo = %d", responceNo)
		fhs, ok := r.MultipartForm.File["selFile"]
		if !ok {
			log.PrintError("Invalid selFile value")
			http.Error(w, "Invalid Form value", http.StatusInternalServerError)
			return
		}

		// copy temp file
		var ofr dialog.OpenFileResult
		for i, fh := range fhs {
			log.PrintInfo("File[%d] %q %d", i, fh.Filename, fh.Size)
			mf, err := fh.Open()
			if err != nil {
				ofr.Cleanup()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var ofi *openFileItem
			if (2 << 20) < fh.Size {
				// > 2MB, temp file
				tf, err := ioutil.TempFile(os.TempDir(), "exciton-web")
				if err != nil {
					ofr.Cleanup()
					mf.Close()
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				log.PrintDebug("addOpenDialogUploadForm: create temp file: %q", tf.Name())
				io.Copy(tf, mf)
				ofi = &openFileItem{
					origName:     fh.Filename,
					tempFilePath: tf.Name(),
					size:         fh.Size,
				}
				tf.Close()
			} else {
				//  <= 2MB, memory
				b := &bytes.Buffer{}
				io.Copy(b, mf)
				ofi = &openFileItem{
					origName: fh.Filename,
					tempData: b.Bytes(),
					size:     fh.Size,
				}
			}
			ofr.Items = append(ofr.Items, ofi)
			mf.Close()
		}
		platform.responceCallbackByValue(&ofr, int(responceNo))
	})
	return nil
}
