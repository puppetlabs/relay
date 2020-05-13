package sdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/puppetlabs/relay/pkg/debug"
)

// githubClient is a...not so great...client that we can use for interacting
// with content inside a public Git repository. For public repositories, you
// can use the API without any authentication which is nice.
type githubClient struct {
	client http.Client
}

// githubFileMetadata represents the metadata for an actual individual file
// within the Github API.
type githubFileMetadata struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	DownloadURL string `json:"download_url"`
}

type githubContent struct {
	// Type indicates the type of content what was found with the request. Will
	// be one of "file", "dir", and "unknown" (if something weird happens).
	Type string

	// Files is the list of files that was found. When `Type` is "file" there
	// will be one entry. When `Type` is "dir" there can be more than one entry.
	Files []githubFileMetadata
}

func (c *githubContent) UnmarshalJSON(buf []byte) error {
	// if this is an array then we know this is a directory.
	if bytes.HasPrefix(buf, []byte("[")) {
		c.Type = "dir"
		return json.Unmarshal(buf, &c.Files)
	}

	if bytes.HasPrefix(buf, []byte("{")) {
		// if it's an object
		c.Type = "file"

		var f githubFileMetadata

		if err := json.Unmarshal(buf, &f); err != nil {
			return err
		}

		c.Files = []githubFileMetadata{f}
		return nil
	}

	// TODO: What should we do h ere I wonder...?
	c.Type = "unknown"
	return nil
}

// Get fetchs the metadata for a given file. Use the Download function to
// actually get content out of Github.
func (g githubClient) Get(name string) (*githubContent, error) {
	debug.Logf("getting %s from Github", url(name))

	req, err := http.NewRequest("GET", url(name), nil)

	// This is catastrophic and client-side.
	if err != nil {
		return nil, err
	}

	resp, err := g.client.Do(req)

	// this error is likely network related and isn't application-level
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	var data githubContent

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// Download fetches the actual content out of a Github repository and returns
// it as a byte array.
func (g githubClient) Download(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", path, nil)

	// This is catastrophic and client-side.
	if err != nil {
		return nil, err
	}

	resp, err := g.client.Do(req)

	// this error is likely network related and isn't application-level
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	return ioutil.ReadAll(resp.Body)
}
