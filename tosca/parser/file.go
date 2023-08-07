package parser

import (
	"sort"
	"strings"

	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// NoEntity
//

type NoEntity struct {
	Context *parsing.Context
}

func NewNoEntity(toscaContext *parsing.Context) *NoEntity {
	return &NoEntity{toscaContext}
}

// parsing.Contextual interface
func (self *NoEntity) GetContext() *parsing.Context {
	return self.Context
}

//
// File
//

type File struct {
	EntityPtr       parsing.EntityPtr
	Container       *File
	Imports         Files
	NameTransformer parsing.NameTransformer

	importsLock util.RWLocker
}

func NewFileNoEntity(toscaContext *parsing.Context, container *File, nameTransformer parsing.NameTransformer) *File {
	return NewFile(NewNoEntity(toscaContext), container, nameTransformer)
}

func NewFile(entityPtr parsing.EntityPtr, container *File, nameTransformer parsing.NameTransformer) *File {
	self := File{
		EntityPtr:       entityPtr,
		Container:       container,
		NameTransformer: nameTransformer,
		importsLock:     util.NewDebugRWLocker(),
	}

	if container != nil {
		container.AddImport(&self)
	}

	return &self
}

func (self *File) AddImport(import_ *File) {
	self.importsLock.Lock()
	defer self.importsLock.Unlock()

	self.Imports = append(self.Imports, import_)
}

func (self *File) GetContext() *parsing.Context {
	return parsing.GetContext(self.EntityPtr)
}

// Print

func (self *File) PrintImports(indent int, treePrefix terminal.TreePrefix) {
	self.importsLock.RLock()
	length := len(self.Imports)
	imports := make(Files, length)
	copy(imports, self.Imports)
	self.importsLock.RUnlock()

	last := length - 1

	// Sort
	sort.Sort(imports)

	for i, file := range imports {
		isLast := i == last
		file.PrintNode(indent, treePrefix, isLast)
		file.PrintImports(indent, append(treePrefix, isLast))
	}
}

func (self *File) PrintNode(indent int, treePrefix terminal.TreePrefix, last bool) {
	treePrefix.Print(indent, last)
	terminal.Printf("%s\n", terminal.DefaultStylist.Value(self.GetContext().URL.String()))
}

//
// Files
//

type Files []*File

// sort.Interface
func (self Files) Len() int {
	return len(self)
}

// sort.Interface
func (self Files) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// sort.Interface
func (self Files) Less(i, j int) bool {
	iName := self[i].GetContext().URL.String()
	jName := self[j].GetContext().URL.String()
	return strings.Compare(iName, jName) < 0
}
