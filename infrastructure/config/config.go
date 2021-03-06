//go:generate mockgen -source config.go -destination config_mock.go -package config

package config

// ConfigLoader is the default configuration loader.
var ConfigLoader = &viperLoader{}

// SecretsLoader is the default secrets loader.
var SecretsLoader = &keyringLoader{}

// Config is the interface for reading and writing configuration.
type Config interface {
	IsSet(key string) bool
	GetString(key string) string
	Set(key string, value interface{})
	Save() error
}

// Loader loads Config by name.
type Loader interface {
	Load(name string) (Config, error)
}
