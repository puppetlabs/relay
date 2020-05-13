package sdk

import (
	"os"
	"path/filepath"
	"time"
)

type githubFileInfo struct {
	name    string
	content *githubContent
}

func (i githubFileInfo) Name() string {
	return filepath.Base(i.name)
}

func (i githubFileInfo) Size() int64 {
	if i.content.Type == "file" {
		return i.content.Files[0].Size
	}

	return 0
}

// Mode returns the FileMode for this resource. I have no idea what the correct
// thing to return here is.
func (i githubFileInfo) Mode() os.FileMode {
	return os.ModePerm
}

func (i githubFileInfo) ModTime() time.Time {
	return time.Now()
}

func (i githubFileInfo) IsDir() bool {
	return i.content.Type == "dir"
}

func (i githubFileInfo) Sys() interface{} {
	return nil
}
