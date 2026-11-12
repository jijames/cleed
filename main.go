package main

import (
	"github.com/radulucut/cleed/cmd/cleed"
	"github.com/radulucut/cleed/internal/storage"
)

var (
	version = "dev"
)

func main() {
	cleed.Version = version
	storage.SetDefaultLocalDirs(`.\.data\config\`, `.\.data\`)
	cleed.Execute()
}
