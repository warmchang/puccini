package parser

import (
	"sync"

	"github.com/tliron/go-kutil/reflection"
	"github.com/tliron/go-kutil/util"
)

//
// Parser
//

type Parser struct {
	readCache          sync.Map // entityPtr or Promise
	lookupFieldsWork   reflection.EntityWork
	addHierarchyWork   reflection.EntityWork
	getInheritTaskWork reflection.EntityWork
	renderWork         reflection.EntityWork
	lock               util.RWLocker
}

func NewParser() *Parser {
	return &Parser{
		lookupFieldsWork:   make(reflection.EntityWork),
		addHierarchyWork:   make(reflection.EntityWork),
		getInheritTaskWork: make(reflection.EntityWork),
		renderWork:         make(reflection.EntityWork),
		lock:               util.NewDefaultRWLocker(),
	}
}
