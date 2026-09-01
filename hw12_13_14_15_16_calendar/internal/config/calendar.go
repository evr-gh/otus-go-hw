package config

type CalendarConfig struct {
	HTTP    HTTPConfig
	RPC     RPCConfig
	Storage StorageConfig
	Logger  LoggerConfig
}

func NewCalendarConfig() *CalendarConfig {
	return &CalendarConfig{}
}
