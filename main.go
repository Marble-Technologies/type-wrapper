package main

import (
	"os"

	"github.com/spf13/afero"

	"github.com/Marble-Technologies/type-wrapper/cmd"
)

func main() {
	cmd.Execute(afero.NewOsFs(), os.Args)
}
