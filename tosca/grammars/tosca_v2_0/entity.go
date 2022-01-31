package tosca_v2_0

import (
	"sync"

	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
)

//
// Entity
//

type Entity struct {
	Context *tosca.Context `traverse:"ignore" json:"-" yaml:"-"`

	lock        util.RWLocker
	inheritOnce sync.Once
	renderOnce  sync.Once
}

func NewEntity(context *tosca.Context) *Entity {
	return &Entity{
		Context: context,
		lock:    util.NewMockRWLocker(),
	}
}

// tosca.Contextual interface
func (self *Entity) GetContext() *tosca.Context {
	return self.Context
}

// util.LockableEntity interface
func (self *Entity) GetEntityLock() util.RWLocker {
	return self.lock
}
