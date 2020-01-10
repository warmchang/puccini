package parser

import (
	"fmt"
	"sort"
	"sync"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/csar"
	"github.com/tliron/puccini/tosca/grammars"
	"github.com/tliron/puccini/tosca/reflection"
	"github.com/tliron/puccini/url"
)

func (self *Context) ReadRoot(url_ url.URL) bool {
	toscaContext := tosca.NewContext(self.Quirks)

	toscaContext.URL = url_

	var ok bool

	self.WG.Add(1)
	self.Root, ok = self.read(nil, toscaContext, nil, nil, "$Root")
	self.WG.Wait()

	sort.Sort(self.Units)

	return ok
}

var readCache sync.Map // entityPtr or Promise

func (self *Context) read(promise Promise, toscaContext *tosca.Context, container *Unit, nameTransfomer tosca.NameTransformer, readerName string) (*Unit, bool) {
	defer self.WG.Done()
	if promise != nil {
		// For the goroutines waiting for our cached entityPtr
		defer promise.Release()
	}

	log.Infof("{read} %s: %s", readerName, toscaContext.URL.Key())

	switch toscaContext.URL.Format() {
	case "csar", "zip":
		var err error
		if toscaContext.URL, err = csar.GetRootURL(toscaContext.URL); err != nil {
			toscaContext.ReportError(err)
			return nil, false
		}
	}

	// Read ARD
	var err error
	if toscaContext.Data, toscaContext.Locator, err = ard.ReadURL(toscaContext.URL, true); err != nil {
		toscaContext.ReportError(err)
		return nil, false
	}

	// Detect grammar
	if !grammars.Detect(toscaContext) {
		return nil, false
	}

	// Read entityPtr
	read, ok := toscaContext.Grammar.Readers[readerName]
	if !ok {
		panic(fmt.Sprintf("grammar does not support reader \"%s\"", readerName))
	}
	entityPtr := read(toscaContext)
	if entityPtr == nil {
		// Even if there are problems, the reader should return an entityPtr
		panic(fmt.Sprintf("reader \"%s\" returned a non-entity: %T", reflection.GetFunctionName(read), entityPtr))
	}

	// Validate required fields
	reflection.Traverse(entityPtr, tosca.ValidateRequiredFields)

	readCache.Store(toscaContext.URL.Key(), entityPtr)

	return self.AddUnit(entityPtr, container, nameTransfomer), true
}

// From Importer interface
func (self *Context) goReadImports(container *Unit) {
	var importSpecs []*tosca.ImportSpec
	if importer, ok := container.EntityPtr.(tosca.Importer); ok {
		importSpecs = importer.GetImportSpecs()
	}

	// Implicit import
	if implicitImportSpec, ok := grammars.GetImplicitImportSpec(container.GetContext()); ok {
		importSpecs = append(importSpecs, implicitImportSpec)
	}

	for _, importSpec := range importSpecs {
		key := importSpec.URL.Key()

		// Skip if causes import loop
		skip := false
		for container_ := container; container_ != nil; container_ = container_.Container {
			url_ := container_.GetContext().URL
			if url_.Key() == key {
				if !importSpec.Implicit {
					// Import loops are considered errors
					container.GetContext().ReportImportLoop(url_)
				}
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		promise := NewPromise()
		if cached, inCache := readCache.LoadOrStore(key, promise); inCache {
			switch cached_ := cached.(type) {
			case Promise:
				// Wait for promise
				log.Debugf("{read} wait for promise: %s", key)
				self.WG.Add(1)
				go self.waitForPromise(cached_, key, container, importSpec.NameTransformer)

			default: // entityPtr
				// Cache hit
				log.Debugf("{read} cache hit: %s", key)
				self.AddUnit(cached, container, importSpec.NameTransformer)
			}
		} else {
			importToscaContext := container.GetContext().NewImportContext(importSpec.URL)

			// Read (concurrently)
			self.WG.Add(1)
			go self.read(promise, importToscaContext, container, importSpec.NameTransformer, "$Unit")
		}
	}
}

func (self *Context) waitForPromise(promise Promise, key string, container *Unit, nameTransformer tosca.NameTransformer) {
	defer self.WG.Done()
	promise.Wait()

	if cached, inCache := readCache.Load(key); inCache {
		switch cached.(type) {
		case Promise:
			log.Debugf("{read} promise broken: %s", key)

		default: // entityPtr
			// Cache hit
			log.Debugf("{read} promise kept: %s", key)
			self.AddUnit(cached, container, nameTransformer)
		}
	} else {
		log.Debugf("{read} promise broken (empty): %s", key)
	}
}