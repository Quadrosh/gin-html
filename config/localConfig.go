package config

// LocalConfig локальный конфиг
type LocalConfig struct {
	InProduction bool
	UseCache     bool
	SSL          bool
	SSLCert      string
	SSLKey       string
	DbHost       string
	DbName       string
	DbUser       string
	DbPass       string
	DbPort       string
	DbSSL        string
	AppPort      string
}
