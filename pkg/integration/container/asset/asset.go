package asset

import (
	"io/ioutil"
	"net/http"
)

var (
	FileSystem = assets
)

func Asset(name string) (http.File, error) {
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

func MustAsset(name string) http.File {
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
