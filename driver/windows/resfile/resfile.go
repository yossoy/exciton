package resfile

import (
	"os"
	"time"
)

type Item struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Size     int64       `json:"size"`
	Mode     os.FileMode `json:"mode"`
	ModTime  time.Time   `json:"mod_time"`
	IsDir    bool        `json:"isDir"`
	Children Items       `json:"children,omitempty"`
}

type Items map[string]*Item

const FileMapJsonID = 500

const FileIDStart = 1000
