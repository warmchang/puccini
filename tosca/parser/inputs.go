package parser

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/parsing"
)

func (self *Context) SetInputs(inputs map[string]ard.Value) {
	self.Parser.lock.Lock()
	defer self.Parser.lock.Unlock()

	parsing.SetInputs(self.Root.EntityPtr, inputs)
}
