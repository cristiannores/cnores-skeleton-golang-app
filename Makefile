# Configuración básica
APP_NAME=cnores-skeleton-golang-app
DOCKER_TAG=latest
CONFIG_FILE=app/shared/utils/config/config.json
PORT=$(shell cat $(CONFIG_FILE) | grep '"port":' | sed 's/.*"port": "\([^"]*\)".*/\1/')

# Construir la imagen Docker
build:
	docker build -t $(APP_NAME):$(DOCKER_TAG) .

# Ejecutar el contenedor
run:
	docker run -p $(PORT):$(PORT) --name $(APP_NAME) $(APP_NAME):$(DOCKER_TAG)

# Detener y remover el contenedor
stop:
	docker stop $(APP_NAME)
	docker rm $(APP_NAME)

# Limpiar imágenes no utilizadas
clean:
	docker system prune -af

# Acceso al contenedor para debugging
shell:
	docker exec -it $(APP_NAME) sh
