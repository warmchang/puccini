package main

import (
	"github.com/tliron/go-kutil/util"
	"github.com/tliron/go-puccini/executables/puccini-csar/commands"

	_ "github.com/tliron/commonlog/simple"
)

func main() {
	util.ExitOnSignals()
	commands.Execute()
	util.Exit(0)
}
