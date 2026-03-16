package config

type ViperConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

const (
	SERVER_MODE_DEV  = "dev"
	SERVER_MODE_PROD = "prod"
)

type ServerConfig struct {
	Mode     string      `mapstructure:"mode"`
	HTTPPort uint        `mapstructure:"http_port"`
	TCPPort  uint        `mapstructure:"tcp_port"`
	Log      string      `mapstructure:"log"`
	Jwt      string      `mapstructure:"jwt"`
	AuthCode string      `mapstructure:"auth_code"`
	Email    EmailConfig `mapstructure:"email"`
}

type EmailConfig struct {
	Template string `mapstructure:"template"`
	Host     string `mapstructure:"host"`
	Port     uint   `mapstructure:"port"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

const (
	DATABASE_DRIVER_SQLITE = "sqlite"
	DATABASE_DRIVER_MYSQL  = "mysql"
)

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}
