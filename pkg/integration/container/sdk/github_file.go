package sdk

import (
	"bytes"
	"net/http"
	"os"
	"sync"

	"github.com/puppetlabs/relay/pkg/debug"
)

type githubFile struct {
	name    string
	content *githubContent
	client  *githubClient

	buf    bytes.Buffer
	filler sync.Once
}

func (gf *githubFile) fillBuffer() error {
	var returnErr error

	gf.filler.Do(func() {
		data := gf.content.Files[0]

		arr, err := gf.client.Download(data.DownloadURL)

		if err != nil {
			returnErr = err
		} else {
			// TODO: technically partial write is possible here. But I'm not
			// entirely sure we want to account for htat?
			gf.buf.Write(arr)
		}
	})

	return returnErr
}

func (gf *githubFile) Close() error {
	return nil
}

func (gf *githubFile) Read(buf []byte) (int, error) {
	if gf.content.Type == "file" {
		err := gf.fillBuffer()

		// not entirely sure what these possible failure modes are.
		if err != nil {
			debug.Logf("error fetching content from Github: %v", err)
			return 0, err
		}

		// Otherwise, we're good to go here!
		return gf.buf.Read(buf)
	} else {
		return 0, ErrNotSupported
	}
}

func (gf *githubFile) Seek(offset int64, whence int) (int64, error) {
	return 0, ErrNotSupported
}

// Readdir does...something. I have no idea what. This implementation actually
// just always returns an ErrNotSupported error because I couldn't figure out
// what behavior for this function should actually be.
func (gf *githubFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, ErrNotSupported
}

func (f *githubFile) Stat() (os.FileInfo, error) {
	return githubFileInfo{f.name, f.content}, nil
}

func newGithubFile(name string, content *githubContent, client *githubClient) (http.File, error) {
	return &githubFile{name: name, content: content, client: client}, nil
}
