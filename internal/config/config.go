package config

import (
	"fmt"
	"os"
	"reflect"
)

// For example Email-Provider and Forget Password or Verify Mail

// This struct is used to store the configuration of the application
type Config struct {
	// App Settings
	App App
	// Server Settings
	Server Server
	// Database Settings
	Database Database
	// Theme Settings
	Theme Theme
	// Auth Settings
	Auth Auth
	// Mail Settings
	Mail Mail
}

type App struct {
	// App Name. Default "Go-Template"
	Name string
	// App URL. Default "http://localhost:8080"
	Url string
	// Shopify App. Default false
	ShopifyApp bool
}

type Server struct {
	// Server Port. Default 8080
	Port int
}

type DatabaseType string

const (
	DatabaseTypeSQLite   DatabaseType = "sqlite"
	DatabaseTypePostgres DatabaseType = "postgres"
)

type Database struct {
	Enabled bool
	// Database Type. Default "sqlite"
	Type DatabaseType
}

type Mail struct {
	// Enable Mail. Default true
	EnableMail bool
	// Mail Provider. Default MailProviderResend
	MailProvider MailProvider
}

type MailProvider string

const (
	// Resend
	MailProviderResend MailProvider = "resend"
	// Print Mails to Console for Debug / Dev
	MailProviderConsole MailProvider = "console"
	// Sent Mails via Imap NOT IMPLEMENTED YET
	//MailProviderIMAP   MailProvider = "imap"
)

type Theme struct {
	// Set Standard Theme. Default ""
	// We use DaisyUI. If you want to add more themes you can do this in tailwind.config.js
	StandardTheme string
	// Enable Theme Switcher. Default true
	EnableThemeSwitcher bool
	// Enable Sidebar. Default true
	EnableSidebar bool
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
	if !c.Database.Enabled {
		c.Auth.EnableAuth = false
	}
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

// New returns the config with the default values
func New(overrides *Config) *Config {
	config := &Config{
		Server: Server{
			Port: 8080, // Default port
		},
		App: App{
			Name:       os.Getenv("APP_NAME"),
			Url:        os.Getenv("APP_URL"),
			ShopifyApp: false,
		},
		Database: Database{
			Enabled: true,
			Type:    DatabaseTypeSQLite,
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
			EnableMail:   true,               // Default to true
			MailProvider: MailProviderResend, // Default to MailProviderResend
		},
	}

	if overrides != nil {
		mergeConfig(config, overrides)
	}

	config.CheckEnvironmentVariables()
	config.validateDependencies()

	return config
}

// Checks if specific environment variables are set
func (c *Config) CheckEnvironmentVariables() error {
	if c.Database.Type == DatabaseTypePostgres {
		if os.Getenv("DB_HOST") == "" || os.Getenv("DB_DATABASE") == "" || os.Getenv("DB_PORT") == "" || os.Getenv("DB_USER") == "" || os.Getenv("DB_PASSWORD") == "" {
			fmt.Println("Warning: DB_HOST, DB_DATABASE, DB_PORT, DB_USER or DB_PASSWORD environment variable is not set")
			c.Database.Enabled = false
			fmt.Println("Database functionality has been disabled")
			return nil
		}
	}
	if c.Database.Type == DatabaseTypeSQLite {
		if os.Getenv("DB_FILE") == "" {
			fmt.Println("Warning: DB_FILE environment variable is not set")
			fmt.Println("Defaulting to 'db/dev.sqlite'")
			c.Database.Type = DatabaseTypeSQLite
			return nil
		}
	}
	if c.Mail.MailProvider == MailProviderResend {
		if os.Getenv("RESEND_API_KEY") == "" || os.Getenv("RESEND_FROM_EMAIL") == "" || os.Getenv("RESEND_FROM_NAME") == "" {
			fmt.Println("Warning: RESEND_API_KEY, RESEND_FROM_EMAIL or RESEND_FROM_NAME environment variable is not set")
			c.Mail.EnableMail = false
			fmt.Println("Mail functionality has been disabled")
			return nil
		}
	}
	// Check for essential environment variables
	essentialEnvVars := []string{"APP_NAME", "APP_URL", "SECRET_KEY"}
	for _, envVar := range essentialEnvVars {
		if os.Getenv(envVar) == "" {
			fmt.Printf("Error: %s environment variable is not set\n", envVar)
			os.Exit(1)
		}
	}
	return nil
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
