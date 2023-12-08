

# Cnores Skeleton Golang App

## Arquitectura Limpia (Clean Architecture)

La aplicación "Cnores Skeleton Golang App" implementa la Clean Architecture para proporcionar una estructura de código bien organizada, fácil de mantener y escalable. Esta arquitectura se centra en la separación de responsabilidades y en la desacoplación de los diferentes componentes de la aplicación.

### Componentes Principales

#### 1. Setup de Dependencias (`setup_dependencies`)

- Ubicación: `app/infrastructure/setup_dependencies/dependencies.go`
- Función: Inicializa todas las dependencias necesarias para la aplicación. Crea un contenedor de dependencias y registra todos los repositorios y handlers necesarios.

#### 2. Repositorios (`repositories`)

- Ubicación: `app/infrastructure/setup_dependencies/repositories/repositories.go`
- Función: Proporciona una abstracción sobre la capa de acceso a datos. Cada repositorio encapsula la lógica para interactuar con una fuente de datos específica (por ejemplo, base de datos, API externa).

#### 3. Handlers (`handlers`)

- Ubicación: `app/infrastructure/setup_dependencies/handlers/handlers.go`
- Función: Gestiona las solicitudes entrantes, invoca la lógica del negocio a través de los controladores y casos de uso, y envía las respuestas adecuadas. Cada handler define cómo se procesa una solicitud específica y cómo se formatea la respuesta.

#### 4. Servicios (`services`)

- Ubicación: `app/infrastructure/setup_dependencies/services/services.go`
- Función: Contiene la lógica de negocio que es específica para un servicio o funcionalidad en particular. Los servicios interactúan con múltiples repositorios y otros servicios según sea necesario.

### Flujo de Trabajo

1. **Inicialización**: Al iniciar la aplicación, se configuran todas las dependencias. Esto incluye la creación de instancias de repositorios, servicios y handlers.

2. **Handlers**: Los handlers, que están registrados durante la inicialización, escuchan las solicitudes entrantes y las delegan a los controladores correspondientes.

3. **Controladores**: Los controladores reciben la solicitud del handler, ejecutan la lógica de negocio a través de los casos de uso y devuelven la respuesta al handler.

4. **Casos de Uso**: Cada caso de uso encapsula una operación de negocio específica, trabajando con los servicios y repositorios necesarios para realizar su función.

5. **Repositorios**: Los repositorios abstraen la capa de acceso a datos, permitiendo a los casos de uso interactuar con bases de datos o servicios externos sin preocuparse por los detalles de implementación.

### Ejemplos de Implementación

- **Handler**: `find_by_id_billing_handler` maneja las solicitudes específicas para encontrar facturas por ID.
- **Controlador**: `find_by_id_billing_controller` procesa la solicitud de encontrar una factura y utiliza el caso de uso correspondiente.
- **Caso de Uso**: `find_by_id_billing_usecase` contiene la lógica para buscar una factura en la base de datos a través del repositorio.

Esta estructura garantiza que la lógica de negocio esté bien organizada, sea fácil de mantener y probar, y que la aplicación sea escalable y adaptable a cambios futuros.



## Configuración

El archivo de configuración se encuentra en `app/shared/utils/config/config.json`. La estructura del archivo de configuración en Go es la siguiente:

```go
package config

type Config struct {
    CurrentStage string             `json:"currentStage" validate:"required"`
    Developers   []string           `json:"developers" validate:"required"`
    Url          string             `json:"url" validate:"required"`
    Port         string             `json:"port" validate:"required"`
    Services     map[string]Service `json:"services" validate:"required"`
    Notifier     Notifier           `json:"notifier" validate:"required"`
    Database     DatabaseSettings   `json:"database" validate:"required"`
}

type DatabaseSettings struct {
    Name string `json:"name" validate:"required"`
    Url  string `json:"url" validate:"required"`
}

type Notifier struct {
    Url         string            `json:"url" validate:"required"`
    Timeout     int               `json:"timeout" validate:"gte=0"`
    Headers     map[string]string `json:"headers" validate:"required"`
    Token       string            `json:"token" validate:"required"`
    EnableFaker bool              `json:"enable-faker"`
}

type Service struct {
    Url         string            `json:"url" validate:"required"`
    EnableFaker bool              `json:"enable-faker"`
    TimeOut     int               `json:"timeout" validate:"gte=0"`
    Headers     map[string]string `json:"headers" validate:"required"`
}
```

## Uso

Para ejecutar la aplicación, puedes utilizar el Makefile incluido en el proyecto o directamente a través de Docker.

### Uso del Makefile

El Makefile proporciona comandos convenientes para construir y ejecutar la aplicación.

#### Construir la Imagen Docker

```bash 
make build
```
#### Ejecutar la Aplicación

```bash 
make run
```
Esto iniciará la aplicación y la expondrá en el puerto especificado en config.json.

#### Detener la Aplicación

```bash 
make stop
``` 

#### Limpiar Imágenes Docker

```bash 
make clean
``` 

### Uso Directo de Docker

Para construir la imagen Docker directamente:

```bash 
docker build -t cnores-skeleton-golang-app .
``` 
Para ejecutar la aplicación:

```bash 
docker run -p <PORT>:<PORT> cnores-skeleton-golang-app
``` 

Reemplaza <PORT> con el puerto especificado en tu config.json.

