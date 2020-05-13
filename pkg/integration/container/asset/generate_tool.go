// +build tools

package main

import (
	"log"
	"os"

	"github.com/puppetlabs/horsehead/v2/httputil/fs"
	"github.com/shurcooL/vfsgen"
)

var h = fs.DirWithoutModTimes("data")

func main() {
	err := vfsgen.Generate(h, vfsgen.Options{
		Filename:    "generate_assets.go",
		PackageName: os.Getenv("GOPACKAGE"),
	})
	if err != nil {
		log.Fatalln(err)
	}
}
