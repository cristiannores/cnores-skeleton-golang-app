package services

import (
	"medical-bff/app/infrastructure/setup_dependencies/dependencies"
)

type ServiceInitializer func(container *dependencies.DependencyContainer) interface{}

var initializers = make(map[string]ServiceInitializer)
var serviceInstances = make(map[string]interface{})

func Register(name string, initializer ServiceInitializer) {
	initializers[name] = initializer
}

func InitializeAll(container *dependencies.DependencyContainer) {
	for name, initializer := range initializers {

		serviceInstances[name] = initializer(container)
	}
}

func GetService(name string) interface{} {
	return serviceInstances[name]
}
