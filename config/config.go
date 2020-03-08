package config

// define the version
var Version string

// Config defines the configuration structure.
type Config struct {
	General struct {
		LogLevel               int    `mapstructure:"log_level"`
		PasswordHashIterations int    `mapstructure:"password_hash_iterations"`
		HostDomain             string `mapstructure:"host_domain"`
	} `mapstructure:"general"`

	PostgreSQL struct {
		DSN                string `mapstructure:"dsn"`
		MaxOpenConnections int `mapstructure:"max_open_connections"`
		MaxIdleConnections int `mapstructure:"max_idle_connections"`
	} `mapstructure:"postgresql"`

	DemoServer struct {
		API struct {
			Bind       string
			CACert     string `mapstructure:"ca_cert"`
			TLSCert    string `mapstructure:"tls_cert"`
			TLSKey     string `mapstructure:"tls_key"`
			PublicHost string `mapstructure:"public_host"`
		} `mapstructure:"api"`

		ExternalAPI struct {
			Bind                       string
			TLSCert                    string `mapstructure:"tls_cert"`
			TLSKey                     string `mapstructure:"tls_key"`
			JWTSecret                  string `mapstructure:"jwt_secret"`
			DisableAssignExistingUsers bool   `mapstructure:"disable_assign_existing_users"`
			CORSAllowOrigin            string `mapstructure:"cors_allow_origin"`
		} `mapstructure:"external_api"`

	} `mapstructure:"demo_server"`
}

// Conf holds the global configuration.
var Conf Config