package manifests

import (
	"io"
	"io/ioutil"
	"path/filepath"
)

func Asset(name string) (io.ReadCloser, error) {
	return assets.Open(name)
}

func AssetString(name string) (string, error) {
	asset, err := Asset(name)
	if err != nil {
		return "", err
	}
	defer asset.Close()

	b, err := ioutil.ReadAll(asset)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func AssetListDir(path string) ([]string, error) {
	dir, err := assets.Open(path)
	if err != nil {
		return nil, err
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	names := []string{}

	for _, fi := range files {
		names = append(names, filepath.Join(path, fi.Name()))
	}

	return names, nil
}

func MustAsset(name string) io.ReadCloser {
	r, err := Asset(name)
	if err != nil {
		panic(err)
	}

	return r
}

func MustAssetString(name string) string {
	data, err := AssetString(name)
	if err != nil {
		panic(err)
	}

	return data
}

func MustAssetListDir(path string) []string {
	files, err := AssetListDir(path)
	if err != nil {
		panic(err)
	}

	return files
}
