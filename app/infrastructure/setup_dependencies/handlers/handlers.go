package handlers

import (
	"cnores-skeleton-golang-app/app/infrastructure/setup_dependencies/dependencies"
)

type HandlerInitializer func(container *dependencies.DependencyContainer)

var handlerInitializers []HandlerInitializer

func RegisterHandler(initializer HandlerInitializer) {
	handlerInitializers = append(handlerInitializers, initializer)
}

func InitializeAllHandlers(container *dependencies.DependencyContainer) {
	for _, initializer := range handlerInitializers {
		initializer(container)
	}
}
