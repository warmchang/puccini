package parser

import (
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/puccini/tosca"
)

func (self *Context) AddHierarchies() {
	self.Root.MergeHierarchies(make(tosca.HierarchyContext), self.HierarchiesWork)
}

func (self *Unit) MergeHierarchies(hierarchyContext tosca.HierarchyContext, work *CoordinatedWork) {
	context := self.GetContext()

	if promise, ok := work.Start(context.URL.Key()); ok {
		defer promise.Release()

		self.importsLock.RLock()
		imports := make(Units, len(self.Imports))
		copy(imports, self.Imports)
		self.importsLock.RUnlock()

		for _, import_ := range imports {
			import_.MergeHierarchies(hierarchyContext, work)
			context.Hierarchy.Merge(import_.GetContext().Hierarchy, hierarchyContext)
		}

		logHierarchies.Debugf("create: %s", context.URL.String())
		hierarchy := tosca.NewHierarchyFor(self.EntityPtr, hierarchyContext)
		context.Hierarchy.Merge(hierarchy, hierarchyContext)
		// TODO: do we need this?
		//context.Hierarchy.AddTo(self.EntityPtr)
	}
}

// Print

func (self *Context) PrintHierarchies(indent int) {
	self.unitsLock.RLock()
	defer self.unitsLock.RUnlock()

	for _, import_ := range self.Units {
		context := import_.GetContext()
		if !context.Hierarchy.Empty() {
			terminal.PrintIndent(indent)
			terminal.Printf("%s\n", terminal.Stylize.Value(context.URL.String()))
			context.Hierarchy.Print(indent)
		}
	}
}
