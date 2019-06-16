package modules

import (
	"sync"

	"github.com/pajbot/pajbot2/pkg"
)

var moduleSpecsMutex sync.Mutex
var moduleSpecs []pkg.ModuleSpec
var moduleSpecsMap map[string]pkg.ModuleSpec

var moduleFactoriesMutex sync.Mutex
var moduleFactories map[string]pkg.ModuleFactory

// Register registers a module by taking its ID and factory (instructions how to create the module)
// Can be called both from within this package, or by an external package. Generally done in any files init function
func Register(moduleID string, factory pkg.ModuleFactory) {
	if factory == nil {
		panic("Trying to register a nil factory")
	}

	moduleFactoriesMutex.Lock()
	if moduleFactories == nil {
		moduleFactories = make(map[string]pkg.ModuleFactory)
	}

	if _, ok := moduleFactories[moduleID]; ok {
		panic("A module factory with the id '" + moduleID + "' was already registered")
	}

	moduleFactories[moduleID] = factory
	moduleFactoriesMutex.Unlock()
}

// List returns a list of module specs
func List() []pkg.ModuleSpec {
	moduleSpecsMutex.Lock()
	defer moduleSpecsMutex.Unlock()

	if moduleSpecs == nil {
		moduleSpecs = generateSpecs()
		initModuleMap(moduleSpecs)
	}

	return moduleSpecs
}

// Map returns a map of module specs, keyed with the modules ID
func Map() map[string]pkg.ModuleSpec {
	moduleSpecsMutex.Lock()
	defer moduleSpecsMutex.Unlock()

	if moduleSpecs == nil {
		moduleSpecs = generateSpecs()
		initModuleMap(moduleSpecs)
	}

	return moduleSpecsMap
}

// GetModuleSpec returns the module spec of the module with the given ID
func GetModuleSpec(moduleID string) (spec pkg.ModuleSpec, ok bool) {
	moduleSpecsMutex.Lock()
	defer moduleSpecsMutex.Unlock()

	if moduleSpecs == nil {
		moduleSpecs = generateSpecs()
		initModuleMap(moduleSpecs)
	}

	spec, ok = moduleSpecsMap[moduleID]
	return
}

// GetModuleFactory returns the module factory with the given ID.
// This is useful for tests where the module spec should not be shared
func GetModuleFactory(moduleID string) (factory pkg.ModuleFactory, ok bool) {
	moduleFactoriesMutex.Lock()
	defer moduleFactoriesMutex.Unlock()

	factory, ok = moduleFactories[moduleID]
	return
}

func generateSpecs() (specs []pkg.ModuleSpec) {
	for _, moduleFactory := range moduleFactories {
		spec := moduleFactory()
		specs = append(specs, spec)
	}

	return
}

func initModuleMap(specs []pkg.ModuleSpec) {
	moduleSpecsMap = make(map[string]pkg.ModuleSpec)

	for _, spec := range specs {
		moduleSpecsMap[spec.ID()] = spec
	}
}
