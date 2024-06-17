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
	TSL  bool   `json:"TSL"`
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
