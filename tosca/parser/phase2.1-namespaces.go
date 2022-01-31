package parser

import (
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/puccini/tosca"
)

func (self *Context) AddNamespaces() {
	self.Root.MergeNamespaces(self.NamespacesWork)
}

func (self *Unit) MergeNamespaces(work *CoordinatedWork) {
	context := self.GetContext()

	if promise, ok := work.Start(context.URL.Key()); ok {
		defer promise.Release()

		self.importsLock.RLock()
		imports := make(Units, len(self.Imports))
		copy(imports, self.Imports)
		self.importsLock.RUnlock()

		for _, import_ := range imports {
			import_.MergeNamespaces(work)
			context.Namespace.Merge(import_.GetContext().Namespace, import_.NameTransformer)
			context.ScriptletNamespace.Merge(import_.GetContext().ScriptletNamespace)
		}

		logNamespaces.Debugf("create: %s", context.URL.String())
		namespace := tosca.NewNamespaceFor(self.EntityPtr)
		context.Namespace.Merge(namespace, nil)
	}
}

// Print

func (self *Context) PrintNamespaces(indent int) {
	self.unitsLock.RLock()
	defer self.unitsLock.RUnlock()

	childIndent := indent + 1
	for _, import_ := range self.Units {
		context := import_.GetContext()
		if !context.Namespace.Empty() {
			terminal.PrintIndent(indent)
			terminal.Printf("%s\n", terminal.Stylize.Value(context.URL.String()))
			context.Namespace.Print(childIndent)
		}
	}
}
