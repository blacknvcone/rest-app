package config

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type (
	HuggingFaceAPIConf struct {
		URL      string
		APIToken string
		Model    string
	}

	GoogleAIAPIConf struct {
		URL      string
		APIToken string
		Model    string
	}

	DB struct {
		DSN             string
		DSNPool         string
		MaxOpenConn     int
		MaxIdleConn     int
		MaxLifetimeConn int
		MaxIdletimeConn int
	}

	app struct {
		Env     string
		Version string
		Name    string
	}

	http struct {
		Port int
	}

	jwt struct {
		SigningKey string
	}

	Config struct {
		DB                 DB
		App                app
		Http               http
		JWT                jwt
		HuggingFaceAPIConf HuggingFaceAPIConf
		GoogleAIAPIConf    GoogleAIAPIConf
	}
)

var (
	configData *Config
)

func InitConfig() {
	viper.SetConfigType("env")
	viper.SetConfigName(".env") // name of Config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/secrets")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Warn("failed to load config file")
	}

	configData = &Config{
		DB: DB{
			DSN:             getRequiredString("DB_DSN"),
			DSNPool:         getRequiredString("DB_POOL_DSN"),
			MaxOpenConn:     getRequiredInt("DB_MAX_OPEN_CONN"),
			MaxIdleConn:     getRequiredInt("DB_MAX_IDLE_CONN"),
			MaxLifetimeConn: getRequiredInt("DB_MAX_LIFETIME_CONN"),
			MaxIdletimeConn: getRequiredInt("DB_MAX_IDLETIME_CONN"),
		},
		App: app{
			Env:     getRequiredString("APP_ENV"),
			Version: viper.GetString("BITBUCKET_TAG"),
			Name:    "rest-app",
		},
		Http: http{
			Port: getRequiredInt("APP_PORT"),
		},
		JWT: jwt{
			SigningKey: getRequiredString("SIGNING_KEY"),
		},
		// HuggingFaceAPIConf: HuggingFaceAPIConf{
		// 	URL:      getRequiredString("HUGGINGFACE_API_URL"),
		// 	Model:    getRequiredString("HUGGINGFACE_API_MODEL"),
		// 	APIToken: getRequiredString("HUGGINGFACE_API_TOKEN"),
		// },
		GoogleAIAPIConf: GoogleAIAPIConf{
			URL:      getRequiredString("GOOGLE_AI_API_URL"),
			Model:    getRequiredString("GOOGLE_AI_API_MODEL"),
			APIToken: getRequiredString("GOOGLE_AI_API_TOKEN"),
		},
	}
}

func getRequiredString(key string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	log.Fatalln(fmt.Errorf("KEY %s IS MISSING", key))
	return ""
}

func getRequiredInt(key string) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}

	panic(fmt.Errorf("KEY %s IS MISSING", key))
}

// func getRequiredBool(key string) bool {
// 	if viper.IsSet(key) {
// 		return viper.GetBool(key)
// 	}

// 	panic(fmt.Errorf("KEY %s IS MISSING", key))
// }

// func getRequiredDuration(key string) time.Duration {
// 	if viper.IsSet(key) {
// 		return viper.GetDuration(key)
// 	}

// 	panic(fmt.Errorf("KEY %s IS MISSING", key))
// }

func GetConfig() Config {
	return *configData
}
