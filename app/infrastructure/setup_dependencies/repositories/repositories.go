package repositories

import (
	"cnores-skeleton-golang-app/app/infrastructure/setup_dependencies/dependencies"
)

type RepositoryInitializer func(container *dependencies.DependencyContainer) interface{}

var initializers = make(map[string]RepositoryInitializer)
var repositoryInstances = make(map[string]interface{})

func Register(name string, initializer RepositoryInitializer) {
	initializers[name] = initializer
}

func InitializeAll(container *dependencies.DependencyContainer) {
	for name, initializer := range initializers {

		repositoryInstances[name] = initializer(container)
	}
}

func GetRepository(name string) interface{} {
	return repositoryInstances[name]
}
