package config

import (
	"reflect"
)

// TODO: Dependencies?
// For example Email-Provider and Forget Password or Verify Mail

// This struct is used to store the configuration of the application
type Config struct {
	// Server Settings
	Server Server
	// Theme Settings
	Theme Theme
	// Auth Settings
	Auth Auth
	// Mail Settings
	Mail Mail
}

// DATABASE -> Config in database.go

type Server struct {
	// Server Port. Default 8080
	Port int
}

type Theme struct {
	// Set Standard Theme. Default ""
	// We use DaisyUI. If you want to add more themes you can do this in tailwind.config.js
	StandardTheme string
	// Enable Theme Switcher. Default true
	EnableThemeSwitcher bool
	// Enable Sidebar. Default true
	EnableSidebar bool
}

type Mail struct {
	// Enable Mail. Default true
	EnableMail bool
}

type Auth struct {
	// Enable Authentication. Default true. Disables the complete user authentification and all routes that need authentication. Removes User menu from the header
	// Default true
	EnableAuth bool
	// Enable Registration. Default true. Already registered users are able to access the site. Only Routes are disabled
	// Default True
	EnableRegistration bool
	// Prevents Login and Logout. Attention: Loggedin are still able to access the website. Only Routes are disabled
	// Default true
	EnableLogin bool
	// Default to true
	// Disable Avatars if you cannot store the images on the server or you don't want to
	EnableAvatar bool
	// Enable Reset Password. Default true
	EnableResetPassword bool
	// Enable Verify Email. Default true
	EnableVerifyEmail bool
}

// validateDependencies checks and adjusts dependent settings
func (c *Config) validateDependencies() {
	// If mail is disabled
	if !c.Mail.EnableMail {
		c.Auth.EnableResetPassword = false
		c.Auth.EnableVerifyEmail = false
	}

	// If registration is disabled
	if !c.Auth.EnableAuth {
		c.Auth.EnableLogin = false
		c.Auth.EnableRegistration = false
		c.Auth.EnableResetPassword = false
		c.Auth.EnableVerifyEmail = false
	}
}

// This function merges the base config with the overrides config set in the server.go
func mergeConfig(base, overrides interface{}) {
	baseVal := reflect.ValueOf(base).Elem()
	overridesVal := reflect.ValueOf(overrides).Elem()

	for i := 0; i < baseVal.NumField(); i++ {
		baseField := baseVal.Field(i)
		overridesField := overridesVal.Field(i)

		if baseField.Kind() == reflect.Struct {
			mergeConfig(baseField.Addr().Interface(), overridesField.Addr().Interface())
		} else {
			// For boolean fields, we want to set them even if they're false
			if overridesField.Kind() == reflect.Bool || !overridesField.IsZero() {
				baseField.Set(overridesField)
			}
		}
	}
}

// New returns the config with the default values
func New(overrides *Config) *Config {
	config := &Config{
		Server: Server{
			Port: 8080, // Default port
		},
		Theme: Theme{
			StandardTheme:       "",
			EnableThemeSwitcher: true,
			EnableSidebar:       true,
		},
		Auth: Auth{
			EnableAuth:          true, // Default to true
			EnableRegistration:  true, // Default to true
			EnableLogin:         true, // Default to true
			EnableAvatar:        true, // Default to true
			EnableResetPassword: true, // Default to true
			EnableVerifyEmail:   true, // Default to true
		},
		Mail: Mail{
			EnableMail: true, // Default to true
		},
	}

	if overrides != nil {
		mergeConfig(config, overrides)
	}

	config.validateDependencies()

	return config
}
