package util

import (
	"time"

	"github.com/spf13/viper"
)

var (
	DevEnv  = "dev"
	ProdEnv = "prod"
)

type Config struct {
	Environment        string   `mapstructure:"ENVIRONMENT"`
	Domain             string   `mapstructure:"DOMAIN"`
	CorsAllowedOrigins []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	// @deprecated
	// DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	MigrationUrl        string        `mapstructure:"MIGRATION_URL"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	HTTPServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress   string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymetricKey    string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefresTokenDuration time.Duration `mapstructure:"REFRES_TOKEN_DURATION"`
	RedisAddress        string        `mapstructure:"REDIS_ADDRESS"`

	EmailSenderName    string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	// EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
	SmtpPassword   string `mapstructure:"SMTP_PASSWORD"`
	SmtpUser       string `mapstructure:"SMTP_USER"`
	SmtpHost       string `mapstructure:"SMTP_HOST"`
	SmtpPort       int    `mapstructure:"SMTP_PORT"`
	SmtpEncryption string `mapstructure:"SMTP_ENCRYPTION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") //json, xml, yaml

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func (c Config) GetSmtpHost() string {
	return c.SmtpHost
}
func (c Config) GetSmtpPort() int {
	return c.SmtpPort
}
func (c Config) GetSmtpUser() string {
	return c.SmtpUser
}
func (c Config) GetSmtpPassword() string {
	return c.SmtpPassword
}
func (c Config) GetSmtpEncryption() string {
	return c.SmtpEncryption
}
func (c Config) GetFromName() string {
	return c.EmailSenderName
}
func (c Config) GetFromEmailAddress() string {
	return c.EmailSenderAddress
}

// func LoadConfigAndWatcing(path string, fn func(config Config)) (err error) {
//     var config Config

// 	viper.AddConfigPath(path)
// 	viper.SetConfigName("app")
// 	viper.SetConfigType("env") //json, xml, yaml

// 	viper.AutomaticEnv()

// 	err = viper.ReadInConfig()
// 	if err != nil {
// 		return
// 	}

// 	err = viper.Unmarshal(&config)

//     viper.WatchConfig()
//     viper.OnConfigChange(func(e fsnotify.Event) {
//         var newConfig Config
//         if err := viper.Unmarshal(&newConfig); err != nil {
//             log.Printf("Error unmarshaling config: %s", err)
//             return
//         }
//         fn(newConfig)
//         // config = newConfig
//     })

//     fn(config)

// 	return
// }
