package parser

import (
	"sync"

	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/puccini/tosca"
)

func (self *Context) TraverseEntities(log logging.Logger, work EntityWork, traverse reflection.EntityTraverser) {
	if work == nil {
		work = make(EntityWork)
	}

	//var traversed tosca.EntityPtrs

	// reflection.EntityTraverser signature
	traverseWrapper := func(entityPtr tosca.EntityPtr) bool {
		if work.Start(log, entityPtr) {
			return false
		}

		// Don't traverse the same entity more than once
		/*for _, entityPtr_ := range traversed {
			if entityPtr_ == entityPtr {
				return false
			}
		}
		traversed = append(traversed, entityPtr)*/

		return traverse(entityPtr)
	}

	// Root
	reflection.TraverseEntities(self.Root.EntityPtr, traverseWrapper)

	// Types
	self.Root.GetContext().Namespace.Range(func(entityPtr tosca.EntityPtr) bool {
		reflection.TraverseEntities(entityPtr, traverseWrapper)
		return true
	})
}

//
// EntityWork
//

type EntityWork map[tosca.EntityPtr]struct{}

func (self EntityWork) Start(log logging.Logger, entityPtr tosca.EntityPtr) bool {
	if _, ok := self[entityPtr]; ok {
		log.Debugf("skip: %s", tosca.GetContext(entityPtr).Path)
		return true
	}
	self[entityPtr] = struct{}{}
	return false
}

//
// CoordinatedWork
//

type CoordinatedWork struct {
	sync.Map
	Log logging.Logger
}

func NewCoordinatedWork(log logging.Logger) *CoordinatedWork {
	return &CoordinatedWork{
		Log: log,
	}
}

func (self *CoordinatedWork) Start(key string) (Promise, bool) {
	promise := NewPromise()
	if existing, loaded := self.LoadOrStore(key, promise); !loaded {
		self.Log.Debugf("start: %s", key)
		return promise, true
	} else {
		self.Log.Debugf("wait for: %s", key)
		promise = existing.(Promise)
		promise.Wait()
		return nil, false
	}
}

//
// Promise
//

type Promise chan struct{}

func NewPromise() Promise {
	return make(Promise)
}

func (self Promise) Release() {
	close(self)
}

func (self Promise) Wait() {
	<-self
}
