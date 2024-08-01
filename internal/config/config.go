package config

// This struct is used to store the configuration of the application
type Config struct {
	Port               int
	EnableAuth         bool
	EnableRegistration bool
	EnableLogin        bool
	EnableAvatar       bool
}

var instance *Config

func New(overrides *Config) *Config {
	if instance == nil {
		instance = &Config{
			Port:               8080, // Default port
			EnableAuth:         true, // Default to true
			EnableRegistration: true, // Default to true
			EnableLogin:        true, // Default to true
			EnableAvatar:       true, // Default to true
		}
	}

	if overrides != nil {
		if overrides.Port != 0 {
			instance.Port = overrides.Port
		}
		instance.EnableAuth = overrides.EnableAuth
		instance.EnableRegistration = overrides.EnableRegistration
		instance.EnableLogin = overrides.EnableLogin
		instance.EnableAvatar = overrides.EnableAvatar
	}

	return instance
}
